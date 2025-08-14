package github

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/bgentry/go-netrc/netrc"
	"github.com/google/go-github/v66/github"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type IssueWithStatus struct {
	ProjectIssue
	GitHubStatus  string
	ProjectItemID string
	UpdatedAt     time.Time
}

type Client struct {
	restClient   *github.Client
	graphqlClient *githubv4.Client
	ctx          context.Context
	org          string
	currentUser  string
}

type ProjectIssue struct {
	ID        int
	Title     string
	Body      string
	Number    int
	State     string
	URL       string
	Labels    []string
	UpdatedAt time.Time
}

func NewClient(org string) (*Client, error) {
	token, err := getGitHubToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get GitHub token: %w", err)
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	restClient := github.NewClient(tc)
	graphqlClient := githubv4.NewClient(tc)

	user, _, err := restClient.Users.Get(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get current user: %w", err)
	}
	

	return &Client{
		restClient:    restClient,
		graphqlClient: graphqlClient,
		ctx:           ctx,
		org:           org,
		currentUser:   user.GetLogin(),
	}, nil
}

func getGitHubToken() (string, error) {
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		return token, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get user home directory: %w", err)
	}

	netrcPath := filepath.Join(homeDir, ".netrc")
	if _, err := os.Stat(netrcPath); os.IsNotExist(err) {
		return "", fmt.Errorf("no GitHub token found: set GITHUB_TOKEN environment variable or add github.com entry to ~/.netrc")
	}

	n, err := netrc.ParseFile(netrcPath)
	if err != nil {
		return "", fmt.Errorf("failed to parse .netrc file: %w", err)
	}

	machine := n.FindMachine("github.com")
	if machine == nil {
		return "", fmt.Errorf("no github.com entry found in .netrc file")
	}

	if machine.Password == "" {
		return "", fmt.Errorf("no password found for github.com in .netrc file")
	}

	return machine.Password, nil
}

func (c *Client) GetProjectIssues(projectNumber int, status string) ([]ProjectIssue, error) {
	var allIssues []ProjectIssue
	
	opts := &github.ListOptions{
		PerPage: 100,
	}

	for {
		items, resp, err := c.restClient.Projects.ListProjectCards(c.ctx, int64(projectNumber), &github.ProjectCardListOptions{
			ListOptions: *opts,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list project cards: %w", err)
		}

		for _, card := range items {
			if card.GetContentURL() == "" {
				continue
			}

			issue, err := c.parseIssueFromCard(card)
			if err != nil {
				continue
			}

			if strings.Contains(strings.ToLower(card.GetColumnName()), strings.ToLower(status)) {
				allIssues = append(allIssues, issue)
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allIssues, nil
}

func (c *Client) parseIssueFromCard(card *github.ProjectCard) (ProjectIssue, error) {
	contentURL := card.GetContentURL()
	parts := strings.Split(contentURL, "/")
	if len(parts) < 2 {
		return ProjectIssue{}, fmt.Errorf("invalid content URL")
	}

	issueNumberStr := parts[len(parts)-1]
	issueNumber, err := strconv.Atoi(issueNumberStr)
	if err != nil {
		return ProjectIssue{}, fmt.Errorf("invalid issue number: %w", err)
	}

	repoName := parts[len(parts)-3]
	issue, _, err := c.restClient.Issues.Get(c.ctx, c.org, repoName, issueNumber)
	if err != nil {
		return ProjectIssue{}, fmt.Errorf("failed to get issue: %w", err)
	}

	var labels []string
	for _, label := range issue.Labels {
		labels = append(labels, label.GetName())
	}

	return ProjectIssue{
		ID:        int(issue.GetID()),
		Title:     issue.GetTitle(),
		Body:      issue.GetBody(),
		Number:    issue.GetNumber(),
		State:     issue.GetState(),
		URL:       issue.GetHTMLURL(),
		Labels:    labels,
		UpdatedAt: issue.GetUpdatedAt().Time,
	}, nil
}

func (c *Client) GetProjectV2Issues(projectNumber int, statusFilters []string) ([]IssueWithStatus, error) {
	// Step 1: Get all issues assigned to current user (much smaller dataset)
	fmt.Printf("Fetching issues assigned to %s...\n", c.currentUser)
	
	assignedIssues, err := c.getAssignedIssues()
	if err != nil {
		return nil, fmt.Errorf("failed to get assigned issues: %w", err)
	}
	
	fmt.Printf("Found %d issues assigned to %s\n", len(assignedIssues), c.currentUser)
	
	if len(assignedIssues) == 0 {
		return []IssueWithStatus{}, nil
	}
	
	// Step 2: Get project status for assigned issues only
	return c.getProjectStatusForIssues(projectNumber, assignedIssues, statusFilters)
}

func (c *Client) getAssignedIssues() ([]ProjectIssue, error) {
	var allIssues []ProjectIssue
	
	// Get both open and closed issues assigned to user
	for _, state := range []string{"open", "closed"} {
		opts := &github.IssueListOptions{
			Filter:    "assigned",
			State:     state,
			ListOptions: github.ListOptions{PerPage: 100},
		}

		for {
			issues, resp, err := c.restClient.Issues.ListByOrg(c.ctx, c.org, opts)
			if err != nil {
				return nil, fmt.Errorf("failed to list %s issues: %w", state, err)
			}

			for _, issue := range issues {
				var labels []string
				for _, label := range issue.Labels {
					labels = append(labels, label.GetName())
				}

				allIssues = append(allIssues, ProjectIssue{
					ID:        int(issue.GetID()),
					Title:     issue.GetTitle(),
					Body:      issue.GetBody(),
					Number:    issue.GetNumber(),
					State:     issue.GetState(),
					URL:       issue.GetHTMLURL(),
					Labels:    labels,
					UpdatedAt: issue.GetUpdatedAt().Time,
				})
			}

			if resp.NextPage == 0 {
				break
			}
			opts.Page = resp.NextPage
		}
	}

	return allIssues, nil
}

func (c *Client) getProjectStatusForIssues(projectNumber int, assignedIssues []ProjectIssue, statusFilters []string) ([]IssueWithStatus, error) {
	fmt.Printf("Checking project status for %d assigned issues...\n", len(assignedIssues))
	
	// Create map of issue URLs for quick lookup
	issueURLMap := make(map[string]ProjectIssue)
	for _, issue := range assignedIssues {
		issueURLMap[issue.URL] = issue
	}
	
	var query struct {
		Organization struct {
			ProjectV2 struct {
				Items struct {
					PageInfo struct {
						HasNextPage bool   `graphql:"hasNextPage"`
						EndCursor   string `graphql:"endCursor"`
					} `graphql:"pageInfo"`
					Nodes []struct {
						ID          string `graphql:"id"`
						FieldValues struct {
							Nodes []struct {
								SingleSelectValue struct {
									Name  string `graphql:"name"`
									Field struct {
										SingleSelectField struct {
											Name string `graphql:"name"`
										} `graphql:"... on ProjectV2SingleSelectField"`
									} `graphql:"field"`
								} `graphql:"... on ProjectV2ItemFieldSingleSelectValue"`
							}
						} `graphql:"fieldValues(first: 10)"`
						Content struct {
							Typename string `graphql:"__typename"`
							Issue    struct {
								URL string `graphql:"url"`
							} `graphql:"... on Issue"`
						} `graphql:"content"`
					}
				} `graphql:"items(first: 100, after: $cursor)"`
			} `graphql:"projectV2(number: $projectNumber)"`
		} `graphql:"organization(login: $org)"`
	}

	var result []IssueWithStatus
	var cursor *string

	for {
		variables := map[string]interface{}{
			"org":           githubv4.String(c.org),
			"projectNumber": githubv4.Int(projectNumber),
		}
		
		if cursor != nil {
			variables["cursor"] = githubv4.String(*cursor)
		} else {
			variables["cursor"] = (*githubv4.String)(nil)
		}

		err := c.graphqlClient.Query(c.ctx, &query, variables)
		if err != nil {
			return nil, fmt.Errorf("failed to query project: %w", err)
		}

		for _, item := range query.Organization.ProjectV2.Items.Nodes {
			if item.Content.Typename != "Issue" {
				continue
			}

			// Check if this is one of our assigned issues
			if assignedIssue, exists := issueURLMap[item.Content.Issue.URL]; exists {
				var currentStatus string
				for _, fieldValue := range item.FieldValues.Nodes {
					if fieldValue.SingleSelectValue.Field.SingleSelectField.Name == "Status" {
						currentStatus = fieldValue.SingleSelectValue.Name
						break
					}
				}

				// Check if status matches our filters
				hasMatchingStatus := false
				for _, status := range statusFilters {
					if strings.EqualFold(currentStatus, status) {
						hasMatchingStatus = true
						break
					}
				}

				if hasMatchingStatus {
					result = append(result, IssueWithStatus{
						ProjectIssue:  assignedIssue,
						GitHubStatus:  currentStatus,
						ProjectItemID: item.ID,
						UpdatedAt:     assignedIssue.UpdatedAt,
					})
				}
			}
		}

		if !query.Organization.ProjectV2.Items.PageInfo.HasNextPage {
			break
		}
		cursor = &query.Organization.ProjectV2.Items.PageInfo.EndCursor
	}

	fmt.Printf("Found %d assigned issues matching status filters\n", len(result))
	return result, nil
}

type ProjectMetadata struct {
	ID            string
	StatusFieldID string
	StatusOptions map[string]string
}

func (c *Client) getProjectMetadata(projectNumber int) (*ProjectMetadata, error) {
	fmt.Printf("DEBUG: Fetching project metadata for project %d\n", projectNumber)
	
	// Single query - get project ID, fields and status options
	var query struct {
		Organization struct {
			ProjectV2 struct {
				ID     string `graphql:"id"`
				Fields struct {
					Nodes []struct {
						Typename         string `graphql:"__typename"`
						ProjectV2Field   struct {
							ID   string `graphql:"id"`
							Name string `graphql:"name"`
						} `graphql:"... on ProjectV2Field"`
						SingleSelectField struct {
							ID      string `graphql:"id"`
							Name    string `graphql:"name"`
							Options []struct {
								ID   string `graphql:"id"`
								Name string `graphql:"name"`
							} `graphql:"options"`
						} `graphql:"... on ProjectV2SingleSelectField"`
					}
				} `graphql:"fields(first: 20)"`
			} `graphql:"projectV2(number: $projectNumber)"`
		} `graphql:"organization(login: $org)"`
	}

	variables := map[string]interface{}{
		"org":           githubv4.String(c.org),
		"projectNumber": githubv4.Int(projectNumber),
	}

	fmt.Printf("DEBUG: Executing GraphQL query for project metadata\n")
	err := c.graphqlClient.Query(c.ctx, &query, variables)
	if err != nil {
		fmt.Printf("DEBUG: GraphQL query failed: %v\n", err)
		return nil, fmt.Errorf("failed to query project metadata: %w", err)
	}
	
	fmt.Printf("DEBUG: GraphQL query successful, found %d fields\n", len(query.Organization.ProjectV2.Fields.Nodes))

	metadata := &ProjectMetadata{
		ID:            query.Organization.ProjectV2.ID,
		StatusOptions: make(map[string]string),
	}

	// Find the Status field and its options
	for _, field := range query.Organization.ProjectV2.Fields.Nodes {
		var fieldID, fieldName string
		
		if field.Typename == "ProjectV2SingleSelectField" {
			fieldID = field.SingleSelectField.ID
			fieldName = field.SingleSelectField.Name
			fmt.Printf("DEBUG: SingleSelectField - Name: %s, ID: %s\n", fieldName, fieldID)
			
			if fieldName == "Status" {
				fmt.Printf("DEBUG: Found Status field with ID: %s\n", fieldID)
				metadata.StatusFieldID = fieldID
				
				fmt.Printf("DEBUG: Found %d status options\n", len(field.SingleSelectField.Options))
				for _, option := range field.SingleSelectField.Options {
					fmt.Printf("DEBUG: Status option - Name: %s, ID: %s\n", option.Name, option.ID)
					metadata.StatusOptions[option.Name] = option.ID
				}
				break
			}
		} else if field.Typename == "ProjectV2Field" {
			fieldID = field.ProjectV2Field.ID
			fieldName = field.ProjectV2Field.Name
			fmt.Printf("DEBUG: ProjectV2Field - Name: %s, ID: %s\n", fieldName, fieldID)
		}
	}

	if metadata.StatusFieldID == "" {
		return nil, fmt.Errorf("Status field not found in project")
	}

	fmt.Printf("DEBUG: Project metadata - ID: %s, StatusFieldID: %s, Options: %v\n", 
		metadata.ID, metadata.StatusFieldID, metadata.StatusOptions)

	return metadata, nil
}

func (c *Client) UpdateProjectItemStatus(projectNumber int, projectItemID string, newStatus string) error {
	// Get fresh project metadata
	metadata, err := c.getProjectMetadata(projectNumber)
	if err != nil {
		return fmt.Errorf("failed to get project metadata: %w", err)
	}

	optionID, exists := metadata.StatusOptions[newStatus]
	if !exists {
		return fmt.Errorf("unsupported status: %s (available: %v)", newStatus, getKeys(metadata.StatusOptions))
	}

	var mutation struct {
		UpdateProjectV2ItemFieldValue struct {
			ProjectV2Item struct {
				ID string
			}
		} `graphql:"updateProjectV2ItemFieldValue(input: $input)"`
	}

	fmt.Printf("DEBUG: Updating project item %s to status %s (optionID: %s)\n", projectItemID, newStatus, optionID)

	input := map[string]interface{}{
		"projectId": githubv4.String(metadata.ID),
		"itemId":    githubv4.String(projectItemID),
		"fieldId":   githubv4.String(metadata.StatusFieldID),
		"value": map[string]interface{}{
			"singleSelectOptionId": githubv4.String(optionID),
		},
	}

	err = c.graphqlClient.Mutate(c.ctx, &mutation, input, nil)
	if err != nil {
		return fmt.Errorf("failed to update project item status: %w", err)
	}

	fmt.Printf("DEBUG: Mutation completed successfully, result: %+v\n", mutation)
	return nil
}

func getKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func (c *Client) CloseIssue(repoName string, issueNumber int) error {
	state := "closed"
	issueRequest := &github.IssueRequest{
		State: &state,
	}

	_, _, err := c.restClient.Issues.Edit(c.ctx, c.org, repoName, issueNumber, issueRequest)
	if err != nil {
		return fmt.Errorf("failed to close issue: %w", err)
	}

	return nil
}

func (c *Client) UpdateIssueFromURL(issueURL string) error {
	parts := strings.Split(issueURL, "/")
	if len(parts) < 2 {
		return fmt.Errorf("invalid issue URL")
	}

	repoName := parts[len(parts)-3]
	issueNumberStr := parts[len(parts)-1]
	issueNumber, err := strconv.Atoi(issueNumberStr)
	if err != nil {
		return fmt.Errorf("invalid issue number: %w", err)
	}

	return c.CloseIssue(repoName, issueNumber)
}

func (c *Client) getIssueUpdateTime(issueURL string) (time.Time, error) {
	parts := strings.Split(issueURL, "/")
	if len(parts) < 7 {
		return time.Time{}, fmt.Errorf("invalid GitHub URL format")
	}

	repoName := parts[4]
	issueNumberStr := parts[6]
	issueNumber, err := strconv.Atoi(issueNumberStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid issue number: %w", err)
	}

	issue, _, err := c.restClient.Issues.Get(c.ctx, c.org, repoName, issueNumber)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to get issue: %w", err)
	}

	return issue.GetUpdatedAt().Time, nil
}

func (c *Client) lookupProjectItemID(projectNumber int, issueURL string) (string, error) {
	// Query the project to find the item ID for this specific issue URL
	var query struct {
		Organization struct {
			ProjectV2 struct {
				Items struct {
					Nodes []struct {
						ID      string `graphql:"id"`
						Content struct {
							Typename string `graphql:"__typename"`
							Issue    struct {
								URL string `graphql:"url"`
							} `graphql:"... on Issue"`
						} `graphql:"content"`
					}
				} `graphql:"items(first: 100)"`
			} `graphql:"projectV2(number: $projectNumber)"`
		} `graphql:"organization(login: $org)"`
	}

	variables := map[string]interface{}{
		"org":           githubv4.String(c.org),
		"projectNumber": githubv4.Int(projectNumber),
	}

	err := c.graphqlClient.Query(c.ctx, &query, variables)
	if err != nil {
		return "", fmt.Errorf("failed to query project for lookup: %w", err)
	}

	// Find the item with matching issue URL
	for _, item := range query.Organization.ProjectV2.Items.Nodes {
		if item.Content.Typename == "Issue" && item.Content.Issue.URL == issueURL {
			return item.ID, nil
		}
	}

	return "", fmt.Errorf("project item not found for issue URL: %s", issueURL)
}
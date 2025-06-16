package cli

import (
	"fmt"
	"strings"

	"github.com/arjungandhi/atp/project"
	"github.com/arjungandhi/atp/todo"
	"github.com/arjungandhi/go-utils/pkg/prompt"
	"github.com/arjungandhi/go-utils/pkg/shell"
	bonzai "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/help"
)

var ProjectCmd = &bonzai.Cmd{
	Name:    "project",
	Aliases: []string{"p"},
	Summary: "manage projects with phases and repository linking",
	Commands: []*bonzai.Cmd{
		help.Cmd,
		projectDocCmd,
		projectEditCmd,
		projectAddCmd,
		projectDeleteCmd,
		projectFinishCmd,
		projectReorgCmd,
		projectActivateCmd,
		projectDeactivateCmd,
	},
}

var projectEditCmd = &bonzai.Cmd{
	Name:     "edit",
	Aliases:  []string{"e"},
	Summary:  "edit active projects in your default editor",
	Commands: []*bonzai.Cmd{help.Cmd, projectEditAllCmd},
	Call: func(cmd *bonzai.Cmd, args ...string) error {
		// get the active projects, path
		path, err := ProjectDir()
		if err != nil {
			return err
		}

		// get active projects
		project_path := project.ActiveFilePath(path)

		// Open the projects file in the editor
		shell.OpenInEditor(project_path)

		return nil
	},
}

var projectEditAllCmd = &bonzai.Cmd{
	Name:     "all",
	Aliases:  []string{"a"},
	Summary:  "edit both active and completed projects",
	Commands: []*bonzai.Cmd{help.Cmd},
	Call: func(cmd *bonzai.Cmd, args ...string) error {

		path, err := ProjectDir()
		if err != nil {
			return err
		}
		// get all projects
		done_path := project.DoneFilePath(path)
		active_path := project.ActiveFilePath(path)

		// Open the projects file in the editor
		shell.OpenInEditor(active_path, done_path)

		return nil
	},
}

var projectAddCmd = &bonzai.Cmd{
	Name:     "add",
	Aliases:  []string{"a"},
	Summary:  "add a new project",
	Commands: []*bonzai.Cmd{help.Cmd},
	Call: func(cmd *bonzai.Cmd, args ...string) (err error) {
		// args 0 should be the todo string
		todo_str := strings.Join(args, " ")
		if len(args) == 0 {
			// prompt for todo
			todo_str, err = prompt.PromptString("Enter a project: ")
			if err != nil {
				return err
			}
		}

		// get repos
		repos, err := GetRepos()
		if err != nil {
			return err
		}

		// parse to string as a todo
		todo := todo.FromString(todo_str)

		// parse todo into project
		new_project, err := project.FromTodo(todo, repos)
		if err != nil {
			return err
		}

		projects, err := GetProjects()
		if err != nil {
			return err
		}

		projects = append(projects, new_project)

		// write the projects to the file
		err = WriteProjects(projects)
		if err != nil {
			return err
		}

		fmt.Printf("Added project: %s\n", new_project.String())

		return nil
	},
}

var projectDeleteCmd = &bonzai.Cmd{
	Name:     "delete",
	Aliases:  []string{"del", "d"},
	Summary:  "delete a project from the project list",
	Commands: []*bonzai.Cmd{help.Cmd},
	Call: func(cmd *bonzai.Cmd, args ...string) error {
		// get input from user
		input := strings.Join(args, " ")

		// load projects
		projects, err := GetProjects()
		if err != nil {
			return err
		}

		// get the project to delete
		index, err := shell.FzfSearch(projects, input)
		if err != nil {
			return err
		}

		// remove the project
		projects = append(projects[:index], projects[index+1:]...)

		// write the projects to the file
		err = WriteProjects(projects)
		if err != nil {
			return err
		}

		fmt.Printf("Deleted project: %s\n", projects[index].String())

		return nil
	},
}

var projectActivateCmd = &bonzai.Cmd{
	Name:     "activate",
	Summary:  "activate an inactive project and optionally link to repository",
	Commands: []*bonzai.Cmd{help.Cmd},
	Call: func(cmd *bonzai.Cmd, args ...string) error {
		// get input from user
		input := strings.Join(args, " ")

		// load projects
		projects, err := GetProjects()
		if err != nil {
			return err
		}

		// get all inactive projects
		inactive_projects := []*project.Project{}
		for _, project := range projects {
			if !project.Active {
				inactive_projects = append(inactive_projects, project)
			}
		}

		// get the project
		index, err := shell.FzfSearch(inactive_projects, input)
		if err != nil {
			return err
		}

		selection := inactive_projects[index]
		// mark the project as active
		selection.Active = true
		// set the phase to 1 if it is not set
		if selection.Phase == "" {
			selection.Phase = "1"
		}

		if selection.Repo == nil {

			// load repos
			repos, err := GetRepos()
			if err != nil {
				return err
			}

			// have the user select the repo
			repo_index, err := shell.FzfSearch(repos, "")
			if err != nil {
				return err
			}

			// set the repo
			selection.Repo = repos[repo_index]
		}

		// write the projects to the file
		err = WriteProjects(projects)
		if err != nil {
			return err
		}

		fmt.Printf("Activated project: %s\n", selection.String())

		return nil
	},
}

var projectDeactivateCmd = &bonzai.Cmd{
	Name:     "deactivate",
	Summary:  "deactivate an active project",
	Commands: []*bonzai.Cmd{help.Cmd},
	Call: func(cmd *bonzai.Cmd, args ...string) error {
		// get input from user
		input := strings.Join(args, " ")

		// load projects
		projects, err := GetProjects()
		if err != nil {
			return err
		}

		// get all active projects
		active_projects := []*project.Project{}
		for _, project := range projects {
			if project.Active {
				active_projects = append(active_projects, project)
			}
		}

		// get the project
		index, err := shell.FzfSearch(active_projects, input)
		if err != nil {
			return err
		}

		// mark the project as inactive
		active_projects[index].Active = false

		err = WriteProjects(projects)

		fmt.Printf("Deactivated project: %s\n", active_projects[index].String())

		return nil
	},
}

var projectFinishCmd = &bonzai.Cmd{
	Name:     "finish",
	Summary:  "mark a project as completed and move to done.txt",
	Commands: []*bonzai.Cmd{help.Cmd},
	Call: func(cmd *bonzai.Cmd, args ...string) error {
		// get input from user
		input := strings.Join(args, " ")

		// load projects
		projects, err := GetProjects()
		if err != nil {
			return err
		}

		// get all not done projects
		not_done_projects := []*project.Project{}
		for _, project := range projects {
			if !project.Done {
				not_done_projects = append(not_done_projects, project)
			}
		}

		// get the project
		index, err := shell.FzfSearch(not_done_projects, input)
		if err != nil {
			return err
		}

		// mark the project as done
		not_done_projects[index].Active = false
		not_done_projects[index].Done = true

		err = WriteProjects(projects)

		fmt.Printf("Finished project: %s\n", not_done_projects[index].String())

		return nil

	},
}

var projectDocCmd = &bonzai.Cmd{
	Name:    "doc",
	Summary: "open project documentation file in editor",
	Commands: []*bonzai.Cmd{
		help.Cmd,
	},
	Call: func(cmd *bonzai.Cmd, args ...string) error {

		projects, err := GetProjects()
		if err != nil {
			return err
		}

		// get the projects with a repo
		repo_projects := []*project.Project{}
		for _, project := range projects {
			if project.Repo != nil {
				repo_projects = append(repo_projects, project)
			}
		}

		// get the project to get the doc for
		index, err := shell.FzfSearch(repo_projects, "")
		if err != nil {
			return err
		}

		// get the doc
		doc_path := repo_projects[index].Repo.GetRepoDoc()

		// open the doc in the editor
		shell.OpenInEditor(doc_path)

		return nil
	},
}

var projectReorgCmd = &bonzai.Cmd{
	Name:    "reorg",
	Summary: "reorganize and sort projects, separating completed from active",
	Commands: []*bonzai.Cmd{
		help.Cmd,
	},
	Call: func(cmd *bonzai.Cmd, args ...string) error {
		// get all the projects
		projects, err := GetProjects()
		if err != nil {
			return err
		}

		// write the projects to the file
		err = WriteProjects(projects)
		if err != nil {
			return err
		}

		fmt.Println("Reorganized projects")

		return nil
	},
}

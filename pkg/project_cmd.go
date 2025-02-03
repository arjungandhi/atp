package atp

import (
	"fmt"

	"github.com/arjungandhi/go-utils/pkg/shell"
	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/help"
)

var ProjectCmd = &Z.Cmd{
	Name:    "project",
	Summary: "project is a command line tool for managing projects",
	Commands: []*Z.Cmd{
		help.Cmd,
		projectListCmd,
		projectDocCmd,
	},
}

var projectListCmd = &Z.Cmd{
	Name:    "list",
	Summary: "list all projects",
	Commands: []*Z.Cmd{
		help.Cmd,
	},
	Call: func(cmd *Z.Cmd, args ...string) error {

		projects, err := getRepoProjects()
		if err != nil {
			return err
		}

		// print the list of projects
		for _, project := range projects {
			fmt.Println(project)
		}

		return nil
	},
}

var projectDocCmd = &Z.Cmd{
	Name:    "doc",
	Summary: "gets the doc for a project and opens it in the editor",
	Commands: []*Z.Cmd{
		help.Cmd,
	},
	Call: func(cmd *Z.Cmd, args ...string) error {

		projects, err := getRepoProjects()
		if err != nil {
			return err
		}

		// are we in a project?
		project, err := userInProject(projects)
		if err != nil {
			return err
		}

		if project == nil {
			// prompt for project
			result_index, err := shell.FzfSearch(projects, "")
			if err != nil {
				return err
			}

			project = projects[result_index]
		}

		// launch the editor
		// editor should make the doc if it does not exist, I've only tested this with vim
		shell.OpenInEditor(project.getProjectDoc())

		return nil
	},
}

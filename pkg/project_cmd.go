package atp

import (
	"github.com/arjungandhi/go-utils/pkg/shell"
	bonzai "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/help"
)

var ProjectCmd = &bonzai.Cmd{
	Name:    "project",
	Aliases: []string{"proj", "p"},
	Summary: "project is a command line tool for managing projects",
	Commands: []*bonzai.Cmd{
		help.Cmd,
		projectDocCmd,
	},
}

var projectEditCmd = &bonzai.Cmd{
	Name:     "edit",
	Aliases:  []string{"e"},
	Summary:  "edit the projects",
	Commands: []*bonzai.Cmd{help.Cmd},
	Call: func(cmd *bonzai.Cmd, args ...string) error {
		//
		return nil
	},
}

var projectEditAllCmd = &bonzai.Cmd{
	Name:     "all",
	Aliases:  []string{"a"},
	Summary:  "edit all projects",
	Commands: []*bonzai.Cmd{help.Cmd},
	Call: func(cmd *bonzai.Cmd, args ...string) error {
		return nil
	},
}

var projectAddCmd = &bonzai.Cmd{
	Name:     "add",
	Aliases:  []string{"a"},
	Summary:  "add a project to the projects.txt",
	Commands: []*bonzai.Cmd{help.Cmd},
	Call: func(cmd *bonzai.Cmd, args ...string) error {
		return nil
	},
}

var projectDeleteCmd = &bonzai.Cmd{
	Name:     "delete",
	Aliases:  []string{"del", "d"},
	Summary:  "del a project to the projects.txt",
	Commands: []*bonzai.Cmd{help.Cmd},
	Call: func(cmd *bonzai.Cmd, args ...string) error {
		return nil
	},
}

var projectActivateCmd = &bonzai.Cmd{
	Name:     "activate",
	Summary:  "activate a project",
	Commands: []*bonzai.Cmd{help.Cmd},
	Call: func(cmd *bonzai.Cmd, args ...string) error {
		return nil
	},
}

var projectDeactivateCmd = &bonzai.Cmd{
	Name:     "deactivate",
	Summary:  "deactivate a project",
	Commands: []*bonzai.Cmd{help.Cmd},
	Call: func(cmd *bonzai.Cmd, args ...string) error {
		return nil
	},
}

var projectFinishCmd = &bonzai.Cmd{
	Name:     "finish",
	Summary:  "finish a project",
	Commands: []*bonzai.Cmd{help.Cmd},
	Call: func(cmd *bonzai.Cmd, args ...string) error {
		return nil
	},
}

var projectDocCmd = &bonzai.Cmd{
	Name:    "doc",
	Summary: "gets the doc for a project and opens it in the editor",
	Commands: []*bonzai.Cmd{
		help.Cmd,
	},
	Call: func(cmd *bonzai.Cmd, args ...string) error {

		projects, err := getRepos()
		if err != nil {
			return err
		}

		// are we in a repo?
		project, err := userInRepo(projects)
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

		return nil
	},
}

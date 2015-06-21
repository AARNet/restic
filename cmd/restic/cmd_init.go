package main

import (
	"errors"

	"github.com/restic/restic/repository"
)

type CmdInit struct {
	global *GlobalOptions
}

func (cmd CmdInit) Execute(args []string) error {
	if cmd.global.Repo == "" {
		return errors.New("Please specify repository location (-r)")
	}

	if cmd.global.password == "" {
		cmd.global.password = cmd.global.ReadPasswordTwice(
			"enter password for new backend: ",
			"enter password again: ")
	}

	be, err := create(cmd.global.Repo)
	if err != nil {
		cmd.global.Exitf(1, "creating backend at %s failed: %v\n", cmd.global.Repo, err)
	}

	s := repository.New(be)
	err = s.Init(cmd.global.password)
	if err != nil {
		cmd.global.Exitf(1, "creating key in backend at %s failed: %v\n", cmd.global.Repo, err)
	}

	cmd.global.Printf("created restic backend %v at %s\n", s.Config.ID[:10], cmd.global.Repo)
	cmd.global.Printf("\n")
	cmd.global.Printf("Please note that knowledge of your password is required to access\n")
	cmd.global.Printf("the repository. Losing your password means that your data is\n")
	cmd.global.Printf("irrecoverably lost.\n")

	return nil
}

func init() {
	_, err := parser.AddCommand("init",
		"create repository",
		"The init command creates a new repository",
		&CmdInit{global: &globalOpts})
	if err != nil {
		panic(err)
	}
}

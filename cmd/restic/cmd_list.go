package main

import (
	"fmt"

	"github.com/restic/restic/internal/errors"
	"github.com/restic/restic/internal/index"
	"github.com/restic/restic/internal/restic"

	"github.com/spf13/cobra"
)

var cmdList = &cobra.Command{
	Use:   "list [blobs|packs|index|snapshots|keys|locks]",
	Short: "List objects in the repository",
	Long: `
The "list" command allows listing objects in the repository based on type.
`,
	DisableAutoGenTag: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runList(globalOptions, args)
	},
}

func init() {
	cmdRoot.AddCommand(cmdList)
}

func runList(opts GlobalOptions, args []string) error {
	if len(args) != 1 {
		return errors.Fatal("type not specified")
	}

	repo, err := OpenRepository(opts)
	if err != nil {
		return err
	}

	if !opts.NoLock {
		lock, err := lockRepo(repo)
		defer unlockRepo(lock)
		if err != nil {
			return err
		}
	}

	var t restic.FileType
	switch args[0] {
	case "packs":
		t = restic.DataFile
	case "index":
		t = restic.IndexFile
	case "snapshots":
		t = restic.SnapshotFile
	case "keys":
		t = restic.KeyFile
	case "locks":
		t = restic.LockFile
	case "blobs":
		idx, err := index.Load(opts.ctx, repo, nil)
		if err != nil {
			return err
		}

		for _, pack := range idx.Packs {
			for _, entry := range pack.Entries {
				fmt.Printf("%v %v\n", entry.Type, entry.ID)
			}
		}

		return nil
	default:
		return errors.Fatal("invalid type")
	}

	for id := range repo.List(opts.ctx, t) {
		Printf("%s\n", id)
	}

	return nil
}

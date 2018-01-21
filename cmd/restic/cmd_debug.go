// +build debug

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/restic/restic/internal/errors"
	"github.com/restic/restic/internal/pack"
	"github.com/restic/restic/internal/repository"
	"github.com/restic/restic/internal/restic"
)

var cmdDebug = &cobra.Command{
	Use:   "debug",
	Short: "Debug commands",
}

var cmdDebugDump = &cobra.Command{
	Use:   "dump [indexes|snapshots|all|packs]",
	Short: "Dump data structures",
	Long: `
The "dump" command dumps data structures from the repository as JSON objects. It
is used for debugging purposes only.`,
	DisableAutoGenTag: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDebugDump(globalOptions, args)
	},
}

func init() {
	cmdRoot.AddCommand(cmdDebug)
	cmdDebug.AddCommand(cmdDebugDump)
}

func prettyPrintJSON(wr io.Writer, item interface{}) error {
	buf, err := json.MarshalIndent(item, "", "  ")
	if err != nil {
		return err
	}

	_, err = wr.Write(append(buf, '\n'))
	return err
}

func debugPrintSnapshots(repo *repository.Repository, wr io.Writer) error {
	return repo.List(context.TODO(), restic.SnapshotFile, func(id restic.ID, size int64) error {
		snapshot, err := restic.LoadSnapshot(context.TODO(), repo, id)
		if err != nil {
			return err
		}

		fmt.Fprintf(wr, "snapshot_id: %v\n", id)

		return prettyPrintJSON(wr, snapshot)
	})
}

// Pack is the struct used in printPacks.
type Pack struct {
	Name string `json:"name"`

	Blobs []Blob `json:"blobs"`
}

// Blob is the struct used in printPacks.
type Blob struct {
	Type   restic.BlobType `json:"type"`
	Length uint            `json:"length"`
	ID     restic.ID       `json:"id"`
	Offset uint            `json:"offset"`
}

func printPacks(repo *repository.Repository, wr io.Writer) error {

	return repo.List(context.TODO(), restic.DataFile, func(id restic.ID, size int64) error {
		h := restic.Handle{Type: restic.DataFile, Name: id.String()}

		blobs, err := pack.List(repo.Key(), restic.ReaderAt(repo.Backend(), h), size)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error for pack %v: %v\n", id.Str(), err)
			return nil
		}

		p := Pack{
			Name:  id.String(),
			Blobs: make([]Blob, len(blobs)),
		}
		for i, blob := range blobs {
			p.Blobs[i] = Blob{
				Type:   blob.Type,
				Length: blob.Length,
				ID:     blob.ID,
				Offset: blob.Offset,
			}
		}

		return prettyPrintJSON(os.Stdout, p)
	})

	return nil
}

func dumpIndexes(repo restic.Repository) error {
	return repo.List(context.TODO(), restic.IndexFile, func(id restic.ID, size int64) error {
		fmt.Printf("index_id: %v\n", id)

		idx, err := repository.LoadIndex(context.TODO(), repo, id)
		if err != nil {
			return err
		}

		return idx.Dump(os.Stdout)
	})
}

func runDebugDump(gopts GlobalOptions, args []string) error {
	if len(args) != 1 {
		return errors.Fatal("type not specified")
	}

	repo, err := OpenRepository(gopts)
	if err != nil {
		return err
	}

	if !gopts.NoLock {
		lock, err := lockRepo(repo)
		defer unlockRepo(lock)
		if err != nil {
			return err
		}
	}

	err = repo.LoadIndex(gopts.ctx)
	if err != nil {
		return err
	}

	tpe := args[0]

	switch tpe {
	case "indexes":
		return dumpIndexes(repo)
	case "snapshots":
		return debugPrintSnapshots(repo, os.Stdout)
	case "packs":
		return printPacks(repo, os.Stdout)
	case "all":
		fmt.Printf("snapshots:\n")
		err := debugPrintSnapshots(repo, os.Stdout)
		if err != nil {
			return err
		}

		fmt.Printf("\nindexes:\n")
		err = dumpIndexes(repo)
		if err != nil {
			return err
		}

		return nil
	default:
		return errors.Fatalf("no such type %q", tpe)
	}
}

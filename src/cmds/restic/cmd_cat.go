package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"restic"
	"restic/backend"
	"restic/debug"
	"restic/errors"
	"restic/repository"
)

var cmdCat = &cobra.Command{
	Use:   "cat [flags] [pack|blob|tree|snapshot|key|masterkey|config|lock] ID",
	Short: "print internal objects to stdout",
	Long: `
The "cat" command is used to print internal objects to stdout.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runCat(globalOptions, args)
	},
}

func init() {
	cmdRoot.AddCommand(cmdCat)
}

func runCat(gopts GlobalOptions, args []string) error {
	if len(args) < 1 || (args[0] != "masterkey" && args[0] != "config" && len(args) != 2) {
		return errors.Fatalf("type or ID not specified")
	}

	repo, err := OpenRepository(gopts)
	if err != nil {
		return err
	}

	lock, err := lockRepo(repo)
	defer unlockRepo(lock)
	if err != nil {
		return err
	}

	tpe := args[0]

	var id restic.ID
	if tpe != "masterkey" && tpe != "config" {
		id, err = restic.ParseID(args[1])
		if err != nil {
			if tpe != "snapshot" {
				return errors.Fatalf("unable to parse ID: %v\n", err)
			}

			// find snapshot id with prefix
			id, err = restic.FindSnapshot(repo, args[1])
			if err != nil {
				return err
			}
		}
	}

	// handle all types that don't need an index
	switch tpe {
	case "config":
		buf, err := json.MarshalIndent(repo.Config(), "", "  ")
		if err != nil {
			return err
		}

		fmt.Println(string(buf))
		return nil
	case "index":
		buf, err := repo.LoadAndDecrypt(restic.IndexFile, id)
		if err != nil {
			return err
		}

		_, err = os.Stdout.Write(append(buf, '\n'))
		return err

	case "snapshot":
		sn := &restic.Snapshot{}
		err = repo.LoadJSONUnpacked(restic.SnapshotFile, id, sn)
		if err != nil {
			return err
		}

		buf, err := json.MarshalIndent(&sn, "", "  ")
		if err != nil {
			return err
		}

		fmt.Println(string(buf))

		return nil
	case "key":
		h := restic.Handle{Type: restic.KeyFile, Name: id.String()}
		buf, err := backend.LoadAll(repo.Backend(), h, nil)
		if err != nil {
			return err
		}

		key := &repository.Key{}
		err = json.Unmarshal(buf, key)
		if err != nil {
			return err
		}

		buf, err = json.MarshalIndent(&key, "", "  ")
		if err != nil {
			return err
		}

		fmt.Println(string(buf))
		return nil
	case "masterkey":
		buf, err := json.MarshalIndent(repo.Key(), "", "  ")
		if err != nil {
			return err
		}

		fmt.Println(string(buf))
		return nil
	case "lock":
		lock, err := restic.LoadLock(repo, id)
		if err != nil {
			return err
		}

		buf, err := json.MarshalIndent(&lock, "", "  ")
		if err != nil {
			return err
		}

		fmt.Println(string(buf))

		return nil
	}

	// load index, handle all the other types
	err = repo.LoadIndex()
	if err != nil {
		return err
	}

	switch tpe {
	case "pack":
		h := restic.Handle{Type: restic.DataFile, Name: id.String()}
		buf, err := backend.LoadAll(repo.Backend(), h, nil)
		if err != nil {
			return err
		}

		hash := restic.Hash(buf)
		if !hash.Equal(id) {
			fmt.Fprintf(stderr, "Warning: hash of data does not match ID, want\n  %v\ngot:\n  %v\n", id.String(), hash.String())
		}

		_, err = os.Stdout.Write(buf)
		return err

	case "blob":
		for _, t := range []restic.BlobType{restic.DataBlob, restic.TreeBlob} {
			list, err := repo.Index().Lookup(id, t)
			if err != nil {
				continue
			}
			blob := list[0]

			buf := make([]byte, blob.Length)
			n, err := repo.LoadBlob(restic.DataBlob, id, buf)
			if err != nil {
				return err
			}
			buf = buf[:n]

			_, err = os.Stdout.Write(buf)
			return err
		}

		return errors.Fatal("blob not found")

	case "tree":
		debug.Log("cat tree %v", id.Str())
		tree, err := repo.LoadTree(id)
		if err != nil {
			debug.Log("unable to load tree %v: %v", id.Str(), err)
			return err
		}

		buf, err := json.MarshalIndent(&tree, "", "  ")
		if err != nil {
			debug.Log("error json.MarshalIndent(): %v", err)
			return err
		}

		_, err = os.Stdout.Write(append(buf, '\n'))
		return nil

	default:
		return errors.Fatal("invalid type")
	}
}

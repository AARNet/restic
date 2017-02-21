package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"restic"
	"strings"

	"github.com/spf13/cobra"
)

var cmdForget = &cobra.Command{
	Use:   "forget [flags] [snapshot ID] [...]",
	Short: "forget removes snapshots from the repository",
	Long: `
The "forget" command removes snapshots according to a policy. Please note that
this command really only deletes the snapshot object in the repository, which
is a reference to data stored there. In order to remove this (now unreferenced)
data after 'forget' was run successfully, see the 'prune' command. `,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runForget(forgetOptions, globalOptions, args)
	},
}

// ForgetOptions collects all options for the forget command.
type ForgetOptions struct {
	Last    int
	Hourly  int
	Daily   int
	Weekly  int
	Monthly int
	Yearly  int

	KeepTags []string

	Hostname string
	Tags     []string

	DryRun bool
	Prune  bool
}

var forgetOptions ForgetOptions

func init() {
	cmdRoot.AddCommand(cmdForget)

	f := cmdForget.Flags()
	f.IntVarP(&forgetOptions.Last, "keep-last", "l", 0, "keep the last `n` snapshots")
	f.IntVarP(&forgetOptions.Hourly, "keep-hourly", "H", 0, "keep the last `n` hourly snapshots")
	f.IntVarP(&forgetOptions.Daily, "keep-daily", "d", 0, "keep the last `n` daily snapshots")
	f.IntVarP(&forgetOptions.Weekly, "keep-weekly", "w", 0, "keep the last `n` weekly snapshots")
	f.IntVarP(&forgetOptions.Monthly, "keep-monthly", "m", 0, "keep the last `n` monthly snapshots")
	f.IntVarP(&forgetOptions.Yearly, "keep-yearly", "y", 0, "keep the last `n` yearly snapshots")

	f.StringSliceVar(&forgetOptions.KeepTags, "keep-tag", []string{}, "always keep snapshots with this `tag` (can be specified multiple times)")
	f.StringVar(&forgetOptions.Hostname, "hostname", "", "only forget snapshots for the given hostname")
	f.StringSliceVar(&forgetOptions.Tags, "tag", []string{}, "only forget snapshots with the `tag` (can be specified multiple times)")

	f.BoolVarP(&forgetOptions.DryRun, "dry-run", "n", false, "do not delete anything, just print what would be done")
	f.BoolVar(&forgetOptions.Prune, "prune", false, "automatically run the 'prune' command if snapshots have been removed")
}

func printSnapshots(w io.Writer, snapshots restic.Snapshots) {
	tab := NewTable()
	tab.Header = fmt.Sprintf("%-8s  %-19s  %-10s  %-10s  %s", "ID", "Date", "Host", "Tags", "Directory")
	tab.RowFormat = "%-8s  %-19s  %-10s  %-10s  %s"

	for _, sn := range snapshots {
		if len(sn.Paths) == 0 {
			continue
		}

		firstTag := ""
		if len(sn.Tags) > 0 {
			firstTag = sn.Tags[0]
		}

		tab.Rows = append(tab.Rows, []interface{}{sn.ID().Str(), sn.Time.Format(TimeFormat), sn.Hostname, firstTag, sn.Paths[0]})

		rows := len(sn.Paths)
		if len(sn.Tags) > rows {
			rows = len(sn.Tags)
		}

		for i := 1; i < rows; i++ {
			path := ""
			if len(sn.Paths) > i {
				path = sn.Paths[i]
			}

			tag := ""
			if len(sn.Tags) > i {
				tag = sn.Tags[i]
			}

			tab.Rows = append(tab.Rows, []interface{}{"", "", "", tag, path})
		}
	}

	tab.Write(w)
}

func runForget(opts ForgetOptions, gopts GlobalOptions, args []string) error {
	repo, err := OpenRepository(gopts)
	if err != nil {
		return err
	}

	lock, err := lockRepoExclusive(repo)
	defer unlockRepo(lock)
	if err != nil {
		return err
	}

	// parse arguments as hex strings
	var ids []string
	for _, s := range args {
		_, err := hex.DecodeString(s)
		if err != nil {
			Warnf("argument %q is not a snapshot ID, ignoring\n", s)
			continue
		}

		ids = append(ids, s)
	}

	// process all snapshot IDs given as arguments
	for _, s := range ids {
		id, err := restic.FindSnapshot(repo, s)
		if err != nil {
			Warnf("could not find a snapshot for ID %q, ignoring\n", s)
			continue
		}

		if !opts.DryRun {
			h := restic.Handle{Type: restic.SnapshotFile, Name: id.String()}
			err = repo.Backend().Remove(h)
			if err != nil {
				return err
			}

			Verbosef("removed snapshot %v\n", id.Str())
		} else {
			Verbosef("would remove snapshot %v\n", id.Str())
		}
	}

	policy := restic.ExpirePolicy{
		Last:    opts.Last,
		Hourly:  opts.Hourly,
		Daily:   opts.Daily,
		Weekly:  opts.Weekly,
		Monthly: opts.Monthly,
		Yearly:  opts.Yearly,
		Tags:    opts.KeepTags,
	}

	if policy.Empty() {
		return nil
	}

	// then, load all remaining snapshots
	snapshots, err := restic.LoadAllSnapshots(repo)
	if err != nil {
		return err
	}

	// group by hostname and dirs
	type key struct {
		Hostname string
		Dirs     string
	}

	snapshotGroups := make(map[key]restic.Snapshots)

	for _, sn := range snapshots {
		if opts.Hostname != "" && sn.Hostname != opts.Hostname {
			continue
		}

		if !sn.HasTags(opts.Tags) {
			continue
		}

		k := key{Hostname: sn.Hostname, Dirs: strings.Join(sn.Paths, ":")}
		list := snapshotGroups[k]
		list = append(list, sn)
		snapshotGroups[k] = list
	}

	removeSnapshots := 0
	for key, snapshotGroup := range snapshotGroups {
		Printf("snapshots for host %v, directories %v:\n\n", key.Hostname, key.Dirs)
		keep, remove := restic.ApplyPolicy(snapshotGroup, policy)

		Printf("keep %d snapshots:\n", len(keep))
		printSnapshots(globalOptions.stdout, keep)
		Printf("\n")

		Printf("remove %d snapshots:\n", len(remove))
		printSnapshots(globalOptions.stdout, remove)
		Printf("\n")

		removeSnapshots += len(remove)

		if !opts.DryRun {
			for _, sn := range remove {
				h := restic.Handle{Type: restic.SnapshotFile, Name: sn.ID().String()}
				err = repo.Backend().Remove(h)
				if err != nil {
					return err
				}
			}
		}
	}

	if removeSnapshots > 0 && opts.Prune {
		Printf("%d snapshots have been removed, running prune\n", removeSnapshots)
		if !opts.DryRun {
			return pruneRepository(gopts, repo)
		}
	}

	return nil
}

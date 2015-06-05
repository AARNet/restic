package main

import (
	"bufio"
	"flag"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/restic/restic/backend"
	. "github.com/restic/restic/test"
)

var TestDataFile = flag.String("test.datafile", "", `specify tar.gz file with test data to backup and restore (required for integration test)`)

func setupTempdir(t testing.TB) (tempdir string) {
	tempdir, err := ioutil.TempDir(*TestTempDir, "restic-test-")
	OK(t, err)

	return tempdir
}

func configureRestic(t testing.TB, tempdir string) {
	// use cache dir within tempdir
	OK(t, os.Setenv("RESTIC_CACHE", filepath.Join(tempdir, "cache")))

	// configure environment
	opts.Repo = filepath.Join(tempdir, "repo")
	OK(t, os.Setenv("RESTIC_PASSWORD", *TestPassword))
}

func cleanupTempdir(t testing.TB, tempdir string) {
	if !*TestCleanup {
		t.Logf("leaving temporary directory %v used for test", tempdir)
		return
	}

	OK(t, os.RemoveAll(tempdir))
}

func setupTarTestFixture(t testing.TB, outputDir, tarFile string) {
	err := system("sh", "-c", `mkdir "$1" && (cd "$1" && tar xz) < "$2"`,
		"sh", outputDir, tarFile)
	OK(t, err)
}

func system(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func parseIDsFromReader(t testing.TB, rd io.Reader) backend.IDs {
	IDs := backend.IDs{}
	sc := bufio.NewScanner(rd)

	for sc.Scan() {
		id, err := backend.ParseID(sc.Text())
		if err != nil {
			t.Logf("parse id %v: %v", sc.Text(), err)
			continue
		}

		IDs = append(IDs, id)
	}

	return IDs
}

func cmdInit(t testing.TB) {
	cmd := &CmdInit{}
	OK(t, cmd.Execute(nil))

	t.Logf("repository initialized at %v", opts.Repo)
}

func cmdBackup(t testing.TB, target []string, parentID backend.ID) {
	cmd := &CmdBackup{}
	cmd.Parent = parentID.String()

	t.Logf("backing up %v", target)

	OK(t, cmd.Execute(target))
}

func cmdList(t testing.TB, tpe string) []backend.ID {

	rd, wr := io.Pipe()

	cmd := &CmdList{w: wr}

	go func() {
		OK(t, cmd.Execute([]string{tpe}))
		OK(t, wr.Close())
	}()

	IDs := parseIDsFromReader(t, rd)

	t.Logf("Listing %v: %v", tpe, IDs)

	return IDs
}

func cmdRestore(t testing.TB, dir string, snapshotID backend.ID) {
	cmd := &CmdRestore{}
	cmd.Execute([]string{snapshotID.String(), dir})
}

func cmdFsck(t testing.TB) {
	cmd := &CmdFsck{CheckData: true, Orphaned: true}
	OK(t, cmd.Execute(nil))
}

func TestBackup(t *testing.T) {
	if !*RunIntegrationTest {
		t.Skip("integration tests disabled, use `-test.integration` to enable")
	}

	if *TestDataFile == "" {
		t.Fatal("no data tar file specified, use flag `-test.datafile`")
	}

	tempdir := setupTempdir(t)
	defer cleanupTempdir(t, tempdir)

	configureRestic(t, tempdir)

	cmdInit(t)

	datadir := filepath.Join(tempdir, "testdata")

	setupTarTestFixture(t, datadir, *TestDataFile)

	// first backup
	cmdBackup(t, []string{datadir}, nil)
	snapshotIDs := cmdList(t, "snapshots")
	Assert(t, len(snapshotIDs) == 1,
		"more than one snapshot ID in repo")

	// second backup, implicit incremental
	cmdBackup(t, []string{datadir}, nil)
	snapshotIDs = cmdList(t, "snapshots")
	Assert(t, len(snapshotIDs) == 2,
		"more than one snapshot ID in repo")

	// third backup, explicit incremental
	cmdBackup(t, []string{datadir}, snapshotIDs[0])
	snapshotIDs = cmdList(t, "snapshots")
	Assert(t, len(snapshotIDs) == 3,
		"more than one snapshot ID in repo")

	// restore all backups and compare
	for _, snapshotID := range snapshotIDs {
		restoredir := filepath.Join(tempdir, "restore", snapshotID.String())
		t.Logf("restoring snapshot %v to %v", snapshotID.Str(), restoredir)
		cmdRestore(t, restoredir, snapshotIDs[0])
		Assert(t, directoriesEqualContents(datadir, filepath.Join(restoredir, "testdata")),
			"directories are not equal")
	}

	cmdFsck(t)
}

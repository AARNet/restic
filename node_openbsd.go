package restic

import (
	"os"
	"syscall"
)

func (node *Node) OpenForReading() (*os.File, error) {
	file, err := os.OpenFile(node.path, os.O_RDONLY, 0)
	if os.IsPermission(err) {
		return os.OpenFile(node.path, os.O_RDONLY, 0)
	}
	return file, err
}

func (node Node) restoreSymlinkTimestamps(path string, utimes [2]syscall.Timespec) error {
	return nil
}

func (s statUnix) atim() syscall.Timespec { return s.Atim }
func (s statUnix) mtim() syscall.Timespec { return s.Mtim }
func (s statUnix) ctim() syscall.Timespec { return s.Ctim }

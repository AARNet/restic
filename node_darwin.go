package restic

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"syscall"
	"time"

	"github.com/restic/restic/debug"
)

func (node *Node) OpenForReading() (*os.File, error) {
	return os.Open(node.path)
}

func (node *Node) fillExtra(path string, fi os.FileInfo) error {
	stat, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return nil
	}

	node.ChangeTime = time.Unix(stat.Ctimespec.Unix())
	node.AccessTime = time.Unix(stat.Atimespec.Unix())
	node.UID = stat.Uid
	node.GID = stat.Gid

	if u, err := user.LookupId(strconv.Itoa(int(stat.Uid))); err == nil {
		node.User = u.Username
	}

	node.Inode = stat.Ino

	var err error

	switch node.Type {
	case "file":
		node.Size = uint64(stat.Size)
		node.Links = uint64(stat.Nlink)
	case "dir":
	case "symlink":
		node.LinkTarget, err = os.Readlink(path)
	case "dev":
		node.Device = uint64(stat.Rdev)
	case "chardev":
		node.Device = uint64(stat.Rdev)
	case "fifo":
	case "socket":
	default:
		err = fmt.Errorf("invalid node type %q", node.Type)
	}

	return err
}

func (node *Node) createDevAt(path string) error {
	return syscall.Mknod(path, syscall.S_IFBLK|0600, int(node.Device))
}

func (node *Node) createCharDevAt(path string) error {
	return syscall.Mknod(path, syscall.S_IFCHR|0600, int(node.Device))
}

func (node *Node) createFifoAt(path string) error {
	return syscall.Mkfifo(path, 0600)
}

func changeTime(stat *syscall.Stat_t) time.Unix {
	return time.Unix(stat.Ctimespec.Unix())
}

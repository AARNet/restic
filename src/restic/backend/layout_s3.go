package backend

import "restic"

// S3Layout implements the old layout used for s3 cloud storage backends, as
// described in the Design document.
type S3Layout struct {
	URL  string
	Path string
	Join func(...string) string
}

var s3LayoutPaths = map[restic.FileType]string{
	restic.DataFile:     "data",
	restic.SnapshotFile: "snapshot",
	restic.IndexFile:    "index",
	restic.LockFile:     "lock",
	restic.KeyFile:      "key",
}

// join calls Join with the first empty elements removed.
func (l *S3Layout) join(url string, items ...string) string {
	for len(items) > 0 && items[0] == "" {
		items = items[1:]
	}

	path := l.Join(items...)
	if path == "" || path[0] != '/' {
		if url != "" && url[len(url)-1] != '/' {
			url += "/"
		}
	}

	return url + path
}

// Dirname returns the directory path for a given file type and name.
func (l *S3Layout) Dirname(h restic.Handle) string {
	if h.Type == restic.ConfigFile {
		return l.URL + l.Join(l.Path, "/")
	}

	return l.join(l.URL, l.Path, s3LayoutPaths[h.Type]) + "/"
}

// Filename returns a path to a file, including its name.
func (l *S3Layout) Filename(h restic.Handle) string {
	name := h.Name

	if h.Type == restic.ConfigFile {
		name = "config"
	}

	return l.join(l.URL, l.Path, s3LayoutPaths[h.Type], name)
}

// Paths returns all directory names
func (l *S3Layout) Paths() (dirs []string) {
	for _, p := range s3LayoutPaths {
		dirs = append(dirs, l.Join(l.Path, p))
	}
	return dirs
}

// Basedir returns the base dir name for type t.
func (l *S3Layout) Basedir(t restic.FileType) string {
	return l.Join(l.Path, s3LayoutPaths[t])
}

package backend

import "restic"

// RESTLayout implements the default layout for the REST protocol.
type RESTLayout struct {
	URL  string
	Path string
	Join func(...string) string
}

var restLayoutPaths = defaultLayoutPaths

// Dirname returns the directory path for a given file type and name.
func (l *RESTLayout) Dirname(h restic.Handle) string {
	if h.Type == restic.ConfigFile {
		return l.URL + l.Join(l.Path, "/")
	}

	return l.URL + l.Join(l.Path, "/", restLayoutPaths[h.Type]) + "/"
}

// Filename returns a path to a file, including its name.
func (l *RESTLayout) Filename(h restic.Handle) string {
	name := h.Name

	if h.Type == restic.ConfigFile {
		name = "config"
	}

	return l.URL + l.Join(l.Path, "/", restLayoutPaths[h.Type], name)
}

// Paths returns all directory names
func (l *RESTLayout) Paths() (dirs []string) {
	for _, p := range restLayoutPaths {
		dirs = append(dirs, l.URL+l.Join(l.Path, p))
	}
	return dirs
}

// Basedir returns the base dir name for files of type t.
func (l *RESTLayout) Basedir(t restic.FileType) string {
	return l.URL + l.Join(l.Path, restLayoutPaths[t])
}

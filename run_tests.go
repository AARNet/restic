// +build ignore

package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func specialDir(name string) bool {
	if name == "." {
		return false
	}

	base := filepath.Base(name)
	return base[0] == '_' || base[0] == '.'
}

func emptyDir(name string) bool {
	dir, err := os.Open(name)
	defer dir.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to open directory %v: %v\n", name, err)
		return true
	}

	fis, err := dir.Readdir(-1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Readdirnames(%v): %v\n", name, err)
		return true
	}

	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}

		if filepath.Ext(fi.Name()) == ".go" {
			return false
		}
	}

	return true
}

func forceRelativeDirname(dirname string) string {
	if dirname == "." {
		return dirname
	}

	if strings.HasPrefix(dirname, "./") {
		return dirname
	}

	return "./" + dirname
}

func mergeCoverprofile(file *os.File, out io.Writer) error {
	_, err := file.Seek(0, 0)
	if err != nil {
		return err
	}

	rd := bufio.NewReader(file)
	_, err = rd.ReadString('\n')
	if err == io.EOF {
		return nil
	}

	if err != nil {
		return err
	}

	_, err = io.Copy(out, rd)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return err
}

func testPackage(pkg string, out io.Writer) error {
	file, err := ioutil.TempFile("", "test-coverage-")
	defer os.Remove(file.Name())
	defer file.Close()
	if err != nil {
		return err
	}

	cmd := exec.Command("go", "test",
		"-cover", "-covermode", "set", "-coverprofile", file.Name(),
		pkg, "-test.integration")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		return err
	}

	return mergeCoverprofile(file, out)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "USAGE: run_tests COVERPROFILE [PATHS]")
		os.Exit(1)
	}

	target := os.Args[1]
	dirs := os.Args[2:]

	if len(dirs) == 0 {
		dirs = append(dirs, ".")
	}

	file, err := os.Create(target)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create coverprofile failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintln(file, "mode: set")

	for _, dir := range dirs {
		err := filepath.Walk(dir,
			func(p string, fi os.FileInfo, e error) error {
				if e != nil {
					return e
				}

				if !fi.IsDir() {
					return nil
				}

				if specialDir(p) || emptyDir(p) {
					return filepath.SkipDir
				}

				return testPackage(forceRelativeDirname(p), file)
			})

		if err != nil {
			fmt.Fprintf(os.Stderr, "walk(%q): %v\n", dir, err)
		}
	}

	err = file.Close()

	fmt.Printf("coverprofile: %v\n", file.Name())
}

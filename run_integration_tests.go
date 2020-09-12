// +build ignore

package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// ForbiddenImports are the packages from the stdlib that should not be used in
// our code.
var ForbiddenImports = map[string]bool{
	"errors": true,
}

// Use a specific version of gofmt (the latest stable, usually) to guarantee
// deterministic formatting. This is used with the GoVersion.AtLeast()
// function (so that we don't forget to update it). This is also used to run
// `go mod tidy`.
var GofmtVersion = ParseGoVersion("go1.14")

// GoVersion is the version of Go used to compile the project.
type GoVersion struct {
	Major int
	Minor int
	Patch int
}

// ParseGoVersion parses the Go version s. If s cannot be parsed, the returned GoVersion is null.
func ParseGoVersion(s string) (v GoVersion) {
	if !strings.HasPrefix(s, "go") {
		return
	}

	s = s[2:]
	data := strings.Split(s, ".")
	if len(data) < 2 || len(data) > 3 {
		// invalid version
		return GoVersion{}
	}

	var err error

	v.Major, err = strconv.Atoi(data[0])
	if err != nil {
		return GoVersion{}
	}

	// try to parse the minor version while removing an eventual suffix (like
	// "rc2" or so)
	for s := data[1]; s != ""; s = s[:len(s)-1] {
		v.Minor, err = strconv.Atoi(s)
		if err == nil {
			break
		}
	}

	if v.Minor == 0 {
		// no minor version found
		return GoVersion{}
	}

	if len(data) >= 3 {
		v.Patch, err = strconv.Atoi(data[2])
		if err != nil {
			return GoVersion{}
		}
	}

	return
}

// AtLeast returns true if v is at least as new as other. If v is empty, true is returned.
func (v GoVersion) AtLeast(other GoVersion) bool {
	var empty GoVersion

	// the empty version satisfies all versions
	if v == empty {
		return true
	}

	if v.Major < other.Major {
		return false
	}

	if v.Minor < other.Minor {
		return false
	}

	if v.Patch < other.Patch {
		return false
	}

	return true
}

func (v GoVersion) String() string {
	return fmt.Sprintf("Go %d.%d.%d", v.Major, v.Minor, v.Patch)
}

// CloudBackends contains a map of backend tests for cloud services to one
// of the essential environment variables which must be present in order to
// test it.
var CloudBackends = map[string]string{
	"restic/backend/s3.TestBackendS3":       "RESTIC_TEST_S3_REPOSITORY",
	"restic/backend/swift.TestBackendSwift": "RESTIC_TEST_SWIFT",
	"restic/backend/b2.TestBackendB2":       "RESTIC_TEST_B2_REPOSITORY",
	"restic/backend/gs.TestBackendGS":       "RESTIC_TEST_GS_REPOSITORY",
	"restic/backend/azure.TestBackendAzure": "RESTIC_TEST_AZURE_REPOSITORY",
}

var runCrossCompile = flag.Bool("cross-compile", true, "run cross compilation tests")

func init() {
	flag.Parse()
}

// CIEnvironment is implemented by environments where tests can be run.
type CIEnvironment interface {
	Prepare() error
	RunTests() error
	Teardown() error
}

// TravisEnvironment is the environment in which Travis tests run.
type TravisEnvironment struct {
	goxOSArch          []string
	env                map[string]string
	gcsCredentialsFile string
}

func (env *TravisEnvironment) getMinio() error {
	tempfile, err := os.Create(filepath.Join(os.Getenv("GOPATH"), "bin", "minio"))
	if err != nil {
		return fmt.Errorf("create tempfile for minio download failed: %v", err)
	}

	url := fmt.Sprintf("https://dl.minio.io/server/minio/release/%s-%s/minio",
		runtime.GOOS, runtime.GOARCH)
	msg("downloading %v\n", url)
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error downloading minio server: %v", err)
	}

	_, err = io.Copy(tempfile, res.Body)
	if err != nil {
		return fmt.Errorf("error saving minio server to file: %v", err)
	}

	err = res.Body.Close()
	if err != nil {
		return fmt.Errorf("error closing HTTP download: %v", err)
	}

	err = tempfile.Close()
	if err != nil {
		msg("closing tempfile failed: %v\n", err)
		return fmt.Errorf("error closing minio server file: %v", err)
	}

	err = os.Chmod(tempfile.Name(), 0755)
	if err != nil {
		return fmt.Errorf("chmod(minio-server) failed: %v", err)
	}

	msg("downloaded minio server to %v\n", tempfile.Name())
	return nil
}

// Prepare installs dependencies and starts services in order to run the tests.
func (env *TravisEnvironment) Prepare() error {
	env.env = make(map[string]string)

	msg("preparing environment for Travis CI\n")

	pkgs := []string{
		"github.com/NebulousLabs/glyphcheck",
		"github.com/restic/rest-server/cmd/rest-server",
		"github.com/restic/calens",
		"github.com/rclone/rclone",
	}

	for _, pkg := range pkgs {
		err := run("go", "get", pkg)
		if err != nil {
			return err
		}
	}

	// reset changes made to go.mod/go.sum by "go get"
	if err := run("git", "checkout", "go.mod", "go.sum"); err != nil {
		return err
	}

	if err := env.getMinio(); err != nil {
		return err
	}

	if *runCrossCompile {
		// only test cross compilation on linux with Travis
		if err := run("go", "get", "github.com/mitchellh/gox"); err != nil {
			return err
		}

		// reset changes made to go.mod/go.sum by "go get"
		if err := run("git", "checkout", "go.mod", "go.sum"); err != nil {
			return err
		}

		if runtime.GOOS == "linux" {
			env.goxOSArch = []string{
				"linux/386", "linux/amd64",
				"windows/386", "windows/amd64",
				"darwin/amd64",
				"freebsd/386", "freebsd/amd64",
				"openbsd/386", "openbsd/amd64",
				"netbsd/386", "netbsd/amd64",
				"linux/arm", "freebsd/arm",
				"linux/ppc64le", "solaris/amd64",
			}
		} else {
			env.goxOSArch = []string{runtime.GOOS + "/" + runtime.GOARCH}
		}

		msg("gox: OS/ARCH %v\n", env.goxOSArch)
	}

	// do not run cloud tests on darwin
	if os.Getenv("RESTIC_TEST_CLOUD_BACKENDS") == "0" {
		msg("skipping cloud backend tests\n")

		for _, name := range CloudBackends {
			err := os.Unsetenv(name)
			if err != nil {
				msg("    error unsetting %v: %v\n", name, err)
			}
		}
	}

	// extract credentials file for GCS tests
	if b64data := os.Getenv("RESTIC_TEST_GS_APPLICATION_CREDENTIALS_B64"); b64data != "" {
		buf, err := base64.StdEncoding.DecodeString(b64data)
		if err != nil {
			return err
		}

		f, err := ioutil.TempFile("", "gcs-credentials-")
		if err != nil {
			return err
		}

		msg("saving GCS credentials to %v\n", f.Name())

		_, err = f.Write(buf)
		if err != nil {
			f.Close()
			return err
		}

		env.gcsCredentialsFile = f.Name()

		if err = f.Close(); err != nil {
			return err
		}
	}

	return nil
}

// Teardown stops backend services and cleans the environment again.
func (env *TravisEnvironment) Teardown() error {
	msg("run travis teardown\n")

	if env.gcsCredentialsFile != "" {
		msg("remove gcs credentials file %v\n", env.gcsCredentialsFile)
		return os.Remove(env.gcsCredentialsFile)
	}

	return nil
}

// RunTests starts the tests for Travis.
func (env *TravisEnvironment) RunTests() error {
	env.env["GOPATH"] = os.Getenv("GOPATH")
	if env.gcsCredentialsFile != "" {
		env.env["GOOGLE_APPLICATION_CREDENTIALS"] = env.gcsCredentialsFile
	}

	// ensure that the following tests cannot be silently skipped on Travis
	ensureTests := []string{
		"restic/backend/rest.TestBackendREST",
		"restic/backend/sftp.TestBackendSFTP",
		"restic/backend/s3.TestBackendMinio",
		"restic/backend/rclone.TestBackendRclone",
	}

	// make sure that cloud backends for which we have credentials are not
	// silently skipped.
	for pkg, env := range CloudBackends {
		if _, ok := os.LookupEnv(env); ok {
			ensureTests = append(ensureTests, pkg)
		} else {
			msg("credentials for %v are not available, skipping\n", pkg)
		}
	}

	env.env["RESTIC_TEST_DISALLOW_SKIP"] = strings.Join(ensureTests, ",")

	if *runCrossCompile {
		// compile for all target architectures with tags
		for _, tags := range []string{"", "debug"} {
			err := runWithEnv(env.env, "gox", "-verbose",
				"-osarch", strings.Join(env.goxOSArch, " "),
				"-tags", tags,
				"-output", "/tmp/{{.Dir}}_{{.OS}}_{{.Arch}}",
				"./cmd/restic")
			if err != nil {
				return err
			}
		}
	}

	v := ParseGoVersion(runtime.Version())
	msg("Detected Go version %v\n", v)

	args := []string{"go", "run", "build.go"}

	// run the build script
	err := run(args[0], args[1:]...)
	if err != nil {
		return err
	}

	// run the tests and gather coverage information
	err = runWithEnv(env.env, "go", "test", "-count", "1", "-coverprofile", "all.cov", "./...")
	if err != nil {
		return err
	}

	// only run gofmt on a specific version of Go.
	if v.AtLeast(GofmtVersion) {
		if err = runGofmt(); err != nil {
			return err
		}

		msg("run go mod tidy\n")
		if err := runGoModTidy(); err != nil {
			return err
		}
	} else {
		msg("Skipping gofmt and mod tidy check for %v\n", v)
	}

	if err = runGlyphcheck(); err != nil {
		return err
	}

	// check for forbidden imports
	deps, err := env.findImports()
	if err != nil {
		return err
	}

	foundForbiddenImports := false
	for name, imports := range deps {
		for _, pkg := range imports {
			if _, ok := ForbiddenImports[pkg]; ok {
				fmt.Fprintf(os.Stderr, "========== package %v imports forbidden package %v\n", name, pkg)
				foundForbiddenImports = true
			}
		}
	}

	if foundForbiddenImports {
		return errors.New("CI: forbidden imports found")
	}

	// check that the entries in changelog/ are valid
	if err := run("calens"); err != nil {
		return errors.New("calens failed, files in changelog/ are not valid")
	}

	return nil
}

// AppveyorEnvironment is the environment on Windows.
type AppveyorEnvironment struct{}

// Prepare installs dependencies and starts services in order to run the tests.
func (env *AppveyorEnvironment) Prepare() error {
	return nil
}

// RunTests start the tests.
func (env *AppveyorEnvironment) RunTests() error {
	return runWithEnv(nil, "go", "run", "build.go", "-v", "-T")
}

// Teardown is a noop.
func (env *AppveyorEnvironment) Teardown() error {
	return nil
}

func msg(format string, args ...interface{}) {
	fmt.Printf("CI: "+format, args...)
}

func updateEnv(env []string, override map[string]string) []string {
	var newEnv []string
	for _, s := range env {
		d := strings.SplitN(s, "=", 2)
		key := d[0]

		if _, ok := override[key]; ok {
			continue
		}

		newEnv = append(newEnv, s)
	}

	for k, v := range override {
		newEnv = append(newEnv, k+"="+v)
	}

	return newEnv
}

func (env *TravisEnvironment) findImports() (map[string][]string, error) {
	res := make(map[string][]string)
	msg("checking for forbidden imports\n")

	cmd := exec.Command("go", "list", "-f", `{{.ImportPath}} {{join .Imports " "}}`, "./internal/...", "./cmd/...")
	cmd.Stderr = os.Stderr

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	sc := bufio.NewScanner(bytes.NewReader(output))
	for sc.Scan() {
		wordScanner := bufio.NewScanner(strings.NewReader(sc.Text()))
		wordScanner.Split(bufio.ScanWords)

		if !wordScanner.Scan() {
			return nil, fmt.Errorf("package name not found in line: %s", output)
		}
		name := wordScanner.Text()
		var deps []string

		for wordScanner.Scan() {
			deps = append(deps, wordScanner.Text())
		}

		res[name] = deps
	}

	return res, nil
}

func runGofmt() error {
	cmd := exec.Command("gofmt", "-l", ".")
	cmd.Stderr = os.Stderr

	buf, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("error running gofmt: %v\noutput: %s", err, buf)
	}

	if len(buf) > 0 {
		return fmt.Errorf("not formatted with `gofmt`:\n%s", buf)
	}

	return nil
}

// run "go mod tidy" so that go.sum and go.mod are updated to reflect all
// dependencies for all OS/Arch combinations, see
// https://github.com/golang/go/wiki/Modules#why-does-go-mod-tidy-put-so-many-indirect-dependencies-in-my-gomod
func runGoModTidy() error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Env = updateEnv(os.Environ(), nil)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error running 'go mod tidy': %v", err)
	}

	// check that "git diff" does not return any output
	cmd = exec.Command("git", "diff", "go.sum", "go.mod")
	cmd.Stderr = os.Stderr

	buf, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("error running 'git diff': %v\noutput: %s", err, buf)
	}

	if len(buf) > 0 {
		return fmt.Errorf("`go.mod` or `go.sum` not up to date (forgot to run `go mod tidy`?):\n%s", buf)
	}

	return nil
}

func runGlyphcheck() error {
	cmd := exec.Command("glyphcheck", "./cmd/...", "./internal/...")
	cmd.Stderr = os.Stderr

	buf, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("error running glyphcheck: %v\noutput: %s", err, buf)
	}

	return nil
}

func run(command string, args ...string) error {
	msg("run %v %v\n", command, strings.Join(args, " "))
	return runWithEnv(nil, command, args...)
}

// runWithEnv calls a command with the current environment, except the entries
// of the env map are set additionally.
func runWithEnv(env map[string]string, command string, args ...string) error {
	msg("runWithEnv %v %v\n", command, strings.Join(args, " "))
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if env != nil {
		cmd.Env = updateEnv(os.Environ(), env)
	}
	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("error running %v %v: %v",
			command, strings.Join(args, " "), err)
	}
	return nil
}

func isTravis() bool {
	return os.Getenv("TRAVIS_BUILD_DIR") != ""
}

func isAppveyor() bool {
	return runtime.GOOS == "windows"
}

func main() {
	// make sure we run in Module mode
	err := os.Setenv("GO111MODULE", "on")
	if err != nil {
		msg("setenv(GO111MODULE=on) return error: %v\n", err)
		os.Exit(1)
	}

	// enable the Go Module Proxy
	err = os.Setenv("GOPROXY", "https://proxy.golang.org")
	if err != nil {
		msg("setenv(GOPROXY) return error: %v\n", err)
		os.Exit(1)
	}

	var env CIEnvironment

	switch {
	case isTravis():
		env = &TravisEnvironment{}
	case isAppveyor():
		env = &AppveyorEnvironment{}
	default:
		fmt.Fprintln(os.Stderr, "unknown CI environment")
		os.Exit(1)
	}

	err = env.Prepare()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error preparing: %v\n", err)
		os.Exit(1)
	}

	err = env.RunTests()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error running tests: %v\n", err)
		os.Exit(2)
	}

	err = env.Teardown()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error during teardown: %v\n", err)
		os.Exit(3)
	}
}

# AARNet Restic Build

## Synopsis

The AARNet Restic Build is intended to add features that we require that are not currently included in the official restic build.

This build includes:
* `--dry-run` for listing packs from a snapshot without restoring them
* `--min-packsize` for configuring the minimum size of a pack.
* `--file-read-concurrency` for configuring number of readers when performing a backup

## Resources

| Name | Info |
| -------- | -------- |
| Restic Git | https://github.com/restic/restic/ |
| Restic Documentation | https://restic.readthedocs.io/en/latest/ |
| Jenkins Build | https://jenkins.aarnet.net.au/job/CloudStor/job/restic/ |
| `--dry-run` source | https://github.com/beer4duke/restic/tree/b12db3c5180352bb3fe739059ddd6a4c32f0aab4 |
| `--min-packsize` and `--file-read-concurrency` source | https://github.com/metalsp0rk/restic/commit/74d19cf2b8be9372860ea74b94be95b4f59e115a |


## Configuration

See [README.rst](README.rst) for the official README for restic.

## Development

### Requirements
- golang 1.15
- fuse (for running tests only)
- make (only required if you use the Makefile)

`go run build.go --test` builds the binary and runs any tests.

`go run build.go` just builds the binary

`go run build.go --help` for more options.










git checkout -b my-first-merge
git pull https://github.com/restic/restic.git v0.10.0


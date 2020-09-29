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


| `--dry-run` source | https://github.com/beer4duke/restic/tree/b12db3c5180352bb3fe739059ddd6a4c32f0aab4 |
| `--min-packsize` and `--file-read-concurrency` source | https://github.com/metalsp0rk/restic/commit/74d19cf2b8be9372860ea74b94be95b4f59e115a |


## Configuration

See [README.rst](README.rst) for the official README for restic.

## Docker Image

There is a Dockerfile included that will build an image to perform the initialisation of a repository (if necessary) and backup to it.

### How to build and run

#### Building the docker iamge

A new build of the restic docker image will be built from master when the restic pipeline completes successfully.

If you need to manually build one, you can just run:


docker build -t restic-build:v0.10.0 .


#### Example running with S3


docker run -e AWS_ACCESS_KEY_ID=minioadmin \
-e AWS_SECRET_ACCESS_KEY=minioadmin \
-e RESTIC_PASSWORD=super_secret_repository_password \
-e RESTIC_REPOSITORY=s3:https://minio.s3.aarnet.edu.au/test_restic_repo/ \
-v `pwd`:/backup -v `pwd`/logs:/var/log/restic restic-build:v0.10.0


#### Example running with WebDAV


docker run -e RCLONE_WEBDAV_USER=webdav_user \
-e RCLONE_WEBDAV_PASS=super_secret \
-e RCLONE_WEBDAV_URL=https://cloudstor-uat.aarnet.edu.au/plus/remote.php/webdav \
-e RESTIC_PASSWORD=super_secret_repository_password \
-e RESTIC_REPOSITORY=rclone::webdav:/test_restic_repo \
-v `pwd`:/backup -v `pwd`/logs:/var/log/restic restic-build:v0.10.0


### Volumes

* `/backup` - The default directory to backup
* `/var/log/restic` - (optional) The log directory. It will contain a file called `restic.log` and `last_status`. `last_status` will container a status and message of the last run.

### Environment Variables

#### Restic
| Variable | Default | Purpose |
| -------- | -------- | -------- |
| RESTIC_REPOSITORY | `none` | The restic repository to use in the format of `s3:https://endpoint.s3.aarnet.edu.au/bucket_name` or `rclone::webdav:/path/for/repository` |
| RESTIC_PASSWORD | `none` | The password for the restic repository |
| RESTIC_BACKUP_SOURCE | `/backup` | The directory to backup |
| RESTIC_LOG_VERBOSITY | `0` | The verbosity of the logs. The higher the number, the more verbose the logs. I have not tried setting it higher than `3`.

#### S3 specific

| Variable | Default | Purpose |
| -------- | -------- | -------- |
| AWS_ACCESS_KEY_ID | `none` | Access key for the S3 endpoint |
| AWS_SECRET_ACCESS_KEY | `none` | Secret key for the S3 endpoint |

#### WebDAV specific

As restic uses rclone for WebDAV, the configuration is done using the rclone environment variables.

| Variable | Default | Purpose |
| -------- | -------- | -------- |
| RCLONE_WEBDAV_URL | `none` | The URL to the WebDAV endpoint (eg. `https://cloudstor.aarnet.edu.au/plus/remote.php/webdav`)|
| RCLONE_WEBDAV_USER | `none` | The username for the webdav endpoint |
| RCLONE_WEBDAV_PASS | `none` | The password for the webdav endpoint |
| RCLONE_WEBDAV_VENDOR | `owncloud` | The vendor for the WebDAV endpoint (can be `owncloud`, `nextcloud` or `other`) |

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


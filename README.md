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
| RESTIC_LOG_VERBOSITY | `0` | The verbosity of the logs. The higher the number, the more verbose the logs. I have not tried setting it higher than `3`. |
| RESTIC_FILE_READ_CONCURRENCY | `2` | The number of files to read concurrently. This is used with `restic backup`. |
| RESTIC_MIN_PACKSIZE | `4` | The minimum pack size. |

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
















### Updating rclone
1. Open `Dockerfile`
2. Replace the download URL for rclone with the new version's URL (or just modify the version numbers)
3. Rebuild the docker iamge

---

[![Documentation](https://readthedocs.org/projects/restic/badge/?version=latest)](https://restic.readthedocs.io/en/latest/?badge=latest)
[![Build Status](https://github.com/restic/restic/workflows/test/badge.svg)](https://github.com/restic/restic/actions?query=workflow%3Atest)
[![Go Report Card](https://goreportcard.com/badge/github.com/restic/restic)](https://goreportcard.com/report/github.com/restic/restic)

# Introduction

restic is a backup program that is fast, efficient and secure. It supports the three major operating systems (Linux, macOS, Windows) and a few smaller ones (FreeBSD, OpenBSD).

For detailed usage and installation instructions check out the [documentation](https://restic.readthedocs.io/en/latest).

You can ask questions in our [Discourse forum](https://forum.restic.net).

Quick start
-----------

Once you've [installed](https://restic.readthedocs.io/en/latest/020_installation.html) restic, start
off with creating a repository for your backups:

    $ restic init --repo /tmp/backup
    enter password for new backend:
    enter password again:
    created restic backend 085b3c76b9 at /tmp/backup
    Please note that knowledge of your password is required to access the repository.
    Losing your password means that your data is irrecoverably lost.

and add some data:

    $ restic --repo /tmp/backup backup ~/work
    enter password for repository:
    scan [/home/user/work]
    scanned 764 directories, 1816 files in 0:00
    [0:29] 100.00%  54.732 MiB/s  1.582 GiB / 1.582 GiB  2580 / 2580 items  0 errors  ETA 0:00
    duration: 0:29, 54.47MiB/s
    snapshot 40dc1520 saved

Next you can either use `restic restore` to restore files or use `restic
mount` to mount the repository via fuse and browse the files from previous
snapshots.

For more options check out the [online documentation](https://restic.readthedocs.io/en/latest/).

# Backends

Saving a backup on the same machine is nice but not a real backup strategy.
Therefore, restic supports the following backends for storing backups natively:

- [Local directory](https://restic.readthedocs.io/en/latest/030_preparing_a_new_repo.html#local)
- [sftp server (via SSH)](https://restic.readthedocs.io/en/latest/030_preparing_a_new_repo.html#sftp)
- [HTTP REST server](https://restic.readthedocs.io/en/latest/030_preparing_a_new_repo.html#rest-server) ([protocol](doc/100_references.rst#rest-backend), [rest-server](https://github.com/restic/rest-server))
- [AWS S3](https://restic.readthedocs.io/en/latest/030_preparing_a_new_repo.html#amazon-s3) (either from Amazon or using the [Minio](https://minio.io) server)
- [OpenStack Swift](https://restic.readthedocs.io/en/latest/030_preparing_a_new_repo.html#openstack-swift)
- [BackBlaze B2](https://restic.readthedocs.io/en/latest/030_preparing_a_new_repo.html#backblaze-b2)
- [Microsoft Azure Blob Storage](https://restic.readthedocs.io/en/latest/030_preparing_a_new_repo.html#microsoft-azure-blob-storage)
- [Google Cloud Storage](https://restic.readthedocs.io/en/latest/030_preparing_a_new_repo.html#google-cloud-storage)
- And many other services via the [rclone](https://rclone.org) [Backend](https://restic.readthedocs.io/en/latest/030_preparing_a_new_repo.html#other-services-via-rclone)

# Design Principles

Restic is a program that does backups right and was designed with the
following principles in mind:

-  **Easy:** Doing backups should be a frictionless process, otherwise
   you might be tempted to skip it. Restic should be easy to configure
   and use, so that, in the event of a data loss, you can just restore
   it. Likewise, restoring data should not be complicated.

-  **Fast**: Backing up your data with restic should only be limited by
   your network or hard disk bandwidth so that you can backup your files
   every day. Nobody does backups if it takes too much time. Restoring
   backups should only transfer data that is needed for the files that
   are to be restored, so that this process is also fast.

-  **Verifiable**: Much more important than backup is restore, so restic
   enables you to easily verify that all data can be restored.

-  **Secure**: Restic uses cryptography to guarantee confidentiality and
   integrity of your data. The location the backup data is stored is
   assumed not to be a trusted environment (e.g. a shared space where
   others like system administrators are able to access your backups).
   Restic is built to secure your data against such attackers.

-  **Efficient**: With the growth of data, additional snapshots should
   only take the storage of the actual increment. Even more, duplicate
   data should be de-duplicated before it is actually written to the
   storage back end to save precious backup space.

# Reproducible Builds

The binaries released with each restic version starting at 0.6.1 are
[reproducible](https://reproducible-builds.org/), which means that you can
reproduce a byte identical version from the source code for that
release. Instructions on how to do that are contained in the
[builder repository](https://github.com/restic/builder).

News
----

You can follow the restic project on Twitter [@resticbackup](https://twitter.com/resticbackup) or by subscribing to
the [project blog](https://restic.net/blog/).

License
-------

Restic is licensed under [BSD 2-Clause License](https://opensource.org/licenses/BSD-2-Clause). You can find the
complete text in [``LICENSE``](LICENSE).

Sponsorship
-----------

Backend integration tests for Google Cloud Storage and Microsoft Azure Blob
Storage are sponsored by [AppsCode](https://appscode.com)!

[![Sponsored by AppsCode](https://cdn.appscode.com/images/logo/appscode/ac-logo-color.png)](https://appscode.com)

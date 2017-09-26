Manual
======

Usage help
----------

Usage help is available:

.. code-block:: console

    $ ./restic --help
    restic is a backup program which allows saving multiple revisions of files and
    directories in an encrypted repository stored on different backends.

    Usage:
      restic [command]

    Available Commands:
      autocomplete  Generate shell autocompletion script
      backup        Create a new backup of files and/or directories
      cat           Print internal objects to stdout
      check         Check the repository for errors
      dump          Dump data structures
      find          Find a file or directory
      forget        Remove snapshots from the repository
      help          Help about any command
      init          Initialize a new repository
      key           Manage keys (passwords)
      list          List items in the repository
      ls            List files in a snapshot
      mount         Mount the repository
      prune         Remove unneeded data from the repository
      rebuild-index Build a new index file
      restore       Extract the data from a snapshot
      snapshots     List all snapshots
      tag           Modify tags on snapshots
      unlock        Remove locks other processes created
      version       Print version information

    Flags:
          --json                   set output mode to JSON for commands that support it
          --no-lock                do not lock the repo, this allows some operations on read-only repos
      -p, --password-file string   read the repository password from a file
      -q, --quiet                  do not output comprehensive progress report
      -r, --repo string            repository to backup to or restore from (default: $RESTIC_REPOSITORY)

    Use "restic [command] --help" for more information about a command.

Similar to programs such as ``git``, restic has a number of
sub-commands. You can see these commands in the listing above. Each
sub-command may have own command-line options, and there is a help
option for each command which lists them, e.g. for the ``backup``
command:

.. code-block:: console

    $ ./restic backup --help
    The "backup" command creates a new snapshot and saves the files and directories
    given as the arguments.

    Usage:
      restic backup [flags] FILE/DIR [FILE/DIR] ...

    Flags:
      -e, --exclude pattern         exclude a pattern (can be specified multiple times)
          --exclude-file string     read exclude patterns from a file
          --files-from string       read the files to backup from file (can be combined with file args)
      -f, --force                   force re-reading the target files/directories. Overrides the "parent" flag
      -x, --one-file-system         Exclude other file systems
          --parent string           use this parent snapshot (default: last snapshot in the repo that has the same target files/directories)
          --stdin                   read backup from stdin
          --stdin-filename string   file name to use when reading from stdin
          --tag tag                 add a tag for the new snapshot (can be specified multiple times)
          --time string             time of the backup (ex. '2012-11-01 22:08:41') (default: now)

    Global Flags:
          --json                   set output mode to JSON for commands that support it
          --no-lock                do not lock the repo, this allows some operations on read-only repos
      -p, --password-file string   read the repository password from a file
      -q, --quiet                  do not output comprehensive progress report
      -r, --repo string            repository to backup to or restore from (default: $RESTIC_REPOSITORY)

Subcommand that support showing progress information such as ``backup``,
``check`` and ``prune`` will do so unless the quiet flag ``-q`` or
``--quiet`` is set. When running from a non-interactive console progress
reporting will be limited to once every 10 seconds to not fill your
logs.

Additionally on Unix systems if ``restic`` receives a SIGUSR signal the
current progress will written to the standard output so you can check up
on the status at will.

Initialize a repository
-----------------------

First, we need to create a "repository". This is the place where your
backups will be saved at.

Local
~~~~~

In order to create a repository at ``/tmp/backup``, run the following
command and enter the same password twice:

.. code-block:: console

    $ restic init --repo /tmp/backup
    enter password for new backend:
    enter password again:
    created restic backend 085b3c76b9 at /tmp/backup
    Please note that knowledge of your password is required to access the repository.
    Losing your password means that your data is irrecoverably lost.

Other backends like sftp and s3 are described in the following
sections.

Remembering your password is important! If you lose it, you won't be
able to access data stored in the repository.

For automated backups, restic accepts the repository location in the
environment variable ``RESTIC_REPOSITORY``. The password can be read
from a file (via the option ``--password-file`` or the environment variable
``RESTIC_PASSWORD_FILE``) or the environment variable ``RESTIC_PASSWORD``.

SFTP
~~~~

In order to backup data via SFTP, you must first set up a server with
SSH and let it know your public key. Passwordless login is really
important since restic fails to connect to the repository if the server
prompts for credentials.

Once the server is configured, the setup of the SFTP repository can
simply be achieved by changing the URL scheme in the ``init`` command:

.. code-block:: console

    $ restic -r sftp:user@host:/tmp/backup init
    enter password for new backend:
    enter password again:
    created restic backend f1c6108821 at sftp:user@host:/tmp/backup
    Please note that knowledge of your password is required to access the repository.
    Losing your password means that your data is irrecoverably lost.

You can also specify a relative (read: no slash (``/``) character at the
beginning) directory, in this case the dir is relative to the remote
user's home directory.

The backend config string does not allow specifying a port. If you need
to contact an sftp server on a different port, you can create an entry
in the ``ssh`` file, usually located in your user's home directory at
``~/.ssh/config`` or in ``/etc/ssh/ssh_config``:

::

    Host foo
        User bar
        Port 2222

Then use the specified host name ``foo`` normally (you don't need to
specify the user name in this case):

::

    $ restic -r sftp:foo:/tmp/backup init

You can also add an entry with a special host name which does not exist,
just for use with restic, and use the ``Hostname`` option to set the
real host name:

::

    Host restic-backup-host
        Hostname foo
        User bar
        Port 2222

Then use it in the backend specification:

::

    $ restic -r sftp:restic-backup-host:/tmp/backup init

Last, if you'd like to use an entirely different program to create the
SFTP connection, you can specify the command to be run with the option
``-o sftp.command="foobar"``.

REST Server
~~~~~~~~~~~

In order to backup data to the remote server via HTTP or HTTPS protocol,
you must first set up a remote `REST
server <https://github.com/restic/rest-server>`__ instance. Once the
server is configured, accessing it is achieved by changing the URL
scheme like this:

.. code-block:: console

    $ restic -r rest:http://host:8000/

Depending on your REST server setup, you can use HTTPS protocol,
password protection, or multiple repositories. Or any combination of
those features, as you see fit. TCP/IP port is also configurable. Here
are some more examples:

.. code-block:: console

    $ restic -r rest:https://host:8000/
    $ restic -r rest:https://user:pass@host:8000/
    $ restic -r rest:https://user:pass@host:8000/my_backup_repo/

If you use TLS, make sure your certificates are signed, 'cause restic
client will refuse to communicate otherwise. It's easy to obtain such
certificates today, thanks to free certificate authorities like `Let’s
Encrypt <https://letsencrypt.org/>`__.

REST server uses exactly the same directory structure as local backend,
so you should be able to access it both locally and via HTTP, even
simultaneously.

Amazon S3
~~~~~~~~~

Restic can backup data to any Amazon S3 bucket. However, in this case,
changing the URL scheme is not enough since Amazon uses special security
credentials to sign HTTP requests. By consequence, you must first setup
the following environment variables with the credentials you obtained
while creating the bucket.

.. code-block:: console

    $ export AWS_ACCESS_KEY_ID=<MY_ACCESS_KEY>
    $ export AWS_SECRET_ACCESS_KEY=<MY_SECRET_ACCESS_KEY>

You can then easily initialize a repository that uses your Amazon S3 as
a backend, if the bucket does not exist yet it will be created in the
default location:

.. code-block:: console

    $ restic -r s3:s3.amazonaws.com/bucket_name init
    enter password for new backend:
    enter password again:
    created restic backend eefee03bbd at s3:s3.amazonaws.com/bucket_name
    Please note that knowledge of your password is required to access the repository.
    Losing your password means that your data is irrecoverably lost.

It is not possible at the moment to have restic create a new bucket in a
different location, so you need to create it using a different program.
Afterwards, the S3 server (``s3.amazonaws.com``) will redirect restic to
the correct endpoint.

For an S3-compatible server that is not Amazon (like Minio, see below),
or is only available via HTTP, you can specify the URL to the server
like this: ``s3:http://server:port/bucket_name``.

Minio Server
~~~~~~~~~~~~

`Minio <https://www.minio.io>`__ is an Open Source Object Storage,
written in Go and compatible with AWS S3 API.

-  Download and Install `Minio
   Server <https://minio.io/downloads/#minio-server>`__.
-  You can also refer to https://docs.minio.io for step by step guidance
   on installation and getting started on Minio Client and Minio Server.

You must first setup the following environment variables with the
credentials of your running Minio Server.

.. code-block:: console

    $ export AWS_ACCESS_KEY_ID=<YOUR-MINIO-ACCESS-KEY-ID>
    $ export AWS_SECRET_ACCESS_KEY= <YOUR-MINIO-SECRET-ACCESS-KEY>

Now you can easily initialize restic to use Minio server as backend with
this command.

.. code-block:: console

    $ ./restic -r s3:http://localhost:9000/restic init
    enter password for new backend:
    enter password again:
    created restic backend 6ad29560f5 at s3:http://localhost:9000/restic1
    Please note that knowledge of your password is required to access
    the repository. Losing your password means that your data is irrecoverably lost.

OpenStack Swift
~~~~~~~~~~~~~~~

Restic can backup data to an OpenStack Swift container. Because Swift supports
various authentication methods, credentials are passed through environment
variables. In order to help integration with existing OpenStack installations,
the naming convention of those variables follows official python swift client:

.. code-block:: console

   # For keystone v1 authentication
   $ export ST_AUTH=<MY_AUTH_URL>
   $ export ST_USER=<MY_USER_NAME>
   $ export ST_KEY=<MY_USER_PASSWORD>

   # For keystone v2 authentication (some variables are optional)
   $ export OS_AUTH_URL=<MY_AUTH_URL>
   $ export OS_REGION_NAME=<MY_REGION_NAME>
   $ export OS_USERNAME=<MY_USERNAME>
   $ export OS_PASSWORD=<MY_PASSWORD>
   $ export OS_TENANT_ID=<MY_TENANT_ID>
   $ export OS_TENANT_NAME=<MY_TENANT_NAME>

   # For keystone v3 authentication (some variables are optional)
   $ export OS_AUTH_URL=<MY_AUTH_URL>
   $ export OS_REGION_NAME=<MY_REGION_NAME>
   $ export OS_USERNAME=<MY_USERNAME>
   $ export OS_PASSWORD=<MY_PASSWORD>
   $ export OS_USER_DOMAIN_NAME=<MY_DOMAIN_NAME>
   $ export OS_PROJECT_NAME=<MY_PROJECT_NAME>
   $ export OS_PROJECT_DOMAIN_NAME=<MY_PROJECT_DOMAIN_NAME>

   # For authentication based on tokens
   $ export OS_STORAGE_URL=<MY_STORAGE_URL>
   $ export OS_AUTH_TOKEN=<MY_AUTH_TOKEN>


Restic should be compatible with [OpenStack RC
file](https://docs.openstack.org/user-guide/common/cli-set-environment-variables-using-openstack-rc.html)
in most cases.

Once environment variables are set up, a new repository can be created. The
name of swift container and optional path can be specified. If
the container does not exist, it will be created automatically:

.. code-block:: console

   $ restic -r swift:container_name:/path init   # path is optional
   enter password for new backend:
   enter password again:
   created restic backend eefee03bbd at swift:container_name:/path
   Please note that knowledge of your password is required to access the repository.
   Losing your password means that your data is irrecoverably lost.

The policy of new container created by restic can be changed using environment variable:

.. code-block:: console

   $ export SWIFT_DEFAULT_CONTAINER_POLICY=<MY_CONTAINER_POLICY>


Backblaze B2
~~~~~~~~~~~~

Restic can backup data to any Backblaze B2 bucket. You need to first setup the
following environment variables with the credentials you obtained when signed
into your B2 account:

.. code-block:: console

    $ export B2_ACCOUNT_ID=<MY_ACCOUNT_ID>
    $ export B2_ACCOUNT_KEY=<MY_SECRET_ACCOUNT_KEY>

You can then easily initialize a repository stored at Backblaze B2. If the
bucket does not exist yet, it will be created:

.. code-block:: console

    $ restic -r b2:bucketname:path/to/repo init
    enter password for new backend:
    enter password again:
    created restic backend eefee03bbd at b2:bucketname:path/to/repo
    Please note that knowledge of your password is required to access the repository.
    Losing your password means that your data is irrecoverably lost.

The number of concurrent connections to the B2 service can be set with the `-o
b2.connections=10`. By default, at most five parallel connections are
established.

Microsoft Azure Blob Storage
~~~~~~~~~~~~~~~~~~~~~~~~~~~~

You can also store backups on Microsoft Azure Blob Storage. Export the Azure
account name and key as follows:

.. code-block:: console

    $ export AZURE_ACCOUNT_NAME=<ACCOUNT_NAME>
    $ export AZURE_ACCOUNT_KEY=<SECRET_KEY>

Afterwards you can initialize a repository in a container called `foo` in the
root path like this:

.. code-block:: console

    $ restic -r azure:foo:/ init
    enter password for new backend:
    enter password again:

    created restic backend a934bac191 at azure:foo:/
    [...]

The number of concurrent connections to the B2 service can be set with the
`-o azure.connections=10`. By default, at most five parallel connections are
established.


Google Cloud Storage
~~~~~~~~~~~~~~~~~~~~

Restic supports Google Cloud Storage as a backend.

Restic connects to Google Cloud Storage via a `service account`_.

For normal restic operation, the service account must have the
``storage.objects.{create,delete,get,list}`` permissions for the bucket. These
are included in the "Storage Object Admin" role.

``restic init`` can create the repository bucket. Doing so requires the
``storage.buckets.create`` permission ("Storage Admin" role). If the bucket
already exists that permission is unnecessary.

To use the Google Cloud Storage backend, first `create a service account key`_
and download the JSON credentials file.

Second, find the Google Project ID that you can see in the Google Cloud
Platform console at the "Storage/Settings" menu. Export the path to the JSON
key file and the project ID as follows:

.. code-block:: console

    $ export GOOGLE_PROJECT_ID=123123123123
    $ export GOOGLE_APPLICATION_CREDENTIALS=$HOME/.config/gs-secret-restic-key.json

Then you can use the ``gs:`` backend type to create a new repository in the
bucket `foo` at the root path:

.. code-block:: console

    $ restic -r gs:foo:/ init
    enter password for new backend:
    enter password again:

    created restic backend bde47d6254 at gs:foo:/
    [...]

The number of concurrent connections to the GCS service can be set with the
`-o gs.connections=10`. By default, at most five parallel connections are
established.

.. _service account: https://cloud.google.com/storage/docs/authentication#service_accounts
.. _create a service account key: https://cloud.google.com/storage/docs/authentication#generating-a-private-key


Password prompt on Windows
~~~~~~~~~~~~~~~~~~~~~~~~~~

At the moment, restic only supports the default Windows console
interaction. If you use emulation environments like
`MSYS2 <https://msys2.github.io/>`__ or
`Cygwin <https://www.cygwin.com/>`__, which use terminals like
``Mintty`` or ``rxvt``, you may get a password error:

You can workaround this by using a special tool called ``winpty`` (look
`here <https://sourceforge.net/p/msys2/wiki/Porting/>`__ and
`here <https://github.com/rprichard/winpty>`__ for detail information).
On MSYS2, you can install ``winpty`` as follows:

.. code-block:: console

    $ pacman -S winpty
    $ winpty restic -r /tmp/backup init

Create a snapshot
-----------------

Now we're ready to backup some data. The contents of a directory at a
specific point in time is called a "snapshot" in restic. Run the
following command and enter the repository password you chose above
again:

.. code-block:: console

    $ restic -r /tmp/backup backup ~/work
    enter password for repository:
    scan [/home/user/work]
    scanned 764 directories, 1816 files in 0:00
    [0:29] 100.00%  54.732 MiB/s  1.582 GiB / 1.582 GiB  2580 / 2580 items  0 errors  ETA 0:00
    duration: 0:29, 54.47MiB/s
    snapshot 40dc1520 saved

As you can see, restic created a backup of the directory and was pretty
fast! The specific snapshot just created is identified by a sequence of
hexadecimal characters, ``40dc1520`` in this case.

If you run the command again, restic will create another snapshot of
your data, but this time it's even faster. This is de-duplication at
work!

.. code-block:: console

    $ restic -r /tmp/backup backup ~/work
    enter password for repository:
    using parent snapshot 40dc1520aa6a07b7b3ae561786770a01951245d2367241e71e9485f18ae8228c
    scan [/home/user/work]
    scanned 764 directories, 1816 files in 0:00
    [0:00] 100.00%  0B/s  1.582 GiB / 1.582 GiB  2580 / 2580 items  0 errors  ETA 0:00
    duration: 0:00, 6572.38MiB/s
    snapshot 79766175 saved

You can even backup individual files in the same repository.

.. code-block:: console

    $ restic -r /tmp/backup backup ~/work.txt
    scan [/home/user/work.txt]
    scanned 0 directories, 1 files in 0:00
    [0:00] 100.00%  0B/s  220B / 220B  1 / 1 items  0 errors  ETA 0:00
    duration: 0:00, 0.03MiB/s
    snapshot 31f7bd63 saved

In fact several hosts may use the same repository to backup directories
and files leading to a greater de-duplication.

Please be aware that when you backup different directories (or the
directories to be saved have a variable name component like a
time/date), restic always needs to read all files and only afterwards
can compute which parts of the files need to be saved. When you backup
the same directory again (maybe with new or changed files) restic will
find the old snapshot in the repo and by default only reads those files
that are new or have been modified since the last snapshot. This is
decided based on the modify date of the file in the file system.

Now is a good time to run ``restic check`` to verify that all data
is properly stored in the repository. You should run this command regularly
to make sure the internal structure of the repository is free of errors.

You can exclude folders and files by specifying exclude-patterns. Either
specify them with multiple ``--exclude``'s or one ``--exclude-file``

.. code-block:: console

    $ cat exclude
    # exclude go-files
    *.go
    # exclude foo/x/y/z/bar foo/x/bar foo/bar
    foo/**/bar
    $ restic -r /tmp/backup backup ~/work --exclude=*.c --exclude-file=exclude

Patterns use `filepath.Glob <https://golang.org/pkg/path/filepath/#Glob>`__ internally,
see `filepath.Match <https://golang.org/pkg/path/filepath/#Match>`__ for syntax.
Additionally ``**`` excludes arbitrary subdirectories.
Environment-variables in exclude-files are expanded with
`os.ExpandEnv <https://golang.org/pkg/os/#ExpandEnv>`__.

By specifying the option ``--one-file-system`` you can instruct restic
to only backup files from the file systems the initially specified files
or directories reside on. For example, calling restic like this won't
backup ``/sys`` or ``/dev`` on a Linux system:

.. code-block:: console

    $ restic -r /tmp/backup backup --one-file-system /

By using the ``--files-from`` option you can read the files you want to
backup from a file. This is especially useful if a lot of files have to
be backed up that are not in the same folder or are maybe pre-filtered
by other software.

For example maybe you want to backup files that have a certain filename
in them:

.. code-block:: console

    $ find /tmp/somefiles | grep 'PATTERN' > /tmp/files_to_backup

You can then use restic to backup the filtered files:

.. code-block:: console

    $ restic -r /tmp/backup backup --files-from /tmp/files_to_backup

Incidentally you can also combine ``--files-from`` with the normal files
args:

.. code-block:: console

    $ restic -r /tmp/backup backup --files-from /tmp/files_to_backup /tmp/some_additional_file


Backing up special items
~~~~~~~~~~~~~~~~~~~~~~~~

**Symlinks** are archieved as symlinks, ``restic``does not follow them.
When you restore, you get the same symlink again, with the same link target
and the same timestamps.

If there is a **bind-mount** below a directory that is to be saved, restic descends into it.

**Device files** are saved and restored as device files. This means that e.g. ``/dev/sda`` is
archived as a block device file and restored as such. This also means that the content of the
corresponding disk is not read, at least not from the device file.


Reading data from stdin
~~~~~~~~~~~~~~~~~~~~~~~

Sometimes it can be nice to directly save the output of a program, e.g.
``mysqldump`` so that the SQL can later be restored. Restic supports
this mode of operation, just supply the option ``--stdin`` to the
``backup`` command like this:

.. code-block:: console

    $ mysqldump [...] | restic -r /tmp/backup backup --stdin

This creates a new snapshot of the output of ``mysqldump``. You can then
use e.g. the fuse mounting option (see below) to mount the repository
and read the file.

By default, the file name ``stdin`` is used, a different name can be
specified with ``--stdin-filename``, e.g. like this:

.. code-block:: console

    $ mysqldump [...] | restic -r /tmp/backup backup --stdin --stdin-filename production.sql

Tags for backup
~~~~~~~~~~~~~~~

Snapshots can have one or more tags, short strings which add identifying
information. Just specify the tags for a snapshot one by one with ``--tag``:

.. code-block:: console

    $ restic -r /tmp/backup backup --tag projectX --tag foo --tag bar ~/work
    [...]

The tags can later be used to keep (or forget) snapshots with the ``forget``
command. The command ``tag`` can be used to modify tags on an existing
snapshot.

List all snapshots
------------------

Now, you can list all the snapshots stored in the repository:

.. code-block:: console

    $ restic -r /tmp/backup snapshots
    enter password for repository:
    ID        Date                 Host    Tags   Directory
    ----------------------------------------------------------------------
    40dc1520  2015-05-08 21:38:30  kasimir        /home/user/work
    79766175  2015-05-08 21:40:19  kasimir        /home/user/work
    bdbd3439  2015-05-08 21:45:17  luigi          /home/art
    590c8fc8  2015-05-08 21:47:38  kazik          /srv
    9f0bc19e  2015-05-08 21:46:11  luigi          /srv

You can filter the listing by directory path:

.. code-block:: console

    $ restic -r /tmp/backup snapshots --path="/srv"
    enter password for repository:
    ID        Date                 Host    Tags   Directory
    ----------------------------------------------------------------------
    590c8fc8  2015-05-08 21:47:38  kazik          /srv
    9f0bc19e  2015-05-08 21:46:11  luigi          /srv

Or filter by host:

.. code-block:: console

    $ restic -r /tmp/backup snapshots --host luigi
    enter password for repository:
    ID        Date                 Host    Tags   Directory
    ----------------------------------------------------------------------
    bdbd3439  2015-05-08 21:45:17  luigi          /home/art
    9f0bc19e  2015-05-08 21:46:11  luigi          /srv

Combining filters is also possible.

Restore a snapshot
------------------

Restoring a snapshot is as easy as it sounds, just use the following
command to restore the contents of the latest snapshot to
``/tmp/restore-work``:

.. code-block:: console

    $ restic -r /tmp/backup restore 79766175 --target /tmp/restore-work
    enter password for repository:
    restoring <Snapshot of [/home/user/work] at 2015-05-08 21:40:19.884408621 +0200 CEST> to /tmp/restore-work

Use the word ``latest`` to restore the last backup. You can also combine
``latest`` with the ``--host`` and ``--path`` filters to choose the last
backup for a specific host, path or both.

.. code-block:: console

    $ restic -r /tmp/backup restore latest --target /tmp/restore-art --path "/home/art" --host luigi
    enter password for repository:
    restoring <Snapshot of [/home/art] at 2015-05-08 21:45:17.884408621 +0200 CEST> to /tmp/restore-art

Use ``--exclude`` and ``--include`` to restrict the restore to a subset of
files in the snapshot. For example, to restore a single file:

.. code-block:: console

    $ restic -r /tmp/backup restore 79766175 --target /tmp/restore-work --include /work/foo
    enter password for repository:
    restoring <Snapshot of [/home/user/work] at 2015-05-08 21:40:19.884408621 +0200 CEST> to /tmp/restore-work

This will restore the file ``foo`` to ``/tmp/restore-work/work/foo``.

Manage repository keys
----------------------

The ``key`` command allows you to set multiple access keys or passwords
per repository. In fact, you can use the ``list``, ``add``, ``remove``
and ``passwd`` sub-commands to manage these keys very precisely:

.. code-block:: console

    $ restic -r /tmp/backup key list
    enter password for repository:
     ID          User        Host        Created
    ----------------------------------------------------------------------
    *eb78040b    username    kasimir   2015-08-12 13:29:57

    $ restic -r /tmp/backup key add
    enter password for repository:
    enter password for new key:
    enter password again:
    saved new key as <Key of username@kasimir, created on 2015-08-12 13:35:05.316831933 +0200 CEST>

    $ restic -r backup key list
    enter password for repository:
     ID          User        Host        Created
    ----------------------------------------------------------------------
     5c657874    username    kasimir   2015-08-12 13:35:05
    *eb78040b    username    kasimir   2015-08-12 13:29:57

Manage tags
-----------

Managing tags on snapshots is done with the ``tag`` command. The
existing set of tags can be replaced completely, tags can be added to
removed. The result is directly visible in the ``snapshots`` command.

Let's say we want to tag snapshot ``590c8fc8`` with the tags ``NL`` and
``CH`` and remove all other tags that may be present, the following
command does that:

.. code-block:: console

    $ restic -r /tmp/backup tag --set NL --set CH 590c8fc8
    Create exclusive lock for repository
    Modified tags on 1 snapshots

Note the snapshot ID has changed, so between each change we need to look
up the new ID of the snapshot. But there is an even better way, the
``tag`` command accepts ``--tag`` for a filter, so we can filter
snapshots based on the tag we just added.

So we can add and remove tags incrementally like this:

.. code-block:: console

    $ restic -r /tmp/backup tag --tag NL --remove CH
    Create exclusive lock for repository
    Modified tags on 1 snapshots

    $ restic -r /tmp/backup tag --tag NL --add UK
    Create exclusive lock for repository
    Modified tags on 1 snapshots

    $ restic -r /tmp/backup tag --tag NL --remove NL
    Create exclusive lock for repository
    Modified tags on 1 snapshots

    $ restic -r /tmp/backup tag --tag NL --add SOMETHING
    No snapshots were modified

Check integrity and consistency
-------------------------------

Imagine your repository is saved on a server that has a faulty hard
drive, or even worse, attackers get privileged access and modify your
backup with the intention to make you restore malicious data:

.. code-block:: console

    $ sudo echo "boom" >> backup/index/d795ffa99a8ab8f8e42cec1f814df4e48b8f49129360fb57613df93739faee97

In order to detect these things, it is a good idea to regularly use the
``check`` command to test whether everything is alright, your precious
backup data is consistent and the integrity is unharmed:

.. code-block:: console

    $ restic -r /tmp/backup check
    Load indexes
    ciphertext verification failed

Trying to restore a snapshot which has been modified as shown above will
yield the same error:

.. code-block:: console

    $ restic -r /tmp/backup restore 79766175 --target /tmp/restore-work
    Load indexes
    ciphertext verification failed

Mount a repository
------------------

Browsing your backup as a regular file system is also very easy. First,
create a mount point such as ``/mnt/restic`` and then use the following
command to serve the repository with FUSE:

.. code-block:: console

    $ mkdir /mnt/restic
    $ restic -r /tmp/backup mount /mnt/restic
    enter password for repository:
    Now serving /tmp/backup at /mnt/restic
    Don't forget to umount after quitting!

Mounting repositories via FUSE is not possible on Windows and OpenBSD.

Restic supports storage and preservation of hard links. However, since
hard links exist in the scope of a filesystem by definition, restoring
hard links from a fuse mount should be done by a program that preserves
hard links. A program that does so is rsync, used with the option
--hard-links.

Removing old snapshots
----------------------

All backup space is finite, so restic allows removing old snapshots.
This can be done either manually (by specifying a snapshot ID to remove)
or by using a policy that describes which snapshots to forget. For all
remove operations, two commands need to be called in sequence:
``forget`` to remove a snapshot and ``prune`` to actually remove the
data that was referenced by the snapshot from the repository. This can
be automated with the ``--prune`` option of the ``forget`` command,
which runs ``prune`` automatically if snapshots have been removed.

It is advisable to run ``restic check`` after pruning, to make sure
you are alerted, should the internal data structures of the repository
be damaged.

Remove a single snapshot
~~~~~~~~~~~~~~~~~~~~~~~~

The command ``snapshots`` can be used to list all snapshots in a
repository like this:

.. code-block:: console

    $ restic -r /tmp/backup snapshots
    enter password for repository:
    ID        Date                 Host      Tags  Directory
    ----------------------------------------------------------------------
    40dc1520  2015-05-08 21:38:30  kasimir         /home/user/work
    79766175  2015-05-08 21:40:19  kasimir         /home/user/work
    bdbd3439  2015-05-08 21:45:17  luigi           /home/art
    590c8fc8  2015-05-08 21:47:38  kazik           /srv
    9f0bc19e  2015-05-08 21:46:11  luigi           /srv

In order to remove the snapshot of ``/home/art``, use the ``forget``
command and specify the snapshot ID on the command line:

.. code-block:: console

    $ restic -r /tmp/backup forget bdbd3439
    enter password for repository:
    removed snapshot d3f01f63

Afterwards this snapshot is removed:

.. code-block:: console

    $ restic -r /tmp/backup snapshots
    enter password for repository:
    ID        Date                 Host     Tags  Directory
    ----------------------------------------------------------------------
    40dc1520  2015-05-08 21:38:30  kasimir        /home/user/work
    79766175  2015-05-08 21:40:19  kasimir        /home/user/work
    590c8fc8  2015-05-08 21:47:38  kazik          /srv
    9f0bc19e  2015-05-08 21:46:11  luigi          /srv

But the data that was referenced by files in this snapshot is still
stored in the repository. To cleanup unreferenced data, the ``prune``
command must be run:

.. code-block:: console

    $ restic -r /tmp/backup prune
    enter password for repository:

    counting files in repo
    building new index for repo
    [0:00] 100.00%  22 / 22 files
    repository contains 22 packs (8512 blobs) with 100.092 MiB bytes
    processed 8512 blobs: 0 duplicate blobs, 0B duplicate
    load all snapshots
    find data that is still in use for 1 snapshots
    [0:00] 100.00%  1 / 1 snapshots
    found 8433 of 8512 data blobs still in use
    will rewrite 3 packs
    creating new index
    [0:00] 86.36%  19 / 22 files
    saved new index as 544a5084
    done

Afterwards the repository is smaller.

You can automate this two-step process by using the ``--prune`` switch
to ``forget``:

.. code-block:: console

    $ restic forget --keep-last 1 --prune
    snapshots for host mopped, directories /home/user/work:

    keep 1 snapshots:
    ID        Date                 Host        Tags        Directory
    ----------------------------------------------------------------------
    4bba301e  2017-02-21 10:49:18  mopped                  /home/user/work

    remove 1 snapshots:
    ID        Date                 Host        Tags        Directory
    ----------------------------------------------------------------------
    8c02b94b  2017-02-21 10:48:33  mopped                  /home/user/work

    1 snapshots have been removed, running prune
    counting files in repo
    building new index for repo
    [0:00] 100.00%  37 / 37 packs
    repository contains 37 packs (5521 blobs) with 151.012 MiB bytes
    processed 5521 blobs: 0 duplicate blobs, 0B duplicate
    load all snapshots
    find data that is still in use for 1 snapshots
    [0:00] 100.00%  1 / 1 snapshots
    found 5323 of 5521 data blobs still in use, removing 198 blobs
    will delete 0 packs and rewrite 27 packs, this frees 22.106 MiB
    creating new index
    [0:00] 100.00%  30 / 30 packs
    saved new index as b49f3e68
    done

Removing snapshots according to a policy
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Removing snapshots manually is tedious and error-prone, therefore restic
allows specifying which snapshots should be removed automatically
according to a policy. You can specify how many hourly, daily, weekly,
monthly and yearly snapshots to keep, any other snapshots are removed.
The most important command-line parameter here is ``--dry-run`` which
instructs restic to not remove anything but print which snapshots would
be removed.

When ``forget`` is run with a policy, restic loads the list of all
snapshots, then groups these by host name and list of directories. The grouping
options can be set with ``--group-by``, to only group snapshots by paths and
tags use ``--group-by paths,tags``. The policy is then applied to each group of
snapshots separately. This is a safety feature.

The ``forget`` command accepts the following parameters:

-  ``--keep-last n`` never delete the ``n`` last (most recent) snapshots
-  ``--keep-hourly n`` for the last ``n`` hours in which a snapshot was
   made, keep only the last snapshot for each hour.
-  ``--keep-daily n`` for the last ``n`` days which have one or more
   snapshots, only keep the last one for that day.
-  ``--keep-weekly n`` for the last ``n`` weeks which have one or more
   snapshots, only keep the last one for that week.
-  ``--keep-monthly n`` for the last ``n`` months which have one or more
   snapshots, only keep the last one for that month.
-  ``--keep-yearly n`` for the last ``n`` years which have one or more
   snapshots, only keep the last one for that year.
-  ``--keep-tag`` keep all snapshots which have all tags specified by
   this option (can be specified multiple times).

Additionally, you can restrict removing snapshots to those which have a
particular hostname with the ``--hostname`` parameter, or tags with the
``--tag`` option. When multiple tags are specified, only the snapshots
which have all the tags are considered. For example, the following command
removes all but the latest snapshot of all snapshots that have the tag ``foo``:

.. code-block:: console

   $ restic forget --tag foo --keep-last 1

This command removes all but the last snapshot of all snapshots that have
either the ``foo`` or ``bar`` tag set:

.. code-block:: console

   $ restic forget --tag foo --tag bar --keep-last 1

To only keep the last snapshot of all snapshots with both the tag ``foo`` and
``bar`` set use:

.. code-block:: console

   $ restic forget --tag foo,tag bar --keep-last 1

All the ``--keep-*`` options above only count
hours/days/weeks/months/years which have a snapshot, so those without a
snapshot are ignored.

All snapshots are evaluated counted against all matching keep-* counts. A
single snapshot on 2017-09-30 (Sun) will count as a daily, weekly and monthly.

Let's explain this with an example: Suppose you have only made a backup
on each Sunday for 12 weeks. Then ``forget --keep-daily 4`` will keep
the last four snapshots for the last four Sundays, but remove the rest.
Only counting the days which have a backup and ignore the ones without
is a safety feature: it prevents restic from removing many snapshots
when no new ones are created. If it was implemented otherwise, running
``forget --keep-daily 4`` on a Friday would remove all snapshots!

Another example: Suppose you make daily backups for 100 years. Then
``forget --keep-daily 7 --keep-weekly 5 --keep-monthly 12 --keep-yearly 75``
will keep the most recent 7 daily snapshots, then 4 (remember, 7 dailies
already include a week!) last-day-of-the-weeks and 11 or 12
last-day-of-the-months (11 or 12 depends if the 5 weeklies cross a month).
And finally 75 last-day-of-the-year snapshots. All other snapshots are
removed.

Autocompletion
--------------

Restic can write out a bash compatible autocompletion script:

.. code-block:: console

    $ ./restic autocomplete --help
    The "autocomplete" command generates a shell autocompletion script.

    NOTE: The current version supports Bash only.
          This should work for *nix systems with Bash installed.

By default, the file is written directly to ``/etc/bash_completion.d/``
for convenience, and the command may need superuser rights, e.g.

.. code-block:: console

    $ sudo restic autocomplete

    Usage:
      restic autocomplete [flags]

    Flags:
          --completionfile string   autocompletion file (default "/etc/bash_completion.d/restic.sh")

    Global Flags:
          --json                   set output mode to JSON for commands that support it
          --no-lock                do not lock the repo, this allows some operations on read-only repos
      -o, --option key=value       set extended option (key=value, can be specified multiple times)
      -p, --password-file string   read the repository password from a file
      -q, --quiet                  do not output comprehensive progress report
      -r, --repo string            repository to backup to or restore from (default: $RESTIC_REPOSITORY)

Debugging
---------

The program can be built with debug support like this:

.. code-block:: console

    $ go run build.go -tags debug

Afterwards, extensive debug messages are written to the file in
environment variable ``DEBUG_LOG``, e.g.:

.. code-block:: console

    $ DEBUG_LOG=/tmp/restic-debug.log restic backup ~/work

If you suspect that there is a bug, you can have a look at the debug
log. Please be aware that the debug log might contain sensitive
information such as file and directory names.

The debug log will always contain all log messages restic generates. You
can also instruct restic to print some or all debug messages to stderr.
These can also be limited to e.g. a list of source files or a list of
patterns for function names. The patterns are globbing patterns (see the
documentation for `path.Glob <https://golang.org/pkg/path/#Glob>`__), multiple
patterns are separated by commas. Patterns are case sensitive.

Printing all log messages to the console can be achieved by setting the
file filter to ``*``:

.. code-block:: console

    $ DEBUG_FILES=* restic check

If you want restic to just print all debug log messages from the files
``main.go`` and ``lock.go``, set the environment variable
``DEBUG_FILES`` like this:

.. code-block:: console

    $ DEBUG_FILES=main.go,lock.go restic check

The following command line instructs restic to only print debug
statements originating in functions that match the pattern ``*unlock*``
(case sensitive):

.. code-block:: console

    $ DEBUG_FUNCS=*unlock* restic check

Under the hood
--------------

Browse repository objects
~~~~~~~~~~~~~~~~~~~~~~~~~

Internally, a repository stores data of several different types
described in the `design
documentation <https://github.com/restic/restic/blob/master/doc/Design.rst>`__.
You can ``list`` objects such as blobs, packs, index, snapshots, keys or
locks with the following command:

.. code-block:: console

    $ restic -r /tmp/backup list snapshots
    d369ccc7d126594950bf74f0a348d5d98d9e99f3215082eb69bf02dc9b3e464c

The ``find`` command searches for a given
`pattern <http://golang.org/pkg/path/filepath/#Match>`__ in the
repository.

.. code-block:: console

    $ restic -r backup find test.txt
    debug log file restic.log
    debug enabled
    enter password for repository:
    found 1 matching entries in snapshot 196bc5760c909a7681647949e80e5448e276521489558525680acf1bd428af36
      -rw-r--r--   501    20      5 2015-08-26 14:09:57 +0200 CEST path/to/test.txt

The ``cat`` command allows you to display the JSON representation of the
objects or its raw content.

.. code-block:: console

    $ restic -r /tmp/backup cat snapshot d369ccc7d126594950bf74f0a348d5d98d9e99f3215082eb69bf02dc9b3e464c
    enter password for repository:
    {
      "time": "2015-08-12T12:52:44.091448856+02:00",
      "tree": "05cec17e8d3349f402576d02576a2971fc0d9f9776ce2f441c7010849c4ff5af",
      "paths": [
        "/home/user/work"
      ],
      "hostname": "kasimir",
      "username": "username",
      "uid": 501,
      "gid": 20
    }

Metadata handling
~~~~~~~~~~~~~~~~~

Restic saves and restores most default attributes, including extended attributes like ACLs.
Sparse files are not handled in a special way yet, and aren't restored.

The following metadata is handled by restic:

- Name
- Type
- Mode
- ModTime
- AccessTime
- ChangeTime
- UID
- GID
- User
- Group
- Inode
- Size
- Links
- LinkTarget
- Device
- Content
- Subtree
- ExtendedAttributes

Scripting
---------

Restic supports the output of some commands in JSON format, the JSON
data can then be processed by other programs (e.g.
`jq <https://stedolan.github.io/jq/>`__). The following example
lists all snapshots as JSON and uses ``jq`` to pretty-print the result:

.. code-block:: console

    $ restic -r /tmp/backup snapshots --json | jq .
    [
      {
        "time": "2017-03-11T09:57:43.26630619+01:00",
        "tree": "bf25241679533df554fc0fd0ae6dbb9dcf1859a13f2bc9dd4543c354eff6c464",
        "paths": [
          "/home/work/doc"
        ],
        "hostname": "kasimir",
        "username": "fd0",
        "uid": 1000,
        "gid": 100,
        "id": "bbeed6d28159aa384d1ccc6fa0b540644b1b9599b162d2972acda86b1b80f89e"
      },
      {
        "time": "2017-03-11T09:58:57.541446938+01:00",
        "tree": "7f8c95d3420baaac28dc51609796ae0e0ecfb4862b609a9f38ffaf7ae2d758da",
        "paths": [
          "/home/user/shared"
        ],
        "hostname": "kasimir",
        "username": "fd0",
        "uid": 1000,
        "gid": 100,
        "id": "b157d91c16f0ba56801ece3a708dfc53791fe2a97e827090d6ed9a69a6ebdca0"
      }
    ]

Temporary files
---------------

During some operations (e.g. ``backup`` and ``prune``) restic uses
temporary files to store data. These files will, by default, be saved to
the system's temporary directory, on Linux this is usually located in
``/tmp/``. The environment variable ``TMPDIR`` can be used to specify a
different directory, e.g. to use the directory ``/var/tmp/restic-tmp``
instead of the default, set the environment variable like this:

.. code-block:: console

    $ export TMPDIR=/var/tmp/restic-tmp
    $ restic -r /tmp/backup backup ~/work

Caching
-------

Restic keeps a cache with some files from the repository on the local machine.
This allows faster operations, since meta data does not need to be loaded from
a remote repository. The cache is automatically created, usually in the
directory ``.cache/restic`` in the user's home directory. The environment
variable ``XDG_CACHE_DIR`` or the command line parameter ``--cache-dir`` can
each be used to specify where the cache is located. The parameter
``--no-cache`` disables the cache entirely. In this case, all data is loaded
from the repo.

The cache is ephemeral: When a file cannot be read from the cache, it is loaded
from the repository.


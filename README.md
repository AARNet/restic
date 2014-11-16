[![Stories in Ready](https://badge.waffle.io/fd0/khepri.png?label=ready&title=Ready)](https://waffle.io/fd0/khepri)
WARNING
=======

WARNING: At the moment, consider khepri as alpha quality software, it is not
yet finished. Do not use it for real data!

Khepri
======

Khepri is a program that does backups right. The design goals are:

 * Easy: Doing backups should be a frictionless process, otherwise you are
   tempted to skip it.  Khepri should be easy to configure and use, so that in
   the unlikely event of a data loss you can just restore it. Likewise,
   restoring data should not be complicated.

 * Fast: Backing up your data with khepri should only be limited by your
   network or harddisk bandwidth so that you can backup your files every day.
   Nobody does backups if it takes too much time. Restoring backups should only
   transfer data that is needed for the files that are to be restored, so that
   this process is also fast.

 * Verifiable: Much more important than backup is restore, so khepri enables
   you to easily verify that all data can be restored.

 * Secure: Khepri uses cryptography to guarantee confidentiality and integrity
   of your data. The location the backup data is stored is assumed not to be a
   trusted environment (e.g. a shared space where others like system
   administrators are able to access your backups). Khepri is built to secure
   your data against such attackers.

 * Efficient: With the growth of data, additional snapshots should only take
   the storage of the actual increment. Even more, duplicate data should be
   de-duplicated before it is actually written to the storage backend to save
   precious backup space.


Building
========

Install Go (at least 1.2), then run:


export GOPATH=~/src/go
go get github.com/fd0/khepri/cmd/khepri
$GOPATH/bin/khepri --help


License
=======

Khepri is licensed under "BSD 2-Clause License". You can find the complete text
in the file `LICENSE`.

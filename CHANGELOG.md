Changelog for restic 0.10.0 (2020-09-19)
=======================================

The following sections list the changes in restic 0.10.0 relevant to
restic users. The changes are ordered by importance.

Summary
-------

 * Fix #1863: Report correct number of directories processed by backup
 * Fix #2254: Fix tar issues when dumping `/`
 * Fix #2281: Handle format verbs like '%' properly in `find` output
 * Fix #2298: Do not hang when run as a background job
 * Fix #2389: Fix mangled json output of backup command
 * Fix #2390: Refresh lock timestamp
 * Fix #2429: Backup --json reports total_bytes_processed as 0
 * Fix #2469: Fix incorrect bytes stats in `diff` command
 * Fix #2518: Do not crash with Synology NAS sftp server
 * Fix #2531: Fix incorrect size calculation in `stats --mode restore-size`
 * Fix #2537: Fix incorrect file counts in `stats --mode restore-size`
 * Fix #2592: SFTP backend supports IPv6 addresses
 * Fix #2607: Honor RESTIC_CACHE_DIR environment variable on Mac and Windows
 * Fix #2668: Don't abort the stats command when data blobs are missing
 * Fix #2674: Add stricter prune error checks
 * Fix #2899: Fix possible crash in the progress bar of check --read-data
 * Chg #2482: Remove vendored dependencies
 * Chg #2546: Return exit code 3 when failing to backup all source data
 * Chg #2600: Update dependencies, require Go >= 1.13
 * Chg #1597: Honor the --no-lock flag in the mount command
 * Enh #1570: Support specifying multiple host flags for various commands
 * Enh #1680: Optimize `restic mount`
 * Enh #2072: Display snapshot date when using `restic find`
 * Enh #2175: Allow specifying user and host when creating keys
 * Enh #2277: Add support for ppc64le
 * Enh #2395: Ignore sync errors when operation not supported by local filesystem
 * Enh #2427: Add flag `--iexclude-file` to backup command
 * Enh #2569: Support excluding files by their size
 * Enh #2571: Self-heal missing file parts during backup of unchanged files
 * Enh #2858: Support filtering snapshots by tag and path in the stats command
 * Enh #323: Add command for copying snapshots between repositories
 * Enh #551: Use optimized library for hash calculation of file chunks
 * Enh #2195: Simplify and improve restore performance
 * Enh #2328: Improve speed of check command
 * Enh #2423: Support user@domain parsing as user
 * Enh #2576: Improve the chunking algorithm
 * Enh #2598: Improve speed of diff command
 * Enh #2599: Slightly reduce memory usage of prune and stats commands
 * Enh #2733: S3 backend: Add support for WebIdentityTokenFile
 * Enh #2773: Optimize handling of new index entries
 * Enh #2781: Reduce memory consumption of in-memory index
 * Enh #2786: Optimize `list blobs` command
 * Enh #2790: Optimized file access in restic mount
 * Enh #2840: Speed-up file deletion in forget, prune and rebuild-index

Details
-------

 * Bugfix #1863: Report correct number of directories processed by backup

   The directory statistics calculation was fixed to report the actual number of processed
   directories instead of always zero.

   https://github.com/restic/restic/issues/1863

 * Bugfix #2254: Fix tar issues when dumping `/`

   We've fixed an issue with dumping either `/` or files on the first sublevel e.g. `/foo` to tar.
   This also fixes tar dumping issues on Windows where this issue could also happen.

   https://github.com/restic/restic/issues/2254
   https://github.com/restic/restic/issues/2357
   https://github.com/restic/restic/pull/2255

 * Bugfix #2281: Handle format verbs like '%' properly in `find` output

   The JSON or "normal" output of the `find` command can now deal with file names that contain
   substrings which the Golang `fmt` package considers "format verbs" like `%s`.

   https://github.com/restic/restic/issues/2281

 * Bugfix #2298: Do not hang when run as a background job

   Restic did hang on exit while restoring the terminal configuration when it was started as a
   background job, for example using `restic ... &`. This has been fixed by only restoring the
   terminal configuration when restic is interrupted while reading a password from the
   terminal.

   https://github.com/restic/restic/issues/2298

 * Bugfix #2389: Fix mangled json output of backup command

   We've fixed a race condition in the json output of the backup command that could cause multiple
   lines to get mixed up. We've also ensured that the backup summary is printed last.

   https://github.com/restic/restic/issues/2389
   https://github.com/restic/restic/pull/2545

 * Bugfix #2390: Refresh lock timestamp

   Long-running operations did not refresh lock timestamp, resulting in locks becoming stale.
   This is now fixed.

   https://github.com/restic/restic/issues/2390

 * Bugfix #2429: Backup --json reports total_bytes_processed as 0

   We've fixed the json output of total_bytes_processed. The non-json output was already fixed
   with pull request #2138 but left the json output untouched.

   https://github.com/restic/restic/issues/2429

 * Bugfix #2469: Fix incorrect bytes stats in `diff` command

   In some cases, the wrong number of bytes (e.g. 16777215.998 TiB) were reported by the `diff`
   command. This is now fixed.

   https://github.com/restic/restic/issues/2469

 * Bugfix #2518: Do not crash with Synology NAS sftp server

   It was found that when restic is used to store data on an sftp server on a Synology NAS with a
   relative path (one which does not start with a slash), it may go into an endless loop trying to
   create directories on the server. We've fixed this bug by using a function in the sftp library
   instead of our own implementation.

   The bug was discovered because the Synology sftp server behaves erratic with non-absolute
   path (e.g. `home/restic-repo`). This can be resolved by just using an absolute path instead
   (`/home/restic-repo`). We've also added a paragraph in the FAQ.

   https://github.com/restic/restic/issues/2518
   https://github.com/restic/restic/issues/2363
   https://github.com/restic/restic/pull/2530

 * Bugfix #2531: Fix incorrect size calculation in `stats --mode restore-size`

   The restore-size mode of stats was counting hard-linked files as if they were independent.

   https://github.com/restic/restic/issues/2531

 * Bugfix #2537: Fix incorrect file counts in `stats --mode restore-size`

   The restore-size mode of stats was failing to count empty directories and some files with hard
   links.

   https://github.com/restic/restic/issues/2537

 * Bugfix #2592: SFTP backend supports IPv6 addresses

   The SFTP backend now supports IPv6 addresses natively, without relying on aliases in the
   external SSH configuration.

   https://github.com/restic/restic/pull/2592

 * Bugfix #2607: Honor RESTIC_CACHE_DIR environment variable on Mac and Windows

   On Mac and Windows, the RESTIC_CACHE_DIR environment variable was ignored. This variable can
   now be used on all platforms to set the directory where restic stores caches.

   https://github.com/restic/restic/pull/2607

 * Bugfix #2668: Don't abort the stats command when data blobs are missing

   Runing the stats command in the blobs-per-file mode on a repository with missing data blobs
   previously resulted in a crash.

   https://github.com/restic/restic/pull/2668

 * Bugfix #2674: Add stricter prune error checks

   Additional checks were added to the prune command in order to improve resiliency to backend,
   hardware and/or networking issues. The checks now detect a few more cases where such outside
   factors could potentially cause data loss.

   https://github.com/restic/restic/pull/2674

 * Bugfix #2899: Fix possible crash in the progress bar of check --read-data

   We've fixed a possible crash while displaying the progress bar for the check --read-data
   command. The crash occurred when the length of the progress bar status exceeded the terminal
   width, which only happened for very narrow terminal windows.

   https://github.com/restic/restic/pull/2899
   https://forum.restic.net/t/restic-rclone-pcloud-connection-issues/2963/15

 * Change #2482: Remove vendored dependencies

   We've removed the vendored dependencies (in the subdir `vendor/`). When building restic, the
   Go compiler automatically fetches the dependencies. It will also cryptographically verify
   that the correct code has been fetched by using the hashes in `go.sum` (see the link to the
   documentation below).

   https://github.com/restic/restic/issues/2482
   https://golang.org/cmd/go/#hdr-Module_downloading_and_verification

 * Change #2546: Return exit code 3 when failing to backup all source data

   The backup command used to return a zero exit code as long as a snapshot could be created
   successfully, even if some of the source files could not be read (in which case the snapshot
   would contain the rest of the files).

   This made it hard for automation/scripts to detect failures/incomplete backups by looking at
   the exit code. Restic now returns the following exit codes for the backup command:

   - 0 when the command was successful - 1 when there was a fatal error (no snapshot created) - 3 when
   some source data could not be read (incomplete snapshot created)

   https://github.com/restic/restic/issues/956
   https://github.com/restic/restic/issues/2064
   https://github.com/restic/restic/issues/2526
   https://github.com/restic/restic/issues/2364
   https://github.com/restic/restic/pull/2546

 * Change #2600: Update dependencies, require Go >= 1.13

   Restic now requires Go to be at least 1.13. This allows simplifications in the build process and
   removing workarounds.

   This is also probably the last version of restic still supporting mounting repositories via
   fuse on macOS. The library we're using for fuse does not support macOS any more and osxfuse is not
   open source any more.

   https://github.com/bazil/fuse/issues/224
   https://github.com/osxfuse/osxfuse/issues/590
   https://github.com/restic/restic/pull/2600
   https://github.com/restic/restic/pull/2852
   https://github.com/restic/restic/pull/2927

 * Change #1597: Honor the --no-lock flag in the mount command

   The mount command now does not lock the repository if given the --no-lock flag. This allows to
   mount repositories which are archived on a read only backend/filesystem.

   https://github.com/restic/restic/issues/1597
   https://github.com/restic/restic/pull/2821

 * Enhancement #1570: Support specifying multiple host flags for various commands

   Previously commands didn't take more than one `--host` or `-H` argument into account, which
   could be limiting with e.g. the `forget` command.

   The `dump`, `find`, `forget`, `ls`, `mount`, `restore`, `snapshots`, `stats` and `tag`
   commands will now take into account multiple `--host` and `-H` flags.

   https://github.com/restic/restic/issues/1570

 * Enhancement #1680: Optimize `restic mount`

   We've optimized the FUSE implementation used within restic. `restic mount` is now more
   responsive and uses less memory.

   https://github.com/restic/restic/issues/1680
   https://github.com/restic/restic/pull/2587
   https://github.com/restic/restic/pull/2787

 * Enhancement #2072: Display snapshot date when using `restic find`

   Added the respective snapshot date to the output of `restic find`.

   https://github.com/restic/restic/issues/2072

 * Enhancement #2175: Allow specifying user and host when creating keys

   When adding a new key to the repository, the username and hostname for the new key can be
   specified on the command line. This allows overriding the defaults, for example if you would
   prefer to use the FQDN to identify the host or if you want to add keys for several different hosts
   without having to run the key add command on those hosts.

   https://github.com/restic/restic/issues/2175

 * Enhancement #2277: Add support for ppc64le

   Adds support for ppc64le, the processor architecture from IBM.

   https://github.com/restic/restic/issues/2277

 * Enhancement #2395: Ignore sync errors when operation not supported by local filesystem

   The local backend has been modified to work with filesystems which doesn't support the `sync`
   operation. This operation is normally used by restic to ensure that data files are fully
   written to disk before continuing.

   For these limited filesystems, saving a file in the backend would previously fail with an
   "operation not supported" error. This error is now ignored, which means that e.g. an SMB mount
   on macOS can now be used as storage location for a repository.

   https://github.com/restic/restic/issues/2395
   https://forum.restic.net/t/sync-errors-on-mac-over-smb/1859

 * Enhancement #2427: Add flag `--iexclude-file` to backup command

   The backup command now supports the flag `--iexclude-file` which is a case-insensitive
   version of `--exclude-file`.

   https://github.com/restic/restic/issues/2427
   https://github.com/restic/restic/pull/2898

 * Enhancement #2569: Support excluding files by their size

   The `backup` command now supports the `--exclude-larger-than` option to exclude files which
   are larger than the specified maximum size. This can for example be useful to exclude
   unimportant files with a large file size.

   https://github.com/restic/restic/issues/2569
   https://github.com/restic/restic/pull/2914

 * Enhancement #2571: Self-heal missing file parts during backup of unchanged files

   We've improved the resilience of restic to certain types of repository corruption.

   For files that are unchanged since the parent snapshot, the backup command now verifies that
   all parts of the files still exist in the repository. Parts that are missing, e.g. from a damaged
   repository, are backed up again. This verification was already run for files that were
   modified since the parent snapshot, but is now also done for unchanged files.

   Note that restic will not backup file parts that are referenced in the index but where the actual
   data is not present on disk, as this situation can only be detected by restic check. Please
   ensure that you run `restic check` regularly.

   https://github.com/restic/restic/issues/2571
   https://github.com/restic/restic/pull/2827

 * Enhancement #2858: Support filtering snapshots by tag and path in the stats command

   We've added filtering snapshots by `--tag tagList` and by `--path path` to the `stats`
   command. This includes filtering of only 'latest' snapshots or all snapshots in a repository.

   https://github.com/restic/restic/issues/2858
   https://github.com/restic/restic/pull/2859
   https://forum.restic.net/t/stats-for-a-host-and-filtered-snapshots/3020

 * Enhancement #323: Add command for copying snapshots between repositories

   We've added a copy command, allowing you to copy snapshots from one repository to another.

   Note that this process will have to read (download) and write (upload) the entire snapshot(s)
   due to the different encryption keys used on the source and destination repository. Also, the
   transferred files are not re-chunked, which may break deduplication between files already
   stored in the destination repo and files copied there using this command.

   To fully support deduplication between repositories when the copy command is used, the init
   command now supports the `--copy-chunker-params` option, which initializes the new
   repository with identical parameters for splitting files into chunks as an already existing
   repository. This allows copied snapshots to be equally deduplicated in both repositories.

   https://github.com/restic/restic/issues/323
   https://github.com/restic/restic/pull/2606
   https://github.com/restic/restic/pull/2928

 * Enhancement #551: Use optimized library for hash calculation of file chunks

   We've switched the library used to calculate the hashes of file chunks, which are used for
   deduplication, to the optimized Minio SHA-256 implementation.

   Depending on the CPU it improves the hashing throughput by 10-30%. Modern x86 CPUs with the SHA
   Extension should be about two to three times faster.

   https://github.com/restic/restic/issues/551
   https://github.com/restic/restic/pull/2709

 * Enhancement #2195: Simplify and improve restore performance

   Significantly improves restore performance of large files (i.e. 50M+):
   https://github.com/restic/restic/issues/2074
   https://forum.restic.net/t/restore-using-rclone-gdrive-backend-is-slow/1112/8
   https://forum.restic.net/t/degraded-restore-performance-s3-backend/1400

   Fixes "not enough cache capacity" error during restore:
   https://github.com/restic/restic/issues/2244

   NOTE: This new implementation does not guarantee order in which blobs are written to the target
   files and, for example, the last blob of a file can be written to the file before any of the
   preceeding file blobs. It is therefore possible to have gaps in the data written to the target
   files if restore fails or interrupted by the user.

   The implementation will try to preallocate space for the restored files on the filesystem to
   prevent file fragmentation. This ensures good read performance for large files, like for
   example VM images. If preallocating space is not supported by the filesystem, then this step is
   silently skipped.

   https://github.com/restic/restic/pull/2195
   https://github.com/restic/restic/pull/2893

 * Enhancement #2328: Improve speed of check command

   We've improved the check command to traverse trees only once independent of whether they are
   contained in multiple snapshots. The check command is now much faster for repositories with a
   large number of snapshots.

   https://github.com/restic/restic/issues/2284
   https://github.com/restic/restic/pull/2328

 * Enhancement #2423: Support user@domain parsing as user

   Added the ability for user@domain-like users to be authenticated over SFTP servers.

   https://github.com/restic/restic/pull/2423

 * Enhancement #2576: Improve the chunking algorithm

   We've updated the chunker library responsible for splitting files into smaller blocks. It
   should improve the chunking throughput by 5-15% depending on the CPU.

   https://github.com/restic/restic/issues/2820
   https://github.com/restic/restic/pull/2576
   https://github.com/restic/restic/pull/2845

 * Enhancement #2598: Improve speed of diff command

   We've improved the performance of the diff command when comparing snapshots with similar
   content. It should run up to twice as fast as before.

   https://github.com/restic/restic/pull/2598

 * Enhancement #2599: Slightly reduce memory usage of prune and stats commands

   The prune and the stats command kept directory identifiers in memory twice while searching for
   used blobs.

   https://github.com/restic/restic/pull/2599

 * Enhancement #2733: S3 backend: Add support for WebIdentityTokenFile

   We've added support for EKS IAM roles for service accounts feature to the S3 backend.

   https://github.com/restic/restic/issues/2703
   https://github.com/restic/restic/pull/2733

 * Enhancement #2773: Optimize handling of new index entries

   Restic now uses less memory for backups which add a lot of data, e.g. large initial backups. In
   addition, we've improved the stability in some edge cases.

   https://github.com/restic/restic/pull/2773

 * Enhancement #2781: Reduce memory consumption of in-memory index

   We've improved how the index is stored in memory. This change can reduce memory usage for large
   repositories by up to 50% (depending on the operation).

   https://github.com/restic/restic/pull/2781
   https://github.com/restic/restic/pull/2812

 * Enhancement #2786: Optimize `list blobs` command

   We've changed the implementation of `list blobs` which should be now a bit faster and consume
   almost no memory even for large repositories.

   https://github.com/restic/restic/pull/2786

 * Enhancement #2790: Optimized file access in restic mount

   Reading large (> 100GiB) files from restic mountpoints is now faster, and the speedup is
   greater for larger files.

   https://github.com/restic/restic/pull/2790

 * Enhancement #2840: Speed-up file deletion in forget, prune and rebuild-index

   We've sped up the file deletion for the commands forget, prune and rebuild-index, especially
   for remote repositories. Deletion was sequential before and is now run in parallel.

   https://github.com/restic/restic/pull/2840


Changelog for restic 0.9.6 (2019-11-22)
=======================================

The following sections list the changes in restic 0.9.6 relevant to
restic users. The changes are ordered by importance.

Summary
-------

 * Fix #2063: Allow absolute path for filename when backing up from stdin
 * Fix #2174: Save files with invalid timestamps
 * Fix #2249: Read fresh metadata for unmodified files
 * Fix #2301: Add upper bound for t in --read-data-subset=n/t
 * Fix #2321: Check errors when loading index files
 * Enh #2179: Use ctime when checking for file changes
 * Enh #2306: Allow multiple retries for interactive password input
 * Enh #2330: Make `--group-by` accept both singular and plural
 * Enh #2350: Add option to configure S3 region

Details
-------

 * Bugfix #2063: Allow absolute path for filename when backing up from stdin

   When backing up from stdin, handle directory path for `--stdin-filename`. This can be used to
   specify the full path for the backed-up file.

   https://github.com/restic/restic/issues/2063

 * Bugfix #2174: Save files with invalid timestamps

   When restic reads invalid timestamps (year is before 0000 or after 9999) it refused to read and
   archive the file. We've changed the behavior and will now save modified timestamps with the
   year set to either 0000 or 9999, the rest of the timestamp stays the same, so the file will be saved
   (albeit with a bogus timestamp).

   https://github.com/restic/restic/issues/2174
   https://github.com/restic/restic/issues/1173

 * Bugfix #2249: Read fresh metadata for unmodified files

   Restic took all metadata for files which were detected as unmodified, not taking into account
   changed metadata (ownership, mode). This is now corrected.

   https://github.com/restic/restic/issues/2249
   https://github.com/restic/restic/pull/2252

 * Bugfix #2301: Add upper bound for t in --read-data-subset=n/t

   256 is the effective maximum for t, but restic would allow larger values, leading to strange
   behavior.

   https://github.com/restic/restic/issues/2301
   https://github.com/restic/restic/pull/2304

 * Bugfix #2321: Check errors when loading index files

   Restic now checks and handles errors which occur when loading index files, the missing check
   leads to odd errors (and a stack trace printed to users) later. This was reported in the forum.

   https://github.com/restic/restic/pull/2321
   https://forum.restic.net/t/check-rebuild-index-prune/1848/13

 * Enhancement #2179: Use ctime when checking for file changes

   Previously, restic only checked a file's mtime (along with other non-timestamp metadata) to
   decide if a file has changed. This could cause restic to not notice that a file has changed (and
   therefore continue to store the old version, as opposed to the modified version) if something
   edits the file and then resets the timestamp. Restic now also checks the ctime of files, so any
   modifications to a file should be noticed, and the modified file will be backed up. The ctime
   check will be disabled if the --ignore-inode flag was given.

   If this change causes problems for you, please open an issue, and we can look in to adding a
   seperate flag to disable just the ctime check.

   https://github.com/restic/restic/issues/2179
   https://github.com/restic/restic/pull/2212

 * Enhancement #2306: Allow multiple retries for interactive password input

   Restic used to quit if the repository password was typed incorrectly once. Restic will now ask
   the user again for the repository password if typed incorrectly. The user will now get three
   tries to input the correct password before restic quits.

   https://github.com/restic/restic/issues/2306

 * Enhancement #2330: Make `--group-by` accept both singular and plural

   One can now use the values `host`/`hosts`, `path`/`paths` and `tag` / `tags` interchangeably
   in the `--group-by` argument.

   https://github.com/restic/restic/issues/2330

 * Enhancement #2350: Add option to configure S3 region

   We've added a new option for setting the region when accessing an S3-compatible service. For
   some providers, it is required to set this to a valid value. You can do that either by setting the
   environment variable `AWS_DEFAULT_REGION` or using the option `s3.region`, e.g. like this:
   `-o s3.region="us-east-1"`.

   https://github.com/restic/restic/pull/2350


Changelog for restic 0.9.5 (2019-04-23)
=======================================

The following sections list the changes in restic 0.9.5 relevant to
restic users. The changes are ordered by importance.

Summary
-------

 * Fix #2135: Return error when no bytes could be read from stdin
 * Fix #2181: Don't cancel timeout after 30 seconds for self-update
 * Fix #2203: Fix reading passwords from stdin
 * Fix #2224: Don't abort the find command when a tree can't be loaded
 * Enh #1895: Add case insensitive include & exclude options
 * Enh #1937: Support streaming JSON output for backup
 * Enh #2155: Add Openstack application credential auth for Swift
 * Enh #2184: Add --json support to forget command
 * Enh #2037: Add group-by option to snapshots command
 * Enh #2124: Ability to dump folders to tar via stdout
 * Enh #2139: Return error if no bytes could be read for `backup --stdin`
 * Enh #2205: Add --ignore-inode option to backup cmd
 * Enh #2220: Add config option to set S3 storage class

Details
-------

 * Bugfix #2135: Return error when no bytes could be read from stdin

   We assume that users reading backup data from stdin want to know when no data could be read, so now
   restic returns an error when `backup --stdin` is called but no bytes could be read. Usually,
   this means that an earlier command in a pipe has failed. The documentation was amended and now
   recommends setting the `pipefail` option (`set -o pipefail`).

   https://github.com/restic/restic/pull/2135
   https://github.com/restic/restic/pull/2139

 * Bugfix #2181: Don't cancel timeout after 30 seconds for self-update

   https://github.com/restic/restic/issues/2181

 * Bugfix #2203: Fix reading passwords from stdin

   Passwords for the `init`, `key add`, and `key passwd` commands can now be read from
   non-terminal stdin.

   https://github.com/restic/restic/issues/2203

 * Bugfix #2224: Don't abort the find command when a tree can't be loaded

   Change the find command so that missing trees don't result in a crash. Instead, the error is
   logged to the debug log, and the tree ID is displayed along with the snapshot it belongs to. This
   makes it possible to recover repositories that are missing trees by forgetting the snapshots
   they are used in.

   https://github.com/restic/restic/issues/2224

 * Enhancement #1895: Add case insensitive include & exclude options

   The backup and restore commands now have --iexclude and --iinclude flags as case insensitive
   variants of --exclude and --include.

   https://github.com/restic/restic/issues/1895
   https://github.com/restic/restic/pull/2032

 * Enhancement #1937: Support streaming JSON output for backup

   We've added support for getting machine-readable status output during backup, just pass the
   flag `--json` for `restic backup` and restic will output a stream of JSON objects which contain
   the current progress.

   https://github.com/restic/restic/issues/1937
   https://github.com/restic/restic/pull/1944

 * Enhancement #2155: Add Openstack application credential auth for Swift

   Since Openstack Queens Identity (auth V3) service supports an application credential auth
   method. It allows to create a technical account with the limited roles. This commit adds an
   application credential authentication method for the Swift backend.

   https://github.com/restic/restic/issues/2155

 * Enhancement #2184: Add --json support to forget command

   The forget command now supports the --json argument, outputting the information about what is
   (or would-be) kept and removed from the repository.

   https://github.com/restic/restic/issues/2184
   https://github.com/restic/restic/pull/2185

 * Enhancement #2037: Add group-by option to snapshots command

   We have added an option to group the output of the snapshots command, similar to the output of the
   forget command. The option has been called "--group-by" and accepts any combination of the
   values "host", "paths" and "tags", separated by commas. Default behavior (not specifying
   --group-by) has not been changed. We have added support of the grouping to the JSON output.

   https://github.com/restic/restic/issues/2037
   https://github.com/restic/restic/pull/2087

 * Enhancement #2124: Ability to dump folders to tar via stdout

   We've added the ability to dump whole folders to stdout via the `dump` command. Restic now
   requires at least Go 1.10 due to a limitation of the standard library for Go <= 1.9.

   https://github.com/restic/restic/issues/2123
   https://github.com/restic/restic/pull/2124

 * Enhancement #2139: Return error if no bytes could be read for `backup --stdin`

   When restic is used to backup the output of a program, like `mysqldump | restic backup --stdin`,
   it now returns an error if no bytes could be read at all. This catches the failure case when
   `mysqldump` failed for some reason and did not output any data to stdout.

   https://github.com/restic/restic/pull/2139

 * Enhancement #2205: Add --ignore-inode option to backup cmd

   This option handles backup of virtual filesystems that do not keep fixed inodes for files, like
   Fuse-based, pCloud, etc. Ignoring inode changes allows to consider the file as unchanged if
   last modification date and size are unchanged.

   https://github.com/restic/restic/issues/1631
   https://github.com/restic/restic/pull/2205
   https://github.com/restic/restic/pull/2047

 * Enhancement #2220: Add config option to set S3 storage class

   The `s3.storage-class` option can be passed to restic (using `-o`) to specify the storage
   class to be used for S3 objects created by restic.

   The storage class is passed as-is to S3, so it needs to be understood by the API. On AWS, it can be
   one of `STANDARD`, `STANDARD_IA`, `ONEZONE_IA`, `INTELLIGENT_TIERING` and
   `REDUCED_REDUNDANCY`. If unspecified, the default storage class is used (`STANDARD` on
   AWS).

   You can mix storage classes in the same bucket, and the setting isn't stored in the restic
   repository, so be sure to specify it with each command that writes to S3.

   https://github.com/restic/restic/issues/706
   https://github.com/restic/restic/pull/2220


Changelog for restic 0.9.4 (2019-01-06)
=======================================

The following sections list the changes in restic 0.9.4 relevant to
restic users. The changes are ordered by importance.

Summary
-------

 * Fix #1989: Google Cloud Storage: Respect bandwidth limit
 * Fix #2040: Add host name filter shorthand flag for `stats` command
 * Fix #2068: Correctly return error loading data
 * Fix #2095: Consistently use local time for snapshots times
 * Enh #1605: Concurrent restore
 * Enh #2089: Increase granularity of the "keep within" retention policy
 * Enh #2097: Add key hinting
 * Enh #2017: Mount: Enforce FUSE Unix permissions with allow-other
 * Enh #2070: Make all commands display timestamps in local time
 * Enh #2085: Allow --files-from to be specified multiple times
 * Enh #2094: Run command to get password

Details
-------

 * Bugfix #1989: Google Cloud Storage: Respect bandwidth limit

   The GCS backend did not respect the bandwidth limit configured, a previous commit
   accidentally removed support for it.

   https://github.com/restic/restic/issues/1989
   https://github.com/restic/restic/pull/2100

 * Bugfix #2040: Add host name filter shorthand flag for `stats` command

   The default value for `--host` flag was set to 'H' (the shorthand version of the flag), this
   caused the lookup for the latest snapshot to fail.

   Add shorthand flag `-H` for `--host` (with empty default so if these flags are not specified the
   latest snapshot will not filter by host name).

   Also add shorthand `-H` for `backup` command.

   https://github.com/restic/restic/issues/2040

 * Bugfix #2068: Correctly return error loading data

   In one case during `prune` and `check`, an error loading data from the backend is not returned
   properly. This is now corrected.

   https://github.com/restic/restic/issues/1999#issuecomment-433737921
   https://github.com/restic/restic/pull/2068

 * Bugfix #2095: Consistently use local time for snapshots times

   By default snapshots created with restic backup were set to local time, but when the --time flag
   was used the provided timestamp was parsed as UTC. With this change all snapshots times are set
   to local time.

   https://github.com/restic/restic/pull/2095

 * Enhancement #1605: Concurrent restore

   This change significantly improves restore performance, especially when using
   high-latency remote repositories like B2.

   The implementation now uses several concurrent threads to download and process multiple
   remote files concurrently. To further reduce restore time, each remote file is downloaded
   using a single repository request.

   https://github.com/restic/restic/issues/1605
   https://github.com/restic/restic/pull/1719

 * Enhancement #2089: Increase granularity of the "keep within" retention policy

   The `keep-within` option of the `forget` command now accepts time ranges with an hourly
   granularity. For example, running `restic forget --keep-within 3d12h` will keep all the
   snapshots made within three days and twelve hours from the time of the latest snapshot.

   https://github.com/restic/restic/issues/2089
   https://github.com/restic/restic/pull/2090

 * Enhancement #2097: Add key hinting

   Added a new option `--key-hint` and corresponding environment variable `RESTIC_KEY_HINT`.
   The key hint is a key ID to try decrypting first, before other keys in the repository.

   This change will benefit repositories with many keys; if the correct key hint is supplied then
   restic only needs to check one key. If the key hint is incorrect (the key does not exist, or the
   password is incorrect) then restic will check all keys, as usual.

   https://github.com/restic/restic/issues/2097

 * Enhancement #2017: Mount: Enforce FUSE Unix permissions with allow-other

   The fuse mount (`restic mount`) now lets the kernel check the permissions of the files within
   snapshots (this is done through the `DefaultPermissions` FUSE option) when the option
   `--allow-other` is specified.

   To restore the old behavior, we've added the `--no-default-permissions` option. This allows
   all users that have access to the mount point to access all files within the snapshots.

   https://github.com/restic/restic/pull/2017

 * Enhancement #2070: Make all commands display timestamps in local time

   Restic used to drop the timezone information from displayed timestamps, it now converts
   timestamps to local time before printing them so the times can be easily compared to.

   https://github.com/restic/restic/pull/2070

 * Enhancement #2085: Allow --files-from to be specified multiple times

   Before, restic took only the last file specified with `--files-from` into account, this is now
   corrected.

   https://github.com/restic/restic/issues/2085
   https://github.com/restic/restic/pull/2086

 * Enhancement #2094: Run command to get password

   We've added the `--password-command` option which allows specifying a command that restic
   runs every time the password for the repository is needed, so it can be integrated with a
   password manager or keyring. The option can also be set via the environment variable
   `$RESTIC_PASSWORD_COMMAND`.

   https://github.com/restic/restic/pull/2094


Changelog for restic 0.9.3 (2018-10-13)
=======================================

The following sections list the changes in restic 0.9.3 relevant to
restic users. The changes are ordered by importance.

Summary
-------

 * Fix #1935: Remove truncated files from cache
 * Fix #1978: Do not return an error when the scanner is slower than backup
 * Enh #1766: Restore: suppress lchown errors when not running as root
 * Enh #1909: Reject files/dirs by name first
 * Enh #1940: Add directory filter to ls command
 * Enh #1967: Use `--host` everywhere
 * Enh #2028: Display size of cache directories
 * Enh #1777: Improve the `find` command
 * Enh #1876: Display reason why forget keeps snapshots
 * Enh #1891: Accept glob in paths loaded via --files-from
 * Enh #1920: Vendor dependencies with Go 1.11 Modules
 * Enh #1949: Add new command `self-update`
 * Enh #1953: Ls: Add JSON output support for restic ls cmd
 * Enh #1962: Stream JSON output for ls command

Details
-------

 * Bugfix #1935: Remove truncated files from cache

   When a file in the local cache is truncated, and restic tries to access data beyond the end of the
   (cached) file, it used to return an error "EOF". This is now fixed, such truncated files are
   removed and the data is fetched directly from the backend.

   https://github.com/restic/restic/issues/1935

 * Bugfix #1978: Do not return an error when the scanner is slower than backup

   When restic makes a backup, there's a background task called "scanner" which collects
   information on how many files and directories are to be saved, in order to display progress
   information to the user. When the backup finishes faster than the scanner, it is aborted
   because the result is not needed any more. This logic contained a bug, where quitting the
   scanner process was treated as an error, and caused restic to print an unhelpful error message
   ("context canceled").

   https://github.com/restic/restic/issues/1978
   https://github.com/restic/restic/pull/1991

 * Enhancement #1766: Restore: suppress lchown errors when not running as root

   Like "cp" and "rsync" do, restic now only reports errors for changing the ownership of files
   during restore if it is run ￼as root, on non-Windows operating systems. On Windows, the error
   is reported as usual.

   https://github.com/restic/restic/issues/1766

 * Enhancement #1909: Reject files/dirs by name first

   The current scanner/archiver code had an architectural limitation: it always ran the
   `lstat()` system call on all files and directories before a decision to include/exclude the
   file/dir was made. This lead to a lot of unnecessary system calls for items that could have been
   rejected by their name or path only.

   We've changed the archiver/scanner implementation so that it now first rejects by name/path,
   and only runs the system call on the remaining items. This reduces the number of `lstat()`
   system calls a lot (depending on the exclude settings).

   https://github.com/restic/restic/issues/1909
   https://github.com/restic/restic/pull/1912

 * Enhancement #1940: Add directory filter to ls command

   The ls command can now be filtered by directories, so that only files in the given directories
   will be shown. If the --recursive flag is specified, then ls will traverse subfolders and list
   their files as well.

   It used to be possible to specify multiple snapshots, but that has been replaced by only one
   snapshot and the possibility of specifying multiple directories.

   Specifying directories constrains the walk, which can significantly speed up the listing.

   https://github.com/restic/restic/issues/1940
   https://github.com/restic/restic/pull/1941

 * Enhancement #1967: Use `--host` everywhere

   We now use the flag `--host` for all commands which need a host name, using `--hostname` (e.g.
   for `restic backup`) still works, but will print a deprecation warning. Also, add the short
   option `-H` where possible.

   https://github.com/restic/restic/issues/1967

 * Enhancement #2028: Display size of cache directories

   The `cache` command now by default shows the size of the individual cache directories. It can be
   disabled with `--no-size`.

   https://github.com/restic/restic/issues/2028
   https://github.com/restic/restic/pull/2033

 * Enhancement #1777: Improve the `find` command

   We've updated the `find` command to support multiple patterns.

   `restic find` is now able to list the snapshots containing a specific tree or blob, or even the
   snapshots that contain blobs belonging to a given pack. A list of IDs can be given, as long as they
   all have the same type.

   The command `find` can also display the pack IDs the blobs belong to, if the `--show-pack-id`
   flag is provided.

   https://github.com/restic/restic/issues/1777
   https://github.com/restic/restic/pull/1780

 * Enhancement #1876: Display reason why forget keeps snapshots

   We've added a column to the list of snapshots `forget` keeps which details the reasons to keep a
   particuliar snapshot. This makes debugging policies for forget much easier. Please remember
   to always try things out with `--dry-run`!

   https://github.com/restic/restic/pull/1876

 * Enhancement #1891: Accept glob in paths loaded via --files-from

   Before that, behaviour was different if paths were appended to command line or from a file,
   because wild card characters were expanded by shell if appended to command line, but not
   expanded if loaded from file.

   https://github.com/restic/restic/issues/1891

 * Enhancement #1920: Vendor dependencies with Go 1.11 Modules

   Until now, we've used `dep` for managing dependencies, we've now switch to using Go modules.
   For users this does not change much, only if you want to compile restic without downloading
   anything with Go 1.11, then you need to run: `go build -mod=vendor build.go`

   https://github.com/restic/restic/pull/1920

 * Enhancement #1949: Add new command `self-update`

   We have added a new command called `self-update` which downloads the latest released version
   of restic from GitHub and replaces the current binary with it. It does not rely on any external
   program (so it'll work everywhere), but still verifies the GPG signature using the embedded
   GPG public key.

   By default, the `self-update` command is hidden behind the `selfupdate` built tag, which is
   only set when restic is built using `build.go` (including official releases). The reason for
   this is that downstream distributions will then not include the command by default, so users
   are encouraged to use the platform-specific distribution mechanism.

   https://github.com/restic/restic/pull/1949

 * Enhancement #1953: Ls: Add JSON output support for restic ls cmd

   We've implemented listing files in the repository with JSON as output, just pass `--json` as an
   option to `restic ls`. This makes the output of the command machine readable.

   https://github.com/restic/restic/pull/1953

 * Enhancement #1962: Stream JSON output for ls command

   The `ls` command now supports JSON output with the global `--json` flag, and this change
   streams out JSON messages one object at a time rather than en entire array buffered in memory
   before encoding. The advantage is it allows large listings to be handled efficiently.

   Two message types are printed: snapshots and nodes. A snapshot object will precede node
   objects which belong to that snapshot. The `struct_type` field can be used to determine which
   kind of message an object is.

   https://github.com/restic/restic/pull/1962


Changelog for restic 0.9.2 (2018-08-06)
=======================================

The following sections list the changes in restic 0.9.2 relevant to
restic users. The changes are ordered by importance.

Summary
-------

 * Fix #1854: Allow saving files/dirs on different fs with `--one-file-system`
 * Fix #1870: Fix restore with --include
 * Fix #1880: Use `--cache-dir` argument for `check` command
 * Fix #1893: Return error when exclude file cannot be read
 * Fix #1861: Fix case-insensitive search with restic find
 * Enh #1906: Add support for B2 application keys
 * Enh #874: Add stats command to get information about a repository
 * Enh #1772: Add restore --verify to verify restored file content
 * Enh #1853: Add JSON output support to `restic key list`
 * Enh #1477: S3 backend: accept AWS_SESSION_TOKEN
 * Enh #1901: Update the Backblaze B2 library

Details
-------

 * Bugfix #1854: Allow saving files/dirs on different fs with `--one-file-system`

   Restic now allows saving files/dirs on a different file system in a subdir correctly even when
   `--one-file-system` is specified.

   The first thing the restic archiver code does is to build a tree of the target
   files/directories. If it detects that a parent directory is already included (e.g. `restic
   backup /foo /foo/bar/baz`), it'll ignore the latter argument.

   Without `--one-file-system`, that's perfectly valid: If `/foo` is to be archived, it will
   include `/foo/bar/baz`. But with `--one-file-system`, `/foo/bar/baz` may reside on a
   different file system, so it won't be included with `/foo`.

   https://github.com/restic/restic/issues/1854
   https://github.com/restic/restic/pull/1855

 * Bugfix #1870: Fix restore with --include

   We fixed a bug which prevented restic to restore files with an include filter.

   https://github.com/restic/restic/issues/1870
   https://github.com/restic/restic/pull/1900

 * Bugfix #1880: Use `--cache-dir` argument for `check` command

   `check` command now uses a temporary sub-directory of the specified directory if set using the
   `--cache-dir` argument. If not set, the cache directory is created in the default temporary
   directory as before. In either case a temporary cache is used to ensure the actual repository is
   checked (rather than a local copy).

   The `--cache-dir` argument was not used by the `check` command, instead a cache directory was
   created in the temporary directory.

   https://github.com/restic/restic/issues/1880

 * Bugfix #1893: Return error when exclude file cannot be read

   A bug was found: when multiple exclude files were passed to restic and one of them could not be
   read, an error was printed and restic continued, ignoring even the existing exclude files.
   Now, an error message is printed and restic aborts when an exclude file cannot be read.

   https://github.com/restic/restic/issues/1893

 * Bugfix #1861: Fix case-insensitive search with restic find

   We've fixed the behavior for `restic find -i PATTERN`, which was broken in v0.9.1.

   https://github.com/restic/restic/pull/1861

 * Enhancement #1906: Add support for B2 application keys

   Restic can now use so-called "application keys" which can be created in the B2 dashboard and
   were only introduced recently. In contrast to the "master key", such keys can be restricted to a
   specific bucket and/or path.

   https://github.com/restic/restic/issues/1906
   https://github.com/restic/restic/pull/1914

 * Enhancement #874: Add stats command to get information about a repository

   https://github.com/restic/restic/issues/874
   https://github.com/restic/restic/pull/1729

 * Enhancement #1772: Add restore --verify to verify restored file content

   Restore will print error message if restored file content does not match expected SHA256
   checksum

   https://github.com/restic/restic/pull/1772

 * Enhancement #1853: Add JSON output support to `restic key list`

   This PR enables users to get the output of `restic key list` in JSON in addition to the existing
   table format.

   https://github.com/restic/restic/pull/1853

 * Enhancement #1477: S3 backend: accept AWS_SESSION_TOKEN

   Before, it was not possible to use s3 backend with AWS temporary security credentials(with
   AWS_SESSION_TOKEN). This change gives higher priority to credentials.EnvAWS credentials
   provider.

   https://github.com/restic/restic/issues/1477
   https://github.com/restic/restic/pull/1479
   https://github.com/restic/restic/pull/1647

 * Enhancement #1901: Update the Backblaze B2 library

   We've updated the library we're using for accessing the Backblaze B2 service to 0.5.0 to
   include support for upcoming so-called "application keys". With this feature, you can create
   access credentials for B2 which are restricted to e.g. a single bucket or even a sub-directory
   of a bucket.

   https://github.com/restic/restic/pull/1901
   https://github.com/kurin/blazer


Changelog for restic 0.9.1 (2018-06-10)
=======================================

The following sections list the changes in restic 0.9.1 relevant to
restic users. The changes are ordered by importance.

Summary
-------

 * Fix #1801: Add limiting bandwidth to the rclone backend
 * Fix #1822: Allow uploading large files to MS Azure
 * Fix #1825: Correct `find` to not skip snapshots
 * Fix #1833: Fix caching files on error
 * Fix #1834: Resolve deadlock

Details
-------

 * Bugfix #1801: Add limiting bandwidth to the rclone backend

   The rclone backend did not respect `--limit-upload` or `--limit-download`. Oftentimes it's
   not necessary to use this, as the limiting in rclone itself should be used because it gives much
   better results, but in case a remote instance of rclone is used (e.g. called via ssh), it is still
   relevant to limit the bandwidth from restic to rclone.

   https://github.com/restic/restic/issues/1801

 * Bugfix #1822: Allow uploading large files to MS Azure

   Sometimes, restic creates files to be uploaded to the repository which are quite large, e.g.
   when saving directories with many entries or very large files. The MS Azure API does not allow
   uploading files larger that 256MiB directly, rather restic needs to upload them in blocks of
   100MiB. This is now implemented.

   https://github.com/restic/restic/issues/1822

 * Bugfix #1825: Correct `find` to not skip snapshots

   Under certain circumstances, the `find` command was found to skip snapshots containing
   directories with files to look for when the directories haven't been modified at all, and were
   already printed as part of a different snapshot. This is now corrected.

   In addition, we've switched to our own matching/pattern implementation, so now things like
   `restic find "/home/user/foo/**/main.go"` are possible.

   https://github.com/restic/restic/issues/1825
   https://github.com/restic/restic/issues/1823

 * Bugfix #1833: Fix caching files on error

   During `check` it may happen that different threads access the same file in the backend, which
   is then downloaded into the cache only once. When that fails, only the thread which is
   responsible for downloading the file signals the correct error. The other threads just assume
   that the file has been downloaded successfully and then get an error when they try to access the
   cached file.

   https://github.com/restic/restic/issues/1833

 * Bugfix #1834: Resolve deadlock

   When the "scanning" process restic runs to find out how much data there is does not finish before
   the backup itself is done, restic stops doing anything. This is resolved now.

   https://github.com/restic/restic/issues/1834
   https://github.com/restic/restic/pull/1835


Changelog for restic 0.9.0 (2018-05-21)
=======================================

The following sections list the changes in restic 0.9.0 relevant to
restic users. The changes are ordered by importance.

Summary
-------

 * Fix #1608: Respect time stamp for new backup when reading from stdin
 * Fix #1652: Ignore/remove invalid lock files
 * Fix #1730: Ignore sockets for restore
 * Fix #1684: Fix backend tests for rest-server
 * Fix #1745: Correctly parse the argument to --tls-client-cert
 * Enh #1433: Support UTF-16 encoding and process Byte Order Mark
 * Enh #1561: Allow using rclone to access other services
 * Enh #1665: Improve cache handling for `restic check`
 * Enh #1721: Add `cache` command to list cache dirs
 * Enh #1758: Allow saving OneDrive folders in Windows
 * Enh #549: Rework archiver code
 * Enh #1552: Use Google Application Default credentials
 * Enh #1477: Accept AWS_SESSION_TOKEN for the s3 backend
 * Enh #1648: Ignore AWS permission denied error when creating a repository
 * Enh #1649: Add illumos/Solaris support
 * Enh #1709: Improve messages `restic check` prints
 * Enh #827: Add --new-password-file flag for non-interactive password changes
 * Enh #1735: Allow keeping a time range of snaphots
 * Enh #1782: Use default AWS credentials chain for S3 backend

Details
-------

 * Bugfix #1608: Respect time stamp for new backup when reading from stdin

   When reading backups from stdin (via `restic backup --stdin`), restic now uses the time stamp
   for the new backup passed in `--time`.

   https://github.com/restic/restic/issues/1608
   https://github.com/restic/restic/pull/1703

 * Bugfix #1652: Ignore/remove invalid lock files

   This corrects a bug introduced recently: When an invalid lock file in the repo is encountered
   (e.g. if the file is empty), the code used to ignore that, but now returns the error. Now, invalid
   files are ignored for the normal lock check, and removed when `restic unlock --remove-all` is
   run.

   https://github.com/restic/restic/issues/1652
   https://github.com/restic/restic/pull/1653

 * Bugfix #1730: Ignore sockets for restore

   We've received a report and correct the behavior in which the restore code aborted restoring a
   directory when a socket was encountered. Unix domain socket files cannot be restored (they are
   created on the fly once a process starts listening). The error handling was corrected, and in
   addition we're now ignoring sockets during restore.

   https://github.com/restic/restic/issues/1730
   https://github.com/restic/restic/pull/1731

 * Bugfix #1684: Fix backend tests for rest-server

   The REST server for restic now requires an explicit parameter (`--no-auth`) if no
   authentication should be allowed. This is fixed in the tests.

   https://github.com/restic/restic/pull/1684

 * Bugfix #1745: Correctly parse the argument to --tls-client-cert

   Previously, the --tls-client-cert method attempt to read ARGV[1] (hardcoded) instead of the
   argument that was passed to it. This has been corrected.

   https://github.com/restic/restic/issues/1745
   https://github.com/restic/restic/pull/1746

 * Enhancement #1433: Support UTF-16 encoding and process Byte Order Mark

   On Windows, text editors commonly leave a Byte Order Mark at the beginning of the file to define
   which encoding is used (oftentimes UTF-16). We've added code to support processing the BOMs in
   text files, like the exclude files, the password file and the file passed via `--files-from`.
   This does not apply to any file being saved in a backup, those are not touched and archived as they
   are.

   https://github.com/restic/restic/issues/1433
   https://github.com/restic/restic/issues/1738
   https://github.com/restic/restic/pull/1748

 * Enhancement #1561: Allow using rclone to access other services

   We've added the ability to use rclone to store backup data on all backends that it supports. This
   was done in collaboration with Nick, the author of rclone. You can now use it to first configure a
   service, then restic manages the rest (starting and stopping rclone). For details, please see
   the manual.

   https://github.com/restic/restic/issues/1561
   https://github.com/restic/restic/pull/1657
   https://rclone.org

 * Enhancement #1665: Improve cache handling for `restic check`

   For safety reasons, restic does not use a local metadata cache for the `restic check` command,
   so that data is loaded from the repository and restic can check it's in good condition. When the
   cache is disabled, restic will fetch each tiny blob needed for checking the integrity using a
   separate backend request. For non-local backends, that will take a long time, and depending on
   the backend (e.g. B2) may also be much more expensive.

   This PR adds a few commits which will change the behavior as follows:

   * When `restic check` is called without any additional parameters, it will build a new cache in a
   temporary directory, which is removed at the end of the check. This way, we'll get readahead for
   metadata files (so restic will fetch the whole file when the first blob from the file is
   requested), but all data is freshly fetched from the storage backend. This is the default
   behavior and will work for almost all users.

   * When `restic check` is called with `--with-cache`, the default on-disc cache is used. This
   behavior hasn't changed since the cache was introduced.

   * When `--no-cache` is specified, restic falls back to the old behavior, and read all tiny blobs
   in separate requests.

   https://github.com/restic/restic/issues/1665
   https://github.com/restic/restic/issues/1694
   https://github.com/restic/restic/pull/1696

 * Enhancement #1721: Add `cache` command to list cache dirs

   The command `cache` was added, it allows listing restic's cache directoriers together with
   the last usage. It also allows removing old cache dirs without having to access a repo, via
   `restic cache --cleanup`

   https://github.com/restic/restic/issues/1721
   https://github.com/restic/restic/pull/1749

 * Enhancement #1758: Allow saving OneDrive folders in Windows

   Restic now contains a bugfix to two libraries, which allows saving OneDrive folders in
   Windows. In order to use the newer versions of the libraries, the minimal version required to
   compile restic is now Go 1.9.

   https://github.com/restic/restic/issues/1758
   https://github.com/restic/restic/pull/1765

 * Enhancement #549: Rework archiver code

   The core archiver code and the complementary code for the `backup` command was rewritten
   completely. This resolves very annoying issues such as 549. The first backup with this release
   of restic will likely result in all files being re-read locally, so it will take a lot longer. The
   next backup after that will be fast again.

   Basically, with the old code, restic took the last path component of each to-be-saved file or
   directory as the top-level file/directory within the snapshot. This meant that when called as
   `restic backup /home/user/foo`, the snapshot would contain the files in the directory
   `/home/user/foo` as `/foo`.

   This is not the case any more with the new archiver code. Now, restic works very similar to what
   `tar` does: When restic is called with an absolute path to save, then it'll preserve the
   directory structure within the snapshot. For the example above, the snapshot would contain
   the files in the directory within `/home/user/foo` in the snapshot. For relative
   directories, it only preserves the relative path components. So `restic backup user/foo`
   will save the files as `/user/foo` in the snapshot.

   While we were at it, the status display and notification system was completely rewritten. By
   default, restic now shows which files are currently read (unless `--quiet` is specified) in a
   multi-line status display.

   The `backup` command also gained a new option: `--verbose`. It can be specified once (which
   prints a bit more detail what restic is doing) or twice (which prints a line for each
   file/directory restic encountered, together with some statistics).

   Another issue that was resolved is the new code only reads two files at most. The old code would
   read way too many files in parallel, thereby slowing down the backup process on spinning discs a
   lot.

   https://github.com/restic/restic/issues/549
   https://github.com/restic/restic/issues/1286
   https://github.com/restic/restic/issues/446
   https://github.com/restic/restic/issues/1344
   https://github.com/restic/restic/issues/1416
   https://github.com/restic/restic/issues/1456
   https://github.com/restic/restic/issues/1145
   https://github.com/restic/restic/issues/1160
   https://github.com/restic/restic/pull/1494

 * Enhancement #1552: Use Google Application Default credentials

   Google provide libraries to generate appropriate credentials with various fallback
   sources. This change uses the library to generate our GCS client, which allows us to make use of
   these extra methods.

   This should be backward compatible with previous restic behaviour while adding the
   additional capabilities to auth from Google's internal metadata endpoints. For users
   running restic in GCP this can make authentication far easier than it was before.

   https://github.com/restic/restic/pull/1552
   https://developers.google.com/identity/protocols/application-default-credentials

 * Enhancement #1477: Accept AWS_SESSION_TOKEN for the s3 backend

   Before, it was not possible to use s3 backend with AWS temporary security credentials(with
   AWS_SESSION_TOKEN). This change gives higher priority to credentials.EnvAWS credentials
   provider.

   https://github.com/restic/restic/issues/1477
   https://github.com/restic/restic/pull/1479
   https://github.com/restic/restic/pull/1647

 * Enhancement #1648: Ignore AWS permission denied error when creating a repository

   It's not possible to use s3 backend scoped to a subdirectory(with specific permissions).
   Restic doesn't try to create repository in a subdirectory, when 'bucket exists' of parent
   directory check fails due to permission issues.

   https://github.com/restic/restic/pull/1648

 * Enhancement #1649: Add illumos/Solaris support

   https://github.com/restic/restic/pull/1649

 * Enhancement #1709: Improve messages `restic check` prints

   Some messages `restic check` prints are not really errors, so from now on restic does not treat
   them as errors any more and exits cleanly.

   https://github.com/restic/restic/pull/1709
   https://forum.restic.net/t/what-is-the-standard-procedure-to-follow-if-a-backup-or-restore-is-interrupted/571/2

 * Enhancement #827: Add --new-password-file flag for non-interactive password changes

   This makes it possible to change a repository password without being prompted.

   https://github.com/restic/restic/issues/827
   https://github.com/restic/restic/pull/1720
   https://forum.restic.net/t/changing-repo-password-without-prompt/591

 * Enhancement #1735: Allow keeping a time range of snaphots

   We've added the `--keep-within` option to the `forget` command. It instructs restic to keep
   all snapshots within the given duration since the newest snapshot. For example, running
   `restic forget --keep-within 5m7d` will keep all snapshots which have been made in the five
   months and seven days since the latest snapshot.

   https://github.com/restic/restic/pull/1735

 * Enhancement #1782: Use default AWS credentials chain for S3 backend

   Adds support for file credentials to the S3 backend (e.g. ~/.aws/credentials), and reorders
   the credentials chain for the S3 backend to match AWS's standard, which is static credentials,
   env vars, credentials file, and finally remote.

   https://github.com/restic/restic/pull/1782


Changelog for restic 0.8.3 (2018-02-26)
=======================================

The following sections list the changes in restic 0.8.3 relevant to
restic users. The changes are ordered by importance.

Summary
-------

 * Fix #1633: Fixed unexpected 'pack file cannot be listed' error
 * Fix #1641: Ignore files with invalid names in the repo
 * Fix #1638: Handle errors listing files in the backend
 * Enh #1497: Add --read-data-subset flag to check command
 * Enh #1560: Retry all repository file download errors
 * Enh #1623: Don't check for presence of files in the backend before writing
 * Enh #1634: Upgrade B2 client library, reduce HTTP requests

Details
-------

 * Bugfix #1633: Fixed unexpected 'pack file cannot be listed' error

   Due to a regression introduced in 0.8.2, the `rebuild-index` and `prune` commands failed to
   read pack files with size of 587, 588, 589 or 590 bytes.

   https://github.com/restic/restic/issues/1633
   https://github.com/restic/restic/pull/1635

 * Bugfix #1641: Ignore files with invalid names in the repo

   The release 0.8.2 introduced a bug: when restic encounters files in the repo which do not have a
   valid name, it tries to load a file with a name of lots of zeroes instead of ignoring it. This is now
   resolved, invalid file names are just ignored.

   https://github.com/restic/restic/issues/1641
   https://github.com/restic/restic/pull/1643
   https://forum.restic.net/t/help-fixing-repo-no-such-file/485/3

 * Bugfix #1638: Handle errors listing files in the backend

   A user reported in the forum that restic completes a backup although a concurrent `prune`
   operation was running. A few error messages were printed, but the backup was attempted and
   completed successfully. No error code was returned.

   This should not happen: The repository is exclusively locked during `prune`, so when `restic
   backup` is run in parallel, it should abort and return an error code instead.

   It was found that the bug was in the code introduced only recently, which retries a List()
   operation on the backend should that fail. It is now corrected.

   https://github.com/restic/restic/pull/1638
   https://forum.restic.net/t/restic-backup-returns-0-exit-code-when-already-locked/484

 * Enhancement #1497: Add --read-data-subset flag to check command

   This change introduces ability to check integrity of a subset of repository data packs. This
   can be used to spread integrity check of larger repositories over a period of time.

   https://github.com/restic/restic/issues/1497
   https://github.com/restic/restic/pull/1556

 * Enhancement #1560: Retry all repository file download errors

   Restic will now retry failed downloads, similar to other operations.

   https://github.com/restic/restic/pull/1560

 * Enhancement #1623: Don't check for presence of files in the backend before writing

   Before, all backend implementations were required to return an error if the file that is to be
   written already exists in the backend. For most backends, that means making a request (e.g. via
   HTTP) and returning an error when the file already exists.

   This is not accurate, the file could have been created between the HTTP request testing for it,
   and when writing starts, so we've relaxed this requeriment, which saves one additional HTTP
   request per newly added file.

   https://github.com/restic/restic/pull/1623

 * Enhancement #1634: Upgrade B2 client library, reduce HTTP requests

   We've upgraded the B2 client library restic uses to access BackBlaze B2. This reduces the
   number of HTTP requests needed to upload a new file from two to one, which should improve
   throughput to B2.

   https://github.com/restic/restic/pull/1634


Changelog for restic 0.8.2 (2018-02-17)
=======================================

The following sections list the changes in restic 0.8.2 relevant to
restic users. The changes are ordered by importance.

Summary
-------

 * Fix #1506: Limit bandwith at the http.RoundTripper for HTTP based backends
 * Fix #1512: Restore directory permissions as the last step
 * Fix #1528: Correctly create missing subdirs in data/
 * Fix #1590: Strip spaces for lines read via --files-from
 * Fix #1589: Complete intermediate index upload
 * Fix #1594: Google Cloud Storage: Use generic HTTP transport
 * Fix #1595: Backup: Remove bandwidth display
 * Enh #1522: Add support for TLS client certificate authentication
 * Enh #1541: Reduce number of remote requests during repository check
 * Enh #1567: Reduce number of backend requests for rebuild-index and prune
 * Enh #1507: Only reload snapshots once per minute for fuse mount
 * Enh #1538: Reduce memory allocations for querying the index
 * Enh #1549: Speed up querying across indices and scanning existing files
 * Enh #1554: Fuse/mount: Correctly handle EOF, add template option
 * Enh #1564: Don't terminate ssh on SIGINT
 * Enh #1579: Retry Backend.List() in case of errors
 * Enh #1584: Limit index file size

Details
-------

 * Bugfix #1506: Limit bandwith at the http.RoundTripper for HTTP based backends

   https://github.com/restic/restic/issues/1506
   https://github.com/restic/restic/pull/1511

 * Bugfix #1512: Restore directory permissions as the last step

   This change allows restoring into directories that were not writable during backup. Before,
   restic created the directory, set the read-only mode and then failed to create files in the
   directory. This change now restores the directory (with its permissions) as the very last
   step.

   https://github.com/restic/restic/issues/1512
   https://github.com/restic/restic/pull/1536

 * Bugfix #1528: Correctly create missing subdirs in data/

   https://github.com/restic/restic/issues/1528
   https://github.com/restic/restic/pull/1529

 * Bugfix #1590: Strip spaces for lines read via --files-from

   Leading and trailing spaces in lines read via `--files-from` are now stripped, so it behaves
   the same as with lines read via `--exclude-file`.

   https://github.com/restic/restic/issues/1590
   https://github.com/restic/restic/pull/1613

 * Bugfix #1589: Complete intermediate index upload

   After a user posted a comprehensive report of what he observed, we were able to find a bug and
   correct it: During backup, restic uploads so-called "intermediate" index files. When the
   backup finishes during a transfer of such an intermediate index, the upload is cancelled, but
   the backup is finished without an error. This leads to an inconsistent state, where the
   snapshot references data that is contained in the repo, but is not referenced in any index.

   The situation can be resolved by building a new index with `rebuild-index`, but looks very
   confusing at first. Since all the data got uploaded to the repo successfully, there was no risk
   of data loss, just minor inconvenience for our users.

   https://github.com/restic/restic/pull/1589
   https://forum.restic.net/t/error-loading-tree-check-prune-and-forget-gives-error-b2-backend/406

 * Bugfix #1594: Google Cloud Storage: Use generic HTTP transport

   It was discovered that the Google Cloud Storage backend did not use the generic HTTP transport,
   so things such as bandwidth limiting with `--limit-upload` did not work. This is resolved now.

   https://github.com/restic/restic/pull/1594

 * Bugfix #1595: Backup: Remove bandwidth display

   This commit removes the bandwidth displayed during backup process. It is misleading and
   seldomly correct, because it's neither the "read bandwidth" (only for the very first backup)
   nor the "upload bandwidth". Many users are confused about (and rightly so), c.f. #1581, #1033,
   #1591

   We'll eventually replace this display with something more relevant when the new archiver code
   is ready.

   https://github.com/restic/restic/pull/1595

 * Enhancement #1522: Add support for TLS client certificate authentication

   Support has been added for using a TLS client certificate for authentication to HTTP based
   backend. A file containing the PEM encoded private key and certificate can be set using the
   `--tls-client-cert` option.

   https://github.com/restic/restic/issues/1522
   https://github.com/restic/restic/pull/1524

 * Enhancement #1541: Reduce number of remote requests during repository check

   This change eliminates redundant remote repository calls and significantly improves
   repository check time.

   https://github.com/restic/restic/issues/1541
   https://github.com/restic/restic/pull/1548

 * Enhancement #1567: Reduce number of backend requests for rebuild-index and prune

   We've found a way to reduce then number of backend requests for the `rebuild-index` and `prune`
   operations. This significantly speeds up the operations for high-latency backends.

   https://github.com/restic/restic/issues/1567
   https://github.com/restic/restic/pull/1574
   https://github.com/restic/restic/pull/1575

 * Enhancement #1507: Only reload snapshots once per minute for fuse mount

   https://github.com/restic/restic/pull/1507

 * Enhancement #1538: Reduce memory allocations for querying the index

   This change reduces the internal memory allocations when the index data structures in memory
   are queried if a blob (part of a file) already exists in the repo. It should speed up backup a bit,
   and maybe even reduce RAM usage.

   https://github.com/restic/restic/pull/1538

 * Enhancement #1549: Speed up querying across indices and scanning existing files

   This change increases the whenever a blob (part of a file) is searched for in a restic
   repository. This will reduce cpu usage some when backing up files already backed up by restic.
   Cpu usage is further decreased when scanning files.

   https://github.com/restic/restic/pull/1549

 * Enhancement #1554: Fuse/mount: Correctly handle EOF, add template option

   We've added the `--snapshot-template` string, which can be used to specify a template for a
   snapshot directory. In addition, accessing data after the end of a file via the fuse mount is now
   handled correctly.

   https://github.com/restic/restic/pull/1554

 * Enhancement #1564: Don't terminate ssh on SIGINT

   We've reworked the code which runs the `ssh` login for the sftp backend so that it can prompt for a
   password (if needed) but does not exit when the user presses CTRL+C (SIGINT) e.g. during
   backup. This allows restic to properly shut down when it receives SIGINT and remove the lock
   file from the repo, afterwards exiting the `ssh` process.

   https://github.com/restic/restic/pull/1564
   https://github.com/restic/restic/pull/1588

 * Enhancement #1579: Retry Backend.List() in case of errors

   https://github.com/restic/restic/pull/1579

 * Enhancement #1584: Limit index file size

   Before, restic would create a single new index file on `prune` or `rebuild-index`, this may
   lead to memory problems when this huge index is created and loaded again. We're now limiting the
   size of the index file, and split newly created index files into several smaller ones. This
   allows restic to be more memory-efficient.

   https://github.com/restic/restic/issues/1412
   https://github.com/restic/restic/issues/979
   https://github.com/restic/restic/issues/526
   https://github.com/restic/restic/pull/1584


Changelog for restic 0.8.1 (2017-12-27)
=======================================

The following sections list the changes in restic 0.8.1 relevant to
restic users. The changes are ordered by importance.

Summary
-------

 * Fix #1457: Improve s3 backend with DigitalOcean Spaces
 * Fix #1454: Correct cache dir location for Windows and Darwin
 * Fix #1459: Disable handling SIGPIPE
 * Chg #1452: Do not save atime by default
 * Enh #1436: Add code to detect old cache directories
 * Enh #1439: Improve cancellation logic
 * Enh #11: Add the `diff` command

Details
-------

 * Bugfix #1457: Improve s3 backend with DigitalOcean Spaces

   https://github.com/restic/restic/issues/1457
   https://github.com/restic/restic/pull/1459

 * Bugfix #1454: Correct cache dir location for Windows and Darwin

   The cache directory on Windows and Darwin was not correct, instead the directory `.cache` was
   used.

   https://github.com/restic/restic/pull/1454

 * Bugfix #1459: Disable handling SIGPIPE

   We've disabled handling SIGPIPE again. Turns out, writing to broken TCP connections also
   raised SIGPIPE, so restic exits on the first write to a broken connection. Instead, restic
   should retry the request.

   https://github.com/restic/restic/issues/1457
   https://github.com/restic/restic/issues/1466
   https://github.com/restic/restic/pull/1459

 * Change #1452: Do not save atime by default

   By default, the access time for files and dirs is not saved any more. It is not possible to
   reliably disable updating the access time during a backup, so for the next backup the access
   time is different again. This means a lot of metadata is saved. If you want to save the access time
   anyway, pass `--with-atime` to the `backup` command.

   https://github.com/restic/restic/pull/1452

 * Enhancement #1436: Add code to detect old cache directories

   We've added code to detect old cache directories of repositories that haven't been used in a
   long time, restic now prints a note when it detects that such dirs exist. Also, the option
   `--cleanup-cache` was added to automatically remove such directories. That's not a problem
   because the cache will be rebuild once a repo is accessed again.

   https://github.com/restic/restic/pull/1436

 * Enhancement #1439: Improve cancellation logic

   The cancellation logic was improved, restic can now shut down cleanly when requested to do so
   (e.g. via ctrl+c).

   https://github.com/restic/restic/pull/1439

 * Enhancement #11: Add the `diff` command

   The command `diff` was added, it allows comparing two snapshots and listing all differences.

   https://github.com/restic/restic/issues/11
   https://github.com/restic/restic/issues/1460
   https://github.com/restic/restic/pull/1462


Changelog for restic 0.8.0 (2017-11-26)
=======================================

The following sections list the changes in restic 0.8.0 relevant to
restic users. The changes are ordered by importance.

Summary
-------

 * Sec #1445: Prevent writing outside the target directory during restore
 * Fix #1256: Re-enable workaround for S3 backend
 * Fix #1291: Reuse backend TCP connections to BackBlaze B2
 * Fix #1317: Run prune when `forget --prune` is called with just snapshot IDs
 * Fix #1437: Remove implicit path `/restic` for the s3 backend
 * Enh #1102: Add subdirectory `ids` to fuse mount
 * Enh #1114: Add `--cacert` to specify TLS certificates to check against
 * Enh #1216: Add upload/download limiting
 * Enh #1271: Cache results for excludes for `backup`
 * Enh #1274: Add `generate` command, replaces `manpage` and `autocomplete`
 * Enh #1367: Allow comments in files read from via `--file-from`
 * Enh #448: Sftp backend prompts for password
 * Enh #510: Add `dump` command
 * Enh #1040: Add local metadata cache
 * Enh #1249: Add `latest` symlink in fuse mount
 * Enh #1269: Add `--compact` to `forget` command
 * Enh #1281: Google Cloud Storage backend needs less permissions
 * Enh #1319: Make `check` print `no errors found` explicitly
 * Enh #1353: Retry failed backend requests

Details
-------

 * Security #1445: Prevent writing outside the target directory during restore

   A vulnerability was found in the restic restorer, which allowed attackers in special
   circumstances to restore files to a location outside of the target directory. Due to the
   circumstances we estimate this to be a low-risk vulnerability, but urge all users to upgrade to
   the latest version of restic.

   Exploiting the vulnerability requires a Linux/Unix system which saves backups via restic and
   a Windows systems which restores files from the repo. In addition, the attackers need to be able
   to create files with arbitrary names which are then saved to the restic repo. For example, by
   creating a file named "..\test.txt" (which is a perfectly legal filename on Linux) and
   restoring a snapshot containing this file on Windows, it would be written to the parent of the
   target directory.

   We'd like to thank Tyler Spivey for reporting this responsibly!

   https://github.com/restic/restic/pull/1445

 * Bugfix #1256: Re-enable workaround for S3 backend

   We've re-enabled a workaround for `minio-go` (the library we're using to access s3 backends),
   this reduces memory usage.

   https://github.com/restic/restic/issues/1256
   https://github.com/restic/restic/pull/1267

 * Bugfix #1291: Reuse backend TCP connections to BackBlaze B2

   A bug was discovered in the library we're using to access Backblaze, it now reuses already
   established TCP connections which should be a lot faster and not cause network failures any
   more.

   https://github.com/restic/restic/issues/1291
   https://github.com/restic/restic/pull/1301

 * Bugfix #1317: Run prune when `forget --prune` is called with just snapshot IDs

   A bug in the `forget` command caused `prune` not to be run when `--prune` was specified without a
   policy, e.g. when only snapshot IDs that should be forgotten are listed manually.

   https://github.com/restic/restic/pull/1317

 * Bugfix #1437: Remove implicit path `/restic` for the s3 backend

   The s3 backend used the subdir `restic` within a bucket if no explicit path after the bucket name
   was specified. Since this version, restic does not use this default path any more. If you
   created a repo on s3 in a bucket without specifying a path within the bucket, you need to add
   `/restic` at the end of the repository specification to access your repo:
   `s3:s3.amazonaws.com/bucket/restic`

   https://github.com/restic/restic/issues/1292
   https://github.com/restic/restic/pull/1437

 * Enhancement #1102: Add subdirectory `ids` to fuse mount

   The fuse mount now has an `ids` subdirectory which contains the snapshots below their (short)
   IDs.

   https://github.com/restic/restic/issues/1102
   https://github.com/restic/restic/pull/1299
   https://github.com/restic/restic/pull/1320

 * Enhancement #1114: Add `--cacert` to specify TLS certificates to check against

   We've added the `--cacert` option which can be used to pass one (or more) CA certificates to
   restic. These are used in addition to the system CA certificates to verify HTTPS certificates
   (e.g. for the REST backend).

   https://github.com/restic/restic/issues/1114
   https://github.com/restic/restic/pull/1276

 * Enhancement #1216: Add upload/download limiting

   We've added support for rate limiting through `--limit-upload` and `--limit-download`
   flags.

   https://github.com/restic/restic/issues/1216
   https://github.com/restic/restic/pull/1336
   https://github.com/restic/restic/pull/1358

 * Enhancement #1271: Cache results for excludes for `backup`

   The `backup` command now caches the result of excludes for a directory.

   https://github.com/restic/restic/issues/1271
   https://github.com/restic/restic/pull/1326

 * Enhancement #1274: Add `generate` command, replaces `manpage` and `autocomplete`

   The `generate` command has been added, which replaces the now removed commands `manpage` and
   `autocomplete`. This release of restic contains the most recent manpages in `doc/man` and the
   auto-completion files for bash and zsh in `doc/bash-completion.sh` and
   `doc/zsh-completion.zsh`

   https://github.com/restic/restic/issues/1274
   https://github.com/restic/restic/pull/1282

 * Enhancement #1367: Allow comments in files read from via `--file-from`

   When the list of files/dirs to be saved is read from a file with `--files-from`, comment lines
   (starting with `#`) are now ignored.

   https://github.com/restic/restic/issues/1367
   https://github.com/restic/restic/pull/1368

 * Enhancement #448: Sftp backend prompts for password

   The sftp backend now prompts for the password if a password is necessary for login.

   https://github.com/restic/restic/issues/448
   https://github.com/restic/restic/pull/1270

 * Enhancement #510: Add `dump` command

   We've added the `dump` command which prints a file from a snapshot to stdout. This can e.g. be
   used to restore files read with `backup --stdin`.

   https://github.com/restic/restic/issues/510
   https://github.com/restic/restic/pull/1346

 * Enhancement #1040: Add local metadata cache

   We've added a local cache for metadata so that restic doesn't need to load all metadata
   (snapshots, indexes, ...) from the repo each time it starts. By default the cache is active, but
   there's a new global option `--no-cache` that can be used to disable the cache. By deafult, the
   cache a standard cache folder for the OS, which can be overridden with `--cache-dir`. The cache
   will automatically populate, indexes and snapshots are saved as they are loaded. Cache
   directories for repos that haven't been used recently can automatically be removed by restic
   with the `--cleanup-cache` option.

   A related change was to by default create pack files in the repo that contain either data or
   metadata, not both mixed together. This allows easy caching of only the metadata files. The
   next run of `restic prune` will untangle mixed files automatically.

   https://github.com/restic/restic/issues/29
   https://github.com/restic/restic/issues/738
   https://github.com/restic/restic/issues/282
   https://github.com/restic/restic/pull/1040
   https://github.com/restic/restic/pull/1287
   https://github.com/restic/restic/pull/1436
   https://github.com/restic/restic/pull/1265

 * Enhancement #1249: Add `latest` symlink in fuse mount

   The directory structure in the fuse mount now exposes a symlink `latest` which points to the
   latest snapshot in that particular directory.

   https://github.com/restic/restic/pull/1249

 * Enhancement #1269: Add `--compact` to `forget` command

   The option `--compact` was added to the `forget` command to provide the same compact view as the
   `snapshots` command.

   https://github.com/restic/restic/pull/1269

 * Enhancement #1281: Google Cloud Storage backend needs less permissions

   The Google Cloud Storage backend no longer requires the service account to have the
   `storage.buckets.get` permission ("Storage Admin" role) in `restic init` if the bucket
   already exists.

   https://github.com/restic/restic/pull/1281

 * Enhancement #1319: Make `check` print `no errors found` explicitly

   The `check` command now explicetly prints `No errors were found` when no errors could be found.

   https://github.com/restic/restic/issues/1303
   https://github.com/restic/restic/pull/1319

 * Enhancement #1353: Retry failed backend requests

   https://github.com/restic/restic/pull/1353


Changelog for restic 0.7.3 (2017-09-20)
=======================================

The following sections list the changes in restic 0.7.3 relevant to
restic users. The changes are ordered by importance.

Summary
-------

 * Fix #1246: List all files stored in Google Cloud Storage

Details
-------

 * Bugfix #1246: List all files stored in Google Cloud Storage

   For large backups stored in Google Cloud Storage, the `prune` command fails because listing
   only returns the first 1000 files. This has been corrected, no data is lost in the process. In
   addition, a plausibility check was added to `prune`.

   https://github.com/restic/restic/issues/1246
   https://github.com/restic/restic/pull/1247


Changelog for restic 0.7.2 (2017-09-13)
=======================================

The following sections list the changes in restic 0.7.2 relevant to
restic users. The changes are ordered by importance.

Summary
-------

 * Fix #1167: Do not create a local repo unless `init` is used
 * Fix #1164: Make the `key remove` command behave as documented
 * Fix #1191: Make sure to write profiling files on interrupt
 * Enh #1132: Make `key` command always prompt for a password
 * Enh #1179: Resolve name conflicts, append a counter
 * Enh #1218: Add `--compact` to `snapshots` command
 * Enh #317: Add `--exclude-caches` and `--exclude-if-present`
 * Enh #697: Automatically generate man pages for all restic commands
 * Enh #1044: Improve `restore`, do not traverse/load excluded directories
 * Enh #1061: Add Dockerfile and official Docker image
 * Enh #1126: Use the standard Go git repository layout, use `dep` for vendoring
 * Enh #1134: Add support for storing backups on Google Cloud Storage
 * Enh #1144: Properly report errors when reading files with exclude patterns
 * Enh #1149: Add support for storing backups on Microsoft Azure Blob Storage
 * Enh #1196: Add `--group-by` to `forget` command for flexible grouping
 * Enh #1203: Print stats on all BSD systems when SIGINFO (ctrl+t) is received
 * Enh #1205: Allow specifying time/date for a backup with `--time`

Details
-------

 * Bugfix #1167: Do not create a local repo unless `init` is used

   When a restic command other than `init` is used with a local repository and the repository
   directory does not exist, restic creates the directory structure. That's an error, only the
   `init` command should create the dir.

   https://github.com/restic/restic/issues/1167
   https://github.com/restic/restic/pull/1182

 * Bugfix #1164: Make the `key remove` command behave as documented

   https://github.com/restic/restic/pull/1164

 * Bugfix #1191: Make sure to write profiling files on interrupt

   Since a few releases restic had the ability to write profiling files for memory and CPU usage
   when `debug` is enabled. It was discovered that when restic is interrupted (ctrl+c is
   pressed), the proper shutdown hook is not run. This is now corrected.

   https://github.com/restic/restic/pull/1191

 * Enhancement #1132: Make `key` command always prompt for a password

   The `key` command now prompts for a password even if the original password to access a repo has
   been specified via the `RESTIC_PASSWORD` environment variable or a password file.

   https://github.com/restic/restic/issues/1132
   https://github.com/restic/restic/pull/1133

 * Enhancement #1179: Resolve name conflicts, append a counter

   https://github.com/restic/restic/issues/1179
   https://github.com/restic/restic/pull/1209

 * Enhancement #1218: Add `--compact` to `snapshots` command

   The option `--compact` was added to the `snapshots` command to get a better overview of the
   snapshots in a repo. It limits each snapshot to a single line.

   https://github.com/restic/restic/issues/1218
   https://github.com/restic/restic/pull/1223

 * Enhancement #317: Add `--exclude-caches` and `--exclude-if-present`

   A new option `--exclude-caches` was added that allows excluding cache directories (that are
   tagged as such). This is a special case of a more generic option `--exclude-if-present` which
   excludes a directory if a file with a specific name (and contents) is present.

   https://github.com/restic/restic/issues/317
   https://github.com/restic/restic/pull/1170
   https://github.com/restic/restic/pull/1224

 * Enhancement #697: Automatically generate man pages for all restic commands

   https://github.com/restic/restic/issues/697
   https://github.com/restic/restic/pull/1147

 * Enhancement #1044: Improve `restore`, do not traverse/load excluded directories

   https://github.com/restic/restic/pull/1044

 * Enhancement #1061: Add Dockerfile and official Docker image

   https://github.com/restic/restic/pull/1061

 * Enhancement #1126: Use the standard Go git repository layout, use `dep` for vendoring

   The git repository layout was changed to resemble the layout typically used in Go projects,
   we're not using `gb` for building restic any more and vendoring the dependencies is now taken
   care of by `dep`.

   https://github.com/restic/restic/pull/1126

 * Enhancement #1134: Add support for storing backups on Google Cloud Storage

   https://github.com/restic/restic/issues/211
   https://github.com/restic/restic/pull/1134
   https://github.com/restic/restic/pull/1052

 * Enhancement #1144: Properly report errors when reading files with exclude patterns

   https://github.com/restic/restic/pull/1144

 * Enhancement #1149: Add support for storing backups on Microsoft Azure Blob Storage

   The library we're using to access the service requires Go 1.8, so restic now needs at least Go
   1.8.

   https://github.com/restic/restic/issues/609
   https://github.com/restic/restic/pull/1149
   https://github.com/restic/restic/pull/1059

 * Enhancement #1196: Add `--group-by` to `forget` command for flexible grouping

   https://github.com/restic/restic/pull/1196

 * Enhancement #1203: Print stats on all BSD systems when SIGINFO (ctrl+t) is received

   https://github.com/restic/restic/pull/1203
   https://github.com/restic/restic/pull/1082#issuecomment-326279920

 * Enhancement #1205: Allow specifying time/date for a backup with `--time`

   https://github.com/restic/restic/pull/1205


Changelog for restic 0.7.1 (2017-07-22)
=======================================

The following sections list the changes in restic 0.7.1 relevant to
restic users. The changes are ordered by importance.

Summary
-------

 * Fix #1115: Fix `prune`, only include existing files in indexes
 * Enh #1055: Create subdirs below `data/` for local/sftp backends
 * Enh #1067: Allow loading credentials for s3 from IAM
 * Enh #1073: Add `migrate` cmd to migrate from `s3legacy` to `default` layout
 * Enh #1081: Clarify semantic for `--tasg` for the `forget` command
 * Enh #1080: Ignore chmod() errors on filesystems which do not support it
 * Enh #1082: Print stats on SIGINFO on Darwin and FreeBSD (ctrl+t)

Details
-------

 * Bugfix #1115: Fix `prune`, only include existing files in indexes

   A bug was found (and corrected) in the index rebuilding after prune, which led to indexes which
   include blobs that were not present in the repo any more. There were already checks in place
   which detected this situation and aborted with an error message. A new run of either `prune` or
   `rebuild-index` corrected the index files. This is now fixed and a test has been added to detect
   this.

   https://github.com/restic/restic/pull/1115

 * Enhancement #1055: Create subdirs below `data/` for local/sftp backends

   The local and sftp backends now create the subdirs below `data/` on open/init. This way, restic
   makes sure that they always exist. This is connected to an issue for the sftp server.

   https://github.com/restic/restic/issues/1055
   https://github.com/restic/rest-server/pull/11#issuecomment-309879710
   https://github.com/restic/restic/pull/1077
   https://github.com/restic/restic/pull/1105

 * Enhancement #1067: Allow loading credentials for s3 from IAM

   When no S3 credentials are specified in the environment variables, restic now tries to load
   credentials from an IAM instance profile when the s3 backend is used.

   https://github.com/restic/restic/issues/1067
   https://github.com/restic/restic/pull/1086

 * Enhancement #1073: Add `migrate` cmd to migrate from `s3legacy` to `default` layout

   The `migrate` command for changing the `s3legacy` layout to the `default` layout for s3
   backends has been improved: It can now be restarted with `restic migrate --force s3_layout`
   and automatically retries operations on error.

   https://github.com/restic/restic/issues/1073
   https://github.com/restic/restic/pull/1075

 * Enhancement #1081: Clarify semantic for `--tasg` for the `forget` command

   https://github.com/restic/restic/issues/1081
   https://github.com/restic/restic/pull/1090

 * Enhancement #1080: Ignore chmod() errors on filesystems which do not support it

   https://github.com/restic/restic/pull/1080
   https://github.com/restic/restic/pull/1112

 * Enhancement #1082: Print stats on SIGINFO on Darwin and FreeBSD (ctrl+t)

   https://github.com/restic/restic/pull/1082


Changelog for restic 0.7.0 (2017-07-01)
=======================================

The following sections list the changes in restic 0.7.0 relevant to
restic users. The changes are ordered by importance.

Summary
-------

 * Fix #1013: Switch back to using the high-level minio-go API for s3
 * Fix #965: Switch to `default` repo layout for the s3 backend
 * Enh #1021: Detect invalid backend name and print error
 * Enh #1029: Remove invalid pack files when `prune` is run
 * Enh #512: Add Backblaze B2 backend
 * Enh #636: Add dirs `tags` and `hosts` to fuse mount
 * Enh #989: Improve performance of the `find` command
 * Enh #975: Add new backend for OpenStack Swift
 * Enh #998: Improve performance of the fuse mount

Details
-------

 * Bugfix #1013: Switch back to using the high-level minio-go API for s3

   For the s3 backend we're back to using the high-level API the s3 client library for uploading
   data, a few users reported dropped connections (which the library will automatically retry
   now).

   https://github.com/restic/restic/issues/1013
   https://github.com/restic/restic/issues/1023
   https://github.com/restic/restic/pull/1025

 * Bugfix #965: Switch to `default` repo layout for the s3 backend

   The default layout for the s3 backend is now `default` (instead of `s3legacy`). Also, there's a
   new `migrate` command to convert an existing repo, it can be run like this: `restic migrate
   s3_layout`

   https://github.com/restic/restic/issues/965
   https://github.com/restic/restic/pull/1004

 * Enhancement #1021: Detect invalid backend name and print error

   Restic now tries to detect when an invalid/unknown backend is used and returns an error
   message.

   https://github.com/restic/restic/issues/1021
   https://github.com/restic/restic/pull/1070

 * Enhancement #1029: Remove invalid pack files when `prune` is run

   The `prune` command has been improved and will now remove invalid pack files, for example files
   that have not been uploaded completely because a backup was interrupted.

   https://github.com/restic/restic/issues/1029
   https://github.com/restic/restic/pull/1036

 * Enhancement #512: Add Backblaze B2 backend

   https://github.com/restic/restic/issues/512
   https://github.com/restic/restic/pull/978

 * Enhancement #636: Add dirs `tags` and `hosts` to fuse mount

   The fuse mount now has two more directories: `tags` contains a subdir for each tag, which in turn
   contains only the snapshots that have this tag. The subdir `hosts` contains a subdir for each
   host that has a snapshot, and the subdir contains the snapshots for that host.

   https://github.com/restic/restic/issues/636
   https://github.com/restic/restic/pull/1050

 * Enhancement #989: Improve performance of the `find` command

   Improved performance for the `find` command: Restic recognizes paths it has already checked
   for the files in question, so the number of backend requests is reduced a lot.

   https://github.com/restic/restic/issues/989
   https://github.com/restic/restic/pull/993

 * Enhancement #975: Add new backend for OpenStack Swift

   https://github.com/restic/restic/pull/975
   https://github.com/restic/restic/pull/648

 * Enhancement #998: Improve performance of the fuse mount

   Listing directories which contain large files now is significantly faster.

   https://github.com/restic/restic/pull/998


Changelog for restic 0.6.1 (2017-06-01)
=======================================

The following sections list the changes in restic 0.6.1 relevant to
restic users. The changes are ordered by importance.

Summary
-------

 * Enh #985: Allow multiple parallel idle HTTP connections
 * Enh #981: Remove temporary path from binary in `build.go`
 * Enh #974: Remove regular status reports

Details
-------

 * Enhancement #985: Allow multiple parallel idle HTTP connections

   Backends based on HTTP now allow several idle connections in parallel. This is especially
   important for the REST backend, which (when used with a local server) may create a lot
   connections and exhaust available ports quickly.

   https://github.com/restic/restic/issues/985
   https://github.com/restic/restic/pull/986

 * Enhancement #981: Remove temporary path from binary in `build.go`

   The `build.go` now strips the temporary directory used for compilation from the binary. This
   is the first step in enabling reproducible builds.

   https://github.com/restic/restic/pull/981

 * Enhancement #974: Remove regular status reports

   Regular status report: We've removed the status report that was printed every 10 seconds when
   restic is run non-interactively. You can still force reporting the current status by sending a
   `USR1` signal to the process.

   https://github.com/restic/restic/pull/974


Changelog for restic 0.6.0 (2017-05-29)
=======================================

The following sections list the changes in restic 0.6.0 relevant to
restic users. The changes are ordered by importance.

Summary
-------

 * Enh #957: Make `forget` consistent
 * Enh #966: Unify repository layout for all backends
 * Enh #962: Improve memory and runtime for the s3 backend

Details
-------

 * Enhancement #957: Make `forget` consistent

   The `forget` command was corrected to be more consistent in which snapshots are to be
   forgotten. It is possible that the new code removes more snapshots than before, so please
   review what would be deleted by using the `--dry-run` option.

   https://github.com/restic/restic/issues/953
   https://github.com/restic/restic/pull/957

 * Enhancement #966: Unify repository layout for all backends

   Up to now the s3 backend used a special repository layout. We've decided to unify the repository
   layout and implemented the default layout also for the s3 backend. For creating a new
   repository on s3 with the default layout, use `restic -o s3.layout=default init`. For further
   commands the option is not necessary any more, restic will automatically detect the correct
   layout to use. A future version will switch to the default layout for new repositories.

   https://github.com/restic/restic/issues/965
   https://github.com/restic/restic/pull/966

 * Enhancement #962: Improve memory and runtime for the s3 backend

   We've updated the library used for accessing s3, switched to using a lower level API and added
   caching for some requests. This lead to a decrease in memory usage and a great speedup. In
   addition, we added benchmark functions for all backends, so we can track improvements over
   time. The Continuous Integration test service we're using (Travis) now runs the s3 backend
   tests not only against a Minio server, but also against the Amazon s3 live service, so we should
   be notified of any regressions much sooner.

   https://github.com/restic/restic/pull/962
   https://github.com/restic/restic/pull/960
   https://github.com/restic/restic/pull/946
   https://github.com/restic/restic/pull/938
   https://github.com/restic/restic/pull/883



..
  Normally, there are no heading levels assigned to certain characters as the structure is
  determined from the succession of headings. However, this convention is used in Python’s
  Style Guide for documenting which you may follow:

########################
Tuning Backup Parameters
########################

************************************
Large Repositories and Min Pack Size
************************************

In certain instances, such as very large repositories, it is desired to have larger pack
sizes to reduce the number of files in the repository.  Notable examples are OpenStack
Swift and some Google Drive Team accounts, where there are hard limits on the total
number of files.  This can be achieved by either using the ``--min-packsize`` flag
or defining the ``$RESTIC_MIN_PACKSIZE`` environment variable.  Restic currently defaults
to a 4MB minimum pack size.

The side effect of increasing the pack size is increased client memory usage.  A bit of
tuning may be required to strike a balance between memory usage and number of pack files.

Restic uses the majority of it's memory according to the pack size, multiplied by the number
of parallel writers. For example, if you have 8 parallel writers, With a minimum pack size of
64 (Megabytes), you'll get a *minimum* of 512MB of memory usage. More details are available in
the Save Blob Concurrency subsection.

*********************
File Read Concurrency
*********************

In some instances, such as backing up traditional spinning disks, reducing the file read
concurrency is desired.  This will help reduce the amount of time spent seeking around
the disk and can increase the overall performance of the backup operation. You can specify
the concurrency of file reads with the ``$RESTIC_FILE_READ_CONCURRENCY`` environment variable
or the ``--file-read-concurrency`` flag for the backup subcommand.

*********************
Save Blob Concurrency
*********************

Restic defaults to creating a concurrent writer for each available CPU.  This will require a
minimum of ``$RESTIC_MIN_PACKSIZE`` or 4MB to allocate per writer.  If you tune the minimum
pack size up to a larger size, you can run into a situation where restic is consuming too much
memory.  This can be adjusted by specifying how many writers are created, by passing the
``$RESTIC_SAVE_BLOB_CONCURRENCY`` environment variable or the ``--save-blob-concurrency`` flag
for the backup subcommand.

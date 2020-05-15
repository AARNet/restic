..
  Normally, there are no heading levels assigned to certain characters as the structure is
  determined from the succession of headings. However, this convention is used in Python’s
  Style Guide for documenting which you may follow:

  # with overline, for parts
  * for chapters
  = for sections
  - for subsections
  ^ for subsubsections
  " for paragraphs

######################
Special Considerations
######################

************************************
Large Repositories and Min Pack Size
************************************

In certain instances, such as very large repositories, it is desired to have larger pack
sizes to reduce the number of files in the repository.  Notable examples are OpenStack
Swift and some Google Drive Team accounts.  This can be achieved by either using the
``--min-packsize`` flag or defining the ``$MIN_PACKSIZE`` environment variable.

The side effect of increasing the pack size is increased client memory usage.  A bit of
tuning may be required to strike a balance between memory usage and number of pack files.

*********************
File Read Concurrency
*********************

In some instances, such as backing up traditional spinning disks, reducing the file read
concurrency is desired.  This will help reduce the amount of time spent seeking around
the disk and can increase the overall performance of the backup operation.
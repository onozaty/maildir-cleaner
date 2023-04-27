# maildir-cleaner

[![GitHub license](https://img.shields.io/github/license/onozaty/maildir-cleaner)](https://github.com/onozaty/maildir-cleaner/blob/main/LICENSE)
[![Test](https://github.com/onozaty/maildir-cleaner/actions/workflows/test.yaml/badge.svg)](https://github.com/onozaty/maildir-cleaner/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/onozaty/maildir-cleaner/branch/main/graph/badge.svg?token=1FPNT0WZAW)](https://codecov.io/gh/onozaty/maildir-cleaner)

`maildir-cleaner` is a tool to clean up maildir.

`maildir-cleaner` has the following subcommands

* [delete](#delete) Delete old mails.
* [archive](#archive) Archive old mails.

## delete

Delete old mails.

### Usage

```
maildir-cleaner delete -d MAIL_DIR_PATH -a AGE
```

```
Usage:
  maildir-cleaner delete [flags]

Flags:
  -d, --dir string   User maildir path.
  -a, --age int      The number of age days to be deleted.
                     If you specify 10, mail that has been in the mailbox for more than 10 days since its arrival will be deleted.
  -h, --help         help for delete
```

### Example

The following is an example of deleting mail that is more than 30 days old, specifying the maildir of `user1`.

```
$ maildir-cleaner delete -d /home/user1/Maildir -a 30
Starts searching for the target mails. maildir: /home/user1/Maildir age: 30
Completed search. The target mails are listed below.
+--------------+-----------------+------------------+
| Name         | Number of mails | Total size(byte) |
+--------------+-----------------+------------------+
|              |               7 |           11,412 |
| A            |               2 |            1,644 |
| INBOX.Drafts |               1 |              507 |
| INBOX.Sent   |               5 |            2,611 |
| INBOX.Trash  |               1 |              506 |
| test         |               1 |              912 |
+--------------+-----------------+------------------+
|        Total |              17 |           17,592 |
+--------------+-----------------+------------------+
Starts deleting mails.
Completed deletion.
```

## archive

Archive old mails.

### Usage

```
maildir-cleaner archive -d MAIL_DIR_PATH -a AGE [--archive-folder ARCHIVE_FOLDER_NAME] [--archive-pattern ARCHIVE_PATTERN]
```

```
Usage:
  maildir-cleaner archive [flags]

Flags:
  -d, --dir string               User maildir path.
  -a, --age int                  The number of age days to be archived.
                                 If you specify 10, mail that has been in the mailbox for more than 10 days since its arrival will be archived.
      --archive-folder string    Archive folder name. (default "Archived")
      --archive-pattern string   Archive pattern. can be specified: keep, year, month (default "keep")
  -h, --help                     help for archive
```

There are three types of `--archive-pattern`.

* `keep` : Archives under the archive folder with the original folder name.  
    * `A` -> `Archived.A`
    * `A.B` -> `Archived.A.B`
* `year` : Archives under the archive folder with each year of mail delivery.
    * `Archived.2022`
    * `Archived.2023`
* `month` : Archives under the archive folder with each month of mail delivery.
    * `Archived.2022.11`
    * `Archived.2022.12`
    * `Archived.2023.01`

### Example

The following is an example of archiving mail that is more than 30 days old by specifying the maildir of `user1`.

```
maildir-cleaner archive -d /home/user1/Maildir -a 30
Starts searching for the target mails. maildir: /home/user1/Maildir age: 30
Completed search. The target mails are listed below.
+--------------+-----------------+------------------+
| Name         | Number of mails | Total size(byte) |
+--------------+-----------------+------------------+
|              |               7 |           11,412 |
| A            |               2 |            1,644 |
| INBOX.B      |               1 |              822 |
| INBOX.Drafts |               1 |              507 |
| INBOX.Sent   |               5 |            2,611 |
| INBOX.Trash  |               1 |              506 |
| test         |               1 |              912 |
+--------------+-----------------+------------------+
|        Total |              18 |           18,414 |
+--------------+-----------------+------------------+
Starts archiving mails.
Completed archive. The archived mails are listed below.
+-----------------------+-----------------+------------------+
| Name                  | Number of mails | Total size(byte) |
+-----------------------+-----------------+------------------+
| Archived              |               7 |           11,412 |
| Archived.A            |               2 |            1,644 |
| Archived.INBOX.B      |               1 |              822 |
| Archived.INBOX.Drafts |               1 |              507 |
| Archived.INBOX.Sent   |               5 |            2,611 |
| Archived.INBOX.Trash  |               1 |              506 |
| Archived.test         |               1 |              912 |
+-----------------------+-----------------+------------------+
|                 Total |              18 |           18,414 |
+-----------------------+-----------------+------------------+
```

If `year` is specified as the `--archive-pattern`, the mails are archived by year.

```
$ maildir-cleaner archive -d /home/user1/Maildir -a 30 --archive-pattern year
Starts searching for the target mails. maildir: /home/user1/Maildir age: 30
Completed search. The target mails are listed below.
+--------------+-----------------+------------------+
| Name         | Number of mails | Total size(byte) |
+--------------+-----------------+------------------+
|              |               7 |           11,412 |
| A            |               2 |            1,644 |
| INBOX.B      |               1 |              822 |
| INBOX.Drafts |               1 |              507 |
| INBOX.Sent   |               5 |            2,611 |
| INBOX.Trash  |               1 |              506 |
| test         |               1 |              912 |
+--------------+-----------------+------------------+
|        Total |              18 |           18,414 |
+--------------+-----------------+------------------+
Starts archiving mails.
Completed archive. The archived mails are listed below.
+---------------+-----------------+------------------+
| Name          | Number of mails | Total size(byte) |
+---------------+-----------------+------------------+
| Archived.2023 |               1 |              102 |
| Archived.2022 |              17 |           18,312 |
+---------------+-----------------+------------------+
|         Total |              18 |           18,414 |
+---------------+-----------------+------------------+
```

If `month` is specified as the `--archive-pattern`, the mails are archived by month.

```
$ maildir-cleaner archive -d /home/user1/Maildir -a 30 --archive-pattern month
Starts searching for the target mails. maildir: /home/user1/Maildir age: 30
Completed search. The target mails are listed below.
+--------------+-----------------+------------------+
| Name         | Number of mails | Total size(byte) |
+--------------+-----------------+------------------+
|              |               7 |           11,412 |
| A            |               2 |            1,644 |
| INBOX.B      |               1 |              822 |
| INBOX.Drafts |               1 |              507 |
| INBOX.Sent   |               5 |            2,611 |
| INBOX.Trash  |               1 |              506 |
| test         |               1 |              912 |
+--------------+-----------------+------------------+
|        Total |              18 |           18,414 |
+--------------+-----------------+------------------+
Starts archiving mails.
Completed archive. The archived mails are listed below.
+------------------+-----------------+------------------+
| Name             | Number of mails | Total size(byte) |
+------------------+-----------------+------------------+
| Archived.2022.10 |               1 |              102 |
| Archived.2022.11 |               5 |            3,120 |
| Archived.2022.12 |               5 |            3,777 |
| Archived.2023.01 |               7 |           11,302 |
| Archived.2023.02 |               1 |              110 |
+------------------+-----------------+------------------+
|            Total |              18 |           18,414 |
+------------------+-----------------+------------------+
```

## Install

`maildir-cleaner` is implemented in golang and runs on all major platforms such as Windows, Mac OS, and Linux.  
You can download the binaries for each OS from the links below.

You can download the binary from the following.

* https://github.com/onozaty/maildir-cleaner/releases/latest

## License

MIT

## Author

[onozaty](https://github.com/onozaty)

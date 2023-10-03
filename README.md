# mkdt

make directory tree, a tool for humans.

## usage

### command

`mkdir [flags]`

| flag | description                                                     |   |   |   |
|------|-----------------------------------------------------------------|---|---|---|
| -f   | input file, use existing input file instead of launching editor |   |   |   |
| -r   | set target root directory, otherwise default to cwd             |   |   |   |
| -d   | dry run, print tree instead of making it                        |   |   |   |

### input file

By default, mkdir will launch the user's default editor to create
a temp file describing the directory tree. Alternatively, a pre-made file can be
specified with `-f`.

Each line of the input file is interpreted as either a directory or a filename.
Filenames are expected to have an `.` followed by an extension.  Otherwise the line is assumed a directory.
Any whitespace at the head of each line is interpreted as a directory level relative to other lines.

Example input file
```
dir1
    file1.txt
    dir2
        file2.txt
        dir3
            file3.txt
    dir4
        file4.txt
file1.txt
```

Resulting tree
```
``````

## installation
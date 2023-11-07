# mkdt

make directory tree, a tool for humans.

## build

go build -o mkdt cmd/mkdt/main.go
mv mkdt /usr/local/bin

## usage

### command

`mkdir [flags]`

| flag | description                                                     |   |   |   |
|------|-----------------------------------------------------------------|---|---|---|
| -f   | input file, use existing input file instead of launching editor |   |   |   |
| -r   | set target root directory, otherwise default to cwd             |   |   |   |
| -d   | dry run, print tree instead of making it                        |   |   |   |
| -v   | verbose, print out exactly whats going on                       |   |   |   |

### input file

By default, mkdir will launch the user's default editor to create
a temp file describing the directory tree. Alternatively, a pre-made file can be
specified with `-f`.

Each line of the input file is interpreted as either a directory or a filename.
Filenames are expected to have an `.` followed by an extension.  Otherwise the line is assumed a directory.
Any whitespace at the head of each line is interpreted as a directory level relative to other lines.

Example input file
```
test
  test.json
one
  test2.json
  two
    test.json
```

Resulting tree
```
```


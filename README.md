# mkdt

make directory tree, a tool for humans.
Jump start a project structure from a text file.

## build

go build -o mkdt main.go
mv mkdt /usr/local/bin

## usage

### command

`mkdir [-fr]`

| flag | description                                                     |
|------|-----------------------------------------------------------------|
| -f   | input file, use existing input file instead of launching editor |
| -r   | set target root directory, otherwise default to cwd             |

### input file

By default, mkdir will launch the user's default editor (using the `EDITOR` environment variable) to create
a temp file describing the directory tree. Alternatively, a pre-made file can be
specified with `-f`.

Each line of the input file is interpreted as either a directory or a filename.
Filenames are expected to have an `.` followed by an extension.  Otherwise the line is assumed a directory.
Any whitespace (`\t` or `space`) at the head of each line is interpreted as a directory level relative to other lines.

Example input file
```
test
 a
  1.txt
 b
  c
   2.txt
 3.txt
```

Resulting tree (`tree ./test`)
```
./test
├── 3.txt
├── a
│   └── 1.txt
└── b
    └── c
        └── 2.txt

3 directories, 3 files
```

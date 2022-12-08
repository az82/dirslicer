# dirslicer

*dirslicer* is a small command line utility that is able to slices and re-slice
directories into slices with a given size. It can process directories with a
very large number of files, flat or divided into sub-directories.
The only condition is that the file names fit into memory.

Given a directory `./source` with the files `a,b,c,d,e`,

```bash
dirslice --size 3 ./source ./target
```

will move the files to this structure in `./target`:

```text
 a-c/ a,b,c
 d-e/ d,e
```


## Usage

```text
Usage: dirslice SOURCE TARGET
  SOURCE must be a directory
  TARGET must either be an empty directory or a new location

  -help
        Show help
  -size int
        Number of files per directory (default 1024)
```

## Build

### Run The Tests

```bash
make test
```

### Build An Executable For The Current Platform

```bash
make build
```

After that, the directory `./bin` will contain the executable for the current
platform.

### Build Distributions for all supported platforms

Provide the desired version with the environment variable VERSION.

```bash
VERSION=1.2.3 make dist
```

If version is not specified, the git tag, git branch or commit ID will be used.

### Clean Build Results

```bash
make clean
```

## License

Copyright (c) 2022 Andreas Zitzelsberger

This software can be under the terms and conditions of the
[MIT License (MIT)](LICENSE).

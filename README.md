# cstar

Show which projects you have contributed to.

![screenshot](./ss.png)

## Description

This is a tool for [taking-star-table with color(ja)](https://note.mu/hyuki/n/n9a6e7c1e0d7b), which devised by [hyuki(tw)](https://twitter.com/hyuki).

## Usage

Simply you can run `cstar` at your project root.

```
$ cstar
```

Cstar will walk through the lower directories to find `.git`.  
With each `.git` directories, cstar counts commits per day, then show how hard you worked with your git repositories.

You can assign some paths which cstar works at.

```
$ cstar /path/to/dir/ /path/to/dir2/
```

### options

- `-d, --day`
  - Days to list up. The default is 7days.
  - ex: `cstar -d 5` or `cstar /path/to/dir/ -d 5`

## Install

### Binary

If you have not installed `go` command, [download cstar binary(from here)](https://github.com/hkdnet/cstar/releases) and place it in your `$PATH`.

### Golang

If you have installed `go` command, use `go get`:

```bash
$ go get -d github.com/hkdnet/cstar
```

## Contribution

1. Fork ([https://github.com/hkdnet/cstar/fork](https://github.com/hkdnet/cstar/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[hkdnet](https://github.com/hkdnet)

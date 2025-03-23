<!-- PROJECT LOGO -->
<br />
<div align="center">
<h3 align="center">ATP (Automated Task Planner)</h3>
  <a href="https://github.com/arjungandhi/ATP">
    <img src="images/logo.png" alt="Logo" height="200">
  </a>

  <p align="center">
    the mitochondria is the power house of the cell 
  </p>
</div>

[![Go Report Card](https://goreportcard.com/badge/github.com/arjungandhi/atp?style=flat-square)](https://goreportcard.com/report/github.com/arjungandhi/atp)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/arjungandhi/atp)](https://pkg.go.dev/github.com/arjungandhi/atp)
[![Release](https://img.shields.io/github/release/arjungandhi/atp.svg?style=flat-square)](https://github.com/arjungandhi/atp/releases/latest)


# About

ATP is a task/ plan manager

## Install

This command can be installed as a standalone program or composed into a
Bonzai command tree.

Standalone

```
go install github.com/arjungandhi/atp/cmd/atp@latest
```

Composed

```go
package z

import (
	Z "github.com/rwxrob/bonzai/z"
	atp "github.com/arjungandhi/atp"
)

var Cmd = &Z.Cmd{
	Name:     `z`,
	Commands: []*Z.Cmd{help.Cmd, atp.Cmd},
}
```

## Tab Completion

To activate bash completion just use the `complete -C` option from your
`.bashrc` or command line. There is no messy sourcing required. All the
completion is done by the program itself.

```
complete -C atp atp
```

If you don't have bash or tab completion check use the shortcut
commands instead.

## Embedded Documentation

All documentation (like manual pages) has been embedded into the source
code of the application. See the source or run the program with help to
access it.


## Setup

Define a storage location for your ATP data. This can be a directory or
a file. The default is `~/.atp` but can be overridden with the '$ATP_DIR' env var


## Design Doc

Stored in [design_doc.md](design_doc.md)



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


## Design Doc

ATP manages tasks, (as opposed to higher level concepts like projects)

Tasks are the lowest unit of work in my life 

ATP as a tool has the following responsibilities 
- Keep a list of all non external* things I need to do, (regardless of work/ personal/ adventure whatever
- Track my time spent on those items
- Do a first pass sort and organize on those tasks (this will almost certainly never be perfect)
- Add those items as events into a calendar while respecting actual calendar events

*Non External means not a calendar event/ involving other people

As Tasks are the lowest level of work in my life, ATP is my lowest level planning tool, as such that there are a few design tenants that arise from this

1. Simple: the core of ATP should be dead simple to hit the jobs described above
2. Tasks as a store of data: as tasks are the lowest level of work in my system there are certainly lots of information that could be useful to store with them that I'm not thinking of now. As such its probably easiest to store that data in the tasks its self



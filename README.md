# Ham related go code

## cmd/termlog

![screenshot](https://raw.githubusercontent.com/tzneal/ham-go/master/_screenshots/screenshot.png)

The main driver for writing the rest of this code.  A cross platform console
based ham contact logger.  I'm developing it for my own use, so it only has
features I need/want.

- Saves logs as ADIF files
- Can auto-commit log files to a git repository
- DX Cluster Monitor
- Controls radio through hamlib (github.com/dh1tw/goHamlib)

## Installation

```
go install github.com/tzneal/ham-go/cmd/termlog
```

## Configuration

1) Run termlog once, then hit Ctrl+Q to quit.  This will create an initial
   config file at ~/.termlog.toml that you can then modify.
2) Fill out the operator section at a minimum

# Commands

c.AddCommand(input.KeyCtrlL, ms.focusQSOList)
	c.AddCommand(input.KeyCtrlN, ms.newQSO)
	c.AddCommand(input.KeyCtrlD, ms.qso.ResetDateTime)
	c.AddCommand(input.KeyCtrlS, ms.saveQSO)
	c.AddCommand(input.KeyCtrlB, ms.saveBookmark)
	c.AddCommand(input.KeyCtrlG, ms.commitLog)
    
| Shortcut | Command |
|----------|---------|
| Ctrl+Q   | Quit termlog |
| Ctrl+H   | Display Help |
| Ctrl+N   | Start a new QSO (clearing the current one if not saved) |
| Ctrl+D   | Set the QSO time on to the current time |
| Ctrl+S   | Save the QSO to the log and start a new one |
| Ctrl+G   | Commit the current logfile to git |



## adif

ADIF parsing and writing

## callsigns

Callsign lookup interface with a couple of supported backends.

## dxcc

Callsign lookup via prefixes/exceptions through the data at
www.country-files.com (works offline).

## dxcluster

The beginnings of a DXCluster client.
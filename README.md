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
- Logs for both WSJT-X and fldigi


## Installation

```
go install github.com/tzneal/ham-go/cmd/termlog
```

## Configuration

1) Run termlog once, then hit Ctrl+Q to quit.  This will create an initial
   config file at ~/.termlog.toml that you can then modify.
2) Fill out the operator section at a minimum

# Commands
    
| Shortcut  | Command |
|-----------|---------|
| Ctrl+Q    | Quit termlog |
| Ctrl+H    | Display Help |
| Ctrl+N    | Start a new QSO (clearing the current one if not saved) |
| Ctrl+D    | Set the QSO time on to the current time |
| Ctrl+S    | Save the QSO to the log and start a new one |
| Ctrl+G    | Commit the current logfile to git |
| Ctrl+B    | Save a bookmark |
| Alt+B     | Open the bookmark list |
| Ctrl+L    | Focus the QSO List |
| Alt+Left  | Tune down 500khz |
| Alt+Right | Tune up 500khz |

## adif

ADIF parsing and writing

## callsigns

Callsign lookup interface with a couple of supported backends.

## db

ADIF indexer used to quickly identify when you last saw a contact and how many
times you've logged him.

## dxcc

Callsign lookup via prefixes/exceptions through the data at
www.country-files.com (works offline).

## dxcluster

The beginnings of a DXCluster client.

## fldigi

Enough code to parse the realtime fldigi emitted logs and save them to termlog.

## wsjtx

Enough code to parse the realtime WSJT-X emitted logs and save them to termlog.
# Ham related go code

## cmd/termlog

![screenshot](https://raw.githubusercontent.com/tzneal/ham-go/master/_screenshots/screenshot.png)

The main driver for writing the rest of this code.  A cross platform console
based ham contact logger.  I'm developing it for my own use, so it only has
features I need/want.

- Saves logs as ADIF files with custom fields
- Supports auto-commiting log files to a git repository
- DX cluster & POTA spot monitoring
- Radio control through hamlib (github.com/dh1tw/goHamlib)
- Logs for both WSJT-X and fldigi
- LoTW integration (syncs QSL information from LoTW to stored ADIF files)

## Installation

```
go install github.com/tzneal/ham-go/cmd/termlog
```

## Configuration

1) Run termlog once, then hit Ctrl+Q to quit.  This will create an initial
   config file at ~/.termlog.toml that you can then modify.
2) Fill out the operator section at a minimum

### Custom Fields
In the configuration file, custom ADIF fields can be defined.  A SOTA field is defined in the default 
configuration and can be removed if not needed.

```
  [[Operator.CustomFields]]
    Name = "sota_ref"
    Label = "SOTA"
    Width = 8
    Default = ""
```

# Command line options
```
Usage of ./termlog:
  -color-test
    	display a color test
  -config string
    	path to the configuration file (default "~/.termlog.toml")
  -hamlib-list
    	list the supported libhamlib devices
  -index
    	index the ADIF files passed in on the command line
  -key-test
    	list keyboard events
  -log string
    	specify a log file to load and write to
  -no-net
    	disable all features that require network access (useful for POTA/SOTA)
  -no-rig
    	disable rig control, even if enabled in the config file
  -search string
    	search the indexed ADIF files and print the results
  -sync-lotw-qsl
    	fetches QSL information from LoTW to update log QSL information in the default log directory
  -upgrade-config
    	upgrade the configuration file to the latest format
```

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
| Ctrl+E    | Display custom user commands |
| Ctrl+R    | Force screen redraw |
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

Enough code to parse the realtime WSJT-X emitted logs and save them to termlog. I use this when running
FT8 to capture logs in real time from WSJT-X.
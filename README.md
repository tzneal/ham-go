# Ham related go code

## cmd/termlog

The driver for the rest of this code.  A cross platform console based ham
contact logger.  I'm developing it for my own use, so it only has features I
need/want.

- Saves logs as ADIF files
- Can auto-commit log files to a git repository
- DX Cluster Monitor
- Controls radio through hamlib (github.com/dh1tw/goHamlib)

## adif

ADIF parsing and writing

## callsigns

Callsign lookup interface with a couple of supported backends.

## dxcc

Callsign lookup via prefixes/exceptions through the data at
www.country-files.com (works offline).

## dxcluster

The beginnings of a DXCluster client.
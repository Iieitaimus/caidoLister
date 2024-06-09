# CaidoLister

CaidoLister is a tool designed to generate lists of hosts, paths, and parameters from Caido CSV export files.

## Usage Example

Suppose you have a CSV export file from [Caido](https://docs.caido.io/reference/features/logging/exports.html). You can use CaidoLister to extract hosts, paths, and parameters from this file.

```
caidoLister -f caido_export.csv
```

This command will output `hosts.txt`, `paths.txt`, and `params.txt` based on the data in `caido_export.csv`.

## Flags

- `-f`: Specifies the Caido CSV export file from which to extract data.

## Install

```
go install github.com/Iieitaimus/caidoLister@latest
``` 

With this tool, you can easily extract useful information from your Caido CSV export files.
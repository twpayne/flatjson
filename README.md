# `flatjson`

[![Build Status](https://travis-ci.org/twpayne/flatjson.svg?branch=master)](https://travis-ci.org/twpayne/flatjson)
[![GoDoc](https://godoc.org/github.com/twpayne/flatjson?status.svg)](https://godoc.org/github.com/twpayne/flatjson)
[![Report Card](https://goreportcard.com/badge/github.com/twpayne/flatjson)](https://goreportcard.com/report/github.com/twpayne/flatjson)

## Overview

`flatjson` converts JSON files to a "flat" representation with one value per
line. For example, given the input:

```json
{
  "menu": {
    "id": "file",
    "value": "File",
    "popup": {
      "menuitem": [
        {
          "value": "New",
          "onclick": "CreateNewDoc()"
        },
        {
          "value": "Open",
          "onclick": "OpenDoc()"
        },
        {
          "value": "Close",
          "onclick": "CloseDoc()"
        }
      ]
    }
  }
}
```

`flatjson` outputs:

```js
root = {};
root.menu = {};
root.menu.id = "file";
root.menu.popup = {};
root.menu.popup.menuitem = [];
root.menu.popup.menuitem[0] = {};
root.menu.popup.menuitem[0].onclick = "CreateNewDoc()";
root.menu.popup.menuitem[0].value = "New";
root.menu.popup.menuitem[1] = {};
root.menu.popup.menuitem[1].onclick = "OpenDoc()";
root.menu.popup.menuitem[1].value = "Open";
root.menu.popup.menuitem[2] = {};
root.menu.popup.menuitem[2].onclick = "CloseDoc()";
root.menu.popup.menuitem[2].value = "Close";
root.menu.value = "File";
```

This format, although verbose, makes it much easier to see the nesting of
values. It also happens to be valid JavaScript that can be used to recreate the
original JSON object.

This "flat" format is very handy for visualizing diffs. For example, comparing
the above JSON object with a second JSON object:

```json
{
  "menu": {
    "id": "file",
    "disabled": true,
    "value": "File menu",
    "popup": {
      "menuitem": [
        {
          "value": "New",
          "onclick": "CreateNewDoc()"
        },
        {
          "value": "Open",
          "onclick": "OpenDoc()"
        }
      ]
    }
  }
}
```

yields the diff:

```diff
--- testdata/a.json
+++ testdata/b.json
@@ -1,5 +1,6 @@
 root = {};
 root.menu = {};
+root.menu.disabled = true;
 root.menu.id = "file";
 root.menu.popup = {};
 root.menu.popup.menuitem = [];
@@ -9,8 +10,5 @@
 root.menu.popup.menuitem[1] = {};
 root.menu.popup.menuitem[1].onclick = "OpenDoc()";
 root.menu.popup.menuitem[1].value = "Open";
-root.menu.popup.menuitem[2] = {};
-root.menu.popup.menuitem[2].onclick = "CloseDoc()";
-root.menu.popup.menuitem[2].value = "Close";
-root.menu.value = "File";
+root.menu.value = "File menu";
```

## Installation

    go install github.com/twpayne/flatjson/cmd/flatjson

## Generating flat JSON

To convert a file to flat JSON, specify it on the command line, for example:

    flatjson vendor/vendor.json

If no filenames are specified, `flatjson` will read JSON from the standard
input.

## Generated a unified diff

To generate a unified diff between two JSON files, specify the `-diff` option
and the filenames on the command line, for example:

    flatjson -diff ./testdata/a.json ./testdata/b.json

An additional `-context` option specifies how many lines of context to show.
The default is three.

## License

MIT

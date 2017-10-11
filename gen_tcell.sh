#!/bin/sh

cp main.go tcell.go
perl -i -pe 's{github.com/nsf/termbox-go}{github.com/gdamore/tcell/termbox}; s{!tcell}{tcell\n// DO NOT EDIT, generated automatically};' tcell.go

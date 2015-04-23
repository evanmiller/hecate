# hecate
the hex editor from hell

Usage:

    go get -u github.com/nsf/termbox-go
    go build
    ./hecate /path/to/binary/file

Hecate is not (yet) a hex editor, only a viewer. It is a terminal program with
Vim-like controls; place the cursor over some bytes and choose a mode (**t**
for text, **p** for a bit pattern, **i** for an integer, **f** for a floating
point) to see what those bytes represent.

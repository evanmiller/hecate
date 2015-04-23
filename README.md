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

Full list of commands:

<table>
<tr><td>h</td><td>left</td></tr>
<tr><td>j</td><td>down</td></tr>
<tr><td>k</td><td>up</td></tr>
<tr><td>l</td><td>right</td></tr>

<tr><td>b</td><td>left 4 bytes</td></tr>
<tr><td>w</td><td>right 4 bytes</td></tr>

<tr><td>g</td><td>first byte</td></tr>
<tr><td>G</td><td>last byte</td></tr>

<tr><td>ctrl-e</td><td>scroll down</td></tr>
<tr><td>ctrl-y</td><td>scroll up</td></tr>

<tr><td>ctrl-f</td><td>page down</td></tr>
<tr><td>ctrl-b</td><td>page up</td></tr>
<tr><td>t</td><td>text mode</td></tr>
<tr><td>p</td><td>bit pattern mode</td></tr>
<tr><td>i</td><td>integer mode</td></tr>
<tr><td>f</td><td>floating-point mode</td></tr>

<tr><td>e</td><td>toggle endianness</td></tr>
<tr><td>u</td><td>toggle signedness</td></tr>

<tr><td>H</td><td>shrink cursor</td></tr>
<tr><td>L</td><td>grow cursor</td></tr>

<tr><td>:</td><td>jump to offset</td></tr>
<tr><td>x</td><td>toggle hex offset</td></tr>

<tr><td>?</td><td>this screen</td></tr>
<tr><td>q</td><td>quit program</td></tr>
</table>

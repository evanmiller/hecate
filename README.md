# hecate
The Hex Editor From Hell!

Download: **[Linux, Mac OS X, and Windows](releases)**

Compile from source:

    go get -u github.com/evanmiller/hecate
    $GOPATH/bin/hecate file1 [file2 [...]]

Hecate is a hex editor designed for reverse engineering file formats without
leaving the comfort of your terminal window. Hecate (pronounced HECK-it)
features tabbed browsing, in-place editing, large-file support, full-file
searching, and Vim-like controls.  Place the cursor over some bytes and choose
a mode (**t** for text, **p** for a bit pattern, **i** for an integer, **f**
for a floating point) to see what those bytes represent. Toggle endianness with
**e** and signedness with **u**. Press **enter** to edit.


### Editing

Pressing **enter** brings up an edit field for the data under the cursor. Make
changes and press **enter** again to write changes to disk. Pressing **esc**
cancels any changes on the current position, otherwise exits edit mode.
Navigating passed the edges of the field moves the cursor. The expected format
depends on the cursor mode when entering edit mode.


Screenshot:
![Hecate screenshot](http://www.evanmiller.org/images/hecate-screenshot2.png)

Full list of commands:


<table>
<tr><td>h</td><td>left</td><td>t</td><td>text mode</td><td>S</td><td>show tabs</td></tr>
<tr><td>j</td><td>down</td><td>p</td><td>bit pattern mode</td><td>W</td><td>hide tabs</td></tr>
<tr><td>k</td><td>up</td><td>i</td><td>integer mode</td><td>A</td><td>previous tab</td></tr>
<tr><td>l</td><td>right</td><td>f</td><td>float mode</td><td>D</td><td>next tab</td></tr>
<tr><td>b</td><td>left 4 bytes</td><td>H</td><td>shrink cursor</td><td>ctrl-t</td><td>new tab</td></tr>
<tr><td>w</td><td>right 4 bytes</td><td>L</td><td>grow cursor</td><td>ctrl-w</td><td>close tab</td></tr>
<tr><td>^</td><td>line start</td><td>e</td><td>toggle endianness</td><td>ctrl-e</td><td>scroll down</td></tr>
<tr><td>$</td><td>line end</td><td>u</td><td>toggle signedness</td><td>ctrl-y</td><td>scroll up</td></tr>
<tr><td>g</td><td>file start</td><td>a</td><td>date decoding</td><td>ctrl-f</td><td>page down</td></tr>
<tr><td>G</td><td>file end</td><td>@</td><td>set date epoch</td><td>ctrl-b</td><td>page up</td></tr>
<tr><td>:</td><td>jump to byte</td><td>/</td><td>search file</td><td>enter</td><td>edit mode</td></tr>
<tr><td>x</td><td>toggle hex</td><td>n</td><td>next match</td><td>?</td><td>help screen</td></tr>
</table>

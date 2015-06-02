# hecate
The Hex Editor From Hell!

Usage:

    go get -u github.com/evanmiller/hecate
    $GOPATH/bin/hecate file1 [file2 [...]]

Hecate is not actually a hex editor, only a viewer. It is a terminal program
(written in Go) with tabbed browsing, large-file support, full-file searching,
and Vim-like controls.  Place the cursor over some bytes and choose a mode
(**t** for text, **p** for a bit pattern, **i** for an integer, **f** for a
floating point) to see what those bytes represent. Toggle endianness with **e**
and signedness with **u**.

Screenshot:
![Hecate screenshot](http://www.evanmiller.org/images/hecate-screenshot2.png)

Full list of commands:

<table>
<tr><td>h</td><td>left</td><td>t</td><td>text mode</td><td>ctrl-j</td><td>show tabs</td></tr>                                                                          
<tr><td>j</td><td>down</td><td>p</td><td>bit pattern mode</td><td>ctrl-k</td><td>hide tabs</td></tr>                                                                   
<tr><td>k</td><td>up</td><td>i</td><td>integer mode</td><td>ctrl-t</td><td>new tab</td></tr>                                                                           
<tr><td>l</td><td>right</td><td>f</td><td>float mode</td><td>ctrl-w</td><td>close tab</td></tr>                                                                        
<tr><td>b</td><td>left 4 bytes</td><td>H</td><td>shrink cursor</td><td>ctrl-h</td><td>previous tab</td></tr>                                                           
<tr><td>w</td><td>right 4 bytes</td><td>L</td><td>grow cursor</td><td>ctrl-l</td><td>next tab</td></tr>                                                                
<tr><td>^</td><td>line start</td><td>e</td><td>toggle endianness</td><td>ctrl-e</td><td>scroll down</td></tr>                                                          
<tr><td>$</td><td>line end</td><td>u</td><td>toggle signedness</td><td>ctrl-y</td><td>scroll up</td></tr>                                                              
<tr><td>g</td><td>file start</td><td>D</td><td>date decoding</td><td>ctrl-f</td><td>page down</td></tr>                                                                
<tr><td>G</td><td>file end</td><td>@</td><td>set date epoch</td><td>ctrl-b</td><td>page up</td></tr>                                                                   
<tr><td>:</td><td>jump to byte</td><td>/</td><td>search file</td><td>?</td><td>help screen</td></tr>                                                                   
<tr><td>x</td><td>toggle hex</td><td>n</td><td>next match</td><td>q</td><td>quit program</td></tr>  
</table>

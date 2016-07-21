:Title: Files for not so sure programmers like me
:Date: 2016-07-21 11:00
:Category: Programming

Recently at work, I needed to replace certain patterns by another in a JPEG file. Specifically,
I wanted to replace // by /, /" by ", /' by ' and //0 by /0. If it was a text file, I wouldn't 
get confused, but hey it's a JPEG file, it's a binary file. How do I read it and replace bytes
in it? I know how to replace strings, but bytes? I should open it in binary mode right? I was
not very happy to realize ow much I lacked a true understanding of what it meant when a file was 
referred to as binary or text.

In a file everything is a byte, It's how your program interprets it that makes it a JPEG file, a text file or a a PNG file

see libmagic and file command
Linux strings command


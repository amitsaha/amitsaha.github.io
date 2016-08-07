:Title: Files for not so sure programmers like me
:Date: 2016-07-21 11:00
:Category: Programming

Recently at work, I needed to replace certain patterns by another in a JPEG file. Specifically,
I wanted to replace // by /, /" by ", /' by ' and //0 by /0. If it was a text file, I wouldn't 
get confused, but hey it's a JPEG file, it's a "binary" file. How do I read it and replace bytes
in it? I know how to replace strings, but bytes? I should open it in binary mode right? I was
not very happy to realize how much I lacked a true understanding of what it meant when a file was 
referred to as binary or text. In this post, I will share my findings post this realization.

After reading a number of blog posts, stack overflow answers, here's the one-liner: "Every file is 
a binary file". A file is just a series of bytes (a byte as in binary byte) stored on disk via your
filesystem. What "happens" to your file when you open it depends on the program you use to open it.
For example:

- If we trying to display the contents of a PDF file using the command line program "cat" or a text editor
, we will mostly see random strings which we probably can't make sense of

- Linux strings command

It's how your program interprets it that makes it a JPEG file, a text file or a PDF file.



see libmagic and file command

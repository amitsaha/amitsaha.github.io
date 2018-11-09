Title: Hard Links and Soft/Symbolic Links on Linux
Date: 2018-11-09
Category: software
Status: Draft

[Much](https://medium.com/@wendymayorgasegura/what-is-the-difference-between-a-hard-link-and-a-symbolic-link-8c0493041b62) 
[has](https://medium.com/@meghamohan/hard-link-and-symbolic-link-3cad74e5b5dc) [been](https://medium.com/meatandmachines/explaining-the-difference-between-hard-links-symbolic-links-using-bruce-lee-32828832e8d3) written (and [asked](https://stackoverflow.com/questions/185899/what-is-the-difference-between-a-symbolic-link-and-a-hard-link)) 
on the topic of hard links and soft links (a.k.a symbolic links) on Linux. I have read a few of those more than once.
However, I end up getting confused between the two, specifically the differences between the two. So, here's 
my post on the topic with the hope that I will stop getting confused ever again.

## Our setup

Let's create a file and write a line into it:

```
$ echo "Hello, I am file1" > file1
```

Next, we create a `hard` link using the `ln` command:

```
$ ln file1 file1-hlink
```

Now, let's create a soft link using `ln -s`:

```
$ ln -s file1 file1-slink
```

At this stage, if we use the `cat` command to display the contents of each of the above, we
will see the same line of text:

```
$ cat file1
Hello, I am file1

$ cat file1-hlink
Hello, I am file1

$ cat file1-slink
Hello, I am file1
```


## Investigation: Inodes

One of the key differences between soft links and hard links is with respect to how they are represented
in the filesystem. If we run `ls` with the `i` switch, it will show the [inode](https://en.wikipedia.org/wiki/Inode)
number of each of the above files:


```
$ ls -il

15481123719144131 -rw-rw-rw- 2 asaha asaha 18 Nov  9 13:52 file1
15481123719144131 -rw-rw-rw- 2 asaha asaha 18 Nov  9 13:52 file1-hlink
29836347531381846 lrwxrwxrwx 1 asaha asaha  5 Nov  9 13:54 file1-slink -> file1
```

We can see that:

- The hardlink `file1-hlink` has the same Inode number as the original file itself (`file1`)
- The softlink `file1-slink` has a different Inode number

This tells us two things straightaway:

When we create a soft link, it is equivalent to creating a new file with its own filename. In the filesystem, 
it is a separate file, with the special property that its contents is the `path` to the real file `file1`.

Graphically:

```
            Soft link ->   FILE CONTENTS -> Path of original file -> FILE CONTENTS -> "Hello, I am file1"
```

A hard link on the other hand is a reference to the original file. It exists on the filesystem, but only as another
reference or a link. Let's explore a bit into what it means. If we execute `ls` with the `-l` (small `L`) switch, the
second column gives the number of `link` counts of a file:

```
$ ls -l file1
-rw-rw-rw- 2 asaha asaha 18 Nov  9 13:52 file1
```

We have created a hard link above, so, the link count now is 2. If we create another hard link, the link count will be 3:

```
$ ln file1 file1-hlink-2
$ ls -l file1
-rw-rw-rw- 3 asaha asaha 18 Nov  9 13:52 file1
```

Graphically:

```
file1-hlink  -----> FILE CONTENTS ("Hello, I am file1") <------ file1
                         /|\ 
                          |
                     file1-hlink-2
```

## Investigation: Size of hard links and soft links

Let's go back to one of the previous output of `ls -l`:

```
$ ls -il

15481123719144131 -rw-rw-rw- 2 asaha asaha 18 Nov  9 13:52 file1
15481123719144131 -rw-rw-rw- 2 asaha asaha 18 Nov  9 13:52 file1-hlink
29836347531381846 lrwxrwxrwx 1 asaha asaha  5 Nov  9 13:54 file1-slink -> file1
```

The sixt column above of the output shows the number of bytes in each of the files. We see `18` as the size
of the original file, `file1` and the hardlink, `file1-hlink`. 18 is the number of characters in `"Hello, I am a file1"`
and a new line character. What are the five bytes in `file1-slink`? The `readlink` command will help us:

```
$ readlink file1-slink
file1
```
It is the "relative" path to the original file.


## Investigation: Deleting the original file

What happens to each kind of link when we delete the original file? From the graphics above, we expect that the symbolic
link will basically be a "dangling" link and hence, we will lose access to the file contents. In the case of hard link, the contents will still continue to be accessible, since all we are doing is deleting one of the links. Even though it is the original file,
it doesn't matter. Other links continue to exist and point to the data.

Let's validate our theory:

```
$ rm file1

$ ls -lrt
total 0
-rw-rw-rw- 2 asaha asaha 18 Nov  9 13:52 file1-hlink-2
-rw-rw-rw- 2 asaha asaha 18 Nov  9 13:52 file1-hlink
lrwxrwxrwx 1 asaha asaha  5 Nov  9 13:54 file1-slink -> file1
```

We delete the original file above. Now the link count of `file1-hlink` and `file-hlink-2` has decreased by 1 and
is now 2.

If we try to display the contents of a hard link:

```
$ cat file1-hlink
Hello, I am file1
```

For the soft link though:

```
$ cat file1-slink
cat: file1-slink: No such file or directory
```

What the above error really says is I am trying to look for a file, `file1`, but it doesn't exist. This also means that
we can essentially do:

```
$ echo "Hello, I am a different file1" > file1
$ cat file1-slink
Hello, I am a different file1
```

I wonder what kind of security risk this may post - may be we need symbolic links with checksums?

## Investigation: Is it a symbolic link or a hard link?

As a program how do I know if a file is a "regular" file, symbolic link or a hard link? The answer lies in the
data that the `stat()` system call returns. Specifically, the `st_mode` field as described [here](http://man7.org/linux/man-pages/man7/inode.7.html).

## Investigation: Directories and Links

## Using links to solve a problem


## Learning more

- [inode man page](http://man7.org/linux/man-pages/man7/inode.7.html)

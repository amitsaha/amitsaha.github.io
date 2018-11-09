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

asaha@DESKTOP-KBVL52S:~/links-demo$ ls -i
15481123719144131 file1  15481123719144131 file1-hlink  29836347531381846 file1-slink  65583669573645372 file2
asaha@DESKTOP-KBVL52S:~/links-demo$ ls -il
total 0
15481123719144131 -rw-rw-rw- 2 asaha asaha 18 Nov  9 13:52 file1
15481123719144131 -rw-rw-rw- 2 asaha asaha 18 Nov  9 13:52 file1-hlink
29836347531381846 lrwxrwxrwx 1 asaha asaha  5 Nov  9 13:54 file1-slink -> file1
65583669573645372 -rw-rw-rw- 1 asaha asaha 18 Nov  9 13:52 file2
asaha@DESKTOP-KBVL52S:~/links-demo$ rm file2

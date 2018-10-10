
[Much](https://medium.com/@wendymayorgasegura/what-is-the-difference-between-a-hard-link-and-a-symbolic-link-8c0493041b62) [has](https://medium.com/@meghamohan/hard-link-and-symbolic-link-3cad74e5b5dc) [been](https://medium.com/meatandmachines/explaining-the-difference-between-hard-links-symbolic-links-using-bruce-lee-32828832e8d3) written (and [asked](https://stackoverflow.com/questions/185899/what-is-the-difference-between-a-symbolic-link-and-a-hard-link)) on the topic of hard links and soft links on Linux. I have read a few of those more than once . 

However, I end up getting confused between the two, specifically the differences between the two. So, here's 
my post on the topic as I experiment with the same topic - with the hope that I will stop getting confused 
ever again.

## Our setup


ubuntu@ip-172-34-57-238:~/links-demo$ ls -i
780973 real-dir1  780972 real-file1  780974 real-file2
ubuntu@ip-172-34-57-238:~/links-demo$ ln -s soft-link-file1 real-file1
ln: failed to create symbolic link 'real-file1': File exists
ubuntu@ip-172-34-57-238:~/links-demo$ ln real-file1 soft-link-file1
ln: failed to create hard link 'soft-link-file1': File existsubuntu@ip-172-34-57-238:~/links-demo$ ln real-file1 hard-link-file1
ubuntu@ip-172-34-57-238:~/links-demo$ ls -lrt
total 0
-rw-rw-r-- 2 ubuntu ubuntu  0 Oct  9 01:29 real-file1
-rw-rw-r-- 2 ubuntu ubuntu  0 Oct  9 01:29 hard-link-file1
-rw-rw-r-- 1 ubuntu ubuntu  0 Oct  9 01:29 real-dir1
-rw-rw-r-- 1 ubuntu ubuntu  0 Oct  9 01:29 real-file2
lrwxrwxrwx 1 ubuntu ubuntu 10 Oct  9 01:31 soft-link-file1 -> real-file1
ubuntu@ip-172-34-57-238:~/links-demo$ ls -lrtitotal 0780972 -rw-rw-r-- 2 ubuntu ubuntu  0 Oct  9 01:29 real-file1780972 -rw-rw-r-- 2 ubuntu ubuntu  0 Oct  9 01:29 hard-link-file1780974 -rw-rw-r-- 1 ubuntu ubuntu  0 Oct  9 01:29 real-file2780975 lrwxrwxrwx 1 ubuntu ubuntu 10 Oct  9 01:31 soft-link-file1 -> real-file1ubuntu@ip-172-34-57-238:~/links-demo$ ls -litotal 0780972 -rw-rw-r-- 2 ubuntu ubuntu  0 Oct  9 01:29 hard-link-file1780973 -rw-rw-r-- 1 ubuntu ubuntu  0 Oct  9 01:29 real-dir1780972 -rw-rw-r-- 2 ubuntu ubuntu  0 Oct  9 01:29 real-file1780974 -rw-rw-r-- 1 ubuntu ubuntu  0 Oct  9 01:29 real-file2780975 lrwxrwxrwx 1 ubuntu ubuntu 10 Oct  9 01:31 soft-link-file1 -> real-file1

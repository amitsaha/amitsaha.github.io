
In a previous post, I shared that we are going to have Fedora Scientific Vagrant boxes with the upcoming Fedora 29 release.
Few weeks back, I wanted to try out a more recent build to script some of the testing I do on Fedora Scientific boxes
to make sure that the expected libraries/programs are installed. Unexpectedly, `vagrant ssh` would not succeed. 

I filed a issue with [rel-eng](https://pagure.io/releng/issue/7814) where I was suggested to see if a package in
Fedora Scientific was mucking around with the SSH config. To do so, I had to find a way to manually build Vagrant
boxes.

## Building Vagrant boxes with chef/bento

Cpacker build -force -only=virtualbox-iso .\fedora-29-scientific-x86_64.json

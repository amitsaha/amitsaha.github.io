Title: Automatic building and publishing DEB packages for Golang applications 
Date: 2018-02-24 23:00
Category: golang
Status: Draft

In my earlier post, [Quick and dirty debian packages for your Golang application](http://echorand.me/quick-and-dirty-debian-packages-for-your-golang-application.html)
I shared a recipe building DEB packages for Golang applications. We are going to see the following things in this post building
upon our recipe in that post:

- Building the DEB packages in [Travis CI](https://travis-ci.org/amitsaha/golang-packaging-demo)
- Publishing the DEB package to [packagecloud.io](https://packagecloud.io)

The primary assumption in my first post was using [dep](https://golang.github.io/dep/) for dependency management. 
That still holds here.

## Building the DEB packages in Travis CI

To let Travis build the DEB package, we add a `.travis.yml` file to our [git repository](https://github.com/amitsaha/golang-packaging-demo)
with the following contents:

```
# This gives us full control over what we intend to do
# in the job
language: minimal
sudo: required
services:
  - docker
addons:
  apt:
    packages:
      - docker-ce
script:
  - make build-deb DEB_PACKAGE_DESCRIPTION="Logrus Demo" DEB_PACKAGE_NAME=demo BINARY_NAME=demo && ls ./artifacts/

```



## 

```
root@c9b3de968621:/# curl -s https://packagecloud.io/install/repositories/amitsaha/logrus-demo/script.deb.sh | bash
Detected operating system as Ubuntu/xenial.
Checking for curl...
Detected curl...
Checking for gpg...
Detected gpg...
Running apt-get update... done.
Installing apt-transport-https... done.
Installing /etc/apt/sources.list.d/amitsaha_logrus-demo.list...done.
Importing packagecloud gpg key... done.
Running apt-get update... done.

The repository is setup! You can now install packages.
root@c9b3de968621:/# apt install demo
Reading package lists... Done
Building dependency tree
Reading state information... Done
The following NEW packages will be installed:
  demo
0 upgraded, 1 newly installed, 0 to remove and 20 not upgraded.
Need to get 842 kB of archives.
After this operation, 2483 kB of additional disk space will be used.
Get:1 https://packagecloud.io/amitsaha/logrus-demo/ubuntu xenial/main amd64 demo amd64 0.1-e7b1650 [842 kB]
Fetched 842 kB in 5s (143 kB/s)
debconf: delaying package configuration, since apt-utils is not installed
Selecting previously unselected package demo.
(Reading database ... 5291 files and directories currently installed.)
Preparing to unpack .../demo_0.1-e7b1650_amd64.deb ...
Unpacking demo (0.1-e7b1650) ...
Setting up demo (0.1-e7b1650) ...
root@c9b3de968621:/# demo
INFO[0000] I love logrus!

```

:Title: Getting started with Golang on Linux
:Date: 2016-06-17 17:00
:Category: Golang

This guide will be how I usually setup and get started with Go development environment on Linux. By the end of this document, we will have seen how to:

- Install the Go compiler and other tools (`gofmt`, for eaxmple), collectively referred to as go tools
- Setup Go workspace
- Working with Go programs using third-party packages

Installing Go tools
===================

The first option to install the Go compiler and other tools from your distro's package manager. On Fedora, you can do ``sudo dnf -y install golang`` for example. However, very likely your distro's packaged version will be behind the upstream release. Hence, we will follow the official `install guide <https://golang.org/doc/install>`__ to get the latest stable version of the Go tools:

- Download the Linux binary tarball from the `Downloads page <https://golang.org/dl/>`__
- ``sudo tar -C /usr/local -xzf <filename-from-above>``
- ``export PATH=$PATH:/usr/local/go/bin`` in your ``.bashrc`` or similar file.

When we now open a new terminal session, we should be able to type in `go version` and get the version we installed:

.. code::
   
   $ go version
   go version go1.6.2 linux/amd64

If we see this, we are all set to go to the next stage

Setting up the Go workspace
===========================

Golang expects us to structure our source code in a certain way. You can read all about it in this `document <https://golang.org/doc/code.html>`__. The summarized version is that:

- All our go code (including those of packages we use) in a single directory
- The environment variable ``GOPATH`` points to this single directory
- This single directory has three sub-directories: ``src``, ``bin``, ``pkg``
- It is in the ``src`` sub-directory where all our Go code will live

For this guide I will assume that the ``GOPATH`` is set to ``$HOME/work/golang``:

.. code::

   $ mkdir -p $HOME/work/golang
   $ mkdir -p $HOME/work/golang/src $HOME/work/golang/bin $HOME/work/golang/pkg
   
At this stage, our $GOPATH directory tree looks like this:

.. code::

   $ tree $GOPATH -L 1
   /home/vagrant/work/golang
   ├── bin
   ├── pkg
   └── src

Next, we will add the line ``export GOPATH=$HOME/work/golang`` in the ``.bashrc`` (or another similar file). If we now start a new terminal session, we should see that ``GOPATH`` is now setup correctly:

.. code::
   
   $ go env GOPATH
   /home/vagrant/work/golang

You can learn more about GOPATH `here <https://golang.org/cmd/go/#hdr-GOPATH_environment_variable>`__.

Writing our first program
=========================

There are two types of Golang programs we can write - one is an application program (output is an executable program) and the other is a package which is meant to be used in other programs. We will first write a program which will be compiled to an executable:








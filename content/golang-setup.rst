:Title: Setting up Golang on Linux
:Date: 2016-06-17 17:00
:Category: Golang

This guide will be how I usually setup my Go development environment on Linux. By the end of this document, we will have seen how to:

- Install the Go compiler and other tools (`gofmt`, for eaxmple), collectively referred to as go tools
- Setup Go workspace
- Working with Go programs using third-party packages

Installing Go tools
===================

The first option to install the Go compiler and other tools from your distro's package manager. On Fedora, you can do ``sudo dnf -y install golang`` for example. However, very likely your distro's packaged version will be behind the upstream release. Hence, we will follow the official `install guide <https://golang.org/doc/install>`__ to get the latest stable version of the Go tools.

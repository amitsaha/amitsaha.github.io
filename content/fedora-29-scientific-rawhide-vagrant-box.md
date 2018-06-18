Title: Pre-release Fedora Scientific Vagrant Boxes
Date: 2018-06-18 11.00
Category: fedora
Status: draft

I am very excited to share that sometime back the Fedora project gave the go ahead on my [idea](https://fedoraproject.org/wiki/Changes/FedoraScientific_VagrantBox) of making Fedora Scientific
available as [Vagrant](https://www.vagrantup.com/) boxes starting with Fedora 29. This basically 
means (I think) that using Fedora Scientific in a virtual machine is even easier. Instead of downloading 
the ISO and then going through the installation process, you can now basically do:

- Download Fedora Scientific Vagrant box
- Initialize VM
- `vagrant up`
    
    
## Trying it out before Fedora 29 release on Mac OS X Host

As of a few days back, Fedora 29 rawhide vagrant boxes for Fedora Scientific are now being published. Thanks to release
engineering for moving this forward. 

If you are keen to try it out, here's what I did using VirtualBox on an OS X host. 

First, install [vagrant](https://www.vagrantup.com/). Then, from a terminal:

```
# Add the box
$ vagrant box add https://kojipkgs.fedoraproject.org//packages/Fedora-Scientific-Vagrant/Rawhide/20180613.n.0/images/Fedora-Scientific-Vagrant-Rawhide-20180613.n.0.x86_64.vagrant-virtualbox.box  --name Fedora-Scientific-Rawhide
==> box: Box file was not detected as metadata. Adding it directly...
==> box: Adding box 'Fedora-Scientific-Rawhide' (v0) for provider: 
    box: Downloading: https://kojipkgs.fedoraproject.org//packages/Fedora-Scientific-Vagrant/Rawhide/20180613.n.0/images/Fedora-Scientific-Vagrant-Rawhide-20180613.n.0.x86_64.vagrant-virtualbox.box
==> box: Box download is resuming from prior download progress
==> box: Successfully added box 'Fedora-Scientific-Rawhide' (v0) for 'virtualbox'!

..
```

Now that the box has been downloaded, we initialize a new VM:

```
$ vagrant init Fedora-Scientific-Rawhide
A `Vagrantfile` has been placed in this directory. You are now
ready to `vagrant up` your first virtual environment! Please read
the comments in the Vagrantfile as well as documentation on
`vagrantup.com` for more information on using Vagrant.
MacBook-Air:Downloads amit$ vagrant up
Bringing machine 'default' up with 'virtualbox' provider...
==> default: Importing base box 'Fedora-Scientific-Rawhide'...
Progress: 70%
Progress: 90%
..


==> default: Machine booted and ready!
..
```

If you see the above message, we are ready to start using Fedora Scientific:

```
$ vagrant ssh
```

The above will drop us into a terminal session on our Fedora Scientific VM. To be able to use graphical programs, we
will change the above command to (please note I also needed to install [xquartz](https://www.xquartz.org/) to be able to see graphical programs from the VM on my host):

```
$ vagrant ssh -- -X
```

## Windows hosts

To be done (If you end up doing it, please let me know - see for a link at the bottom of this post).

## Linux hosts

To be done (If you end up doing it, please let me know - see for a link at the bottom of this post).

## Explore

Now that we have Fedora Scientific setup, head over to the [docs](http://fedora-scientific.readthedocs.io/en/latest/) to 
explore what's in Fedora Scientific.

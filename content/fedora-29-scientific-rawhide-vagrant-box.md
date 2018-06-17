Title: Fedora Scientific Vagrant Boxes
Date: 2018-06-18 11.00
Category: fedora
Status: draft

I am very excited to share that sometime back the Fedora project gave the go ahead on my [idea](https://fedoraproject.org/wiki/Changes/FedoraScientific_VagrantBox) of making Fedora Scientific
available as [Vagrant](https://www.vagrantup.com/) boxes starting with Fedora 29. This basically means (I think) that using Fedora Scientific in a virtual
machine is even easier. Instead of downloading the ISO and then going through the installation process, you can now basically
do:

- Download Fedora Scientific Vagrant box
- <initialization step>
- `vagrant up`
    
    
## Trying it out

As of a few days back, Fedora 29 rawhide vagrant boxes for Fedora Scientific are now being published. Thanks to release
engineering for moving this forward. 

If you are keen to try it out, here's what I did using VirtualBox on an OS X host. 

First, install [vagrant](https://www.vagrantup.com/). Then:

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

Now that the box has been downloaded:

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




==> default: Matching MAC address for NAT networking...
==> default: Setting the name of the VM: Downloads_default_1528978126658_94807
==> default: Clearing any previously set network interfaces...
==> default: Preparing network interfaces based on configuration...
    default: Adapter 1: nat
==> default: Forwarding ports...
    default: 22 (guest) => 2222 (host) (adapter 1)
==> default: Booting VM...

==> default: Waiting for machine to boot. This may take a few minutes...
    default: SSH address: 127.0.0.1:2222
    default: SSH username: vagrant
    default: SSH auth method: private key
    default: 
    default: Vagrant insecure key detected. Vagrant will automatically replace
    default: this with a newly generated keypair for better security.
    default: 
    default: Inserting generated public key within guest...
    default: Removing insecure key from the guest if it's present...
    default: Key inserted! Disconnecting and reconnecting using new SSH key...
==> default: Machine booted and ready!


# -*- mode: ruby -*-
# vi: set ft=ruby :

# All Vagrant configuration is done below. The "2" in Vagrant.configure
# configures the configuration version (we support older styles for
# backwards compatibility). Please don't change it unless you know what
# you're doing.
Vagrant.configure("2") do |config|
  # The most common configuration options are documented and commented below.
  # For a complete reference, please see the online documentation at
  # https://docs.vagrantup.com.

  # Every Vagrant development environment requires a box. You can search for
  # boxes at https://vagrantcloud.com/search.
  config.vm.box = "Fedora-Scientific-Rawhide"

  # Disable automatic box update checking. If you disable this, then
  # boxes will only be checked for updates when the user runs
  # `vagrant box outdated`. This is not recommended.
  # config.vm.box_check_update = false

  # Create a forwarded port mapping which allows access to a specific port
  # within the machine from a port on the host machine. In the example below,
  # accessing "localhost:8080" will access port 80 on the guest machine.
  # NOTE: This will enable public access to the opened port
  # config.vm.network "forwarded_port", guest: 80, host: 8080

  # Create a forwarded port mapping which allows access to a specific port
  # within the machine from a port on the host machine and only allow access
  # via 127.0.0.1 to disable public access
  # config.vm.network "forwarded_port", guest: 80, host: 8080, host_ip: "127.0.0.1"

  # Create a private network, which allows host-only access to the machine
  # using a specific IP.
  # config.vm.network "private_network", ip: "192.168.33.10"

  # Create a public network, which generally matched to bridged network.
  # Bridged networks make the machine appear as another physical device on
  # your network.
  # config.vm.network "public_network"

  # Share an additional folder to the guest VM. The first argument is
  # the path on the host to the actual folder. The second argument is
  # the path on the guest to mount the folder. And the optional third
  # argument is a set of non-required options.
  # config.vm.synced_folder "../data", "/vagrant_data"

  # Provider-specific configuration so you can fine-tune various
  # backing providers for Vagrant. These expose provider-specific options.
  # Example for VirtualBox:
  #
  # config.vm.provider "virtualbox" do |vb|
  #   # Display the VirtualBox GUI when booting the machine
  #   vb.gui = true
  #
  #   # Customize the amount of memory on the VM:
  #   vb.memory = "1024"
  # end
  #
  # View the documentation for the provider you are using for more
  # information on available options.

  # Enable provisioning with a shell script. Additional provisioners such as
  # Puppet, Chef, Ansible, Salt, and Docker are also available. Please see the
  # documentation for more information about their specific syntax and use.
  # config.vm.provision "shell", inline: <<-SHELL
  #   apt-get update
  #   apt-get install -y apache2
  # SHELL
end



```
vagrant ssh -- -X
```

https://www.xquartz.org/


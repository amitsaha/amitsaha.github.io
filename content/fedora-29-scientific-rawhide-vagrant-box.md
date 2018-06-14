$ vagrant box add https://kojipkgs.fedoraproject.org//packages/Fedora-Scientific-Vagrant/Rawhide/20180613.n.0/images/Fedora-Scientific-Vagrant-Rawhide-20180613.n.0.x86_64.vagrant-virtualbox.box  --name Fedora-Scientific-Rawhide
==> box: Box file was not detected as metadata. Adding it directly...
==> box: Adding box 'Fedora-Scientific-Rawhide' (v0) for provider: 
    box: Downloading: https://kojipkgs.fedoraproject.org//packages/Fedora-Scientific-Vagrant/Rawhide/20180613.n.0/images/Fedora-Scientific-Vagrant-Rawhide-20180613.n.0.x86_64.vagrant-virtualbox.box
==> box: Box download is resuming from prior download progress
==> box: Successfully added box 'Fedora-Scientific-Rawhide' (v0) for 'virtualbox'!
MacBook-Air:Downloads amit$ vagrant init Fedora-Scientific-Rawhide
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

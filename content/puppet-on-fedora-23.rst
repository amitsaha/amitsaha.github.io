:Title: Setting up a test puppet system on Fedora
:Date: 2015-10-01 10:20
:Category: Fedora

Setup
=====

.. code::
  
   # `dnf -y install puppet`

Setup the host name:

.. code::

   $ cat /etc/hostname 
   fedora-23.node

Reboot, verify:

.. code::
  $ facter | grep node
  domain => node
  fqdn => fedora-23.node

.. code::

   # tree /etc/puppet/
   /etc/puppet/
    ├── auth.conf
    ├── modules
    └── puppet.conf

.. code::

   # mkdir /etc/puppet/manifests

Create our first manifest `/etc/puppet/manifests/nginx.pp`:

.. code::

   node "fedora-23.node" {
      package { "nginx":
        ensure => installed
      }
   }


Apply with ``--noop``:

.. code::

  # puppet apply nginx.pp --noop
  Notice: Compiled catalog for fedora-23.node in environment production in 0.66 seconds
  Notice: /Stage[main]/Main/Node[fedora-23.node]/Package[nginx]/ensure: current_value purged, should be present (noop)
  Notice: Node[fedora-23.node]: Would have triggered 'refresh' from 1 events
  Notice: Class[Main]: Would have triggered 'refresh' from 1 events
  Notice: Stage[main]: Would have triggered 'refresh' from 1 events
  Notice: Applied catalog in 0.26 seconds

Really apply:

.. code::

   # puppet apply nginx.pp
   Notice: Compiled catalog for fedora-23.node in environment production in 0.60 seconds
   Notice: /Stage[main]/Main/Node[fedora-23.node]/Package[nginx]/ensure: created
   Notice: Applied catalog in 5.67 seconds


.. code::
   # rpm -q nginx
   nginx-1.8.0-13.fc23.x86_64

.. code::
   
   # puppet apply nginx.pp --noop
   Notice: Compiled catalog for fedora-23.node in environment production in 0.59 seconds
   Notice: Applied catalog in 0.25 seconds




Resources
=========

- https://docs.puppetlabs.com/references/latest/type.html#package
- https://www.digitalocean.com/community/tutorials/how-to-install-puppet-in-standalone-mode-on-centos-7



Title: Puppet
Date: 2018-01-05 10:00


**Puppet version**

```
$ puppet --version
5.3.3
```

**Show diff**

```
$ sudo puppet apply --noop --show_diff /examples/file_hello.pp
```

**Query resources**

```
$ sudo puppet resource package openssl
```

`puppet resource package` will display a list of all packages

`puppet resource -e package <name>` will bring up an interactive editor

**require and notify**

`require` and `notify` both takes the resource type and name, but the resource type must be specified as `Resource`. For example, `Package`, `Service`, `File`.

**puppet_gem provider**

The `puppet_gem` provider for the `package` resource is used to install a Ruby gem for puppet's context and doesn't install the gem in for system Ruby.

**service type**

The `service` type has attributed such as `pattern` which allows us to check as an alternative to using the service commands to check for service status and instead check the process table for running processes of the same name as that of the service.

**cron resource**

`cron` resource supports the `fqdn_rand` function so that the cron job doesn't run at the exact same time on all nodes.


**refreshonly**

Setting the `refreshonly` attribute to a resource such as `exec` only triggers the execution of the `exec` resource if notified by another event.

**notice**

The `notice()` function can be used to print strings/values of variables useful for debugging:

```
$name = "John"
notice("Hello ${name}")
```

**Arrays and resources**

Specifying an array of strings to a `resource` declaration as the title creates a resource for each element of the array. Example (from book):

```
$dependencies = [
  'php7.0-cgi',
  'php7.0-cli',
  'php7.0-common',
  'php7.0-gd',
  'php7.0-json',
  'php7.0-mcrypt',
  'php7.0-mysql',
  'php7.0-soap',
]

package { $dependencies:
  ensure => installed,
}
```

**Declaring a hash**

```
$heights = {
  'john'    => 193,
  ..
}

notice("John's height is ${heights['john']}cm.")
```

**Attribute splat operator**

Deconstructing a hash and passing it as parameters to a hash:

```
$attributes = {
  'owner' => 'john',
  'permissions' => 0777.
}


file {'/tmp/hello':
  ensure => present,
  * => $attributes,
 }
```

**Comparison operators**

- `==`: Equal to
- `!=`: Not equal to
- `>`, '>=`, '<`, `<=`: Greater/Greater or equal, Less/lesser or equal
- `=~`: Match a [regular expression](http://ruby-doc.org/core/Regexp.html)


**Control flow constructs**


Example of `if`:

```
$install_package = true
if $install_package {
  package { $package_name:
    ensure => installed,
  }
} else {
  package { $package_name:
    ensure => absent,
  }
}
```

Puppet also has "switch" case:

```
$os_name = 'ubuntu'

case $os_name {
  'ubuntu': {
    $pkg_manager = 'apt-get'
  }
  'fedora': {
    $pkg_manager = 'dnf'
   }
   default: {
     $pkg_manager = false
   }
 }
..
```

**Facter and Retrieving facts**

Puppet makes the system information available via `facter` as an array `$facts` which is available to all puppet manifests.

Using the `$facts` array is as follows:

```
if $facts['os']['selinux']['enabled'] {
  notice('SELinux is enabled')
} else {
  notice('SELinux is disabled')
}
```

The `facter` command line program can be used to query various system facts:

```
$ facter os
{
  architecture => "amd64",
  distro => {
    codename => "xenial",
    description => "Ubuntu 16.04.3 LTS",
    id => "Ubuntu",
    release => {
      full => "16.04",
      major => "16.04"
    }
  },
  family => "Debian",
  hardware => "x86_64",
  name => "Ubuntu",
  release => {
    full => "16.04",
    major => "16.04"
  },
  selinux => {
    enabled => false
  }
}
```

Another example:

```
$ facter kernel
Linux
```

Use `$ puppet facts` to show the facts available to puppet manifests. This includes custom facts loaded from modules.

**Setting custom facts**

Quick way to set custom fact is via an environment variable, like so:

```
$ FACTER_myapp_version=0.1 facter myapp_version 
0.1
```

External facts can be specified via dropping files in `/opt/puppetlabs/facter/facts.d/`, like so:

```
myapp_version=0.1
```

If we mark the [external file](https://puppet.com/docs/facter/3.9/custom_facts.html#executable-facts-----unix) as executable, it will be executed - the only condition is that the output must be a list
of key value pairs. For example:

```
key1=value1
key2=value2
key3=value3
```


**Iteration**

Example code:

```
$ sudo cat /etc/puppetlabs/code/environments/production/manifests/iteration_example.pp
$my_partitions = $facts['partitions']
$my_partitions.each | String $device, Hash $attrs | {
   notice("Device ${device} has attributes ${attrs}")
}

# apply
Notice: Scope(Class[main]): Device /dev/sda1 has attributes {filesystem => ext4, label => cloudimg-rootfs, mount => /, partuuid => c5572b29-01, size => 10.00 GiB, size_bytes => 10736352768, uuid => 09cea002-3bc0-41b2-81b5-2cbf8b138eb7}

```


Note in the above I can't use `$partitions` as the variable name, since `partitions` is a fact and is a [classic fact name fact](https://puppet.com/docs/puppet/5.3/lang_facts_and_builtin_vars.html#classic-factname-facts)


**Look up hiera data**

```
$ sudo puppet lookup --environment pbg test
..
```

In manifests, we use the `lookup()` function to look up values of keys:

```
$numworkers = lookup('apache_workers', int)

```

The second argument allows us to specify an expected type of the value we are looking up. Examples including `int`, `Boolean`, `Hash`.

To apply the manifest to a certain environment:

```
$ sudo puppet apply --environment <env> <manifest>
..
```

Hiera values can be: single values, arrays and hashes. Single values:

```
package_version: 1.1
```

Array:

```
packages:
  - vim
  - docker
```

Hash:

```
services:
  webapp: 8080
  db: 3306
  master: true
```

Directly looking up a value in a hash from hiera:

```
$value = lookup('services.webapp', String)
..
```

**Interpolation in Hiera data**

We can look up another variable in a hiera file:

```
 backup_path: "/backup/%{facts.hostname}"
```

We can also look up another variable value defined in hiera:

```
package_name: 'mypackage'
package_version: '1.1'
package: "%{lookup('package_name')}-%{lookup('package_version')}"
```

**Creating resources with Hiera data**

Let's consider the following hiera data:

```
packages:
  'htop':
    version: latest
  'nginx':
    version: 1.6
```

We can then write the following manifest to install the above packages.

```
lookup('packages', Hash, 'hash').each | String $package, Hash $attrs | {
  package { $package:
    ensure => $attrs['version'],
  }
}
```


**Storing secrets in Hiera**

- Use gnupg
- Create a key
- Encrypt your secrets using that key
- Distribute the decrypting key across all nodes that needs access to the secret

**Using puppet forge modules**

`r10k` is the defactor standard module manager for Puppet deployments.

- Create a `Puppetfile` listing the modules we want to use
- Use `r10k` to install it at "build" time and deploy it

```
$ sudo r10k puppetfile install --verbose
$ sudo r10k puppetfile install --verbose
INFO	 -> Updating module /etc/puppetlabs/code/environments/pbg/modules/docker
INFO	 -> Updating module /etc/puppetlabs/code/environments/pbg/modules/archive
INFO	 -> Updating module /etc/puppetlabs/code/environments/pbg/modules/staging
INFO	 -> Updating module /etc/puppetlabs/code/environments/pbg/modules/apache
..
```

This will create a `modules` sub-directory with the downloaded modules. Using the `--moduledir`
option, we can specify a different sub-directory to download the modules.

`r10k` doesn't automatically manage dependencies between modules. That is if `module1` depends on
`module2`, it will not automatically pull it. [puppet-generate-puppetfile](https://github.com/rnelson0/puppet-generate-puppetfile)
is a useful tool to help to manage this. We can use it to update dependency versions as well.

Learn about [module path](https://puppet.com/docs/puppet/5.2/dirs_modulepath.html) here.

The puppet [stdlib](https://github.com/puppetlabs/puppetlabs-stdlib) module provides functions to provide
various functionalities.

**Classes**

Definition:

```
class mymodule {
}
```

Declaration:

```
include mymodule
```


Defining class parameters:

```
class mymodule(
  String $version = 'installed',
```

Specifying a default value here means that if a declaration doesn't specify the value, `puppet` doesn't
error out. The type `String` also ensures non-string values error out.

A `String[1]` ensures a string value of atleast 1 character.

We can specify a list of allowed values via the `Enum` datatype:

```
Enum['running', 'stopped'] $service_state = 'running',
```

**Parameter data types**

- `Integer`
- `Float`
- `Numeric`
- `Array`
- `Hash`
- `Regexp`
- `Undef`
- `Type`

Hashes/Arrays can also specify the type of the elements: `Array[Integer]` or `Hash[String, Integer]`

Abstract data types

- `Optional`
- `Pattern`
- `Scalar`
- `Data`
- `Variant`

[documentation](https://puppet.com/docs/puppet/5.3/lang_data_abstract.html)

**Defined resource types**

A defined resource type allows grouping of resources:

```
define myresource {
...
}
```

The difference between a defined resource type and `class` is that a class can be declared
only once.

**Type alias**

```
type ServiceName = Pattern[..]
class myclass(
  Array[ServiceName] $services,
){
...
}
```

**Roles and Profiles**

`Role` identify a particular function for a node - an app server, a DB server, etc.

`Profile` identify a particular piece of software or functionality which contributes
to a role.

**Templates**

- `Embedded Puppet` - has the extension `epp`
- [epp v$ erb](http://puppet-on-the-edge.blogspot.com.au/2014/03/templating-with-embedded-puppet.html)

Templates from files:

```
file { '/etc/systemd/system/myunit.service':
  content => epp('module_name/myunit.service.epp'),
}
```

Puppet will look for the `myunit.service.epp` file in the `templates` sub-directory
of the module automatically.


**Validating template syntax**: `puppet epp validate /examples/template.epp`

**Rendering template on the command line**: `puppet epp render --values_file <params.pp> <template>`



 


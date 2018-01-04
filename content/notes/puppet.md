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




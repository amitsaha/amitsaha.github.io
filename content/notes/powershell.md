Title: Powershell
Date: 2018-04-01 10:00

# Introduction

If you open powershell on Windows and are coming from a Linux background, it's mildly pleasantly surprising to find
commands like `cat`, `ls`, `mkdir` and `cd` work. However, there's no point in continuing to use those since learning
Powershell itself is an intellectually rewarding experience I think. Let's jump right in with the things I usually
end up doing.

# Environment variables

Set an environment variable from the command line or script:

```
$Env:AWS_PROFILE = "profilename"
```

Read an environment variable from the command line or script:

```
"The AWS profile is $Env:AWS_PROFILE"
```

Learn all about it [here](https://docs.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_environment_variables?view=powershell-6)

Show all environment variables:

```
> Get-ChildItem Env:
```

# Variable substitution

Learn all about it [here](https://kevinmarquette.github.io/2017-01-13-powershell-variable-substitution-in-strings/)

# Create an empty file

```
> New-Item <path>  -ItemType file
```

# Time a command

```
> Measure-Command { Get-ChildItem | Out-Host }
```

The `Out-Host` pipe is required so as to not gobble up the output of your command.

# Tail follow a file

```
> Get-Content -Path <file> -Wait
```

# Display first/last N lines of a file

```
> Get-Content -Path -Head 10
> Get-Content -Path -Tail 10
```

# Where is my current command/program being invoked from?

```
> Get-Command <cmd>
```

It's also alised to `gcm`.


Learn all about it [here](https://docs.microsoft.com/en-us/powershell/module/microsoft.powershell.core/get-command?view=powershell-6)

# Display directory contents

```
> Get-ChildItem
```

Show only directories:

```
> Get-ChildItem -Directory
```

Show only files:

```
> Get-ChildItem -File
```

# Remove a file/directory

```
> Remove-Item <path>
```

To delete the contents of a directory as well recursively, use `-Recurse`

Learn all about it [here](https://docs.microsoft.com/en-us/powershell/module/microsoft.powershell.core/providers/filesystem-provider/remove-item-for-filesystem?view=powershell-6)

# Run an external EXE and wait for it to finish

```
> <path to exe> | Out-Null
```

Learn [more](https://stackoverflow.com/questions/1741490/how-to-tell-powershell-to-wait-for-each-command-to-end-before-starting-the-next)

# System DNS

Get DNS settings:

```
> Get-DnsClientServerAddress

InterfaceAlias               Interface Address ServerAddresses
                             Index     Family
--------------               --------- ------- ---------------
Ethernet 2                          11 IPv4    {172.34.0.2}
```

# Current location management

Change to a `path`:

```
> Set-Location -Path <path>
```

Get current location:

```
> Get-Location
```

Push to a location temporarily:

```
> Push-Location -Path <path>
```

Pop out:

```
> Pop-Location
```

Learn all about it [here](https://docs.microsoft.com/en-us/powershell/scripting/getting-started/cookbooks/managing-current-location?view=powershell-6).

# Redirecting output

Learn all about it [here](https://docs.microsoft.com/en-us/powershell/scripting/getting-started/cookbooks/redirecting-data-with-out---cmdlets?view=powershell-6).

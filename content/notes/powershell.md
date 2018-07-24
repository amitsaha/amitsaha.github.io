Title: Powershell
Date: 2018-04-01 10:00

# Introduction

If you open powershell on Windows and are coming from a Linux background, it's mildly pleasantly surprising to find
commands like `cat`, `ls`, `mkdir` and `cd` work. However, there's no point in continuing to use those since learning
Powershell itself is an intellectually rewarding experience I think. Let's jump right in with the things I usually
end up doing.

# Static Analysis

Use [PSScriptAnalyzer](https://github.com/PowerShell/PSScriptAnalyzer)

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

# Strings

A simple string:

```
> 'Hello'
Hello
```

Substitute variable values:

```
> '{0}-{1}' -f "hello", "world"
hello-world

```

Substitute the result of a command in a string:

```
> '{0}' -f $(Get-Location)
...
```

Learn more [here](http://jon.netdork.net/2015/09/08/powershell-and-single-vs-double-quotes/)

# Variable substitution

Learn all about it [here](https://kevinmarquette.github.io/2017-01-13-powershell-variable-substitution-in-strings/)

For property substitution:

```
$Results = Get-ChildItem  -Filter *.xml | Select-Object FullName
Write-Output "Parsing:  $($result.FullName)"
```

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


# Reading/parsing text/xml/json

```
> $lines = Get-Content <path>
```

`lines` is an array with each line of the file an array element

See [this post](https://blog.stangroome.com/2014/02/10/powershell-select-xml-versus-get-content/) for parsing XML files
and [this post](https://blogs.technet.microsoft.com/heyscriptingguy/2015/10/08/playing-with-json-and-powershell/) for parsing
JSON files.

# Write output

```
> Write-Host "Hello World"
```

See [this](https://docs.microsoft.com/en-us/powershell/module/microsoft.powershell.utility/write-host?view=powershell-6) 
and [this](https://stackoverflow.com/questions/8755497/which-should-i-use-write-host-write-output-or-consolewriteline) to learn more.

# TLS/SSL Error

```
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
```

# Extracting a ZIP archive

```
> Expand-Archive <path to zip>
```

Learn all about it [here](https://docs.microsoft.com/en-us/powershell/module/Microsoft.PowerShell.Archive/Expand-Archive?view=powershell-5.1)

# Find the hash value of a file

Calculate the sha256 sum:

```
> Get-FileHash <path to file>
```

Also supports other algorithms and input stream. Learn all about it [here](https://docs.microsoft.com/en-us/powershell/module/microsoft.powershell.utility/get-filehash?view=powershell-6)


# Sleep

Sleep for n seconds:

```
> Start-Sleep -s n
```
Learn all about it [here](https://docs.microsoft.com/en-us/powershell/module/microsoft.powershell.utility/start-sleep?view=powershell-6)

# Working with DNS clients/servers

Clear DNS cache:

```
> Clear-DnsClientCache
```

Get current DNS servers:

```
> Get-DnsClientServerAddress
```
https://docs.microsoft.com/en-us/powershell/module/dnsclient/get-dnsclientserveraddress?view=win10-ps


Set DNS address:

```
> Set-DnsClientServerAddress -InterfaceAlias ..
```

Learn all about it [here](https://docs.microsoft.com/en-us/powershell/module/dnsclient/set-dnsclientserveraddress?view=win10-ps)

Add a DNS server to the existing addresses:

```
> $result = Get-DnsClientServerAddress -InterfaceAlias "vEthernet" # or InterfaceIndex
> Set-DnsClientServerAddress -InterfaceAlias "vEthernet" -ServerAddresses ("172.34.0.2", $result.ServerAddresses)
```


# Working with list of things

Samples script to go over a list of things and doing something:

```
$repositories = @(
    "demo/repo1",
    "demo/repo2"    
)

foreach ($repo in $repositories)
{
    aws ecr create-repository --repository-name $repo
}
```

# Sourcing stuff from other files

```
. ./<path-to-ps1> file
```

# Update a AWS ECR policy using AWS CLI

The reason I have it here is the escaping for the JSON policy for powershell:

```
# Create repository policy to allow public pulling
# https://docs.aws.amazon.com/cli/latest/userguide/cli-using-param.html#quoting-strings
$policy=@"
{
     \"Version\": \"2008-10-17\",
     \"Statement\": [
         {
             \"Sid\": \"Allow public pulling of images\",
             \"Effect\": \"Allow\",
             \"Principal\": \"*\",
             \"Action\": [
                 \"ecr:GetDownloadUrlForLayer\",
                 \"ecr:BatchGetImage\",
                 \"ecr:BatchCheckLayerAvailability\"
             ]
         }
     ]
 }
"@

foreach ($repo in $repositories)
{
    aws ecr set-repository-policy --repository-name $repo --policy-text $policy
}
```


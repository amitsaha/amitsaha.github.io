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





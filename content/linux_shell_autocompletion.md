What happens when you press `<your favorite Linux command> <TAB>`? You get a bunch of suggestions with one of them
just the one you had in mind you were going to type. Auto completions in shells are super helpful and that's
probably the most boring sales pitch. Anyway, as we will find out, quite a bit goes on behind the scenes. 

Before we go too further into my investigations and findings, I  must credit the author of this [blog post](https://www.joshmcguigan.com/blog/shell-completions-pure-rust/) - it triggered
my curiosity and made me find time out to learn more about something I use everyday, but don't know much about how
it worked.

Let's learn more about the minions that gets to work when we press the `<TAB>` key on a popular Linux shell, `BASH`.

## Setting up

Let's get a fresh Fedora 30 VM in a Vagrant box and set it up:

```
$ vagrant up
...
$ vagrant ssh
 $ sudo dnf remove bash-completion
 $ sudo dnf install bpftrace
 $ curl https://raw.githubusercontent.com/iovisor/bpftrace/master/tools/execsnoop.bt -o execsnoop.bt
```

On a terminal, (Terminal 1), type in `$ git <TAB>`. The only suggestions you get will be the files in the
directory you are in. Not very helpful suggestions, you say. I know - and that's because we uninstalled a 
RPM package which would have magically done that for us. (`bash-completion`). Let's keep it that way for 
us now.

## DIY bash completion

When we type in `$git <TAB>`, we expect to see suggestions of the git sub-commands, such as `status`, `clone`,
and `checkout`. Let's aim for achieving that.

First, create a file, `/tmp/git_suggestions`, put in the following and make it executeable (`chmod +x /tmp/git_suggestions`):

```

#!/bin/bash

echo "status"
echo "checkout"
echo "clone"
echo "branch"
printenv COMP_KEY COMP_LINE COMP_POINT COMP_TYPE

```

This script prints four git subcommands - one on each line. Now, execute the command from Terminal 1:

```
$  complete -C "/tmp/git_suggestions" git
```

Next, type in `git <TAB>`, you will see that you are now suggested four options for your next command/option:

```
[vagrant@ip-10-0-2-15 temp]$ git  <TAB>
branch    checkout  clone     status    
[vagrant@ip-10-0-2-15 temp]$ git 

```

Note that each of the suggestion is a line printed by the above script. Let's delve further into this.

## Exec snooping 

Open a new terminal (Terminal 2) and execute the following:


```
$ sudo bpftrace ./execnsoop.bt
```

Now, go back to Terminal 1, and type `$git <TAB>`. 

On Terminal 2, you will see something like:

```
49916      15949 /tmp/git_suggestions git
```

Let's break this down:

Ignore the first column, which is how long the external program executed for. The second column gives us
the process ID and the third column shows us the external program along with the arguments it was executed it.
We can see that the script `/tmp/git_suggestions` is supplied


Terminal 2:

$ git <TAB>

On Terminal 1, we will see nothing. The suggestions we get in Terminal 2 are the boring current files (which seem
to be some sort of built-bash default) and not the git sub-commands we are so used to.


Now, on Terminal 2, let's install the `bash-completion` package and then exec bash and type, in `$ git <TAB>`


On Terminal 1, now we will see something like:

```
95607      1470  git --list-cmds=list-mainporcelain,others,nohelpers,alias,list-complete,config
```

Let's try another command:

$ systemctl status <TAB>

Terminal 2:

143342     1479  systemctl --system --full --no-legend --no-pager list-unit-files *
143607     1480  systemctl --system --full --no-legend --no-pager list-units --all *



And a third command om terminal 1:

$ service <TAB>

On Terminal 2:


134569     2530  awk $1 ~ /\.service$/ { sub("\\.service$", "", $1); print $1 }
134571     2529  systemctl list-units --full --all
134776     2536  awk $1 ~ /\.service$/ { sub("\\.service$", "", $1); print $1 }
134780     2535  systemctl list-units --full --all





https://www.joshmcguigan.com/blog/shell-completions-pure-rust/

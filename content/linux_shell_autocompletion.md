What happens when you press `<command> <TAB>` or `<command> <TAB><TAB>`? You get a bunch of suggestions with one of them
just the one you had in mind you were going to type. Auto completions in shells are super helpful and as we will find
out, quite a bit goes on behind the scenes. Before we go too further into my investigations and findings, I  must credit the author of 
this [blog post](https://www.joshmcguigan.com/blog/shell-completions-pure-rust/) - it triggered
my curiosity and made me find time out to learn more about something I use everyday, but didn't know much about how
it worked.

Let's learn more about the minions that gets to work when we press `<TAB>` or `<TAB><TAB>` on a popular Linux shell, `BASH`.

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

On a terminal, (Terminal 1), type in `$ git <TAB><TAB>` (the space between `git` and `<TAB><TAB>` is important). 
The only suggestions you get will be the files in the
directory you are in. Not very helpful suggestions, you say. I know - and that's because we uninstalled a 
package which would have magically done that for us. (`bash-completion`). Let's keep it that way for 
now.

## DIY bash completion

When we type in `$git <TAB><TAB>`, we expect to see suggestions of the git sub-commands, such as `status`, `clone`,
and `checkout`. Let's aim to achieve that.

First, create a file, `/tmp/git_suggestions`, put in the following snippet and make it executeable 
(`chmod +x /tmp/git_suggestions`):

```

#!/bin/bash

echo "status"
echo "checkout"
echo "clone"
echo "branch"
```

This script prints four git subcommands - one on each line. Now, execute the command from Terminal 1:

```
$  complete -C "/tmp/git_suggestions" git
```

The `complete` bash built-in is a way to tell bash what we want it to do when we ask for completing a command.
The `-C` option asks it to call a specified external program. There are various other options some of which we
will learn later on in this post.

Next, type in `git <TAB><TAB>`, you will see that you are now suggested four options for your next command/option:

```
[vagrant@ip-10-0-2-15 temp]$ git <TAB><TAB>
branch    checkout  clone     status    

```

Note that each of the suggestion is a line printed by the above script. Let's delve further into this.

## Exec snooping 

Open a new terminal (Terminal 2) and execute the following:

```
$ sudo bpftrace ./execnsoop.bt
```

Now, go back to Terminal 1, and type `$git <TAB><TAB>`. 

On Terminal 2, you will see something like:

```
49916      15949 /tmp/git_suggestions git
```

The first column is how long the external program executed for in milliseconds (the number seem weird to me, but that's
a different problem). The second column gives us the process ID and the third column shows us the external program 
along with the arguments it was executed it. We can see that the script `/tmp/git_suggestions` is executed and the
command for which the auto-completion suggestions are being shown is provided as the first argument.


Now, go back to Terminal 1, and type:

```
$ git chec<TAB>
```

On Terminal 2, we will see:

```
189064     16701/tmp/git_suggestions git chec git
```

We see that the script `/tmp/git_suggestions` is now being called with three arguments:

- `git`: The command we are trying to ask BASH to suggest auto-completions for
- `chec`: The word we are asking for completions for
- `git`: The word before the word we are asking for completions for

Let's discuss this a bit. When we press `TAB`, bash tries to find out a matching auto-completion "handler" for the command 
and if it finds one, invokes the handler. When invoking the handler, it calls the handler providing data about the command
it is attempting to provide auto-completion suggestions for. The output of the handler is then parsed by bash and each separate 
line in the output is suggested as possible candidates for the auto-completion.

Let's discuss more about the data provided to completion handlers.

## Data provided to completion handlers

We will update our `/tmp/git_suggestions` script to be as follows:

```
#!/bin/bash

echo "status"
echo "checkout"
echo "clone"
echo "branch"

# Print the envrionment variables relevant to bash autocompletion
env | grep COMP
```

We add a line at the end to print all environment variables having COMP in them. Now, if we go back to Terminal 1,
and type `$git checkout <TAB><TAB>`, we will see:

```

$ git checkout <TAB>
branch                   clone                    COMP_LINE=git checkout   COMP_TYPE=63             
checkout                 COMP_KEY=9               COMP_POINT=13            status 

```

Besides the earlier four suggestions, we have some new ones:

- `COMP_LINE`: This is the entire line of the command we pressed `<TAB><TAB>` on
- `COMP_TYPE`: This is the type of completion that is being done, `63` is the ASCII code for `?`. According to the Bash programmable completion manual,
this is the operation which will list completions after successive tabs
- `COMP_KEY`: This is the ASCII code of the key which triggered the auto-completion. 9 [stands for](https://en.wikipedia.org/wiki/Tab_key#Tab_characters)
the TAB key
- `COMP_POINT`: This is the cursor position where the `<TAB><TAB>` was pressed

https://www.gnu.org/software/bash/manual/html_node/Bash-Variables.html#Bash-Variables

([show-key](https://linux.die.net/man/1/showkey)) is pretty cool.

## Single `<TAB>` and double `<TAB><TAB>`

It's only while working on this article, I realized, that a double `<TAB><TAB>` was necessary to display all the possible
auto-completions in most cases. Let's look into that. A single `<TAB>` by default 


## `complete` built-in command

## `compgen` built-in command


## Getting good old BASH completion back

Now, on Terminal 1, let's install the `bash-completion` package and then exec bash and type, in `$ git <TAB>`


On Terminal 2, now we will see something like:

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

## Putting the pieces together







https://www.joshmcguigan.com/blog/shell-completions-pure-rust/

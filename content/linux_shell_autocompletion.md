What happens when you press `<your favorite Linux command> <TAB>`? You get a bunch of suggestions with one of them
just the one you had in mind you were going to type, or did not. Auto completions in shells are super helpful and that's
probably the most boring sales pitch. Anyway, as you will see, quite a bit goes on behind the scenes. In this post, we will
find out about the minions that gets to work when we press the `<TAB>` key on a popular Linux shell, `BASH`.

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

## Look, no bash completion

On a terminal, (Terminal 1), type in `$ git <TAB>`. The only suggestions you get will be the files in the
directory you are in. Not very helpful suggestions, you say. I know - and that's because we uninstalled a 
RPM package which would have magically done that for us. (`bash-completion`). Let's keep it that way for us now.


External program invocations for auto complete suggestions

Terminal 1:

$ sudo bpftrace ./execnsoop.bt



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






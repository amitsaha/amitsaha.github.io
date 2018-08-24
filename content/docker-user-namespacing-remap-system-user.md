Title: Using Docker `userns-remap` with a system user on Linux
Date: 2018-08-24 20:00
Category: infrastructure
Status: Draft

In this post, we learn how we can make use of `docker`'s user namespacing feature on Linux in a CI/build environment
to avoid running into permission issues while also keeping things a bit sane from the security standpoint.

# Introduction

Let's consider that we are leveraging `docker` in a continuous integration (CI)/build environment and the scenario
looks as follows:

1. CI agent/slave runs as an unpriviliged user `agent` on the host
2. `agent` clones the repository during a build on the host
3. The build happens in a `docker` container spawned by scripts running as `agent` with the repository volume mounted

On a new build, the agent doesn't do a fresh clone, but instead does a `git clean` followed by `git fetch` of the commit.

Here's is what's going to happen: the `agent` is going to get a permission denied. In Step 3 above, when the build was
done in the container, the build process was running as `root` user. Since the repository was volume mounted, contents
written to the repository directory will show up as being owned by the `root` user on the host. Hence, when `agent` tries
to cleanup the directory on the next build, it gets a permission denied.

What do we do? We could run the CI agent as `root` user - avoid it. Or, figure out some way of changing back the permissions
after the build. However, `User namespaces` via `userns-remap` is better than both these workarounds.

Before we get into configuring `docker` engine, we have a bit to learn about Linux `system` users and entries in
`/etc/subuid` and `/etc/subgid`.

## System users and entries in `/etc/subuid` and `/etc/subgid`





```
echo "--- Updating subuid and subguid maps to include buildkite-agent user"
username="buildkite-agent"
uid=$(id -u "$username")
gid=$(id -g "$username")
lastuid=$(( uid + 65536 ))
lastgid=$(( gid + 65536 ))

sudo usermod --add-subuids "$uid"-"$lastuid" "$username"
sudo usermod --add-subgids "$gid"-"$lastgid" "$username"

sudo cat /etc/subuid
sudo cat /etc/subgid
```

## Enabling `docker's` userns-remap

sudo mkdir -p /etc/systemd/system/docker.service.d
echo "[Service]" | sudo tee /etc/systemd/system/docker.service.d/docker-userns-remap.conf > /dev/null
# First clear ExecStart (https://github.com/moby/moby/issues/14491)
echo "ExecStart=" | sudo tee --append  /etc/systemd/system/docker.service.d/docker-userns-remap.conf > /dev/null
# Now, override to apply userns-remap
echo "ExecStart=/usr/bin/dockerd -H fd:// --userns-remap=\"buildkite-agent:buildkite-agent\"" | sudo tee --append  /etc/systemd/system/docker.service.d/docker-userns-remap.conf > /dev/null

sudo systemctl daemon-reload
sudo systemctl restart docker


## Using third party images

root@735b8dd78674:~#  find / \( -uid 1004 \)  -ls 2>/dev/null | head -1
   772345      4 drwxr-xr-x   1 1004     sudo         4096 Aug 15 02:29 /usr/share/dotnet

https://circleci.com/docs/2.0/high-uid-error/




https://github.com/moby/moby/pull/21266/commits/c18e7f3a0419e35aeab4eefa51f3c17fbd72381f


https://success.docker.com/article/introduction-to-user-namespaces-in-docker-engine


```

Title: Docker userns-remap and Container ID cannot be mapped to Host ID
Date: 2018-08-14 20:00
Category: infrastructure
Status: Draft


root@735b8dd78674:~#  find / \( -uid 1004 \)  -ls 2>/dev/null | head -1
   772345      4 drwxr-xr-x   1 1004     sudo         4096 Aug 15 02:29 /usr/share/dotnet

https://circleci.com/docs/2.0/high-uid-error/


https://docs.docker.com/engine/security/userns-remap/

https://github.com/moby/moby/pull/21266/commits/c18e7f3a0419e35aeab4eefa51f3c17fbd72381f


https://success.docker.com/article/introduction-to-user-namespaces-in-docker-engine

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

echo "--- Updating docker systemd config file to use userns remap"

sudo mkdir -p /etc/systemd/system/docker.service.d
echo "[Service]" | sudo tee /etc/systemd/system/docker.service.d/docker-userns-remap.conf > /dev/null
# First clear ExecStart (https://github.com/moby/moby/issues/14491)
echo "ExecStart=" | sudo tee --append  /etc/systemd/system/docker.service.d/docker-userns-remap.conf > /dev/null
# Now, override to apply userns-remap
echo "ExecStart=/usr/bin/dockerd -H fd:// --userns-remap=\"buildkite-agent:buildkite-agent\"" | sudo tee --append  /etc/systemd/system/docker.service.d/docker-userns-remap.conf > /dev/null

sudo systemctl daemon-reload
sudo systemctl restart docker
```

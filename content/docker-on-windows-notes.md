Title: Notes on running Windows docker containers on Windows
Date: 2017-10-26 15:00
Category: infrastructure
Status: draft

I have been working with Windows docker containers running on Windows for the past three months with
the goal to have isolated environments for each build in a continuous integration pipeline. 
That is, each build happens on an exclusive build host (AWS EC2 VM instance) and every database and service the application 
needs access to for the integration tests (including selenium tets) are run on docker containers on the 
same host.

All docker features I was familiar with on Linux and needed access on Windows to just worked. The experience was 
definitely 100x better (faster and reliable) on Windows Server than on Windows 10 (more on this soon). But, 
considering that this was for a CI environment, it was a good thing. I wish I had moved to Windows Server earlier 
for my experimentation.

Next, I share some of my findings in the hope that it may be useful to others.

## Docker versions

On Windows 10, when we install docker from the [docker store](https://www.docker.com/docker-windows), we get the community
version of docker. Once you install it, and run `docker version`, you will get the version stated as `18.03.1-ce` or similar.

On Windows server, we have to use the docker enterprise version which can be installed as per [this page](https://docs.docker.com/install/windows/docker-ee/). When we run `docker version` on a Windows server host, you will 
see the version stated as `18.03.1-ee-2` or simiar.

The following `PowerShell` snippet will install docker engine and `docker-compose` on Windows server:

```
# docker engine and docker compose
Install-PackageProvider -Name NuGet -MinimumVersion 2.8.5.201 -Force
Install-Module -Name DockerMsftProvider -Repository PSGallery -Force
Install-Package -Name docker -ProviderName DockerMsftProvider -Force
choco install docker-compose # this needs chocolatey installed
```

## Container Isolation

On Windows 10, `docker` uses `hyperv` isolation:

```
> docker info -f '{{ .Isolation }}'
hyperv
```

On Windows server, it uses `process` isolation:

```
> docker info -f '{{ .Isolation }}'
process
```

(Please note, use `powershell` for the above command, not `cmd` [issue](https://github.com/moby/moby/issues/33959))

This basically means that on Windows 10, a container is running within a tiny VM. I suspect this is a major reason for 
the slow container startup and image build times on Windows 10. To learn more, see [here](https://docs.microsoft.com/en-us/virtualization/windowscontainers/manage-containers/hyperv-container) and a relevant [issue](https://github.com/docker/for-win/issues/1822).

## Docker commit

`docker commit` allows us to create a image from a running container. However, on Windows we cannot do that:

```
> docker commit hopeful_clarke myimage
Error response from daemon: windows does not support commit of a running container
```

We will have to stop the container and then commit it. 

```
$ docker stop hopeful_clarke
$ docker commit hopeful_clarke myimage
```

## User-defined networks

On Windows 10, multiple `nat` networks are supported, but on Windows Server with 17.06 EE docker engine, only 
one `nat` network is supported. Hence, when using `docker-compose`, we must specify the following:

```
networks:
  default:
    external:
      name: nat
      
```

The reason for the above is the default behavior of `docker-compose` is to create a new network for the services which will
fail with an error: `Problem : Error response from daemon: HNS failed with error : The parameter is incorrect`.

See [here](https://docs.microsoft.com/en-us/virtualization/windowscontainers/container-networking/network-drivers-topologies)
to learn more.

## `hcsshim::PrepareLayer failed in Win32`

While building images and running containers, an error as follows appeared with varying frequency:

```
hcsshim::PrepareLayer failed in Win32: This operation returned because the timeout period expired. (0x5b4) layerId=ae8414f34fb31faea64b7bee078b295023db93c8505c67da686750843c853629 flavour=1
```

On Windows 10, it was very frequent. On Windows server, the frequency is almost never. A situation in which it 
can happen as discussed in this [issue](https://github.com/moby/moby/issues/27588) is when
an attempt is made to run a number of containers simulataneously. In my case, when I got the issue on Windows server,
the number of containers was around 10. However, the most common occurence for me was on Windows 10 while building
images.

A retry of the operation usually fixed it. Relevant project - [hccshim](https://github.com/Microsoft/hcsshim).

## Writing Dockerfiles for Windows containers

One of the first things I had to tackle while writing a `Dockerfile` was that of which slash to use in Path names.
Basically, for Dockerfile instructions such as `ADD`, `WORKDIR`, `COPY`, continue using the "Linux" style, "/".
For example:

```
ADD . C:/app
WORKDIR C:/app
```

For `RUN` instruction, you must specify any path using ":\", like so:

```
RUN Expand-Archive -Path c:\mysql.zip -DestinationPath C:\
..
```

When writing multi-line instructions in a `Dockerfile`, you may be tempted to use "\", but, that can cause problems, since
"\" is also a path separator in Windows, so declare a `escape` character as explained 
[here](https://docs.microsoft.com/en-us/virtualization/windowscontainers/manage-docker/manage-windows-dockerfile#escape-character).

However, I seemed to have success without it - the following snippet works just fine to build a mysql server Windows 
container:

```
FROM microsoft/nanoserver
..
RUN powershell -command \
    Expand-Archive -Path c:\mysql.zip -DestinationPath C:\ ; \
    ren C:\mysql-5.7.22-winx64 C:\MySQL ; \
    New-Item -Path C:\MySQL\data -ItemType directory ; \
    C:\MySQL\bin\mysqld.exe --initialize --init-file=C:\mysql-init.txt --console --explicit_defaults_for_timestamp ; \
    Remove-Item c:\mysql.zip -Force
..
```

## Volume mounting in NodeJS applications

Trying to do a `npm install` (for example) on a NojdeJS application inside a container with the directory
volume mounted from host wouldn't work. The issue is described [here](https://github.com/nodejs/node/issues/8897).

My solution was to basically:

- Volume mount host directory: `docker run -v ${pwd}:C:\app --name venus-legacy-front-end-builder ..`
- Inside the container, copy the files to another directory - i.e. copy from `C:\app` to `C:\workspace`
- Build within `C:\workspace`
- Copy built assets back from `C:\workspace` to `C:\app`



## Dockerizing a dotnet framework application

```
FROM microsoft/dotnet-framework:4.7.2-sdk
RUN Add-WindowsFeature Web-Server; `
    Add-WindowsFeature NET-Framework-45-ASPNET; `
    Add-WindowsFeature Web-Asp-Net45; `
    Remove-Item -Recurse C:\inetpub\wwwroot\*; `
    Invoke-WebRequest -Uri https://dotnetbinaries.blob.core.windows.net/servicemonitor/2.0.1.3/ServiceMonitor.exe -OutFile C:\ServiceMonitor.exe
# Install the rewrite module
# Open issue: https://github.com/Microsoft/aspnet-docker/issues/109
WORKDIR /install
ADD https://download.microsoft.com/download/C/9/E/C9E8180D-4E51-40A6-A9BF-776990D8BCA9/rewrite_amd64.msi rewrite_amd64.msi
RUN Write-Host 'Installing URL Rewrite' ; `
    Start-Process msiexec.exe -ArgumentList '/i', 'rewrite_amd64.msi', '/quiet', '/norestart' -NoNewWindow -Wait

# https://www.richard-banks.org/2017/02/debug-net-in-windows-container.html
# https://docs.microsoft.com/en-us/visualstudio/debugger/remote-debugger-port-assignments
EXPOSE 4022 4023
ADD https://download.visualstudio.microsoft.com/download/pr/12210758/3732c1fb2e37696edab25c565695c1b0/VS_RemoteTools.exe rtools_setup.exe
RUN & './rtools_setup.exe' /install /quiet

WORKDIR C:\
RUN nuget.exe install WebConfigTransformRunner -Version 1.0.0.1

VOLUME C:\app
EXPOSE 51033 51034 51035
WORKDIR /app
HEALTHCHECK --interval=30s --timeout=10s --retries=20 CMD powershell -Command Invoke-WebRequest -UseBasicParsing http://localhost:51033;Invoke-WebRequest -UseBasicParsing http://localhost:51034;Invoke-WebRequest -UseBasicParsing http://localhost:51035;
```


```
$ErrorActionPreference = "Stop"

# Setup Venus websites

Write-Output "--- Setting up IIS websites"

cp -r RSAdmin C:\inetpub\wwwroot\
cp -r RSMembers C:\inetpub\wwwroot\
cp -r RSApi C:\inetpub\wwwroot\


New-IISSite -Name 'site1' -PhysicalPath c:\inetpub\wwwroot\site1 -BindingInformation "*:51033:"
New-IISSite -Name 'site2' -PhysicalPath c:\inetpub\wwwroot\site2 -BindingInformation "*:51034:"

Set-Location C:\inetpub\wwwroot\site1
C:\WebConfigTransformRunner.1.0.0.1\Tools\WebConfigTransformRunner.exe .\Configs\AppSetting.config .\Configs\AppSetting.Docker.config .\Configs\AppSetting.config
C:\app\UpdateConfigsDocker.ps1

cat .\Configs\AppSetting.config

Set-Location C:\inetpub\wwwroot\site2
C:\WebConfigTransformRunner.1.0.0.1\Tools\WebConfigTransformRunner.exe .\Configs\AppSetting.config .\Configs\AppSetting.Docker.config .\Configs\AppSetting.config
C:\app\UpdateConfigsDocker.ps1

..
```

https://docs.microsoft.com/en-us/virtualization/windowscontainers/manage-docker/manage-windows-dockerfile
https://blog.alexellis.io/run-iis-asp-net-on-windows-10-with-docker/
https://anthonychu.ca/post/aspnet-web-config-transforms-windows-containers-revisited/
https://anthonychu.ca/post/overriding-web-config-settings-environment-variables-containerized-aspnet-apps/

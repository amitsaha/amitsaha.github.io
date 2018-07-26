Title: Notes on running Windows docker containers on Windows
Date: 2017-10-26 15:00
Category: infrastructure
Status: draft

## Docker versions

## Isolation

## Docker commit

```
$ docker stop containerid
$ docker commit containerid new_image
```

    

## User-defined networks

On Windows 10, multiple `nat` networks are supported, but on Windows Server with 17.06 EE docker engine, only one `nat` network
is supported. Hence, when using `docker-compose`, we must specify the following:

```
networks:
  default:
    external:
      name: nat
      
```

The reason for the above is the default behavior of `docker-compose` is to create a new network for the services which will
fail with an error: `Problem : Error response from daemon: HNS failed with error : The parameter is incorrect`.
https://docs.microsoft.com/en-us/virtualization/windowscontainers/container-networking/network-drivers-topologies

## Miscellaneous

```
hcsshim::PrepareLayer failed in Win32: This operation returned because the timeout period expired. (0x5b4) layerId=ae8414f34fb31faea64b7bee078b295023db93c8505c67da686750843c853629 flavour=1
```

https://github.com/moby/moby/issues/27588

## Dockerfiles


```
FROM microsoft/windowsservercore
SHELL ["powershell", "-Command", "$ErrorActionPreference = 'Stop'; $ProgressPreference = 'SilentlyContinue';"]
RUN Set-ExecutionPolicy Bypass -Scope Process -Force; iex ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))

RUN choco install -y googlechrome phantomjs git nodejs yarn
CMD [ "node.exe" ]
```

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

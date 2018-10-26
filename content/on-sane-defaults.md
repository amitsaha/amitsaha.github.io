Title: On sane defaults
Date: 2018-07-26 16:00
Category: infrastructure
Status: Draft

My task at hand was simple. Build a Docker image of a ASP.NET application (full framework) hosted in IIS on
a build host (**host1**) and move it to a deployment host (**host2**) and run it. This is a story of how I spent close to two full working days
trying to debug a simple issue which sane default behavior of a tool would have cut it to seconds.

## Key details

The key details that are important to my story are:

- **host1** and **host2** lives in two different AWS VPC subnets
- **host2** has access to some external services that the web application communicates with when the 
  homepage is hit


## Observations on build host

I built the image on build host, and ran it in a docker container, like so:

```
$ docker run -d test/image
```

My web application is configured to run on port 51034. From the host, I find out the container IP using `docker inspect`
and make a GET request using PowerShell's `Invoke-WebRequest`:

```
$ Invoke-WebRequest -UseBasicParsing http://ip:51034
```

I get back errors saying there was a error in connecting to the external services. This is expected, since
**host1** doesn't have connectivity to these services.

Great, so I push the docker image to AWS ECR.

## Observations on deployment host


I spent around 16 hours trying to debug a issue 
> curl 172.29.170.207:51034
curl : Unable to connect to the remote server
At line:1 char:1
+ curl 172.29.170.207:51034
+ ~~~~~~~~~~~~~~~~~~~~~~~~~
    + CategoryInfo          : InvalidOperation: (System.Net.HttpWebRequest:HttpWebRequest) [Invoke-WebRequest], WebException
    + FullyQualifiedErrorId : WebCmdletWebResponseException,Microsoft.PowerShell.Commands.InvokeWebRequestCommand

PS C:\Users\Administrator\work\curl> docker ps
CONTAINER ID        IMAGE                                                                          COMMAND                   CREATED              STATUS                         PORTS               NAMES
1fc75d1cb667        519527725796.dkr.ecr.ap-southeast-2.amazonaws.com/rsaus/venus-members:latest   "powershell .\\Star..."   About a minute ago   Up About a minute (healthy)    51034/tcp           adoring_dijkstra
15bb9cff5dd7        fff28fca3ca8                                                                   "powershell .\\Star..."   About an hour ago    Up About an hour (unhealthy)   51034/tcp           musing_sinoussi
PS C:\Users\Administrator\work\curl> .\curl-7.61.1-win64-mingw\bin\curl.exe  172.29.170.207:51034
<html><head><title>Object moved</title></head><body>
<h2>Object moved to <a href="https://localhost:51034/">here</a>.</h2>
</body></html>
PS C:\Users\Administrator\work\curl> curl 172.29.170.207:51034 -Verbose
VERBOSE: GET http://172.29.170.207:51034/ with 0-byte payload
curl : Unable to connect to the remote server
At line:1 char:1
+ curl 172.29.170.207:51034 -Verbose
+ ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
    + CategoryInfo          : InvalidOperation: (System.Net.HttpWebRequest:HttpWebRequest) [Invoke-WebRequest], WebException
    + FullyQualifiedErrorId : WebCmdletWebResponseException,Microsoft.PowerShell.Commands.InvokeWebRequestCommand

PS C:\Users\Administrator\work\curl> curl 172.29.170.207:51034 -Verbose -MaximumRedirection 0
VERBOSE: GET http://172.29.170.207:51034/ with 0-byte payload
VERBOSE: received 141-byte response of content type text/html; charset=utf-8
curl : The response content cannot be parsed because the Internet Explorer engine is not available, or Internet Explorer's first-launch configuration is not complete. Specify the UseBasicParsing parameter and try again.
At line:1 char:1
+ curl 172.29.170.207:51034 -Verbose -MaximumRedirection 0
+ ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
    + CategoryInfo          : NotImplemented: (:) [Invoke-WebRequest], NotSupportedException
    + FullyQualifiedErrorId : WebCmdletIEDomNotSupportedException,Microsoft.PowerShell.Commands.InvokeWebRequestCommand

PS C:\Users\Administrator\work\curl> curl -UseBasicParsing 172.29.170.207:51034 -Verbose -MaximumRedirection 0
VERBOSE: GET http://172.29.170.207:51034/ with 0-byte payload
VERBOSE: received 141-byte response of content type text/html; charset=utf-8


StatusCode        : 301
StatusDescription : Moved Permanently
Content           : <html><head><title>Object moved</title></head><body>
                    <h2>Object moved to <a href="https://localhost:51034/">here</a>.</h2>
                    </body></html>

RawContent        : HTTP/1.1 301 Moved Permanently
                    X-Content-Type-Options: nosniff
                    X-UA-Compatible: IE=Edge,chrome=1
                    Content-Length: 141
                    Content-Type: text/html; charset=utf-8
                    Date: Fri, 26 Oct 2018 02:38:00 GMT
                    Lo...
Forms             :
Headers           : {[X-Content-Type-Options, nosniff], [X-UA-Compatible, IE=Edge,chrome=1], [Content-Length, 141], [Content-Type, text/html; charset=utf-8]...}
Images            : {}
InputFields       : {}
Links             : {@{outerHTML=<a href="https://localhost:51034/">here</a>; tagName=A; href=https://localhost:51034/}}
ParsedHtml        :
RawContentLength  : 141

curl : The maximum redirection count has been exceeded. To increase the number of redirections allowed, supply a higher value to the -MaximumRedirection parameter.
At line:1 char:1
+ curl -UseBasicParsing 172.29.170.207:51034 -Verbose -MaximumRedirecti ...
+ ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
    + CategoryInfo          : InvalidOperation: (System.Net.HttpWebRequest:HttpWebRequest) [Invoke-WebRequest], InvalidOperationException
    + FullyQualifiedErrorId : MaximumRedirectExceeded,Microsoft.PowerShell.Commands.InvokeWebRequestCommand

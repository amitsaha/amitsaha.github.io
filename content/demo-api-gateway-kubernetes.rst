:Title: Getting started with deploying services in a local Kubernetes cluster
:Date: 2017-05-06 19:00
:Category: kubernetes

I wrote up a demo `API Gateway <https://github.com/amitsaha/apigatewaydemo>`__ a while back which demonstrated how we can
write an API gateway which would forward incoming requests to downstream services using HTTP and a gRPC service as
examples. In this post, I attempt to deploy the services in a Kubernetes cluster as a way to learn about Kubernetes.

To follow along, we will use the `kubernetes branch <https://github.com/amitsaha/apigatewaydemo/tree/kubernetes>`__ of
the code. One notable change to the code for the API gateway itself is removing the use of 
`go kit <https://github.com/go-kit/kit>`__, since at this point it does seem like, I should be able to get the features i was looking for 
from ``kubernetes`` and other software I plan to try (such as `linkerd <https://linkerd.io/>`__ for service proxying, circuit breaking, distributed tracing, etc). But, I will see how that goes. If I succeed, it means my API gateway code will be simpler.

Le'ts get started.

Setting up Kubernetes
=====================

My host configuration is Mac OS X and I am using VirtualBox. Install `minkube <https://github.com/kubernetes/minikube>`__ and  `kubectl <https://coreos.com/kubernetes/docs/latest/configure-kubectl.html>`__ and we will do the following:

.. code::

   $ minikube start
   Starting local Kubernetes cluster...
   Starting VM...
   Downloading Minikube ISO
   89.51 MB / 89.51 MB [==============================================] 100.00% 0s
   SSH-ing files into VM...
   Setting up certs...
   Starting cluster components...
   Connecting to cluster...
   Setting up kubeconfig...
   Kubectl is now configured to use the cluster.
   
Let's perfomr some sanity checking:

.. code::

  $ kubectl version
  Client Version: version.Info{Major:"1", Minor:"6", GitVersion:"v1.6.1", GitCommit:"b0b7a323cc5a4a2019b2e9520c21c7830b7f708e", GitTreeState:"clean", BuildDate:"2017-04-03T20:44:38Z", GoVersion:"go1.7.5", Compiler:"gc", Platform:"darwin/amd64"}
  Server Version: version.Info{Major:"1", Minor:"6", GitVersion:"v1.6.0", GitCommit:"fff5156092b56e6bd60fff75aad4dc9de6b6ef37", GitTreeState:"dirty", BuildDate:"2017-04-07T20:46:46Z", GoVersion:"go1.7.3", Compiler:"gc", Platform:"linux/amd64"}
  
  $ kubectl cluster-info
  Kubernetes master is running at https://192.168.99.100:8443
  KubeDNS is running at https://192.168.99.100:8443/api/v1/proxy/namespaces/kube-system/services/kube-dns
  kubernetes-dashboard is running at https://192.168.99.100:8443/api/v1/proxy/namespaces/kube-system/services/kubernetes-dashboard
  
  $ kubectl get services
  NAME         CLUSTER-IP   EXTERNAL-IP   PORT(S)   AGE
  kubernetes   10.0.0.1     <none>        443/TCP   14m
  
  $ kubectl get nodes
  NAME       STATUS    AGE       VERSION
  minikube   Ready     14m       v1.6.0
  
  $ kubectl get pods
  No resources found.

The last output is interesting. Applications/services live inside a pod in Kubernetes and we currently don't have any running,
hence no pods are shown. Similarly, `kubectl get services` runs only the `kubernetes` service running on port 443. I believe this is the kubernetes API server.


Before we delve into deploying all the services, I will discuss three concepts which we will apply this post.

Kubernetes concepts
===================

**Pods**

A `pod <https://kubernetes.io/docs/concepts/workloads/pods/pod/>`__ is a logical group of one or more containers, together providing the functionality of a single application. In our case, each of the API gateway and the services will be a pod each with a single docker container. We won't be creating a pod directly, but via a *deployment*.

**Deployment**

A `deployment <https://kubernetes.io/docs/concepts/workloads/controllers/deployment/>`__ is a declaration of the desired
state of the *pods*. Here's where we can specify how many of the *pods* we want running, what should be running in those
pods (container image, name), ports to expose, healthcheck and others.

**Service**

A `service <https://kubernetes.io/docs/concepts/services-networking/service/>`__ is an abstraction to access the application
running in pods. To access the application running inside a pod, we access it via the service which automatically load balances the requests among the pods. Another helpful document I found was `this one <https://kubernetes.io/docs/concepts/services-networking/connect-applications-service/>`__.


Deploying the services
======================

At this stage, we have a Kubernetes cluster up and running. Now a bit about what we are going to deploy in it? We are going to deploy an "API gateway" and two other services. The API gateway forwards requests it gets to one of these services - `webapp-1` via HTTP 1.1 and `grpc-app-1` via `gRPC <http://www.grpc.io/>`__. 

What are the features we want to have?

- We will be running 3 instances of the API gateway service and 3 instances each of the other services
- API gateway should not have any hard coded IP address for any service it talks to. If an instance of a service goes up or down, the API gateway shouldn't have to know about it
- API gateway should load balance between the services
- External requests should be load balanced between the multiple API gateway instances
- We should have circuit breaking between API gateway and any of the other services
- We should have rate limitting across the API gateway
- We should have metrics on each service
- We should have correlated logging
- We should have distributed tracing
- We should be able to scale up or down automatically based on incoming requests
- Upgrades and downgrades without downtime
- Only healthy instances should get traffic
- Monitoring our service system usage

**Setting up docker**

We will be using the `docker` engine running in the minikube VM so that we can build the images in the VM and deploy them without having to push to a remote registry. In the shell on the host system, where we will build docker images, doing the following will suffice:

.. code::

    $ eval $(minikube docker-env)


Service #1: Deploying the HTTP service
======================================

First, let's build the image for the `webapp-1` service:

.. code::

    $ cd webapp-1
    $ docker build -t amitsaha/webapp-1 .
    
Next, we will create a kubernetes `deployment`:

.. code::

    $ cat kubernetes/deployment.yml

      apiVersion: apps/v1beta1
      kind: Deployment
      metadata:
        name: webapp-1-deployment
      spec:
        replicas: 3
        template:
          metadata:
            labels:
              app: webapp-1
          spec:
            containers:
            - name: webapp-1
              image: amitsaha/webapp1
              ports:
              - containerPort: 5000
        
To create the deployment:

.. code::
    
    $ kubectl create -f deployment.yaml
    deployment "webapp-1-deployment" created
    

.. code::
    
   $ kubectl describe deployment webapp-1-deployment

   Name:                   webapp-1-deployment
   Namespace:              default
   CreationTimestamp:      Wed, 03 May 2017 13:46:46 +1000
   Labels:                 app=webapp-1
   Annotations:            deployment.kubernetes.io/revision=1
   Selector:               app=webapp-1
   Replicas:               3 desired | 3 updated | 3 total | 1 available | 2 unavailable
   StrategyType:           RollingUpdate
   MinReadySeconds:        0
   RollingUpdateStrategy:  25% max unavailable, 25% max surge
   Pod Template:
     Labels:       app=webapp-1
     Containers:
      webapp-1:
       Image:              amitsaha/webapp1:latest
       Port:               5000/TCP
       Liveness:           http-get http://:80/_status/healthcheck/ delay=30s timeout=1s period=10s #success=1 #failure=3
       Environment:        <none>
       Mounts:             <none>
     Volumes:              <none>
   Conditions:
     Type          Status  Reason
     ----          ------  ------
     Progressing   True    NewReplicaSetAvailable
     Available     False   MinimumReplicasUnavailable
   OldReplicaSets: <none>
   NewReplicaSet:  webapp-1-deployment-4250575981 (3/3 replicas created)
   Events:         <none>


.. code::

   $ kubectl get pods -l app=webapp-1
   NAME                                 READY     STATUS    RESTARTS   AGE
   webapp1-deployment-536678510-dtmjb   1/1       Running   0          4m
   webapp1-deployment-536678510-kt1zs   1/1       Running   0          4m
   webapp1-deployment-536678510-wkmkq   1/1       Running   0          4m


.. code::
    $ cat kubernetes/service.yml
      apiVersion: v1
      kind: Service
      metadata:
        name: webapp-1
      spec:
        selector:
          app: webapp-1
        ports:
          - protocol: TCP
            port: 80
            targetPort: 5000


.. code::

    $ kubectl create -f kubernetes/service.yaml
    service "webapp-1" created
      
.. code::

      $ kubectl describe svc webapp1
      Name:			webapp1
      Namespace:		default
      Labels:			<none>
      Annotations:		<none>
      Selector:		app=webapp1
      Type:			ClusterIP
      IP:			10.0.0.91
      Port:			<unset>	80/TCP
      Endpoints:		172.17.0.5:5000,172.17.0.8:5000,172.17.0.9:5000
      Session Affinity:	None
      Events:			<none>


Interacting with the service:

.. code::

   $ minikube ssh
   $ curl 10.0.0.91/create
   <!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 3.2 Final//EN">
   <title>405 Method Not Allowed</title>
   <h1>Method Not Allowed</h1>
   <p>The method is not allowed for the requested URL.</p>

We will also be able to talk to our webapp1 service using "webapp-1" from another *pod*.

Service #2: Deploying the RPC service
=====================================

$ cd apigatewaydemo/grpc-app-1/server
$ docker build -t amitsaha/rpc-app-1 .
..

$ cat kubernetes/deployment.yaml

apiVersion: apps/v1beta1
kind: Deployment
metadata:
  name: rpc-app-1-deployment
spec:
  replicas: 3
  template:
    metadata:
      labels:
        app: rpc-app-1
    spec:
      containers:
      - name: rpc-app-1
        image: amitsaha/rpc-app-1:latest
        imagePullPolicy: Never
        ports:
        - containerPort: 6000
        livenessProbe:
          tcpSocket:
            port: 6000
          initialDelaySeconds: 30
          timeoutSeconds: 1
$ kubectl create -f kubernetes/deployment.yaml
deployment "rpc-app-1-deployment" created


$ cat kubernetes/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: rpc-app-1
spec:
  selector:
    app: rpc-app-1
  ports:
    - protocol: TCP
      port: 6000
      targetPort: 6000

$ kubectl create -f kubernetes/service.yaml
service "rpc-app-1" created



$ kubectl get services
NAME         CLUSTER-IP   EXTERNAL-IP   PORT(S)    AGE
kubernetes   10.0.0.1     <none>        443/TCP    8d
rpc-app-1    10.0.0.30    <none>        6000/TCP   17s
webapp-1     10.0.0.46    <none>        80/TCP     8d
webapp1      10.0.0.91    <none>        80/TCP     8d

API gateway: Deploying the API gateway
=====================================

$ cd apigatewaydemo/apigateway
$ docker build -t amitsaha/apigateway .
..


$ kubectl create -f kubernetes/deployment.yaml
deployment "apigateway" created

$ kubectl create -f kubernetes/service.yaml
service "apigateway" created

..


➜  apigateway git:(kubernetes) ✗ kubectl get pod | grep 'apigateway' | cut -d " " -f1 - | xargs -n1 -P 10 kubectl delete pod
p


$ kubectl get services
NAME         CLUSTER-IP   EXTERNAL-IP   PORT(S)    AGE
apigateway   10.0.0.153   <none>        80/TCP     21h
kubernetes   10.0.0.1     <none>        443/TCP    23d
rpc-app-1    10.0.0.30    <none>        6000/TCP   14d
webapp-1     10.0.0.46    <none>        80/TCP     22d

$ curl -q -H "Content-type: application/json" -X POST -d '{"title1":"My project hello hello11"}' 10.0.0.153/projects/
{
  "id": 123,
  "url": "Project-None"
}


diff --git a/apigateway/kubernetes/service.yaml b/apigateway/kubernetes/service.yaml
index 8c32a97..819ae25 100644
--- a/apigateway/kubernetes/service.yaml
+++ b/apigateway/kubernetes/service.yaml
@@ -9,3 +9,4 @@ spec:
     - protocol: TCP
       port: 80
       targetPort: 8000
+  type: NodePort
(END)


➜  apigateway git:(kubernetes) ✗ kubectl describe services apigateway
Name:                   apigateway
Namespace:              default
Labels:                 <none>
Annotations:            <none>
Selector:               app=apigateway
Type:                   ClusterIP
IP:                     10.0.0.153
Port:                   <unset> 80/TCP
Endpoints:              172.17.0.11:8000,172.17.0.14:8000,172.17.0.15:8000
Session Affinity:       None
Events:                 <none>
➜  apigateway git:(kubernetes) ✗ kubectl apply -f kubernetes/service.yaml
Warning: kubectl apply should be used on resource created by either kubectl create --save-config or kubectl apply
service "apigateway" configured
➜  apigateway git:(kubernetes) ✗ kubectl describe services apigateway
Name:                   apigateway
Namespace:              default
Labels:                 <none>
Annotations:            kubectl.kubernetes.io/last-applied-configuration={"apiVersion":"v1","kind":"Service","metadata":{"annotations":{},"name":"apigateway","namespace":"default"},"spec":{"ports":[{"port":80,"protocol":"TCP...
Selector:               app=apigateway
Type:                   NodePort
IP:                     10.0.0.153
Port:                   <unset> 80/TCP
NodePort:               <unset> 30638/TCP
Endpoints:              172.17.0.11:8000,172.17.0.14:8000,172.17.0.15:8000
Session Affinity:       None
Events:                 <none>
➜

➜  apigateway git:(kubernetes) ✗ minikube service --url apigateway
http://192.168.99.100:30638
➜


$ curl -q -H "Content-type: application/json" -X POST -d '{"title1":"My project hello hello11"}' 10.0.0.153/projects/
{
  "id": 123,
  "url": "Project-None"
}
$
$ curl -q -H "Content-type: application/json" -X POST -d '{"title1":"My project hello hello11"}' 10.0.0.153/verify/
{"message":"Verified: 12321"}
$
$ curl -q -H "Content-type: application/json" -X POST -d '{"id": 111, "token": "a$$" }'
10.0.0.153/verify/
{"message":"Verified: 111"}
$ curl -q -H "Content-type: application/json" -X POST -d '{"id": 121, "token": "a$$" }'
10.0.0.153/verify/
{"message":"Verified: 121"}
$
$ curl -q -H "Content-type: application/json" -X POST -d '{"id1": 121, "token": "a$$" }'
 10.0.0.153/verify/
{"message":"Verified: 0"}
$
$


➜  apigateway git:(kubernetes) ✗ curl -q -H "Content-type: application/json" -X POST -d '{"id": 121, "token": "a$$" }' `minikube service --url apigateway`/verify/
{"message":"Verified: 121"}
➜  apigateway git:(kubernetes) ✗
➜  apigateway git:(kubernetes) ✗
➜  apigateway git:(kubernetes) ✗ curl -q -H "Content-type: application/json" -X POST -d '{"title1":"My project hello hello11"}'  `minikube service --url apigateway`/projects/
{
  "id": 123,
  "url": "Project-None"
}
➜  apigateway git:(kubernetes) ✗ curl -q -H "Content-type: application/json" -X POST -d '{"title":"An awesome project"}'  `minikube service --url apigateway`/projects/
{
  "id": 123,
  "url": "Project-An awesome project"
}



What's running on port 443 kubernetes?
======================================

**How to update service config changes**


References
==========

- https://kubernetes.io/docs/concepts/services-networking/connect-applications-service/
- https://medium.com/google-cloud/running-workloads-in-kubernetes-86194d133593
- https://kubernetes.io/docs/concepts/services-networking/service/#defining-a-service

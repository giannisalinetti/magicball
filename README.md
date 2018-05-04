# MagicBall - A tutorial about Init Containers in Kubernetes

This project has the purpose of demonstrating how Init Containers work inside a 
Kubernetes Pod. Pods are the minimimal unit of execution inside a Kubernetes cluster.
A Pod can embed one or more containers, sharing common namespaces. Thanks to this
feature a developer can build different parts of code in separated containers and have
them executed together using, for example, the same network namespace.
A special kind of containers in the Pod are the **Init container**, whose only purpose is to 
execute init related tasks and then exit before the execution of the main container(s).
More than one init container can me executed in the same Pod.

![](https://blog.openshift.com/wp-content/uploads/loap.png)

Image from the OpenShift Blog article : [Kubernetes: A Pod's Life](https://blog.openshift.com/kubernetes-pods-life/)

## Build

To demonstrate the behaviour of init containers, a simple application was written using
the Go programming language. It is a simple web appliation that writes the Magic 8 Ball random
answers. To do this the main app reads sinlge records from a MariaDB database randomly and 
writes the answer.
The purpose of the init container here is to initialize the database in an indempotent way.
So, when the pod starts for the very first time, the init container creates the table and 
populates it with the magic ball answers. 
If we scale up more pods or simply restart it the init container will simply contact the database,
find the table already filled, and exit.
Only after the execution of the init container the main http server process will be started in a 
second container.

The subfolders *magicball-init* and *magicball-server* contains respectively the init container 
source code and the http server source code, along with their own Dockerfiles.

**NOTE**: trying to build and run code locally will result in a failure since the database
connection parameters are strictly bound to Kubernetes service environment variables. To do so, modify
the *myConn* struct instance using fixed strings, flags or different environment variables.

Images for both magicball-init and magicball-server are already available: 
- docker.io/gbsal/magicball-init
- docker.io/gbsal/magicball-server

If you need to rebuild you own use the **docker build** command:

```
cd magicball-init
docker build -t customname/magicball-init .
```

Alternatively, you can also modify the Makefile to suit your needs.

The public **Golang 1.8** Docker image has been used. This is good for a proof-of-concept. For production
environment you may need a different images. For example, if you need to build and deploy in 
OpenShift the most suitable image is ths [golang-container](https://github.com/sclorg/golang-container) provided
by Software Collections, whose purpose is to build application using [Source2Image](https://github.com/openshift/source-to-image)
strategy.


## Install

To install the application you need to have **MiniKube** up and running on you machine.
Please refer to the [Getting Started Guide](https://kubernetes.io/docs/getting-started-guides/minikube/) to install and run
MiniKube.

To deploy the application in Kubernetes the following resources must be created:

- **Namespace**. We choose the name **magic8** for this example:
```
kubectl create namespace magic8
```

- **Database Pod**. This is a Deployment that delivers a MariaDB image along
with custom environment variables.
```
kubectl apply -f k8s/appdb-deployment.yml -n magic8
```
Please note that the database is **NOT** persistent.

- **Database Service**. The service used to expose the database in the cluster:
```
kubectl apply -f k8s/appdb-svc.yml -n magic8
```

- **Application Pod**. The Deployment that creates the pod with init and http container.
This Deployment also provides some extra variables for data source management.
```
kubectl apply -f k8s/magicball-deployment.yml -n magic8
```

During init stage you should see an output like this.
```
$ kubectl get pods -n magic8
NAME                         READY     STATUS        RESTARTS   AGE
appdb-696878b6b7-m8xw4       1/1       Running       0          3m
magicball-656dcc854f-kvzc5   0/1       Init:0/1      0          12s
```

This shows that the init container is running and doing his job.

- **Aplication Service**. The service used to expose the application. Note that this service
if of type **LoadBalancer**. This will enable the usage of an external Load Balancer, which,
in this case, will be the **minikube** process:
```
kubectl apply -f k8s/magicball-svc.yml -n magic8
```

After creating all the resources the service can be exposed outside the cluster with MiniKube:
```
minikube service magicball -n magic8 --url
```

This will print out the reachable url.

## Usage

Just test the application using the curl command to retrieve different answers. At the moment the 
application exposes a simple GET method. In the future it could be extended with a plain simple
UI to write an answer with a fake POST.
```
curl http://192.168.39.163:30021
```

The magicball will answer with sentences like these:
```
Magic 8 Ball said: Most likely.
```

## Author

Gianni Salinetti (gbsalinetti@gmail.com)

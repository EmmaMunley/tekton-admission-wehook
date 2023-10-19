# tekton-webhook-admission-webhook
This is a webhook admission webhook that adds validation for Tekton pipelines and tasks.

## Installation
This project can fully run locally and includes automation to deploy a local Kubernetes cluster (using Minkikube).

### Requirements
* Docker
* kubectl
* minikube
* Go >=1.19

## Usage
### Create Cluster
First, we need to create a Kubernetes cluster:
```
❯ minikube start

🔧 Creating Kubernetes cluster...
😄  minikube v1.30.1 on Darwin 13.6 (arm64)
🎉  minikube 1.31.2 is available! Download it: https://github.com/kubernetes/minikube/releases/tag/v1.31.2
💡  To disable this notice, run: 'minikube config set WantUpdateNotification false'

✨  Using the docker driver based on existing profile
👍  Starting control plane node minikube in cluster minikube
🚜  Pulling base image ...
🤷  docker "minikube" container is missing, will recreate.
🔥  Creating docker container (CPUs=2, Memory=7803MB) ...

🐳  Preparing Kubernetes v1.26.3 on Docker 23.0.2 ...

🔗  Configuring bridge CNI (Container Networking Interface) ...
🔎  Verifying Kubernetes components...
    ▪ Using image gcr.io/k8s-minikube/storage-provisioner:v5
🌟  Enabled addons: storage-provisioner, default-storageclass
🏄  Done! kubectl is now configured to use "minikube" cluster and "default" namespace by default
```

Make sure that the Kubernetes node is ready:
```
❯ kubectl get nodes
NAME                 STATUS   ROLES                  AGE     VERSION
minikube             Ready    control-plane,master   3m25s   v1.26.3
```

And that system pods are running happily:
```
❯ kubectl -n kube-system get pods
NAME                                         READY   STATUS    RESTARTS   AGE
coredns-558bd4d5db-thwvj                     1/1     Running   0          3m39s
etcd-minikube                                1/1     Running   0          3m56s
kube-apiserver-minikube                      1/1     Running   0          3m54s
kube-controller-manager-minikube             1/1     Running   0          3m56s
kube-proxy-4h6sj                             1/1     Running   0          3m40s
kube-scheduler-minikube                      1/1     Running   0          3m54s
```

### Deploy Admission Webhook
To configure the cluster to use the admission webhook and to deploy said webhook, simply run:
```
❯ make deploy

📦 Building tekton-webhook Docker image...
docker build -t tekton-webhook:latest .
[+] Building 14.3s (13/13) FINISHED
...

📦 Pushing admission-webhook image into Minikube's Docker daemon...
minikube image load tekton-webhook:latest

⚙️  Applying cluster config...
kubectl apply -f dev/manifests/cluster-config/
namespace/apps created
mutatingwebhookconfiguration.admissionregistration.k8s.io/tekton.webhook.config created
validatingwebhookconfiguration.admissionregistration.k8s.io/tekton.webhook.config created

🚀 Deploying tekton-webhook...
kubectl apply -f dev/manifests/webhook/
deployment.apps/tekton-webhook created
service/tekton-webhook created
secret/tekton-webhook-tls created
```

Then, make sure the admission webhook pod is running (in the `default` namespace):
```
❯ kubectl get pods
NAME                                        READY   STATUS    RESTARTS   AGE
tekton-webhook-77444566b7-wzwmx   1/1     Running   0          2m21s
```

You can stream logs from it:
```
❯ make logs

🔍 Streaming tekton-webhook logs...
kubectl logs -l app=tekton-webhook -f
time="2021-09-03T04:59:10Z" level=info msg="Listening on port 443..."
time="2021-09-03T05:02:21Z" level=debug msg=healthy uri=/health
```

And hit it's health endpoint from your local machine:
```
❯ curl -k https://localhost:8443/health
OK
```

### Deploying tasks
Deploy a valid task that gets successfully created:
```
❯ make valid-task

🚀 Deploying valid pod...
kubectl apply -f dev/manifests/tasks/valid-task.yaml
tasks/valid-task created
```
You should see in the admission webhook logs that the task was validated and created.

Deploy an invalid task that gets rejected:
```
❯ make invalid-task

🚀 Deploying "invalid" task...
kubectl apply -f dev/manifests/tasks/invalid-task.yaml
Error from server: error when creating "dev/manifests/tasks/invalid-task.yaml": admission webhook "tekton.webhook.config" denied the request: pod name contains "offensive"
```
You should see in the admission webhook logs that the pod validation failed.


## Admission Logic
A set of validations for pipelines and tasks are implemented in an extensible framework. Those happen on the fly when a pipeline/task is created and no further resources are tracked and updated (ie. no controller logic).

### Validating Webhooks
#### Implemented
- [pipeline name validation](pkg/validation/name_validator.go): validates that a pipeline name doesn't contain any offensive string
- [task name validation](pkg/validation/name_validator.go): validates that a task name doesn't contain any offensive string



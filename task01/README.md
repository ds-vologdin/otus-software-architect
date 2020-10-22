# Task 01

## Build the app

```bash
docker build . -t dsvologdin/app:1.0.1
docker login
docker push dsvologdin/app:1.0.1
```

## Manual local test the app
Run
```bash
docker run --rm -p 8000:8000 dsvologdin/app:1.0.1
```

Run in the other console
```bash
curl http://127.0.0.1:8000/healthz/
```

We must get
```
{"status": "OK"}
```

## Deploy in minikube

Activate ingress addon
```bash
minikube addons enable ingress
```

Deploy the app
```bash
cd k8s
kubectl apply -f .
```
### Check the deploy
```bash
kubectl get pod
```
Output
```
NAME                          READY   STATUS    RESTARTS   AGE
simple-app-6f55cc9c7c-8qnnw   1/1     Running   0          19m
simple-app-6f55cc9c7c-ppvs8   1/1     Running   0          20m
```
```bash
kubectl get svc
```
Output
```
bud@bud-notebook:~/go/src/github.com/ds-vologdin/otus-software-architect/task01/k8s$ kubectl get svc
NAME         TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)    AGE
simple-app   ClusterIP   10.97.27.163   <none>        8000/TCP   19h
```
```bash
kubectl get ingress
```
Output
```
NAME         CLASS    HOSTS           ADDRESS         PORTS   AGE
simple-app   <none>   arch.homework   192.168.39.78   80      19h
```

## Manual test the endpoint

```bash
curl -H "Host: arch.homework" http://192.168.39.78/otusapp/petrovich/healthz/
```
Output
```
{"status": "OK"}
```

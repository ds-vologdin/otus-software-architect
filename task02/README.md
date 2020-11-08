# Task 02

## Build the accounts app

```bash
docker build -t dsvologdin/accounts:0.0.3 .
docker build -t dsvologdin/migrate-account:0.0.2 .
docker login
docker push dsvologdin/accounts:0.0.3
docker push dsvologdin/migrate-account:0.0.2
```

## Deploy in minikube

Activate ingress addon
```bash
minikube addons enable ingress
```

Deploy the app
```bash
helm install accounts helm/accounts --atomic
```

Upgrade the app
```bash
helm upgrade accounts helm/accounts --atomic
```

If you don't want to deploy the chart to the default namespace, you can use flag `-n your_namespace`.

### Check the deploy

```bash
kubectl get all
```
Output
```
NAME                            READY   STATUS    RESTARTS   AGE
pod/accounts-79c944645c-rw8pk   1/1     Running   0          34m
pod/accounts-postgresql-0       1/1     Running   0          3h18m

NAME                                   TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
service/accounts                       ClusterIP   10.107.12.186   <none>        8000/TCP   3h19m
service/accounts-postgresql            ClusterIP   10.99.44.149    <none>        5432/TCP   3h19m
service/accounts-postgresql-headless   ClusterIP   None            <none>        5432/TCP   3h19m
service/kubernetes                     ClusterIP   10.96.0.1       <none>        443/TCP    3h22m

NAME                       READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/accounts   1/1     1            1           3h19m

NAME                                  DESIRED   CURRENT   READY   AGE
replicaset.apps/accounts-79c944645c   1         1         1       34m
replicaset.apps/accounts-7d8c7d969d   0         0         0       3h18m
replicaset.apps/accounts-99c594456    0         0         0       112m

NAME                                   READY   AGE
statefulset.apps/accounts-postgresql   1/1     3h18m

```

# Postman

You can use postman collection for test endpoints.
Before that, you should set the arch.homework IP address in `/etc/hosts`.

```bash
newman run postman/accounts.postman_collection.json
```

Output

```
→ Create user
  POST http://arch.homework/otusapp/user/ [201 Created, 169B, 111ms]
  ✓  status test

→ Get user
  GET http://arch.homework/otusapp/user/11 [200 OK, 274B, 10ms]
  ✓  status test
  ✓  body test

→ Edit user
  PUT http://arch.homework/otusapp/user/11 [200 OK, 170B, 41ms]
  ✓  status test
  ✓  body test

→ Get user
  GET http://arch.homework/otusapp/user/11 [200 OK, 278B, 13ms]
  ✓  status test
  ✓  body test

→ Delete user
  DELETE http://arch.homework/otusapp/user/11 [200 OK, 170B, 25ms]
  ✓  status test
  ✓  body test

→ Get deleted user
  GET http://arch.homework/otusapp/user/11 [404 Not Found, 188B, 7ms]
  ✓  status test
  ✓  body test

┌─────────────────────────┬───────────────────┬──────────────────┐
│                         │          executed │           failed │
├─────────────────────────┼───────────────────┼──────────────────┤
│              iterations │                 1 │                0 │
├─────────────────────────┼───────────────────┼──────────────────┤
│                requests │                 6 │                0 │
├─────────────────────────┼───────────────────┼──────────────────┤
│            test-scripts │                12 │                0 │
├─────────────────────────┼───────────────────┼──────────────────┤
│      prerequest-scripts │                 6 │                0 │
├─────────────────────────┼───────────────────┼──────────────────┤
│              assertions │                11 │                0 │
├─────────────────────────┴───────────────────┴──────────────────┤
│ total run duration: 540ms                                      │
├────────────────────────────────────────────────────────────────┤
│ total data received: 311B (approx)                             │
├────────────────────────────────────────────────────────────────┤
│ average response time: 34ms [min: 7ms, max: 111ms, s.d.: 36ms] │
└────────────────────────────────────────────────────────────────┘

```

# Prometheus

## Install

```bash
minikube addons disable ingress

helm install prom stable/prometheus-operator -f prometheus.yaml --atomic

helm install nginx stable/nginx-ingress -f helm/nginx-ingress/values.yaml --atomic
```

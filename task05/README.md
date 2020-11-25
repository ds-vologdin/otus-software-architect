# Deploy

You should install and start minikube.
```bash
minikube start --cpus=4 --memory=6g --driver=kvm2
```

We will deploy all our charts to default namespace. **Don't use the default namespace in production.**

Install prometheus. It is an optional step. If you don't want deploy the prometheus, you need to disable metrics in value files of all services.

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm install prom prometheus-community/prometheus-operator -f infra/helm/prometheus/values.yaml --atomic
```

Install nginx-ingress.

```bash
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm install nginx ingress-nginx/ingress-nginx -f infra/helm/nginx-ingress/values.yaml --atomic
```

Install account service. We will deploy a postgresql as a dependency and migrate a database scheme.

```bash
helm install accounts account/helm --atomic
```

Install auth service.

```bash
helm install auth auth/helm --atomic
```

# Test

Should run postman's test

```bash
cd postman
newman run accounts_auth.postman_collection.json
```

### Rusult

```
newman

accounts auth

→ register tom
  POST http://arch.homework/accounts/profile/ [201 Created, 147B, 457ms]
  ✓  status test

→ register: user already exists
  POST http://arch.homework/accounts/profile/ [409 Conflict, 200B, 23ms]
  ✓  status test

→ Unathorized get profile
  GET http://arch.homework/accounts/profile/{{userId}} [401 Unauthorized, 205B, 9ms]
  ✓  status test

→ Unathorized change profile
  POST http://arch.homework/accounts/profile/{{userId}} [401 Unauthorized, 205B, 7ms]
  ✓  status test

→ login tom
  POST http://arch.homework/auth/token/refresh [201 Created, 1.03KB, 21ms]
  ✓  status test

→ get profile tom
  GET http://arch.homework/accounts/profile/36 [200 OK, 263B, 11ms]
  ✓  status test
  ✓  body test

→ edit profile tom
  PATCH http://arch.homework/accounts/profile/36 [200 OK, 148B, 83ms]
  ✓  status test

→ check edited profile tom
  GET http://arch.homework/accounts/profile/36 [200 OK, 263B, 11ms]
  ✓  status test
  ✓  body test

→ register user bob
  POST http://arch.homework/accounts/profile/ [201 Created, 147B, 39ms]
  ✓  status test

→ login bob
  POST http://arch.homework/auth/token/refresh [201 Created, 1.03KB, 89ms]
  ✓  status test

→ get alien profile
  GET http://arch.homework/accounts/profile/36 [403 Forbidden, 191B, 9ms]
  ✓  status test

→ delete profile tom
  DELETE http://arch.homework/accounts/profile/36 [200 OK, 148B, 352ms]
  ✓  status test

→ delete profile bob
  DELETE http://arch.homework/accounts/profile/37 [200 OK, 148B, 166ms]
  ✓  status test

┌─────────────────────────┬───────────────────┬───────────────────┐
│                         │          executed │            failed │
├─────────────────────────┼───────────────────┼───────────────────┤
│              iterations │                 1 │                 0 │
├─────────────────────────┼───────────────────┼───────────────────┤
│                requests │                13 │                 0 │
├─────────────────────────┼───────────────────┼───────────────────┤
│            test-scripts │                26 │                 0 │
├─────────────────────────┼───────────────────┼───────────────────┤
│      prerequest-scripts │                13 │                 0 │
├─────────────────────────┼───────────────────┼───────────────────┤
│              assertions │                15 │                 0 │
├─────────────────────────┴───────────────────┴───────────────────┤
│ total run duration: 1981ms                                      │
├─────────────────────────────────────────────────────────────────┤
│ total data received: 2.19KB (approx)                            │
├─────────────────────────────────────────────────────────────────┤
│ average response time: 98ms [min: 7ms, max: 457ms, s.d.: 139ms] │
└─────────────────────────────────────────────────────────────────┘

```


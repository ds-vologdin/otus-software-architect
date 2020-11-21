# Deploy

# Test deploy

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
# Progressive Delivery with Aristi, Argo Rollouts, and Istio

This tutorial will guide you through deploying Aristi with Argo Rollouts and Istio for progressive delivery in Kubernetes.

## Prerequisites

Ensure you have the following installed and configured:

- Kubernetes cluster (v1.26+ recommended)

- kubectl CLI

- Argo Rollouts installed: Installation Guide

- Istio service mesh installed: Istio Installation

## Step 1: Install Aristi 

Follow the steps to install Aristi from [README](./README.md)

Verify the CRD is installed:

```bash
kubectl get crds | grep aristi
```

## Step 2: Deploy the Aristi Resource

Create a file named `aristi.yaml` with the following content:

```yaml
apiVersion: aristi.cloudstation/v1alpha1
kind: Aristi
metadata:
  labels:
    app.kubernetes.io/name: aristi
    app.kubernetes.io/managed-by: kustomize
  name: aristi-sample
  namespace: default
spec:
  gateway:
    name: argo-gateway
    spec:
      selector:
        istio: ingressgateway # use the default IngressGateway
      servers:
        - port:
            number: 80
            name: http
            protocol: HTTP
          hosts:
            - "*"
  istio:
    hosts:
      - "*"
    gateways:
      - argo-gateway
    virtualService:
      name: rollout-vsvc
      routes:
        - destination:
            host: weather-test-app-hyd
          weight: 10
        - destination:
            host: weather-test-app-ny
          weight: 90
  rollout:
    services:
      stable:
        name: weather-test-app-hyd
        ports:
          - protocol: "TCP"
            port: 80
            targetPort: 5000
        type: ClusterIP
      canary:
        name: weather-test-app-ny
        ports:
          - protocol: "TCP"
            port: 80
            targetPort: 5000
        type: ClusterIP
    selector:
      matchLabels:
        app: rollout-istio
    replicas: 2
    template:
      metadata:
        labels:
          app: rollout-istio
      spec:
        containers:
          - name: weather-app-falcon
            image: docker.io/atulinfracloud/weathersample:v1
    strategy:
      canary:
        canaryService: weather-test-app-ny
        stableService: weather-test-app-hyd
        steps:
          - setWeight: 50
          - pause:
              duration: "10m"
```

Apply the Aristi resource:

```bash
kubectl apply -f aristi.yaml
```

## Step 3: Verify the Deployment

### Check Aristi Resource Status

Ensure the Aristi resource is created successfully:

```bash
kubectl get aristi -n default
```

You should see output similar to:

```
NAME            AGE
aristi-sample   5m
```

### Validate Istio Gateway and VirtualService

Check the Istio Gateway:

```bash
kubectl get gateway -n default
```

Check the VirtualService:

```bash
kubectl get virtualservice -n default
```

### Verify Argo Rollouts

Check the rollout status:

```bash
kubectl argo rollouts get rollout aristi-sample-rollout -n default
```

Monitor the rollout progress:

```bash
kubectl argo rollouts dashboard
```

## Step 4: Test the Progressive Delivery

The next step is to promote a canary version of the app using the following command:

```bash
kubectl argo rollouts set image aristi-sample-rollout weather-app-falcon=docker.io/atulinfracloud/weathersample:v2
```

Use `curl` or a web browser to access your Ingress Gateway:

```bash
kubectl port-forward svc/istio-ingressgateway -n istio-system 8080:80
```
We now have v2 deployed in canary while v1 in stable. Visit the service at the http://localhost:8080 and refresh the page. You’ll see stable (Hyderabad) and canary (New York) versions of the apps displayed randomly based on the weights provided.

You should see traffic being routed to both `weather-test-app-hyd` and `weather-test-app-ny` according to the defined weights (90% and 10% respectively).

## Step 5: Clean Up

To remove the Aristi deployment and associated resources:

```bash
kubectl delete -f aristi.yaml

 kubectl delete svc,deploy,po,ingress,virtualservice,gateway,rollout -n default -l aristi.io/managed-by=aristi

```

## Conclusion

You have successfully deployed Aristi with Argo Rollouts and Istio, enabling progressive delivery with traffic shifting and advanced rollout strategies.


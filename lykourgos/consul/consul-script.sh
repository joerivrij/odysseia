#! /bin/bash

minikube start --kubernetes-version v1.19.4  --cpus 8 --memory 8292
minikube addons enable ingress
minikube addons enable ingress-dns

kubectl create ns system-apps
kubens system-apps
helm repo add hashicorp https://helm.releases.hashicorp.com
kubectl create secret generic consul-gossip-encryption-key --from-literal=key=$(consul keygen)
helm install -f config.yaml consul hashicorp/consul --version "0.32.0"

export CONSUL_HTTP_ADDR=https://127.0.0.1:8501
export CONSUL_HTTP_TOKEN=$(kubectl get secrets/consul-bootstrap-acl-token --template={{.data.token}} | base64 -d)

minikube service list
minikube ip

kubectl create ns monitoring-apps
kubens monitoring-apps

helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add grafana https://grafana.github.io/helm-charts

helm upgrade --install prometheus prometheus-community/prometheus
helm upgrade --install grafana grafana/grafana
helm upgrade --install loki grafana/loki-stack

kubectl get secret --namespace monitoring-apps grafana -o jsonpath="{.data.admin-password}" | base64 --decode ; echo

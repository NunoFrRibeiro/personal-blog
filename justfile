timeout := "300s"

default:
  just --list

create-cluster:
  kind create cluster
  kind kubectl create namespace argo

destroy-cluster:
  kind delete cluster

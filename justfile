timeout := "400s"
ARGO_WORKFLOWS_VERSION := "v3.5.10"
wf_url := """
  ARGO WORKFLOWS
  https://localhost:2746
"""

# Creates the kind cluster and deploys Ar
create-cluster:
  kind create cluster
  kubectl cluster-info --context kind-kind

# Deploys Argo-Worklows to the kind cluster
wf: create-cluster
  kubectl create namespace argo
  kubectl apply -n argo -f "https://github.com/argoproj/argo-workflows/releases/download/{{ARGO_WORKFLOWS_VERSION}}/quick-start-minimal.yaml" 
  kubectl wait --for=jsonpath='{.status.phase}'=Running po -n argo --all --timeout={{timeout}}
  kubectl -n argo port-forward service/argo-server 2746:2746 &
  @echo "{{wf_url}}"

events: create-cluster wf
  @echo "TBD"
destroy:
  kind delete cluster

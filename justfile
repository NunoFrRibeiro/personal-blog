timeout := "400s"
ARGO_WORKFLOWS_VERSION := "v3.5.10"
wf_url := """
  ARGO WORKFLOWS
  http://localhost:2746
"""

# Creates the kind cluster and deploys
create-cluster:
  kind create cluster
  kubectl cluster-info --context kind-kind

# Deploys Argo-Worklows to the kind cluster
wf:
  -helm repo add argo https://argoproj.github.io/argo-helm
  -helm upgrade --install argo-workflows argo/argo-workflows --namespace argo --create-namespace \
    --values="./values/argo-workflow.yaml" --wait
  -kubectl -n argo port-forward service/argo-workflows-server 2746:2746 &

# Deploys Argo-Events to the kind cluster
events:
  helm upgrade --install argo-events argo/argo-events --namespace argo --create-namespace \
    --values="./values/argo-events.yaml" --wait

# Start the demo
demo:
  -kubectl create secret generic -n argo github --from-literal token=$(infisical secrets get GITHUB_API_TOKEN --plain)
  -kubectl create secret generic -n argo dagger-cloud --from-literal=token=$(infisical secrets get DAGGER_CLOUD --plain)
  kubectl create -f dagger-workflow.yaml

# Deploy the full suite for the demo
start: create-cluster wf events
  @echo "All Deployed"
  @echo "{{wf_url}}"

# Clean all up
destroy:
  kind delete cluster

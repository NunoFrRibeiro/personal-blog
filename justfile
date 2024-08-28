timeout := "400s"
ARGO_WORKFLOWS_VERSION := "v3.5.10"
wf_url := """
  ARGO WORKFLOWS
  https://localhost:2746
"""
cd_url := """
  ARGO WORKFLOWS
  https://localhost:8080
"""

# Creates the kind cluster and deploys Ar
create-cluster:
  kind create cluster
  kubectl cluster-info --context kind-kind

# Deploys Argo-CD to the cluster
cd:
  -helm repo add argo https://argoproj.github.io/argo-helm
  -helm upgrade --install argo-cd argo/argo-cd --namespace argocd --create-namespace \
    --values ./values/argocd.yaml --wait
  kubectl port-forward service/argo-cd-argocd-server -n argocd 8080:443 &

# Deploys Argo-Worklows to the kind cluster
wf:
  -helm upgrade --install argo-workflows argo/argo-workflows --namespace argo --create-namespace \
    --values="./values/argo-workflow.yaml"
  -kubectl wait --for=condition=Ready po -n argo --all --timeout={{timeout}}
  -kubectl -n argo port-forward service/argo-workflows-server 2746:2746 &

# Deploys Argo-Events to the kind cluster
events:
  helm upgrade --install argo-events argo/argo-events --namespace argo-events --create-namespace \
    --values="./values/argo-events.yaml" --wait

# Deploy the full suite for the demo
start: create-cluster cd wf events
  @echo "All Deployed"
  @echo "{{wf_url}}"
  @echo "{{cd_url}}"
  @echo admin pass=$(kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d)

# Clean all up
destroy:
  kind delete cluster

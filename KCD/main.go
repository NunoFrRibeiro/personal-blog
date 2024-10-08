package main

import (
	"context"

	"dagger/kcd/internal/dagger"
)

type Kcd struct {
	Name      string
	KCDServer *dagger.Service
	// +private
	Source *dagger.Directory
}

func New(
	ctx context.Context,
	// k3s cluster name
	// +default="default"
	name string,
	// source directory to build
	// +private
	// +defaultPath="../"
	// +optional
	source *dagger.Directory,
) *Kcd {
	return &Kcd{
		Name:   name,
		Source: source,
	}
}

func (c *Kcd) CreateCluster(ctx context.Context) *Kcd {
	// c.KCDServer = dag.K3S(c.Name).Server()
	kc := dag.K3S(c.Name).Container()
	kc = kc.WithMountedCache("/var/lib/dagger", dag.CacheVolume("varlibdagger"))

	c.KCDServer = dag.K3S(c.Name).WithContainer(kc).Server()
	return c
}

func (c *Kcd) GetConfig() *dagger.File {
	return dag.K3S(c.Name).Config(dagger.K3SConfigOpts{
		Local: false,
	})
}

func (c *Kcd) KNS() *dagger.Container {
	return dag.K3S(c.Name).Kns().Terminal()
}

func (c *Kcd) TestWF(
	ctx context.Context,
	// Auth Client Id to fetch credentials
	// +required
	authClientId *dagger.Secret,
	// Auth Client Id to fetch credentials
	// +required
	authClientSecret *dagger.Secret,
) (*dagger.Container, error) {
	_, err := c.CreateCluster(ctx).KCDServer.Start(ctx)
	if err != nil {
		return nil, err
	}

	kubeConfig := c.GetConfig()

	daggerCloud := dag.Infisical(authClientId, authClientSecret).
		GetSecret("DAGGER_CLOUD", "495b60ca-a6c5-46e9-bc08-6e37b1d715de", "dev")

	return dag.Container().From("bitnami/kubectl:1.31.0-debian-12-r4").
		WithUser("root").
		WithFile("/.kube/config", kubeConfig).
		WithEnvVariable("KUBECONFIG", "/.kube/config").
		WithSecretVariable("DAGGER_CLOUD", daggerCloud).
		WithSecretVariable("INFISICAL_SECRET", authClientSecret).
		WithSecretVariable("INFISICAL_ID", authClientId).
		WithExec([]string{"chown", "1001:0", "/.kube/config"}).
		WithExec([]string{
			"bash",
			"-c",
			"apt update && apt install -y curl",
		}).
		WithDirectory("/demo", c.Source).
		WithExec([]string{
			"bash",
			"-c",
			"curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash",
		}).
		WithExec([]string{
			"bash",
			"-c",
			"helm repo add argo https://argoproj.github.io/argo-helm",
		}).
		WithExec([]string{
			"bash",
			"-c",
			"helm upgrade --install --namespace=dagger --create-namespace dagger oci://registry.dagger.io/dagger-helm --wait",
		}).
		WithExec([]string{
			"bash",
			"-c",
			"helm upgrade --install argo-wf argo/argo-workflows --namespace argo --create-namespace --values=/demo/values/argo-workflow.yaml --wait",
		}).
		WithExec([]string{
			"bash",
			"-c",
			"kubectl create secret generic -n argo dagger-cloud --from-literal=token=$DAGGER_CLOUD",
		}).
		WithExec([]string{
			"bash",
			"-c",
			"kubectl create secret generic -n argo infisical-secret --from-literal=infisical_secret=$INFISICAL_SECRET --from-literal=infisical_id=$INFISICAL_ID",
		}).
		WithExec([]string{
			"bash",
			"-c",
			"kubectl create rolebinding default-admin --clusterrole=admin --serviceaccount=argo:dagger-workflow",
		}).
		WithExec([]string{
			"bash",
			"-c",
			"kubectl create -f /demo/dagger-test-workflow.yaml",
		}), nil
}

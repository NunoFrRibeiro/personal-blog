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
	c.KCDServer = dag.K3S(c.Name).Server()
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

func (c *Kcd) DemoStart(
	ctx context.Context,
) (string, error) {
	_, err := c.CreateCluster(ctx).KCDServer.Start(ctx)
	if err != nil {
		return "", err
	}

	kubeConfig := c.GetConfig()

	return dag.Container().From("alpine:latest").
		WithUser("root").
		WithFile("/.kube/config", kubeConfig).
		WithEnvVariable("KUBECONFIG", "/.kube/config").
		WithExec([]string{
			"sh",
			"-c",
			"apk update && apk add curl openssl",
		}).
		WithDirectory("/demo", c.Source).
		WithExec([]string{
			"sh",
			"-c",
			"curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | sh",
		}).
		WithExec([]string{
			"sh",
			"-c",
			"helm repo add argo https://argoproj.github.io/argo-helm",
		}).
		WithExec([]string{
			"sh",
			"-c",
			"helm upgrade --install argo-cd argo/argo-cd --namespace argocd --create-namespace --values=/demo/values/argocd.yaml --wait",
		}).
		WithExec([]string{
			"sh",
			"-c",
			"helm upgrade --install argo-wf argo/argo-workflows --namespace argo --create-namespace --values=/demo/values/argo-workflow.yaml",
		}).
		WithExec([]string{
			"sh",
			"-c",
			"helm upgrade --install argo-events argo/argo-events --namespace argo-events --create-namespace --values=/demo/values/argo-events.yaml --wait",
		}).Stdout(ctx)
}

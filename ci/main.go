package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"dagger/backend/internal/dagger"
)

var (
	APP     = "personal-blog"
	GH_REPO = "https://github.com/NunoFrRibeiro/personal-blog"
	DH_REPO = "index.docker.io"
	IMAGE   = "nunofilribeiro/go-blog:v0.1.0"
)

type Goblog struct {
	// +private
	Source *dagger.Directory
}

func New(
	// Project source directory
	// +optional
	source *dagger.Directory,

	// Checkout the repository (at the designated ref) and use it as the source directory instead of the local one.
	// +optional
	ref string,
) (*Goblog, error) {
	if source == nil && ref != "" {
		source = dag.Git("https://github.com/NunoFrRibeiro/personal-blog.git", dagger.GitOpts{
			KeepGitDir: true,
		}).Ref(ref).Tree()
	}

	if source == nil {
		return nil, errors.New("either source or ref is needed")
	}

	return &Goblog{
		Source: source,
	}, nil
}

// Run unit tests on the Project
func (g *Goblog) RunUnitTests(
	ctx context.Context,
) (string, error) {
	result, err := dag.Backend().RunUnitTests(ctx, g.Source)
	if err != nil {
		return "", err
	}

	return result, nil
}

// Lint the Project
func (g *Goblog) Lint(
	ctx context.Context,
) (string, error) {
	result, err := dag.Backend().Lint(ctx, g.Source)
	if err != nil {
		return "", err
	}

	return result, nil
}

// Serve the blog on port 8080
func (g *Goblog) Serve() *dagger.Service {
	backendService := dag.Backend().Serve(g.Source)
	numInt32, _ := strconv.Atoi("8081")

	return dag.Proxy().
		WithService(backendService, "backend", numInt32, numInt32, false).
		Service()
}

// Deploy the blog to fly.io
func (g *Goblog) Deploy(
	ctx context.Context,
	flyToken *dagger.Secret,
	registryUser string,
	registryPass *dagger.Secret,
) (string, error) {

	source := g.Source

	blogAmd64 := dag.Backend().Container(source, dagger.BackendContainerOpts{
		Arch: "amd64",
	})

	blogArm64 := dag.Backend().Container(source, dagger.BackendContainerOpts{
		Arch: "arm64",
	})

	_, err := dag.Container().
		WithRegistryAuth(DH_REPO, registryUser, registryPass).
		Publish(ctx, IMAGE, dagger.ContainerPublishOpts{
			PlatformVariants: []*dagger.Container{
				blogAmd64,
				blogArm64,
			},
		})

	if err != nil {
		return "", err
	}

	result, err := dag.Flyio(flyToken, dagger.FlyioOpts{
		Version: "latest",
		Regions: "mad",
		Org:     "nuno-ribeiro",
	}).Deploy(ctx, source, dagger.FlyioDeployOpts{
		Image: IMAGE,
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Blog deployed to Fly.io %s", result), nil
}

func (g *Goblog) RunAll(
	ctx context.Context,
	// Infisical Auth Client ID
	// +required
	infisicalClientId *dagger.Secret,
	// Infisical Auth Client Secret
	// +required
	infisicalClientSecret *dagger.Secret,
	// Infisical Project to fetch secrets
	// +required
	infisicalProject string,
) (string, error) {

	// Lint the source
	result, err := g.Lint(ctx)
	if err != nil {
		return "", err
	}

	// Run all Unit Tests
	testResult, err := g.RunUnitTests(ctx)
	if err != nil {
		return "", err
	}

	result = result + "\n" + testResult

	// Deploy to Fly.io
	if infisicalClientId != nil && infisicalProject != "" {

		flyToken := dag.Infisical(infisicalClientId, infisicalClientSecret).GetSecret("FLY_TOKEN", infisicalProject, "dev", dagger.InfisicalGetSecretOpts{
			SecretPath: "/flyio",
		})

		registryUser, err := dag.Infisical(infisicalClientId, infisicalClientSecret).GetSecret("DH_USER", infisicalProject, "dev", dagger.InfisicalGetSecretOpts{
			SecretPath: "/flyio",
		}).Plaintext(ctx)
		if err != nil {
			return "", err
		}

		registryPass := dag.Infisical(infisicalClientId, infisicalClientSecret).GetSecret("DH_PASS", infisicalProject, "dev", dagger.InfisicalGetSecretOpts{
			SecretPath: "/",
		})

		deployResult, err := g.Deploy(ctx, flyToken, registryUser, registryPass)
		if err != nil {
			return "", err
		}

		result = result + "\n" + deployResult
	}

	return result, nil
}

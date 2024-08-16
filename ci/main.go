package main

import (
	"context"
	"fmt"
	"strconv"

	"dagger/backend/internal/dagger"
)

var (
	APP     = "personal-blog"
	DH_REPO = "https://hub.docker.com/r/nunofilribeiro/go-blog"
	GH_REPO = "https://github.com/NunoFrRibeiro/personal-blog"
	IMAGE   = "nunofilribeiro/go-blog:latest"
)

type Goblog struct{}

// Run unit tests on the Project
func (g *Goblog) RunUnitTests(
	ctx context.Context,
	// Point to the host directory where the project is located
	// +required
	dir *dagger.Directory,
) (string, error) {
	result, err := dag.Backend().RunUnitTests(ctx, dir)
	if err != nil {
		return "", err
	}

	return result, nil
}

// Lint the Project
func (g *Goblog) Lint(
	ctx context.Context,
	// Point to the host directory where the project is located
	// +required
	dir *dagger.Directory,
) (string, error) {
	result, err := dag.Backend().Lint(ctx, dir)
	if err != nil {
		return "", err
	}

	return result, nil
}

// Serve the blog on port 8080
func (g *Goblog) Serve(
	// Point to the host directory where the project is located
	// +required
	source *dagger.Directory,
) *dagger.Service {
	backendService := dag.Backend().Serve(source)
	numInt32, _ := strconv.Atoi("8081")

	return dag.Proxy().
		WithService(backendService, "backend", numInt32, numInt32, false).
		Service()
}

// Deploy the blog to fly.io
func (g *Goblog) Deploy(
	ctx context.Context,
	// Point to the host directory where the project is located
	// +required
	source *dagger.Directory,
	flyToken *dagger.Secret,
	registryUser string,
	registryPass *dagger.Secret,
) (string, error) {

	blogAmd64 := dag.Backend().Container(source, dagger.BackendContainerOpts{
		Arch: "amd64",
	})

	blogArm64 := dag.Backend().Container(source, dagger.BackendContainerOpts{
		Arch: "arm64",
	})

	_, err := dag.Container().
		//WithRegistryAuth(DH_REPO, registryUser, registryPass).
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
	// Point to the host directory where the project is located
	// +required
	source *dagger.Directory,
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
	result, err := g.Lint(ctx, source)
	if err != nil {
		return "", err
	}

	// Run all Unit Tests
	testResult, err := g.RunUnitTests(ctx, source)
	if err != nil {
		return "", err
	}

	result = result + "\n" + testResult

	// Deploy to Fly.io
	if infisicalClientId != nil && infisicalProject != "" {

		flyTokenStr, err := dag.Infisical(infisicalClientId, infisicalClientSecret).GetSecret("FLY_TOKEN", infisicalProject, "dev", dagger.InfisicalGetSecretOpts{
			SecretPath: "/flyio",
		}).Plaintext(ctx)
		if err != nil {
			return "", err
		}
		flyToken := dag.SetSecret("val_1", flyTokenStr)

		registryUser, err := dag.Infisical(infisicalClientId, infisicalClientSecret).GetSecret("DH_USER", infisicalProject, "dev", dagger.InfisicalGetSecretOpts{
			SecretPath: "/flyio",
		}).Plaintext(ctx)
		if err != nil {
			return "", err
		}

		registryPassStr, err := dag.Infisical(infisicalClientId, infisicalClientSecret).GetSecret("DH_PASS", infisicalProject, "dev", dagger.InfisicalGetSecretOpts{
			SecretPath: "/flyio",
		}).Plaintext(ctx)
		if err != nil {
			return "", err
		}
		registryPass := dag.SetSecret("val_2", registryPassStr)

		deployResult, err := g.Deploy(ctx, source, flyToken, registryUser, registryPass)
		if err != nil {
			return "", err
		}

		result = result + "\n" + deployResult
	}

	return result, nil
}

package main

import (
	"context"
	"strconv"

	"dagger/backend/internal/dagger"
)

var (
	APP     = "go-blog"
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

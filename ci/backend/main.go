package main

import (
	"context"
	"runtime"

	"dagger/backend/internal/dagger"
)

type Backend struct{}

// Run the projects unit tests
func (b *Backend) RunUnitTests(
	ctx context.Context,
	// Point to the host directory where the project is located
	// +required
	dir *dagger.Directory,
) (string, error) {
	return dag.Golang().
		WithProject(dir).
		Test(ctx, dagger.GolangTestOpts{
      Source: dir,
      Args: []string{"./..."},
    })
}

// Lint the project
func (b *Backend) Lint(
	ctx context.Context,
	// Point to the host directory where the project is located
	// +required
	dir *dagger.Directory,
) (string, error) {
	return dag.Golang().
		WithProject(dir).
		GolangciLint(ctx)
}

// Builds the project as a binary
func (b *Backend) BuildProject(
	// Point to the host directory where the project is located
	// +required
	dir *dagger.Directory,
	// If needded, specify the archtecture of the binary
	// +optional
	arch string,
) *dagger.Directory {
	if arch == "" {
		arch = runtime.GOARCH
	}

	return dag.Golang().
		WithProject(dir).
		Build([]string{}, dagger.GolangBuildOpts{
			Arch: arch,
		})
}

// Returns the built binary
func (b *Backend) Binary(
	// Point to the host directory where the project is located
	// +required
	dir *dagger.Directory,
	// If needded, specify the archtecture of the binary
	// +optional
	arch string,
) *dagger.File {
	binary := b.BuildProject(dir, arch)

	return binary.File("personal-blog")
}

// Create a container to run the binary
func (b *Backend) Container(
	// Point to the host directory where the project is located
	// +required
	dir *dagger.Directory,
	// If needded, specify the archtecture of the binary
	// +optional
	arch string,
) *dagger.Container {
	if arch == "" {
		arch = runtime.GOARCH
	}

	binary := b.Binary(dir, arch)

	return dag.
		Container(dagger.ContainerOpts{
			Platform: dagger.Platform(arch),
		}).
		From("ubuntu:24.10").
		WithWorkdir("/opt/blog").
		WithDirectory("posts", dir.Directory("posts")).
		WithDirectory("static", dir.Directory("static")).
		WithDirectory("templates", dir.Directory("templates")).
		WithFile("blog-bin", binary).
		WithEntrypoint([]string{"./blog-bin"}).
		WithExposedPort(8081)
}

// Run a service to test the go-blog
func (b *Backend) Serve(
	// Point to the host directory where the project is located
	// +required
	dir *dagger.Directory,
) *dagger.Service {
	return b.Container(dir, runtime.GOARCH).AsService()
}

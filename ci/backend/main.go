// A generated module for Backend functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"runtime"
)

type Backend struct{}

// Run the projects unit tests
func (b *Backend) RunUnitTests(
	ctx context.Context,
	// Point to the host directory where the project is located
	// +required
	dir *Directory,
) (string, error) {
	return dag.Golang().
		WithProject(dir).
		Test(ctx)
}

// Lint the project
func (b *Backend) Lint(
	ctx context.Context,
	// Point to the host directory where the project is located
	// +required
	dir *Directory,
) (string, error) {
	return dag.Golang().
		WithProject(dir).
		GolangciLint(ctx)
}

// Builds the project as a binary
func (b *Backend) BuildProject(
	// Point to the host directory where the project is located
	// +required
	dir *Directory,
	// If needded, specify the archtecture of the binary
	// +optional
	arch string,
) *Directory {

	if arch == "" {
		arch = runtime.GOARCH
	}

	return dag.Golang().
		WithProject(dir).
		Build([]string{}, GolangBuildOpts{
			Arch: arch,
		})
}

// Returns the built binary
func (b *Backend) Binary(
	// Point to the host directory where the project is located
	// +required
	dir *Directory,
	// If needded, specify the archtecture of the binary
	// +optional
	arch string,
) *File {

	binary := b.BuildProject(dir, arch)

	return binary.File("go-blog")
}

// Create a container to run the binary
func (b *Backend) Container(
	// Point to the host directory where the project is located
	// +required
	dir *Directory,
	// If needded, specify the archtecture of the binary
	// +optional
	arch string,
) *Container {

	if arch == "" {
		arch = runtime.GOARCH
	}

	binary := b.Binary(dir, arch)

	return dag.
		Container(ContainerOpts{
			Platform: Platform(arch),
		}).
		From("golang:1.21.7-alpine3.18@sha256:2a690031ca4e8affec4c25b1dbbafd5e734527b6ac6c8b58001cd6f342af7a44").
		WithFile("/bin/go-blog", binary).
		WithEntrypoint([]string{"/bin/go-blog"}).
		WithExposedPort(8080)
}

// Run a service to test the go-blog
func (b *Backend) Serve(
	// Point to the host directory where the project is located
	// +required
	dir *Directory,
) *Service {
	return b.Container(dir, runtime.GOARCH).AsService()
}

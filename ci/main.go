// A generated module for Goblog functions
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
	"fmt"
)

var (
  APP = "go-blog"
  GH_REPO = "https://github.com/NunoFrRibeiro/personal-go-blog"
  IMAGE = "nunofilribeiro/go-blog:tagname"
)
type Goblog struct{}

// Run unit tests on the Project
func (g *Goblog) RunUnitTests(
  ctx context.Context,
	// Point to the host directory where the project is located
	// +required
	dir *Directory,
) (string, error) {
  return fmt.Sprintf("Hola"), nil
}

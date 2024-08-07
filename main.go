package main

import (
	app "github.com/NunoFrRibeiro/personal-blog/cmd"
	"github.com/NunoFrRibeiro/personal-blog/pkg/logger"
)

func main() {
	logger.Infof("starting server")
	app.Start()
}

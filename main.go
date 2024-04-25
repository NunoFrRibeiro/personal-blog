package main

import (
	"fmt"

	app "github.com/NunoFrRibeiro/personal-blog/cmd"
	"github.com/NunoFrRibeiro/personal-blog/pkg/logger"
)

func main() {
	logger.Infof(fmt.Sprintf("starting server"))
	app.Start()
}

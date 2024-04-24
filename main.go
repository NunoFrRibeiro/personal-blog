package main

import (
	"fmt"

	app "github.com/NunoFrRibeiro/personal-go-blog/cmd"
	"github.com/NunoFrRibeiro/personal-go-blog/pkg/logger"
)

func main() {
	logger.Infof(fmt.Sprintf("starting server"))
	app.Start()
}

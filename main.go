package main

import (
	"github.com/ashep/go-apprun"

	"github.com/ashep/ujds-cli/internal/app"
)

func main() {
	apprun.Run(app.New, app.Config{})
}

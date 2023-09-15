package main

import (
	"github.com/PossibleLlama/golang-helpers/logging"
)

func main() {
	logging.InitLogger("version", "project", "service")

	logging.LogInfo("token", "message")
}

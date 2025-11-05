package main

import (
	"log/slog"

	"github.com/team-nino/iam_service/cmd"
)

func main() {
	err := cmd.RunIAM()
	if err != nil {
		slog.Error(err.Error())
	}
}

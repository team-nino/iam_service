package cmd

import "github.com/team-nino/iam_service/internal/app"

func RunIAM() error {
	application, err := app.NewApp()
	if err != nil {
		return err
	}

	return application.Run()
}

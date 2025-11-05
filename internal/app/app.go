package app

type App struct {
}

func NewApp() (*App, error) {
	return &App{}, nil
}
func (a *App) Run() error {
	return nil
}

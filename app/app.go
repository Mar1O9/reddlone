package app

type IApp interface {
	Run()
}

type App struct{
}

func (a *App) RunApp(app IApp) {
	app.Run()
}

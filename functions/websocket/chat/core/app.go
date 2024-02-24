package wschat

import (
	"blinders/packages/db"
	"blinders/packages/session"
)

var app App

type App struct {
	Session *session.Manager
	DB      *db.MongoManager
}

// init app construct an app instance for internal use
func InitApp(sessionManager *session.Manager, database *db.MongoManager) *App {
	app = App{
		Session: sessionManager,
		DB:      database,
	}

	return &app
}

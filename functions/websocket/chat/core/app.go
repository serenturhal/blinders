package wschat

import (
	"blinders/packages/db"
	"blinders/packages/session"
)

var app *App

type App struct {
	Session *session.Manager
	DB      *db.MongoManager
}

// init app construct an app instance for internal use
// is that violate stateless of functional design? app instance is used in a func
func InitApp(sm *session.Manager, dbm *db.MongoManager) *App {
	app = &App{
		Session: sm,
		DB:      dbm,
	}

	return app
}

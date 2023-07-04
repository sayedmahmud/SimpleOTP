package main

import (
	"context"

	"github.com/acheong08/SimpleOTP/internal/database"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

}

func (a *App) shutdown(ctx context.Context) {
	database.FileStore.Save()
}

// Login sets the password for the encrypted database
func (a *App) Login(password string) string {
	return database.SetPassword(password)
}

// Lists all entries in the database
func (a *App) List() []database.Entry {
	database.FileStore.Entries.Add(database.Entry{
		Name:        "Example 2",
		Description: "Example entry 2",
		URL:         "https://example.com/login",
		Secret:      "G7CQJT55IREL6NBWRJI1KLIX332YPTQN",
	})
	ret, err := database.FileStore.Entries.List()
	if err != nil {
		panic(err)
	}
	return ret
}

// Adds a new entry to the database
func (a *App) AddEntry(entry database.Entry) {
	database.FileStore.Entries.Add(entry)
}

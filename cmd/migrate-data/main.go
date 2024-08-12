package main

import (
	"context"
	"fmt"
	"log/slog"
	"sso-service/internal/storage"
	"sso-service/prisma/db"
)

const (
	appName   = "Test App"
	appSecret = "test_secret"
)

func main() {
	storage := storage.New(slog.Default())
	storage.Connect()
	prisma := storage.DB
	ctx := context.Background()

	app, err := prisma.App.FindUnique(db.App.ID.Equals(1)).Exec(ctx)
	if err != nil {
		if !db.IsErrNotFound(err) {
			panic(err)
		}

		app, err = prisma.App.CreateOne(
			db.App.Name.Set(appName),
			db.App.Secret.Set(appSecret),
		).Exec(context.Background())

		if err != nil {
			panic(err)
		}
		if app.ID != 1 {
			panic(fmt.Errorf("app id must be 1, got app id: %d", app.ID))
		}
		fmt.Println("Test app created")
	}

	if app != nil {
		fmt.Printf("Found app: %+v\n", app.InnerApp)
	}
	storage.Close()
}

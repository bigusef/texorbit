package main

import (
	"context"
	"fmt"
	"github.com/bigusef/texorbit/internal/database"
	"github.com/bigusef/texorbit/internal/router"
	"github.com/bigusef/texorbit/pkg/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
)

func main() {
	ctx := context.Background()
	setting := config.NewSetting()

	//region Database Setup
	conn, err := pgxpool.New(ctx, setting.ConnString)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	queries := database.New(conn)
	//endregion

	handler := router.New(setting, queries)

	// start application server
	log.Println("[Restfull server] Start Store application")
	log.Println(fmt.Sprintf("Serving starting on port :%s", setting.Port))

	serverPort := fmt.Sprintf(":%s", setting.Port)
	if err := http.ListenAndServe(serverPort, handler); err != nil {
		log.Fatal(err)
	}
}

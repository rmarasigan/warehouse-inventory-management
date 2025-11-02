package main

import (
	"os"

	"github.com/rmarasigan/warehouse-inventory-management/internal/app"
	"github.com/rmarasigan/warehouse-inventory-management/internal/app/db"
)

func main() {
	// Check if the argument for initialize db is set.
	if len(os.Args) > 1 && os.Args[1] == "--db=init" {
		db.Initialize()
	}

	app.StartServer()
}

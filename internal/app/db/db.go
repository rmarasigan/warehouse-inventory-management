package db

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/rmarasigan/warehouse-inventory-management/internal/utils/trail"
)

func database(ctx context.Context, db *sqlx.DB, name string) error {
	trail.Info("Initializing database...")
	_, err := db.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", name))
	if err != nil {
		trail.Error("failed to create database: %s", name)
		return err
	}

	_, err = db.ExecContext(ctx, fmt.Sprintf("USE %s;", name))
	if err != nil {
		trail.Error("failed using database: %s", name)
		return err
	}

	trail.OK("Successfully created %s database...", name)

	return nil
}

func tables(ctx context.Context, db *sqlx.DB, name, query string) error {
	_, err := db.ExecContext(ctx, query)
	if err != nil {
		trail.Warn("failed to create table: %s", name)
		return err
	}

	trail.OK("Successfully created %s table...", name)

	return nil
}

func insert[T any](db *sqlx.DB, schema T, query, tableName string) error {
	_, err := db.NamedExec(query, schema)
	if err != nil {
		trail.Warn("failed to insert item in %s table", tableName)
		return err
	}

	return nil
}

func insertList[T any](db *sqlx.DB, items []T, query, tableName string) error {
	for _, item := range items {
		err := insert(db, item, query, tableName)
		if err != nil {
			return err
		}
	}

	trail.OK("Successfully imported items to %s table...", tableName)

	return nil
}

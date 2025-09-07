package dbutils

import (
	"database/sql"
	"strings"
	"time"
)

func GetBool(input sql.NullBool) bool {
	if !input.Valid {
		return false
	}

	return input.Bool
}

func GetTime(input sql.NullTime) time.Time {
	if !input.Valid {
		return time.Time{}
	}

	return input.Time
}

func GetString(input sql.NullString) string {
	if !input.Valid {
		return ""
	}

	return input.String
}

func GetFloat(input sql.NullFloat64) float64 {
	if !input.Valid {
		return 0
	}

	return input.Float64
}

func GetAsInt(input sql.NullInt32) int {
	if !input.Valid {
		return 0
	}

	return int(input.Int32)
}

func SetString(value string) sql.NullString {
	if len(strings.TrimSpace(value)) == 0 {
		return sql.NullString{}
	}

	return sql.NullString{Valid: true, String: value}
}

func SetBool(value bool) sql.NullBool {
	return sql.NullBool{Valid: true, Bool: value}
}

func SetInt(value int32) sql.NullInt32 {
	return sql.NullInt32{Valid: true, Int32: value}
}

func SetFloat(value float64) sql.NullFloat64 {
	return sql.NullFloat64{Valid: true, Float64: value}
}

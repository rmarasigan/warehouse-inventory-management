package schema

type Currency struct {
	ID     int    `db:"id"`
	Code   string `db:"code"`
	Symbol string `db:"symbol"`
	Active bool   `db:"is_active"`
}

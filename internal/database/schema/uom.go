package schema

type UOM struct {
	ID   int    `db:"id"`
	Code string `db:"code"`
	Name string `db:"name"`
}

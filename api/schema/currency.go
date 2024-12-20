package apischema

type Currency struct {
	ID     int    `json:"id"`
	Code   string `json:"code"`
	Symbol string `json:"symbol"`
	Active bool   `json:"active"`
}

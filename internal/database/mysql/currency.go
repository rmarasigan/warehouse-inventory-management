package mysql

import (
	"github.com/rmarasigan/warehouse-inventory-management/internal/database/schema"
)

// ListCurrency returns all the currencies.
func ListCurrency() ([]schema.Currency, error) { return FetchItems[schema.Currency](CurrencyTable) }

// GetCurrency returns a currency based on ID passed.
func GetCurrency(id int) (schema.Currency, error) {
	return RetrieveItemByField[schema.Currency](CurrencyTable, "id", id)
}

// GetActiveCurrency returns the active currency.
func GetActiveCurrency() (schema.Currency, error) {
	return RetrieveItemByField[schema.Currency](CurrencyTable, "active", true)
}

// ActivateCurrency activate a currency by code.
func ActivateCurrency(code string) error {
	active, err := GetActiveCurrency()
	if err != nil {
		return err
	}

	const disableQuery = "UPDATE currency SET active = false WHERE id = ?;"
	_, err = Exec(disableQuery, active.ID)
	if err != nil {
		return err
	}

	const enableQuery = "UPDATE currency SET active = true WHERE code = ?;"
	_, err = Exec(enableQuery, code)
	if err != nil {
		return err
	}

	return nil
}

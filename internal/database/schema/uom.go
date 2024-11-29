package schema

import (
	"strings"

	apischema "github.com/rmarasigan/warehouse-inventory-management/api/schema"
)

type UOM struct {
	ID   int    `db:"id"`
	Code string `db:"code"`
	Name string `db:"name"`
}

func (u *UOM) UpdateValues(uom apischema.UOM) {
	if strings.TrimSpace(uom.Code) != "" {
		u.Code = uom.Code
	}

	if strings.TrimSpace(uom.Name) != "" {
		u.Name = uom.Name
	}
}

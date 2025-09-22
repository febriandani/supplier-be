package supplier

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type SupplierFilter struct {
	SupplierName null.String `json:"supplier_name"`
	Status       null.String `json:"status"`
	Offset       null.Int    `json:"offset"`
	Limit        null.Int    `json:"limit"`
}

type SupplierList struct {
	SupplierName string `db:"supplier_name" json:"supplier_name"`
	Logo         string `db:"logo" json:"logo"`
	Address      string `db:"address" json:"address"`
	Contact      string `db:"name" json:"contact"`
	Status       string `db:"status" json:"status"`
}

type SupplierListResponse struct {
	TotalData int64          `json:"total_data"`
	Data      []SupplierList `json:"data"`
}

type SupplierRequest struct {
	SupplierName string      `json:"supplier_name" validate:"required"`
	Logo         string      `json:"logo" validate:"required"`
	Nickname     string      `json:"nickname" validate:"required"`
	Address      []Address   `json:"address" validate:"required,dive"`
	Contacts     []Contacts  `json:"contacts" validate:"required,dive"`
	Groups       []Groups    `json:"groups" validate:"required,dive"`
	Materials    []Materials `json:"materials" validate:"required,dive"`
	Others       []Others    `json:"others" validate:"required,dive"`
}

type Supplier struct {
	SupplierID   int64     `db:"supplier_id" json:"supplier_id"`
	SupplierName string    `db:"supplier_name" json:"supplier_name"`
	Logo         string    `db:"logo" json:"logo"`
	Nickname     string    `db:"nickname" json:"nickname"`
	Status       string    `db:"status" json:"status"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
	Address      string    `db:"address" json:"address"`
}

type Address struct {
	AddressID  int64  `db:"address_id" json:"address_id"`
	SupplierID int64  `db:"supplier_id" json:"supplier_id"`
	Name       string `db:"name" json:"name"`
	Address    string `db:"address" json:"address"`
	IsMain     bool   `db:"is_main" json:"is_main"`
}

type Contacts struct {
	ContactID    int64  `db:"contact_id" json:"contact_id"`
	SupplierID   int64  `db:"supplier_id" json:"supplier_id"`
	Name         string `db:"name" json:"name"`
	JobPosition  string `db:"job_position" json:"job_position"`
	Email        string `db:"email" json:"email"`
	PhoneNumber  string `db:"phone" json:"phone_number"`
	MobileNumber string `db:"mobile" json:"mobile_number"`
	IsMain       bool   `db:"is_main" json:"is_main"`
}

type Groups struct {
	GroupID    int64  `db:"group_id" json:"group_id"`
	SupplierID int64  `db:"supplier_id" json:"supplier_id"`
	GroupName  string `db:"group_name" json:"group_name"`
	Value      string `db:"value" json:"value"`
	IsActive   bool   `db:"active" json:"is_active"`
}

type Materials struct {
	MaterialID    int64  `db:"material_id" json:"material_id"`
	SupplierID    int64  `db:"supplier_id" json:"supplier_id"`
	MaterialGroup string `db:"material_group" json:"material_group"`
	MaterialName  string `db:"material_name" json:"material_name"`
	MaterialCode  string `db:"material_code" json:"material_code"`
	IsActive      bool   `db:"active" json:"is_active"`
}

type Others struct {
	OtherID       int64  `db:"other_id" json:"other_id"`
	SupplierID    int64  `db:"supplier_id" json:"supplier_id"`
	AttributeName string `db:"attribute_name" json:"attribute_name"`
	Value         string `db:"value" json:"value"`
	IsActive      bool   `db:"active" json:"is_active"`
}

var (
	StatusSupplier = map[int]string{
		1: "Draft",
		2: "In Review",
		3: "In Assessment",
		4: "Active",
		5: "Blocked",
	}
)

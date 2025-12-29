package model

import "time"

type Part struct {
	UUID      string `json:"uuid"`
	Partinfo  PartInfo
	CreatedAt time.Time  `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type PartInfo struct {
	Name           string           `json:"name"`
	Description    string           `json:"description"`
	Price          float64          `json:"price"`
	Stock_quantity int64            `json:"stock_quantity"`
	Category       Category         `json:"category"`
	Dimensions     Dimensions       `json:"dimensions"`
	Manufacturer   Manufacturer     `json:"manufacturer"`
	Tags           []string         `json:"tags"`
	Metadata       map[string]Value `json:"metadata"`
}

type PartInfoUpdate struct {
	Name          *string
	Description   *string
	Price         *float64
	StockQuantity *int64
	Category      *Category
	Dimensions    *Dimensions
	Manufacturer  *Manufacturer
}

type Category int32

const (
	CategoryUnknown  Category = 0
	CategoryEngine   Category = 1
	CategoryFuel     Category = 2
	CategoryPorthole Category = 3
	CategoryWing     Category = 4
)

type Dimensions struct {
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Weight float64 `json:"weight"`
}

type Manufacturer struct {
	Name    string `json:"name"`
	Country string `json:"country"`
	Website string `json:"website"`
}

type Value struct {
	StringValue *string  `json:"string_value,omitempty"`
	Int64Value  *int64   `json:"int64_value,omitempty"`
	DoubleValue *float64 `json:"double_value,omitempty"`
	BoolValue   *bool    `json:"bool_value,omitempty"`
}

type PartsFilter struct {
	Uuids                 *[]string
	Names                 *[]string
	Categories            *[]Category
	ManufacturerCountries *[]string
	Tags                  *[]string
}

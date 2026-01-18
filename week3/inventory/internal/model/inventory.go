package model

import (
	"time"
)

// Inventory - основной объект
type Inventory struct {
	UUID          string        `json:"uuid"`
	InventoryInfo InventoryInfo `json:"inventory_info"`
	CreatedAt     time.Time     `json:"created_at,omitempty"`
	UpdatedAt     *time.Time    `json:"updated_at,omitempty"`
	DeletedAt     *time.Time    `json:"deleted_at,omitempty"`
}

// InventoryInfo - информация о конкретном инвентаре
type InventoryInfo struct {
	Name          string           `json:"name"`
	Description   string           `json:"description"`
	Price         float64          `json:"price"`
	StockQuantity int64            `json:"stock_quantity"`
	Category      Category         `json:"category"`
	Dimensions    Dimensions       `json:"dimensions"`
	Manufacturer  Manufacturer     `json:"manufacturer"`
	Tags          []string         `json:"tags"`
	Metadata      map[string]Value `json:"metadata"`
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

// Category - перечисление категорий
type Category int32

const (
	CategoryUnspecified Category = 0
	CategoryEngine      Category = 1
	CategoryFuel        Category = 2
	CategoryPorthole    Category = 3
	CategoryWing        Category = 4
)

// Dimensions - размеры и вес
type Dimensions struct {
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Weight float64 `json:"weight"`
}

// Manufacturer - производитель
type Manufacturer struct {
	Name    string `json:"name"`
	Country string `json:"country"`
	Website string `json:"website"`
}

// Value - значение для метаданных (oneof)
type Value struct {
	StringValue *string  `json:"string_value,omitempty"`
	Int64Value  *int64   `json:"int64_value,omitempty"`
	DoubleValue *float64 `json:"double_value,omitempty"`
	BoolValue   *bool    `json:"bool_value,omitempty"`
}

type ListParts struct {
	UUID                  *[]string
	Names                 *[]string
	Categories            *[]Category
	ManufacturerCountries *[]string
	Tags                  *[]string
}

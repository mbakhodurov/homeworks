package converter

import (
	"github.com/mbakhodurov/homeworks/week6/inventory/internal/model"
	inventory_v1 "github.com/mbakhodurov/homeworks/week6/shared/pkg/proto/inventory/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func InventortListPartProtoToModel(from *inventory_v1.PartsFilter) model.InventoryFilter {
	if from == nil {
		return model.InventoryFilter{}
	}
	partsInfo := model.InventoryFilter{}
	if len(from.GetUuid()) > 0 {
		uuids := from.GetUuid()
		partsInfo.UUID = &uuids
	}
	if len(from.GetNames()) > 0 {
		names := from.GetNames()
		partsInfo.Names = &names
	}

	if len(from.GetCategories()) > 0 {
		categories := make([]model.Category, 0, len(from.GetCategories()))
		for _, c := range from.GetCategories() {
			categories = append(categories, model.Category(c))
		}
		partsInfo.Categories = &categories
	}

	if len(from.GetManufacturerCountries()) > 0 {
		countries := from.GetManufacturerCountries()
		partsInfo.ManufacturerCountries = &countries
	}

	if len(from.GetTags()) > 0 {
		tags := from.GetTags()
		partsInfo.Tags = &tags
	}

	return partsInfo
}

func metadataProtoToModel(meta map[string]*inventory_v1.Value) map[string]model.Value {
	if meta == nil {
		return nil
	}
	res := make(map[string]model.Value, len(meta))
	for k, v := range meta {
		if v == nil {
			continue
		}

		mv := model.Value{}
		switch val := v.Kind.(type) {
		case *inventory_v1.Value_StringValue:
			mv.StringValue = &val.StringValue
		case *inventory_v1.Value_Int64Value:
			mv.Int64Value = &val.Int64Value
		case *inventory_v1.Value_DoubleValue:
			mv.DoubleValue = &val.DoubleValue
		case *inventory_v1.Value_BoolValue:
			mv.BoolValue = &val.BoolValue
		}

		res[k] = mv
	}
	return res
}

func InventoryInfoProtoToModel(from *inventory_v1.InventoryInfo) model.InventoryInfo {
	return model.InventoryInfo{
		Name:           from.Name,
		Description:    from.Description,
		Price:          from.Price,
		Stock_quantity: from.StockQuantity,
		Category:       model.Category(from.Category),
		Dimensions: model.Dimensions{
			Length: from.Dimensions.Length,
			Width:  from.Dimensions.Width,
			Height: from.Dimensions.Height,
			Weight: from.Dimensions.Weight,
		},
		Manufacturer: model.Manufacturer{
			Name:    from.Manufacturer.Name,
			Country: from.Manufacturer.Country,
			Website: from.Manufacturer.Website,
		},
		Tags:     from.Tags,
		Metadata: metadataProtoToModel(from.Metadata),
	}
}

func UpdateInventoryProtoToModel(from *inventory_v1.InventoryUpdateInfo) model.InventoryInfoUpdate {
	var res model.InventoryInfoUpdate

	if from.Name != nil {
		tmp := from.Name.Value
		res.Name = &tmp
	}

	if from.Description != nil {
		tmp := from.Description.Value
		res.Description = &tmp
	}

	if from.Price != nil {
		tmp := from.Price.Value
		res.Price = &tmp
	}

	if from.StockQuantity != nil {
		tmp := from.StockQuantity.Value
		res.StockQuantity = &tmp
	}

	if from.Category != inventory_v1.Category_CATEGORY_UNSPECIFIED {
		tmp := model.Category(from.Category)
		res.Category = &tmp
	}

	if from.Dimensions != nil {
		res.Dimensions = &model.Dimensions{
			Length: from.Dimensions.Length,
			Width:  from.Dimensions.Width,
			Height: from.Dimensions.Height,
			Weight: from.Dimensions.Weight,
		}
	}

	if from.Manufacturer != nil {
		res.Manufacturer = &model.Manufacturer{
			Name:    from.Manufacturer.Name,
			Country: from.Manufacturer.Country,
			Website: from.Manufacturer.Website,
		}
	}

	return res
}

func InventoryInfoModelToProto(from model.InventoryInfo) *inventory_v1.InventoryInfo {
	return &inventory_v1.InventoryInfo{
		Name:          from.Name,
		Description:   from.Description,
		Price:         from.Price,
		StockQuantity: from.Stock_quantity,
		Category:      inventory_v1.Category(from.Category),
		Dimensions: &inventory_v1.Dimensions{
			Length: from.Dimensions.Length,
			Width:  from.Dimensions.Width,
			Height: from.Dimensions.Height,
			Weight: from.Dimensions.Weight,
		},
		Manufacturer: &inventory_v1.Manufacturer{
			Name:    from.Manufacturer.Name,
			Country: from.Manufacturer.Country,
			Website: from.Manufacturer.Website,
		},
		Tags:     from.Tags,
		Metadata: metadataModelToProto(from.Metadata),
	}
}
func metadataModelToProto(meta map[string]model.Value) map[string]*inventory_v1.Value {
	if meta == nil {
		return nil
	}

	res := make(map[string]*inventory_v1.Value, len(meta))
	for k, v := range meta {
		pv := &inventory_v1.Value{}

		switch {
		case v.StringValue != nil:
			pv.Kind = &inventory_v1.Value_StringValue{
				StringValue: *v.StringValue,
			}
		case v.Int64Value != nil:
			pv.Kind = &inventory_v1.Value_Int64Value{
				Int64Value: *v.Int64Value,
			}
		case v.DoubleValue != nil:
			pv.Kind = &inventory_v1.Value_DoubleValue{
				DoubleValue: *v.DoubleValue,
			}
		case v.BoolValue != nil:
			pv.Kind = &inventory_v1.Value_BoolValue{
				BoolValue: *v.BoolValue,
			}
		default:
			continue // пустое Value — пропускаем
		}

		res[k] = pv
	}

	return res
}

func InventoryModelToProto(from model.Inventory) *inventory_v1.Inventory {
	var updatedAt *timestamppb.Timestamp
	if from.UpdatedAt != nil {
		updatedAt = timestamppb.New(*from.UpdatedAt)
	}

	var deletedAt *timestamppb.Timestamp
	if from.DeletedAt != nil {
		deletedAt = timestamppb.New(*from.DeletedAt)
	}

	return &inventory_v1.Inventory{
		Uuid:          from.UUID,
		InventoryInfo: InventoryInfoModelToProto(from.InventoryInfo),
		CreatedAt:     timestamppb.New(from.CreatedAt),
		UpdatedAt:     updatedAt,
		DeletedAt:     deletedAt,
	}
}

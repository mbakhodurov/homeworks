package converter

import (
	"github.com/mbakhodurov/homeworks/week3/inventory/internal/model"
	inventory_v1 "github.com/mbakhodurov/homeworks/week3/shared/pkg/proto/inventory/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func InventoryModelToProto(inventory model.Inventory) *inventory_v1.Inventory {

	var updatedAt *timestamppb.Timestamp
	if inventory.UpdatedAt != nil {
		updatedAt = timestamppb.New(*inventory.UpdatedAt)
	}

	var deletedAt *timestamppb.Timestamp
	if inventory.DeletedAt != nil {
		deletedAt = timestamppb.New(*inventory.DeletedAt)
	}

	return &inventory_v1.Inventory{
		Uuid:          inventory.UUID,
		InventoryInfo: InventoryInfoModelToProto(inventory.InventoryInfo),
		CreatedAt:     timestamppb.New(inventory.CreatedAt),
		UpdatedAt:     updatedAt,
		DeletedAt:     deletedAt,
	}
}

func InventoryInfoModelToProto(inventorypart model.InventoryInfo) *inventory_v1.InventoryInfo {
	return &inventory_v1.InventoryInfo{
		Name:          inventorypart.Name,
		Description:   inventorypart.Description,
		Price:         inventorypart.Price,
		StockQuantity: inventorypart.StockQuantity,
		Category:      inventory_v1.Category(inventorypart.Category),
		Dimensions: &inventory_v1.Dimensions{
			Length: inventorypart.Dimensions.Length,
			Width:  inventorypart.Dimensions.Width,
			Height: inventorypart.Dimensions.Height,
			Weight: inventorypart.Dimensions.Weight,
		},
		Manufacturer: &inventory_v1.Manufacturer{
			Name:    inventorypart.Manufacturer.Name,
			Country: inventorypart.Manufacturer.Country,
			Website: inventorypart.Manufacturer.Website,
		},
		Tags:     inventorypart.Tags,
		Metadata: metadataModelToProto(inventorypart.Metadata),
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

func InventoryInfoProtoToModel(partInfo *inventory_v1.InventoryInfo) model.InventoryInfo {
	if partInfo == nil {
		return model.InventoryInfo{}
	}
	return model.InventoryInfo{
		Name:          partInfo.Name,
		Description:   partInfo.Description,
		Price:         partInfo.Price,
		StockQuantity: partInfo.StockQuantity,
		Category:      model.Category(partInfo.Category),
		Dimensions: model.Dimensions{
			Length: partInfo.Dimensions.Length,
			Width:  partInfo.Dimensions.Width,
			Height: partInfo.Dimensions.Height,
			Weight: partInfo.Dimensions.Weight,
		},
		Manufacturer: model.Manufacturer{
			Name:    partInfo.Manufacturer.Name,
			Country: partInfo.Manufacturer.Country,
			Website: partInfo.Manufacturer.Website,
		},
		Tags:     partInfo.Tags,
		Metadata: metadataProtoToModel(partInfo.Metadata),
	}
}

func InventoryUpdateInfoProtoToModel(updatePartInfo *inventory_v1.InventoryUpdateInfo) model.PartInfoUpdate {
	var res model.PartInfoUpdate

	if updatePartInfo.Name != nil {
		tmp := updatePartInfo.Name.Value
		res.Name = &tmp
	}

	if updatePartInfo.Description != nil {
		tmp := updatePartInfo.Description.Value
		res.Description = &tmp
	}

	if updatePartInfo.Price != nil {
		tmp := updatePartInfo.Price.Value
		res.Price = &tmp
	}

	if updatePartInfo.StockQuantity != nil {
		tmp := updatePartInfo.StockQuantity.Value
		res.StockQuantity = &tmp
	}

	if updatePartInfo.Category != inventory_v1.Category_CATEGORY_UNSPECIFIED {
		tmp := model.Category(updatePartInfo.Category)
		res.Category = &tmp
	}

	if updatePartInfo.Dimensions != nil {
		res.Dimensions = &model.Dimensions{
			Length: updatePartInfo.Dimensions.Length,
			Width:  updatePartInfo.Dimensions.Width,
			Height: updatePartInfo.Dimensions.Height,
			Weight: updatePartInfo.Dimensions.Weight,
		}
	}

	if updatePartInfo.Manufacturer != nil {
		res.Manufacturer = &model.Manufacturer{
			Name:    updatePartInfo.Manufacturer.Name,
			Country: updatePartInfo.Manufacturer.Country,
			Website: updatePartInfo.Manufacturer.Website,
		}
	}

	return res
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

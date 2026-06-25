package converter

import (
	"github.com/mbakhodurov/homeworks/week7/inventory/internal/model"
	repoModel "github.com/mbakhodurov/homeworks/week7/inventory/internal/repository/model"
)

func InventoryInfoModelToRepoInfoModel(from model.InventoryInfo) repoModel.InventoryInfo {
	var meta map[string]repoModel.Value
	if from.Metadata != nil {
		meta = make(map[string]repoModel.Value, len(from.Metadata))
		for k, v := range from.Metadata {
			meta[k] = repoModel.Value{
				StringValue: v.StringValue,
				Int64Value:  v.Int64Value,
				DoubleValue: v.DoubleValue,
				BoolValue:   v.BoolValue,
			}
		}
	}

	return repoModel.InventoryInfo{
		Name:           from.Name,
		Description:    from.Description,
		Price:          from.Price,
		Stock_quantity: from.Stock_quantity,
		Category:       repoModel.Category(from.Category),
		Dimensions: repoModel.Dimensions{
			Length: from.Dimensions.Length,
			Width:  from.Dimensions.Width,
			Height: from.Dimensions.Height,
			Weight: from.Dimensions.Weight,
		},
		Manufacturer: repoModel.Manufacturer{
			Name:    from.Manufacturer.Name,
			Country: from.Manufacturer.Country,
			Website: from.Manufacturer.Website,
		},
		Tags:     from.Tags,
		Metadata: meta,
	}
}

func InventoryModelToRepoModel(from model.Inventory) repoModel.Inventory {
	return repoModel.Inventory{
		UUID:          from.UUID,
		InventoryInfo: InventoryInfoModelToRepoInfoModel(from.InventoryInfo),
		CreatedAt:     from.CreatedAt,
		UpdatedAt:     from.UpdatedAt,
		DeletedAt:     from.DeletedAt,
	}
}

func InventoryInfoRepoModelToInventoryInfoModel(inventoryInfoRepoModel repoModel.InventoryInfo) model.InventoryInfo {
	var meta map[string]model.Value
	if inventoryInfoRepoModel.Metadata != nil {
		meta = make(map[string]model.Value, len(inventoryInfoRepoModel.Metadata))
		for k, v := range inventoryInfoRepoModel.Metadata {
			meta[k] = model.Value{
				StringValue: v.StringValue,
				Int64Value:  v.Int64Value,
				DoubleValue: v.DoubleValue,
				BoolValue:   v.BoolValue,
			}
		}
	}

	return model.InventoryInfo{
		Name:           inventoryInfoRepoModel.Name,
		Description:    inventoryInfoRepoModel.Description,
		Price:          inventoryInfoRepoModel.Price,
		Stock_quantity: inventoryInfoRepoModel.Stock_quantity,
		Category:       model.Category(inventoryInfoRepoModel.Category),
		Dimensions: model.Dimensions{
			Length: inventoryInfoRepoModel.Dimensions.Length,
			Width:  inventoryInfoRepoModel.Dimensions.Width,
			Height: inventoryInfoRepoModel.Dimensions.Height,
			Weight: inventoryInfoRepoModel.Dimensions.Weight,
		},
		Manufacturer: model.Manufacturer{
			Name:    inventoryInfoRepoModel.Manufacturer.Name,
			Country: inventoryInfoRepoModel.Manufacturer.Country,
			Website: inventoryInfoRepoModel.Manufacturer.Website,
		},
		Tags: inventoryInfoRepoModel.Tags,
	}
}

func InventoryRepoModelToModel(from repoModel.Inventory) model.Inventory {
	return model.Inventory{
		UUID:          from.UUID,
		InventoryInfo: InventoryInfoRepoModelToInventoryInfoModel(from.InventoryInfo),
		CreatedAt:     from.CreatedAt,
		UpdatedAt:     from.UpdatedAt,
		DeletedAt:     from.DeletedAt,
	}
}

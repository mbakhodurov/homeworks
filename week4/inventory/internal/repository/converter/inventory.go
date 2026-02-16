package converter

import (
	"github.com/mbakhodurov/homeworks/week4/inventory/internal/model"
	repoModel "github.com/mbakhodurov/homeworks/week4/inventory/internal/repository/model"
)

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
		Name:          inventoryInfoRepoModel.Name,
		Description:   inventoryInfoRepoModel.Description,
		Price:         inventoryInfoRepoModel.Price,
		StockQuantity: inventoryInfoRepoModel.StockQuantity,
		Category:      model.Category(inventoryInfoRepoModel.Category),
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

func InventoryRepoModelToModel(inventoryRepoModel repoModel.Inventory) model.Inventory {
	return model.Inventory{
		UUID:          inventoryRepoModel.UUID,
		InventoryInfo: InventoryInfoRepoModelToInventoryInfoModel(inventoryRepoModel.InventoryInfo),
		CreatedAt:     inventoryRepoModel.CreatedAt,
		UpdatedAt:     inventoryRepoModel.UpdatedAt,
		DeletedAt:     inventoryRepoModel.DeletedAt,
	}
}

func InventoryInfoModelToRepoModel(inventoryInfoModel model.InventoryInfo) repoModel.InventoryInfo {
	var meta map[string]repoModel.Value
	if inventoryInfoModel.Metadata != nil {
		meta = make(map[string]repoModel.Value, len(inventoryInfoModel.Metadata))
		for k, v := range inventoryInfoModel.Metadata {
			meta[k] = repoModel.Value{
				StringValue: v.StringValue,
				Int64Value:  v.Int64Value,
				DoubleValue: v.DoubleValue,
				BoolValue:   v.BoolValue,
			}
		}
	}

	return repoModel.InventoryInfo{
		Name:          inventoryInfoModel.Name,
		Description:   inventoryInfoModel.Description,
		Price:         inventoryInfoModel.Price,
		StockQuantity: inventoryInfoModel.StockQuantity,
		Category:      repoModel.Category(inventoryInfoModel.Category),
		Dimensions: repoModel.Dimensions{
			Length: inventoryInfoModel.Dimensions.Length,
			Width:  inventoryInfoModel.Dimensions.Width,
			Height: inventoryInfoModel.Dimensions.Height,
			Weight: inventoryInfoModel.Dimensions.Weight,
		},
		Manufacturer: repoModel.Manufacturer{
			Name:    inventoryInfoModel.Manufacturer.Name,
			Country: inventoryInfoModel.Manufacturer.Country,
			Website: inventoryInfoModel.Manufacturer.Website,
		},
		Tags:     inventoryInfoModel.Tags,
		Metadata: meta,
	}
}

func InventoryModelToRepoModel(from model.Inventory) repoModel.Inventory {
	return repoModel.Inventory{
		UUID:          from.UUID,
		InventoryInfo: InventoryInfoModelToRepoModel(from.InventoryInfo),
		CreatedAt:     from.CreatedAt,
		UpdatedAt:     from.UpdatedAt,
		DeletedAt:     from.DeletedAt,
	}
}

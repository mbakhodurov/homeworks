package converter

import (
	"github.com/mbakhodurov/homeworks/week2/inventory/internal/model"
	repoModel "github.com/mbakhodurov/homeworks/week2/inventory/internal/repository/model"
)

func InventoryInfoToModel(info repoModel.PartInfo) model.PartInfo {
	var meta map[string]model.Value
	if info.Metadata != nil {
		meta = make(map[string]model.Value, len(info.Metadata))
		for k, v := range info.Metadata {
			meta[k] = model.Value{
				StringValue: v.StringValue,
				Int64Value:  v.Int64Value,
				DoubleValue: v.DoubleValue,
				BoolValue:   v.BoolValue,
			}
		}
	}

	return model.PartInfo{
		Name:           info.Name,
		Description:    info.Description,
		Price:          info.Price,
		Stock_quantity: info.Stock_quantity,
		Category:       model.Category(info.Category),
		Dimensions: model.Dimensions{
			Length: info.Dimensions.Length,
			Width:  info.Dimensions.Width,
			Height: info.Dimensions.Height,
			Weight: info.Dimensions.Weight,
		},
		Manufacturer: model.Manufacturer{
			Name:    info.Manufacturer.Name,
			Country: info.Manufacturer.Country,
			Website: info.Manufacturer.Website,
		},
		Tags:     info.Tags,
		Metadata: meta,
	}
}
func InventoryPartToModel(info repoModel.Part) model.Part {
	return model.Part{
		UUID:      info.UUID,
		Partinfo:  InventoryInfoToModel(info.Partinfo),
		CreatedAt: info.CreatedAt,
		UpdatedAt: info.UpdatedAt,
		DeletedAt: info.DeletedAt,
	}
}

func InventoryInfoToRepoModel(info model.PartInfo) repoModel.PartInfo {
	var meta map[string]repoModel.Value
	if info.Metadata != nil {
		meta = make(map[string]repoModel.Value, len(info.Metadata))
		for k, v := range info.Metadata {
			meta[k] = repoModel.Value{
				StringValue: v.StringValue,
				Int64Value:  v.Int64Value,
				DoubleValue: v.DoubleValue,
				BoolValue:   v.BoolValue,
			}
		}
	}

	return repoModel.PartInfo{
		Name:           info.Name,
		Description:    info.Description,
		Price:          info.Price,
		Stock_quantity: info.Stock_quantity,
		Category:       repoModel.Category(info.Category),
		Dimensions: repoModel.Dimensions{
			Length: info.Dimensions.Length,
			Width:  info.Dimensions.Width,
			Height: info.Dimensions.Height,
			Weight: info.Dimensions.Weight,
		},
		Manufacturer: repoModel.Manufacturer{
			Name:    info.Manufacturer.Name,
			Country: info.Manufacturer.Country,
			Website: info.Manufacturer.Website,
		},
		Tags:     info.Tags,
		Metadata: meta,
	}
}

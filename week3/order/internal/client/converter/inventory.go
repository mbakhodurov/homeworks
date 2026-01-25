package converter

import (
	"time"

	"github.com/mbakhodurov/homeworks/week3/order/internal/model"
	inventory_v1 "github.com/mbakhodurov/homeworks/week3/shared/pkg/proto/inventory/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ModelFilterToProtoFilter(from model.InventoryFilter) *inventory_v1.PartsFilter {

	var categories []inventory_v1.Category
	if from.Categories != nil {
		categories = make([]inventory_v1.Category, 0, len(*from.Categories))
		for _, category := range *from.Categories {
			categories = append(categories, inventory_v1.Category(int32(category))) //nolint:gosec
		}
	}

	var uuids, names, manufacturerCountries, tags []string
	if from.UUID != nil {
		uuids = *from.UUID
	}

	if from.Names != nil {
		names = *from.Names
	}
	if from.ManufacturerCountries != nil {
		manufacturerCountries = *from.ManufacturerCountries
	}
	if from.Tags != nil {
		tags = *from.Tags
	}

	return &inventory_v1.PartsFilter{
		Uuids:                 uuids,
		Names:                 names,
		Categories:            categories,
		ManufacturerCountries: manufacturerCountries,
		Tags:                  tags,
	}
}

func DimensionsToModel(d *inventory_v1.Dimensions) model.Dimensions {
	if d == nil {
		return model.Dimensions{}
	}
	return model.Dimensions{
		Length: d.Length,
		Width:  d.Width,
		Height: d.Height,
		Weight: d.Weight,
	}
}

func manufacturerToModel(m *inventory_v1.Manufacturer) model.Manufacturer {
	if m == nil {
		return model.Manufacturer{}
	}

	return model.Manufacturer{
		Name:    m.GetName(),
		Country: m.GetCountry(),
		Website: m.GetWebsite(),
	}
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
		Name:          from.Name,
		Description:   from.Description,
		Price:         from.Price,
		StockQuantity: from.StockQuantity,
		Category:      model.Category(from.Category),
		Dimensions:    DimensionsToModel(from.Dimensions),
		Manufacturer:  manufacturerToModel(from.Manufacturer),
		Tags:          from.Tags,
		Metadata:      metadataProtoToModel(from.Metadata),
	}
}

func timestampToPtr(ts *timestamppb.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}

	t := ts.AsTime()
	return &t
}

func ResponseProtoToInventoryModel(from *inventory_v1.ListPartInventoryResponse) []model.Inventory {
	if from == nil {
		return []model.Inventory{}
	}

	var resPart []model.Inventory
	for _, part := range from.Inventory {
		if part == nil {
			continue
		}

		part := model.Inventory{
			UUID:          part.Uuid,
			InventoryInfo: InventoryInfoProtoToModel(part.InventoryInfo),
			CreatedAt:     part.CreatedAt.AsTime(),
			UpdatedAt:     timestampToPtr(part.UpdatedAt),
			DeletedAt:     timestampToPtr(part.UpdatedAt),
		}
		resPart = append(resPart, part)
	}
	return resPart
}

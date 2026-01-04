package converter

import (
	"time"

	"github.com/mbakhodurov/homeworks/week2/order/internal/model"
	inventory_v1 "github.com/mbakhodurov/homeworks/week2/shared/pkg/proto/inventory/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func PartsFilterToProto(filter model.PartsFilter) *inventory_v1.PartsFilter {
	var categories []inventory_v1.Category
	if filter.Categories != nil {
		categories = make([]inventory_v1.Category, 0, len(*filter.Categories))
		for _, category := range *filter.Categories {
			categories = append(categories, inventory_v1.Category(int32(category))) //nolint:gosec
		}
	}

	var uuids, names, manufacturerCountries, tags []string
	if filter.Uuids != nil {
		uuids = *filter.Uuids
	}

	if filter.Names != nil {
		names = *filter.Names
	}
	if filter.ManufacturerCountries != nil {
		manufacturerCountries = *filter.ManufacturerCountries
	}
	if filter.Tags != nil {
		tags = *filter.Tags
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

func timestampToPtr(ts *timestamppb.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}

	t := ts.AsTime()
	return &t
}

func PartInfoToModel(info *inventory_v1.PartInfo) model.PartInfo {
	if info == nil {
		return model.PartInfo{}
	}
	return model.PartInfo{
		Name:           info.Name,
		Description:    info.Description,
		Price:          info.Price,
		Stock_quantity: info.StockQuantity,
		Category:       model.Category(info.Category),
		Dimensions:     DimensionsToModel(info.Dimensions),
		Manufacturer:   manufacturerToModel(info.Manufacturer),
		Tags:           info.Tags,
		Metadata:       metadataProtoToModel(info.Metadata),
	}
}

func PartListToModel(r *inventory_v1.ListPartsResponse) []model.Part {
	if r == nil {
		return []model.Part{}
	}

	var resPart []model.Part
	for _, part := range r.Part {
		if part == nil {
			continue
		}

		part := model.Part{
			UUID:      part.Uuid,
			Partinfo:  PartInfoToModel(part.Info),
			CreatedAt: part.CreatedAt.AsTime(),
			UpdatedAt: timestampToPtr(part.UpdatedAt),
			DeletedAt: timestampToPtr(part.UpdatedAt),
		}
		resPart = append(resPart, part)
	}
	return resPart
}

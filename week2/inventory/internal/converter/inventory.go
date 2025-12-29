package converter

import (
	"github.com/mbakhodurov/homeworks/week2/inventory/internal/model"
	inventory_v1 "github.com/mbakhodurov/homeworks/week2/shared/pkg/proto/inventory/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

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

func InventoryPartInfoProtoToModel(info *inventory_v1.PartInfo) model.PartInfo {
	if info == nil {
		return model.PartInfo{}
	}

	return model.PartInfo{
		Name:           info.Name,
		Description:    info.Description,
		Price:          info.Price,
		Stock_quantity: info.StockQuantity,
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
		Metadata: metadataProtoToModel(info.Metadata),
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

func InvertoryPartToProto(info model.Part) *inventory_v1.Part {
	var updatedAt *timestamppb.Timestamp
	if info.UpdatedAt != nil {
		updatedAt = timestamppb.New(*info.UpdatedAt)
	}

	var deletedAt *timestamppb.Timestamp
	if info.DeletedAt != nil {
		deletedAt = timestamppb.New(*info.DeletedAt)
	}

	return &inventory_v1.Part{
		Uuid:      info.UUID,
		Info:      InvertoryPartInfoToProto(info.Partinfo),
		CreatedAt: timestamppb.New(info.CreatedAt),
		UpdatedAt: updatedAt,
		DeletedAt: deletedAt,
	}
}

func InvertoryPartInfoToProto(info model.PartInfo) *inventory_v1.PartInfo {
	return &inventory_v1.PartInfo{
		Name:          info.Name,
		Description:   info.Description,
		Price:         info.Price,
		StockQuantity: info.Stock_quantity,
		Category:      inventory_v1.Category(info.Category),
		Dimensions: &inventory_v1.Dimensions{
			Length: info.Dimensions.Length,
			Width:  info.Dimensions.Width,
			Height: info.Dimensions.Height,
			Weight: info.Dimensions.Weight,
		},
		Manufacturer: &inventory_v1.Manufacturer{
			Name:    info.Manufacturer.Name,
			Country: info.Manufacturer.Country,
			Website: info.Manufacturer.Website,
		},
		Tags:     info.Tags,
		Metadata: metadataModelToProto(info.Metadata),
	}
}

func InventoryUpdateInfoToModel(info *inventory_v1.InventoryUpdateInfo) model.PartInfoUpdate {
	var res model.PartInfoUpdate

	if info.Name != nil {
		tmp := info.Name.Value
		res.Name = &tmp
	}

	if info.Description != nil {
		tmp := info.Description.Value
		res.Description = &tmp
	}

	if info.Price != nil {
		tmp := info.Price.Value
		res.Price = &tmp
	}

	if info.StockQuantity != nil {
		tmp := info.StockQuantity.Value
		res.StockQuantity = &tmp
	}

	if info.Category != inventory_v1.Category_UNKNOWN {
		tmp := model.Category(info.Category)
		res.Category = &tmp
	}

	if info.Dimensions != nil {
		res.Dimensions = &model.Dimensions{
			Length: info.Dimensions.Length,
			Width:  info.Dimensions.Width,
			Height: info.Dimensions.Height,
			Weight: info.Dimensions.Weight,
		}
	}

	if info.Manufacturer != nil {
		res.Manufacturer = &model.Manufacturer{
			Name:    info.Manufacturer.Name,
			Country: info.Manufacturer.Country,
			Website: info.Manufacturer.Website,
		}
	}

	return res
}

func ProtoToPartsFilter(info *inventory_v1.PartsFilter) model.PartsFilter {
	var res model.PartsFilter
	if info == nil {
		return model.PartsFilter{}
	}

	if len(info.Uuids) > 0 {
		tmp := make([]string, len(info.Uuids))
		copy(tmp, info.Uuids)
		res.Uuids = &tmp
	}

	if len(info.Names) > 0 {
		tmp := make([]string, len(info.Names))
		copy(tmp, info.Names)
		res.Names = &tmp
	}

	if len(info.Categories) > 0 {
		tmp := make([]model.Category, len(info.Categories))
		for i, c := range info.Categories {
			tmp[i] = model.Category(c)
		}
		res.Categories = &tmp
	}

	if len(info.ManufacturerCountries) > 0 {
		tmp := make([]string, len(info.ManufacturerCountries))
		copy(tmp, info.ManufacturerCountries)
		res.ManufacturerCountries = &tmp
	}

	if len(info.Tags) > 0 {
		tmp := make([]string, len(info.Tags))
		copy(tmp, info.Tags)
		res.Tags = &tmp
	}

	return res
}

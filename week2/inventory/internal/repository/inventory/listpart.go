package inventory

import (
	"context"

	"github.com/mbakhodurov/homeworks/week2/inventory/internal/model"
	repoconverter "github.com/mbakhodurov/homeworks/week2/inventory/internal/repository/converter"
)

func (r *repository) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	result := make([]model.Part, 0, len(r.data))
	for _, part := range r.data {
		repoPart := repoconverter.InventoryPartToModel(part)
		if matchPart(&repoPart, &filter) {
			result = append(result, repoPart)
		}
	}

	return result, nil
}

func matchPart(part *model.Part, filter *model.PartsFilter) bool {
	if filter == nil {
		return true
	}

	info := part.Partinfo
	// UUIDs
	if filter.Uuids != nil && len(*filter.Uuids) > 0 {
		found := false
		for _, u := range *filter.Uuids {
			if part.UUID == u {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Names
	if filter.Names != nil && len(*filter.Names) > 0 {
		found := false
		for _, n := range *filter.Names {
			if part.Partinfo.Name == n {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Categories
	if filter.Categories != nil && len(*filter.Categories) > 0 {
		found := false
		for _, c := range *filter.Categories {
			if info.Category == c {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Manufacturer countries
	if filter.ManufacturerCountries != nil && len(*filter.ManufacturerCountries) > 0 {
		found := false
		for _, country := range *filter.ManufacturerCountries {
			if info.Manufacturer.Country == country {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Tags
	if filter.Tags != nil && len(*filter.Tags) > 0 {
		found := false
		for _, tag := range *filter.Tags {
			for _, pt := range info.Tags {
				if pt == tag {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

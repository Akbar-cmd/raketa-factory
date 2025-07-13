package part

import (
	"context"
	"strings"

	"github.com/Akbar-cmd/raketa-factory/inventory/internal/model"
	"github.com/Akbar-cmd/raketa-factory/inventory/internal/repository/converter"
	repoModel "github.com/Akbar-cmd/raketa-factory/inventory/internal/repository/model"
)

func (r *repository) ListParts(_ context.Context, filter model.PartsFilter) ([]model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var parts []model.Part

	repoFilter := converter.FilterToRepoModel(filter)

	for _, repoPart := range r.data {
		if matchesFilter(repoPart, repoFilter) {
			parts = append(parts, converter.PartToModel(repoPart))
		}
	}

	return parts, nil
}

// matchesFilter фильтрует детали
func matchesFilter(part repoModel.Part, filter repoModel.PartsFilter) bool {

	return matchUUID(part, filter) &&
		matchName(part, filter) &&
		matchCategory(part, filter) &&
		matchManufacturerCountry(part, filter) &&
		matchTags(part, filter)
}

// Фильтрация по UUID
func matchUUID(part repoModel.Part, filter repoModel.PartsFilter) bool {
	if len(filter.Uuids) == 0 {
		return true
	}
	for _, u := range filter.Uuids {
		if part.Uuid == u {
			return true
		}
	}
	return false
}

// Фильтрация по имени
func matchName(part repoModel.Part, filter repoModel.PartsFilter) bool {
	if len(filter.Names) == 0 {
		return true
	}
	lower := strings.ToLower(part.Name)
	for _, name := range filter.Names {
		if strings.Contains(lower, strings.ToLower(name)) {
			return true
		}
	}
	return false
}

// Фильтрация по категориям
func matchCategory(part repoModel.Part, filter repoModel.PartsFilter) bool {
	if len(filter.Categories) == 0 {
		return true
	}
	for _, cat := range filter.Categories {
		if part.Category == cat {
			return true
		}
	}
	return false
}

// Фильтрация по странам производителям
func matchManufacturerCountry(part repoModel.Part, filter repoModel.PartsFilter) bool {
	if len(filter.ManufacturerCountries) == 0 {
		return true
	}

	for _, country := range filter.ManufacturerCountries {
		if strings.EqualFold(part.Manufacturer.Country, country) {
			return true
		}
	}
	return false
}

// Фильтрация по тегам
func matchTags(part repoModel.Part, filter repoModel.PartsFilter) bool {
	if len(filter.Tags) == 0 {
		return true
	}
	for _, pTag := range part.Tags {
		for _, fTag := range filter.Tags {
			if pTag == fTag {
				return true
			}
		}
	}
	return false
}

// copyDimensions создает копию объекта Dimensions
func copyDimensions(src repoModel.Dimensions) repoModel.Dimensions {

	return repoModel.Dimensions{
		Length: src.Length,
		Width:  src.Width,
		Height: src.Height,
		Weight: src.Weight,
	}
}

// copyManufacturer создает копию объекта Manufacturer
func copyManufacturer(src repoModel.Manufacturer) repoModel.Manufacturer {

	return repoModel.Manufacturer{
		Name:    src.Name,
		Country: src.Country,
		Website: src.Website,
	}
}

// copyMetadata создаёт глубокую копию карты metadata.
// Каждый Value внутри тоже копируется.
func copyMetadata(src map[string]repoModel.Value) map[string]repoModel.Value {
	if src == nil {
		return nil
	}
	dst := make(map[string]repoModel.Value, len(src))
	for key, val := range src {
		dst[key] = copyValue(val)
	}
	return dst
}

// copyValue копирует repoModel.Value
func copyValue(val repoModel.Value) repoModel.Value {
	var v repoModel.Value
	if val.StringValue != nil {
		s := *val.StringValue
		v.StringValue = &s
	}
	if val.Int64Value != nil {
		i := *val.Int64Value
		v.Int64Value = &i
	}
	if val.DoubleValue != nil {
		d := *val.DoubleValue
		v.DoubleValue = &d
	}
	if val.BoolValue != nil {
		b := *val.BoolValue
		v.BoolValue = &b
	}
	return v
}

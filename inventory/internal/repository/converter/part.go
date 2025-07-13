package converter

import (
	"github.com/Akbar-cmd/raketa-factory/inventory/internal/model"
	repoModel "github.com/Akbar-cmd/raketa-factory/inventory/internal/repository/model"
)

// PartToRepoModel конвертирует модель Part в репозиторную модель
func PartToRepoModel(part model.Part) repoModel.Part {
	return repoModel.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      CategoryToRepoModel(part.Category),
		Dimensions:    DimensionsToRepoModel(part.Dimensions),
		Manufacturer:  ManufacturerToRepoModel(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      MetadataToRepoModel(part.Metadata),
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}
}

// RepoModelToPart конвертирует репозиторную модель Part во внутреннюю модель
func PartToModel(repo repoModel.Part) model.Part {
	return model.Part{
		Uuid:          repo.Uuid,
		Name:          repo.Name,
		Description:   repo.Description,
		Price:         repo.Price,
		StockQuantity: repo.StockQuantity,
		Category:      CategoryToModel(repo.Category),
		Dimensions:    DimensionsToModel(repo.Dimensions),
		Manufacturer:  ManufacturerToModel(repo.Manufacturer),
		Tags:          repo.Tags,
		Metadata:      MetadataToModel(repo.Metadata),
		CreatedAt:     repo.CreatedAt,
		UpdatedAt:     repo.UpdatedAt,
	}
}

// CategoryToRepoModel конвертирует model.Category в repoModel.Category
func CategoryToRepoModel(c model.Category) repoModel.Category {
	switch c {
	case model.CategoryEngine:
		return repoModel.CategoryEngine
	case model.CategoryFuel:
		return repoModel.CategoryFuel
	case model.CategoryPorthole:
		return repoModel.CategoryPorthole
	case model.CategoryWing:
		return repoModel.CategoryWing
	default:
		return repoModel.CategoryUnknown
	}
}

// CategoryToModel конвертирует repoModel.Category в model.Category
func CategoryToModel(c repoModel.Category) model.Category {
	switch c {
	case repoModel.CategoryEngine:
		return model.CategoryEngine
	case repoModel.CategoryFuel:
		return model.CategoryFuel
	case repoModel.CategoryPorthole:
		return model.CategoryPorthole
	case repoModel.CategoryWing:
		return model.CategoryWing
	default:
		return model.CategoryUnknown
	}
}

// DimensionsToRepoModel конвертирует model.Dimensions в repoModel.Dimensions
func DimensionsToRepoModel(d model.Dimensions) repoModel.Dimensions {
	return repoModel.Dimensions{
		Length: d.Length,
		Width:  d.Width,
		Height: d.Height,
		Weight: d.Weight,
	}
}

// DimensionsToModel конвертирует repoModel.Dimensions в model.Dimensions
func DimensionsToModel(d repoModel.Dimensions) model.Dimensions {
	return model.Dimensions{
		Length: d.Length,
		Width:  d.Width,
		Height: d.Height,
		Weight: d.Weight,
	}
}

// ManufacturerToRepoModel конвертирует model.Manufacturer в repoModel.Manufacturer
func ManufacturerToRepoModel(m model.Manufacturer) repoModel.Manufacturer {
	return repoModel.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

// ManufacturerToModel конвертирует repoModel.Manufacturer в model.Manufacturer
func ManufacturerToModel(m repoModel.Manufacturer) model.Manufacturer {
	return model.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

// TagsToRepoModel поверхностно копирует срез тегов
func TagsToRepoModel(tags []string) []string {
	if tags == nil {
		return nil
	}
	out := make([]string, len(tags))
	copy(out, tags)
	return out
}

// TagsToModel поверхностно копирует срез тегов из репо
func TagsToModel(tags []string) []string {
	if tags == nil {
		return nil
	}
	out := make([]string, len(tags))
	copy(out, tags)
	return out
}

// MetadataToRepoModel копирует карту metadata и её вложенные значения
func MetadataToRepoModel(md map[string]model.Value) map[string]repoModel.Value {
	if md == nil {
		return nil
	}
	out := make(map[string]repoModel.Value, len(md))
	for k, v := range md {
		out[k] = ValueToRepoModel(v)
	}
	return out
}

// MetadataToModel копирует репозиторную карту metadata обратно в model
func MetadataToModel(md map[string]repoModel.Value) map[string]model.Value {
	if md == nil {
		return nil
	}
	out := make(map[string]model.Value, len(md))
	for k, v := range md {
		out[k] = ValueToModel(v)
	}
	return out
}

// ValueToRepoModel конвертирует model.Value в repoModel.Value
func ValueToRepoModel(v model.Value) repoModel.Value {
	var rv repoModel.Value
	if v.StringValue != nil {
		s := *v.StringValue
		rv.StringValue = &s
	}
	if v.Int64Value != nil {
		i := *v.Int64Value
		rv.Int64Value = &i
	}
	if v.DoubleValue != nil {
		f := *v.DoubleValue
		rv.DoubleValue = &f
	}
	if v.BoolValue != nil {
		b := *v.BoolValue
		rv.BoolValue = &b
	}
	return rv
}

// ValueToModel конвертирует repoModel.Value в model.Value
func ValueToModel(v repoModel.Value) model.Value {
	var mv model.Value
	if v.StringValue != nil {
		s := *v.StringValue
		mv.StringValue = &s
	}
	if v.Int64Value != nil {
		i := *v.Int64Value
		mv.Int64Value = &i
	}
	if v.DoubleValue != nil {
		f := *v.DoubleValue
		mv.DoubleValue = &f
	}
	if v.BoolValue != nil {
		b := *v.BoolValue
		mv.BoolValue = &b
	}
	return mv
}

// FilterToRepoModel конвертирует доменный PartsFilter в репозиторный
func FilterToRepoModel(f model.PartsFilter) repoModel.PartsFilter {
	return repoModel.PartsFilter{
		Uuids:                 copyStrings(f.Uuids),
		Names:                 copyStrings(f.Names),
		Categories:            CategoriesToRepoModel(f.Categories),
		ManufacturerCountries: copyStrings(f.ManufacturerCountries),
		Tags:                  copyStrings(f.Tags),
	}
}

// FilterToModel конвертирует репозиторный PartsFilter обратно в доменный.
func FilterToModel(r repoModel.PartsFilter) model.PartsFilter {
	return model.PartsFilter{
		Uuids:                 copyStrings(r.Uuids),
		Names:                 copyStrings(r.Names),
		Categories:            CategoriesToModel(r.Categories),
		ManufacturerCountries: copyStrings(r.ManufacturerCountries),
		Tags:                  copyStrings(r.Tags),
	}
}

// copyStrings возвращает новый срез строк, копируя src.
func copyStrings(src []string) []string {
	if src == nil {
		return nil
	}
	dst := make([]string, len(src))
	copy(dst, src)
	return dst
}

// CategoriesToRepoModel конвертирует []model.Categories в []repoModel.Categories.
func CategoriesToRepoModel(src []model.Category) []repoModel.Category {
	if src == nil {
		return nil
	}
	dst := make([]repoModel.Category, len(src))
	for i, c := range src {
		dst[i] = repoModel.Category(c)
	}
	return dst
}

// CategoriesToModel конвертирует []repoModel.Categories в []model.Categories.
func CategoriesToModel(src []repoModel.Category) []model.Category {
	if src == nil {
		return nil
	}
	dst := make([]model.Category, len(src))
	for i, c := range src {
		dst[i] = model.Category(c)
	}
	return dst
}

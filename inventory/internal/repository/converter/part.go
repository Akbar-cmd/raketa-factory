package converter

import (
	"github.com/Akbar-cmd/raketa-factory/inventory/internal/model"
	repoModel "github.com/Akbar-cmd/raketa-factory/inventory/internal/repository/model"
)

// PartsFilterToRepoModel конвертирует доменный PartsFilter в репозиторный
func PartsFilterToRepoModel(f model.PartsFilter) repoModel.PartsFilter {
	return repoModel.PartsFilter{
		Uuids:                 f.Uuids,
		Names:                 f.Names,
		Categories:            CategoriesToRepoModel(f.Categories),
		ManufacturerCountries: f.ManufacturerCountries,
		Tags:                  f.Tags,
	}
}

// CategoriesToRepoModel конвертирует []model.Categories в []repoModel.Categories.
func CategoriesToRepoModel(categories []model.Category) []repoModel.Category {
	if categories == nil {
		return nil
	}
	data := make([]repoModel.Category, len(categories))
	for _, c := range categories {
		data = append(data, repoModel.Category(c))
	}
	return data
}

// PartToModel конвертирует репозиторную модель Part во внутреннюю модель
func PartToModel(repo repoModel.Part) model.Part {
	return model.Part{
		Uuid:          repo.Uuid,
		Name:          repo.Name,
		Description:   repo.Description,
		Price:         repo.Price,
		StockQuantity: repo.StockQuantity,
		Category:      model.Category(repo.Category),
		Dimensions:    DimensionsToModel(repo.Dimensions),
		Manufacturer:  ManufacturerToModel(repo.Manufacturer),
		Tags:          repo.Tags,
		Metadata:      MetadataToModel(repo.Metadata),
		CreatedAt:     repo.CreatedAt,
		UpdatedAt:     repo.UpdatedAt,
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

// ManufacturerToModel конвертирует repoModel.Manufacturer в model.Manufacturer
func ManufacturerToModel(m repoModel.Manufacturer) model.Manufacturer {
	return model.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

// MetadataToModel конвертирует repoModel.Metadata в model.Metadata
func MetadataToModel(metadata repoModel.Metadata) model.Metadata {
	return model.Metadata{
		StringValue: metadata.StringValue,
		Int64Value:  metadata.Int64Value,
		DoubleValue: metadata.DoubleValue,
		BoolValue:   metadata.BoolValue,
	}
}

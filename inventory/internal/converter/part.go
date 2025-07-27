package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Akbar-cmd/raketa-factory/inventory/internal/model"
	inventoryV1 "github.com/Akbar-cmd/raketa-factory/shared/pkg/proto/inventory/v1"
)

// PartsFilterToModel funcs
func PartsFilterToModel(filter *inventoryV1.PartsFilter) model.PartsFilter {
	return model.PartsFilter{
		Uuids:                 filter.GetUuids(),
		Names:                 filter.GetNames(),
		Categories:            CategoriesToModel(filter.GetCategories()),
		ManufacturerCountries: filter.GetManufacturerCountries(),
		Tags:                  filter.GetTags(),
	}
}

func CategoriesToModel(categories []inventoryV1.Category) []model.Category {
	data := make([]model.Category, len(categories))
	for _, category := range categories {
		data = append(data, model.Category(category))
	}
	return data
}

// PartToProto funcs
func PartToProto(part model.Part) *inventoryV1.Part {
	var updatedAt *timestamppb.Timestamp
	if part.UpdatedAt != nil {
		updatedAt = timestamppb.New(*part.UpdatedAt)
	}

	return &inventoryV1.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      CategoryToProto(part.Category),
		Dimensions:    dimensionsToProto(part.Dimensions),
		Manufacturer:  manufacturerToProto(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      metadataToProto(part.Metadata),
		CreatedAt:     timestamppb.New(part.CreatedAt),
		UpdatedAt:     updatedAt,
	}
}

func CategoryToProto(category model.Category) inventoryV1.Category {
	switch category {
	case model.CategoryEngine:
		return inventoryV1.Category_ENGINE
	case model.CategoryFuel:
		return inventoryV1.Category_FUEL
	case model.CategoryPorthole:
		return inventoryV1.Category_PORTHOLE
	case model.CategoryWing:
		return inventoryV1.Category_WING
	default:
		return inventoryV1.Category_UNKNOWN
	}
}

func dimensionsToProto(d model.Dimensions) *inventoryV1.Dimensions {
	return &inventoryV1.Dimensions{
		Length: d.Length,
		Width:  d.Width,
		Height: d.Height,
		Weight: d.Weight,
	}
}

func manufacturerToProto(m model.Manufacturer) *inventoryV1.Manufacturer {
	return &inventoryV1.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

func metadataToProto(meta model.Metadata) map[string]*inventoryV1.Value {
	var val *inventoryV1.Value
	switch {
	case meta.StringValue != nil:
		val = &inventoryV1.Value{
			Kind: &inventoryV1.Value_StringValue{StringValue: *meta.StringValue},
		}
	case meta.Int64Value != nil:
		val = &inventoryV1.Value{
			Kind: &inventoryV1.Value_Int64Value{Int64Value: *meta.Int64Value},
		}
	case meta.DoubleValue != nil:
		val = &inventoryV1.Value{
			Kind: &inventoryV1.Value_DoubleValue{DoubleValue: *meta.DoubleValue},
		}
	case meta.BoolValue != nil:
		val = &inventoryV1.Value{
			Kind: &inventoryV1.Value_BoolValue{BoolValue: *meta.BoolValue},
		}
	default:
		val = &inventoryV1.Value{}
	}
	return map[string]*inventoryV1.Value{"value": val}
}

// PartsToProto func
func PartsToProto(parts []model.Part) []*inventoryV1.Part {
	data := make([]*inventoryV1.Part, len(parts))
	for i, part := range parts {
		data[i] = PartToProto(part)
	}

	return data
}

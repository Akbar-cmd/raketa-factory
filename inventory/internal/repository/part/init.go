package part

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"math"
	"time"

	repoModel "github.com/Akbar-cmd/raketa-factory/inventory/internal/repository/model"
)

// Генератор деталей

func generateParts() []repoModel.Part {
	names := []string{
		"Main Engine",
		"Reserve Engine",
		"Thruster",
		"Fuel Tank",
		"Left Wing",
		"Right Wing",
		"Window A",
		"Window B",
		"Control Module",
		"Stabilizer",
	}

	descriptions := []string{
		"Primary propulsion unit",
		"Backup propulsion unit",
		"Thruster for fine adjustments",
		"Main fuel tank",
		"Left aerodynamic wing",
		"Right aerodynamic wing",
		"Front viewing window",
		"Side viewing window",
		"Flight control module",
		"Stabilization fin",
	}

	var parts []repoModel.Part
	for i := 0; i < gofakeit.Number(1, 50); i++ {
		idx := gofakeit.Number(0, len(names)-1)
		parts = append(parts, repoModel.Part{
			Uuid:          uuid.NewString(),
			Name:          names[idx],
			Description:   descriptions[idx],
			Price:         roundTo(gofakeit.Float64Range(100, 10_000)),
			StockQuantity: int64(gofakeit.Number(1, 100)),
			Category:      repoModel.Category(gofakeit.Number(1, 4)), //nolint:gosec // safe: gofakeit.Number returns 1..4
			Dimensions:    generateDimensions(),
			Manufacturer:  generateManufacturer(),
			Tags:          generateTags(),
			Metadata:      generateMetadata(),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		})
	}

	return parts
}

func generateDimensions() repoModel.Dimensions {
	return repoModel.Dimensions{
		Length: roundTo(gofakeit.Float64Range(1, 1000)),
		Width:  roundTo(gofakeit.Float64Range(1, 1000)),
		Height: roundTo(gofakeit.Float64Range(1, 1000)),
		Weight: roundTo(gofakeit.Float64Range(1, 1000)),
	}
}

func generateManufacturer() repoModel.Manufacturer {
	return repoModel.Manufacturer{
		Name:    gofakeit.Name(),
		Country: gofakeit.Country(),
		Website: gofakeit.URL(),
	}
}

func generateTags() []string {
	var tags []string
	for i := 0; i < gofakeit.Number(1, 10); i++ {
		tags = append(tags, gofakeit.EmojiTag())
	}

	return tags
}

func generateMetadata() map[string]repoModel.Value {
	metadata := make(map[string]repoModel.Value)

	for i := 0; i < gofakeit.Number(1, 10); i++ {
		metadata[gofakeit.Word()] = generateMetadataValue()
	}

	return metadata
}

func generateMetadataValue() repoModel.Value {
	switch gofakeit.Number(0, 3) {
	case 0:
		s := gofakeit.Word()
		return repoModel.Value{StringValue: &s}

	case 1:
		i := int64(gofakeit.Number(1, 100))
		return repoModel.Value{Int64Value: &i}

	case 2:
		f := roundTo(gofakeit.Float64Range(1, 100))
		return repoModel.Value{DoubleValue: &f}

	case 3:
		b := gofakeit.Bool()
		return repoModel.Value{BoolValue: &b}

	default:
		return repoModel.Value{}
	}
}

func roundTo(x float64) float64 {
	return math.Round(x*100) / 100
}

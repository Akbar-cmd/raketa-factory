package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	inventoryV1 "github.com/Akbar-cmd/raketa-factory/shared/pkg/proto/inventory/v1"
)

const (
	grpcPort = 50051
	httpPort = 8081
)

// InventoryServer реализует gRPC сервис для работы с наблюдениями НЛО
type InventoryService struct {
	inventoryV1.UnimplementedInventoryServiceServer

	mu    sync.RWMutex
	parts map[string]*inventoryV1.Part
}

func NewInventoryService() *InventoryService {
	svc := &InventoryService{
		parts: make(map[string]*inventoryV1.Part),
	}
	svc.initParts()
	return svc
}

// Генератор деталей
func (s *InventoryService) initParts() {
	parts := generateParts()

	for _, part := range parts {
		s.parts[part.Uuid] = part
	}

	log.Printf("✅ Инициализировано %d запчастей в inventory", len(parts))
}

func generateParts() []*inventoryV1.Part {
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

	var parts []*inventoryV1.Part
	for i := 0; i < gofakeit.Number(1, 50); i++ {
		idx := gofakeit.Number(0, len(names)-1)
		parts = append(parts, &inventoryV1.Part{
			Uuid:          uuid.NewString(),
			Name:          names[idx],
			Description:   descriptions[idx],
			Price:         roundTo(gofakeit.Float64Range(100, 10_000)),
			StockQuantity: int64(gofakeit.Number(1, 100)),
			Category:      inventoryV1.Category(gofakeit.Number(1, 4)), //nolint:gosec // safe: gofakeit.Number returns 1..4
			Dimensions:    generateDimensions(),
			Manufacturer:  generateManufacturer(),
			Tags:          generateTags(),
			Metadata:      generateMetadata(),
			CreatedAt:     timestamppb.Now(),
		})
	}

	return parts
}

func generateDimensions() *inventoryV1.Dimensions {
	return &inventoryV1.Dimensions{
		Length: roundTo(gofakeit.Float64Range(1, 1000)),
		Width:  roundTo(gofakeit.Float64Range(1, 1000)),
		Height: roundTo(gofakeit.Float64Range(1, 1000)),
		Weight: roundTo(gofakeit.Float64Range(1, 1000)),
	}
}

func generateManufacturer() *inventoryV1.Manufacturer {
	return &inventoryV1.Manufacturer{
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

func generateMetadata() map[string]*inventoryV1.Value {
	metadata := make(map[string]*inventoryV1.Value)

	for i := 0; i < gofakeit.Number(1, 10); i++ {
		metadata[gofakeit.Word()] = generateMetadataValue()
	}

	return metadata
}

func generateMetadataValue() *inventoryV1.Value {
	switch gofakeit.Number(0, 3) {
	case 0:
		return &inventoryV1.Value{
			Kind: &inventoryV1.Value_StringValue{
				StringValue: gofakeit.Word(),
			},
		}

	case 1:
		return &inventoryV1.Value{
			Kind: &inventoryV1.Value_Int64Value{
				Int64Value: int64(gofakeit.Number(1, 100)),
			},
		}

	case 2:
		return &inventoryV1.Value{
			Kind: &inventoryV1.Value_DoubleValue{
				DoubleValue: roundTo(gofakeit.Float64Range(1, 100)),
			},
		}

	case 3:
		return &inventoryV1.Value{
			Kind: &inventoryV1.Value_BoolValue{
				BoolValue: gofakeit.Bool(),
			},
		}

	default:
		return nil
	}
}

func roundTo(x float64) float64 {
	return math.Round(x*100) / 100
}

// matchesFilter фильтрует детали
func matchesFilter(part *inventoryV1.Part, filter *inventoryV1.PartsFilter) bool {
	if filter == nil {
		return true
	}
	return matchUUID(part, filter) &&
		matchName(part, filter) &&
		matchCategory(part, filter) &&
		matchManufacturerCountry(part, filter) &&
		matchTags(part, filter)
}

// Фильтрация по UUID
func matchUUID(part *inventoryV1.Part, filter *inventoryV1.PartsFilter) bool {
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
func matchName(part *inventoryV1.Part, filter *inventoryV1.PartsFilter) bool {
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
func matchCategory(part *inventoryV1.Part, filter *inventoryV1.PartsFilter) bool {
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
func matchManufacturerCountry(part *inventoryV1.Part, filter *inventoryV1.PartsFilter) bool {
	if len(filter.ManufacturerCountries) == 0 {
		return true
	}
	country := part.Manufacturer.GetCountry()
	for _, c := range filter.ManufacturerCountries {
		if strings.EqualFold(country, c) {
			return true
		}
	}
	return false
}

// Фильтрация по тегам
func matchTags(part *inventoryV1.Part, filter *inventoryV1.PartsFilter) bool {
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
func copyDimensions(src *inventoryV1.Dimensions) *inventoryV1.Dimensions {
	if src == nil {
		return nil
	}

	return &inventoryV1.Dimensions{
		Length: src.Length,
		Width:  src.Width,
		Height: src.Height,
		Weight: src.Weight,
	}
}

// copyManufacturer создает копию объекта Manufacturer
func copyManufacturer(src *inventoryV1.Manufacturer) *inventoryV1.Manufacturer {
	if src == nil {
		return nil
	}

	return &inventoryV1.Manufacturer{
		Name:    src.Name,
		Country: src.Country,
		Website: src.Website,
	}
}

// copyMetadata создаёт глубокую копию карты metadata.
// Каждый Value внутри тоже копируется.
func copyMetadata(src map[string]*inventoryV1.Value) map[string]*inventoryV1.Value {
	if src == nil {
		return nil
	}
	dst := make(map[string]*inventoryV1.Value, len(src))
	for key, val := range src {
		if val == nil {
			dst[key] = nil
			continue
		}
		// Копируем по Kind
		switch kind := val.Kind.(type) {
		case *inventoryV1.Value_StringValue:
			dst[key] = &inventoryV1.Value{
				Kind: &inventoryV1.Value_StringValue{
					StringValue: kind.StringValue,
				},
			}
		case *inventoryV1.Value_Int64Value:
			dst[key] = &inventoryV1.Value{
				Kind: &inventoryV1.Value_Int64Value{
					Int64Value: kind.Int64Value,
				},
			}
		case *inventoryV1.Value_DoubleValue:
			dst[key] = &inventoryV1.Value{
				Kind: &inventoryV1.Value_DoubleValue{
					DoubleValue: kind.DoubleValue,
				},
			}
		case *inventoryV1.Value_BoolValue:
			dst[key] = &inventoryV1.Value{
				Kind: &inventoryV1.Value_BoolValue{
					BoolValue: kind.BoolValue,
				},
			}
		default:
			// на случай расширения oneof
			dst[key] = nil
		}
	}
	return dst
}

// GetPart возвращает информацию по детали по ее UUID
func (s *InventoryService) GetPart(_ context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	part, ok := s.parts[req.GetUuid()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "parts with UUID %s not found", req.GetUuid())
	}

	return &inventoryV1.GetPartResponse{
		Part: part,
	}, nil
}

func (s *InventoryService) ListParts(_ context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var parts []*inventoryV1.Part
	filter := req.GetFilter()

	for _, part := range s.parts {
		if matchesFilter(part, filter) {
			parts = append(parts, &inventoryV1.Part{
				Uuid:          part.Uuid,
				Name:          part.Name,
				Description:   part.Description,
				Price:         part.Price,
				StockQuantity: part.StockQuantity,
				Category:      part.Category,
				Dimensions:    copyDimensions(part.Dimensions),
				Manufacturer:  copyManufacturer(part.Manufacturer),
				Tags:          append([]string{}, part.Tags...),
				Metadata:      copyMetadata(part.Metadata),
				CreatedAt:     part.CreatedAt,
			})
		}
	}

	return &inventoryV1.ListPartsResponse{
		Parts: parts,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	// Создаем gRPC сервер
	s := grpc.NewServer()

	// Регистрируем наш сервис
	service := NewInventoryService()

	inventoryV1.RegisterInventoryServiceServer(s, service)

	// включаем рефлексию для отладки
	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}

package transactiondatarepository

import (
	"context"

	"github.com/google/uuid"
	"github.com/thdoikn/sihp-be/internal/entity"
	entitybase "github.com/thdoikn/sihp-be/internal/entity/base"
)

type TransactionDataRepository interface {
	CreatePengumpulanData(ctx context.Context, input *entity.PengumpulanData) (*entity.PengumpulanData, error)
	GetPengumpulanDataByID(ctx context.Context, id uuid.UUID) (*entity.PengumpulanData, error)
	GetPengumpulanDataByFilter(ctx context.Context, filter *entity.PengumpulanDataFilter) ([]entity.PengumpulanData, entitybase.BasePaginationResult, error)
	UpdatePengumpulanData(ctx context.Context, id uuid.UUID, updateMap map[string]any) (*entity.PengumpulanData, error)
	DeletePengumpulanData(ctx context.Context, id uuid.UUID) error
	FinalizePengumpulanData(ctx context.Context, id uuid.UUID) (int64, error)

	CreateHargaRutin(ctx context.Context, input *entity.HargaRutin) (*entity.HargaRutin, error)
	GetHargaRutinByID(ctx context.Context, id uuid.UUID) (*entity.HargaRutin, error)
	GetHargaRutinByFilter(ctx context.Context, filter *entity.HargaRutinFilter) ([]entity.HargaRutin, entitybase.BasePaginationResult, error)
	UpdateHargaRutin(ctx context.Context, id uuid.UUID, updateMap map[string]any) (*entity.HargaRutin, error)
	DeleteHargaRutin(ctx context.Context, id uuid.UUID) error

	GetHargaPelaporanByID(ctx context.Context, id uuid.UUID) (*entity.HargaPelaporanEnriched, error)
	GetHargaPelaporanByFilter(ctx context.Context, filter *entity.HargaPelaporanFilter) ([]entity.HargaPelaporanEnriched, entitybase.BasePaginationResult, error)
}

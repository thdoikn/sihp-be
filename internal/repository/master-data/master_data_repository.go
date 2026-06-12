package masterdatarepository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/thdoikn/sihp-be/internal/entity"
	entitybase "github.com/thdoikn/sihp-be/internal/entity/base"
)

type MasterDataRepository interface {
	GetOverview(ctx context.Context) (pasarCount, tempatUsahaCount, komoditasCount int64, err error)
	CreatePasar(ctx context.Context, input *entity.Pasar) (*entity.Pasar, error)
	GetPasarByID(ctx context.Context, id uuid.UUID) (*entity.Pasar, error)
	GetPasarByFilter(ctx context.Context, filter *entity.PasarFilter) ([]entity.Pasar, entitybase.BasePaginationResult, error)
	UpdatePasar(ctx context.Context, id uuid.UUID, updateMap map[string]any) (*entity.Pasar, error)

	CreateKomoditas(ctx context.Context, input *entity.Komoditas) (*entity.Komoditas, error)
	GetKomoditasByID(ctx context.Context, id uuid.UUID) (*entity.Komoditas, error)
	GetKomoditasByFilter(ctx context.Context, filter *entity.KomoditasFilter) ([]entity.Komoditas, entitybase.BasePaginationResult, error)
	UpdateKomoditas(ctx context.Context, id uuid.UUID, updateMap map[string]any) (*entity.Komoditas, error)

	CreateTempatUsaha(ctx context.Context, input *entity.TempatUsaha) (*entity.TempatUsaha, error)
	GetTempatUsahaByID(ctx context.Context, id uuid.UUID) (*entity.TempatUsaha, error)
	GetTempatUsahaByFilter(ctx context.Context, filter *entity.TempatUsahaFilter) ([]entity.TempatUsaha, entitybase.BasePaginationResult, error)
	UpdateTempatUsaha(ctx context.Context, id uuid.UUID, updateMap map[string]any) (*entity.TempatUsaha, error)

	CreateKomoditasDijual(ctx context.Context, input *entity.KomoditasDijual) (*entity.KomoditasDijual, error)
	GetKomoditasDijualByID(ctx context.Context, id uuid.UUID) (*entity.KomoditasDijual, error)
	GetKomoditasDijualByFilter(ctx context.Context, filter *entity.KomoditasDijualFilter) ([]entity.KomoditasDijual, entitybase.BasePaginationResult, error)
	UpdateKomoditasDijual(ctx context.Context, id uuid.UUID, updateMap map[string]any) (*entity.KomoditasDijual, error)

	GetPublicKomoditasList(ctx context.Context, filter *entity.PublicKomoditasListFilter) ([]entity.PublicKomoditasListItem, entitybase.BasePaginationResult, error)
	GetPublicKomoditasStats(ctx context.Context, komoditasID uuid.UUID, days int, idPasar *uuid.UUID) (*time.Time, *float64, *float64, *float64, *float64, error)
	GetPublicKomoditasTrend(ctx context.Context, komoditasID uuid.UUID, days int, idPasar *uuid.UUID) ([]map[string]any, error)
	GetPublicPasarList(ctx context.Context) ([]entity.PublicPasarListItem, error)
	GetPublicPasarDetail(ctx context.Context, id uuid.UUID, filter *entity.TempatUsahaFilter) (*entity.Pasar, []entity.TempatUsaha, entitybase.BasePaginationResult, error)
	GetPublicTempatUsahaList(ctx context.Context, filter *entity.PublicTempatUsahaListFilter) ([]entity.PublicTempatUsahaListItem, entitybase.BasePaginationResult, error)
	GetPublicTempatUsahaDetail(ctx context.Context, id uuid.UUID, filter entity.KomoditasFilter) (*entity.TempatUsaha, *entity.Pasar, []map[string]any, entitybase.BasePaginationResult, error)
}

package sihpserializer

import (
	"github.com/thdoikn/sihp-be/internal/entity"
	entitybase "github.com/thdoikn/sihp-be/internal/entity/base"
	"github.com/thdoikn/sihp-be/pkg/dto"
	dtobase "github.com/thdoikn/sihp-be/pkg/dto/base"
)

type SIHPSerializer interface {
	ToPasar(entity.Pasar) dto.ResPasar
	ToKomoditas(entity.Komoditas) dto.ResKomoditas
	ToTempatUsaha(entity.TempatUsaha) dto.ResTempatUsaha
	ToKomoditasDijual(entity.KomoditasDijual) dto.ResKomoditasDijual
	ToPengumpulanData(entity.PengumpulanData) dto.ResPengumpulanData
	ToHargaRutin(entity.HargaRutin) dto.ResHargaRutin
	ToHargaPelaporan(entity.HargaPelaporanEnriched) dto.ResHargaPelaporan
	ToPage(entitybase.BasePaginationResult) dtobase.BasePagination
}

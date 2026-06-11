package entity

import (
	"github.com/google/uuid"
	entitybase "github.com/thdoikn/sihp-be/internal/entity/base"
)

type HargaRutin struct {
	entitybase.Base
	IDPengumpulanData uuid.UUID `gorm:"column:id_pengumpulan_data;type:uuid;not null"`
	IDTempatUsaha     uuid.UUID `gorm:"column:id_tempat_usaha;type:uuid;not null"`
	IDKomoditas       uuid.UUID `gorm:"column:id_komoditas;type:uuid;not null"`
	KelasKomoditas    string    `gorm:"column:kelas_komoditas;type:text;not null"`
	Harga             int64     `gorm:"column:harga;type:bigint;not null"`
}

func (h *HargaRutin) TableName() string { return "sihp.harga_rutin" }
func (h *HargaRutin) OrderMap() map[string]bool {
	out := entitybase.GenerateBaseOrderMap()
	out["kelas_komoditas"] = true
	out["harga"] = true
	return out
}

type HargaRutinFilter struct {
	IDPengumpulanData *uuid.UUID
	IDKomoditas       *uuid.UUID
	IDTempatUsaha     *uuid.UUID
	PaginationFilter  entitybase.BasePaginationFilter
}

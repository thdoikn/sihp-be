package entity

import (
	"time"

	"github.com/google/uuid"
	entitybase "github.com/thdoikn/sihp-be/internal/entity/base"
	"github.com/thdoikn/sihp-be/pkg/constant"
)

type HargaPelaporan struct {
	entitybase.Base
	IDPengumpulanData uuid.UUID `gorm:"column:id_pengumpulan_data;type:uuid;not null"`
	IDKomoditas       uuid.UUID `gorm:"column:id_komoditas;type:uuid;not null"`
	Tanggal           time.Time `gorm:"column:tanggal;type:date;not null"`
	Harga             int64     `gorm:"column:harga;type:bigint;not null"`
}

func (h *HargaPelaporan) TableName() string { return "sihp.harga_pelaporan" }
func (h *HargaPelaporan) OrderMap() map[string]bool {
	out := entitybase.GenerateBaseOrderMap()
	out["tanggal"] = true
	out["harga"] = true
	return out
}

type HargaPelaporanFilter struct {
	IDPasar          *uuid.UUID
	IDKomoditas      *uuid.UUID
	From             *time.Time
	To               *time.Time
	Status           *constant.PengumpulanDataStatus
	PaginationFilter entitybase.BasePaginationFilter
}

type HargaPelaporanEnriched struct {
	HargaPelaporan
	IDPasar       uuid.UUID
	HargaBesar    *int64
	HargaMenengah *int64
	HargaKecil    *int64
}

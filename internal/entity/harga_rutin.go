package entity

import (
	"time"

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
	HargaInput        *int64    `gorm:"column:harga_input;type:bigint"`
	JumlahInput       *float64  `gorm:"column:jumlah_input;type:decimal(10,2)"`
	SatuanInput       *string   `gorm:"column:satuan_input;type:varchar(100)"`
	Status            int16     `gorm:"column:status;type:smallint;not null;default:0"`
	NamaEnumerator    *string   `gorm:"column:nama_enumerator;type:varchar(255)"`
	// Denormalized fields from JOIN (not stored in DB)
	Tanggal time.Time `gorm:"column:tanggal;type:date" json:"tanggal"`
	IDPasar uuid.UUID `gorm:"column:id_pasar;type:uuid" json:"id_pasar"`
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

package entity

import (
	"github.com/google/uuid"
	entitybase "github.com/thdoikn/sihp-be/internal/entity/base"
	"github.com/thdoikn/sihp-be/pkg/constant"
)

type Pasar struct {
	entitybase.Base
	Nama   string                        `gorm:"column:nama;type:varchar(255);not null"`
	Alamat *string                       `gorm:"column:alamat;type:text"`
	Status constant.ActiveInactiveStatus `gorm:"column:status;type:smallint;not null;default:1"`

	// Relation
	TempatUsaha []TempatUsaha `gorm:"foreignKey:IDPasar;references:ID"`
}

func (p *Pasar) TableName() string {
	return "sihp.pasar"
}

func (p *Pasar) OrderMap() map[string]bool {
	out := entitybase.GenerateBaseOrderMap()
	out["nama"] = true
	out["status"] = true
	return out
}

type Komoditas struct {
	entitybase.Base
	Nama      string  `gorm:"column:nama;type:varchar(255);not null"`
	Satuan    *string `gorm:"column:satuan;type:varchar(100)"`
	GambarURL *string `gorm:"column:gambar_url;type:text"`
}

func (k *Komoditas) TableName() string {
	return "sihp.komoditas"
}

func (k *Komoditas) OrderMap() map[string]bool {
	out := entitybase.GenerateBaseOrderMap()
	out["nama"] = true
	return out
}

type TempatUsaha struct {
	entitybase.Base
	IDPasar uuid.UUID                     `gorm:"column:id_pasar;type:uuid;not null"`
	Nama    string                        `gorm:"column:nama;type:varchar(255);not null"`
	Pemilik *string                       `gorm:"column:pemilik;type:varchar(255)"`
	Status  constant.ActiveInactiveStatus `gorm:"column:status;type:smallint;not null;default:1"`

	// Relation
	KomoditasDijual []KomoditasDijual `gorm:"foreignKey:IDTempatUsaha;references:ID"`
}

func (t *TempatUsaha) TableName() string {
	return "sihp.tempat_usaha"
}

func (t *TempatUsaha) OrderMap() map[string]bool {
	out := entitybase.GenerateBaseOrderMap()
	out["nama"] = true
	out["status"] = true
	return out
}

type KomoditasDijual struct {
	entitybase.Base
	IDTempatUsaha            uuid.UUID                     `gorm:"column:id_tempat_usaha;type:uuid;not null"`
	IDKomoditas              uuid.UUID                     `gorm:"column:id_komoditas;type:uuid;not null"`
	HargaNormal              float64                       `gorm:"column:harga_normal;type:double precision;not null;default:0"`
	HargaMahal               float64                       `gorm:"column:harga_mahal;type:double precision;not null;default:0"`
	HargaAvg                 *float64                      `gorm:"column:harga_avg;type:double precision;->"`
	SatuanStok               string                        `gorm:"column:satuan_stok;type:sihp.satuan_stok_enum;not null;default:'kg'"`
	NilaiStok                float64                       `gorm:"column:nilai_stok;type:numeric(18,4);not null;default:0"`
	SatuanPeriode            string                        `gorm:"column:satuan_periode;type:sihp.satuan_periode_enum;not null;default:'minggu'"`
	NilaiPeriode             int                           `gorm:"column:nilai_periode;type:integer;not null;default:1"`
	LokasiSupplier           string                        `gorm:"column:lokasi_supplier;type:varchar(255);not null;default:''"`
	PolaDistribusi           *string                       `gorm:"column:pola_distribusi;type:varchar(50)"`
	StandardizedStockPeriode float64                       `gorm:"column:standardized_stock_periode;type:numeric(18,4);not null;default:0"`
	KelasKomoditas           *string                       `gorm:"column:kelas_komoditas;type:sihp.kelas_komoditas_enum"`
	Status                   constant.ActiveInactiveStatus `gorm:"column:status;type:smallint;not null;default:1"`

	// Relation
	Komoditas   Komoditas   `gorm:"foreignKey:IDKomoditas;references:ID"`
	TempatUsaha TempatUsaha `gorm:"foreignKey:IDTempatUsaha;references:ID"`
}

func (k *KomoditasDijual) TableName() string {
	return "sihp.komoditas_dijual"
}

func (k *KomoditasDijual) OrderMap() map[string]bool {
	out := entitybase.GenerateBaseOrderMap()
	out["status"] = true
	return out
}

type PasarFilter struct {
	Name             *string
	Status           *constant.ActiveInactiveStatus
	PaginationFilter entitybase.BasePaginationFilter
}

type KomoditasFilter struct {
	Name             *string
	IDTempatUsaha    *uuid.UUID
	IDPasar          *uuid.UUID
	PaginationFilter entitybase.BasePaginationFilter
}

type TempatUsahaFilter struct {
	Name             *string
	IDPasar          *uuid.UUID
	Status           *constant.ActiveInactiveStatus
	PaginationFilter entitybase.BasePaginationFilter
}

type KomoditasDijualFilter struct {
	IDTempatUsaha    *uuid.UUID
	IDKomoditas      *uuid.UUID
	Status           *constant.ActiveInactiveStatus
	PaginationFilter entitybase.BasePaginationFilter
}

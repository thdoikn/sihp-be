package dto

import (
	"time"

	"github.com/google/uuid"
	dtobase "github.com/thdoikn/sihp-be/pkg/dto/base"
)

type ResAuthLogin struct {
	User  ResAdmin     `json:"user"`
	Token ResAuthToken `json:"token"`
}

type ResAuthToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type ReqAuthLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type ResAdmin struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Name  string    `json:"name"`
}

type ResAuthLoginEnvelope struct {
	dtobase.BaseRes
	Data *ResAuthLogin `json:"data,omitempty"`
}

type ReqGetPasar struct {
	Name   *string `query:"name"`
	Status *int16  `query:"status"`
	dtobase.BaseReqQueryPagination
}

type ReqCreatePasar struct {
	Nama   string  `json:"nama" validate:"required"`
	Alamat *string `json:"alamat"`
}

type ReqUpdatePasar struct {
	Nama   *string `json:"nama"`
	Alamat *string `json:"alamat"`
	Status *int16  `json:"status"`
}

type ResPasar struct {
	ID     uuid.UUID `json:"id"`
	Nama   string    `json:"nama"`
	Alamat *string   `json:"alamat"`
	Status int16     `json:"status"`
}

type ResPasarSingle struct {
	dtobase.BaseRes
	Data *ResPasar `json:"data,omitempty"`
}

type ResPasarList struct {
	dtobase.BaseResPagination
	Data []ResPasar `json:"data"`
}

type ReqGetKomoditas struct {
	Name          *string    `query:"name"`
	IDTempatUsaha *uuid.UUID `query:"id_tempat_usaha"`
	IDPasar       *uuid.UUID `query:"id_pasar"`
	dtobase.BaseReqQueryPagination
}

type ReqCreateKomoditas struct {
	Nama   string  `json:"nama" validate:"required"`
	Satuan *string `json:"satuan"`
}

type ReqUpdateKomoditas struct {
	Nama   *string `json:"nama"`
	Satuan *string `json:"satuan"`
}

type ResKomoditas struct {
	ID        uuid.UUID `json:"id"`
	Nama      string    `json:"nama"`
	Satuan    *string   `json:"satuan"`
	GambarURL *string   `json:"gambar_url"`
}

type ResKomoditasSingle struct {
	dtobase.BaseRes
	Data *ResKomoditas `json:"data,omitempty"`
}

type ResKomoditasList struct {
	dtobase.BaseResPagination
	Data []ResKomoditas `json:"data"`
}

type ReqGetTempatUsaha struct {
	Name    *string    `query:"name"`
	IDPasar *uuid.UUID `query:"id_pasar"`
	Status  *int16     `query:"status"`
	dtobase.BaseReqQueryPagination
}

type ReqCreateTempatUsaha struct {
	IDPasar uuid.UUID `json:"id_pasar" validate:"required"`
	Nama    string    `json:"nama" validate:"required"`
	Pemilik *string   `json:"pemilik"`
}

type ReqUpdateTempatUsaha struct {
	IDPasar *uuid.UUID `json:"id_pasar"`
	Nama    *string    `json:"nama"`
	Pemilik *string    `json:"pemilik"`
	Status  *int16     `json:"status"`
}

type ResTempatUsaha struct {
	ID      uuid.UUID `json:"id"`
	IDPasar uuid.UUID `json:"id_pasar"`
	Nama    string    `json:"nama"`
	Pemilik *string   `json:"pemilik"`
	Status  int16     `json:"status"`
}

type ResTempatUsahaSingle struct {
	dtobase.BaseRes
	Data *ResTempatUsaha `json:"data,omitempty"`
}

type ResTempatUsahaList struct {
	dtobase.BaseResPagination
	Data []ResTempatUsaha `json:"data"`
}

type ReqGetKomoditasDijual struct {
	IDTempatUsaha *uuid.UUID `query:"id_tempat_usaha"`
	IDKomoditas   *uuid.UUID `query:"id_komoditas"`
	Status        *int16     `query:"status"`
	dtobase.BaseReqQueryPagination
}

type ReqCreateKomoditasDijual struct {
	IDTempatUsaha uuid.UUID `json:"id_tempat_usaha" validate:"required"`
	IDKomoditas   uuid.UUID `json:"id_komoditas" validate:"required"`
}

type ReqUpdateKomoditasDijual struct {
	Status *int16 `json:"status"`
}

type ResKomoditasDijual struct {
	ID            uuid.UUID `json:"id"`
	IDTempatUsaha uuid.UUID `json:"id_tempat_usaha"`
	IDKomoditas   uuid.UUID `json:"id_komoditas"`
	Status        int16     `json:"status"`
}

type ResKomoditasDijualSingle struct {
	dtobase.BaseRes
	Data *ResKomoditasDijual `json:"data,omitempty"`
}

type ResKomoditasDijualList struct {
	dtobase.BaseResPagination
	Data []ResKomoditasDijual `json:"data"`
}

type ReqGetPengumpulanData struct {
	IDPasar *uuid.UUID `query:"id_pasar"`
	Status  *int16     `query:"status"`
	From    *time.Time `query:"from"`
	To      *time.Time `query:"to"`
	dtobase.BaseReqQueryPagination
}

type ReqCreatePengumpulanData struct {
	IDPasar uuid.UUID `json:"id_pasar" validate:"required"`
	Tanggal time.Time `json:"tanggal" validate:"required"`
	Catatan *string   `json:"catatan"`
}

type ReqUpdatePengumpulanData struct {
	Tanggal *time.Time `json:"tanggal"`
	Catatan *string    `json:"catatan"`
}

type ResPengumpulanData struct {
	ID      uuid.UUID `json:"id"`
	IDPasar uuid.UUID `json:"id_pasar"`
	Tanggal time.Time `json:"tanggal"`
	Status  int16     `json:"status"`
	Catatan *string   `json:"catatan"`
}

type ResPengumpulanDataSingle struct {
	dtobase.BaseRes
	Data *ResPengumpulanData `json:"data,omitempty"`
}

type ResPengumpulanDataList struct {
	dtobase.BaseResPagination
	Data []ResPengumpulanData `json:"data"`
}

type ResFinalizePengumpulanData struct {
	FinalizedKomoditasCount int64 `json:"finalized_komoditas_count"`
}

type ResFinalizePengumpulanDataEnvelope struct {
	dtobase.BaseRes
	Data *ResFinalizePengumpulanData `json:"data,omitempty"`
}

type ReqGetHargaRutin struct {
	IDPengumpulanData *uuid.UUID `query:"id_pengumpulan_data"`
	IDKomoditas       *uuid.UUID `query:"id_komoditas"`
	IDTempatUsaha     *uuid.UUID `query:"id_tempat_usaha"`
	dtobase.BaseReqQueryPagination
}

type ReqCreateHargaRutin struct {
	IDPengumpulanData uuid.UUID `json:"id_pengumpulan_data" validate:"required"`
	IDTempatUsaha     uuid.UUID `json:"id_tempat_usaha" validate:"required"`
	IDKomoditas       uuid.UUID `json:"id_komoditas" validate:"required"`
	KelasKomoditas    string    `json:"kelas_komoditas" validate:"required"`
	Harga             int64     `json:"harga" validate:"required,min=1"`
}

type ReqUpdateHargaRutin struct {
	IDTempatUsaha  *uuid.UUID `json:"id_tempat_usaha"`
	KelasKomoditas *string    `json:"kelas_komoditas"`
	Harga          *int64     `json:"harga"`
}

type ResHargaRutin struct {
	ID                uuid.UUID `json:"id"`
	IDPengumpulanData uuid.UUID `json:"id_pengumpulan_data"`
	IDTempatUsaha     uuid.UUID `json:"id_tempat_usaha"`
	IDKomoditas       uuid.UUID `json:"id_komoditas"`
	KelasKomoditas    string    `json:"kelas_komoditas"`
	Harga             int64     `json:"harga"`
}

type ResHargaRutinSingle struct {
	dtobase.BaseRes
	Data *ResHargaRutin `json:"data,omitempty"`
}

type ResHargaRutinList struct {
	dtobase.BaseResPagination
	Data []ResHargaRutin `json:"data"`
}

type ReqGetHargaPelaporan struct {
	IDPasar     *uuid.UUID `query:"id_pasar"`
	IDKomoditas *uuid.UUID `query:"id_komoditas"`
	From        *time.Time `query:"from"`
	To          *time.Time `query:"to"`
	Status      *int16     `query:"status"`
	dtobase.BaseReqQueryPagination
}

type ResHargaPelaporan struct {
	ID                uuid.UUID `json:"id"`
	IDPengumpulanData uuid.UUID `json:"id_pengumpulan_data"`
	IDKomoditas       uuid.UUID `json:"id_komoditas"`
	Tanggal           time.Time `json:"tanggal"`
	Harga             int64     `json:"harga"`
}

type ResHargaPelaporanSingle struct {
	dtobase.BaseRes
	Data *ResHargaPelaporan `json:"data,omitempty"`
}

type ResHargaPelaporanList struct {
	dtobase.BaseResPagination
	Data []ResHargaPelaporan `json:"data"`
}

type ReqPublicKomoditasDetail struct {
	Days    *int       `query:"days"`
	IDPasar *uuid.UUID `query:"id_pasar"`
}

type ResPublicOverview struct {
	PasarActiveCount       int64 `json:"pasar_active_count"`
	TempatUsahaActiveCount int64 `json:"tempat_usaha_active_count"`
	KomoditasCount         int64 `json:"komoditas_count"`
}

type ResPublicOverviewEnvelope struct {
	dtobase.BaseRes
	Data *ResPublicOverview `json:"data,omitempty"`
}

type ResPublicHargaStat struct {
	Tanggal       *time.Time `json:"tanggal,omitempty"`
	HargaRataRata *float64   `json:"harga_rata_rata,omitempty"`
}

type ResPublicKomoditasDetail struct {
	Komoditas ResKomoditas       `json:"komoditas"`
	Latest    ResPublicHargaStat `json:"latest"`
	AvgND     *float64           `json:"avg_nd"`
	MinND     *float64           `json:"min_nd"`
	MaxND     *float64           `json:"max_nd"`
	Days      int                `json:"days"`
}

type ResPublicKomoditasDetailEnvelope struct {
	dtobase.BaseRes
	Data *ResPublicKomoditasDetail `json:"data,omitempty"`
}

type ResPublicTrendPoint struct {
	Tanggal       time.Time `json:"tanggal"`
	HargaRataRata float64   `json:"harga_rata_rata"`
}

type ResPublicTrendEnvelope struct {
	dtobase.BaseRes
	Data []ResPublicTrendPoint `json:"data"`
}

type ResPublicPasarDetail struct {
	Pasar       ResPasar               `json:"pasar"`
	TempatUsaha []ResTempatUsaha       `json:"tempat_usaha"`
	Page        dtobase.BasePagination `json:"page"`
}

type ResPublicPasarDetailEnvelope struct {
	dtobase.BaseRes
	Data *ResPublicPasarDetail `json:"data,omitempty"`
}

type ResPublicTempatUsahaKomoditas struct {
	ResKomoditas
	Latest ResPublicHargaStat `json:"latest"`
}

type ResPublicTempatUsahaDetail struct {
	TempatUsaha ResTempatUsaha                  `json:"tempat_usaha"`
	Pasar       ResPasar                        `json:"pasar"`
	Komoditas   []ResPublicTempatUsahaKomoditas `json:"komoditas"`
	Page        dtobase.BasePagination          `json:"page"`
}

type ResPublicTempatUsahaDetailEnvelope struct {
	dtobase.BaseRes
	Data *ResPublicTempatUsahaDetail `json:"data,omitempty"`
}

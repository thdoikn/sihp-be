package sihpserializerimplementation

import (
	"github.com/thdoikn/sihp-be/internal/entity"
	entitybase "github.com/thdoikn/sihp-be/internal/entity/base"
	sihpserializer "github.com/thdoikn/sihp-be/internal/serializer/sihp"
	"github.com/thdoikn/sihp-be/pkg/dto"
	dtobase "github.com/thdoikn/sihp-be/pkg/dto/base"
)

type serializer struct{}

func NewSIHPSerializer() sihpserializer.SIHPSerializer { return &serializer{} }

func (s *serializer) ToPasar(in entity.Pasar) dto.ResPasar {
	return dto.ResPasar{ID: in.ID, Nama: in.Nama, Alamat: in.Alamat, Status: int16(in.Status)}
}
func (s *serializer) ToKomoditas(in entity.Komoditas) dto.ResKomoditas {
	return dto.ResKomoditas{ID: in.ID, Nama: in.Nama, Satuan: in.Satuan, GambarURL: in.GambarURL}
}
func (s *serializer) ToTempatUsaha(in entity.TempatUsaha) dto.ResTempatUsaha {
	return dto.ResTempatUsaha{ID: in.ID, IDPasar: in.IDPasar, Nama: in.Nama, Pemilik: in.Pemilik, Status: int16(in.Status)}
}
func (s *serializer) ToKomoditasDijual(in entity.KomoditasDijual) dto.ResKomoditasDijual {
	return dto.ResKomoditasDijual{
		ID:                       in.ID,
		IDTempatUsaha:            in.IDTempatUsaha,
		IDKomoditas:              in.IDKomoditas,
		HargaNormal:              in.HargaNormal,
		HargaMahal:               in.HargaMahal,
		HargaAvg:                 in.HargaAvg,
		SatuanStok:               in.SatuanStok,
		NilaiStok:                in.NilaiStok,
		SatuanPeriode:            in.SatuanPeriode,
		NilaiPeriode:             in.NilaiPeriode,
		LokasiSupplier:           in.LokasiSupplier,
		PolaDistribusi:           in.PolaDistribusi,
		StandardizedStockPeriode: in.StandardizedStockPeriode,
		KelasKomoditas:           in.KelasKomoditas,
		Status:                   int16(in.Status),
	}
}
func (s *serializer) ToPengumpulanData(in entity.PengumpulanData) dto.ResPengumpulanData {
	return dto.ResPengumpulanData{ID: in.ID, IDPasar: in.IDPasar, Tanggal: in.Tanggal, Status: int16(in.Status), Catatan: in.Catatan}
}
func (s *serializer) ToHargaRutin(in entity.HargaRutin) dto.ResHargaRutin {
	return dto.ResHargaRutin{ID: in.ID, IDPengumpulanData: in.IDPengumpulanData, IDTempatUsaha: in.IDTempatUsaha, IDKomoditas: in.IDKomoditas, KelasKomoditas: in.KelasKomoditas, Harga: in.Harga}
}
func (s *serializer) ToHargaPelaporan(in entity.HargaPelaporan) dto.ResHargaPelaporan {
	return dto.ResHargaPelaporan{ID: in.ID, IDPengumpulanData: in.IDPengumpulanData, IDKomoditas: in.IDKomoditas, Tanggal: in.Tanggal, Harga: in.Harga}
}
func (s *serializer) ToPage(in entitybase.BasePaginationResult) dtobase.BasePagination {
	return dtobase.BasePagination{Offset: in.Offset, Limit: in.Limit, Count: in.Count, OrderBy: in.OrderBy}
}

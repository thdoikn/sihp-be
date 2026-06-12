package sihpusecase

import (
	"context"
	"io"

	"github.com/google/uuid"
	"github.com/thdoikn/sihp-be/pkg/dto"
	dtobase "github.com/thdoikn/sihp-be/pkg/dto/base"
)

type SIHPUsecase interface {
	GetPublicOverview(ctx context.Context) dto.ResPublicOverviewEnvelope
	GetPublicKomoditas(ctx context.Context, req *dto.ReqPublicGetKomoditas) dto.ResPublicKomoditasList
	GetPublicKomoditasDetail(ctx context.Context, id uuid.UUID, req *dto.ReqPublicKomoditasDetail) dto.ResPublicKomoditasDetailEnvelope
	GetPublicKomoditasTrend(ctx context.Context, id uuid.UUID, req *dto.ReqPublicKomoditasDetail) dto.ResPublicTrendEnvelope
	GetPublicPasar(ctx context.Context) dto.ResPublicPasarList
	GetPublicPasarDetail(ctx context.Context, id uuid.UUID, req *dto.ReqGetTempatUsaha) dto.ResPublicPasarDetailEnvelope
	GetPublicTempatUsaha(ctx context.Context, req *dto.ReqPublicGetTempatUsaha) dto.ResPublicTempatUsahaList
	GetPublicTempatUsahaDetail(ctx context.Context, id uuid.UUID, req *dto.ReqGetKomoditas) dto.ResPublicTempatUsahaDetailEnvelope

	CreatePasar(ctx context.Context, req *dto.ReqCreatePasar) dto.ResPasarSingle
	GetPasarByID(ctx context.Context, id uuid.UUID) dto.ResPasarSingle
	GetPasarByFilter(ctx context.Context, req *dto.ReqGetPasar) dto.ResPasarList
	UpdatePasar(ctx context.Context, id uuid.UUID, req *dto.ReqUpdatePasar) dto.ResPasarSingle
	DeletePasar(ctx context.Context, id uuid.UUID) dtobase.BaseRes

	CreateKomoditas(ctx context.Context, req *dto.ReqCreateKomoditas) dto.ResKomoditasSingle
	GetKomoditasByID(ctx context.Context, id uuid.UUID) dto.ResKomoditasSingle
	GetKomoditasByFilter(ctx context.Context, req *dto.ReqGetKomoditas) dto.ResKomoditasList
	UpdateKomoditas(ctx context.Context, id uuid.UUID, req *dto.ReqUpdateKomoditas) dto.ResKomoditasSingle
	UploadKomoditasGambar(ctx context.Context, id uuid.UUID, reader io.Reader, size int64, contentType string) dto.ResKomoditasSingle

	CreateTempatUsaha(ctx context.Context, req *dto.ReqCreateTempatUsaha) dto.ResTempatUsahaSingle
	GetTempatUsahaByID(ctx context.Context, id uuid.UUID) dto.ResTempatUsahaSingle
	GetTempatUsahaByFilter(ctx context.Context, req *dto.ReqGetTempatUsaha) dto.ResTempatUsahaList
	UpdateTempatUsaha(ctx context.Context, id uuid.UUID, req *dto.ReqUpdateTempatUsaha) dto.ResTempatUsahaSingle
	DeleteTempatUsaha(ctx context.Context, id uuid.UUID) dtobase.BaseRes

	CreateKomoditasDijual(ctx context.Context, req *dto.ReqCreateKomoditasDijual) dto.ResKomoditasDijualSingle
	GetKomoditasDijualByID(ctx context.Context, id uuid.UUID) dto.ResKomoditasDijualSingle
	GetKomoditasDijualByFilter(ctx context.Context, req *dto.ReqGetKomoditasDijual) dto.ResKomoditasDijualList
	UpdateKomoditasDijual(ctx context.Context, id uuid.UUID, req *dto.ReqUpdateKomoditasDijual) dto.ResKomoditasDijualSingle
	DeleteKomoditasDijual(ctx context.Context, id uuid.UUID) dtobase.BaseRes

	CreatePengumpulanData(ctx context.Context, req *dto.ReqCreatePengumpulanData) dto.ResPengumpulanDataSingle
	GetPengumpulanDataByID(ctx context.Context, id uuid.UUID) dto.ResPengumpulanDataSingle
	GetPengumpulanDataByFilter(ctx context.Context, req *dto.ReqGetPengumpulanData) dto.ResPengumpulanDataList
	UpdatePengumpulanData(ctx context.Context, id uuid.UUID, req *dto.ReqUpdatePengumpulanData) dto.ResPengumpulanDataSingle
	DeletePengumpulanData(ctx context.Context, id uuid.UUID) dtobase.BaseRes
	FinalizePengumpulanData(ctx context.Context, id uuid.UUID) dto.ResFinalizePengumpulanDataEnvelope
	UploadPengumpulanTandaTangan(ctx context.Context, id uuid.UUID, reader io.Reader, size int64, contentType string) dto.ResPengumpulanDataSingle

	CreateHargaRutin(ctx context.Context, req *dto.ReqCreateHargaRutin) dto.ResHargaRutinSingle
	GetHargaRutinByID(ctx context.Context, id uuid.UUID) dto.ResHargaRutinSingle
	GetHargaRutinByFilter(ctx context.Context, req *dto.ReqGetHargaRutin) dto.ResHargaRutinList
	UpdateHargaRutin(ctx context.Context, id uuid.UUID, req *dto.ReqUpdateHargaRutin) dto.ResHargaRutinSingle
	DeleteHargaRutin(ctx context.Context, id uuid.UUID) dtobase.BaseRes

	GetHargaPelaporanByID(ctx context.Context, id uuid.UUID) dto.ResHargaPelaporanSingle
	GetHargaPelaporanByFilter(ctx context.Context, req *dto.ReqGetHargaPelaporan) dto.ResHargaPelaporanList
}

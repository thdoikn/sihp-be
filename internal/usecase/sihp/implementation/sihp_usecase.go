package sihpusecaseimplementation

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/thdoikn/sihp-be/config"
	"github.com/thdoikn/sihp-be/internal/entity"
	entitybase "github.com/thdoikn/sihp-be/internal/entity/base"
	masterdatarepository "github.com/thdoikn/sihp-be/internal/repository/master-data"
	transactiondatarepository "github.com/thdoikn/sihp-be/internal/repository/transaction-data"
	sihpserializer "github.com/thdoikn/sihp-be/internal/serializer/sihp"
	sihpusecase "github.com/thdoikn/sihp-be/internal/usecase/sihp"
	"github.com/thdoikn/sihp-be/pkg/constant"
	"github.com/thdoikn/sihp-be/pkg/dto"
	dtobase "github.com/thdoikn/sihp-be/pkg/dto/base"
	"github.com/thdoikn/sihp-be/pkg/storage"
	miniostorage "github.com/thdoikn/sihp-be/pkg/storage/minio"
	httphelper "github.com/thdoikn/sihp-be/pkg/helper/http"
	queryhelper "github.com/thdoikn/sihp-be/pkg/helper/query"
	"gorm.io/gorm"
)

var _ sihpusecase.SIHPUsecase = (*usecase)(nil)
var _ = entitybase.BasePaginationResult{}
var _ = gorm.ErrRecordNotFound

type usecase struct {
	masterDataRepo      masterdatarepository.MasterDataRepository
	transactionDataRepo transactiondatarepository.TransactionDataRepository
	serializer          sihpserializer.SIHPSerializer
	storage             storage.KomoditasStorage
	cfg                 *config.Config
	validator           *validator.Validate
}

func NewSIHPUsecase(masterDataRepo masterdatarepository.MasterDataRepository, transactionDataRepo transactiondatarepository.TransactionDataRepository, serializer sihpserializer.SIHPSerializer, komoditasStorage storage.KomoditasStorage, cfg *config.Config) sihpusecase.SIHPUsecase {
	return &usecase{
		masterDataRepo:      masterDataRepo,
		transactionDataRepo: transactionDataRepo,
		serializer:          serializer,
		storage:             komoditasStorage,
		cfg:                 cfg,
		validator:           validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (u *usecase) baseRes(code int, message string, err error) dtobase.BaseRes {
	return dtobase.BaseRes{Success: code >= 200 && code < 300, Code: code, Message: message, Stacktrace: httphelper.Stacktrace(u.cfg, err)}
}

func (u *usecase) GetPublicOverview(ctx context.Context) dto.ResPublicOverviewEnvelope {
	pasarCount, tuCount, komCount, err := u.masterDataRepo.GetOverview(ctx)
	if err != nil {
		return dto.ResPublicOverviewEnvelope{BaseRes: u.baseRes(http.StatusInternalServerError, err.Error(), err)}
	}
	return dto.ResPublicOverviewEnvelope{BaseRes: u.baseRes(http.StatusOK, "ok", nil), Data: &dto.ResPublicOverview{PasarActiveCount: pasarCount, TempatUsahaActiveCount: tuCount, KomoditasCount: komCount}}
}

func (u *usecase) komoditasFilter(req *dto.ReqGetKomoditas) entity.KomoditasFilter {
	return entity.KomoditasFilter{Name: req.Name, IDTempatUsaha: req.IDTempatUsaha, IDPasar: req.IDPasar, PaginationFilter: queryhelper.SerializeFilterPaginationDtoToEntity(req.BaseReqQueryPagination)}
}
func (u *usecase) pasarFilter(req *dto.ReqGetPasar, defaultActive bool) entity.PasarFilter {
	f := entity.PasarFilter{Name: req.Name, PaginationFilter: queryhelper.SerializeFilterPaginationDtoToEntity(req.BaseReqQueryPagination)}
	if req.Status != nil {
		s := constant.ActiveInactiveStatus(*req.Status)
		f.Status = &s
	} else if defaultActive {
		s := constant.StatusActive
		f.Status = &s
	}
	return f
}
func (u *usecase) tempatUsahaFilter(req *dto.ReqGetTempatUsaha, defaultActive bool) entity.TempatUsahaFilter {
	f := entity.TempatUsahaFilter{Name: req.Name, IDPasar: req.IDPasar, PaginationFilter: queryhelper.SerializeFilterPaginationDtoToEntity(req.BaseReqQueryPagination)}
	if req.Status != nil {
		s := constant.ActiveInactiveStatus(*req.Status)
		f.Status = &s
	} else if defaultActive {
		s := constant.StatusActive
		f.Status = &s
	}
	return f
}

func (u *usecase) GetPublicKomoditas(ctx context.Context, req *dto.ReqGetKomoditas) dto.ResKomoditasList {
	return u.GetKomoditasByFilter(ctx, req)
}
func (u *usecase) GetPublicKomoditasDetail(ctx context.Context, id uuid.UUID, req *dto.ReqPublicKomoditasDetail) dto.ResPublicKomoditasDetailEnvelope {
	days := 30
	if req != nil && req.Days != nil && *req.Days > 0 {
		days = *req.Days
	}
	komoditas, err := u.masterDataRepo.GetKomoditasByID(ctx, id)
	if err != nil {
		return dto.ResPublicKomoditasDetailEnvelope{BaseRes: u.baseRes(http.StatusNotFound, "komoditas not found", err)}
	}
	latestTanggal, latestHarga, avg, min, max, err := u.masterDataRepo.GetPublicKomoditasStats(ctx, id, days, req.IDPasar)
	if err != nil {
		return dto.ResPublicKomoditasDetailEnvelope{BaseRes: u.baseRes(http.StatusInternalServerError, err.Error(), err)}
	}
	return dto.ResPublicKomoditasDetailEnvelope{BaseRes: u.baseRes(http.StatusOK, "ok", nil), Data: &dto.ResPublicKomoditasDetail{Komoditas: u.serializer.ToKomoditas(*komoditas), Latest: dto.ResPublicHargaStat{Tanggal: latestTanggal, HargaRataRata: latestHarga}, AvgND: avg, MinND: min, MaxND: max, Days: days}}
}
func (u *usecase) GetPublicKomoditasTrend(ctx context.Context, id uuid.UUID, req *dto.ReqPublicKomoditasDetail) dto.ResPublicTrendEnvelope {
	days := 30
	if req != nil && req.Days != nil && *req.Days > 0 {
		days = *req.Days
	}
	rows, err := u.masterDataRepo.GetPublicKomoditasTrend(ctx, id, days, req.IDPasar)
	if err != nil {
		return dto.ResPublicTrendEnvelope{BaseRes: u.baseRes(http.StatusInternalServerError, err.Error(), err)}
	}
	out := make([]dto.ResPublicTrendPoint, 0, len(rows))
	for _, row := range rows {
		tanggal, _ := row["tanggal"].(time.Time)
		harga, _ := row["harga_rata_rata"].(float64)
		out = append(out, dto.ResPublicTrendPoint{Tanggal: tanggal, HargaRataRata: harga})
	}
	return dto.ResPublicTrendEnvelope{BaseRes: u.baseRes(http.StatusOK, "ok", nil), Data: out}
}
func (u *usecase) GetPublicPasar(ctx context.Context, req *dto.ReqGetPasar) dto.ResPasarList {
	filter := u.pasarFilter(req, true)
	rows, page, err := u.masterDataRepo.GetPasarByFilter(ctx, &filter)
	if err != nil {
		return dto.ResPasarList{BaseResPagination: dtobase.BaseResPagination{BaseRes: u.baseRes(http.StatusInternalServerError, err.Error(), err)}}
	}
	res := make([]dto.ResPasar, 0, len(rows))
	for _, row := range rows {
		res = append(res, u.serializer.ToPasar(row))
	}
	return dto.ResPasarList{BaseResPagination: dtobase.BaseResPagination{BaseRes: u.baseRes(http.StatusOK, "ok", nil), Page: u.serializer.ToPage(page)}, Data: res}
}
func (u *usecase) GetPublicPasarDetail(ctx context.Context, id uuid.UUID, req *dto.ReqGetTempatUsaha) dto.ResPublicPasarDetailEnvelope {
	filter := u.tempatUsahaFilter(req, true)
	pasar, tus, page, err := u.masterDataRepo.GetPublicPasarDetail(ctx, id, &filter)
	if err != nil {
		return dto.ResPublicPasarDetailEnvelope{BaseRes: u.baseRes(http.StatusInternalServerError, err.Error(), err)}
	}
	arr := make([]dto.ResTempatUsaha, 0, len(tus))
	for _, row := range tus {
		arr = append(arr, u.serializer.ToTempatUsaha(row))
	}
	return dto.ResPublicPasarDetailEnvelope{BaseRes: u.baseRes(http.StatusOK, "ok", nil), Data: &dto.ResPublicPasarDetail{Pasar: u.serializer.ToPasar(*pasar), TempatUsaha: arr, Page: u.serializer.ToPage(page)}}
}
func (u *usecase) GetPublicTempatUsaha(ctx context.Context, req *dto.ReqGetTempatUsaha) dto.ResTempatUsahaList {
	return u.GetTempatUsahaByFilter(ctx, req)
}
func (u *usecase) GetPublicTempatUsahaDetail(ctx context.Context, id uuid.UUID, req *dto.ReqGetKomoditas) dto.ResPublicTempatUsahaDetailEnvelope {
	filter := u.komoditasFilter(req)
	tu, pasar, komoditasRows, page, err := u.masterDataRepo.GetPublicTempatUsahaDetail(ctx, id, filter)
	if err != nil {
		return dto.ResPublicTempatUsahaDetailEnvelope{BaseRes: u.baseRes(http.StatusInternalServerError, err.Error(), err)}
	}
	kom := make([]dto.ResPublicTempatUsahaKomoditas, 0, len(komoditasRows))
	for _, row := range komoditasRows {
		item := dto.ResPublicTempatUsahaKomoditas{ResKomoditas: dto.ResKomoditas{ID: row["id"].(uuid.UUID), Nama: row["nama"].(string)}}
		if v, ok := row["satuan"]; ok {
			if val, ok2 := v.(*string); ok2 {
				item.Satuan = val
			}
		}
		if v, ok := row["gambar_url"]; ok {
			if val, ok2 := v.(*string); ok2 {
				item.GambarURL = val
			}
		}
		if latestMap, ok := row["latest"].(map[string]any); ok {
			if tanggal, ok2 := latestMap["tanggal"].(time.Time); ok2 {
				item.Latest.Tanggal = &tanggal
			}
			if harga, ok2 := latestMap["harga_rata_rata"].(float64); ok2 {
				item.Latest.HargaRataRata = &harga
			}
		}
		kom = append(kom, item)
	}
	return dto.ResPublicTempatUsahaDetailEnvelope{BaseRes: u.baseRes(http.StatusOK, "ok", nil), Data: &dto.ResPublicTempatUsahaDetail{TempatUsaha: u.serializer.ToTempatUsaha(*tu), Pasar: u.serializer.ToPasar(*pasar), Komoditas: kom, Page: u.serializer.ToPage(page)}}
}

func (u *usecase) CreatePasar(ctx context.Context, req *dto.ReqCreatePasar) dto.ResPasarSingle {
	if err := u.validator.Struct(req); err != nil {
		return dto.ResPasarSingle{BaseRes: u.baseRes(http.StatusBadRequest, err.Error(), err)}
	}
	obj, err := u.masterDataRepo.CreatePasar(ctx, &entity.Pasar{Nama: req.Nama, Alamat: req.Alamat, Status: constant.StatusActive})
	if err != nil {
		return dto.ResPasarSingle{BaseRes: u.baseRes(http.StatusInternalServerError, err.Error(), err)}
	}
	res := u.serializer.ToPasar(*obj)
	return dto.ResPasarSingle{BaseRes: u.baseRes(http.StatusCreated, "created", nil), Data: &res}
}
func (u *usecase) GetPasarByID(ctx context.Context, id uuid.UUID) dto.ResPasarSingle {
	obj, err := u.masterDataRepo.GetPasarByID(ctx, id)
	if err != nil {
		return dto.ResPasarSingle{BaseRes: u.baseRes(http.StatusNotFound, "not found", err)}
	}
	res := u.serializer.ToPasar(*obj)
	return dto.ResPasarSingle{BaseRes: u.baseRes(http.StatusOK, "ok", nil), Data: &res}
}
func (u *usecase) GetPasarByFilter(ctx context.Context, req *dto.ReqGetPasar) dto.ResPasarList {
	filter := u.pasarFilter(req, false)
	rows, page, err := u.masterDataRepo.GetPasarByFilter(ctx, &filter)
	if err != nil {
		return dto.ResPasarList{BaseResPagination: dtobase.BaseResPagination{BaseRes: u.baseRes(http.StatusInternalServerError, err.Error(), err)}}
	}
	res := make([]dto.ResPasar, 0, len(rows))
	for _, row := range rows {
		res = append(res, u.serializer.ToPasar(row))
	}
	return dto.ResPasarList{BaseResPagination: dtobase.BaseResPagination{BaseRes: u.baseRes(http.StatusOK, "ok", nil), Page: u.serializer.ToPage(page)}, Data: res}
}
func (u *usecase) UpdatePasar(ctx context.Context, id uuid.UUID, req *dto.ReqUpdatePasar) dto.ResPasarSingle {
	update := map[string]any{}
	if req.Nama != nil {
		update["nama"] = *req.Nama
	}
	if req.Alamat != nil {
		update["alamat"] = req.Alamat
	}
	if req.Status != nil {
		update["status"] = *req.Status
	}
	obj, err := u.masterDataRepo.UpdatePasar(ctx, id, update)
	if err != nil {
		return dto.ResPasarSingle{BaseRes: u.baseRes(http.StatusInternalServerError, err.Error(), err)}
	}
	res := u.serializer.ToPasar(*obj)
	return dto.ResPasarSingle{BaseRes: u.baseRes(http.StatusOK, "updated", nil), Data: &res}
}
func (u *usecase) DeletePasar(ctx context.Context, id uuid.UUID) dtobase.BaseRes {
	_, err := u.masterDataRepo.UpdatePasar(ctx, id, map[string]any{"status": constant.StatusInactive})
	if err != nil {
		return u.baseRes(http.StatusInternalServerError, err.Error(), err)
	}
	return u.baseRes(http.StatusOK, "deleted", nil)
}

func (u *usecase) CreateKomoditas(ctx context.Context, req *dto.ReqCreateKomoditas) dto.ResKomoditasSingle {
	if err := u.validator.Struct(req); err != nil {
		return dto.ResKomoditasSingle{BaseRes: u.baseRes(http.StatusBadRequest, err.Error(), err)}
	}
	obj, err := u.masterDataRepo.CreateKomoditas(ctx, &entity.Komoditas{Nama: req.Nama, Satuan: req.Satuan})
	if err != nil {
		return dto.ResKomoditasSingle{BaseRes: u.baseRes(http.StatusInternalServerError, err.Error(), err)}
	}
	res := u.serializer.ToKomoditas(*obj)
	return dto.ResKomoditasSingle{BaseRes: u.baseRes(http.StatusCreated, "created", nil), Data: &res}
}
func (u *usecase) GetKomoditasByID(ctx context.Context, id uuid.UUID) dto.ResKomoditasSingle {
	obj, err := u.masterDataRepo.GetKomoditasByID(ctx, id)
	if err != nil {
		return dto.ResKomoditasSingle{BaseRes: u.baseRes(http.StatusNotFound, "not found", err)}
	}
	res := u.serializer.ToKomoditas(*obj)
	return dto.ResKomoditasSingle{BaseRes: u.baseRes(http.StatusOK, "ok", nil), Data: &res}
}
func (u *usecase) GetKomoditasByFilter(ctx context.Context, req *dto.ReqGetKomoditas) dto.ResKomoditasList {
	f := u.komoditasFilter(req)
	rows, page, err := u.masterDataRepo.GetKomoditasByFilter(ctx, &f)
	if err != nil {
		return dto.ResKomoditasList{BaseResPagination: dtobase.BaseResPagination{BaseRes: u.baseRes(http.StatusInternalServerError, err.Error(), err)}}
	}
	res := make([]dto.ResKomoditas, 0, len(rows))
	for _, row := range rows {
		res = append(res, u.serializer.ToKomoditas(row))
	}
	return dto.ResKomoditasList{BaseResPagination: dtobase.BaseResPagination{BaseRes: u.baseRes(http.StatusOK, "ok", nil), Page: u.serializer.ToPage(page)}, Data: res}
}
func (u *usecase) UpdateKomoditas(ctx context.Context, id uuid.UUID, req *dto.ReqUpdateKomoditas) dto.ResKomoditasSingle {
	update := map[string]any{}
	if req.Nama != nil {
		update["nama"] = *req.Nama
	}
	if req.Satuan != nil {
		update["satuan"] = req.Satuan
	}
	obj, err := u.masterDataRepo.UpdateKomoditas(ctx, id, update)
	if err != nil {
		return dto.ResKomoditasSingle{BaseRes: u.baseRes(http.StatusInternalServerError, err.Error(), err)}
	}
	res := u.serializer.ToKomoditas(*obj)
	return dto.ResKomoditasSingle{BaseRes: u.baseRes(http.StatusOK, "updated", nil), Data: &res}
}

const maxKomoditasImageSize = 2 * 1024 * 1024 // 2MB

func (u *usecase) UploadKomoditasGambar(ctx context.Context, id uuid.UUID, reader io.Reader, size int64, contentType string) dto.ResKomoditasSingle {
	if u.storage == nil {
		return dto.ResKomoditasSingle{BaseRes: u.baseRes(http.StatusServiceUnavailable, "storage not configured", errors.New("storage not configured"))}
	}
	if size <= 0 || size > maxKomoditasImageSize {
		return dto.ResKomoditasSingle{BaseRes: u.baseRes(http.StatusBadRequest, "file size must be between 1 byte and 2MB", errors.New("invalid file size"))}
	}

	existing, err := u.masterDataRepo.GetKomoditasByID(ctx, id)
	if err != nil {
		return dto.ResKomoditasSingle{BaseRes: u.baseRes(http.StatusNotFound, "komoditas not found", err)}
	}

	publicURL, err := u.storage.UploadKomoditasImage(ctx, id, reader, size, contentType)
	if err != nil {
		return dto.ResKomoditasSingle{BaseRes: u.baseRes(http.StatusBadRequest, err.Error(), err)}
	}

	if existing.GambarURL != nil && *existing.GambarURL != "" {
		oldKey := miniostorage.ObjectKeyFromURL(*existing.GambarURL, u.cfg.MinIO.PublicURL)
		_ = u.storage.DeleteKomoditasImage(ctx, oldKey)
	}

	obj, err := u.masterDataRepo.UpdateKomoditas(ctx, id, map[string]any{"gambar_url": publicURL})
	if err != nil {
		return dto.ResKomoditasSingle{BaseRes: u.baseRes(http.StatusInternalServerError, err.Error(), err)}
	}

	res := u.serializer.ToKomoditas(*obj)
	return dto.ResKomoditasSingle{BaseRes: u.baseRes(http.StatusOK, "uploaded", nil), Data: &res}
}

func (u *usecase) CreateTempatUsaha(ctx context.Context, req *dto.ReqCreateTempatUsaha) dto.ResTempatUsahaSingle {
	if err := u.validator.Struct(req); err != nil {
		return dto.ResTempatUsahaSingle{BaseRes: u.baseRes(http.StatusBadRequest, err.Error(), err)}
	}
	obj, err := u.masterDataRepo.CreateTempatUsaha(ctx, &entity.TempatUsaha{IDPasar: req.IDPasar, Nama: req.Nama, Pemilik: req.Pemilik, Status: constant.StatusActive})
	if err != nil {
		return dto.ResTempatUsahaSingle{BaseRes: u.baseRes(http.StatusInternalServerError, err.Error(), err)}
	}
	res := u.serializer.ToTempatUsaha(*obj)
	return dto.ResTempatUsahaSingle{BaseRes: u.baseRes(http.StatusCreated, "created", nil), Data: &res}
}
func (u *usecase) GetTempatUsahaByID(ctx context.Context, id uuid.UUID) dto.ResTempatUsahaSingle {
	obj, err := u.masterDataRepo.GetTempatUsahaByID(ctx, id)
	if err != nil {
		return dto.ResTempatUsahaSingle{BaseRes: u.baseRes(http.StatusNotFound, "not found", err)}
	}
	res := u.serializer.ToTempatUsaha(*obj)
	return dto.ResTempatUsahaSingle{BaseRes: u.baseRes(http.StatusOK, "ok", nil), Data: &res}
}
func (u *usecase) GetTempatUsahaByFilter(ctx context.Context, req *dto.ReqGetTempatUsaha) dto.ResTempatUsahaList {
	f := u.tempatUsahaFilter(req, false)
	rows, page, err := u.masterDataRepo.GetTempatUsahaByFilter(ctx, &f)
	if err != nil {
		return dto.ResTempatUsahaList{BaseResPagination: dtobase.BaseResPagination{BaseRes: u.baseRes(http.StatusInternalServerError, err.Error(), err)}}
	}
	res := make([]dto.ResTempatUsaha, 0, len(rows))
	for _, row := range rows {
		res = append(res, u.serializer.ToTempatUsaha(row))
	}
	return dto.ResTempatUsahaList{BaseResPagination: dtobase.BaseResPagination{BaseRes: u.baseRes(http.StatusOK, "ok", nil), Page: u.serializer.ToPage(page)}, Data: res}
}
func (u *usecase) UpdateTempatUsaha(ctx context.Context, id uuid.UUID, req *dto.ReqUpdateTempatUsaha) dto.ResTempatUsahaSingle {
	update := map[string]any{}
	if req.IDPasar != nil {
		update["id_pasar"] = *req.IDPasar
	}
	if req.Nama != nil {
		update["nama"] = *req.Nama
	}
	if req.Pemilik != nil {
		update["pemilik"] = req.Pemilik
	}
	if req.Status != nil {
		update["status"] = *req.Status
	}
	obj, err := u.masterDataRepo.UpdateTempatUsaha(ctx, id, update)
	if err != nil {
		return dto.ResTempatUsahaSingle{BaseRes: u.baseRes(http.StatusInternalServerError, err.Error(), err)}
	}
	res := u.serializer.ToTempatUsaha(*obj)
	return dto.ResTempatUsahaSingle{BaseRes: u.baseRes(http.StatusOK, "updated", nil), Data: &res}
}

func (u *usecase) DeleteTempatUsaha(ctx context.Context, id uuid.UUID) dtobase.BaseRes {
	_, err := u.masterDataRepo.UpdateTempatUsaha(ctx, id, map[string]any{"status": constant.StatusInactive})
	if err != nil {
		return u.baseRes(http.StatusInternalServerError, err.Error(), err)
	}
	return u.baseRes(http.StatusOK, "deleted", nil)
}

func komoditasDijualFromCreateReq(req *dto.ReqCreateKomoditasDijual) entity.KomoditasDijual {
	satuanStok := req.SatuanStok
	if satuanStok == "" {
		satuanStok = "kg"
	}
	satuanPeriode := req.SatuanPeriode
	if satuanPeriode == "" {
		satuanPeriode = "minggu"
	}
	nilaiPeriode := req.NilaiPeriode
	if nilaiPeriode <= 0 {
		nilaiPeriode = 1
	}
	var standardized float64
	if req.StandardizedStockPeriode != nil {
		standardized = *req.StandardizedStockPeriode
	}
	status := constant.StatusActive
	if req.Status != nil {
		status = constant.ActiveInactiveStatus(*req.Status)
	}
	return entity.KomoditasDijual{
		IDTempatUsaha:            req.IDTempatUsaha,
		IDKomoditas:              req.IDKomoditas,
		HargaNormal:              req.HargaNormal,
		HargaMahal:               req.HargaMahal,
		SatuanStok:               satuanStok,
		NilaiStok:                req.NilaiStok,
		SatuanPeriode:            satuanPeriode,
		NilaiPeriode:             nilaiPeriode,
		LokasiSupplier:           req.LokasiSupplier,
		PolaDistribusi:           req.PolaDistribusi,
		StandardizedStockPeriode: standardized,
		Status:                   status,
	}
}

func (u *usecase) CreateKomoditasDijual(ctx context.Context, req *dto.ReqCreateKomoditasDijual) dto.ResKomoditasDijualSingle {
	if err := u.validator.Struct(req); err != nil {
		return dto.ResKomoditasDijualSingle{BaseRes: u.baseRes(http.StatusBadRequest, err.Error(), err)}
	}
	input := komoditasDijualFromCreateReq(req)
	obj, err := u.masterDataRepo.CreateKomoditasDijual(ctx, &input)
	if err != nil {
		return dto.ResKomoditasDijualSingle{BaseRes: u.baseRes(http.StatusInternalServerError, err.Error(), err)}
	}
	res := u.serializer.ToKomoditasDijual(*obj)
	return dto.ResKomoditasDijualSingle{BaseRes: u.baseRes(http.StatusCreated, "created", nil), Data: &res}
}
func (u *usecase) GetKomoditasDijualByID(ctx context.Context, id uuid.UUID) dto.ResKomoditasDijualSingle {
	obj, err := u.masterDataRepo.GetKomoditasDijualByID(ctx, id)
	if err != nil {
		return dto.ResKomoditasDijualSingle{BaseRes: u.baseRes(http.StatusNotFound, "not found", err)}
	}
	res := u.serializer.ToKomoditasDijual(*obj)
	return dto.ResKomoditasDijualSingle{BaseRes: u.baseRes(http.StatusOK, "ok", nil), Data: &res}
}
func (u *usecase) GetKomoditasDijualByFilter(ctx context.Context, req *dto.ReqGetKomoditasDijual) dto.ResKomoditasDijualList {
	f := entity.KomoditasDijualFilter{IDTempatUsaha: req.IDTempatUsaha, IDKomoditas: req.IDKomoditas, PaginationFilter: queryhelper.SerializeFilterPaginationDtoToEntity(req.BaseReqQueryPagination)}
	if req.Status != nil {
		s := constant.ActiveInactiveStatus(*req.Status)
		f.Status = &s
	} else {
		s := constant.StatusActive
		f.Status = &s
	}
	rows, page, err := u.masterDataRepo.GetKomoditasDijualByFilter(ctx, &f)
	if err != nil {
		return dto.ResKomoditasDijualList{BaseResPagination: dtobase.BaseResPagination{BaseRes: u.baseRes(http.StatusInternalServerError, err.Error(), err)}}
	}
	res := make([]dto.ResKomoditasDijual, 0, len(rows))
	for _, row := range rows {
		res = append(res, u.serializer.ToKomoditasDijual(row))
	}
	return dto.ResKomoditasDijualList{BaseResPagination: dtobase.BaseResPagination{BaseRes: u.baseRes(http.StatusOK, "ok", nil), Page: u.serializer.ToPage(page)}, Data: res}
}

func (u *usecase) UpdateKomoditasDijual(ctx context.Context, id uuid.UUID, req *dto.ReqUpdateKomoditasDijual) dto.ResKomoditasDijualSingle {
	_, err := u.masterDataRepo.GetKomoditasDijualByID(ctx, id)
	if err != nil {
		return dto.ResKomoditasDijualSingle{BaseRes: u.baseRes(http.StatusNotFound, "not found", err)}
	}

	update := map[string]any{}
	if req.HargaNormal != nil {
		update["harga_normal"] = *req.HargaNormal
	}
	if req.HargaMahal != nil {
		update["harga_mahal"] = *req.HargaMahal
	}
	if req.SatuanStok != nil {
		update["satuan_stok"] = *req.SatuanStok
	}
	if req.NilaiStok != nil {
		update["nilai_stok"] = *req.NilaiStok
	}
	if req.SatuanPeriode != nil {
		update["satuan_periode"] = *req.SatuanPeriode
	}
	if req.NilaiPeriode != nil {
		update["nilai_periode"] = *req.NilaiPeriode
	}
	if req.LokasiSupplier != nil {
		update["lokasi_supplier"] = *req.LokasiSupplier
	}
	if req.PolaDistribusi != nil {
		update["pola_distribusi"] = req.PolaDistribusi
	}
	if req.KelasKomoditas != nil {
		update["kelas_komoditas"] = req.KelasKomoditas
	}
	if req.Status != nil {
		update["status"] = *req.Status
	}

	if req.StandardizedStockPeriode != nil {
		update["standardized_stock_periode"] = *req.StandardizedStockPeriode
	}

	obj, err := u.masterDataRepo.UpdateKomoditasDijual(ctx, id, update)
	if err != nil {
		return dto.ResKomoditasDijualSingle{BaseRes: u.baseRes(http.StatusInternalServerError, err.Error(), err)}
	}
	res := u.serializer.ToKomoditasDijual(*obj)
	return dto.ResKomoditasDijualSingle{BaseRes: u.baseRes(http.StatusOK, "updated", nil), Data: &res}
}

func (u *usecase) DeleteKomoditasDijual(ctx context.Context, id uuid.UUID) dtobase.BaseRes {
	_, err := u.masterDataRepo.UpdateKomoditasDijual(ctx, id, map[string]any{"status": constant.StatusInactive})
	if err != nil {
		return u.baseRes(http.StatusInternalServerError, err.Error(), err)
	}
	return u.baseRes(http.StatusOK, "deleted", nil)
}

func (u *usecase) CreatePengumpulanData(ctx context.Context, req *dto.ReqCreatePengumpulanData) dto.ResPengumpulanDataSingle {
	if err := u.validator.Struct(req); err != nil {
		return dto.ResPengumpulanDataSingle{BaseRes: u.baseRes(http.StatusBadRequest, err.Error(), err)}
	}
	obj, err := u.transactionDataRepo.CreatePengumpulanData(ctx, &entity.PengumpulanData{IDPasar: req.IDPasar, Tanggal: req.Tanggal, Status: constant.PengumpulanDataDraft, Catatan: req.Catatan})
	if err != nil {
		return dto.ResPengumpulanDataSingle{BaseRes: u.baseRes(http.StatusInternalServerError, err.Error(), err)}
	}
	res := u.serializer.ToPengumpulanData(*obj)
	return dto.ResPengumpulanDataSingle{BaseRes: u.baseRes(http.StatusCreated, "created", nil), Data: &res}
}

func (u *usecase) GetPengumpulanDataByID(ctx context.Context, id uuid.UUID) dto.ResPengumpulanDataSingle {
	obj, err := u.transactionDataRepo.GetPengumpulanDataByID(ctx, id)
	if err != nil {
		return dto.ResPengumpulanDataSingle{BaseRes: u.baseRes(http.StatusNotFound, "not found", err)}
	}
	res := u.serializer.ToPengumpulanData(*obj)
	return dto.ResPengumpulanDataSingle{BaseRes: u.baseRes(http.StatusOK, "ok", nil), Data: &res}
}

func (u *usecase) GetPengumpulanDataByFilter(ctx context.Context, req *dto.ReqGetPengumpulanData) dto.ResPengumpulanDataList {
	f := entity.PengumpulanDataFilter{IDPasar: req.IDPasar, From: req.From, To: req.To, PaginationFilter: queryhelper.SerializeFilterPaginationDtoToEntity(req.BaseReqQueryPagination)}
	if req.Status != nil {
		s := constant.PengumpulanDataStatus(*req.Status)
		f.Status = &s
	}
	rows, page, err := u.transactionDataRepo.GetPengumpulanDataByFilter(ctx, &f)
	if err != nil {
		return dto.ResPengumpulanDataList{BaseResPagination: dtobase.BaseResPagination{BaseRes: u.baseRes(http.StatusInternalServerError, err.Error(), err)}}
	}
	res := make([]dto.ResPengumpulanData, 0, len(rows))
	for _, row := range rows {
		res = append(res, u.serializer.ToPengumpulanData(row))
	}
	return dto.ResPengumpulanDataList{BaseResPagination: dtobase.BaseResPagination{BaseRes: u.baseRes(http.StatusOK, "ok", nil), Page: u.serializer.ToPage(page)}, Data: res}
}

func (u *usecase) UpdatePengumpulanData(ctx context.Context, id uuid.UUID, req *dto.ReqUpdatePengumpulanData) dto.ResPengumpulanDataSingle {
	current, err := u.transactionDataRepo.GetPengumpulanDataByID(ctx, id)
	if err != nil {
		return dto.ResPengumpulanDataSingle{BaseRes: u.baseRes(http.StatusNotFound, "not found", err)}
	}
	if current.Status != constant.PengumpulanDataDraft {
		e := errors.New("only draft can be updated")
		return dto.ResPengumpulanDataSingle{BaseRes: u.baseRes(http.StatusBadRequest, e.Error(), e)}
	}
	update := map[string]any{}
	if req.Tanggal != nil {
		update["tanggal"] = *req.Tanggal
	}
	if req.Catatan != nil {
		update["catatan"] = req.Catatan
	}
	obj, err := u.transactionDataRepo.UpdatePengumpulanData(ctx, id, update)
	if err != nil {
		return dto.ResPengumpulanDataSingle{BaseRes: u.baseRes(http.StatusInternalServerError, err.Error(), err)}
	}
	res := u.serializer.ToPengumpulanData(*obj)
	return dto.ResPengumpulanDataSingle{BaseRes: u.baseRes(http.StatusOK, "updated", nil), Data: &res}
}

func (u *usecase) DeletePengumpulanData(ctx context.Context, id uuid.UUID) dtobase.BaseRes {
	current, err := u.transactionDataRepo.GetPengumpulanDataByID(ctx, id)
	if err != nil {
		return u.baseRes(http.StatusNotFound, "not found", err)
	}
	if current.Status != constant.PengumpulanDataDraft {
		e := errors.New("only draft can be deleted")
		return u.baseRes(http.StatusBadRequest, e.Error(), e)
	}
	if err := u.transactionDataRepo.DeletePengumpulanData(ctx, id); err != nil {
		return u.baseRes(http.StatusInternalServerError, err.Error(), err)
	}
	return u.baseRes(http.StatusOK, "deleted", nil)
}

func (u *usecase) FinalizePengumpulanData(ctx context.Context, id uuid.UUID) dto.ResFinalizePengumpulanDataEnvelope {
	count, err := u.transactionDataRepo.FinalizePengumpulanData(ctx, id)
	if err != nil {
		return dto.ResFinalizePengumpulanDataEnvelope{BaseRes: u.baseRes(http.StatusBadRequest, err.Error(), err)}
	}
	return dto.ResFinalizePengumpulanDataEnvelope{BaseRes: u.baseRes(http.StatusOK, "finalized", nil), Data: &dto.ResFinalizePengumpulanData{FinalizedKomoditasCount: count}}
}

func (u *usecase) ensureHargaRutinDraft(ctx context.Context, idPengumpulanData uuid.UUID) error {
	pd, err := u.transactionDataRepo.GetPengumpulanDataByID(ctx, idPengumpulanData)
	if err != nil {
		return err
	}
	if pd.Status != constant.PengumpulanDataDraft {
		return errors.New("parent pengumpulan_data must be draft")
	}
	return nil
}
func (u *usecase) CreateHargaRutin(ctx context.Context, req *dto.ReqCreateHargaRutin) dto.ResHargaRutinSingle {
	if err := u.validator.Struct(req); err != nil {
		return dto.ResHargaRutinSingle{BaseRes: u.baseRes(http.StatusBadRequest, err.Error(), err)}
	}
	if req.KelasKomoditas == "" {
		e := errors.New("kelas_komoditas is required")
		return dto.ResHargaRutinSingle{BaseRes: u.baseRes(http.StatusBadRequest, e.Error(), e)}
	}
	if err := u.ensureHargaRutinDraft(ctx, req.IDPengumpulanData); err != nil {
		return dto.ResHargaRutinSingle{BaseRes: u.baseRes(http.StatusBadRequest, err.Error(), err)}
	}
	obj, err := u.transactionDataRepo.CreateHargaRutin(ctx, &entity.HargaRutin{IDPengumpulanData: req.IDPengumpulanData, IDTempatUsaha: req.IDTempatUsaha, IDKomoditas: req.IDKomoditas, KelasKomoditas: req.KelasKomoditas, Harga: req.Harga})
	if err != nil {
		return dto.ResHargaRutinSingle{BaseRes: u.baseRes(http.StatusInternalServerError, err.Error(), err)}
	}
	res := u.serializer.ToHargaRutin(*obj)
	return dto.ResHargaRutinSingle{BaseRes: u.baseRes(http.StatusCreated, "created", nil), Data: &res}
}

func (u *usecase) GetHargaRutinByID(ctx context.Context, id uuid.UUID) dto.ResHargaRutinSingle {
	obj, err := u.transactionDataRepo.GetHargaRutinByID(ctx, id)
	if err != nil {
		return dto.ResHargaRutinSingle{BaseRes: u.baseRes(http.StatusNotFound, "not found", err)}
	}
	res := u.serializer.ToHargaRutin(*obj)
	return dto.ResHargaRutinSingle{BaseRes: u.baseRes(http.StatusOK, "ok", nil), Data: &res}
}
func (u *usecase) GetHargaRutinByFilter(ctx context.Context, req *dto.ReqGetHargaRutin) dto.ResHargaRutinList {
	f := entity.HargaRutinFilter{IDPengumpulanData: req.IDPengumpulanData, IDKomoditas: req.IDKomoditas, IDTempatUsaha: req.IDTempatUsaha, PaginationFilter: queryhelper.SerializeFilterPaginationDtoToEntity(req.BaseReqQueryPagination)}
	rows, page, err := u.transactionDataRepo.GetHargaRutinByFilter(ctx, &f)
	if err != nil {
		return dto.ResHargaRutinList{BaseResPagination: dtobase.BaseResPagination{BaseRes: u.baseRes(http.StatusInternalServerError, err.Error(), err)}}
	}
	res := make([]dto.ResHargaRutin, 0, len(rows))
	for _, row := range rows {
		res = append(res, u.serializer.ToHargaRutin(row))
	}
	return dto.ResHargaRutinList{BaseResPagination: dtobase.BaseResPagination{BaseRes: u.baseRes(http.StatusOK, "ok", nil), Page: u.serializer.ToPage(page)}, Data: res}
}

func (u *usecase) UpdateHargaRutin(ctx context.Context, id uuid.UUID, req *dto.ReqUpdateHargaRutin) dto.ResHargaRutinSingle {
	current, err := u.transactionDataRepo.GetHargaRutinByID(ctx, id)
	if err != nil {
		return dto.ResHargaRutinSingle{BaseRes: u.baseRes(http.StatusNotFound, "not found", err)}
	}
	if err := u.ensureHargaRutinDraft(ctx, current.IDPengumpulanData); err != nil {
		return dto.ResHargaRutinSingle{BaseRes: u.baseRes(http.StatusBadRequest, err.Error(), err)}
	}
	update := map[string]any{}
	if req.IDTempatUsaha != nil {
		update["id_tempat_usaha"] = *req.IDTempatUsaha
	}
	if req.KelasKomoditas != nil {
		if *req.KelasKomoditas == "" {
			e := errors.New("kelas_komoditas cannot be empty")
			return dto.ResHargaRutinSingle{BaseRes: u.baseRes(http.StatusBadRequest, e.Error(), e)}
		}
		update["kelas_komoditas"] = *req.KelasKomoditas
	}
	if req.Harga != nil {
		update["harga"] = *req.Harga
	}
	obj, err := u.transactionDataRepo.UpdateHargaRutin(ctx, id, update)
	if err != nil {
		return dto.ResHargaRutinSingle{BaseRes: u.baseRes(http.StatusInternalServerError, err.Error(), err)}
	}
	res := u.serializer.ToHargaRutin(*obj)
	return dto.ResHargaRutinSingle{BaseRes: u.baseRes(http.StatusOK, "updated", nil), Data: &res}
}

func (u *usecase) DeleteHargaRutin(ctx context.Context, id uuid.UUID) dtobase.BaseRes {
	current, err := u.transactionDataRepo.GetHargaRutinByID(ctx, id)
	if err != nil {
		return u.baseRes(http.StatusNotFound, "not found", err)
	}
	if err := u.ensureHargaRutinDraft(ctx, current.IDPengumpulanData); err != nil {
		return u.baseRes(http.StatusBadRequest, err.Error(), err)
	}
	if err := u.transactionDataRepo.DeleteHargaRutin(ctx, id); err != nil {
		return u.baseRes(http.StatusInternalServerError, err.Error(), err)
	}
	return u.baseRes(http.StatusOK, "deleted", nil)
}

func (u *usecase) GetHargaPelaporanByID(ctx context.Context, id uuid.UUID) dto.ResHargaPelaporanSingle {
	obj, err := u.transactionDataRepo.GetHargaPelaporanByID(ctx, id)
	if err != nil {
		return dto.ResHargaPelaporanSingle{BaseRes: u.baseRes(http.StatusNotFound, "not found", err)}
	}
	res := u.serializer.ToHargaPelaporan(*obj)
	return dto.ResHargaPelaporanSingle{BaseRes: u.baseRes(http.StatusOK, "ok", nil), Data: &res}
}

func (u *usecase) GetHargaPelaporanByFilter(ctx context.Context, req *dto.ReqGetHargaPelaporan) dto.ResHargaPelaporanList {
	f := entity.HargaPelaporanFilter{IDPasar: req.IDPasar, IDKomoditas: req.IDKomoditas, From: req.From, To: req.To, PaginationFilter: queryhelper.SerializeFilterPaginationDtoToEntity(req.BaseReqQueryPagination)}
	if req.Status != nil {
		s := constant.PengumpulanDataStatus(*req.Status)
		f.Status = &s
	} else {
		s := constant.PengumpulanDataFinal
		f.Status = &s
	}
	rows, page, err := u.transactionDataRepo.GetHargaPelaporanByFilter(ctx, &f)
	if err != nil {
		return dto.ResHargaPelaporanList{BaseResPagination: dtobase.BaseResPagination{BaseRes: u.baseRes(http.StatusInternalServerError, err.Error(), err)}}
	}
	res := make([]dto.ResHargaPelaporan, 0, len(rows))
	for _, row := range rows {
		res = append(res, u.serializer.ToHargaPelaporan(row))
	}
	return dto.ResHargaPelaporanList{BaseResPagination: dtobase.BaseResPagination{BaseRes: u.baseRes(http.StatusOK, "ok", nil), Page: u.serializer.ToPage(page)}, Data: res}
}

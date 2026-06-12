package masterdatarepositoryimplementation

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/thdoikn/sihp-be/internal/entity"
	entitybase "github.com/thdoikn/sihp-be/internal/entity/base"
	masterdatarepository "github.com/thdoikn/sihp-be/internal/repository/master-data"
	"github.com/thdoikn/sihp-be/pkg/constant"
	"gorm.io/gorm"
)

type masterDataRepository struct {
	db              *gorm.DB
	pasar           entity.Pasar
	komoditas       entity.Komoditas
	tempatUsaha     entity.TempatUsaha
	komoditasDijual entity.KomoditasDijual
}

func NewMasterDataRepository(db *gorm.DB) masterdatarepository.MasterDataRepository {
	return &masterDataRepository{
		db:              db,
		pasar:           entity.Pasar{},
		komoditas:       entity.Komoditas{},
		tempatUsaha:     entity.TempatUsaha{},
		komoditasDijual: entity.KomoditasDijual{},
	}
}

func (r *masterDataRepository) GetOverview(ctx context.Context) (int64, int64, int64, error) {
	var pasarCount, tuCount, komoditasCount int64
	if r.db == nil {
		return pasarCount, tuCount, komoditasCount, errors.New("database connection is not initialized")
	}
	if err := r.db.WithContext(ctx).Table(r.pasar.TableName()).Count(&pasarCount).Error; err != nil {
		return pasarCount, tuCount, komoditasCount, err
	}
	if err := r.db.WithContext(ctx).Table(r.tempatUsaha.TableName()).Count(&tuCount).Error; err != nil {
		return pasarCount, tuCount, komoditasCount, err
	}
	if err := r.db.WithContext(ctx).Table(r.komoditas.TableName()).Count(&komoditasCount).Error; err != nil {
		return pasarCount, tuCount, komoditasCount, err
	}
	return pasarCount, tuCount, komoditasCount, nil
}

func (r *masterDataRepository) CreatePasar(ctx context.Context, input *entity.Pasar) (*entity.Pasar, error) {
	if r.db == nil {
		return nil, errors.New("database connection is not initialized")
	}
	if err := r.db.WithContext(ctx).Create(input).Error; err != nil {
		return nil, err
	}
	return r.GetPasarByID(ctx, input.ID)
}

func (r *masterDataRepository) GetPasarByID(ctx context.Context, id uuid.UUID) (*entity.Pasar, error) {
	if r.db == nil {
		return nil, errors.New("database connection is not initialized")
	}
	var out entity.Pasar
	if err := r.db.WithContext(ctx).Table(r.pasar.TableName()).First(&out, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *masterDataRepository) GetPasarByFilter(ctx context.Context, filter *entity.PasarFilter) ([]entity.Pasar, entitybase.BasePaginationResult, error) {
	if r.db == nil {
		return nil, entitybase.BasePaginationResult{}, errors.New("database connection is not initialized")
	}
	var out []entity.Pasar
	var pagination entitybase.BasePaginationResult
	query := r.db.WithContext(ctx).Table(r.pasar.TableName())
	if filter.Name != nil {
		query = query.Where("sihp.pasar.nama ILIKE ?", "%"+*filter.Name+"%")
	}
	if filter.Status != nil {
		query = query.Where("sihp.pasar.status = ?", *filter.Status)
	}
	query = entitybase.PaginateEntityQuery(query, r.pasar.TableName(), r.pasar.OrderMap(), &filter.PaginationFilter, &pagination)
	if err := query.Find(&out).Error; err != nil {
		return nil, pagination, err
	}
	return out, pagination, nil
}

func (r *masterDataRepository) UpdatePasar(ctx context.Context, id uuid.UUID, updateMap map[string]any) (*entity.Pasar, error) {
	if r.db == nil {
		return nil, errors.New("database connection is not initialized")
	}
	if err := r.db.WithContext(ctx).Table(r.pasar.TableName()).Where("id = ?", id).Updates(updateMap).Error; err != nil {
		return nil, err
	}
	return r.GetPasarByID(ctx, id)
}

func (r *masterDataRepository) CreateKomoditas(ctx context.Context, input *entity.Komoditas) (*entity.Komoditas, error) {
	if r.db == nil {
		return nil, errors.New("database connection is not initialized")
	}
	if err := r.db.WithContext(ctx).Create(input).Error; err != nil {
		return nil, err
	}
	return r.GetKomoditasByID(ctx, input.ID)
}
func (r *masterDataRepository) GetKomoditasByID(ctx context.Context, id uuid.UUID) (*entity.Komoditas, error) {
	var out entity.Komoditas
	if err := r.db.WithContext(ctx).Table(r.komoditas.TableName()).First(&out, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &out, nil
}
func (r *masterDataRepository) GetKomoditasByFilter(ctx context.Context, filter *entity.KomoditasFilter) ([]entity.Komoditas, entitybase.BasePaginationResult, error) {
	if r.db == nil {
		return nil, entitybase.BasePaginationResult{}, errors.New("database connection is not initialized")
	}
	var out []entity.Komoditas
	var pagination entitybase.BasePaginationResult
	query := r.db.WithContext(ctx).Table(r.komoditas.TableName())
	if filter.Name != nil {
		query = query.Where("sihp.komoditas.nama ILIKE ?", "%"+*filter.Name+"%")
	}
	if filter.IDTempatUsaha != nil {
		query = query.Joins("JOIN sihp.komoditas_dijual kd ON kd.id_komoditas = sihp.komoditas.id AND kd.deleted_at IS NULL").Where("kd.id_tempat_usaha = ?", *filter.IDTempatUsaha)
	}
	if filter.IDPasar != nil {
		query = query.Joins("JOIN sihp.komoditas_dijual kd2 ON kd2.id_komoditas = sihp.komoditas.id AND kd2.deleted_at IS NULL").
			Joins("JOIN sihp.tempat_usaha tu ON tu.id = kd2.id_tempat_usaha AND tu.deleted_at IS NULL").
			Where("tu.id_pasar = ?", *filter.IDPasar)
	}
	query = query.Distinct(
		"sihp.komoditas.id",
		"sihp.komoditas.nama",
		"sihp.komoditas.satuan",
		"sihp.komoditas.gambar_url",
		"sihp.komoditas.created_at",
		"sihp.komoditas.updated_at",
		"sihp.komoditas.deleted_at",
	)
	query = entitybase.PaginateEntityQuery(query, (&entity.Komoditas{}).TableName(), (&entity.Komoditas{}).OrderMap(), &filter.PaginationFilter, &pagination)
	if err := query.Find(&out).Error; err != nil {
		return nil, pagination, err
	}
	return out, pagination, nil
}

func (r *masterDataRepository) UpdateKomoditas(ctx context.Context, id uuid.UUID, updateMap map[string]any) (*entity.Komoditas, error) {
	if r.db == nil {
		return nil, errors.New("database connection is not initialized")
	}
	if err := r.db.WithContext(ctx).Table(r.komoditas.TableName()).Where("id = ?", id).Updates(updateMap).Error; err != nil {
		return nil, err
	}
	return r.GetKomoditasByID(ctx, id)
}
func (r *masterDataRepository) DeleteKomoditas(ctx context.Context, id uuid.UUID) error {
	if r.db == nil {
		return errors.New("database connection is not initialized")
	}
	return r.db.WithContext(ctx).Table(r.komoditas.TableName()).Delete(&entity.Komoditas{}, "id = ?", id).Error
}

func (r *masterDataRepository) CreateTempatUsaha(ctx context.Context, input *entity.TempatUsaha) (*entity.TempatUsaha, error) {
	if r.db == nil {
		return nil, errors.New("database connection is not initialized")
	}
	if err := r.db.WithContext(ctx).Create(input).Error; err != nil {
		return nil, err
	}
	return r.GetTempatUsahaByID(ctx, input.ID)
}
func (r *masterDataRepository) GetTempatUsahaByID(ctx context.Context, id uuid.UUID) (*entity.TempatUsaha, error) {
	if r.db == nil {
		return nil, errors.New("database connection is not initialized")
	}
	var out entity.TempatUsaha
	if err := r.db.WithContext(ctx).Table(r.tempatUsaha.TableName()).First(&out, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &out, nil
}
func (r *masterDataRepository) GetTempatUsahaByFilter(ctx context.Context, filter *entity.TempatUsahaFilter) ([]entity.TempatUsaha, entitybase.BasePaginationResult, error) {
	if r.db == nil {
		return nil, entitybase.BasePaginationResult{}, errors.New("database connection is not initialized")
	}
	var out []entity.TempatUsaha
	var pagination entitybase.BasePaginationResult
	query := r.db.WithContext(ctx).Table(r.tempatUsaha.TableName())
	if filter.Name != nil {
		query = query.Where("nama ILIKE ?", "%"+*filter.Name+"%")
	}
	if filter.IDPasar != nil {
		query = query.Where("id_pasar = ?", *filter.IDPasar)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	query = entitybase.PaginateEntityQuery(query, (&entity.TempatUsaha{}).TableName(), (&entity.TempatUsaha{}).OrderMap(), &filter.PaginationFilter, &pagination)
	if err := query.Find(&out).Error; err != nil {
		return nil, pagination, err
	}
	return out, pagination, nil
}
func (r *masterDataRepository) UpdateTempatUsaha(ctx context.Context, id uuid.UUID, updateMap map[string]any) (*entity.TempatUsaha, error) {
	if r.db == nil {
		return nil, errors.New("database connection is not initialized")
	}
	if err := r.db.WithContext(ctx).Table(r.tempatUsaha.TableName()).Where("id = ?", id).Updates(updateMap).Error; err != nil {
		return nil, err
	}
	return r.GetTempatUsahaByID(ctx, id)
}

func (r *masterDataRepository) CreateKomoditasDijual(ctx context.Context, input *entity.KomoditasDijual) (*entity.KomoditasDijual, error) {
	if r.db == nil {
		return nil, errors.New("database connection is not initialized")
	}
	if err := r.db.WithContext(ctx).Create(input).Error; err != nil {
		return nil, err
	}
	return r.GetKomoditasDijualByID(ctx, input.ID)
}
func (r *masterDataRepository) GetKomoditasDijualByID(ctx context.Context, id uuid.UUID) (*entity.KomoditasDijual, error) {
	if r.db == nil {
		return nil, errors.New("database connection is not initialized")
	}
	var out entity.KomoditasDijual
	if err := r.db.WithContext(ctx).Table(r.komoditasDijual.TableName()).First(&out, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &out, nil
}
func (r *masterDataRepository) GetKomoditasDijualByFilter(ctx context.Context, filter *entity.KomoditasDijualFilter) ([]entity.KomoditasDijual, entitybase.BasePaginationResult, error) {
	if r.db == nil {
		return nil, entitybase.BasePaginationResult{}, errors.New("database connection is not initialized")
	}
	var out []entity.KomoditasDijual
	var pagination entitybase.BasePaginationResult
	query := r.db.WithContext(ctx).Table(r.komoditasDijual.TableName())
	if filter.IDTempatUsaha != nil {
		query = query.Where("id_tempat_usaha = ?", *filter.IDTempatUsaha)
	}
	if filter.IDKomoditas != nil {
		query = query.Where("id_komoditas = ?", *filter.IDKomoditas)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	query = entitybase.PaginateEntityQuery(query, (&entity.KomoditasDijual{}).TableName(), (&entity.KomoditasDijual{}).OrderMap(), &filter.PaginationFilter, &pagination)
	if err := query.Find(&out).Error; err != nil {
		return nil, pagination, err
	}
	return out, pagination, nil
}
func (r *masterDataRepository) UpdateKomoditasDijual(ctx context.Context, id uuid.UUID, updateMap map[string]any) (*entity.KomoditasDijual, error) {
	if r.db == nil {
		return nil, errors.New("database connection is not initialized")
	}
	if err := r.db.WithContext(ctx).Table(r.komoditasDijual.TableName()).Where("id = ?", id).Updates(updateMap).Error; err != nil {
		return nil, err
	}
	return r.GetKomoditasDijualByID(ctx, id)
}

func (r *masterDataRepository) GetPublicKomoditasList(ctx context.Context, filter *entity.PublicKomoditasListFilter) ([]entity.PublicKomoditasListItem, entitybase.BasePaginationResult, error) {
	if r.db == nil {
		return nil, entitybase.BasePaginationResult{}, errors.New("database connection is not initialized")
	}
	if filter == nil {
		return nil, entitybase.BasePaginationResult{}, errors.New("filter is required")
	}

	var baseSQL string
	var args []any

	switch {
	case filter.IDTempatUsaha != nil:
		baseSQL = `
SELECT k.id, k.nama, k.satuan, k.gambar_url,
  lat.harga_terbaru, agg.harga_terkecil, agg.harga_terbesar, agg.harga_avg, lat.tanggal_terbaru
FROM sihp.komoditas k
INNER JOIN sihp.komoditas_dijual kd ON kd.id_komoditas = k.id AND kd.deleted_at IS NULL AND kd.id_tempat_usaha = ?
LEFT JOIN (
  SELECT hr.id_komoditas,
    AVG(hr.harga)::float8 AS harga_terbaru,
    MAX(pd.tanggal) AS tanggal_terbaru
  FROM sihp.harga_rutin hr
  INNER JOIN sihp.pengumpulan_data pd ON pd.id = hr.id_pengumpulan_data AND pd.deleted_at IS NULL AND pd.status = ?
  WHERE hr.id_tempat_usaha = ?
    AND pd.tanggal = (
      SELECT MAX(pd2.tanggal)
      FROM sihp.pengumpulan_data pd2
      INNER JOIN sihp.harga_rutin hr2 ON hr2.id_pengumpulan_data = pd2.id AND hr2.id_tempat_usaha = ?
      WHERE hr2.id_komoditas = hr.id_komoditas AND pd2.deleted_at IS NULL AND pd2.status = ?
    )
  GROUP BY hr.id_komoditas
) lat ON lat.id_komoditas = k.id
LEFT JOIN (
  SELECT hr.id_komoditas,
    MIN(hr.harga)::float8 AS harga_terkecil,
    MAX(hr.harga)::float8 AS harga_terbesar,
    AVG(hr.harga)::float8 AS harga_avg
  FROM sihp.harga_rutin hr
  INNER JOIN sihp.pengumpulan_data pd ON pd.id = hr.id_pengumpulan_data AND pd.deleted_at IS NULL AND pd.status = ?
  WHERE hr.id_tempat_usaha = ?
  GROUP BY hr.id_komoditas
) agg ON agg.id_komoditas = k.id
WHERE k.deleted_at IS NULL`
		args = []any{
			*filter.IDTempatUsaha,
			constant.PengumpulanDataFinal, *filter.IDTempatUsaha,
			*filter.IDTempatUsaha, constant.PengumpulanDataFinal,
			constant.PengumpulanDataFinal, *filter.IDTempatUsaha,
		}
	case filter.IDPasar != nil:
		baseSQL = `
SELECT k.id, k.nama, k.satuan, k.gambar_url,
  lat.harga_terbaru, agg.harga_terkecil, agg.harga_terbesar, agg.harga_avg, lat.tanggal_terbaru
FROM sihp.komoditas k
INNER JOIN (
  SELECT DISTINCT ON (hp.id_komoditas)
    hp.id_komoditas, hp.harga::float8 AS harga_terbaru, hp.tanggal AS tanggal_terbaru
  FROM sihp.harga_pelaporan hp
  INNER JOIN sihp.pengumpulan_data pd ON pd.id = hp.id_pengumpulan_data AND pd.deleted_at IS NULL AND pd.status = ? AND pd.id_pasar = ?
  ORDER BY hp.id_komoditas, hp.tanggal DESC, hp.created_at DESC
) lat ON lat.id_komoditas = k.id
INNER JOIN (
  SELECT hp.id_komoditas,
    MIN(hp.harga)::float8 AS harga_terkecil,
    MAX(hp.harga)::float8 AS harga_terbesar,
    AVG(hp.harga)::float8 AS harga_avg
  FROM sihp.harga_pelaporan hp
  INNER JOIN sihp.pengumpulan_data pd ON pd.id = hp.id_pengumpulan_data AND pd.deleted_at IS NULL AND pd.status = ? AND pd.id_pasar = ?
  GROUP BY hp.id_komoditas
) agg ON agg.id_komoditas = k.id
WHERE k.deleted_at IS NULL`
		args = []any{constant.PengumpulanDataFinal, *filter.IDPasar, constant.PengumpulanDataFinal, *filter.IDPasar}
	default:
		baseSQL = `
SELECT k.id, k.nama, k.satuan, k.gambar_url,
  lat.harga_terbaru, agg.harga_terkecil, agg.harga_terbesar, agg.harga_avg, lat.tanggal_terbaru
FROM sihp.komoditas k
LEFT JOIN (
  SELECT DISTINCT ON (hp.id_komoditas)
    hp.id_komoditas, hp.harga::float8 AS harga_terbaru, hp.tanggal AS tanggal_terbaru
  FROM sihp.harga_pelaporan hp
  INNER JOIN sihp.pengumpulan_data pd ON pd.id = hp.id_pengumpulan_data AND pd.deleted_at IS NULL AND pd.status = ?
  ORDER BY hp.id_komoditas, hp.tanggal DESC, hp.created_at DESC
) lat ON lat.id_komoditas = k.id
LEFT JOIN (
  SELECT hp.id_komoditas,
    MIN(hp.harga)::float8 AS harga_terkecil,
    MAX(hp.harga)::float8 AS harga_terbesar,
    AVG(hp.harga)::float8 AS harga_avg
  FROM sihp.harga_pelaporan hp
  INNER JOIN sihp.pengumpulan_data pd ON pd.id = hp.id_pengumpulan_data AND pd.deleted_at IS NULL AND pd.status = ?
  GROUP BY hp.id_komoditas
) agg ON agg.id_komoditas = k.id
WHERE k.deleted_at IS NULL`
		args = []any{constant.PengumpulanDataFinal, constant.PengumpulanDataFinal}
	}

	if filter.Nama != nil && *filter.Nama != "" {
		baseSQL += " AND k.nama ILIKE ?"
		args = append(args, "%"+*filter.Nama+"%")
	}

	countSQL := "SELECT COUNT(*) FROM (" + baseSQL + ") counted"
	var total int64
	if err := r.db.WithContext(ctx).Raw(countSQL, args...).Scan(&total).Error; err != nil {
		return nil, entitybase.BasePaginationResult{}, err
	}

	limit := 10
	if filter.PaginationFilter.Limit != nil && *filter.PaginationFilter.Limit > 0 {
		limit = *filter.PaginationFilter.Limit
		if limit > 50 {
			limit = 50
		}
	}
	offset := 0
	if filter.PaginationFilter.Offset != nil && *filter.PaginationFilter.Offset > 0 {
		offset = *filter.PaginationFilter.Offset
	}

	listSQL := baseSQL + " ORDER BY k.nama ASC LIMIT ? OFFSET ?"
	listArgs := append(append([]any{}, args...), limit, offset)

	rows := []entity.PublicKomoditasListItem{}
	if err := r.db.WithContext(ctx).Raw(listSQL, listArgs...).Scan(&rows).Error; err != nil {
		return nil, entitybase.BasePaginationResult{}, err
	}

	return rows, entitybase.BasePaginationResult{
		Offset:  offset,
		Limit:   limit,
		Count:   int(total),
		OrderBy: "nama asc",
	}, nil
}

func (r *masterDataRepository) publicHargaPelaporanQuery(
	ctx context.Context,
	komoditasID uuid.UUID,
	idPasar *uuid.UUID,
) *gorm.DB {
	query := r.db.WithContext(ctx).Table("sihp.harga_pelaporan hp").
		Joins("JOIN sihp.pengumpulan_data pd ON pd.id = hp.id_pengumpulan_data AND pd.deleted_at IS NULL").
		Where("hp.id_komoditas = ? AND pd.status = ?", komoditasID, constant.PengumpulanDataFinal)
	if idPasar != nil {
		query = query.Where("pd.id_pasar = ?", *idPasar)
	}
	return query
}

func (r *masterDataRepository) GetPublicKomoditasStats(ctx context.Context, komoditasID uuid.UUID, days int, idPasar *uuid.UUID) (*time.Time, *float64, *float64, *float64, *float64, error) {
	var latest struct {
		Tanggal time.Time
		Harga   float64
	}
	if err := r.publicHargaPelaporanQuery(ctx, komoditasID, idPasar).
		Select("hp.tanggal, hp.harga::float8 as harga").
		Order("hp.tanggal desc").
		Limit(1).
		Scan(&latest).Error; err != nil {
		return nil, nil, nil, nil, nil, err
	}
	var latestTanggal *time.Time
	var latestHarga *float64
	if !latest.Tanggal.IsZero() {
		latestTanggal = &latest.Tanggal
		latestHarga = &latest.Harga
	}

	from := time.Now().AddDate(0, 0, -days)
	stats := struct{ Avg, Min, Max *float64 }{}
	if err := r.publicHargaPelaporanQuery(ctx, komoditasID, idPasar).
		Where("hp.tanggal >= ?", from.Format("2006-01-02")).
		Select("AVG(hp.harga)::float8 as avg, MIN(hp.harga)::float8 as min, MAX(hp.harga)::float8 as max").
		Scan(&stats).Error; err != nil {
		return nil, nil, nil, nil, nil, err
	}
	return latestTanggal, latestHarga, stats.Avg, stats.Min, stats.Max, nil
}

func (r *masterDataRepository) GetPublicKomoditasTrend(ctx context.Context, komoditasID uuid.UUID, days int, idPasar *uuid.UUID) ([]map[string]any, error) {
	from := time.Now().AddDate(0, 0, -days)
	query := r.db.WithContext(ctx).Table("sihp.harga_pelaporan hp").
		Joins("JOIN sihp.pengumpulan_data pd ON pd.id = hp.id_pengumpulan_data AND pd.deleted_at IS NULL").
		Where("hp.id_komoditas = ? AND pd.status = ? AND hp.tanggal >= ?", komoditasID, constant.PengumpulanDataFinal, from.Format("2006-01-02"))
	if idPasar != nil {
		query = query.Where("pd.id_pasar = ?", *idPasar)
	}
	out := []map[string]any{}
	err := query.Select("hp.tanggal as tanggal, AVG(hp.harga)::float8 as harga_rata_rata").Group("hp.tanggal").Order("hp.tanggal asc").Scan(&out).Error
	return out, err
}

func (r *masterDataRepository) GetPublicPasarList(ctx context.Context) ([]entity.PublicPasarListItem, error) {
	if r.db == nil {
		return nil, errors.New("database connection is not initialized")
	}

	const sql = `
SELECT p.id, p.nama, p.alamat, p.status, p.longitude, p.latitude,
  COALESCE(tu_count.cnt, 0) AS total_tempat_usaha,
  COALESCE(kom_count.cnt, 0) AS total_komoditas
FROM sihp.pasar p
LEFT JOIN (
  SELECT id_pasar, COUNT(*) AS cnt
  FROM sihp.tempat_usaha
  WHERE deleted_at IS NULL AND status = ?
  GROUP BY id_pasar
) tu_count ON tu_count.id_pasar = p.id
LEFT JOIN (
  SELECT pd.id_pasar, COUNT(DISTINCT hp.id_komoditas) AS cnt
  FROM sihp.harga_pelaporan hp
  INNER JOIN sihp.pengumpulan_data pd ON pd.id = hp.id_pengumpulan_data AND pd.deleted_at IS NULL AND pd.status = ?
  GROUP BY pd.id_pasar
) kom_count ON kom_count.id_pasar = p.id
WHERE p.deleted_at IS NULL AND p.status = ?
ORDER BY p.nama ASC`

	rows := []entity.PublicPasarListItem{}
	err := r.db.WithContext(ctx).Raw(sql, constant.StatusActive, constant.PengumpulanDataFinal, constant.StatusActive).Scan(&rows).Error
	return rows, err
}

func (r *masterDataRepository) GetPublicTempatUsahaList(ctx context.Context, filter *entity.PublicTempatUsahaListFilter) ([]entity.PublicTempatUsahaListItem, entitybase.BasePaginationResult, error) {
	if r.db == nil {
		return nil, entitybase.BasePaginationResult{}, errors.New("database connection is not initialized")
	}
	if filter == nil {
		return nil, entitybase.BasePaginationResult{}, errors.New("filter is required")
	}

	baseSQL := `
SELECT tu.id, tu.nama, tu.id_pasar AS pasar_id, p.nama AS pasar_nama
FROM sihp.tempat_usaha tu
INNER JOIN sihp.pasar p ON p.id = tu.id_pasar AND p.deleted_at IS NULL
WHERE tu.deleted_at IS NULL AND tu.status = ?`
	args := []any{constant.StatusActive}

	if filter.Nama != nil && *filter.Nama != "" {
		baseSQL += " AND tu.nama ILIKE ?"
		args = append(args, "%"+*filter.Nama+"%")
	}
	if filter.IDPasar != nil {
		baseSQL += " AND tu.id_pasar = ?"
		args = append(args, *filter.IDPasar)
	}

	countSQL := "SELECT COUNT(*) FROM (" + baseSQL + ") counted"
	var total int64
	if err := r.db.WithContext(ctx).Raw(countSQL, args...).Scan(&total).Error; err != nil {
		return nil, entitybase.BasePaginationResult{}, err
	}

	limit := 10
	if filter.PaginationFilter.Limit != nil && *filter.PaginationFilter.Limit > 0 {
		limit = *filter.PaginationFilter.Limit
		if limit > 50 {
			limit = 50
		}
	}
	offset := 0
	if filter.PaginationFilter.Offset != nil && *filter.PaginationFilter.Offset > 0 {
		offset = *filter.PaginationFilter.Offset
	}

	listSQL := baseSQL + " ORDER BY tu.nama ASC LIMIT ? OFFSET ?"
	listArgs := append(append([]any{}, args...), limit, offset)

	rows := []entity.PublicTempatUsahaListItem{}
	if err := r.db.WithContext(ctx).Raw(listSQL, listArgs...).Scan(&rows).Error; err != nil {
		return nil, entitybase.BasePaginationResult{}, err
	}

	return rows, entitybase.BasePaginationResult{
		Offset:  offset,
		Limit:   limit,
		Count:   int(total),
		OrderBy: "nama asc",
	}, nil
}

func (r *masterDataRepository) GetPublicPasarDetail(ctx context.Context, id uuid.UUID, filter *entity.TempatUsahaFilter) (*entity.Pasar, []entity.TempatUsaha, entitybase.BasePaginationResult, error) {
	pasar, err := r.GetPasarByID(ctx, id)
	if err != nil {
		return nil, nil, entitybase.BasePaginationResult{}, err
	}
	local := *filter
	local.IDPasar = &id
	tus, page, err := r.GetTempatUsahaByFilter(ctx, &local)
	return pasar, tus, page, err
}

func (r *masterDataRepository) GetPublicTempatUsahaDetail(ctx context.Context, id uuid.UUID, filter entity.KomoditasFilter) (*entity.TempatUsaha, *entity.Pasar, []map[string]any, entitybase.BasePaginationResult, error) {
	tu, err := r.GetTempatUsahaByID(ctx, id)
	if err != nil {
		return nil, nil, nil, entitybase.BasePaginationResult{}, err
	}
	pasar, err := r.GetPasarByID(ctx, tu.IDPasar)
	if err != nil {
		return nil, nil, nil, entitybase.BasePaginationResult{}, err
	}

	limit := 50
	if filter.PaginationFilter.Limit != nil && *filter.PaginationFilter.Limit > 0 {
		limit = *filter.PaginationFilter.Limit
		if limit > 1000 {
			limit = 1000
		}
	}
	offset := 0
	if filter.PaginationFilter.Offset != nil && *filter.PaginationFilter.Offset > 0 {
		offset = *filter.PaginationFilter.Offset
	}

	baseSQL := `
SELECT DISTINCT ON (k.id) k.id, k.nama, k.satuan, k.gambar_url
FROM sihp.komoditas k
INNER JOIN sihp.komoditas_dijual kd ON kd.id_komoditas = k.id AND kd.deleted_at IS NULL
WHERE kd.id_tempat_usaha = ? AND k.deleted_at IS NULL`
	args := []any{id}
	if filter.Name != nil && *filter.Name != "" {
		baseSQL += " AND k.nama ILIKE ?"
		args = append(args, "%"+*filter.Name+"%")
	}
	baseSQL += " ORDER BY k.id, k.nama ASC"

	countSQL := `SELECT COUNT(*) FROM (` + baseSQL + `) counted`
	var total int64
	if err := r.db.WithContext(ctx).Raw(countSQL, args...).Scan(&total).Error; err != nil {
		return nil, nil, nil, entitybase.BasePaginationResult{}, err
	}

	listSQL := `SELECT * FROM (` + baseSQL + `) items ORDER BY nama ASC LIMIT ? OFFSET ?`
	listArgs := append(append([]any{}, args...), limit, offset)

	type row struct {
		ID        uuid.UUID
		Nama      string
		Satuan    *string
		GambarURL *string
	}
	rows := []row{}
	if err := r.db.WithContext(ctx).Raw(listSQL, listArgs...).Scan(&rows).Error; err != nil {
		return nil, nil, nil, entitybase.BasePaginationResult{}, err
	}

	pagination := entitybase.BasePaginationResult{
		Offset:  offset,
		Limit:   limit,
		Count:   int(total),
		OrderBy: "nama asc",
	}

	result := make([]map[string]any, 0, len(rows))
	for _, item := range rows {
		latest := map[string]any{"tanggal": nil, "harga_rata_rata": nil}
		var temp struct {
			Tanggal time.Time
			Harga   float64
		}
		err := r.db.WithContext(ctx).Raw(`
SELECT pd.tanggal, AVG(hr.harga)::float8 AS harga
FROM sihp.harga_rutin hr
INNER JOIN sihp.pengumpulan_data pd ON pd.id = hr.id_pengumpulan_data AND pd.deleted_at IS NULL AND pd.status = ?
WHERE hr.id_komoditas = ? AND hr.id_tempat_usaha = ?
GROUP BY pd.tanggal
ORDER BY pd.tanggal DESC
LIMIT 1`, constant.PengumpulanDataFinal, item.ID, id).Scan(&temp).Error
		if err != nil {
			return nil, nil, nil, pagination, err
		}
		if !temp.Tanggal.IsZero() {
			latest["tanggal"] = temp.Tanggal
			latest["harga_rata_rata"] = temp.Harga
		}
		result = append(result, map[string]any{"id": item.ID, "nama": item.Nama, "satuan": item.Satuan, "gambar_url": item.GambarURL, "latest": latest})
	}
	return tu, pasar, result, pagination, nil
}

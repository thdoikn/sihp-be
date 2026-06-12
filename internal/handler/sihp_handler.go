package handler

import (
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	sihpusecase "github.com/thdoikn/sihp-be/internal/usecase/sihp"
	"github.com/thdoikn/sihp-be/pkg/dto"
)

type sihpHandler struct{ usecase sihpusecase.SIHPUsecase }

func NewSIHPHandler(usecase sihpusecase.SIHPUsecase) *sihpHandler {
	return &sihpHandler{usecase: usecase}
}

func parseUUIDParam(c *fiber.Ctx) uuid.UUID { return uuid.MustParse(c.Params("id")) }

// PublicOverview godoc
// @Summary Public overview
// @Description Get SIHP public summary counts
// @Tags Public
// @Produce json
// @Success 200 {object} dto.ResPublicOverviewEnvelope
// @Router /public/overview [get]
func (h *sihpHandler) PublicOverview(c *fiber.Ctx) error {
	res := h.usecase.GetPublicOverview(c.Context())
	return c.Status(res.Code).JSON(res)
}

// PublicKomoditas godoc
// @Summary Public komoditas list
// @Description Get public komoditas list
// @Tags Public
// @Produce json
// @Param nama query string false "Komoditas name"
// @Param id_tempat_usaha query string false "Tempat usaha ID"
// @Param id_pasar query string false "Pasar ID"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} dto.ResPublicKomoditasList
// @Router /public/komoditas [get]
func (h *sihpHandler) PublicKomoditas(c *fiber.Ctx) error {
	var req dto.ReqPublicGetKomoditas
	_ = c.QueryParser(&req)
	res := h.usecase.GetPublicKomoditas(c.Context(), &req)
	return c.Status(res.Code).JSON(res)
}

// PublicKomoditasDetail godoc
// @Summary Public komoditas detail
// @Description Get public komoditas detail and stats
// @Tags Public
// @Produce json
// @Param id path string true "Komoditas ID"
// @Param days query int false "Stats range in days"
// @Param id_pasar query string false "Pasar ID"
// @Success 200 {object} dto.ResPublicKomoditasDetailEnvelope
// @Router /public/komoditas/{id} [get]
func (h *sihpHandler) PublicKomoditasDetail(c *fiber.Ctx) error {
	var req dto.ReqPublicKomoditasDetail
	_ = c.QueryParser(&req)
	res := h.usecase.GetPublicKomoditasDetail(c.Context(), parseUUIDParam(c), &req)
	return c.Status(res.Code).JSON(res)
}

// PublicKomoditasTrend godoc
// @Summary Public komoditas trend
// @Description Get public komoditas trend points
// @Tags Public
// @Produce json
// @Param id path string true "Komoditas ID"
// @Param days query int false "Trend range in days"
// @Param id_pasar query string false "Pasar ID"
// @Success 200 {object} dto.ResPublicTrendEnvelope
// @Router /public/komoditas/{id}/trend [get]
func (h *sihpHandler) PublicKomoditasTrend(c *fiber.Ctx) error {
	var req dto.ReqPublicKomoditasDetail
	_ = c.QueryParser(&req)
	res := h.usecase.GetPublicKomoditasTrend(c.Context(), parseUUIDParam(c), &req)
	return c.Status(res.Code).JSON(res)
}

// PublicPasar godoc
// @Summary Public pasar list
// @Description Get public active pasar list with stats
// @Tags Public
// @Produce json
// @Success 200 {object} dto.ResPublicPasarList
// @Router /public/pasar [get]
func (h *sihpHandler) PublicPasar(c *fiber.Ctx) error {
	res := h.usecase.GetPublicPasar(c.Context())
	return c.Status(res.Code).JSON(res)
}

// PublicPasarDetail godoc
// @Summary Public pasar detail
// @Description Get public pasar detail with tempat usaha list
// @Tags Public
// @Produce json
// @Param id path string true "Pasar ID"
// @Param name query string false "Tempat usaha name"
// @Param status query int false "Tempat usaha status"
// @Param show-count query bool false "Show count"
// @Param offset query int false "Offset"
// @Param limit query int false "Limit"
// @Param order-by query string false "Order by"
// @Success 200 {object} dto.ResPublicPasarDetailEnvelope
// @Router /public/pasar/{id} [get]
func (h *sihpHandler) PublicPasarDetail(c *fiber.Ctx) error {
	var req dto.ReqGetTempatUsaha
	_ = c.QueryParser(&req)
	res := h.usecase.GetPublicPasarDetail(c.Context(), parseUUIDParam(c), &req)
	return c.Status(res.Code).JSON(res)
}

// PublicTempatUsaha godoc
// @Summary Public tempat usaha list
// @Description Get public active tempat usaha list
// @Tags Public
// @Produce json
// @Param nama query string false "Tempat usaha name"
// @Param id_pasar query string false "Pasar ID"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} dto.ResPublicTempatUsahaList
// @Router /public/tempat-usaha [get]
func (h *sihpHandler) PublicTempatUsaha(c *fiber.Ctx) error {
	var req dto.ReqPublicGetTempatUsaha
	_ = c.QueryParser(&req)
	res := h.usecase.GetPublicTempatUsaha(c.Context(), &req)
	return c.Status(res.Code).JSON(res)
}

// PublicTempatUsahaDetail godoc
// @Summary Public tempat usaha detail
// @Description Get public tempat usaha detail with komoditas list
// @Tags Public
// @Produce json
// @Param id path string true "Tempat usaha ID"
// @Param name query string false "Komoditas name"
// @Param show-count query bool false "Show count"
// @Param offset query int false "Offset"
// @Param limit query int false "Limit"
// @Param order-by query string false "Order by"
// @Success 200 {object} dto.ResPublicTempatUsahaDetailEnvelope
// @Router /public/tempat-usaha/{id} [get]
func (h *sihpHandler) PublicTempatUsahaDetail(c *fiber.Ctx) error {
	var req dto.ReqGetKomoditas
	_ = c.QueryParser(&req)
	res := h.usecase.GetPublicTempatUsahaDetail(c.Context(), parseUUIDParam(c), &req)
	return c.Status(res.Code).JSON(res)
}

// CreatePasar godoc
// @Summary Create pasar
// @Description Create new pasar
// @Tags Pasar
// @Security Authorization
// @Accept json
// @Produce json
// @Param payload body dto.ReqCreatePasar true "Create pasar payload"
// @Success 201 {object} dto.ResPasarSingle
// @Router /admin/pasar [post]
func (h *sihpHandler) CreatePasar(c *fiber.Ctx) error {
	var req dto.ReqCreatePasar
	_ = c.BodyParser(&req)
	res := h.usecase.CreatePasar(c.Context(), &req)
	return c.Status(res.Code).JSON(res)
}

// GetPasarByID godoc
// @Summary Get pasar by ID
// @Description Get pasar detail by ID
// @Tags Pasar
// @Security Authorization
// @Produce json
// @Param id path string true "Pasar ID"
// @Success 200 {object} dto.ResPasarSingle
// @Router /admin/pasar/{id} [get]
func (h *sihpHandler) GetPasarByID(c *fiber.Ctx) error {
	res := h.usecase.GetPasarByID(c.Context(), parseUUIDParam(c))
	return c.Status(res.Code).JSON(res)
}

// GetPasarByFilter godoc
// @Summary Get pasar list
// @Description Get pasar list with filters
// @Tags Pasar
// @Security Authorization
// @Produce json
// @Param name query string false "Pasar name"
// @Param status query int false "Pasar status"
// @Param show-count query bool false "Show count"
// @Param offset query int false "Offset"
// @Param limit query int false "Limit"
// @Param order-by query string false "Order by"
// @Success 200 {object} dto.ResPasarList
// @Router /admin/pasar [get]
func (h *sihpHandler) GetPasarByFilter(c *fiber.Ctx) error {
	var req dto.ReqGetPasar
	_ = c.QueryParser(&req)
	res := h.usecase.GetPasarByFilter(c.Context(), &req)
	return c.Status(res.Code).JSON(res)
}

// UpdatePasar godoc
// @Summary Update pasar
// @Description Update pasar by ID
// @Tags Pasar
// @Security Authorization
// @Accept json
// @Produce json
// @Param id path string true "Pasar ID"
// @Param payload body dto.ReqUpdatePasar true "Update pasar payload"
// @Success 200 {object} dto.ResPasarSingle
// @Router /admin/pasar/{id} [put]
func (h *sihpHandler) UpdatePasar(c *fiber.Ctx) error {
	var req dto.ReqUpdatePasar
	_ = c.BodyParser(&req)
	res := h.usecase.UpdatePasar(c.Context(), parseUUIDParam(c), &req)
	return c.Status(res.Code).JSON(res)
}

// DeletePasar godoc
// @Summary Delete pasar
// @Description Set pasar status to inactive
// @Tags Pasar
// @Security Authorization
// @Produce json
// @Param id path string true "Pasar ID"
// @Success 200 {object} dto.ResPasarSingle
// @Router /admin/pasar/{id} [delete]
func (h *sihpHandler) DeletePasar(c *fiber.Ctx) error {
	res := h.usecase.DeletePasar(c.Context(), parseUUIDParam(c))
	return c.Status(res.Code).JSON(res)
}

// CreateKomoditas godoc
// @Summary Create komoditas
// @Description Create new komoditas
// @Tags Komoditas
// @Security Authorization
// @Accept json
// @Produce json
// @Param payload body dto.ReqCreateKomoditas true "Create komoditas payload"
// @Success 201 {object} dto.ResKomoditasSingle
// @Router /admin/komoditas [post]
func (h *sihpHandler) CreateKomoditas(c *fiber.Ctx) error {
	var req dto.ReqCreateKomoditas
	_ = c.BodyParser(&req)
	res := h.usecase.CreateKomoditas(c.Context(), &req)
	return c.Status(res.Code).JSON(res)
}

// GetKomoditasByID godoc
// @Summary Get komoditas by ID
// @Description Get komoditas detail by ID
// @Tags Komoditas
// @Security Authorization
// @Produce json
// @Param id path string true "Komoditas ID"
// @Success 200 {object} dto.ResKomoditasSingle
// @Router /admin/komoditas/{id} [get]
func (h *sihpHandler) GetKomoditasByID(c *fiber.Ctx) error {
	res := h.usecase.GetKomoditasByID(c.Context(), parseUUIDParam(c))
	return c.Status(res.Code).JSON(res)
}

// GetKomoditasByFilter godoc
// @Summary Get komoditas list
// @Description Get komoditas list with filters
// @Tags Komoditas
// @Security Authorization
// @Produce json
// @Param name query string false "Komoditas name"
// @Param id_tempat_usaha query string false "Tempat usaha ID"
// @Param id_pasar query string false "Pasar ID"
// @Param show-count query bool false "Show count"
// @Param offset query int false "Offset"
// @Param limit query int false "Limit"
// @Param order-by query string false "Order by"
// @Success 200 {object} dto.ResKomoditasList
// @Router /admin/komoditas [get]
func (h *sihpHandler) GetKomoditasByFilter(c *fiber.Ctx) error {
	var req dto.ReqGetKomoditas
	_ = c.QueryParser(&req)
	res := h.usecase.GetKomoditasByFilter(c.Context(), &req)
	return c.Status(res.Code).JSON(res)
}

// UpdateKomoditas godoc
// @Summary Update komoditas
// @Description Update komoditas by ID
// @Tags Komoditas
// @Security Authorization
// @Accept json
// @Produce json
// @Param id path string true "Komoditas ID"
// @Param payload body dto.ReqUpdateKomoditas true "Update komoditas payload"
// @Success 200 {object} dto.ResKomoditasSingle
// @Router /admin/komoditas/{id} [put]
func (h *sihpHandler) UpdateKomoditas(c *fiber.Ctx) error {
	var req dto.ReqUpdateKomoditas
	_ = c.BodyParser(&req)
	res := h.usecase.UpdateKomoditas(c.Context(), parseUUIDParam(c), &req)
	return c.Status(res.Code).JSON(res)
}

// UploadKomoditasGambar godoc
// @Summary Upload komoditas image
// @Description Upload komoditas image to MinIO and save public URL
// @Tags Komoditas
// @Security Authorization
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Komoditas ID"
// @Param gambar formData file true "Komoditas image (jpeg/png/webp, max 2MB)"
// @Success 200 {object} dto.ResKomoditasSingle
// @Router /admin/komoditas/{id}/gambar [post]
func (h *sihpHandler) UploadKomoditasGambar(c *fiber.Ctx) error {
	file, err := c.FormFile("gambar")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"code":    fiber.StatusBadRequest,
			"message": "field gambar is required",
		})
	}

	f, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"code":    fiber.StatusBadRequest,
			"message": "failed to read uploaded file",
		})
	}
	defer f.Close()

	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		switch strings.ToLower(filepath.Ext(file.Filename)) {
		case ".jpg", ".jpeg":
			contentType = "image/jpeg"
		case ".png":
			contentType = "image/png"
		case ".webp":
			contentType = "image/webp"
		}
	}
	res := h.usecase.UploadKomoditasGambar(c.Context(), parseUUIDParam(c), f, file.Size, contentType)
	return c.Status(res.Code).JSON(res)
}

// CreateTempatUsaha godoc
// @Summary Create tempat usaha
// @Description Create new tempat usaha
// @Tags Tempat Usaha
// @Security Authorization
// @Accept json
// @Produce json
// @Param payload body dto.ReqCreateTempatUsaha true "Create tempat usaha payload"
// @Success 201 {object} dto.ResTempatUsahaSingle
// @Router /admin/tempat-usaha [post]
func (h *sihpHandler) CreateTempatUsaha(c *fiber.Ctx) error {
	var req dto.ReqCreateTempatUsaha
	_ = c.BodyParser(&req)
	res := h.usecase.CreateTempatUsaha(c.Context(), &req)
	return c.Status(res.Code).JSON(res)
}

// GetTempatUsahaByID godoc
// @Summary Get tempat usaha by ID
// @Description Get tempat usaha detail by ID
// @Tags Tempat Usaha
// @Security Authorization
// @Produce json
// @Param id path string true "Tempat usaha ID"
// @Success 200 {object} dto.ResTempatUsahaSingle
// @Router /admin/tempat-usaha/{id} [get]
func (h *sihpHandler) GetTempatUsahaByID(c *fiber.Ctx) error {
	res := h.usecase.GetTempatUsahaByID(c.Context(), parseUUIDParam(c))
	return c.Status(res.Code).JSON(res)
}

// GetTempatUsahaByFilter godoc
// @Summary Get tempat usaha list
// @Description Get tempat usaha list with filters
// @Tags Tempat Usaha
// @Security Authorization
// @Produce json
// @Param name query string false "Tempat usaha name"
// @Param id_pasar query string false "Pasar ID"
// @Param status query int false "Status"
// @Param show-count query bool false "Show count"
// @Param offset query int false "Offset"
// @Param limit query int false "Limit"
// @Param order-by query string false "Order by"
// @Success 200 {object} dto.ResTempatUsahaList
// @Router /admin/tempat-usaha [get]
func (h *sihpHandler) GetTempatUsahaByFilter(c *fiber.Ctx) error {
	var req dto.ReqGetTempatUsaha
	_ = c.QueryParser(&req)
	res := h.usecase.GetTempatUsahaByFilter(c.Context(), &req)
	return c.Status(res.Code).JSON(res)
}

// UpdateTempatUsaha godoc
// @Summary Update tempat usaha
// @Description Update tempat usaha by ID
// @Tags Tempat Usaha
// @Security Authorization
// @Accept json
// @Produce json
// @Param id path string true "Tempat usaha ID"
// @Param payload body dto.ReqUpdateTempatUsaha true "Update tempat usaha payload"
// @Success 200 {object} dto.ResTempatUsahaSingle
// @Router /admin/tempat-usaha/{id} [put]
func (h *sihpHandler) UpdateTempatUsaha(c *fiber.Ctx) error {
	var req dto.ReqUpdateTempatUsaha
	_ = c.BodyParser(&req)
	res := h.usecase.UpdateTempatUsaha(c.Context(), parseUUIDParam(c), &req)
	return c.Status(res.Code).JSON(res)
}

// DeleteTempatUsaha godoc
// @Summary Delete tempat usaha
// @Description Set tempat usaha status to inactive
// @Tags Tempat Usaha
// @Security Authorization
// @Produce json
// @Param id path string true "Tempat usaha ID"
// @Success 200 {object} dto.ResTempatUsahaSingle
// @Router /admin/tempat-usaha/{id} [delete]
func (h *sihpHandler) DeleteTempatUsaha(c *fiber.Ctx) error {
	res := h.usecase.DeleteTempatUsaha(c.Context(), parseUUIDParam(c))
	return c.Status(res.Code).JSON(res)
}

// CreateKomoditasDijual godoc
// @Summary Create komoditas dijual
// @Description Create komoditas dijual row
// @Tags Komoditas Dijual
// @Security Authorization
// @Accept json
// @Produce json
// @Param payload body dto.ReqCreateKomoditasDijual true "Create komoditas dijual payload"
// @Success 201 {object} dto.ResKomoditasDijualSingle
// @Router /admin/komoditas-dijual [post]
func (h *sihpHandler) CreateKomoditasDijual(c *fiber.Ctx) error {
	var req dto.ReqCreateKomoditasDijual
	_ = c.BodyParser(&req)
	res := h.usecase.CreateKomoditasDijual(c.Context(), &req)
	return c.Status(res.Code).JSON(res)
}

// GetKomoditasDijualByID godoc
// @Summary Get komoditas dijual by ID
// @Description Get komoditas dijual detail by ID
// @Tags Komoditas Dijual
// @Security Authorization
// @Produce json
// @Param id path string true "Komoditas dijual ID"
// @Success 200 {object} dto.ResKomoditasDijualSingle
// @Router /admin/komoditas-dijual/{id} [get]
func (h *sihpHandler) GetKomoditasDijualByID(c *fiber.Ctx) error {
	res := h.usecase.GetKomoditasDijualByID(c.Context(), parseUUIDParam(c))
	return c.Status(res.Code).JSON(res)
}

// GetKomoditasDijualByFilter godoc
// @Summary Get komoditas dijual list
// @Description Get komoditas dijual list with filters
// @Tags Komoditas Dijual
// @Security Authorization
// @Produce json
// @Param id_tempat_usaha query string false "Tempat usaha ID"
// @Param id_komoditas query string false "Komoditas ID"
// @Param status query int false "Status"
// @Param show-count query bool false "Show count"
// @Param offset query int false "Offset"
// @Param limit query int false "Limit"
// @Param order-by query string false "Order by"
// @Success 200 {object} dto.ResKomoditasDijualList
// @Router /admin/komoditas-dijual [get]
func (h *sihpHandler) GetKomoditasDijualByFilter(c *fiber.Ctx) error {
	var req dto.ReqGetKomoditasDijual
	_ = c.QueryParser(&req)
	res := h.usecase.GetKomoditasDijualByFilter(c.Context(), &req)
	return c.Status(res.Code).JSON(res)
}

// UpdateKomoditasDijual godoc
// @Summary Update komoditas dijual
// @Description Update komoditas dijual by ID
// @Tags Komoditas Dijual
// @Security Authorization
// @Accept json
// @Produce json
// @Param id path string true "Komoditas dijual ID"
// @Param payload body dto.ReqUpdateKomoditasDijual true "Update komoditas dijual payload"
// @Success 200 {object} dto.ResKomoditasDijualSingle
// @Router /admin/komoditas-dijual/{id} [put]
func (h *sihpHandler) UpdateKomoditasDijual(c *fiber.Ctx) error {
	var req dto.ReqUpdateKomoditasDijual
	_ = c.BodyParser(&req)
	res := h.usecase.UpdateKomoditasDijual(c.Context(), parseUUIDParam(c), &req)
	return c.Status(res.Code).JSON(res)
}

// DeleteKomoditasDijual godoc
// @Summary Delete komoditas dijual
// @Description Set komoditas dijual status to inactive
// @Tags Komoditas Dijual
// @Security Authorization
// @Produce json
// @Param id path string true "Komoditas dijual ID"
// @Success 200 {object} dto.ResKomoditasDijualSingle
// @Router /admin/komoditas-dijual/{id} [delete]
func (h *sihpHandler) DeleteKomoditasDijual(c *fiber.Ctx) error {
	res := h.usecase.DeleteKomoditasDijual(c.Context(), parseUUIDParam(c))
	return c.Status(res.Code).JSON(res)
}

// CreatePengumpulanData godoc
// @Summary Create pengumpulan data
// @Description Create draft pengumpulan data
// @Tags Pengumpulan Data
// @Security Authorization
// @Accept json
// @Produce json
// @Param payload body dto.ReqCreatePengumpulanData true "Create pengumpulan data payload"
// @Success 201 {object} dto.ResPengumpulanDataSingle
// @Router /admin/pengumpulan-data [post]
func (h *sihpHandler) CreatePengumpulanData(c *fiber.Ctx) error {
	var req dto.ReqCreatePengumpulanData
	_ = c.BodyParser(&req)
	res := h.usecase.CreatePengumpulanData(c.Context(), &req)
	return c.Status(res.Code).JSON(res)
}

// GetPengumpulanDataByID godoc
// @Summary Get pengumpulan data by ID
// @Description Get pengumpulan data detail by ID
// @Tags Pengumpulan Data
// @Security Authorization
// @Produce json
// @Param id path string true "Pengumpulan data ID"
// @Success 200 {object} dto.ResPengumpulanDataSingle
// @Router /admin/pengumpulan-data/{id} [get]
func (h *sihpHandler) GetPengumpulanDataByID(c *fiber.Ctx) error {
	res := h.usecase.GetPengumpulanDataByID(c.Context(), parseUUIDParam(c))
	return c.Status(res.Code).JSON(res)
}

// GetPengumpulanDataByFilter godoc
// @Summary Get pengumpulan data list
// @Description Get pengumpulan data list with filters
// @Tags Pengumpulan Data
// @Security Authorization
// @Produce json
// @Param id_pasar query string false "Pasar ID"
// @Param status query int false "Status"
// @Param from query string false "From date"
// @Param to query string false "To date"
// @Param show-count query bool false "Show count"
// @Param offset query int false "Offset"
// @Param limit query int false "Limit"
// @Param order-by query string false "Order by"
// @Success 200 {object} dto.ResPengumpulanDataList
// @Router /admin/pengumpulan-data [get]
func (h *sihpHandler) GetPengumpulanDataByFilter(c *fiber.Ctx) error {
	var req dto.ReqGetPengumpulanData
	_ = c.QueryParser(&req)
	res := h.usecase.GetPengumpulanDataByFilter(c.Context(), &req)
	return c.Status(res.Code).JSON(res)
}

// UpdatePengumpulanData godoc
// @Summary Update pengumpulan data
// @Description Update pengumpulan data by ID when draft
// @Tags Pengumpulan Data
// @Security Authorization
// @Accept json
// @Produce json
// @Param id path string true "Pengumpulan data ID"
// @Param payload body dto.ReqUpdatePengumpulanData true "Update pengumpulan data payload"
// @Success 200 {object} dto.ResPengumpulanDataSingle
// @Router /admin/pengumpulan-data/{id} [put]
func (h *sihpHandler) UpdatePengumpulanData(c *fiber.Ctx) error {
	var req dto.ReqUpdatePengumpulanData
	_ = c.BodyParser(&req)
	res := h.usecase.UpdatePengumpulanData(c.Context(), parseUUIDParam(c), &req)
	return c.Status(res.Code).JSON(res)
}

// DeletePengumpulanData godoc
// @Summary Delete pengumpulan data
// @Description Delete pengumpulan data when draft
// @Tags Pengumpulan Data
// @Security Authorization
// @Produce json
// @Param id path string true "Pengumpulan data ID"
// @Success 200 {object} dto.ResPengumpulanDataSingle
// @Router /admin/pengumpulan-data/{id} [delete]
func (h *sihpHandler) DeletePengumpulanData(c *fiber.Ctx) error {
	res := h.usecase.DeletePengumpulanData(c.Context(), parseUUIDParam(c))
	return c.Status(res.Code).JSON(res)
}

// FinalizePengumpulanData godoc
// @Summary Finalize pengumpulan data
// @Description Finalize draft and materialize harga pelaporan
// @Tags Pengumpulan Data
// @Security Authorization
// @Produce json
// @Param id path string true "Pengumpulan data ID"
// @Success 200 {object} dto.ResFinalizePengumpulanDataEnvelope
// @Router /admin/pengumpulan-data/{id}/finalize [post]
func (h *sihpHandler) FinalizePengumpulanData(c *fiber.Ctx) error {
	res := h.usecase.FinalizePengumpulanData(c.Context(), parseUUIDParam(c))
	return c.Status(res.Code).JSON(res)
}

// UploadPengumpulanTandaTangan godoc
// @Summary Upload pengumpulan signature
// @Description Upload enumerator signature image to MinIO and save URL in catatan
// @Tags Pengumpulan Data
// @Security Authorization
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Pengumpulan data ID"
// @Param tanda_tangan formData file true "Signature image (jpeg/png/webp, max 512KB)"
// @Success 200 {object} dto.ResPengumpulanDataSingle
// @Router /admin/pengumpulan-data/{id}/tanda-tangan [post]
func (h *sihpHandler) UploadPengumpulanTandaTangan(c *fiber.Ctx) error {
	file, err := c.FormFile("tanda_tangan")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"code":    fiber.StatusBadRequest,
			"message": "field tanda_tangan is required",
		})
	}

	f, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"code":    fiber.StatusBadRequest,
			"message": "failed to read uploaded file",
		})
	}
	defer f.Close()

	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		switch strings.ToLower(filepath.Ext(file.Filename)) {
		case ".jpg", ".jpeg":
			contentType = "image/jpeg"
		case ".png":
			contentType = "image/png"
		case ".webp":
			contentType = "image/webp"
		}
	}
	res := h.usecase.UploadPengumpulanTandaTangan(c.Context(), parseUUIDParam(c), f, file.Size, contentType)
	return c.Status(res.Code).JSON(res)
}

// CreateHargaRutin godoc
// @Summary Create harga rutin
// @Description Create harga rutin row for draft batch
// @Tags Harga Rutin
// @Security Authorization
// @Accept json
// @Produce json
// @Param payload body dto.ReqCreateHargaRutin true "Create harga rutin payload"
// @Success 201 {object} dto.ResHargaRutinSingle
// @Router /admin/harga-rutin [post]
func (h *sihpHandler) CreateHargaRutin(c *fiber.Ctx) error {
	var req dto.ReqCreateHargaRutin
	_ = c.BodyParser(&req)
	res := h.usecase.CreateHargaRutin(c.Context(), &req)
	return c.Status(res.Code).JSON(res)
}

// GetHargaRutinByID godoc
// @Summary Get harga rutin by ID
// @Description Get harga rutin detail by ID
// @Tags Harga Rutin
// @Security Authorization
// @Produce json
// @Param id path string true "Harga rutin ID"
// @Success 200 {object} dto.ResHargaRutinSingle
// @Router /admin/harga-rutin/{id} [get]
func (h *sihpHandler) GetHargaRutinByID(c *fiber.Ctx) error {
	res := h.usecase.GetHargaRutinByID(c.Context(), parseUUIDParam(c))
	return c.Status(res.Code).JSON(res)
}

// GetHargaRutinByFilter godoc
// @Summary Get harga rutin list
// @Description Get harga rutin list with filters
// @Tags Harga Rutin
// @Security Authorization
// @Produce json
// @Param id_pengumpulan_data query string false "Pengumpulan data ID"
// @Param id_tempat_usaha query string false "Tempat usaha ID"
// @Param id_komoditas query string false "Komoditas ID"
// @Param show-count query bool false "Show count"
// @Param offset query int false "Offset"
// @Param limit query int false "Limit"
// @Param order-by query string false "Order by"
// @Success 200 {object} dto.ResHargaRutinList
// @Router /admin/harga-rutin [get]
func (h *sihpHandler) GetHargaRutinByFilter(c *fiber.Ctx) error {
	var req dto.ReqGetHargaRutin
	_ = c.QueryParser(&req)
	res := h.usecase.GetHargaRutinByFilter(c.Context(), &req)
	return c.Status(res.Code).JSON(res)
}

// UpdateHargaRutin godoc
// @Summary Update harga rutin
// @Description Update harga rutin by ID for draft batch
// @Tags Harga Rutin
// @Security Authorization
// @Accept json
// @Produce json
// @Param id path string true "Harga rutin ID"
// @Param payload body dto.ReqUpdateHargaRutin true "Update harga rutin payload"
// @Success 200 {object} dto.ResHargaRutinSingle
// @Router /admin/harga-rutin/{id} [put]
func (h *sihpHandler) UpdateHargaRutin(c *fiber.Ctx) error {
	var req dto.ReqUpdateHargaRutin
	_ = c.BodyParser(&req)
	res := h.usecase.UpdateHargaRutin(c.Context(), parseUUIDParam(c), &req)
	return c.Status(res.Code).JSON(res)
}

// DeleteHargaRutin godoc
// @Summary Delete harga rutin
// @Description Delete harga rutin by ID for draft batch
// @Tags Harga Rutin
// @Security Authorization
// @Produce json
// @Param id path string true "Harga rutin ID"
// @Success 200 {object} dto.ResHargaRutinSingle
// @Router /admin/harga-rutin/{id} [delete]
func (h *sihpHandler) DeleteHargaRutin(c *fiber.Ctx) error {
	res := h.usecase.DeleteHargaRutin(c.Context(), parseUUIDParam(c))
	return c.Status(res.Code).JSON(res)
}

// GetHargaPelaporanByID godoc
// @Summary Get harga pelaporan by ID
// @Description Get harga pelaporan detail by ID
// @Tags Harga Pelaporan
// @Security Authorization
// @Produce json
// @Param id path string true "Harga pelaporan ID"
// @Success 200 {object} dto.ResHargaPelaporanSingle
// @Router /admin/harga-pelaporan/{id} [get]
func (h *sihpHandler) GetHargaPelaporanByID(c *fiber.Ctx) error {
	res := h.usecase.GetHargaPelaporanByID(c.Context(), parseUUIDParam(c))
	return c.Status(res.Code).JSON(res)
}

// GetHargaPelaporanByFilter godoc
// @Summary Get harga pelaporan list
// @Description Get harga pelaporan list with filters
// @Tags Harga Pelaporan
// @Security Authorization
// @Produce json
// @Param id_pasar query string false "Pasar ID"
// @Param id_komoditas query string false "Komoditas ID"
// @Param from query string false "From date"
// @Param to query string false "To date"
// @Param status query int false "Pengumpulan data status"
// @Param show-count query bool false "Show count"
// @Param offset query int false "Offset"
// @Param limit query int false "Limit"
// @Param order-by query string false "Order by"
// @Success 200 {object} dto.ResHargaPelaporanList
// @Router /admin/harga-pelaporan [get]
func (h *sihpHandler) GetHargaPelaporanByFilter(c *fiber.Ctx) error {
	var req dto.ReqGetHargaPelaporan
	_ = c.QueryParser(&req)
	res := h.usecase.GetHargaPelaporanByFilter(c.Context(), &req)
	return c.Status(res.Code).JSON(res)
}

package router

import (
	"github.com/thdoikn/sihp-be/internal/handler"
	masterdatarepository "github.com/thdoikn/sihp-be/internal/repository/master-data/implementation"
	transactiondatarepository "github.com/thdoikn/sihp-be/internal/repository/transaction-data/implementation"
	sihpserializer "github.com/thdoikn/sihp-be/internal/serializer/sihp/implementation"
	restmiddleware "github.com/thdoikn/sihp-be/internal/server/rest/middleware"
	sihpusecase "github.com/thdoikn/sihp-be/internal/usecase/sihp/implementation"
)

func AdminRouter(deps *Dependencies) {
	masterDataRepo := masterdatarepository.NewMasterDataRepository(deps.DB)
	transactionDataRepo := transactiondatarepository.NewTransactionDataRepository(deps.DB)
	serializer := sihpserializer.NewSIHPSerializer()
	usecase := sihpusecase.NewSIHPUsecase(masterDataRepo, transactionDataRepo, serializer, deps.Storage, deps.Cfg)
	h := handler.NewSIHPHandler(usecase)

	protected := deps.App.Group("/v1/admin", restmiddleware.JWTRequired(deps.Cfg))

	pasar := protected.Group("/pasar")
	pasar.Post("/", h.CreatePasar)
	pasar.Get("/:id", h.GetPasarByID)
	pasar.Get("/", h.GetPasarByFilter)
	pasar.Put("/:id", h.UpdatePasar)
	pasar.Delete("/:id", h.DeletePasar)

	komoditas := protected.Group("/komoditas")
	komoditas.Post("/", h.CreateKomoditas)
	komoditas.Post("/:id/gambar", h.UploadKomoditasGambar)
	komoditas.Get("/:id", h.GetKomoditasByID)
	komoditas.Get("/", h.GetKomoditasByFilter)
	komoditas.Put("/:id", h.UpdateKomoditas)

	tu := protected.Group("/tempat-usaha")
	tu.Post("/", h.CreateTempatUsaha)
	tu.Get("/:id", h.GetTempatUsahaByID)
	tu.Get("/", h.GetTempatUsahaByFilter)
	tu.Put("/:id", h.UpdateTempatUsaha)
	tu.Delete("/:id", h.DeleteTempatUsaha)

	kd := protected.Group("/komoditas-dijual")
	kd.Post("/", h.CreateKomoditasDijual)
	kd.Get("/:id", h.GetKomoditasDijualByID)
	kd.Get("/", h.GetKomoditasDijualByFilter)
	kd.Put("/:id", h.UpdateKomoditasDijual)
	kd.Delete("/:id", h.DeleteKomoditasDijual)

	pd := protected.Group("/pengumpulan-data")
	pd.Post("/", h.CreatePengumpulanData)
	pd.Get("/:id", h.GetPengumpulanDataByID)
	pd.Get("/", h.GetPengumpulanDataByFilter)
	pd.Put("/:id", h.UpdatePengumpulanData)
	pd.Delete("/:id", h.DeletePengumpulanData)
	pd.Post("/:id/finalize", h.FinalizePengumpulanData)

	hr := protected.Group("/harga-rutin")
	hr.Post("/", h.CreateHargaRutin)
	hr.Get("/:id", h.GetHargaRutinByID)
	hr.Get("/", h.GetHargaRutinByFilter)
	hr.Put("/:id", h.UpdateHargaRutin)
	hr.Delete("/:id", h.DeleteHargaRutin)

	hp := protected.Group("/harga-pelaporan")
	hp.Get("/:id", h.GetHargaPelaporanByID)
	hp.Get("/", h.GetHargaPelaporanByFilter)
}

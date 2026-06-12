package router

import (
	"github.com/thdoikn/sihp-be/internal/handler"
	masterdatarepository "github.com/thdoikn/sihp-be/internal/repository/master-data/implementation"
	transactiondatarepository "github.com/thdoikn/sihp-be/internal/repository/transaction-data/implementation"
	sihpserializer "github.com/thdoikn/sihp-be/internal/serializer/sihp/implementation"
	sihpusecase "github.com/thdoikn/sihp-be/internal/usecase/sihp/implementation"
)

func PublicRouter(deps *Dependencies) {
	masterDataRepo := masterdatarepository.NewMasterDataRepository(deps.DB)
	transactionDataRepo := transactiondatarepository.NewTransactionDataRepository(deps.DB)
	serializer := sihpserializer.NewSIHPSerializer()
	usecase := sihpusecase.NewSIHPUsecase(masterDataRepo, transactionDataRepo, serializer, deps.Storage, deps.Cfg)
	h := handler.NewSIHPHandler(usecase)

	group := deps.App.Group("/v1/public")
	group.Get("/overview", h.PublicOverview)
	group.Get("/komoditas", h.PublicKomoditas)
	group.Get("/komoditas/:id", h.PublicKomoditasDetail)
	group.Get("/komoditas/:id/trend", h.PublicKomoditasTrend)
	group.Get("/pasar", h.PublicPasar)
	group.Get("/pasar/:id", h.PublicPasarDetail)
	group.Get("/tempat-usaha", h.PublicTempatUsaha)
	group.Get("/tempat-usaha/:id", h.PublicTempatUsahaDetail)
}

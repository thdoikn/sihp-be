package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thdoikn/sihp-be/config"
	"github.com/thdoikn/sihp-be/pkg/storage"
	"gorm.io/gorm"
)

type Dependencies struct {
	App     *fiber.App
	DB      *gorm.DB
	Cfg     *config.Config
	Storage storage.KomoditasStorage
}

func NewDependencies(app *fiber.App, db *gorm.DB, cfg *config.Config, komoditasStorage storage.KomoditasStorage) *Dependencies {
	return &Dependencies{
		App:     app,
		DB:      db,
		Cfg:     cfg,
		Storage: komoditasStorage,
	}
}

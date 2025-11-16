package main

import (
	endpoints "internal-transfers/constants"
	database "internal-transfers/db"
	"internal-transfers/db/models"
	"internal-transfers/dtos"
	account_service "internal-transfers/service"
	"internal-transfers/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func init() {
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.ReadInConfig()
	zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))
	db, err := database.GetDb()
	if err != nil {
		zap.L().Panic("Error occurred while trying to migrate", zap.Error(err))
	}
	zap.L().Debug("Migrating schema")
	db.AutoMigrate(&models.Account{}, &models.Transaction{})
}

func main() {
	app := fiber.New()

	app.Get(endpoints.HEALTH_CHECK, func(c *fiber.Ctx) error {
		return c.SendString("Welcome to Internal Transfers API")
	})

	app.Get(endpoints.GET_ACCOUNT, func(c *fiber.Ctx) error {
		accountId, err := c.ParamsInt("accountId")
		if err != nil {
			return fiber.NewError(400, "Invalid type for accountId")
		}
		response, err := account_service.GetAccount(int64(accountId))
		if err != nil {
			return err
		}
		return c.Status(200).JSON(response)
	})

	app.Post(endpoints.CREATE_ACCOUNT, func(c *fiber.Ctx) error {
		body := new(dtos.CreateAccountDto)
		if err := c.BodyParser(body); err != nil {
			zap.L().Error("Error occurred while parsing body", zap.Error(err))
			return fiber.ErrBadRequest
		}
		if validation_err := utils.ValidateStruct(body); len(validation_err) > 0 {
			return c.Status(400).JSON(validation_err)
		}
		if err := account_service.CreateAccount(body); err != nil {
			zap.L().Error("CreateAccount failed", zap.Error(err))
			return err
		}
		return c.SendStatus(201)
	})

	app.Post(endpoints.TRANSFER, func(c *fiber.Ctx) error {
		body := new(dtos.CreateTransactionDto)
		if err := c.BodyParser(body); err != nil {
			zap.L().Error("Error occurred while parsing body", zap.Error(err))
			return fiber.ErrInternalServerError
		}
		if validation_err := utils.ValidateStruct(body); len(validation_err) > 0 {
			return c.Status(400).JSON(validation_err)
		}
		if err := account_service.TransferFunds(body); err != nil {
			return err
		}
		return c.SendStatus(201)
	})

	port := viper.GetString("app.port")
	zap.L().Info("Starting server", zap.String("port", port))
	if err := app.Listen(":" + port); err != nil {
		zap.L().Fatal("Failed to start server", zap.Error(err))
	}
}

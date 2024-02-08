package chatapi

import "github.com/gofiber/fiber/v2"

func handlePostMessage(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{})
}

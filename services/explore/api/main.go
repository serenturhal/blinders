package exploreapi

import (
	"context"
	"fmt"

	"blinders/packages/auth"
	"blinders/packages/db/models"
	"blinders/packages/explore"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	App         *fiber.App
	Auth        auth.Manager
	Core        explore.Explorer
	RedisClient *redis.Client
	MatchCh     chan models.MatchInfo
}

func NewService(
	app *fiber.App,
	auth auth.Manager,
	exploreCore explore.Explorer,
	redisClient *redis.Client,
) *Service {
	return &Service{
		App:         app,
		Auth:        auth,
		Core:        exploreCore,
		RedisClient: redisClient,
		MatchCh:     make(chan models.MatchInfo),
	}
}

func (s *Service) InitRoute() {
	s.App.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("service healthy")
	})

	matchRoute := s.App.Group("/match", auth.FiberAuthMiddleware(s.Auth))
	matchRoute.Get("/suggest", s.HandleGetMatches)
	// Temporarily expose this method, it must be call internal, or we will listen to user update-match-information event
	matchRoute.Post("/", s.HandleAddUserMatch)
}

func (s *Service) Loop() {
	for match := range s.MatchCh {
		value := map[string]string{
			"id": match.UserID.Hex(),
		}

		// fire user-match-entry-created event to the stream
		if err := s.RedisClient.XAdd(context.Background(), &redis.XAddArgs{
			Stream:     "match:embed",
			NoMkStream: false,
			Values:     value,
		}).Err(); err != nil {
			fmt.Println("error", err)
			break
		}
		fmt.Println("fired event to stream", value)
	}
}

/*
This service is responsible for notifying any event to users via websocket or push notification
*/
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"blinders/packages/apigateway"
	"blinders/packages/session"
	"blinders/packages/transport"
	"blinders/packages/utils"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/redis/go-redis/v9"
)

var (
	APIGatewayClient *apigateway.Client
	SessionManager   *session.Manager
)

func init() {
	APIGatewayClient = apigateway.NewClient(context.Background(), apigateway.CustomEndpointResolve{
		Domain:     os.Getenv("API_GATEWAY_DOMAIN"),
		PathPrefix: os.Getenv("API_GATEWAY_PATH_PREFIX"),
	})

	SessionManager = session.NewManager(redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}))
}

func HandleRequest(_ context.Context, payload []byte) error {
	event, err := utils.ParseJSON[transport.Event](payload)
	if err != nil {
		log.Println("can not parse request payload, require type in payload:", err)
		return err
	}

	log.Println("handle event:", event.Type)

	switch event.Type {
	case transport.AddFriend:
		event, err := utils.ParseJSON[transport.AddFriendEvent](payload)
		if err != nil {
			log.Println("can not parse request payload:", err)
			return err
		}
		userConID, err := SessionManager.GetSessions(event.UserID)
		if err != nil {
			log.Println("can not get session:", err)
			return err
		}

		wg := sync.WaitGroup{}
		for _, conID := range userConID {
			wg.Add(1)
			go func(conID string, payload []byte) {
				_ = APIGatewayClient.Publish(context.Background(), conID, payload)
				wg.Done()
			}(conID, payload)
		}
		wg.Wait()
	default:
		log.Print("does not support event type:", event.Type)
	}

	return nil
}

func main() {
	log.Println(APIGatewayClient)
	lambda.Start(HandleRequest)
}

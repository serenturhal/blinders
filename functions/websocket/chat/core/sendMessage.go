package wschat

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"blinders/packages/db/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func HandleSendMessage(
	rawUserID string, // for all case, userID must be valid and user existed
	connectionID string,
	payload UserSendMessagePayload,
) (<-chan *DistributeEvent, <-chan error, error) {
	dCh := make(chan *DistributeEvent)
	errCh := make(chan error)
	wg := sync.WaitGroup{}

	userID, _ := primitive.ObjectIDFromHex(rawUserID)
	conversationID, err := primitive.ObjectIDFromHex(payload.ConversationID)
	if err != nil {
		return dCh, errCh, fmt.Errorf("invalid conversationId: %s", payload.ConversationID)
	}
	replyTo, err := primitive.ObjectIDFromHex(payload.ReplyTo)
	if err != nil {
		return dCh, errCh, fmt.Errorf("invalid replyTo: %s", payload.ReplyTo)
	}

	conversation, err := queryConversationOfUser(conversationID, userID)
	if err != nil {
		return dCh, errCh, fmt.Errorf("failed to query conversation: %v", err)
	}

	if !replyTo.IsZero() {
		err := checkValidReplyTo(replyTo, conversationID)
		if err != nil {
			return dCh, errCh, fmt.Errorf("cannot reply to message %s", payload.ReplyTo)
		}
	}

	message := app.DB.Messages.ConstructNewMessage(
		userID,
		conversationID,
		replyTo,
		payload.Content,
	)

	wg.Add(1)
	go func() {
		distributeAckMessage(message, connectionID, dCh)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		distributeMessageToRecipients(message, *conversation, dCh)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		_, err := app.DB.Messages.InsertNewMessage(message)
		if err != nil {
			errCh <- fmt.Errorf("[important] failed to insert message %v", err)
		}
		wg.Done()
	}()

	go func() {
		wg.Wait()
		dCh <- nil
	}()

	return dCh, errCh, nil
}

// query conversation by id
// just send conversation via channel if user is a member
func queryConversationOfUser(
	conversationID primitive.ObjectID,
	userID primitive.ObjectID,
) (*models.Conversation, error) {
	conversation, err := app.DB.Conversations.GetConversationByID(conversationID)
	if err != nil {
		return nil, err
	}

	for _, m := range conversation.Members {
		if m.UserID == userID {
			return &conversation, nil
		}
	}

	return nil, fmt.Errorf("user %s is not a member of conversation %s", userID.Hex(), conversationID.Hex())
}

func checkValidReplyTo(replyTo primitive.ObjectID, conversationID primitive.ObjectID) error {
	repliedMessage, err := app.DB.Messages.GetMessageByID(replyTo)
	if err != nil {
		return err
	} else if repliedMessage.ConversationID != conversationID {
		return fmt.Errorf("reply to message %s is not in conversation %s", replyTo.Hex(), conversationID.Hex())
	}

	return nil
}

func distributeAckMessage(
	message models.Message,
	connectionID string,
	dCh chan *DistributeEvent,
) {
	dCh <- &DistributeEvent{
		ConnectionID: connectionID,
		Payload: ServerAckSendMessagePayload{
			ChatEvent: ChatEvent{Type: ServerAckSendMessage},
			ResolveID: message.ID.Hex(),
			Message:   message,
		},
	}
}

func distributeMessageToRecipients(
	message models.Message,
	conversation models.Conversation,
	dCh chan *DistributeEvent,
) {
	wg := sync.WaitGroup{}
	for _, m := range conversation.Members {
		if m.UserID == message.SenderID {
			continue
		}
		wg.Add(1)
		// TODO: use go 1.22 to resolve loop with goroutine
		go func(m models.Member) {
			sessions, err := app.Session.GetSessions(m.UserID.Hex())
			if err == nil {
				log.Println("failed to query sessions for user", m.UserID.Hex())
				wg.Done()
				return
			}

			for _, s := range sessions {
				connectionID := strings.Split(s, ":")[1]
				dCh <- &DistributeEvent{
					ConnectionID: connectionID,
					Payload: ServerSendMessagePayload{
						ChatEvent: ChatEvent{Type: ServerSendMessage},
						Message:   message,
					},
				}
			}

			wg.Done()
		}(m)
	}

	wg.Wait()
}

// updateMessageStatus
// TODO: receive events
// conversationID, messageID, status -> ""
// store to database
// distribute to user

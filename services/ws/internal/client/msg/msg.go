package msg

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"MyMessenger/services/ws/internal/config"
	"MyMessenger/services/ws/internal/service"

	"github.com/google/uuid"
)

const (
	clientTimeout = 3 * time.Second
)

type MsgClient struct {
	client  *http.Client
	fullURL string
}

func NewMsgClient(conf *config.MsgClientConfig) *MsgClient {
	return &MsgClient{
		client:  &http.Client{Timeout: clientTimeout},
		fullURL: conf.FullURL,
	}
}

func (c *MsgClient) GetAllUserChats(ctx context.Context, userId, accessToken string) ([]uuid.UUID, error) {
	var chats []uuid.UUID

	r, err := http.NewRequestWithContext(ctx, "GET", c.fullURL, nil)
	if err != nil {
		log.Printf("GetAllUserChats: Trouble with NewRequest, err: %q", err)
		return chats, service.ErrInternal
	}
	r.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.client.Do(r)
	if err != nil {
		log.Printf("GetAllUserChats: Trouble with client.Do(r), err: %q", err)
		return chats, service.ErrInternal
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Printf("GetAllUserChats: Trouble with resp.Body.Close(), err: %q", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			return chats, service.ErrUnauthorized
		}
		log.Printf("GetAllUserChats: Trouble with client.Do(r), StatusCode: %d", resp.StatusCode)
		return chats, service.ErrInternal
	}

	err = json.NewDecoder(resp.Body).Decode(&chats)
	if err != nil {
		log.Printf("GetAllUserChats: Trouble with decoding, err: %q", err)
		return chats, service.ErrInternal
	}

	return chats, nil
}

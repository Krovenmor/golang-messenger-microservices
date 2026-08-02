package ws

import "context"

type Subscriber interface {
	Subscribe(ctx context.Context, channel string) (<-chan []byte, func(), error)
}

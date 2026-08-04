package service

import "context"

type Subscriber interface {
	Subscribe(ctx context.Context, userId string) (<-chan []byte, func(), error)
}

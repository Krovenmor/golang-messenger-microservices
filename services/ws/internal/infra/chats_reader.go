package infra

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
)

const (
	bufferLength = 10
)

type chatChannel struct {
	readerContext context.Context
	readerCancel  context.CancelFunc

	ch     <-chan []byte
	cancel func()

	readersMu sync.Mutex
	readers   map[uuid.UUID]chan []byte
}

type chatsReader struct {
	mu       sync.Mutex
	channels map[string]*chatChannel

	sub *RedisSubscriber
}

func newChatsReader(sub *RedisSubscriber) *chatsReader {
	return &chatsReader{
		sub:      sub,
		channels: make(map[string]*chatChannel),
	}
}

func (r *chatsReader) tryRemove(chatId string, channel *chatChannel) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	ch, isExists := r.channels[chatId]
	if isExists {
		if channel == ch {
			ch.readersMu.Lock()
			defer ch.readersMu.Unlock()

			if len(ch.readers) == 0 {
				delete(r.channels, chatId)
				return nil
			} else {
				log.Printf("tryRemove: Chat %v, There still %d subscribers", chatId, len(ch.readers))
				ch.readerContext, ch.readerCancel = context.WithCancel(context.Background())
				return fmt.Errorf("can't remove")
			}

		} else {
			log.Printf("tryRemove: Chat %v, channel != ch", chatId)
			return nil
		}
	} else {
		log.Printf("tryRemove: Chat %v, doesn't exists", chatId)
		return nil
	}
}

func (r *chatsReader) remove(chatId string, channel *chatChannel) {
	r.mu.Lock()
	defer r.mu.Unlock()

	ch, isExists := r.channels[chatId]
	if isExists && ch == channel {
		ch.readersMu.Lock()
		defer ch.readersMu.Unlock()

		for _, readerCh := range ch.readers {
			close(readerCh)
		}

		delete(r.channels, chatId)
	}
}

func (r *chatsReader) startReading(chatId string, channel *chatChannel) {
	for {
		select {
		case <-channel.readerContext.Done():
			err := r.tryRemove(chatId, channel)
			if err == nil {
				channel.cancel()
				return
			}

		case msg, ok := <-channel.ch:
			if !ok {
				r.remove(chatId, channel)
				return
			}
			channel.readersMu.Lock()
			for id, ch := range channel.readers {
				select {
				case ch <- msg:
				default:
					log.Printf("Trouble with writing to %v", id)
				}
			}
			channel.readersMu.Unlock()
		}
	}
}

func (r *chatsReader) newReader(channel *chatChannel) (<-chan []byte, func()) {
	readerId := uuid.New()
	readerChan := make(chan []byte, bufferLength)

	channel.readersMu.Lock()
	channel.readers[readerId] = readerChan
	channel.readersMu.Unlock()

	cancel := func() {
		var cancelFn context.CancelFunc

		channel.readersMu.Lock()
		delete(channel.readers, readerId)
		if len(channel.readers) == 0 {
			cancelFn = channel.readerCancel
		}
		channel.readersMu.Unlock()

		if cancelFn != nil {
			cancelFn()
		}
	}

	return readerChan, cancel
}

func (r *chatsReader) newSub(chatId string) (*chatChannel, error) {
	ch, cancel, err := r.sub.SubscribeOnChatEventsInternal(context.Background(), chatId)
	if err != nil {
		return nil, err
	}
	ctx, cCtx := context.WithCancel(context.Background())
	channel := &chatChannel{
		ch:            ch,
		cancel:        cancel,
		readers:       make(map[uuid.UUID]chan []byte),
		readerContext: ctx,
		readerCancel:  cCtx,
	}
	r.channels[chatId] = channel
	return channel, nil
}

func (r *chatsReader) Subscribe(chatId string) (<-chan []byte, func(), error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	channel, isExists := r.channels[chatId]
	if !isExists {
		var err error
		channel, err = r.newSub(chatId)
		if err != nil {
			return nil, nil, err
		}
		go r.startReading(chatId, channel)
	}

	ch, close := r.newReader(channel)
	return ch, close, nil
}

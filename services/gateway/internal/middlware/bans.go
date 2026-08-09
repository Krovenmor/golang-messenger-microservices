package middlware

import (
	"MyMessenger/pkg/broker"
	"log"
	"time"
)

type BanReason string

const (
	TooManyRequestsUser BanReason = "too_many_user_reqs"
	TooManyRequestsIp   BanReason = "too_many_ip_reqs"
)

func FromBroker(reason broker.BanReason) BanReason {
	switch reason {
	case broker.TooManyRequests:
		return TooManyRequestsUser
	default:
		log.Printf("FromBroker: broker.BanReason=%q not implemented", reason)
		return TooManyRequestsUser
	}
}

func (m *Middleware) getBanTtl(reason BanReason) time.Duration {
	switch reason {
	case TooManyRequestsUser:
		return m.conf.BanTooManyRequestsDurationUser
	case TooManyRequestsIp:
		return m.conf.BanTooManyRequestsDurationIp
	default:
		log.Printf("calcBanTime: not implemented BanReason: %q", reason)
		return DefBanTime
	}
}

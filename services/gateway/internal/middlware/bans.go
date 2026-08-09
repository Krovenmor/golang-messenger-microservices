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

func (m *Middleware) getBanTime(reason BanReason) time.Time {
	var banDur time.Duration
	switch reason {
	case TooManyRequestsUser:
		banDur = m.conf.BanTooManyRequestsDurationUser
	case TooManyRequestsIp:
		banDur = m.conf.BanTooManyRequestsDurationIp
	default:
		log.Printf("calcBanTime: not implemented BanReason: %q", reason)
		banDur = DefBanTime
	}
	return time.Now().Add(banDur)
}

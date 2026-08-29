package service

import (
	"fmt"
)

func (m *messageService) checkQ(q int) error {
	if q > m.conf.MaxQuantity {
		return fmt.Errorf("too big quantity")
	}
	if q < m.conf.MinQuantity {
		return fmt.Errorf("too small quantity")
	}
	return nil
}

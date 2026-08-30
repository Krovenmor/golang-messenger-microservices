package service

import "strings"

func (s *profileService) checkString(str string, min, max int) (string, error) {
	st := strings.TrimSpace(str)
	l := len(st)
	if l < min || l > max {
		return "", ErrBadData
	}
	return st, nil
}

func (s *profileService) checkName(name string) (string, error) {
	return s.checkString(name, s.conf.ProfileMinNameLen, s.conf.ProfileMaxNameLen)
}

func (s *profileService) checkBio(bio string) (string, error) {
	return s.checkString(bio, s.conf.ProfileMinBioLen, s.conf.ProfileMaxBioLen)
}

package config

import "time"

type ConfigReader struct {
	err error
}

func NewConfigReader() *ConfigReader {
	return &ConfigReader{}
}

func (r *ConfigReader) GetString(key string) string {
	if r.err != nil {
		return ""
	}
	var val string
	val, r.err = GetEnvVar(key)
	return val
}

func (r *ConfigReader) GetInt(key string) int {
	if r.err != nil {
		return -1
	}
	var val int
	val, r.err = GetEnvVarInt(key)
	return val
}

func (r *ConfigReader) GetDuration(key string) time.Duration {
	if r.err != nil {
		return -1
	}
	var val time.Duration
	val, r.err = GetEnvVarDuration(key)
	return val
}

func (r *ConfigReader) GetBytes(key string) []byte {
	if r.err != nil {
		return nil
	}
	var val []byte
	val, r.err = GetEnvVarBytes(key)
	return val
}

func (r *ConfigReader) GetFloat64(key string) float64 {
	if r.err != nil {
		return -1
	}
	var val float64
	val, r.err = GetEnvVarFloat64(key)
	return val
}

func (r *ConfigReader) Err() error {
	return r.err
}

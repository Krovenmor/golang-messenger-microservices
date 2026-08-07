package config

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

func (r *ConfigReader) Err() error {
	return r.err
}

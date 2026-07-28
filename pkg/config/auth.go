package config

type JwtCheckerConf struct {
	PubKey []byte
}

func GetJwtCheckerConf() (JwtCheckerConf, error) {
	pubKey, err := GetEnvVarBytes("JWT_PUBKEY_PATH")
	if err != nil {
		return JwtCheckerConf{}, err
	}
	return JwtCheckerConf{
		PubKey: pubKey,
	}, nil
}

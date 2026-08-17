package config

type JwtCheckerConf struct {
	PubKey    []byte
	ATokenLen int
}

func GetJwtCheckerConf() (JwtCheckerConf, error) {
	pubKey, err := GetEnvVarBytes("JWT_PUBKEY_PATH")
	if err != nil {
		return JwtCheckerConf{}, err
	}
	aTokenL, err := GetEnvVarInt("JWT_ACCESS_TOKEN_LEN")
	if err != nil {
		return JwtCheckerConf{}, err
	}
	return JwtCheckerConf{
		PubKey:    pubKey,
		ATokenLen: aTokenL,
	}, nil
}

package tests

import (
	"MyMessenger/services/auth/internal/config"
	"MyMessenger/services/auth/internal/service"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	prvKeyFile = "privateTests.pem"
	pubKeyFile = "publicTests.pem"

	accessTTL  = time.Second * 1
	refreshTTL = time.Second * 2
)

func printRepo(t *testing.T, repo *RepoMock) {
	if repo == nil {
		return
	}
	for idx, entry := range repo.commands {
		t.Logf("%d: %q\n", idx, entry)
	}
}

func getConfig() (*config.AuthConfig, error) {
	prvKey, err := os.ReadFile(prvKeyFile)
	if err != nil {
		return nil, err
	}
	pubKey, err := os.ReadFile(pubKeyFile)
	if err != nil {
		return nil, err
	}
	return &config.AuthConfig{
		MinPassLength: 5,
		MaxPassLength: 10,

		MinLoginLength: 5,
		MaxLoginLength: 10,

		PrvKey: prvKey,
		PubKey: pubKey,

		AccessTokenTTL:  accessTTL,
		RefreshTokenTTL: refreshTTL,
	}, nil
}

func getAuth(t *testing.T) (service.AuthService, *RepoMock) {
	repo := NewMockRepo()
	conf, err := getConfig()
	if err != nil {
		t.Fatal(err)
	}

	auth := service.NewAuth(repo, *conf)
	if auth == nil {
		t.Fatal("auth == nil")
	}

	return service.AuthService(auth), repo
}

func TestAuth(t *testing.T) {
	auth, repo := getAuth(t)

	testLogin, testPass := "aboba", "password"
	ctx := context.Background()

	tokens, err := auth.LogIn(ctx, testLogin, testPass)
	if err == nil {
		t.Logf("Commands slice: %v\n", repo.commands)
		t.Fatal("After LogIn: err == nil, when user not exists")
	}
	if tokens != nil {
		t.Logf("Commands slice: %v\n", repo.commands)
		t.Fatal("After LogIn: tokens != nil, when user not exists")
	}

	err = auth.Register(ctx, testLogin, testPass)
	if err != nil {
		t.Logf("Commands slice: %v\n", repo.commands)
		t.Fatalf("After Register: err != nil, when user not exists, err: %q", err)
	}

	tokens, err = auth.LogIn(ctx, testLogin, testPass)
	if err != nil {
		t.Logf("Commands slice: %v\n", repo.commands)
		t.Fatalf("After LogIn: err != nil, when user exists, err: %q", err)
	}
	if tokens == nil {
		t.Logf("Commands slice: %v\n", repo.commands)
		t.Fatal("After LogIn: tokens == nil, when user exists")
	}

	err = auth.IsValidAccess(ctx, tokens.AccessToken)
	if err != nil {
		t.Logf("Commands slice: %v\n", repo.commands)
		t.Fatalf("After IsValidAccess: err != nil, when token must be valid, err: %q", err)
	}

	nTokens, err := auth.UpdateTokens(ctx, tokens.RefreshToken)
	if err != nil {
		t.Logf("Commands slice: %v\n", repo.commands)
		t.Fatalf("After UpdateTokens: err != nil, when user exists, err: %q", err)
	}
	if nTokens == nil {
		t.Logf("Commands slice: %v\n", repo.commands)
		t.Fatal("After UpdateTokens: nTokens == nil, when user exists")
	}
	tokens = nTokens

	time.Sleep(refreshTTL)
	nTokens, err = auth.UpdateTokens(ctx, tokens.RefreshToken)
	if err == nil {
		t.Logf("Commands slice: %v\n", repo.commands)
		t.Fatalf("After UpdateTokens: err == nil, when rToken must be dead by TTL")
	}
	if nTokens != nil {
		t.Logf("Commands slice: %v\n", repo.commands)
		t.Fatal("After UpdateTokens: nTokens != nil, when rToken dead by TTL")
	}

	err = auth.IsValidAccess(ctx, tokens.AccessToken)
	if err == nil {
		t.Logf("Commands slice: %v\n", repo.commands)
		t.Fatalf("After IsValidAccess: err == nil, when token must be not valid")
	}
}

func TestRegister_Validation(t *testing.T) {
	tests := []struct {
		name     string
		login    string
		password string
		wantErr  bool
	}{
		{
			name:     "valid registration",
			login:    "normal",
			password: "validPass",
			wantErr:  false,
		},
		{
			name:     "login too short",
			login:    "usr",
			password: "validPass",
			wantErr:  true,
		},
		{
			name:     "login too long",
			login:    "verylongloginthatexceedsmaxlength",
			password: "validPass",
			wantErr:  true,
		},
		{
			name:     "password too short",
			login:    "normal",
			password: "123",
			wantErr:  true,
		},
		{
			name:     "empty credentials",
			login:    "",
			password: "",
			wantErr:  true,
		},
		{
			name:     "spaces login",
			login:    "     ",
			password: "normal",
			wantErr:  true,
		},
		{
			name:     "spaces pass",
			login:    "normal",
			password: "     ",
			wantErr:  true,
		},
	}

	auth, _ := getAuth(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			err := auth.Register(ctx, tt.login, tt.password)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func generateTestRSAKeys(t *testing.T) (*rsa.PrivateKey, *rsa.PublicKey) {
	t.Helper()
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate private key: %v", err)
	}
	return privateKey, &privateKey.PublicKey
}

func TestJWT_Validation(t *testing.T) {
	validPrvKey, validPubKey := generateTestRSAKeys(t)

	attackerPrvKey, _ := generateTestRSAKeys(t)

	testUserID := uuid.New()

	validToken, err := service.GenToken(testUserID, 1*time.Hour, validPrvKey)
	if err != nil {
		t.Fatalf("failed to create valid token: %v", err)
	}

	expiredToken, err := service.GenToken(testUserID, -1*time.Minute, validPrvKey)
	if err != nil {
		t.Fatalf("failed to create expired token: %v", err)
	}

	forgedToken, err := service.GenToken(testUserID, 1*time.Hour, attackerPrvKey)
	if err != nil {
		t.Fatalf("failed to create forged token: %v", err)
	}

	tamperedToken := validToken[:len(validToken)-5] + "XXXXX"

	tests := []struct {
		name       string
		tokenStr   string
		pubKey     *rsa.PublicKey
		wantErr    bool
		wantUserID string
	}{
		{
			name:       "success: valid token",
			tokenStr:   validToken,
			pubKey:     validPubKey,
			wantErr:    false,
			wantUserID: testUserID.String(),
		},
		{
			name:     "fail: expired token",
			tokenStr: expiredToken,
			pubKey:   validPubKey,
			wantErr:  true,
		},
		{
			name:     "fail: forged token (signed by attacker key)",
			tokenStr: forgedToken,
			pubKey:   validPubKey,
			wantErr:  true,
		},
		{
			name:     "fail: tampered token signature",
			tokenStr: tamperedToken,
			pubKey:   validPubKey,
			wantErr:  true,
		},
		{
			name:     "fail: random garbage string",
			tokenStr: "invalid.jwt.string",
			pubKey:   validPubKey,
			wantErr:  true,
		},
		{
			name:     "fail: empty string",
			tokenStr: "",
			pubKey:   validPubKey,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := service.IsValidToken(tt.tokenStr, tt.pubKey)

			if (err != nil) != tt.wantErr {
				t.Fatalf("IsValidToken() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if claims == nil {
					t.Fatal("expected claims, got nil")
				}
				if claims.Subject != tt.wantUserID {
					t.Errorf("got userID = %s, want %s", claims.Subject, tt.wantUserID)
				}
			}
		})
	}
}

func TestAuthService_TokenSecurityAndEdgeCases(t *testing.T) {
	auth, repo := getAuth(t)
	ctx := context.Background()

	testLogin, testPass := "securUser", "securePass"
	err := auth.Register(ctx, testLogin, testPass)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	tokens, err := auth.LogIn(ctx, testLogin, testPass)
	if err != nil || tokens == nil {
		t.Fatalf("LogIn failed: %v", err)
	}

	t.Run("fail: update tokens using access token instead of refresh", func(t *testing.T) {
		_, err := auth.UpdateTokens(ctx, tokens.AccessToken)
		if err == nil {
			t.Fatal("expected error when passing AccessToken into UpdateTokens, got nil")
		}
	})

	var newTokens *service.Tokens
	t.Run("success: valid refresh token rotation", func(t *testing.T) {
		var err error
		newTokens, err = auth.UpdateTokens(ctx, tokens.RefreshToken)
		if err != nil || newTokens == nil {
			t.Fatalf("UpdateTokens failed with valid refresh token: %v", err)
		}
	})

	t.Run("fail: reuse old (rotated) refresh token", func(t *testing.T) {
		_, err := auth.UpdateTokens(ctx, tokens.RefreshToken)
		if err == nil {
			printRepo(t, repo)
			t.Fatal("expected error when reusing old refresh token, got nil")
		}
	})

	t.Run("fail: generate token for empty uuid", func(t *testing.T) {
		conf, err := getConfig()
		if err != nil {
			t.Fatal(err)
		}
		prvKey, err := jwt.ParseRSAPrivateKeyFromPEM(conf.PrvKey)

		nilUUIDToken, err := service.GenToken(uuid.Nil, 1*time.Hour, prvKey)
		if err == nil {
			errValid := auth.IsValidAccess(ctx, nilUUIDToken)
			if errValid == nil {
				t.Error("expected IsValidAccess to reject token with uuid.Nil")
			}
		}
	})
}

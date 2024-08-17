package tests

import (
	"sso-service/tests/suite"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/brianvoe/gofakeit/v7"
	ssov1 "github.com/sollidy/go-sso-protos/gen/go/proto/sso"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	emptyAppId     = 0
	appId          = 1
	appSecret      = "test_secret"
	passDefaultLen = 10
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)
	pass := gofakeit.Password(true, true, true, true, false, passDefaultLen)
	email := gofakeit.Email()

	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: pass,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respLogin, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: pass,
		AppId:    appId,
	})
	require.NoError(t, err)

	loginTime := time.Now()

	token := respLogin.GetToken()
	require.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	require.True(t, ok)
	assert.Equal(t, int64(respReg.GetUserId()), int64(claims["uid"].(float64)))
	assert.Equal(t, email, claims["email"])
	assert.Equal(t, appId, int(claims["app_id"].(float64)))

	const deltaSeconds = 2

	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)
}

func TestRegisterLogin_DublicatedRegistration(t *testing.T) {
	ctx, st := suite.New(t)
	pass := gofakeit.Password(true, true, true, true, false, passDefaultLen)
	email := gofakeit.Email()

	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: pass,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respReg, err = st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: pass,
	})

	require.Error(t, err)
	assert.Empty(t, respReg.GetUserId())
	assert.ErrorContains(t, err, "user already exists")
}

func TestRegisterLogin_TableTest(t *testing.T) {
	ctx, suite := suite.New(t)

	testCases := []struct {
		name          string
		email         string
		password      string
		expectedError string
	}{
		{
			"Both empty",
			"",
			"",
			"email is required",
		},
		{
			"Empty email",
			"",
			gofakeit.Password(true, true, true, true, false, passDefaultLen),
			"email is required",
		},
		{
			"Invalid email",
			"test",
			gofakeit.Password(true, true, true, true, false, passDefaultLen),
			"invalid email",
		},
		{
			"Empty password",
			gofakeit.Email(),
			"",
			"password is required",
		},
	}

	t.Run("Parallel", func(t *testing.T) {
		t.Parallel()
		for _, testCase := range testCases {
			testCase := testCase
			t.Run(testCase.name, func(t *testing.T) {
				_, err := suite.AuthClient.Register(ctx, &ssov1.RegisterRequest{
					Email:    testCase.email,
					Password: testCase.password,
				})

				require.Error(t, err)
				assert.ErrorContains(t, err, testCase.expectedError)
			})
		}
	})
}

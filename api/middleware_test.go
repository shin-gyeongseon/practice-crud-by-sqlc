package api

import (
	"fmt"
	"go-practice/token"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddlerWare(t *testing.T) {
	user, _ := randomUser(t)

	testCases := []struct {
		name string
		setupAuth func(*testing.T, *http.Request, token.Maker)
		checkResponse func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "ok",
			setupAuth: func(t *testing.T, r *http.Request, m token.Maker) {
				duration := time.Duration.Minutes(1)
				token, _, err := m.CreateToken(user.Username, time.Duration(duration))
				require.NoError(t, err)

				authorization := fmt.Sprintf("%s %s", authorizationTypeBearer, token)
				r.Header.Set(authorizationHeaderKey, authorization)
			},
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rr.Code)
			},
		},
		{
			name: "headerNotProvided",
			setupAuth: func(t *testing.T, r *http.Request, m token.Maker) {
				authorization := ""
				r.Header.Set(authorizationHeaderKey, authorization)
			},
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, rr.Code)
			},
		},
		{
			name: "invalidAuthorizationHeaderFormat",
			setupAuth: func(t *testing.T, r *http.Request, m token.Maker) {
				duration := time.Duration.Minutes(1)
				token, _, err := m.CreateToken(user.Username, time.Duration(duration))
				require.NoError(t, err)

				authorization := fmt.Sprintf("%s %s", "", token)
				r.Header.Set(authorizationHeaderKey, authorization)
			},
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, rr.Code)
			},
		},
		{
			name: "unsupportedAuthorizationType",
			setupAuth: func(t *testing.T, r *http.Request, m token.Maker) {
				duration := time.Duration.Minutes(1)
				token, _, err := m.CreateToken(user.Username, time.Duration(duration))
				require.NoError(t, err)

				authorization := fmt.Sprintf("%s %s", "zone", token)
				r.Header.Set(authorizationHeaderKey, authorization)
			},
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, rr.Code)
			},
		},
		{
			name: "invalidAuthorizationHeaderFormat",
			setupAuth: func(t *testing.T, r *http.Request, m token.Maker) {
				duration := time.Duration.Minutes(1)
				token, _, err := m.CreateToken(user.Username, time.Duration(duration))
				require.NoError(t, err)

				authorization := fmt.Sprintf("%s %s", authorizationTypeBearer, token + "failed_token")
				r.Header.Set(authorizationHeaderKey, authorization)
			},
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, rr.Code)
			},
		},
		
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			server := newTestServer(t, nil)
			authPath := "/auth"
			server.router.GET(
				authPath,
				authMiddleware(server.tokenMaker),
				func(ctx *gin.Context) {
					ctx.JSON(http.StatusOK, gin.H{})
				},
			)

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

package auth

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManager_GetTokenFromCookie(t *testing.T) {
	manager := New("test-secret-key")

	tests := []struct {
		name         string
		setupRequest func() *http.Request
		wantToken    string
		wantError    bool
		errorMsg     string
	}{
		{
			name: "Valid cookie with token",
			setupRequest: func() *http.Request {
				return NewMockRequestWithCookie("valid-token")
			},
			wantToken: "valid-token",
			wantError: false,
		},
		{
			name: "Empty cookie value",
			setupRequest: func() *http.Request {
				return NewMockRequestWithCookie("")
			},
			wantToken: "",
			wantError: true,
			errorMsg:  "unauthorized: cookie not found",
		},
		{
			name: "No cookie in request",
			setupRequest: func() *http.Request {
				return NewMockRequestWithoutCookie()
			},
			wantToken: "",
			wantError: true,
			errorMsg:  "unauthorized: cookie not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.setupRequest()

			token, err := manager.GetTokenFromCookie(req)

			if tt.wantError {
				require.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantToken, token)
			}
		})
	}
}

func TestManager_BuildJWTAndParseUserID(t *testing.T) {
	manager := New("test-secret-key")

	tests := []struct {
		name      string
		userID    int
		wantError bool
	}{
		{
			name:      "Valid user ID",
			userID:    123,
			wantError: false,
		},
		{
			name:      "Another valid user ID",
			userID:    456,
			wantError: false,
		},
		{
			name:      "Zero user ID",
			userID:    0,
			wantError: false,
		},
		{
			name:      "Negative user ID",
			userID:    -1,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenString, err := manager.BuildJWT(tt.userID)
			if tt.wantError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, tokenString)

			parsedUserID, err := manager.ParseUserID(tokenString)
			require.NoError(t, err)
			assert.Equal(t, tt.userID, parsedUserID)
		})
	}
}

func TestManager_ParseUserID_InvalidTokens(t *testing.T) {
	manager := New("test-secret-key")

	tests := []struct {
		name      string
		token     string
		wantError bool
		errorMsg  string
	}{
		{
			name:      "Empty token",
			token:     "",
			wantError: true,
			errorMsg:  "invalid token",
		},
		{
			name:      "Malformed token",
			token:     "not.a.valid.token",
			wantError: true,
			errorMsg:  "invalid token",
		},
		{
			name:      "Token with wrong signature",
			token:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			wantError: true,
			errorMsg:  "invalid token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := manager.ParseUserID(tt.token)

			require.Error(t, err)
			if tt.errorMsg != "" {
				assert.Contains(t, err.Error(), tt.errorMsg)
			}
			assert.Equal(t, -1, userID)
		})
	}
}

func TestManager_CreateToken(t *testing.T) {
	manager := New("test-secret-key")

	tests := []struct {
		name      string
		userID    int
		wantError bool
	}{
		{
			name:      "Create token for valid user",
			userID:    123,
			wantError: false,
		},
		{
			name:      "Create token for another user",
			userID:    456,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWriter := NewMockResponseWriter()

			tokenString, err := manager.CreateToken(mockWriter, tt.userID)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, tokenString)

			cookies := mockWriter.Header().Get("Set-Cookie")
			assert.Contains(t, cookies, "ID="+tokenString)
			assert.Contains(t, cookies, "HttpOnly")
			assert.Contains(t, cookies, "Path=/")

			parsedUserID, err := manager.ParseUserID(tokenString)
			require.NoError(t, err)
			assert.Equal(t, tt.userID, parsedUserID)
		})
	}
}

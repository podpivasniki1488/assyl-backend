package service

import (
	"context"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/podpivasniki1488/assyl-backend/protopb"
	"golang.org/x/crypto/bcrypt"
)

// Mock tracer for testing
type mockTracer struct{}

func (m *mockTracer) Start(ctx context.Context, spanName string, opts ...interface{}) (context.Context, mockSpan) {
	return ctx, mockSpan{}
}

type mockSpan struct{}

func (m mockSpan) End(opts ...interface{}) {}

func TestAuthService_isEmail(t *testing.T) {
	// Create a minimal auth service for testing pure functions
	as := &authService{}
	
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "valid email",
			input: "test@example.com",
			want:  true,
		},
		{
			name:  "valid email with subdomain",
			input: "user@mail.example.com",
			want:  true,
		},
		{
			name:  "valid email with plus",
			input: "user+tag@example.com",
			want:  true,
		},
		{
			name:  "invalid email - no @",
			input: "testexample.com",
			want:  false,
		},
		{
			name:  "invalid email - no domain",
			input: "test@",
			want:  false,
		},
		{
			name:  "invalid email - no username",
			input: "@example.com",
			want:  false,
		},
		{
			name:  "empty string",
			input: "",
			want:  false,
		},
		{
			name:  "phone number",
			input: "+1234567890",
			want:  false,
		},
		{
			name:  "invalid email - multiple @",
			input: "test@@example.com",
			want:  false,
		},
		{
			name:  "invalid email - spaces",
			input: "test @example.com",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := as.isEmail(tt.input)
			if got != tt.want {
				t.Errorf("isEmail(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestAuthService_isPhone(t *testing.T) {
	as := &authService{}
	
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "valid phone starting with 13",
			input: "13912345678",
			want:  true,
		},
		{
			name:  "valid phone starting with 14",
			input: "14912345678",
			want:  true,
		},
		{
			name:  "valid phone starting with 15",
			input: "15912345678",
			want:  true,
		},
		{
			name:  "valid phone starting with 17",
			input: "17912345678",
			want:  true,
		},
		{
			name:  "valid phone starting with 18",
			input: "18912345678",
			want:  true,
		},
		{
			name:  "invalid phone - too short",
			input: "139123456",
			want:  false,
		},
		{
			name:  "invalid phone - too long",
			input: "139123456789",
			want:  false,
		},
		{
			name:  "invalid phone - starts with 12",
			input: "12912345678",
			want:  false,
		},
		{
			name:  "invalid phone - starts with 19",
			input: "19912345678",
			want:  false,
		},
		{
			name:  "email address",
			input: "test@example.com",
			want:  false,
		},
		{
			name:  "empty string",
			input: "",
			want:  false,
		},
		{
			name:  "letters mixed with numbers",
			input: "13a12345678",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := as.isPhone(tt.input)
			if got != tt.want {
				t.Errorf("isPhone(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestAuthService_comparePasswords(t *testing.T) {
	as := &authService{}
	
	// Generate a test hash
	testPassword := "TestPassword123!"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to generate test hash: %v", err)
	}
	
	tests := []struct {
		name     string
		hashedPw string
		plainPw  string
		want     bool
	}{
		{
			name:     "correct password",
			hashedPw: string(hashedPassword),
			plainPw:  testPassword,
			want:     true,
		},
		{
			name:     "incorrect password",
			hashedPw: string(hashedPassword),
			plainPw:  "WrongPassword",
			want:     false,
		},
		{
			name:     "empty plain password",
			hashedPw: string(hashedPassword),
			plainPw:  "",
			want:     false,
		},
		{
			name:     "case sensitive password check",
			hashedPw: string(hashedPassword),
			plainPw:  "testpassword123!",
			want:     false,
		},
		{
			name:     "password with extra characters",
			hashedPw: string(hashedPassword),
			plainPw:  testPassword + "extra",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := as.comparePasswords(tt.hashedPw, tt.plainPw)
			if got != tt.want {
				t.Errorf("comparePasswords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthService_generateJwtToken(t *testing.T) {
	secretKey := "test-secret-key-123"
	as := &authService{
		secretKey: secretKey,
	}
	
	tests := []struct {
		name     string
		username string
		role     protopb.Role
		wantErr  bool
	}{
		{
			name:     "generate token for guest user",
			username: "testuser",
			role:     protopb.Role_GUEST,
			wantErr:  false,
		},
		{
			name:     "generate token for admin user",
			username: "admin",
			role:     protopb.Role_ADMIN,
			wantErr:  false,
		},
		{
			name:     "generate token with empty username",
			username: "",
			role:     protopb.Role_GUEST,
			wantErr:  false, // JWT allows empty claims
		},
		{
			name:     "generate token with special characters in username",
			username: "user@example.com",
			role:     protopb.Role_GUEST,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := as.generateJwtToken(tt.username, tt.role)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("generateJwtToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr {
				if token == "" {
					t.Error("generateJwtToken() returned empty token")
				}
				
				// Verify the token can be parsed
				parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
					return []byte(secretKey), nil
				})
				
				if err != nil {
					t.Errorf("Failed to parse generated token: %v", err)
				}
				
				if !parsedToken.Valid {
					t.Error("Generated token is not valid")
				}
				
				// Verify claims
				claims, ok := parsedToken.Claims.(jwt.MapClaims)
				if !ok {
					t.Fatal("Failed to extract claims from token")
				}
				
				if claims["username"] != tt.username {
					t.Errorf("Token username = %v, want %v", claims["username"], tt.username)
				}
				
				if claims["role"] != tt.role.String() {
					t.Errorf("Token role = %v, want %v", claims["role"], tt.role.String())
				}
				
				if claims["issuer"] != "jeffry's backend" {
					t.Errorf("Token issuer = %v, want 'jeffry's backend'", claims["issuer"])
				}
			}
		})
	}
}

func TestAuthService_generateJwtToken_SignatureVerification(t *testing.T) {
	secretKey := "test-secret-key"
	wrongSecretKey := "wrong-secret-key"
	
	as := &authService{
		secretKey: secretKey,
	}
	
	token, err := as.generateJwtToken("testuser", protopb.Role_GUEST)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	
	// Try to parse with wrong secret - should fail
	_, err = jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(wrongSecretKey), nil
	})
	
	if err == nil {
		t.Error("Token parsed successfully with wrong secret key, should have failed")
	}
	
	// Parse with correct secret - should succeed
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	
	if err != nil {
		t.Errorf("Failed to parse token with correct secret: %v", err)
	}
	
	if !parsedToken.Valid {
		t.Error("Token is not valid with correct secret")
	}
}

func TestUsernameTypeConstants(t *testing.T) {
	// Test that username type constants have expected values
	tests := []struct {
		name  string
		value int
		want  int
	}{
		{
			name:  "UsernameTypeNone",
			value: UsernameTypeNone,
			want:  1,
		},
		{
			name:  "UsernameTypeEmail",
			value: UsernameTypeEmail,
			want:  2,
		},
		{
			name:  "UsernameTypePhone",
			value: UsernameTypePhone,
			want:  3,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != tt.want {
				t.Errorf("%s = %d, want %d", tt.name, tt.value, tt.want)
			}
		})
	}
	
	// Test that values are unique
	values := map[int]bool{
		UsernameTypeNone:  true,
		UsernameTypeEmail: true,
		UsernameTypePhone: true,
	}
	
	if len(values) != 3 {
		t.Error("Username type constants are not unique")
	}
}

func TestPasswordHashing_ConsistencyCheck(t *testing.T) {
	// Test that bcrypt hashing is consistent with comparison
	password := "TestPassword123"
	
	hash1, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to generate hash: %v", err)
	}
	
	hash2, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to generate second hash: %v", err)
	}
	
	// Hashes should be different (bcrypt uses salt)
	if string(hash1) == string(hash2) {
		t.Error("Two hashes of the same password should be different due to salt")
	}
	
	as := &authService{}
	
	// But both should verify successfully
	if !as.comparePasswords(string(hash1), password) {
		t.Error("First hash did not verify")
	}
	
	if !as.comparePasswords(string(hash2), password) {
		t.Error("Second hash did not verify")
	}
}

func TestAuthService_UserRegistrationFlow(t *testing.T) {
	// Integration-style test for username type detection
	as := &authService{}
	
	testCases := []struct {
		username     string
		expectedType int
	}{
		{
			username:     "user@example.com",
			expectedType: UsernameTypeEmail,
		},
		{
			username:     "13912345678",
			expectedType: UsernameTypePhone,
		},
		{
			username:     "randomuser123",
			expectedType: UsernameTypeNone,
		},
	}
	
	for _, tc := range testCases {
		t.Run("username_"+tc.username, func(t *testing.T) {
			var actualType int
			
			switch {
			case as.isEmail(tc.username):
				actualType = UsernameTypeEmail
			case as.isPhone(tc.username):
				actualType = UsernameTypePhone
			default:
				actualType = UsernameTypeNone
			}
			
			if actualType != tc.expectedType {
				t.Errorf("Username type for %s = %d, want %d", tc.username, actualType, tc.expectedType)
			}
		})
	}
}
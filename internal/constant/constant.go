package constant

// Error codes for authentication
const (
	InvalidCredentials    string = "invalid_credentials"
	InvalidToken          string = "invalid_token"
	InvalidSigningMethod  string = "invalid_signing_method"
	MissingToken          string = "missing_token"
	InvalidTokenSignature string = "invalid_token_signature"
	ExpiredToken          string = "expired_token"
)

// Success codes for authentication
const (
	LoginSuccess             string = "login_success"
	LogoutSuccess            string = "logout_success"
	RefreshSuccess           string = "refresh_success"
	ForgotPasswordOTPSuccess string = "forgot_password_otp_success"
	ResetPasswordOTPSuccess  string = "reset_password_otp_success"
)

// Gin context keys
const (
	ContextResponseKey        string = "x-response"
	ContextUserIDKey          string = "x-user-id"
	ContextUserRolesKey       string = "x-user-roles"
	ContextUserPermissionsKey string = "x-user-permissions"
	ContextAccessTokenKey     string = "x-access-token"
	ContextDeviceIDKey        string = "x-device-id"
)

// Prefix for redis keys
const (
	CacheKeyPrefix            string = "cache"
	ResetPasswordOTPKeyPrefix string = "reset_password_otp"
	RefreshTokenKeyPrefix     string = "refresh_token"
)

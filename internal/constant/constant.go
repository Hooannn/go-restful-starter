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
	LoginSuccess   string = "login_success"
	LogoutSuccess  string = "logout_success"
	RefreshSuccess string = "refresh_success"
)

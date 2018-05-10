package adapter

type AuthProvider interface {
	RequireAuth() bool
	Verify(bearerToken string) error
}

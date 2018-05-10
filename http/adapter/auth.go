package adapter

type AuthProvider interface {
	RequireAuth() bool
	Verify(bearerToken string) error
}

func RequireAuth() bool {
	return false
}

func Verify(bearerToken string) error {
	return nil
}

package authentication

type RefreshTokenGenerator interface {
	Generate() (string, error)
}

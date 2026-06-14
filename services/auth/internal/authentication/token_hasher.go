package authentication

type TokenHasher interface {
	Hash(
		token string,
	) string
}

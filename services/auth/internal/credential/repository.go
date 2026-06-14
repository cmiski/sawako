package credential

import "context"

type Repository interface {
	Create(
		ctx context.Context,
		credential *Credential,
	) error
}

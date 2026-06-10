package uuidx

import (
	"fmt"

	"github.com/google/uuid"
)

func NewV7() string {
	id, err := uuid.NewV7()
	if err != nil {
		panic(
			fmt.Sprintf(
				"uuidx: failed to generate UUIDv7: %v", err,
			),
		)
	}
	return id.String()
}

func IsValid(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}

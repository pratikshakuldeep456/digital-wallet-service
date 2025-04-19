package dws

import "github.com/google/uuid"

func GenerateId() int {

	return int(uuid.New().ID())
}

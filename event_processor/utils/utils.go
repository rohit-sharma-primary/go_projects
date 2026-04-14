package utils

import (
	"event_processor/model"
	"math/rand/v2"

	"github.com/google/uuid"
)

func GenerateEventId() string {
	return uuid.NewString()
}

func GenerateUserId() int {
	return rand.IntN(100) + 1
}

func GenerateEventType() model.EventType {
	switch rand.IntN(3) {
	case 0:
		return model.Click
	case 1:
		return model.View
	default:
		return model.Purchase
	}
}

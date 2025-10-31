package kafka

import "github.com/major1ink/simple-notification-telegram/internal/model"

type OrderAssembledDecoder interface {
	DecodeAssembled(data []byte) (model.AssembledEvent, error)
}

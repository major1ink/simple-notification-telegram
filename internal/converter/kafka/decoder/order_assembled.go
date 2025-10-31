package decoder

import (
	"encoding/json"
	"fmt"

	"github.com/major1ink/simple-notification-telegram/internal/model"
)

type decoderAssembled struct{}

func NewOrderDecoderAssembled() *decoderAssembled {
	return &decoderAssembled{}
}

func (d *decoderAssembled) DecodeAssembled(data []byte) (model.AssembledEvent, error) {
	var event model.AssembledEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return model.AssembledEvent{}, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return event, nil
}

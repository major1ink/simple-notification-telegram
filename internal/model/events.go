package model

type AssembledEvent struct {
	EventUuid string `json:"event_uuid"`
	TypeEvent string `json:"type_event"`
	App       string `json:"app"`
	Message   string `json:"message"`
}

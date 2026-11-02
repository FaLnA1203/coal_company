package api

import (
	"encoding/json"
	"time"
)

type ErrorDTO struct {
	Error string
	Time  time.Time
}

func NewErrorDTO(message error) ErrorDTO {
	return ErrorDTO{
		Error: message.Error(),
		Time:  time.Now(),
	}
}

func (e ErrorDTO) ToString() string {
	b, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		panic(err)
	}

	return string(b)
}

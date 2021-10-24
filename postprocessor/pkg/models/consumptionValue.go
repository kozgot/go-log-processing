package models

import (
	"encoding/json"
	"fmt"
	"time"
)

// ConsumtionValue contains cunsumption data in a given time range.
type ConsumtionValue struct {
	ReceiveTime  time.Time
	StartTime    time.Time
	EndTime      time.Time
	Value        int
	ServiceLevel int
	SmcUID       string
}

// Serialize serlializes a consumption value to JSON format and returns a byte array.
func (c ConsumtionValue) Serialize() []byte {
	bytes, err := json.Marshal(c)
	if err != nil {
		fmt.Println("Can't serialize consumption ", c)
	}

	return bytes
}

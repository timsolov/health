package health

import "encoding/json"

type Status string

const (
	Up           Status = "UP"
	Down         Status = "DOWN"
	OutOfService Status = "OUT_OF_SERVICE"
	NotReady     Status = "NOT_READY"
	Unknown      Status = "UNKNOWN"
)

// Health is a health status struct
type Health struct {
	status Status
	info   map[string]interface{}
}

// MarshalJSON is a custom JSON marshaller
func (h Health) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}{}

	for k, v := range h.info {
		data[k] = v
	}

	data["status"] = h.status

	return json.Marshal(data)
}

// NewHealth return a new Health with status Down
func NewHealth() Health {
	h := Health{
		info: make(map[string]interface{}),
	}

	h.Unknown()

	return h
}

// AddInfo adds a info value to the Info map
func (h *Health) AddInfo(key string, value interface{}) *Health {
	if h.info == nil {
		h.info = make(map[string]interface{})
	}

	h.info[key] = value

	return h
}

// GetInfo returns a value from the info map
func (h Health) GetInfo(key string) interface{} {
	return h.info[key]
}

// IsUnknown returns true if Status is Unknown
func (h Health) IsUnknown() bool {
	return h.status == Unknown
}

// IsUp returns true if Status is Up
func (h Health) IsUp() bool {
	return h.status == Up
}

// IsDown returns true if Status is Down
func (h Health) IsDown() bool {
	return h.status == Down
}

// IsOutOfService returns true if Status is IsOutOfService
func (h Health) IsOutOfService() bool {
	return h.status == OutOfService
}

// IsNotReady returns true if Status is IsNotReady
func (h Health) IsNotReady() bool {
	return h.status == NotReady
}

// Down set the status to Down
func (h *Health) Down() *Health {
	h.status = Down
	return h
}

// OutOfService set the status to OutOfService
func (h *Health) OutOfService() *Health {
	h.status = OutOfService
	return h
}

// NotReady set the status to NotReady
func (h *Health) NotReady() *Health {
	h.status = NotReady
	return h
}

// Unknown set the status to Unknown
func (h *Health) Unknown() *Health {
	h.status = Unknown
	return h
}

// Up set the status to Up
func (h *Health) Up() *Health {
	h.status = Up
	return h
}

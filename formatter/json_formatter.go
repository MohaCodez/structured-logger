package formatter

import (
	"encoding/json"

	"github.com/MohaCodez/structured-logger/logger"
)

type JSONFormatter struct{}

func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{}
}

func (f *JSONFormatter) Format(entry *logger.Entry) ([]byte, error) {
	m := map[string]interface{}{
		"timestamp": entry.Timestamp,
		"level":     entry.Level,
		"message":   entry.Message,
	}
	if entry.Caller != "" {
		m["caller"] = entry.Caller
	}
	for k, v := range entry.Fields {
		m[k] = v
	}
	return json.Marshal(m)
}

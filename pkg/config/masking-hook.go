package config

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

type MaskingHook struct {
	toMask      []string
	MaskedValue string
}

// Adds the specified secrets to the logging secret mask so that it's not emitted in the output
func (h *MaskingHook) AddToMaskList(secret string) {
	if len(strings.TrimSpace(secret)) > 0 {
		h.toMask = append(h.toMask, secret)
	}
}

func (h *MaskingHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *MaskingHook) Fire(entry *logrus.Entry) error {
	for _, val := range h.toMask {
		entry.Message = h.replace(entry.Message, val)

		for k, dataVal := range entry.Data {
			entry.Data[k] = h.replace(fmt.Sprint(dataVal), val)
		}
	}
	return nil
}

func (h *MaskingHook) replace(message string, secret string) string {
	return strings.ReplaceAll(message, secret, h.MaskedValue)
}

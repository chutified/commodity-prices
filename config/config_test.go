package config

import (
	"fmt"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestGetConfig(t *testing.T) {

	tests := []struct {
		name   string
		file   string
		errMsg string
	}{
		{
			name:   "ok",
			file:   "config.yaml",
			errMsg: "",
		},
		{
			name:   "file not found",
			file:   "config_not_exist.yaml",
			errMsg: "could not read config file",
		},
		{
			name:   "invalid file",
			file:   "config/tests/config_invalid.yaml",
			errMsg: "could not decode config file content",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t1 *testing.T) {

			cfg, err := GetConfig(test.file)
			if err != nil {

				exp := fmt.Sprintf(".*%s.*", test.errMsg)
				assert.MatchRegex(t1, err.Error(), exp)

			} else {

				assert.NotEqual(t1, cfg.Host, "")
				assert.NotEqual(t1, cfg.Port, 0)
			}
		})
	}
}

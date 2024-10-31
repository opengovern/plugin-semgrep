package provider

import (
	"encoding/json"
)


// Import the any type from the config package.
// You should write you own Authorization Function.



type AccountConfig struct {
	// You should provide Credentials for any Provider.
}

// AccountConfigFromMap converts a map to an AccountConfig.
func AccountConfigFromMap(m map[string]any) (AccountConfig, error) {
	mj, err := json.Marshal(m)
	if err != nil {
		return AccountConfig{}, err
	}

	var c AccountConfig
	err = json.Unmarshal(mj, &c)
	if err != nil {
		return AccountConfig{}, err
	}

	return c, nil
}

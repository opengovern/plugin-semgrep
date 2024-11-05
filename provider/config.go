package provider

import (
	"encoding/json"
	model "github.com/opengovern/og-describer-template/pkg/SDK/models"
	"github.com/opengovern/og-util/pkg/describe"
)

// Import the any type from the config package.
// You should write you own Authorization Function.

type AccountConfig struct {
	// You should provide Credentials for any Provider.
}

// AccountConfigFromMap TODO: converts a map to an AccountConfig.
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

// GetResourceMetadata TODO: Get metadata as a map to add to the resources
func GetResourceMetadata(job describe.DescribeJob, resource model.Resource) (map[string]string, error) {
	metadata := make(map[string]string)

	return metadata, nil
}

// AdjustResource TODO: Do any needed adjustment on resource object before storing
func AdjustResource(job describe.DescribeJob, resource *model.Resource) error {
	return nil
}

// GetAdditionalParameters TODO: pass additional parameters needed in describer wrappers in /provider/describer_wrapper.go
func GetAdditionalParameters(job describe.DescribeJob) (map[string]string, error) {
	additionalParameters := make(map[string]string)

	return additionalParameters, nil
}

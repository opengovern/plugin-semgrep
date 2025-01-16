package main

import (
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/opengovern/og-describer-template/platform/constants"
	"github.com/opengovern/og-util/pkg/integration/interfaces"
	"os"
)

func main() {
	i := Integration{}
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Debug,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	var pluginMap = map[string]plugin.Plugin{
		constants.IntegrationName.String(): &interfaces.IntegrationTypePlugin{Impl: &i},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: interfaces.HandshakeConfig,
		Plugins:         pluginMap,
		Logger:          logger,
	})
}

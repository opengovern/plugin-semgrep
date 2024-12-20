package main

import (
	"github.com/opengovern/og-describer-template/plugin/cohereai"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: cohere.Plugin})
}

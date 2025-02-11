package main

import (
	semgrep "github.com/opengovern/og-describer-semgrep/cloudql/semgrep"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: semgrep.Plugin})
}

package adapter

import (
	pb "github.com/katallaxie/autobot/pkg/proto/v1"

	p "github.com/hashicorp/go-plugin"
)

const (
	PluginName             = "plugin"
	DefaultProtocolVersion = 1
)

var VersionedPlugins = map[int]p.PluginSet{
	1: {
		"plugin": &GRPCAdapter{},
	},
}

// Handshake ...
var Handshake = p.HandshakeConfig{
	ProtocolVersion: DefaultProtocolVersion,

	MagicCookieKey:   "AUTOBOT_ADAPTER_MAGIC_COOKIE",
	MagicCookieValue: "iaafij5485d5utqh",
}

// GRPCAdapterFunc ...
type GRPCAdapterFunc func() pb.AdapterServer

// ServeOpts ...
type ServeOpts struct {
	GRPCAdapterFunc GRPCAdapterFunc
}

// Serve ...
func Serve(opts *ServeOpts) {
	p.Serve(&p.ServeConfig{
		GRPCServer:       p.DefaultGRPCServer,
		HandshakeConfig:  Handshake,
		VersionedPlugins: pluginSet(opts),
	})
}

func pluginSet(opts *ServeOpts) map[int]p.PluginSet {
	plugins := map[int]p.PluginSet{}

	// add the new protocol versions if they're configured
	if opts.GRPCAdapterFunc != nil {
		plugins[1] = p.PluginSet{
			"plugin": &GRPCAdapter{
				GRPCPlugin: opts.GRPCAdapterFunc,
			},
		}
	}

	return plugins
}

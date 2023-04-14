package adapter

import (
	"context"
	"os"
	"os/exec"

	"github.com/katallaxie/autobot/internal/ports"
	pb "github.com/katallaxie/autobot/pkg/proto/v1"

	"github.com/hashicorp/go-hclog"
	p "github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

var enablePluginAutoMTLS = os.Getenv("RUN_DISABLE_PLUGIN_TLS") == ""

// Meta ...
type Meta struct {
	// Path ...
	Path string

	// Arguments ...
	Arguments []string
}

// ExecutableFile ...
func (m *Meta) ExecutableFile() (string, error) {
	// TODO: make this check for the executable file
	return m.Path, nil
}

func (m *Meta) Factory(ctx context.Context) Factory {
	return adapterFactory(ctx, m)
}

// GRPCAdapter ...
type GRPCAdapter struct {
	p.Plugin
	GRPCPlugin func() pb.AdapterServer
}

// GRPCClient ...
func (p *GRPCAdapter) GRPCClient(ctx context.Context, broker *p.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCPlugin{
		client: pb.NewAdapterClient(c),
		broker: broker,
		ctx:    ctx,
	}, nil
}

func (p *GRPCAdapter) GRPCServer(broker *p.GRPCBroker, s *grpc.Server) error {
	pb.RegisterAdapterServer(s, p.GRPCPlugin())

	return nil
}

// GRPCPlugin ...
type GRPCPlugin struct {
	PluginClient *p.Client

	ctx    context.Context
	client pb.AdapterClient
	broker *p.GRPCBroker
}

// Start ...
func (p *GRPCPlugin) Close() error {
	if p.PluginClient != nil {
		return nil
	}

	p.PluginClient.Kill()

	return nil
}

// Subscribe ...
func (p *GRPCPlugin) Subscribe(req ports.SubscribeRequest, queue chan<- *pb.Event) (ports.SubscribeResponse, error) {
	r := new(pb.Subscribe_Request)

	c, err := p.client.Subscribe(p.ctx, r)
	if err != nil {
		return ports.SubscribeResponse{}, err
	}

	go func() {
		for {
			resp, err := c.Recv()
			if err != nil {
				close(queue)
			}

			queue <- resp.GetEvent()
		}
	}()

	return ports.SubscribeResponse{}, nil
}

// Publish ...
func (p *GRPCPlugin) Publish(req ports.PublishRequest) (ports.PublishResponse, error) {
	r := new(pb.Publish_Request)
	r.Reply = req.Reply

	_, err := p.client.Publish(p.ctx, r)
	if err != nil {
		return ports.PublishResponse{}, err
	}

	return ports.PublishResponse{}, nil
}

// Factory ...
type Factory func() (ports.MessageStore, error)

func adapterFactory(ctx context.Context, meta *Meta) Factory {
	return func() (ports.MessageStore, error) {
		f, err := meta.ExecutableFile()
		if err != nil {
			return nil, err
		}

		l := hclog.New(&hclog.LoggerOptions{
			Name:  meta.Path,
			Level: hclog.LevelFromString("DEBUG"),
		})

		cfg := &p.ClientConfig{
			AllowedProtocols: []p.Protocol{p.ProtocolGRPC},
			AutoMTLS:         enablePluginAutoMTLS,
			Cmd:              exec.CommandContext(ctx, f, meta.Arguments...), //nolint:gosec
			HandshakeConfig:  Handshake,
			Logger:           l,
			Managed:          true,
			SyncStderr:       l.StandardWriter(&hclog.StandardLoggerOptions{}),
			SyncStdout:       l.StandardWriter(&hclog.StandardLoggerOptions{}),
			VersionedPlugins: VersionedPlugins,
		}
		client := p.NewClient(cfg)

		rpc, err := client.Client()
		if err != nil {
			return nil, err
		}

		raw, err := rpc.Dispense(PluginName)
		if err != nil {
			return nil, err
		}

		p := raw.(*GRPCPlugin)
		p.PluginClient = client

		return p, nil
	}
}

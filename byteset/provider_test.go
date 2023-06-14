package byteset

import (
	"context"
	"io"
	"strings"
	"testing"
	"text/template"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/stretchr/testify/assert"

	"github.com/seal-io/terraform-provider-byteset/utils/version"
)

func TestProvider_metadata(t *testing.T) {
	var (
		ctx  = context.TODO()
		req  = provider.MetadataRequest{}
		resp = &provider.MetadataResponse{}
	)

	p := NewProvider()
	p.Metadata(ctx, req, resp)
	assert.Equal(t, resp.TypeName, ProviderType)
	assert.Equal(t, resp.Version, version.Version)
}

var testAccProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"byteset": providerserver.NewProtocol6WithError(NewProvider()),
}

func renderConfigTemplate(ct string, keyValuePairs ...any) string {
	t, err := template.New("ct").Parse(ct)
	if err != nil {
		panic(err)
	}

	d := make(map[any]any, len(keyValuePairs)/2)
	for i := 0; i < len(keyValuePairs); i += 2 {
		d[keyValuePairs[i]] = keyValuePairs[i+1]
	}

	var s strings.Builder

	err = t.Execute(&s, d)
	if err != nil {
		panic(err)
	}

	return s.String()
}

type dockerContainer struct {
	Name  string
	Image string
	Env   []string
	Port  []string

	cli *client.Client
	id  string
}

func (c *dockerContainer) Start(t *testing.T, ctx context.Context) error {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return err
	}
	c.cli = cli

	t.Logf("Pulling container %s image: %s", c.Name, c.Image)

	pulled, err := cli.ImagePull(ctx, c.Image,
		types.ImagePullOptions{})
	if err != nil {
		return err
	}
	_, _ = io.Copy(io.Discard, pulled)
	_ = pulled.Close()

	var (
		cCfg = &container.Config{
			Image:        c.Image,
			Env:          c.Env,
			ExposedPorts: nat.PortSet{},
		}
		chCfg = &container.HostConfig{
			PortBindings: nat.PortMap{},
		}
		nCfg = &network.NetworkingConfig{}
	)

	for _, p := range c.Port {
		var host, external, internal string

		ps := strings.Split(p, ":")
		switch len(ps) {
		case 3:
			host, external, internal = ps[0], ps[1], ps[2]
		case 2:
			external, internal = ps[0], ps[1]
		case 1:
			internal = ps[0]
		}

		if internal == "" {
			continue
		}

		cCfg.ExposedPorts[nat.Port(internal)] = struct{}{}

		if external == "" {
			continue
		}

		chCfg.PortBindings[nat.Port(internal)] = append(chCfg.PortBindings[nat.Port(internal)],
			nat.PortBinding{
				HostPort: external,
				HostIP:   host,
			})
	}

	t.Logf("Creating container %s", c.Name)

	created, err := cli.ContainerCreate(ctx, cCfg, chCfg, nCfg, nil, c.Name)
	if err != nil {
		return err
	}
	c.id = created.ID

	t.Logf("Starting container %s", c.Name)

	err = cli.ContainerStart(ctx, c.id,
		types.ContainerStartOptions{})
	if err != nil {
		_ = c.Stop(t, ctx)
		return err
	}

	return nil
}

func (c *dockerContainer) Stop(t *testing.T, ctx context.Context) error {
	if c.id == "" || c.cli == nil {
		return nil
	}

	t.Logf("Stopping container %s", c.Name)
	err := c.cli.ContainerRemove(ctx, c.id,
		types.ContainerRemoveOptions{
			Force: true,
		})
	_ = c.cli.Close()

	return err
}

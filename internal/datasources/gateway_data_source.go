package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/virtualmo/terraform-provider-netskopesdwan/internal/client"
)

var (
	_ datasource.DataSource              = &gatewayDataSource{}
	_ datasource.DataSourceWithConfigure = &gatewayDataSource{}
)

type gatewayDataSource struct {
	client *client.Client
}

func NewGatewayDataSource() datasource.DataSource {
	return &gatewayDataSource{}
}

func (d *gatewayDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gateway"
}

func (d *gatewayDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasourceschema.Schema{
		Description: "Gateway data source placeholder.",
	}
}

func (d *gatewayDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	apiClient, ok := req.ProviderData.(*client.Client)
	if !ok {
		return
	}

	d.client = apiClient
}

func (d *gatewayDataSource) Read(_ context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured provider client",
			"The provider client was not configured before reading netskopesdwan_gateway.",
		)
		return
	}

	resp.Diagnostics.AddError(
		"Data source not implemented",
		"TODO: The Netskope SD-WAN gateway API response shape and lookup arguments are not confirmed yet, so this data source intentionally avoids inventing Terraform schema fields.",
	)
}

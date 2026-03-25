package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/virtualmo/terraform-provider-netskopesdwan/internal/client"
)

var (
	_ datasource.DataSource              = &gatewaysDataSource{}
	_ datasource.DataSourceWithConfigure = &gatewaysDataSource{}
)

type gatewaysDataSource struct {
	client *client.Client
}

func NewGatewaysDataSource() datasource.DataSource {
	return &gatewaysDataSource{}
}

func (d *gatewaysDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gateways"
}

func (d *gatewaysDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasourceschema.Schema{
		Description: "Gateways data source placeholder.",
	}
}

func (d *gatewaysDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	apiClient, ok := req.ProviderData.(*client.Client)
	if !ok {
		return
	}

	d.client = apiClient
}

func (d *gatewaysDataSource) Read(_ context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured provider client",
			"The provider client was not configured before reading netskopesdwan_gateways.",
		)
		return
	}

	resp.Diagnostics.AddError(
		"Data source not implemented",
		"TODO: The Netskope SD-WAN gateways API response shape is not confirmed yet, so this data source intentionally avoids inventing Terraform schema fields.",
	)
}

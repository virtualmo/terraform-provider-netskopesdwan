package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/virtualmo/terraform-provider-netskopesdwan/internal/client"
)

var (
	_ datasource.DataSource              = &gatewaysDataSource{}
	_ datasource.DataSourceWithConfigure = &gatewaysDataSource{}
)

const gatewaysCollectionID = "gateways"

type gatewaysDataSource struct {
	client *client.Client
}

type gatewaysDataSourceModel struct {
	ID types.String `tfsdk:"id"`
}

func NewGatewaysDataSource() datasource.DataSource {
	return &gatewaysDataSource{}
}

func (d *gatewaysDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gateways"
}

func (d *gatewaysDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasourceschema.Schema{
		Description: "Reads the Netskope SD-WAN gateways collection.",
		Attributes: map[string]datasourceschema.Attribute{
			"id": datasourceschema.StringAttribute{
				Computed:    true,
				Description: "Synthetic identifier for the gateways collection.",
			},
		},
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

func (d *gatewaysDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured provider client",
			"The provider client was not configured before reading netskopesdwan_gateways.",
		)
		return
	}

	_, err := d.client.ListGateways(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read GET /v2/gateways",
			"The Netskope SD-WAN gateways list request failed: "+err.Error(),
		)
		return
	}

	// TODO: Add Terraform-facing gateway attributes only after the real response
	// item shape is confirmed. The client already isolates envelope uncertainty.
	state := gatewaysDataSourceModel{
		ID: types.StringValue(gatewaysCollectionID),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

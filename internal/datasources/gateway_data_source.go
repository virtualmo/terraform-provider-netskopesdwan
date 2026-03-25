package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/virtualmo/terraform-provider-netskopesdwan/internal/client"
)

var (
	_ datasource.DataSource              = &gatewayDataSource{}
	_ datasource.DataSourceWithConfigure = &gatewayDataSource{}
)

type gatewayDataSource struct {
	client *client.Client
}

type gatewayDataSourceModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	CreatedAt            types.String `tfsdk:"created_at"`
	ModifiedAt           types.String `tfsdk:"modified_at"`
	ConfigUpdatesEnabled types.Bool   `tfsdk:"config_updates_enabled"`
	Managed              types.Bool   `tfsdk:"managed"`
}

func NewGatewayDataSource() datasource.DataSource {
	return &gatewayDataSource{}
}

func (d *gatewayDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gateway"
}

func (d *gatewayDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasourceschema.Schema{
		Description: "Reads a single Netskope SD-WAN gateway by ID.",
		Attributes: map[string]datasourceschema.Attribute{
			"id": datasourceschema.StringAttribute{
				Required:    true,
				Description: "Gateway identifier.",
			},
			"name": datasourceschema.StringAttribute{
				Computed:    true,
				Description: "Gateway name.",
			},
			"created_at": datasourceschema.StringAttribute{
				Computed:    true,
				Description: "Gateway creation timestamp returned by the API.",
			},
			"modified_at": datasourceschema.StringAttribute{
				Computed:    true,
				Description: "Gateway modification timestamp returned by the API.",
			},
			"config_updates_enabled": datasourceschema.BoolAttribute{
				Computed:    true,
				Description: "Whether config updates are enabled for the gateway.",
			},
			"managed": datasourceschema.BoolAttribute{
				Computed:    true,
				Description: "Whether the gateway is managed.",
			},
		},
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

func (d *gatewayDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured provider client",
			"The provider client was not configured before reading netskopesdwan_gateway.",
		)
		return
	}

	var config gatewayDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.ID.IsUnknown() {
		resp.Diagnostics.AddError(
			"Unknown gateway id",
			"The id attribute must be known before reading netskopesdwan_gateway.",
		)
		return
	}

	gatewayID := config.ID.ValueString()
	if gatewayID == "" {
		resp.Diagnostics.AddError(
			"Missing gateway id",
			"The id attribute is required to read netskopesdwan_gateway.",
		)
		return
	}

	gateway, err := d.client.GetGateway(ctx, gatewayID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read GET /v2/gateways/{id}",
			"The Netskope SD-WAN gateway request failed: "+err.Error(),
		)
		return
	}

	state := gatewayDataSourceModel{
		ID:                   types.StringValue(gateway.ID),
		Name:                 types.StringValue(gateway.Name),
		CreatedAt:            types.StringValue(gateway.CreatedAt),
		ModifiedAt:           types.StringValue(gateway.ModifiedAt),
		ConfigUpdatesEnabled: types.BoolValue(gateway.ConfigUpdatesEnabled),
		Managed:              types.BoolValue(gateway.Managed),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

package datasources

import (
	"context"
	"slices"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func TestGatewayDataSourceSchemaIsConservative(t *testing.T) {
	t.Parallel()

	ds := NewGatewayDataSource()
	resp := &datasource.SchemaResponse{}

	ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)

	gotAttributes := make([]string, 0, len(resp.Schema.Attributes))
	for name := range resp.Schema.Attributes {
		gotAttributes = append(gotAttributes, name)
	}
	slices.Sort(gotAttributes)

	wantAttributes := []string{
		"config_updates_enabled",
		"created_at",
		"id",
		"managed",
		"modified_at",
		"name",
	}
	if !slices.Equal(gotAttributes, wantAttributes) {
		t.Fatalf("attributes = %v, want %v", gotAttributes, wantAttributes)
	}

	idAttribute, ok := resp.Schema.Attributes["id"].(schema.StringAttribute)
	if !ok {
		t.Fatalf("id attribute type = %T, want schema.StringAttribute", resp.Schema.Attributes["id"])
	}
	if !idAttribute.Required {
		t.Fatal("id attribute Required = false, want true")
	}
}

package provider

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"gopkg.in/yaml.v3"
)

type dataSourceYamlMergeType struct{}

func (t dataSourceYamlMergeType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Merge a list of YAML strings into a single YAML string, where maps are deep merged and list entries are compared against existing list entries and if all primitive values match, the entries are deep merged. ",

		Attributes: map[string]tfsdk.Attribute{
			"input": {
				Description: "A list of YAML strings that is merged into the `output` attribute.",
				Type:        types.ListType{ElemType: types.StringType},
				Required:    true,
			},
			"output": {
				Description: "The merged output.",
				Type:        types.StringType,
				Computed:    true,
			},
			"merge_list_items": {
				Description: "Merge list entries if all primitive values match. Default value is `true`.",
				Type:        types.BoolType,
				Optional:    true,
			},
		},
	}, nil
}

type YamlMerge struct {
	Input          []string     `tfsdk:"input"`
	Output         types.String `tfsdk:"output"`
	MergeListItems types.Bool   `tfsdk:"merge_list_items"`
}

func (t dataSourceYamlMergeType) NewDataSource(ctx context.Context, in tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return dataSourceYamlMerge{
		provider: provider,
	}, diags
}

type dataSourceYamlMerge struct {
	provider provider
}

func (d dataSourceYamlMerge) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var config YamlMerge

	// Read config
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.MergeListItems.IsUnknown() || config.MergeListItems.IsNull() {
		config.MergeListItems.Value = true
	}

	merged := map[interface{}]interface{}{}
	vMerged := reflect.ValueOf(merged)
	for _, input := range config.Input {
		var data map[interface{}]interface{}
		b := []byte(input)

		err := yaml.Unmarshal(b, &data)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading YAML string",
				fmt.Sprintf("Error reading YAML string: %s", err),
			)
			return
		}

		vData := reflect.ValueOf(data)

		err = MergeMaps(vMerged, vData, config.MergeListItems.Value)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error merging YAML",
				fmt.Sprintf("Error merging YAML: %s", err),
			)
			return
		}
	}

	output, err := yaml.Marshal(merged)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting result to YAML",
			fmt.Sprintf("Error converting result to YAML: %s", err),
		)
		return
	}

	config.Output = types.String{Value: string(output)}

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}

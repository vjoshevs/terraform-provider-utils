package provider

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"gopkg.in/yaml.v3"
)

var _ datasource.DataSource = (*yamlMergeDataSource)(nil)

func NewYamlMergeDataSource() datasource.DataSource {
	return &yamlMergeDataSource{}
}

type yamlMergeDataSource struct{}

func (d *yamlMergeDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "utils_yaml_merge"
}

func (t yamlMergeDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Merge a list of YAML strings into a single YAML string, where maps are deep merged and list entries are compared against existing list entries and if all primitive values match, the entries are deep merged. ",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description: "Hexadecimal encoding of the checksum of the output.",
				Type:        types.StringType,
				Computed:    true,
			},
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
	Id             types.String `tfsdk:"id"`
	Input          []string     `tfsdk:"input"`
	Output         types.String `tfsdk:"output"`
	MergeListItems types.Bool   `tfsdk:"merge_list_items"`
}

func (d yamlMergeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config YamlMerge

	// Read config
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.MergeListItems.IsUnknown() || config.MergeListItems.IsNull() {
		config.MergeListItems = types.BoolValue(true)
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

		err = MergeMaps(vMerged, vData, config.MergeListItems.ValueBool())
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

	config.Output = types.StringValue(string(output))

	checksum := sha1.Sum(output)
	config.Id = types.StringValue(hex.EncodeToString(checksum[:]))

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}

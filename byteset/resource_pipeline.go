package byteset

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/seal-io/terraform-provider-byteset/pipeline"
	"github.com/seal-io/terraform-provider-byteset/utils/strx"
)

var _ resource.Resource = (*ResourcePipeline)(nil)

type ResourcePipelineSource struct {
	Address     types.String `tfsdk:"address"`
	ConnMaxOpen types.Int64  `tfsdk:"conn_max_open"`
	ConnMaxIdle types.Int64  `tfsdk:"conn_max_idle"`
	ConnMaxLife types.Int64  `tfsdk:"conn_max_life"`
}

func (r ResourcePipelineSource) Reflect(ctx context.Context) (pipeline.Source, error) {
	return pipeline.NewSource(
		ctx,
		r.Address.ValueString(),
		pipeline.WithConnMaxOpen(int(r.ConnMaxOpen.ValueInt64())),
		pipeline.WithConnMaxIdle(int(r.ConnMaxIdle.ValueInt64())),
		pipeline.WithConnMaxLife(
			time.Duration(r.ConnMaxLife.ValueInt64())*time.Second,
		),
	)
}

type ResourcePipelineDestination struct {
	Address     types.String `tfsdk:"address"`
	ConnMaxOpen types.Int64  `tfsdk:"conn_max_open"`
	ConnMaxIdle types.Int64  `tfsdk:"conn_max_idle"`
	Salt        types.String `tfsdk:"salt"`
}

func (r ResourcePipelineDestination) Reflect(ctx context.Context) (pipeline.Destination, error) {
	return pipeline.NewDestination(
		ctx,
		r.Address.ValueString(),
		pipeline.WithConnMaxOpen(int(r.ConnMaxOpen.ValueInt64())),
		pipeline.WithConnMaxIdle(int(r.ConnMaxIdle.ValueInt64())),
	)
}

type ResourcePipeline struct {
	Source      ResourcePipelineSource      `tfsdk:"source"`
	Destination ResourcePipelineDestination `tfsdk:"destination"`
	ID          types.String                `tfsdk:"id"`
}

func (r ResourcePipeline) Corrupted() bool {
	return r.ID.ValueString() != r.Hash()
}

func (r ResourcePipeline) Equal(l ResourcePipeline) bool {
	return r.Source.Address.Equal(l.Source.Address) &&
		r.Destination.Address.Equal(l.Destination.Address) &&
		r.Destination.Salt.Equal(l.Destination.Salt)
}

func (r ResourcePipeline) Hash() string {
	return strx.Sum(
		r.Source.Address.ValueString(),
		r.Destination.Address.ValueString(),
		r.Destination.Salt.ValueString())
}

func NewResourcePipeline() resource.Resource {
	return ResourcePipeline{}
}

func (r ResourcePipeline) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = strings.Join([]string{req.ProviderTypeName, "pipeline"}, "_")
}

func (r ResourcePipeline) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `Specify the pipeline to seed database.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"source": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"address": schema.StringAttribute{
						Required: true,
						Description: `The address of source, which to provide the dataset, 
choose from local/remote SQL file or database.

  - Local/Remote SQL file format:
	  - file:///path/to/filename
	  - http(s)://...

  - Database address format:
	  - mysql://[username:[password]@][protocol[(address)]][:port][/dbname][?param1=value1&...]
	  - maria://[username:[password]@][protocol[(address)]][:port][/dbname][?param1=value1&...]
	  - postgres://[username:[password]@][address][:port][/dbname][?param1=value1&...]
	  - sqlite:///path/to/filename.db[?param1=value1&...]
	  - oracle://[username:[password]@][address][:port][/service][?param1=value1&...]
	  - mssql://[username:[password]@][address][:port][/instance][?database=dbname&param1=value1&...]`,
					},
					"conn_max_open": schema.Int64Attribute{
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(5),
						Description: `The maximum opening connectors of source database.`,
						Validators: []validator.Int64{
							int64validator.AtLeast(1),
						},
					},
					"conn_max_idle": schema.Int64Attribute{
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(5),
						Description: `The maximum idling connections of source database.`,
						Validators: []validator.Int64{
							int64validator.AtLeast(1),
							int64validator.AtMostSumOf(
								path.MatchRelative().AtParent().AtName("conn_max_open")),
						},
					},
					"conn_max_life": schema.Int64Attribute{
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(5 * 60),
						Description: `The maximum lifetime in seconds of source database.`,
					},
				},
			},
			"destination": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"address": schema.StringAttribute{
						Required: true,
						Description: `The address of destination database, which to receive the dataset.

  - Database address format:
	  - mysql://[username:[password]@][protocol[(address)]][:port][/dbname][?param1=value1&...]
	  - maria://[username:[password]@][protocol[(address)]][:port][/dbname][?param1=value1&...]
	  - postgres://[username:[password]@][address][:port][/dbname][?param1=value1&...]
	  - sqlite:///path/to/filename.db[?param1=value1&...]
	  - oracle://[username:[password]@][address][:port][/service][?param1=value1&...]
	  - mssql://[username:[password]@][address][:port][/instance][?database=dbname&param1=value1&...]`,
					},
					"conn_max_open": schema.Int64Attribute{
						Optional: true,
						Computed: true,
						Default:  int64default.StaticInt64(5),
						Description: `The maximum opening connectors of destination database, 
if the given SQL file is using single transaction, should turn down the "conn_max_open" to 1. 
`,
						Validators: []validator.Int64{
							int64validator.AtLeast(1),
						},
					},
					"conn_max_idle": schema.Int64Attribute{
						Optional: true,
						Computed: true,
						Default:  int64default.StaticInt64(5),
						Description: `The maximum idling connections of destination database, 
if the given SQL file is using single transaction, should turn down the "conn_max_idle" to 1.
`,
						Validators: []validator.Int64{
							int64validator.AtLeast(1),
							int64validator.AtMostSumOf(
								path.MatchRelative().AtParent().AtName("conn_max_open")),
						},
					},
					"salt": schema.StringAttribute{
						Optional: true,
						Description: `The salt assist calculating the destination database has changed 
but the address not, like the database Terraform Managed Resource ID.`,
					},
				},
			},
		},
	}
}

func (r ResourcePipeline) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := r

	if plan.ID.IsNull() {
		resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

		if resp.Diagnostics.HasError() {
			return
		}

		plan.ID = types.StringValue(plan.Hash())
	}

	src, err := plan.Source.Reflect(ctx)
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("source"),
			"Invalid Source",
			"Cannot reflect from source: "+err.Error())

		return
	}

	defer func() { _ = src.Close() }()

	dst, err := plan.Destination.Reflect(ctx)
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("destination"),
			"Invalid Destination",
			"Cannot reflect from destination: "+err.Error())

		return
	}

	defer func() { _ = dst.Close() }()

	err = src.Pipe(ctx, dst)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed Pipe",
			"Cannot pipe from source to destination: "+err.Error())

		return
	}

	plan.Read(
		ctx,
		resource.ReadRequest{
			State:        resp.State,
			Private:      resp.Private,
			ProviderMeta: req.ProviderMeta,
		},
		(*resource.ReadResponse)(resp))
}

func (r ResourcePipeline) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := r

	if state.ID.IsNull() {
		resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

		if resp.Diagnostics.HasError() {
			return
		}
	}

	if state.Corrupted() {
		tflog.Debug(ctx, "State is changed, recreating...")
		resp.State.RemoveResource(ctx)

		return
	}

	// TODO Record something from Destination.

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r ResourcePipeline) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ResourcePipeline

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.Equal(state) {
		tflog.Debug(ctx, "Plan is changed, recreating...")

		plan.ID = types.StringValue(plan.Hash())

		plan.Create(
			ctx,
			resource.CreateRequest{
				Config:       req.Config,
				Plan:         req.Plan,
				ProviderMeta: req.ProviderMeta,
			},
			(*resource.CreateResponse)(resp))
	}
}

func (r ResourcePipeline) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

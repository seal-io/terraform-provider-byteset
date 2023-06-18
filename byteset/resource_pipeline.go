package byteset

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/seal-io/terraform-provider-byteset/pipeline"
	"github.com/seal-io/terraform-provider-byteset/utils/strx"
)

var _ resource.Resource = (*ResourcePipeline)(nil)

type ResourcePipelineSource struct {
	Address types.String `tfsdk:"address"`
	ConnMax types.Int64  `tfsdk:"conn_max"`
}

func (r ResourcePipelineSource) Reflect(ctx context.Context) (pipeline.Source, error) {
	return pipeline.NewSource(
		ctx,
		r.Address.ValueString(),
		int(r.ConnMax.ValueInt64()),
	)
}

type ResourcePipelineDestination struct {
	Address  types.String `tfsdk:"address"`
	ConnMax  types.Int64  `tfsdk:"conn_max"`
	BatchCap types.Int64  `tfsdk:"batch_cap"`
	Salt     types.String `tfsdk:"salt"`
}

func (r ResourcePipelineDestination) Reflect(ctx context.Context) (pipeline.Destination, error) {
	return pipeline.NewDestination(
		ctx,
		r.Address.ValueString(),
		int(r.ConnMax.ValueInt64()),
		int(r.BatchCap.ValueInt64()),
	)
}

type ResourcePipeline struct {
	ID          types.String                `tfsdk:"id"`
	Source      ResourcePipelineSource      `tfsdk:"source"`
	Destination ResourcePipelineDestination `tfsdk:"destination"`
	Timeouts    timeouts.Value              `tfsdk:"timeouts"`
	Cost        types.String                `tfsdk:"cost"`
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
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
						Description: `The address of source, which to provide the dataset, 
choose from local/remote SQL file or database.

  - Local/Remote SQL file format:
	  - file:///path/to/filename
	  - http(s)://...
	  - raw://...
	  - raw+base64://...

  - Database address format:
	  - mysql://[username:[password]@][protocol[(address)]][:port][/dbname][?param1=value1&...]
	  - maria|mariadb://[username:[password]@][protocol[(address)]][:port][/dbname][?param1=value1&...]
	  - postgres|postgresql://[username:[password]@][address][:port][/dbname][?param1=value1&...]
	  - oracle://[username:[password]@][address][:port][/service][?param1=value1&...]
	  - mssql|sqlserver://[username:[password]@][address][:port][/instance][?database=dbname&param1=value1&...]`,
					},
					"conn_max": schema.Int64Attribute{
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(5),
						Description: `The maximum connections of source database.`,
						Validators: []validator.Int64{
							int64validator.AtLeast(1),
						},
					},
				},
			},
			"destination": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"address": schema.StringAttribute{
						Required: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
						Description: `The address of destination database, which to receive the dataset.

  - Database address format:
	  - mysql://[username:[password]@][protocol[(address)]][:port][/dbname][?param1=value1&...]
	  - maria|mariadb://[username:[password]@][protocol[(address)]][:port][/dbname][?param1=value1&...]
	  - postgres|postgresql://[username:[password]@][address][:port][/dbname][?param1=value1&...]
	  - oracle://[username:[password]@][address][:port][/service][?param1=value1&...]
	  - mssql|sqlserver://[username:[password]@][address][:port][/instance][?database=dbname&param1=value1&...]`,
					},
					"conn_max": schema.Int64Attribute{
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(5),
						Description: `The maximum opening connectors of destination database.`,
						Validators: []validator.Int64{
							int64validator.AtLeast(1),
						},
					},
					"batch_cap": schema.Int64Attribute{
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(500),
						Description: `The maximum value statement number for once insert statements.`,
						Validators: []validator.Int64{
							int64validator.AtLeast(1),
						},
					},
					"salt": schema.StringAttribute{
						Optional: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
						Description: `The salt assist calculating the destination database has changed 
but the address not, like the database Terraform Managed Resource ID.`,
					},
				},
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Update: true,
			}),
			"cost": schema.StringAttribute{
				Computed:    true,
				Description: `The time spent on this transfer.`,
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

		timeout, diags := r.Timeouts.Create(ctx, 30*time.Minute)
		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		var cancel func()

		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
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

	start := time.Now()

	if err = src.Pipe(ctx, dst); err != nil {
		resp.Diagnostics.AddError(
			"Failed Pipe",
			"Cannot pipe from source to destination: "+err.Error())

		return
	}
	plan.Cost = types.StringValue(time.Since(start).String())

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

		timeout, diags := r.Timeouts.Create(ctx, 30*time.Minute)
		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		var cancel func()

		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()

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

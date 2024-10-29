//go:generate go run ../inventory-data/resource_types_generator.go --provider aws --output resource_types.go --index-map ../pkg/steampipe/table_index_map.go && gofmt -w -s resource_types.go  && goimports -w resource_types.go

package aws

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"sort"
	"strings"

	"github.com/aws/smithy-go"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/opengovern/og-aws-describer/aws/describer"
	"github.com/opengovern/og-util/pkg/describe/enums"
	"github.com/opengovern/og-util/pkg/source"
)

type ResourceDescriber func(context.Context, aws.Config, string, []string, string, enums.DescribeTriggerType, *describer.StreamSender) (*Resources, error)
type SingleResourceDescriber func(context.Context, aws.Config, string, []string, string, map[string]string, enums.DescribeTriggerType) (*Resources, error)

type ResourceType struct {
	Connector source.Type

	ResourceName  string
	ResourceLabel string
	ServiceName   string

	Tags map[string][]string

	ListDescriber ResourceDescriber
	GetDescriber  SingleResourceDescriber

	TerraformName        []string
	TerraformServiceName string

	FastDiscovery bool
	CostDiscovery bool
	Summarize     bool
}

func (r ResourceType) GetConnector() source.Type {
	return r.Connector
}

func (r ResourceType) GetResourceName() string {
	return r.ResourceName
}

func (r ResourceType) GetResourceLabel() string {
	return r.ResourceLabel
}

func (r ResourceType) GetServiceName() string {
	return r.ServiceName
}

func (r ResourceType) GetTags() map[string][]string {
	return r.Tags
}

func (r ResourceType) GetTerraformName() []string {
	return r.TerraformName
}

func (r ResourceType) GetTerraformServiceName() string {
	return r.TerraformServiceName
}

func (r ResourceType) IsFastDiscovery() bool {
	return r.FastDiscovery
}

func (r ResourceType) IsCostDiscovery() bool {
	return r.CostDiscovery
}

func (r ResourceType) IsSummarized() bool {
	return r.Summarize
}

func ListResourceTypes() []string {
	var list []string
	for k := range resourceTypes {
		list = append(list, k)
	}

	sort.Strings(list)
	return list
}

func ListFastDiscoveryResourceTypes() []string {
	var list []string
	for k, v := range resourceTypes {
		if v.FastDiscovery {
			list = append(list, k)
		}
	}

	sort.Strings(list)
	return list
}

func ListSummarizeResourceTypes() []string {
	var list []string
	for k, v := range resourceTypes {
		if v.Summarize {
			list = append(list, k)
		}
	}

	sort.Strings(list)
	return list
}

func GetResourceType(resourceType string) (*ResourceType, error) {
	if r, ok := resourceTypes[resourceType]; ok {
		return &r, nil
	}
	resourceType = strings.ToLower(resourceType)
	for k, v := range resourceTypes {
		k := strings.ToLower(k)
		v := v
		if k == resourceType {
			return &v, nil
		}
	}
	return nil, fmt.Errorf("resource type %s not found", resourceType)
}

func GetResourceTypesMap() map[string]ResourceType {
	return resourceTypes
}

type Resources struct {
	Resources map[string][]describer.Resource
	Errors    map[string]string
	ErrorCode string
}

func GetResources(ctx context.Context, logger *zap.Logger,
	resourceType string, triggerType enums.DescribeTriggerType,
	accountId string, regions []string,
	credAccountId, accessKey, secretKey, sessionToken, assumeRoleName, assumeAdminRoleName string, externalId *string,
	includeDisabledRegions bool, stream *describer.StreamSender) (*Resources, error) {
	var err error
	var cfg aws.Config

	needToRunOnOrgMaster := false
	if strings.HasPrefix(strings.ToLower(resourceType), "aws::costexplorer") {
		needToRunOnOrgMaster = true
	} else if strings.HasPrefix(strings.ToLower(resourceType), "aws::ssoadmin::") {
		needToRunOnOrgMaster = true
	}

	if accountId != credAccountId && !needToRunOnOrgMaster {
		assumeRoleArn := GetRoleArnFromName(accountId, assumeRoleName)
		cfg, err = GetConfig(ctx, accessKey, secretKey, sessionToken, assumeRoleArn, externalId)
	} else if accountId != credAccountId && needToRunOnOrgMaster {
		assumeAdminRoleArn := GetRoleArnFromName(credAccountId, assumeAdminRoleName)
		cfg, err = GetConfig(ctx, accessKey, secretKey, sessionToken, assumeAdminRoleArn, externalId)
	} else {
		assumeAdminRoleArn := GetRoleArnFromName(accountId, assumeAdminRoleName)
		cfg, err = GetConfig(ctx, accessKey, secretKey, sessionToken, assumeAdminRoleArn, externalId)
	}
	if err != nil {
		return nil, err
	}

	if len(regions) == 0 {
		cfgClone := cfg.Copy()
		cfgClone.Region = "us-east-1"

		rs, err := getAllRegions(ctx, cfgClone, includeDisabledRegions)
		if err != nil {
			return nil, err
		}

		for _, r := range rs {
			regions = append(regions, *r.RegionName)
		}
	}

	sort.Slice(regions, func(i, j int) bool {
		if regions[i] == "us-east-1" {
			return true
		}
		if regions[j] == "us-east-1" {
			return false
		}

		return regions[i] < regions[j]
	})

	logger.Info("Running the describer started")
	resources, err := describe(ctx, logger, cfg, accountId, regions, resourceType, triggerType, stream)
	if err != nil {
		return nil, err
	}
	logger.Info("Running the describer finished")

	return resources, nil
}

func GetSingleResource(
	ctx context.Context,
	resourceType string,
	triggerType enums.DescribeTriggerType,
	accountId string,
	regions []string,
	accessKey,
	secretKey,
	sessionToken,
	assumeRoleName string,
	externalId *string,
	includeDisabledRegions bool,
	fields map[string]string,
) (*Resources, error) {
	assumeRoleArn := GetRoleArnFromName(accountId, assumeRoleName)
	cfg, err := GetConfig(ctx, accessKey, secretKey, sessionToken, assumeRoleArn, externalId)
	if err != nil {
		return nil, err
	}

	if len(regions) == 0 {
		cfgClone := cfg.Copy()
		cfgClone.Region = "us-east-1"

		rs, err := getAllRegions(ctx, cfgClone, includeDisabledRegions)
		if err != nil {
			return nil, err
		}

		for _, r := range rs {
			regions = append(regions, *r.RegionName)
		}
	}

	resources, err := describeSingle(ctx, cfg, accountId, regions, resourceType, fields, triggerType)
	if err != nil {
		return nil, err
	}

	return resources, nil
}

func describeSingle(
	ctx context.Context,
	cfg aws.Config,
	account string,
	regions []string,
	resourceType string,
	fields map[string]string,
	triggerType enums.DescribeTriggerType) (*Resources, error) {
	resourceTypeObject, ok := resourceTypes[resourceType]
	if !ok {
		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
	}

	return resourceTypeObject.GetDescriber(ctx, cfg, account, regions, resourceType, fields, triggerType)
}

func describe(ctx context.Context, logger *zap.Logger, cfg aws.Config, account string, regions []string, resourceType string, triggerType enums.DescribeTriggerType, stream *describer.StreamSender) (*Resources, error) {
	resourceTypeObject, ok := resourceTypes[resourceType]
	if !ok {
		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
	}
	ctx = describer.WithLogger(ctx, logger)

	return resourceTypeObject.ListDescriber(ctx, cfg, account, regions, resourceType, triggerType, stream)
}

func ParallelDescribeRegionalSingleResource(describe func(context.Context, aws.Config, map[string]string) ([]describer.Resource, error)) SingleResourceDescriber {
	type result struct {
		region    string
		resources []describer.Resource
		err       error
		errorCode string
	}
	return func(ctx context.Context, cfg aws.Config, account string, regions []string, rType string, fields map[string]string, triggerType enums.DescribeTriggerType) (*Resources, error) {
		input := make(chan result, len(regions))
		for _, region := range regions {
			go func(r string) {
				defer func() {
					if err := recover(); err != nil {
						//stack := debug.Stack()
						//input <- result{region: r, resources: nil, err: fmt.Errorf("paniced: %v\n%s", err, string(stack))}
						input <- result{region: r, resources: nil, err: fmt.Errorf("paniced: %v", err)}
					}
				}()
				// Make a shallow copy and override the default region
				rCfg := cfg.Copy()
				rCfg.Region = r

				partition, _ := PartitionOf(r)
				ctx := describer.WithDescribeContext(ctx, describer.DescribeContext{
					AccountID:   account,
					Region:      r,
					KaytuRegion: r,
					Partition:   partition,
				})
				ctx = describer.WithTriggerType(ctx, triggerType)
				resources, err := describe(ctx, rCfg, fields)
				errCode := ""
				if err != nil {
					var ae smithy.APIError
					if errors.As(err, &ae) {
						errCode = ae.ErrorCode()
					}
				}
				input <- result{region: r, resources: resources, err: err, errorCode: errCode}
			}(region)
		}

		output := Resources{
			Resources: make(map[string][]describer.Resource, len(regions)),
			Errors:    make(map[string]string, len(regions)),
		}
		for range regions {
			resp := <-input
			if resp.err != nil {
				if !IsUnsupportedOrInvalidError(rType, resp.region, resp.err) {
					output.Errors[resp.region] = resp.err.Error()
					output.ErrorCode = resp.errorCode
					continue
				}
			}

			if resp.resources == nil {
				resp.resources = []describer.Resource{}
			}

			partition, _ := PartitionOf(resp.region)
			for i := range resp.resources {
				resp.resources[i].Account = account
				resp.resources[i].Region = resp.region
				resp.resources[i].Partition = partition
				resp.resources[i].Type = rType
			}

			output.Resources[resp.region] = resp.resources
		}

		return &output, nil
	}
}

func SequentialDescribeRegional(describe func(context.Context, aws.Config, *describer.StreamSender) ([]describer.Resource, error)) ResourceDescriber {
	return func(ctx context.Context, cfg aws.Config, account string, regions []string, rType string, triggerType enums.DescribeTriggerType, stream *describer.StreamSender) (*Resources, error) {
		output := Resources{
			Resources: make(map[string][]describer.Resource, len(regions)),
			Errors:    make(map[string]string, len(regions)),
		}

		for _, region := range regions {
			// Make a shallow copy and override the default region
			rCfg := cfg.Copy()
			rCfg.Region = region

			partition, _ := PartitionOf(region)
			ctx := describer.WithDescribeContext(ctx, describer.DescribeContext{
				AccountID:   account,
				Region:      region,
				KaytuRegion: region,
				Partition:   partition,
			})
			ctx = describer.WithTriggerType(ctx, triggerType)
			resources, err := describe(ctx, rCfg, stream)
			if err != nil {
				if !IsUnsupportedOrInvalidError(rType, region, err) {
					errCode := ""
					if err != nil {
						var ae smithy.APIError
						if errors.As(err, &ae) {
							errCode = ae.ErrorCode()
						}
					}
					output.Errors[region] = err.Error()
					output.ErrorCode = errCode
				}
				continue
			}

			if resources == nil {
				resources = []describer.Resource{}
			}

			for i := range resources {
				resources[i].Account = account
				resources[i].Region = region
				resources[i].Partition = partition
				resources[i].Type = rType
			}

			output.Resources[region] = resources
		}
		return &output, nil
	}
}

// Parallel describe the resources across the reigons. Failure in one regions won't affect
// other regions.
func ParallelDescribeRegional(describe func(context.Context, aws.Config, *describer.StreamSender) ([]describer.Resource, error)) ResourceDescriber {
	type result struct {
		region    string
		resources []describer.Resource
		err       error
		errorCode string
	}
	return func(ctx context.Context, cfg aws.Config, account string, regions []string, rType string, triggerType enums.DescribeTriggerType, stream *describer.StreamSender) (*Resources, error) {
		fmt.Println("ParallelDescribeRegional")
		input := make(chan result, len(regions))
		for _, region := range regions {
			go func(r string) {
				defer func() {
					if err := recover(); err != nil {
						//stack := debug.Stack()
						//input <- result{region: r, resources: nil, err: fmt.Errorf("paniced: %v\n%s", err, string(stack))}
						input <- result{region: r, resources: nil, err: fmt.Errorf("paniced: %v", err)}
					}
				}()
				// Make a shallow copy and override the default region
				rCfg := cfg.Copy()
				rCfg.Region = r

				partition, _ := PartitionOf(r)
				describeCtx := describer.DescribeContext{
					AccountID:   account,
					Region:      r,
					KaytuRegion: r,
					Partition:   partition,
				}

				fmt.Println("ParallelDescribeRegional for region", r, rCfg.Region, describeCtx.Region)
				ctx := describer.WithDescribeContext(ctx, describeCtx)
				ctx = describer.WithTriggerType(ctx, triggerType)
				fmt.Println("running describe")
				resources, err := describe(ctx, rCfg, stream)
				fmt.Println("describe finished", err)
				errCode := ""
				if err != nil {
					var ae smithy.APIError
					if errors.As(err, &ae) {
						errCode = ae.ErrorCode()
					}
				}
				input <- result{region: r, resources: resources, err: err, errorCode: errCode}
			}(region)
		}

		output := Resources{
			Resources: make(map[string][]describer.Resource, len(regions)),
			Errors:    make(map[string]string, len(regions)),
			ErrorCode: "",
		}
		for range regions {
			fmt.Println("ParallelDescribeRegional waiting for result")
			resp := <-input
			fmt.Println("ParallelDescribeRegional got a result")
			if resp.err != nil {
				if !IsUnsupportedOrInvalidError(rType, resp.region, resp.err) {
					output.Errors[resp.region] = resp.err.Error()
					output.ErrorCode = resp.errorCode
					continue
				}
			}

			if resp.resources == nil {
				resp.resources = []describer.Resource{}
			}

			partition, _ := PartitionOf(resp.region)
			for i := range resp.resources {
				resp.resources[i].Account = account
				resp.resources[i].Region = resp.region
				resp.resources[i].Partition = partition
				resp.resources[i].Type = rType
			}

			output.Resources[resp.region] = resp.resources
		}
		fmt.Println("ParallelDescribeRegional finished")

		return &output, nil
	}
}

// Sequentially describe the resources. If anyone of the regions fails, it will move on to the next region.
func SequentialDescribeGlobal(describe func(context.Context, aws.Config, *describer.StreamSender) ([]describer.Resource, error)) ResourceDescriber {
	return func(ctx context.Context, cfg aws.Config, account string, regions []string, rType string, triggerType enums.DescribeTriggerType, stream *describer.StreamSender) (*Resources, error) {
		output := Resources{
			Resources: make(map[string][]describer.Resource, len(regions)),
			Errors:    make(map[string]string, len(regions)),
		}

		for _, region := range regions {
			// Make a shallow copy and override the default region
			rCfg := cfg.Copy()
			rCfg.Region = region

			partition, _ := PartitionOf(region)
			ctx := describer.WithDescribeContext(ctx, describer.DescribeContext{
				AccountID:   account,
				Region:      region,
				KaytuRegion: "global",
				Partition:   partition,
			})
			ctx = describer.WithTriggerType(ctx, triggerType)
			resources, err := describe(ctx, rCfg, stream)
			if err != nil {
				if !IsUnsupportedOrInvalidError(rType, region, err) {
					errCode := ""
					if err != nil {
						var ae smithy.APIError
						if errors.As(err, &ae) {
							errCode = ae.ErrorCode()
						}
					}
					output.Errors[region] = err.Error()
					output.ErrorCode = errCode
				}
				continue
			}

			if resources == nil {
				resources = []describer.Resource{}
			}

			for i := range resources {
				resources[i].Account = account
				resources[i].Region = "global"
				resources[i].Partition = partition
				resources[i].Type = rType
			}

			output.Resources[region] = resources

			break
		}

		m := map[string]interface{}{}
		for k, v := range output.Resources {
			var newV []describer.Resource
			for _, r := range v {
				if _, ok := m[r.UniqueID()]; ok {
					continue
				}

				m[r.UniqueID()] = struct{}{}
				newV = append(newV, r)
			}
			output.Resources[k] = newV
		}

		return &output, nil
	}
}

func GetResourceTypeByTerraform(terraformType string) string {
	for t, v := range resourceTypes {
		for _, name := range v.TerraformName {
			if name == terraformType {
				return t
			}
		}
	}
	return ""
}

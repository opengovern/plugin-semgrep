package provider

import (
	"errors"
	"regexp"
	"strings"

	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/smithy-go"
)

const (
	ErrSubscriptionRequired = "SubscriptionRequiredException"
)

func IsUnsupportedOrInvalidError(resource, region string, err error) bool {
	var ae smithy.APIError
	if errors.As(err, &ae) {
		switch ae.ErrorCode() {
		case "InvalidAction":
			return true
		case "UnsupportedOperation":
			return true
		case "AccessDeniedException":
			re := &awshttp.ResponseError{}
			if errors.As(err, &re) {
				if re.HTTPStatusCode() == 400 && strings.Contains(strings.ToLower(ae.ErrorMessage()), strings.ToLower("UnknownError")) {
					return true
				}
			}
			return false
		}
	}
	var responseError *awshttp.ResponseError
	regex := regexp.MustCompile(`dial tcp: lookup .* on .*: no such host`)
	if errors.As(err, &responseError) {
		if responseError.HTTPStatusCode() == 0 && regex.MatchString(responseError.Error()) {
			return true
		}
	}

	// The following resources types describe calls are not supported
	// in the corresponding regions. The error message is usually
	// not very clear about the error. For now just ignoring them.
	switch resource {
	case "AWS::MemoryDb::Cluster":
		if isInRegion(region, "ap-northeast-3") {
			return true
		}
	case "AWS::CloudSearch::Domain":
		if isInRegion(region, "ap-northeast-3", "ca-central-1", "eu-west-2", "eu-north-1", "eu-west-3", "us-east-2", "ap-south-1") {
			return true
		}
	case "AWS::Amplify::App":
		if isInRegion(region, "ap-northeast-3") {
			return true
		}
	case "AWS::CodeArtifact::Domain":
		if isInRegion(region, "ap-northeast-3", "us-west-1", "ap-northeast-2", "sa-east-1", "ca-central-1") {
			return true
		}
	case "AWS::Inspector::AssessmentTemplate",
		"AWS::Inspector::Exclusion",
		"AWS::Inspector::Finding",
		"AWS::Inspector::AssessmentRun",
		"AWS::Inspector::AssessmentTarget":
		if isInRegion(region, "eu-west-3", "ca-central-1", "ap-southeast-1", "sa-east-1", "ap-northeast-3") {
			return true
		}
	case "AWS::GlobalAccelerator::EndpointGroup",
		"AWS::GlobalAccelerator::Listener",
		"AWS::GlobalAccelerator::Accelerator":
		if isInRegion(region, "eu-west-1", "ap-northeast-3", "ap-northeast-2", "eu-north-1", "ca-central-1", "eu-west-2", "sa-east-1", "ap-northeast-1", "eu-west-3", "eu-central-1", "us-west-1", "ap-south-1", "ap-southeast-1", "us-east-2", "ap-southeast-2", "us-east-1") {
			return true
		}
	case "AWS::Route53Resolver::ResolverDNSSECConfig",
		"AWS::Route53Resolver::ResolverQueryLoggingConfigAssociation",
		"AWS::Route53Resolver::ResolverQueryLoggingConfig":
		if isInRegion(region, "ap-northeast-3") {
			return true
		}
	case "AWS::RedshiftServerless::Namespace", "AWS::RedshiftServerless::Snapshot":
		if isInRegion(region, "ap-northeast-3", "sa-east-1") {
			return true
		}
	case "AWS::RDS::DBProxy",
		"AWS::RDS::DBProxyTargetGroup",
		"AWS::RDS::DBProxyEndpoint",
		"AWS::Lambda::CodeSigningConfig",
		"AWS::S3::StorageLens":
		if isInRegion(region, "ap-northeast-3", "eu-north-1", "eu-west-3", "sa-east-1") {
			return true
		}
	case "AWS::Workspaces::ConnectionAlias",
		"AWS::Workspaces::Workspace",
		"AWS::Workspaces::Bundle":
		if isInRegion(region, "ap-northeast-3", "eu-north-1", "eu-west-3", "us-east-2", "us-west-1") {
			return true
		}
	case "AWS::DAX::Cluster",
		"AWS::DAX::ParameterGroup",
		"AWS::DAX::Parameter",
		"AWS::DAX::SubnetGroup":
		// Region whitelist https://docs.aws.amazon.com/general/latest/gr/ddb.html
		if !isInRegion(region, "cn-north-1", "cn-northwest-1", "eu-west-2", "eu-west-3", "ap-northeast-1",
			"ap-south-1", "ap-southeast-1", "ap-southeast-2", "eu-central-1", "eu-west-1",
			"sa-east-1", "us-east-1", "us-east-2", "us-west-1", "us-west-2") {
			return true
		}
	case "AWS::AppStream::Application",
		"AWS::AppStream::Stack",
		"AWS::AppStream::Fleet":
		// Region whitelist https://docs.aws.amazon.com/general/latest/gr/aas2.html#aas2_region
		if !isInRegion(region, "us-east-2", "us-east-1", "us-west-2", "ap-south-1", "ap-northeast-2", "ap-southeast-1",
			"ap-southeast-2", "ap-northeast-1", "ca-central-1", "eu-central-1", "eu-west-1", "eu-west-2", "us-gov-west-1") {
			return true
		}
	case "AWS::Keyspaces::Keyspace", "AWS::Keyspaces::Table":
		if isInRegion(region, "ap-northeast-3") {
			return true
		}
	case "AWS::Grafana::Workspace":
		// Region whitelist https://docs.aws.amazon.com/grafana/latest/userguide/what-is-Amazon-Managed-Service-Grafana.html
		if !isInRegion(region, "us-east-2", "us-east-1", "us-west-2", "ap-northeast-2", "ap-southeast-1",
			"ap-southeast-2", "ap-northeast-1", "eu-central-1", "eu-west-1", "eu-west-2") {
			return true
		}
	case "AWS::AMP::Workspace":
		// Region whitelist https://docs.aws.amazon.com/prometheus/latest/userguide/what-is-Amazon-Managed-Service-Prometheus.html
		if !isInRegion(region, "us-east-2", "us-east-1", "us-west-2", "ap-southeast-1", "ap-southeast-2",
			"ap-northeast-1", "eu-central-1", "eu-west-1", "eu-west-2", "eu-north-1") {
			return true
		}
	case "AWS::MWAA::Environment":
		// Region whitelist https://docs.aws.amazon.com/mwaa/latest/userguide/what-is-mwaa.html#regions-mwaa
		if !isInRegion(region, "eu-central-1", "eu-west-1", "eu-west-2", "eu-west-3", "ap-south-1", "ap-southeast-1",
			"ap-southeast-2", "ap-northeast-1", "ap-northeast-2", "us-east-1", "us-east-2", "us-west-2", "ca-central-1", "sa-east-1") {
			return true
		}
	case "AWS::CodeArtifact::Repository":
		// Region whitelist https://docs.aws.amazon.com/general/latest/gr/codeartifact.html
		if !isInRegion(region, "us-east-2", "us-east-1", "us-west-2", "ap-south-1", "ap-southeast-1", "ap-southeast-2",
			"ap-northeast-1", "eu-central-1", "eu-west-1", "eu-west-2", "eu-south-1", "eu-west-3", "eu-north-1") {
			return true
		}
	case "AWS::CodeStar::Project":
		// Region whitelist https://docs.aws.amazon.com/general/latest/gr/codestar.html
		if !isInRegion(region, "us-east-2", "us-east-1", "us-west-1", "us-west-2", "ap-northeast-2", "ap-southeast-1",
			"ap-southeast-2", "ap-northeast-1", "ca-central-1", "eu-central-1", "eu-west-1", "eu-west-2", "eu-north-1") {
			return true
		}
	case "AWS::OpsWorksCM::Server":
		//Region whitelist https://docs.aws.amazon.com/general/latest/gr/opsworks-cm.html
		if !isInRegion(region, "us-east-2", "us-east-1", "us-west-1", "us-west-2", "ap-southeast-1", "ap-southeast-2",
			"ap-northeast-1", "eu-central-1", "eu-west-1") {
			return true
		}
	case "AWS::CodePipeline::Pipeline":
		// Region whitelist https://docs.aws.amazon.com/general/latest/gr/codepipeline.html
		if !isInRegion(region, "us-east-2", "us-east-1", "us-west-1", "us-west-2", "ap-east-1", "ap-south-1", "ap-northeast-2",
			"ap-southeast-1", "ap-southeast-2", "ap-northeast-1", "ca-central-1", "eu-central-1", "eu-west-1", "eu-west-2",
			"eu-south-1", "eu-west-3", "eu-north-1", "sa-east-1", "us-gov-west-1") {
			return true
		}
	}
	return false
}

func isInRegion(region string, regions ...string) bool {
	for _, r := range regions {
		if strings.EqualFold(region, r) {
			return true
		}
	}

	return false
}

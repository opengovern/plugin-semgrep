//go:generate go run ./gen/main.go --file $GOFILE --output ../../pkg/opengovernance-es-sdk/aws_resources_clients.go --type aws

package model

import (
	dynamodb2 "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	lakeformationTypes "github.com/aws/aws-sdk-go-v2/service/lakeformation/types"
	networkfirewall2 "github.com/aws/aws-sdk-go-v2/service/networkfirewall"
	"github.com/aws/aws-sdk-go-v2/service/opensearchserverless"
	types6 "github.com/aws/aws-sdk-go-v2/service/opensearchserverless/types"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroups"
	types4 "github.com/aws/aws-sdk-go-v2/service/resourcegroups/types"
	"github.com/aws/aws-sdk-go-v2/service/timestreamwrite"
	types5 "github.com/aws/aws-sdk-go-v2/service/timestreamwrite/types"
	"time"

	accessanalyzer "github.com/aws/aws-sdk-go-v2/service/accessanalyzer/types"
	account "github.com/aws/aws-sdk-go-v2/service/account/types"
	acm "github.com/aws/aws-sdk-go-v2/service/acm/types"
	acmpca "github.com/aws/aws-sdk-go-v2/service/acmpca/types"
	amp "github.com/aws/aws-sdk-go-v2/service/amp/types"
	amplify "github.com/aws/aws-sdk-go-v2/service/amplify/types"
	apigateway "github.com/aws/aws-sdk-go-v2/service/apigateway/types"
	apigatewayv2 "github.com/aws/aws-sdk-go-v2/service/apigatewayv2/types"
	appconfig "github.com/aws/aws-sdk-go-v2/service/appconfig/types"
	applicationautoscaling "github.com/aws/aws-sdk-go-v2/service/applicationautoscaling/types"
	appstream "github.com/aws/aws-sdk-go-v2/service/appstream/types"
	athena "github.com/aws/aws-sdk-go-v2/service/athena/types"
	auditmanager "github.com/aws/aws-sdk-go-v2/service/auditmanager/types"
	autoscaling "github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	backupop "github.com/aws/aws-sdk-go-v2/service/backup"
	backupservice "github.com/aws/aws-sdk-go-v2/service/backup"
	backup "github.com/aws/aws-sdk-go-v2/service/backup/types"
	batch "github.com/aws/aws-sdk-go-v2/service/batch/types"
	cloudcontrol "github.com/aws/aws-sdk-go-v2/service/cloudcontrol/types"
	cloudformationop "github.com/aws/aws-sdk-go-v2/service/cloudformation"
	cloudformation "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	cloudfrontop "github.com/aws/aws-sdk-go-v2/service/cloudfront"
	cloudfront "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	cloudsearch "github.com/aws/aws-sdk-go-v2/service/cloudsearch/types"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	cloudtrailtypes "github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	cloudwatch "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	cloudwatchlogs2 "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	cloudwatchlogs "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	codeartifact "github.com/aws/aws-sdk-go-v2/service/codeartifact/types"
	codebuild "github.com/aws/aws-sdk-go-v2/service/codebuild/types"
	codecommit "github.com/aws/aws-sdk-go-v2/service/codecommit/types"
	codedeploy "github.com/aws/aws-sdk-go-v2/service/codedeploy/types"
	codepipeline "github.com/aws/aws-sdk-go-v2/service/codepipeline/types"
	codestarop "github.com/aws/aws-sdk-go-v2/service/codestar"
	configservice "github.com/aws/aws-sdk-go-v2/service/configservice/types"
	dms "github.com/aws/aws-sdk-go-v2/service/databasemigrationservice/types"
	dax "github.com/aws/aws-sdk-go-v2/service/dax/types"
	directconnect "github.com/aws/aws-sdk-go-v2/service/directconnect/types"
	directoryservice "github.com/aws/aws-sdk-go-v2/service/directoryservice/types"
	dlm "github.com/aws/aws-sdk-go-v2/service/dlm/types"
	docdb "github.com/aws/aws-sdk-go-v2/service/docdb/types"
	drs2 "github.com/aws/aws-sdk-go-v2/service/drs"
	drs "github.com/aws/aws-sdk-go-v2/service/drs/types"
	dynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	dynamodbstream "github.com/aws/aws-sdk-go-v2/service/dynamodbstreams/types"
	ec2op "github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2 "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	ecrop "github.com/aws/aws-sdk-go-v2/service/ecr"
	ecr "github.com/aws/aws-sdk-go-v2/service/ecr/types"
	ecrpublicop "github.com/aws/aws-sdk-go-v2/service/ecrpublic"
	ecrpublic "github.com/aws/aws-sdk-go-v2/service/ecrpublic/types"
	ecs "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	efs "github.com/aws/aws-sdk-go-v2/service/efs/types"
	eks "github.com/aws/aws-sdk-go-v2/service/eks/types"
	elasticache "github.com/aws/aws-sdk-go-v2/service/elasticache/types"
	elasticbeanstalk "github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk/types"
	elasticloadbalancing "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing/types"
	elasticloadbalancingv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	es "github.com/aws/aws-sdk-go-v2/service/elasticsearchservice/types"
	emr "github.com/aws/aws-sdk-go-v2/service/emr/types"
	eventbridgeop "github.com/aws/aws-sdk-go-v2/service/eventbridge"
	eventbridge "github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	firehose "github.com/aws/aws-sdk-go-v2/service/firehose/types"
	fms "github.com/aws/aws-sdk-go-v2/service/fms/types"
	fsx "github.com/aws/aws-sdk-go-v2/service/fsx/types"
	glacier "github.com/aws/aws-sdk-go-v2/service/glacier/types"
	globalaccelerator "github.com/aws/aws-sdk-go-v2/service/globalaccelerator/types"
	glueop "github.com/aws/aws-sdk-go-v2/service/glue"
	glue "github.com/aws/aws-sdk-go-v2/service/glue/types"
	grafana "github.com/aws/aws-sdk-go-v2/service/grafana/types"
	guarddutyop "github.com/aws/aws-sdk-go-v2/service/guardduty"
	guardduty "github.com/aws/aws-sdk-go-v2/service/guardduty/types"
	health "github.com/aws/aws-sdk-go-v2/service/health/types"
	iamop "github.com/aws/aws-sdk-go-v2/service/iam"
	iam "github.com/aws/aws-sdk-go-v2/service/iam/types"
	identitystore2 "github.com/aws/aws-sdk-go-v2/service/identitystore"
	identitystore "github.com/aws/aws-sdk-go-v2/service/identitystore/types"
	imagebuilder "github.com/aws/aws-sdk-go-v2/service/imagebuilder/types"
	inspector "github.com/aws/aws-sdk-go-v2/service/inspector/types"
	inspector2 "github.com/aws/aws-sdk-go-v2/service/inspector2/types"
	kafkaop "github.com/aws/aws-sdk-go-v2/service/kafka"
	kafka "github.com/aws/aws-sdk-go-v2/service/kafka/types"
	keyspaces "github.com/aws/aws-sdk-go-v2/service/keyspaces/types"
	kinesis "github.com/aws/aws-sdk-go-v2/service/kinesis/types"
	kinesisanalyticsv2 "github.com/aws/aws-sdk-go-v2/service/kinesisanalyticsv2/types"
	kinesisvideo "github.com/aws/aws-sdk-go-v2/service/kinesisvideo/types"
	kms "github.com/aws/aws-sdk-go-v2/service/kms/types"
	lambdaop "github.com/aws/aws-sdk-go-v2/service/lambda"
	lambda "github.com/aws/aws-sdk-go-v2/service/lambda/types"
	lightsail "github.com/aws/aws-sdk-go-v2/service/lightsail/types"
	macie2op "github.com/aws/aws-sdk-go-v2/service/macie2"
	mediastoreop "github.com/aws/aws-sdk-go-v2/service/mediastore"
	mediastore "github.com/aws/aws-sdk-go-v2/service/mediastore/types"
	memorydb "github.com/aws/aws-sdk-go-v2/service/memorydb/types"
	mgn "github.com/aws/aws-sdk-go-v2/service/mgn/types"
	"github.com/aws/aws-sdk-go-v2/service/mq"
	mwaa "github.com/aws/aws-sdk-go-v2/service/mwaa/types"
	neptune "github.com/aws/aws-sdk-go-v2/service/neptune/types"
	networkfirewall "github.com/aws/aws-sdk-go-v2/service/networkfirewall/types"
	oamop "github.com/aws/aws-sdk-go-v2/service/oam"
	oam "github.com/aws/aws-sdk-go-v2/service/oam/types"
	opensearch "github.com/aws/aws-sdk-go-v2/service/opensearch/types"
	opsworkscm "github.com/aws/aws-sdk-go-v2/service/opsworkscm/types"
	organizations "github.com/aws/aws-sdk-go-v2/service/organizations/types"
	pinpoint "github.com/aws/aws-sdk-go-v2/service/pinpoint/types"
	pipesop "github.com/aws/aws-sdk-go-v2/service/pipes"
	pipes "github.com/aws/aws-sdk-go-v2/service/pipes/types"
	ram "github.com/aws/aws-sdk-go-v2/service/ram/types"
	rds_sdkv2 "github.com/aws/aws-sdk-go-v2/service/rds"
	rds "github.com/aws/aws-sdk-go-v2/service/rds/types"
	redshiftop "github.com/aws/aws-sdk-go-v2/service/redshift"
	redshift "github.com/aws/aws-sdk-go-v2/service/redshift/types"
	redshiftserverlesstypes "github.com/aws/aws-sdk-go-v2/service/redshiftserverless/types"
	resourceexplorer2 "github.com/aws/aws-sdk-go-v2/service/resourceexplorer2/types"
	types2 "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi/types"
	route53op "github.com/aws/aws-sdk-go-v2/service/route53"
	route53 "github.com/aws/aws-sdk-go-v2/service/route53/types"
	route53domainsop "github.com/aws/aws-sdk-go-v2/service/route53domains"
	route53domains "github.com/aws/aws-sdk-go-v2/service/route53domains/types"
	route53resolverop "github.com/aws/aws-sdk-go-v2/service/route53resolver"
	route53resolver "github.com/aws/aws-sdk-go-v2/service/route53resolver/types"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3"
	s3 "github.com/aws/aws-sdk-go-v2/service/s3/types"
	s3controlop "github.com/aws/aws-sdk-go-v2/service/s3control"
	s3control "github.com/aws/aws-sdk-go-v2/service/s3control/types"
	sagemakerop "github.com/aws/aws-sdk-go-v2/service/sagemaker"
	sagemaker "github.com/aws/aws-sdk-go-v2/service/sagemaker/types"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	securityhubop "github.com/aws/aws-sdk-go-v2/service/securityhub"
	securityhub "github.com/aws/aws-sdk-go-v2/service/securityhub/types"
	securitylake "github.com/aws/aws-sdk-go-v2/service/securitylake/types"
	serverlessapplicationrepositoryop "github.com/aws/aws-sdk-go-v2/service/serverlessapplicationrepository"
	serverlessapplicationrepository "github.com/aws/aws-sdk-go-v2/service/serverlessapplicationrepository/types"
	serviceCatalog "github.com/aws/aws-sdk-go-v2/service/servicecatalog/types"
	serviceDiscovery "github.com/aws/aws-sdk-go-v2/service/servicediscovery/types"
	servicequotas "github.com/aws/aws-sdk-go-v2/service/servicequotas/types"
	ses "github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	sfnop "github.com/aws/aws-sdk-go-v2/service/sfn"
	sfn "github.com/aws/aws-sdk-go-v2/service/sfn/types"
	shield "github.com/aws/aws-sdk-go-v2/service/shield/types"
	simspaceweaverop "github.com/aws/aws-sdk-go-v2/service/simspaceweaver"
	simspaceweaver "github.com/aws/aws-sdk-go-v2/service/simspaceweaver/types"
	sns "github.com/aws/aws-sdk-go-v2/service/sns/types"
	ssm_sdkv2 "github.com/aws/aws-sdk-go-v2/service/ssm"
	ssm "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	ssoadmin "github.com/aws/aws-sdk-go-v2/service/ssoadmin/types"
	storagegateway "github.com/aws/aws-sdk-go-v2/service/storagegateway/types"
	waf2 "github.com/aws/aws-sdk-go-v2/service/waf"
	waf "github.com/aws/aws-sdk-go-v2/service/waf/types"
	wafregional "github.com/aws/aws-sdk-go-v2/service/wafregional/types"
	wafv2op "github.com/aws/aws-sdk-go-v2/service/wafv2"
	wafv2 "github.com/aws/aws-sdk-go-v2/service/wafv2/types"
	types3 "github.com/aws/aws-sdk-go-v2/service/wellarchitected/types"
	workspaces "github.com/aws/aws-sdk-go-v2/service/workspaces/types"
)

type Metadata struct {
	Name         string
	AccountID    string
	SourceID     string
	Region       string
	Partition    string
	ResourceType string
}

//  ===================  Access Analyzer ==================

//index:aws_accessanalyzer_analyzer
//getfilter:name=description.Analyzer.Name
//listfilter:type=description.Analyzer.Type
type AccessAnalyzerAnalyzerDescription struct {
	Analyzer accessanalyzer.AnalyzerSummary
	Findings []accessanalyzer.FindingSummary
}

type AccessAnalyzerAnalyzerFindingDescription struct {
	AnalyzerArn string
	Finding     accessanalyzer.FindingSummary
}

//  ===================   ApiGateway   ===================

//index:aws_apigateway_stage
//getfilter:rest_api_id=description.RestApiId
//getfilter:name=description.Stage.StageName
type ApiGatewayStageDescription struct {
	RestApiId *string
	Stage     apigateway.Stage
}

//index:aws_apigatewayv2_stage
//getfilter:api_id=description.ApiId
//getfilter:name=description.Stage.StageName
type ApiGatewayV2StageDescription struct {
	ApiId *string
	Stage apigatewayv2.Stage
}

//index:aws_apigateway_restapi
//getfilter:api_id=description.RestAPI.Id
type ApiGatewayRestAPIDescription struct {
	RestAPI apigateway.RestApi
}

//index:aws_apigateway_apikey
//getfilter:id=description.ApiKey.Id
//listfilter:customer_id=description.ApiKey.CustomerId
type ApiGatewayApiKeyDescription struct {
	ApiKey apigateway.ApiKey
}

//index:aws_apigateway_usageplan
//getfilter:id=description.UsagePlan.Id
type ApiGatewayUsagePlanDescription struct {
	UsagePlan apigateway.UsagePlan
}

//index:aws_apigateway_authorizer
//getfilter:id=description.Authorizer.Id
//getfilter:rest_api_id=description.RestApiId
type ApiGatewayAuthorizerDescription struct {
	Authorizer apigateway.Authorizer
	RestApiId  string
}

//index:aws_apigatewayv2_api
//getfilter:api_id=description.API.ApiId
type ApiGatewayV2APIDescription struct {
	API apigatewayv2.Api
}

//index:aws_apigatewayv2_domainname
//getfilter:domain_name=description.DomainName.DomainName
type ApiGatewayV2DomainNameDescription struct {
	DomainName apigatewayv2.DomainName
}

//index:aws_apigateway_domainname
//getfilter:domain_name=description.DomainName.DomainName
type ApiGatewayDomainNameDescription struct {
	DomainName apigateway.DomainName
}

//index:aws_apigateway_domainname
//getfilter:domain_name=description.DomainName.DomainName
type ApiGatewayV2RouteDescription struct {
	Route apigatewayv2.Route
}

//index:aws_apigatewayv2_integration
//getfilter:integration_id=description.Integration.IntegrationId
//getfilter:api_id=description.ApiId
type ApiGatewayV2IntegrationDescription struct {
	Integration apigatewayv2.Integration
	ApiId       string
}

//  ===================   ElasticBeanstalk   ===================

//index:aws_elasticbeanstalk_environment
//getfilter:environment_name=description.EnvironmentDescription.EnvironmentName
type ElasticBeanstalkEnvironmentDescription struct {
	EnvironmentDescription elasticbeanstalk.EnvironmentDescription
	ManagedAction          []elasticbeanstalk.ManagedAction
	Tags                   []elasticbeanstalk.Tag
	ConfigurationSetting   []elasticbeanstalk.ConfigurationSettingsDescription
}

//index:aws_elasticbeanstalk_application
//getfilter:name=description.Application.ApplicationName
type ElasticBeanstalkApplicationDescription struct {
	Application elasticbeanstalk.ApplicationDescription
	Tags        []elasticbeanstalk.Tag
}

//index:aws_elasticbeanstalk_platform
//getfilter:platform_name=description.Platform.PlatformName
type ElasticBeanstalkPlatformDescription struct {
	Platform elasticbeanstalk.PlatformDescription
}

type ElasticBeanstalkApplicationVersionDescription struct {
	ApplicationVersion elasticbeanstalk.ApplicationVersionDescription
	Tags               []elasticbeanstalk.Tag
}

//  ===================   ElastiCache   ===================

//index:aws_elasticache_replicationgroup
//getfilter:replication_group_id=description.ReplicationGroup.ReplicationGroupId
type ElastiCacheReplicationGroupDescription struct {
	ReplicationGroup elasticache.ReplicationGroup
}

//index:aws_elasticache_cluster
//getfilter:cache_cluster_id=description.Cluster.CacheClusterId
type ElastiCacheClusterDescription struct {
	Cluster elasticache.CacheCluster
	TagList []elasticache.Tag
}

//index:aws_elasticache_parametergroup
//getfilter:cache_parameter_group_name=description.ParameterGroup.CacheParameterGroupName
type ElastiCacheParameterGroupDescription struct {
	ParameterGroup elasticache.CacheParameterGroup
}

//index:aws_elasticache_reservedcachenode
//getfilter:reserved_cache_node_id=description.ReservedCacheNode.ReservedCacheNodeId
//listfilter:cache_node_type=description.ReservedCacheNode.CacheNodeType
//listfilter:duration=description.ReservedCacheNode.Duration
//listfilter:offering_type=description.ReservedCacheNode.OfferingType
//listfilter:reserved_cache_nodes_offering_id=description.ReservedCacheNode.ReservedCacheNodesOfferingId
type ElastiCacheReservedCacheNodeDescription struct {
	ReservedCacheNode elasticache.ReservedCacheNode
}

//index:aws_elasticache_subnetgroup
//getfilter:cache_subnet_group_name=description.SubnetGroup.CacheSubnetGroupName
type ElastiCacheSubnetGroupDescription struct {
	SubnetGroup elasticache.CacheSubnetGroup
}

//  ===================   ElasticSearch   ===================

//index:aws_elasticsearch_domain
//getfilter:domain_name=description.Domain.DomainName
type ESDomainDescription struct {
	Domain es.ElasticsearchDomainStatus
	Tags   []es.Tag
}

//  ===================   EMR   ===================

//index:aws_emr_cluster
//getfilter:id=description.Cluster.Id
type EMRClusterDescription struct {
	Cluster *emr.Cluster
}

//index:aws_emr_instance
//listfilter:cluster_id=description.ClusterID
//listfilter:instance_fleet_id=description.Instance.InstanceFleetId
//listfilter:instance_group_id=description.Instance.InstanceGroupId
type EMRInstanceDescription struct {
	Instance  emr.Instance
	ClusterID string
}

//index:aws_emr_instancefleet
type EMRInstanceFleetDescription struct {
	InstanceFleet emr.InstanceFleet
	ClusterID     string
}

//index:aws_emr_instancegroup
type EMRInstanceGroupDescription struct {
	InstanceGroup emr.InstanceGroup
	ClusterID     string
}

//index:aws_emr_instancegroup
type EMRBlockPublicAccessConfigurationDescription struct {
	Configuration         emr.BlockPublicAccessConfiguration
	ConfigurationMetadata emr.BlockPublicAccessConfigurationMetadata
}

//  ===================   GuardDuty   ===================

//index:aws_guardduty_finding
type GuardDutyFindingDescription struct {
	Finding guardduty.Finding
}

//index:aws_guardduty_detector
//getfilter:detector_id=description.DetectorId
type GuardDutyDetectorDescription struct {
	DetectorId string
	Detector   *guarddutyop.GetDetectorOutput
}

//index:aws_guardduty_filter
//getfilter:name=description.Filter.Name
//getfilter:detector_id=description.DetectorId
//listfilter:detector_id=description.DetectorId
type GuardDutyFilterDescription struct {
	Filter     guarddutyop.GetFilterOutput
	DetectorId string
}

//index:aws_guardduty_ipset
//getfilter:ipset_id=description.IPSetId
//getfilter:detector_id=description.DetectorId
//listfilter:detector_id=description.DetectorId
type GuardDutyIPSetDescription struct {
	IPSet      guarddutyop.GetIPSetOutput
	IPSetId    string
	DetectorId string
}

//index:aws_guardduty_member
//getfilter:member_account_id=description.Member.AccountId
//getfilter:detector_id=description.Member.DetectorId
//listfilter:detector_id=description.Member.DetectorId
type GuardDutyMemberDescription struct {
	Member guardduty.Member
}

//index:aws_guardduty_publishingdestination
//getfilter:destination_id=description.PublishingDestination.DestinationId
//getfilter:detector_id=description.DetectorId
//listfilter:detector_id=description.DetectorId
type GuardDutyPublishingDestinationDescription struct {
	PublishingDestination guarddutyop.DescribePublishingDestinationOutput
	DetectorId            string
}

//index:aws_guardduty_threatintelset
//getfilter:threat_intel_set_id=description.ThreatIntelSetID
//getfilter:detector_id=description.DetectorId
//listfilter:detector_id=description.DetectorId
type GuardDutyThreatIntelSetDescription struct {
	ThreatIntelSet   guarddutyop.GetThreatIntelSetOutput
	DetectorId       string
	ThreatIntelSetID string
}

//  ===================   Backup   ===================

//index:aws_backup_plan
//getfilter:backup_plan_id=description.BackupPlan.BackupPlanId
type BackupPlanDescription struct {
	BackupPlan  backup.BackupPlansListMember
	PlanDetails backup.BackupPlan
}

//index:aws_backup_selection
//getfilter:backup_plan_id=description.BackupSelection.BackupPlanId
//getfilter:selection_id=description.BackupSelection.SelectionId
type BackupSelectionDescription struct {
	BackupSelection backup.BackupSelectionsListMember
	ListOfTags      []backup.Condition
	Resources       []string
}

//index:aws_backup_vault
//getfilter:name=description.BackupVault.BackupVaultName
type BackupVaultDescription struct {
	BackupVault       backup.BackupVaultListMember
	Policy            *string
	BackupVaultEvents []backup.BackupVaultEvent
	SNSTopicArn       *string
	Tags              map[string]string
}

//index:aws_backup_recoverypoint
//getfilter:backup_vault_name=description.RecoveryPoint.BackupVaultName
//getfilter:recovery_point_arn=description.RecoveryPoint.RecoveryPointArn
//listfilter:recovery_point_arn=description.RecoveryPoint.RecoveryPointArn
//listfilter:resource_type=description.RecoveryPoint.ResourceType
//listfilter:completion_date=description.RecoveryPoint.CompletionDate
type BackupRecoveryPointDescription struct {
	RecoveryPoint *backupservice.DescribeRecoveryPointOutput
	Tags          map[string]string
}

//index:aws_backup_protectedresource
//getfilter:resource_arn=description.ProtectedResource.ResourceArn
type BackupProtectedResourceDescription struct {
	ProtectedResource backup.ProtectedResource
}

//index:aws_backup_framework
//getfilter:framework_name=description.Framework.FrameworkName
type BackupFrameworkDescription struct {
	Framework backupop.DescribeFrameworkOutput
	Tags      map[string]string
}

//index:aws_backup_legalhold
//getfilter:legal_hold_id=description.Framework.LegalHoldId
type BackupLegalHoldDescription struct {
	LegalHold backupop.GetLegalHoldOutput
}

//index:aws_backup_reportplan
//getfilter:framework_name=description.Framework.FrameworkName
type BackupReportPlanDescription struct {
	ReportPlan backup.ReportPlan
	Tags       map[string]string
}

//index:aws_backup_regionsetting
//getfilter:framework_name=description.Framework.FrameworkName
type BackupRegionSettingDescription struct {
	Region                           string
	ResourceTypeManagementPreference map[string]bool
	ResourceTypeOptInPreference      map[string]bool
}

//  ===================   CloudFront   ===================

//index:aws_cloudfront_distribution
//getfilter:id=description.Distribution.Id
type CloudFrontDistributionDescription struct {
	Distribution *cloudfront.Distribution
	ETag         *string
	Tags         []cloudfront.Tag
}

//index:aws_cloudfront_streamingdistribution
type CloudFrontStreamingDistributionDescription struct {
	StreamingDistribution *cloudfront.StreamingDistribution
	ETag                  *string
	Tags                  []cloudfront.Tag
}

//index:aws_cloudfront_originaccesscontrol
//getfilter:id=description.OriginAccessControl.Id
type CloudFrontOriginAccessControlDescription struct {
	OriginAccessControl cloudfront.OriginAccessControlSummary
	Tags                []cloudfront.Tag
}

//index:aws_cloudfront_cachepolicy
//getfilter:id=description.CachePolicy.Id
type CloudFrontCachePolicyDescription struct {
	CachePolicy cloudfrontop.GetCachePolicyOutput
}

//index:aws_cloudfront_function
//getfilter:name=description.Function.FunctionSummary.Name
type CloudFrontFunctionDescription struct {
	Function cloudfrontop.DescribeFunctionOutput
}

//index:aws_cloudfront_originaccessidentity
//getfilter:id=description.OriginAccessIdentity.CloudFrontOriginAccessIdentity.Id
type CloudFrontOriginAccessIdentityDescription struct {
	OriginAccessIdentity cloudfrontop.GetCloudFrontOriginAccessIdentityOutput
}

//index:aws_cloudfront_originrequestpolicy
//getfilter:id=description.OriginRequestPolicy.OriginRequestPolicy.Id
type CloudFrontOriginRequestPolicyDescription struct {
	OriginRequestPolicy cloudfrontop.GetOriginRequestPolicyOutput
}

//index:aws_cloudfront_responseheaderspolicy
type CloudFrontResponseHeadersPolicyDescription struct {
	ResponseHeadersPolicy cloudfrontop.GetResponseHeadersPolicyOutput
}

//  ===================   CloudWatch   ===================

type CloudWatchMetricRow struct {
	// The (single) metric Dimension name
	DimensionName *string

	// The value for the (single) metric Dimension
	DimensionValue *string

	// The namespace of the metric
	Namespace *string

	// The name of the metric
	MetricName *string

	// The average of the metric values that correspond to the data point.
	Average *float64

	// The percentile statistic for the data point.
	//ExtendedStatistics map[string]*float64 `type:"map"`

	// The maximum metric value for the data point.
	Maximum *float64

	// The minimum metric value for the data point.
	Minimum *float64

	// The number of metric values that contributed to the aggregate value of this
	// data point.
	SampleCount *float64

	// The sum of the metric values for the data point.
	Sum *float64

	// The time stamp used for the data point.
	Timestamp *time.Time

	// The standard unit for the data point.
	Unit *string
}

//index:aws_cloudwatch_alarm
//getfilter:name=description.MetricAlarm.AlarmName
//listfilter:name=description.MetricAlarm.AlarmName
//listfilter:state_value=description.MetricAlarm.StateValue
type CloudWatchAlarmDescription struct {
	MetricAlarm cloudwatch.MetricAlarm
	Tags        []cloudwatch.Tag
}

//index:aws_cloudwatch_logevent
//listfilter:log_stream_name=description.LogEvent.LogStreamName
//listfilter:log_group_name=description.LogGroupName
//listfilter:timestamp=description.LogEvent.Timestamp
type CloudWatchLogEventDescription struct {
	LogEvent     cloudwatchlogs.FilteredLogEvent
	LogGroupName string
}

//index:aws_cloudwatch_logresourcepolicy
type CloudWatchLogResourcePolicyDescription struct {
	ResourcePolicy cloudwatchlogs.ResourcePolicy
}

//index:aws_cloudwatch_logstream
//getfilter:name=description.LogStream.LogStreamName
//listfilter:name=description.LogStream.LogStreamName
type CloudWatchLogStreamDescription struct {
	LogStream    cloudwatchlogs.LogStream
	LogGroupName string
}

//index:aws_cloudwatch_logsubscriptionfilter
//getfilter:name=description.SubscriptionFilter.FilterName
//getfilter:log_group_name=description.SubscriptionFilter.LogGroupName
//listfilter:name=description.SubscriptionFilter.FilterName
//listfilter:log_group_name=description.SubscriptionFilter.LogGroupName
type CloudWatchLogSubscriptionFilterDescription struct {
	SubscriptionFilter cloudwatchlogs.SubscriptionFilter
	LogGroupName       string
}

//index:aws_cloudwatch_metric
//listfilter:metric_name=description.Metric.MetricName
//listfilter:namespace=description.Metric.Namespace
type CloudWatchMetricDescription struct {
	Metric cloudwatch.Metric
}

//index:aws_cloudwatch_metricdata
//listfilter:metric_data_id=description.MetricDataQuery.Id
type CloudWatchMetricDataPointDescription struct {
	MetricDataResult cloudwatch.MetricDataResult
	MetricDataQuery  cloudwatch.MetricDataQuery
	TimeStamp        time.Time
	Value            float64
}

//index:aws_cloudwatch_loggroup
//getfilter:name=description.LogGroup.LogGroupName
//listfilter:name=description.LogGroup.LogGroupName
type CloudWatchLogsLogGroupDescription struct {
	LogGroup       cloudwatchlogs.LogGroup
	DataProtection *cloudwatchlogs2.GetDataProtectionPolicyOutput
	Tags           map[string]string
}

//index:aws_logs_metricfilter
//getfilter:name=decsription.MetricFilter.FilterName
//listfilter:name=decsription.MetricFilter.FilterName
//listfilter:log_group_name=decsription.MetricFilter.LogGroupName
//listfilter:metric_transformation_name=decsription.MetricFilter.MetricTransformations.MetricName
//listfilter:metric_transformation_namespace=decsription.MetricFilter.MetricTransformations.MetricNamespace
type CloudWatchLogsMetricFilterDescription struct {
	MetricFilter cloudwatchlogs.MetricFilter
}

//  ===================   CodeBuild   ===================

//index:aws_codebuild_project
//getfilter:name=description.Project.Name
type CodeBuildProjectDescription struct {
	Project codebuild.Project
}

//index:aws_codebuild_sourcecredential
type CodeBuildSourceCredentialDescription struct {
	SourceCredentialsInfo codebuild.SourceCredentialsInfo
}

//index:aws_codebuild_build
//getfilter:id=description.Build.Id
type CodeBuildBuildDescription struct {
	Build codebuild.Build
}

//  ===================   Config   ===================

//index:aws_config_configurationrecorder
//getfilter:name=description.ConfigurationRecorder.Name
//listfilter:name=description.ConfigurationRecorder.Name
type ConfigConfigurationRecorderDescription struct {
	ConfigurationRecorder        configservice.ConfigurationRecorder
	ConfigurationRecordersStatus configservice.ConfigurationRecorderStatus
}

//index:aws_config_aggregationauthorization
type ConfigAggregationAuthorizationDescription struct {
	AggregationAuthorization configservice.AggregationAuthorization
	Tags                     []configservice.Tag
}

//index:aws_config_conformancepack
//getfilter:name=description.ConformancePack.ConformancePackName
type ConfigConformancePackDescription struct {
	ConformancePack configservice.ConformancePackDetail
}

//index:aws_config_rule
//getfilter:name=description.Rule.ConfigRuleName
type ConfigRuleDescription struct {
	Rule       configservice.ConfigRule
	Compliance configservice.ComplianceByConfigRule
	Tags       []configservice.Tag
}

//index:aws_config_retentionconfiguration
//getfilter:name=description.ConformancePack.ConformancePackName
type ConfigRetentionConfigurationDescription struct {
	RetentionConfiguration configservice.RetentionConfiguration
}

//  ===================   Dax   ===================

//index:aws_dax_cluster
//getfilter:cluster_name=description.Cluster.ClusterName
//listfilter:cluster_name=description.Cluster.ClusterName
type DAXClusterDescription struct {
	Cluster dax.Cluster
	Tags    []dax.Tag
}

//index:aws_dax_parametergroup
//listfilter:parameter_group_name=description.ParameterGroup.ParameterGroupName
type DAXParameterGroupDescription struct {
	ParameterGroup dax.ParameterGroup
}

//index:aws_dax_parameter
//listfilter:parameter_group_name=description.ParameterGroupName
type DAXParameterDescription struct {
	Parameter          dax.Parameter
	ParameterGroupName string
}

//index:aws_dax_subnetgroup
//listfilter:subnet_group_name=description.SubnetGroup.SubnetGroupName
type DAXSubnetGroupDescription struct {
	SubnetGroup dax.SubnetGroup
}

//  ===================   Database Migration Service   ===================

//index:aws_dms_replicationinstance
//getfilter:arn=description.ReplicationInstance.ReplicationInstanceArn
//listfilter:replication_instance_identifier=description.ReplicationInstance.ReplicationInstanceIdentifier
//listfilter:arn=description.ReplicationInstance.ReplicationInstanceArn
//listfilter:replication_instance_class=description.ReplicationInstance.ReplicationInstanceClass
//listfilter:engine_version=description.ReplicationInstance.EngineVersion
type DMSReplicationInstanceDescription struct {
	ReplicationInstance dms.ReplicationInstance
	Tags                []dms.Tag
}

type DMSEndpointDescription struct {
	Endpoint dms.Endpoint
	Tags     []dms.Tag
}

type DMSReplicationTaskDescription struct {
	ReplicationTask dms.ReplicationTask
	Tags            []dms.Tag
}

//  ===================   DynamoDb   ===================

//index:aws_dynamodb_table
//getfilter:name=description.Table.TableName
//listfilter:name=description.Table.TableName
type DynamoDbTableDescription struct {
	Table                *dynamodb.TableDescription
	ContinuousBackup     *dynamodb.ContinuousBackupsDescription
	Tags                 []dynamodb.Tag
	StreamingDestination *dynamodb2.DescribeKinesisStreamingDestinationOutput
}

//index:aws_dynamodb_globalsecondaryindex
//getfilter:index_arn=description.GlobalSecondaryIndex.IndexArn
type DynamoDbGlobalSecondaryIndexDescription struct {
	GlobalSecondaryIndex dynamodb.GlobalSecondaryIndexDescription
}

//index:aws_dynamodb_localsecondaryindex
//getfilter:index_arn=description.LocalSecondaryIndex.IndexArn
type DynamoDbLocalSecondaryIndexDescription struct {
	LocalSecondaryIndex dynamodb.LocalSecondaryIndexDescription
}

//index:aws_dynamodbstreams_stream
//getfilter:stream_arn=description.Stream.StreamArn
type DynamoDbStreamDescription struct {
	Stream dynamodbstream.Stream
}

//index:aws_dynamodb_backup
//getfilter:arn=description.Backup.BackupArn
//listfilter:backup_type=description.Backup.BackupType
//listfilter:arn=description.Backup.BackupArn
//listfilter:table_name=description.Backup.TableName
type DynamoDbBackupDescription struct {
	Backup dynamodb.BackupSummary
}

//index:aws_dynamodb_globaltable
//getfilter:global_table_name=description.GlobalTable.GlobalTableName
//listfilter:global_table_name=description.GlobalTable.GlobalTableName
type DynamoDbGlobalTableDescription struct {
	GlobalTable dynamodb.GlobalTableDescription
}

//index:aws_dynamodb_tableexport
//getfilter:arn=description.Export.ExportArn
//listfilter:arn=description.Export.ExportArn
type DynamoDbTableExportDescription struct {
	Export dynamodb.ExportDescription
}

//index:aws_dynamodb_metricaccountprovisionedreadcapacityutilization
type DynamoDBMetricAccountProvisionedReadCapacityUtilizationDescription struct {
	CloudWatchMetricRow
}

//index:aws_dynamodb_metricaccountprovisionedwritecapacityutilization
type DynamoDBMetricAccountProvisionedWriteCapacityUtilizationDescription struct {
	CloudWatchMetricRow
}

//  ===================   OAM   ===================

//index:aws_oam_link
//getfilter:arn=description.Link.Arn
type OAMLinkDescription struct {
	Link *oamop.GetLinkOutput
}

//index:aws_oam_sink
//getfilter:arn=description.Sink.Arn
type OAMSinkDescription struct {
	Sink oam.ListSinksItem
	Tags map[string]string
}

//  ===================   EC2   ===================

//index:aws_ec2_volumesnapshot
//getfilter:snapshot_id=description.Snapshot.SnapshotId
//listfilter:description=description.Snapshot.Description
//listfilter:encrypted=description.Snapshot.Encrypted
//listfilter:owner_alias=description.Snapshot.OwnerAlias
//listfilter:owner_id=description.Snapshot.OwnerId
//listfilter:snapshot_id=description.Snapshot.SnapshotId
//listfilter:state=description.Snapshot.State
//listfilter:progress=description.Snapshot.Progress
//listfilter:volume_id=description.Snapshot.VolumeId
//listfilter:volume_size=description.Snapshot.VolumeSize
type EC2VolumeSnapshotDescription struct {
	Snapshot                *ec2.Snapshot
	CreateVolumePermissions []ec2.CreateVolumePermission
}

//index:aws_ec2_elasticip
type EC2ElasticIPDescription struct {
	Address ec2.Address
}

//index:aws_ec2_customergateway
//getfilter:customer_gateway_id=description.CustomerGateway.CustomerGatewayId
//listfilter:ip_address=description.CustomerGateway.IpAddress
//listfilter:bgp_asn=description.CustomerGateway.BgpAsn
//listfilter:state=description.CustomerGateway.State
//listfilter:type=description.CustomerGateway.Type
type EC2CustomerGatewayDescription struct {
	CustomerGateway ec2.CustomerGateway
}

//index:aws_ec2_verifiedaccessinstance
//listfilter:verified_access_instance_id=description.VerifiedAccountInstance.VerifiedAccessInstanceId
type EC2VerifiedAccessInstanceDescription struct {
	VerifiedAccountInstance ec2.VerifiedAccessInstance
}

//index:aws_ec2_verifiedaccessendpoint
//getfilter:verified_access_endpoint_id=description.VerifiedAccountEndpoint.VerifiedAccessEndpointId
//listfilter:verified_access_group_id=description.VerifiedAccountEndpoint.VerifiedAccessGroupId
//listfilter:verified_access_instance_id=description.VerifiedAccountEndpoint.VerifiedAccessInstanceId
type EC2VerifiedAccessEndpointDescription struct {
	VerifiedAccountEndpoint ec2.VerifiedAccessEndpoint
}

//index:aws_ec2_verifiedaccessgroup
//getfilter:verified_access_group_id=description.VerifiedAccountEndpoint.VerifiedAccessGroupId
//listfilter:verified_access_instance_id=description.VerifiedAccountGroup.VerifiedAccessInstanceId
type EC2VerifiedAccessGroupDescription struct {
	VerifiedAccountGroup ec2.VerifiedAccessGroup
}

//index:aws_ec2_verifiedaccesstrustprovider
//listfilter:verified_access_trust_provider_id=description.VerifiedAccessTrustProvider.VerifiedAccessTrustProviderId
type EC2VerifiedAccessTrustProviderDescription struct {
	VerifiedAccessTrustProvider ec2.VerifiedAccessTrustProvider
}

//index:aws_ec2_vpngateway
//getfilter:vpn_gateway_id=description.VPNGateway.VpnGatewayId
//listfilter:amazon_side_asn=description.VPNGateway.AmazonSideAsn
//listfilter:availability_zone=description.VPNGateway.AvailabilityZone
//listfilter:state=description.VPNGateway.State
//listfilter:type=description.VPNGateway.Type
type EC2VPNGatewayDescription struct {
	VPNGateway ec2.VpnGateway
}

//index:aws_ec2_volume
//getfilter:volume_id=description.Volume.VolumeId
type EC2VolumeDescription struct {
	Volume     *ec2.Volume
	Attributes struct {
		AutoEnableIO bool
		ProductCodes []ec2.ProductCode
	}
}

//index:aws_ec2_volume
//getfilter:volume_id=description.Volume.VolumeId
type EC2ClientVpnEndpointDescription struct {
	ClientVpnEndpoint ec2.ClientVpnEndpoint
}

//index:aws_ec2_instance
//getfilter:instance_id=description.Instance.InstanceId
//listfilter:hypervisor=description.Instance.Hypervisor
//listfilter:iam_instance_profile_arn=description.Instance.IamInstanceProfile.Arn
//listfilter:image_id=description.Instance.ImageId
//listfilter:instance_lifecycle=description.Instance.InstanceLifecycle
//listfilter:instance_state=description.Instance.State.Name
//listfilter:instance_type=description.Instance.InstanceType
//listfilter:monitoring_state=description.Instance.Monitoring.State
//listfilter:outpost_arn=description.Instance.OutpostArn
//listfilter:placement_availability_zone=description.Instance.Placement.AvailabilityZone
//listfilter:placement_group_name=description.Instance.Placement.GroupName
//listfilter:public_dns_name=description.Instance.PublicDnsName
//listfilter:ram_disk_id=description.Instance.RamdiskId
//listfilter:root_device_name=description.Instance.RootDeviceName
//listfilter:root_device_type=description.Instance.RootDeviceType
//listfilter:subnet_id=description.Instance.SubnetId
//listfilter:placement_tenancy=description.Instance.Placement.Tenancy
//listfilter:virtualization_type=description.Instance.VirtualizationType
//listfilter:vpc_id=description.Instance.VpcId
type EC2InstanceDescription struct {
	Instance       *ec2.Instance
	InstanceStatus *ec2.InstanceStatus
	Attributes     struct {
		InstanceInitiatedShutdownBehavior string
		DisableApiTermination             bool
	}
	LaunchTemplateData ec2.ResponseLaunchTemplateData
}

//index:aws_ec2_vpc
//getfilter:vpc_id=description.Vpc.VpcId
type EC2VpcDescription struct {
	Vpc ec2.Vpc
}

//index:aws_ec2_networkinterface
//getfilter:network_interface_id=description.NetworkInterface.NetworkInterfaceId
type EC2NetworkInterfaceDescription struct {
	NetworkInterface ec2.NetworkInterface
}

//index:aws_ec2_regionalsettings
type EC2RegionalSettingsDescription struct {
	EbsEncryptionByDefault         *bool
	KmsKeyId                       *string
	SnapshotBlockPublicAccessState ec2.SnapshotBlockPublicAccessState
}

//index:aws_ec2_ebsvolumemetricreadops
type EbsVolumeMetricReadOpsDescription struct {
	CloudWatchMetricRow
}

//index:aws_ec2_ebsvolumemetricreadopsdaily
type EbsVolumeMetricReadOpsDailyDescription struct {
	CloudWatchMetricRow
}

//index:aws_ec2_ebsvolumemetricreadopshourly
type EbsVolumeMetricReadOpsHourlyDescription struct {
	CloudWatchMetricRow
}

//index:aws_ec2_ebsvolumemetricwriteops
type EbsVolumeMetricWriteOpsDescription struct {
	CloudWatchMetricRow
}

//index:aws_ec2_ebsvolumemetricwriteopsdaily
type EbsVolumeMetricWriteOpsDailyDescription struct {
	CloudWatchMetricRow
}

//index:aws_ec2_ebsvolumemetricwriteopshourly
type EbsVolumeMetricWriteOpsHourlyDescription struct {
	CloudWatchMetricRow
}

//index:aws_ec2_subnet
//getfilter:subnet_id=description.Subnet.SubnetId
type EC2SubnetDescription struct {
	Subnet ec2.Subnet
}

//index:aws_ec2_vpcendpoint
//getfilter:vpc_endpoint_id=description.VpcEndpoint.VpcEndpointId
type EC2VPCEndpointDescription struct {
	VpcEndpoint ec2.VpcEndpoint
}

//index:aws_ec2_securitygroup
//getfilter:group_id=description.SecurityGroup.GroupId
type EC2SecurityGroupDescription struct {
	SecurityGroup ec2.SecurityGroup
}

//index:aws_ec2_eip
//getfilter:allocation_id=description.SecurityGroup.AllocationId
type EC2EIPDescription struct {
	Address ec2.Address
}

//index:aws_ec2_internetgateway
//getfilter:internet_gateway_id=description.InternetGateway.InternetGatewayId
type EC2InternetGatewayDescription struct {
	InternetGateway ec2.InternetGateway
}

//index:aws_ec2_networkacl
//getfilter:network_acl_id=description.NetworkAcl.NetworkAclId
type EC2NetworkAclDescription struct {
	NetworkAcl ec2.NetworkAcl
}

//index:aws_ec2_vpnconnection
//getfilter:vpn_connection_id=description.VpnConnection.VpnConnectionId
type EC2VPNConnectionDescription struct {
	VpnConnection ec2.VpnConnection
}

//index:aws_ec2_routetable
//getfilter:route_table_id=description.RouteTable.RouteTableId
type EC2RouteTableDescription struct {
	RouteTable ec2.RouteTable
}

//index:aws_ec2_natgateway
//getfilter:nat_gateway_id=description.NatGateway.NatGatewayId
type EC2NatGatewayDescription struct {
	NatGateway ec2.NatGateway
}

//index:aws_ec2_localgateway
//getfilter:id=description.LocalGateway.LocalGatewayId
type EC2LocalGatewayDescription struct {
	LocalGateway ec2.LocalGateway
}

//index:aws_ec2_region
//getfilter:name=description.Region.RegionName
type EC2RegionDescription struct {
	Region ec2.Region
}

//index:aws_ec2_availabilityzone
//getfilter:name=description.AvailabilityZone.ZoneName
//getfilter:region_name=description.AvailabilityZone.RegionName
//listfilter:name=description.AvailabilityZone.ZoneName
//listfilter:zone_id=description.AvailabilityZone.ZoneId
type EC2AvailabilityZoneDescription struct {
	AvailabilityZone ec2.AvailabilityZone
}

//index:aws_ec2_flowlog
//getfilter:flow_log_id=description.FlowLog.FlowLogId
type EC2FlowLogDescription struct {
	FlowLog ec2.FlowLog
}

//index:aws_ec2_capacityreservation
//getfilter:capacity_reservation_id=description.CapacityReservation.CapacityReservationId
type EC2CapacityReservationDescription struct {
	CapacityReservation ec2.CapacityReservation
}

//index:aws_ec2_keypair
//getfilter:key_name=description.KeyPair.KeyName
type EC2KeyPairDescription struct {
	KeyPair ec2.KeyPairInfo
}

//index:aws_ec2_image
//getfilter:image_id=description.AMI.ImageId
type EC2AMIDescription struct {
	AMI               ec2.Image
	LaunchPermissions ec2op.DescribeImageAttributeOutput
}

//index:aws_ec2_reservedinstance
//getfilter:reserved_instance_id=description.ReservedInstance.ReservedInstancesId
type EC2ReservedInstancesDescription struct {
	ReservedInstances   ec2.ReservedInstances
	ModificationDetails []ec2.ReservedInstancesModification
}

//index:aws_ec2_capacityreservationfleet
//getfilter:capacity_reservation_fleet_id=description.CapacityReservationFleet.CapacityReservationFleetId
type EC2CapacityReservationFleetDescription struct {
	CapacityReservationFleet ec2.CapacityReservationFleet
}

//index:aws_ec2_fleet
//getfilter:fleet_id=description.Fleet.FleetId
type EC2FleetDescription struct {
	Fleet ec2.FleetData
}

//index:aws_ec2_host
//getfilter:host_id=description.Host.HostId
type EC2HostDescription struct {
	Host ec2.Host
}

//index:aws_ec2_placementgroup
//getfilter:group_name=description.PlacementGroup.GroupName
type EC2PlacementGroupDescription struct {
	PlacementGroup ec2.PlacementGroup
}

//index:aws_ec2_transitgateway
//getfilter:transit_gateway_id=description.TransitGateway.TransitGatewayId
type EC2TransitGatewayDescription struct {
	TransitGateway ec2.TransitGateway
}

//index:aws_ec2_transitgatewayroutetable
//getfilter:transit_gateway_route_table_id=description.TransitGatewayRouteTable.TransitGatewayRouteTableId
type EC2TransitGatewayRouteTableDescription struct {
	TransitGatewayRouteTable ec2.TransitGatewayRouteTable
}

//index:aws_ec2_dhcpoptions
//getfilter:dhcp_options_id=description.DhcpOptions.DhcpOptionsId
type EC2DhcpOptionsDescription struct {
	DhcpOptions ec2.DhcpOptions
}

//index:aws_ec2_egressonlyinternetgateway
//getfilter:id=description.EgressOnlyInternetGateway.EgressOnlyInternetGatewayId
type EC2EgressOnlyInternetGatewayDescription struct {
	EgressOnlyInternetGateway ec2.EgressOnlyInternetGateway
}

//index:aws_ec2_vpcpeeringconnection
type EC2VpcPeeringConnectionDescription struct {
	VpcPeeringConnection ec2.VpcPeeringConnection
}

//index:aws_ec2_securitygrouprule
type EC2SecurityGroupRuleDescription struct {
	Group           ec2.SecurityGroup
	Permission      ec2.IpPermission
	IPRange         *ec2.IpRange
	Ipv6Range       *ec2.Ipv6Range
	UserIDGroupPair *ec2.UserIdGroupPair
	PrefixListId    *ec2.PrefixListId
	Type            string
}

//index:aws_ec2_ipampool
//getfilter:ipam_pool_id=description.IpamPool.IpamPoolId
type EC2IpamPoolDescription struct {
	IpamPool ec2.IpamPool
}

//index:aws_ec2_ipam
//getfilter:ipam_id=description.Ipam.IpamId
type EC2IpamDescription struct {
	Ipam ec2.Ipam
}

//index:aws_ec2_vpcendpointservice
//getfilter:service_name=description.VPCEndpoint.ServiceName
type EC2VPCEndpointServiceDescription struct {
	VpcEndpointService     ec2.ServiceDetail
	AllowedPrincipals      []ec2.AllowedPrincipal
	VpcEndpointConnections []ec2.VpcEndpointConnection
}

//index:aws_ec2_instanceavailability
//listfilter:instance_type=description.InstanceAvailability.InstanceType
type EC2InstanceAvailabilityDescription struct {
	InstanceAvailability ec2.InstanceTypeOffering
}

//index:aws_ec2_instancetype
//getfilter:instance_type=description.InstanceType.InstanceType
type EC2InstanceTypeDescription struct {
	InstanceType ec2.InstanceTypeInfo
}

//index:aws_ec2_managedprefixlist
//listfilter:name=description.ManagedPrefixList.PrefixListName
//listfilter:id=description.ManagedPrefixList.PrefixListId
//listfilter:owner_id=description.ManagedPrefixList.OwnerId
type EC2ManagedPrefixListDescription struct {
	ManagedPrefixList ec2.ManagedPrefixList
}

//index:aws_ec2_managedprefixlistentry
type EC2ManagedPrefixListEntryDescription struct {
	PrefixListEntry ec2.PrefixListEntry
	PrefixListId    string
}

//index:aws_ec2_spotprice
//listfilter:availability_zone=description.SpotPrice.AvailabilityZone
//listfilter:instance_type=description.SpotPrice.InstanceType
//listfilter:product_description=description.SpotPrice.ProductDescription
type EC2SpotPriceDescription struct {
	SpotPrice ec2.SpotPrice
}

//index:aws_ec2_transitgatewayroute
//listfilter:prefix_list_id=description.TransitGatewayRoute.PrefixListId
//listfilter:state=description.TransitGatewayRoute.State
//listfilter:type=description.TransitGatewayRoute.Type
type EC2TransitGatewayRouteDescription struct {
	TransitGatewayRoute        ec2.TransitGatewayRoute
	TransitGatewayRouteTableId string
}

//index:aws_ec2_transitgatewayvpcattachment
//getfilter:transit_gateway_attachment_id=description.TransitGatewayAttachment.TransitGatewayAttachmentId
//listfilter:association_state=description.TransitGatewayAttachment.Association.State
//listfilter:association_transit_gateway_route_table_id=description.TransitGatewayAttachment.Association.TransitGatewayRouteTableId
//listfilter:resource_id=description.TransitGatewayAttachment.ResourceId
//listfilter:resource_owner_id=description.TransitGatewayAttachment.ResourceOwnerId
//listfilter:resource_type=description.TransitGatewayAttachment.ResourceType
//listfilter:state=description.TransitGatewayAttachment.State
//listfilter:transit_gateway_id=description.TransitGatewayAttachment.TransitGatewayId
//listfilter:transit_gateway_owner_id=description.TransitGatewayAttachment.TransitGatewayOwnerId
type EC2TransitGatewayAttachmentDescription struct {
	TransitGatewayAttachment ec2.TransitGatewayAttachment
}

type EC2LaunchTemplateDescription struct {
	LaunchTemplate ec2.LaunchTemplate
}

type EC2LaunchTemplateVersionDescription struct {
	LaunchTemplateVersion ec2.LaunchTemplateVersion
}

//index:aws_vpc_nat_gateway_metric_bytes_out_to_destination
type EC2NatGatewayMetricBytesOutToDestinationDescription struct {
	NatGateway ec2.NatGateway
}

//index:aws_vpc_eip_address_transfer
type EC2EIPAddressTransferDescription struct {
	AddressTransfer ec2.AddressTransfer
}

type EC2InstanceMetricCpuUtilizationHourlyDescription struct {
	InstanceId  *string
	Timestamp   *time.Time
	Average     *float64
	Sum         *float64
	Maximum     *float64
	Minimum     *float64
	SampleCount *float64
}

// ===================  STS Caller  =================

//index:aws_stscaller_identity
type STSCallerIdentityDescription struct {
	UsrId   string
	Account string
	Arn     string
}

//  ===================  Elastic Load Balancing  ===================

//index:aws_elasticloadbalancingv2_sslpolicy
//getfilter:name=description.SslPolicy.Name
//getfilter:region=metadata.Region
type ElasticLoadBalancingV2SslPolicyDescription struct {
	SslPolicy elasticloadbalancingv2.SslPolicy
}

//index:aws_elasticloadbalancingv2_targetgroup
//getfilter:target_group_arn=description.TargetGroup.TargetGroupArn
//listfilter:target_group_name=description.TargetGroup.TargetGroupName
type ElasticLoadBalancingV2TargetGroupDescription struct {
	TargetGroup elasticloadbalancingv2.TargetGroup
	Health      []elasticloadbalancingv2.TargetHealthDescription
	Tags        []elasticloadbalancingv2.Tag
}

//index:aws_elasticloadbalancingv2_loadbalancer
//getfilter:arn=description.LoadBalancer.LoadBalancerArn
//getfilter:type=description.LoadBalancer.Type
//listfilter:type=description.LoadBalancer.Type
type ElasticLoadBalancingV2LoadBalancerDescription struct {
	LoadBalancer elasticloadbalancingv2.LoadBalancer
	Attributes   []elasticloadbalancingv2.LoadBalancerAttribute
	Tags         []elasticloadbalancingv2.Tag
}

//index:aws_elasticloadbalancing_loadbalancer
//getfilter:name=description.LoadBalancer.LoadBalancerName
type ElasticLoadBalancingLoadBalancerDescription struct {
	LoadBalancer elasticloadbalancing.LoadBalancerDescription
	Attributes   *elasticloadbalancing.LoadBalancerAttributes
	Tags         []elasticloadbalancing.Tag
}

//index:aws_elasticloadbalancingv2_listener
//getfilter:arn=description.Listener.ListenerArn
type ElasticLoadBalancingV2ListenerDescription struct {
	Listener elasticloadbalancingv2.Listener
}

//index:aws_elasticloadbalancingv2_rule
//getfilter:arn=description.Rule.RuleArn
type ElasticLoadBalancingV2RuleDescription struct {
	Rule elasticloadbalancingv2.Rule
}

//index:aws_elasticloadbalancingv2_applicationloadbalancermetricrequestcount
type ApplicationLoadBalancerMetricRequestCountDescription struct {
	CloudWatchMetricRow
}

//index:aws_elasticloadbalancingv2_applicationloadbalancermetricrequestcountdaily
type ApplicationLoadBalancerMetricRequestCountDailyDescription struct {
	CloudWatchMetricRow
}

//index:aws_elasticloadbalancingv2_networkloadbalancermetricnetflowcount
type NetworkLoadBalancerMetricNetFlowCountDescription struct {
	CloudWatchMetricRow
}

//index:aws_elasticloadbalancingv2_networkloadbalancermetricnetflowcountdaily
type NetworkLoadBalancerMetricNetFlowCountDailyDescription struct {
	CloudWatchMetricRow
}

//  ===================  FSX  ===================

//index:aws_fsx_filesystem
//getfilter:file_system_id=description.FileSystem.FileSystemId
type FSXFileSystemDescription struct {
	FileSystem fsx.FileSystem
}

//index:aws_fsx_storagevirtualmachine
//getfilter:storage_virtual_machine_id=description.StorageVirtualMachine.StorageVirtualMachineId
type FSXStorageVirtualMachineDescription struct {
	StorageVirtualMachine fsx.StorageVirtualMachine
}

//index:aws_fsx_task
//getfilter:task_id=description.Task.TaskId
type FSXTaskDescription struct {
	Task fsx.DataRepositoryTask
}

//index:aws_fsx_volume
//getfilter:volume_id=description.Volume.VolumeId
type FSXVolumeDescription struct {
	Volume fsx.Volume
}

//index:aws_fsx_snapshot
//getfilter:snapshot_id=description.Snapshot.SnapshotId
type FSXSnapshotDescription struct {
	Snapshot fsx.Snapshot
}

//  ===================  Application Auto Scaling  ===================

//index:aws_applicationautoscaling_target
//getfilter:service_namespace=description.ScalableTarget.ServiceNamespace
//getfilter:resource_id=description.ScalableTarget.ResourceId
//listfilter:service_namespace=description.ScalableTarget.ServiceNamespace
//listfilter:resource_id=description.ScalableTarget.ResourceId
//listfilter:scalable_dimension=description.ScalableTarget.ScalableDimension
type ApplicationAutoScalingTargetDescription struct {
	ScalableTarget applicationautoscaling.ScalableTarget
}

//index:aws_applicationautoscaling_target
//getfilter:service_namespace=description.ScalablePolicy.ServiceNamespace
//getfilter:resource_id=description.ScalablePolicy.ResourceId
//listfilter:service_namespace=description.ScalablePolicy.ServiceNamespace
//listfilter:resource_id=description.ScalablePolicy.ResourceId
//listfilter:scalable_dimension=description.ScalablePolicy.ScalableDimension
type ApplicationAutoScalingPolicyDescription struct {
	ScalablePolicy applicationautoscaling.ScalingPolicy
}

//  ===================  Auto Scaling  ===================

//index:aws_autoscaling_autoscalinggroup
//getfilter:name=description.AutoScalingGroup.AutoScalingGroupName
type AutoScalingGroupDescription struct {
	AutoScalingGroup *autoscaling.AutoScalingGroup
	Policies         []autoscaling.ScalingPolicy
}

//index:aws_autoscaling_launchconfiguration
//getfilter:name=description.LaunchConfiguration.LaunchConfigurationName
type AutoScalingLaunchConfigurationDescription struct {
	LaunchConfiguration autoscaling.LaunchConfiguration
}

// ======================== ACM ==========================

//index:aws_certificatemanager_certificate
//getfilter:certificate_arn=description.Certificate.CertificateArn
//listfilter:status=description.Certificate.Status
type CertificateManagerCertificateDescription struct {
	Certificate acm.CertificateDetail
	Attributes  struct {
		Certificate      *string
		CertificateChain *string
	}
	Tags []acm.Tag
}

// =====================  CloudTrail  =====================

//index:aws_cloudtrail_trail
//getfilter:name=description.Trail.Name
//getfilter:arn=description.Trail.TrailARN
type CloudTrailTrailDescription struct {
	Trail                  cloudtrailtypes.Trail
	TrailStatus            cloudtrail.GetTrailStatusOutput
	EventSelectors         []cloudtrailtypes.EventSelector
	AdvancedEventSelectors []cloudtrailtypes.AdvancedEventSelector
	Tags                   []cloudtrailtypes.Tag
}

//index:aws_cloudtrail_channel
//getfilter:arn=description.Channel.ChannelArn
type CloudTrailChannelDescription struct {
	Channel cloudtrail.GetChannelOutput
}

//index:aws_cloudtrail_eventdatastore
//getfilter:arn=description.EventDataStore.EventDataStoreArn
type CloudTrailEventDataStoreDescription struct {
	EventDataStore cloudtrail.GetEventDataStoreOutput
}

//index:aws_cloudtrail_import
//getfilter:import_id=description.Import.ImportId
//listfilter:import_status=description.Import.ImportStatus
type CloudTrailImportDescription struct {
	Import cloudtrail.GetImportOutput
}

//index:aws_cloudtrail_query
//getfilter:event_data_store_arn=description.EventDataStoreARN
//getfilter:query_id=description.Query.QueryId
//listfilter:event_data_store_arn=description.EventDataStoreARN
//listfilter:query_status=description.Query.QueryStatus
//listfilter:creation_time=description.Query.QueryStatistics.CreationTime
type CloudTrailQueryDescription struct {
	Query             cloudtrail.DescribeQueryOutput
	EventDataStoreARN string
}

//index:aws_cloudtrail_trailevent
//listfilter:log_stream_name=description.TrailEvent.LogStreamName
//listfilter:timestamp=description.TrailEvent.Timestamp
type CloudTrailTrailEventDescription struct {
	TrailEvent   cloudwatchlogs.FilteredLogEvent
	LogGroupName string
}

// ====================== IAM =========================

//index:aws_iam_account
type IAMAccountDescription struct {
	Aliases      []string
	Account      *organizations.Account
	Organization *organizations.Organization
}

//index:aws_iam_access_advisor
type IAMAccessAdvisorDescription struct {
	PrincipalARN        string
	ServiceLastAccessed iam.ServiceLastAccessed
}

type AccountSummary struct {
	AccountMFAEnabled                 int32
	AccessKeysPerUserQuota            int32
	AccountAccessKeysPresent          int32
	AccountSigningCertificatesPresent int32
	AssumeRolePolicySizeQuota         int32
	AttachedPoliciesPerGroupQuota     int32
	AttachedPoliciesPerRoleQuota      int32
	AttachedPoliciesPerUserQuota      int32
	GlobalEndpointTokenVersion        int32
	GroupPolicySizeQuota              int32
	Groups                            int32
	GroupsPerUserQuota                int32
	GroupsQuota                       int32
	InstanceProfiles                  int32
	InstanceProfilesQuota             int32
	MFADevices                        int32
	MFADevicesInUse                   int32
	Policies                          int32
	PoliciesQuota                     int32
	PolicySizeQuota                   int32
	PolicyVersionsInUse               int32
	PolicyVersionsInUseQuota          int32
	Providers                         int32
	RolePolicySizeQuota               int32
	Roles                             int32
	RolesQuota                        int32
	ServerCertificates                int32
	ServerCertificatesQuota           int32
	SigningCertificatesPerUserQuota   int32
	UserPolicySizeQuota               int32
	Users                             int32
	UsersQuota                        int32
	VersionsPerPolicyQuota            int32
}

//index:aws_iam_accountsummary
type IAMAccountSummaryDescription struct {
	AccountSummary AccountSummary
}

//index:aws_iam_accesskey
type IAMAccessKeyDescription struct {
	AccessKey         iam.AccessKeyMetadata
	AccessKeyLastUsed *iam.AccessKeyLastUsed
}

//index:aws_iam_sshpublickey
type IAMSSHPublicKeyDescription struct {
	SSHPublicKeyKey iam.SSHPublicKeyMetadata
}

//index:aws_iam_accountpasswordpolicy
type IAMAccountPasswordPolicyDescription struct {
	PasswordPolicy iam.PasswordPolicy
}

type InlinePolicy struct {
	PolicyName     string
	PolicyDocument string
}

//index:aws_iam_user
//getfilter:name=description.User.UserName
//getfilter:arn=description.User.Arn
type IAMUserDescription struct {
	User               iam.User
	Groups             []iam.Group
	LoginProfile       iam.LoginProfile
	InlinePolicies     []InlinePolicy
	AttachedPolicyArns []string
	MFADevices         []iam.MFADevice
}

//index:aws_iam_group
//getfilter:name=description.Group.GroupName
//getfilter:arn=description.Group.Arn
type IAMGroupDescription struct {
	Group              iam.Group
	Users              []iam.User
	InlinePolicies     []InlinePolicy
	AttachedPolicyArns []string
}

//index:aws_iam_role
//getfilter:name=description.Role.RoleName
//getfilter:arn=description.Role.Arn
type IAMRoleDescription struct {
	Role                iam.Role
	InstanceProfileArns []string
	InlinePolicies      []InlinePolicy
	AttachedPolicyArns  []string
}

//index:aws_iam_servercertificate
//getfilter:name=description.ServerCertificate.ServerCertificateMetadata.ServerCertificateName
type IAMServerCertificateDescription struct {
	ServerCertificate iam.ServerCertificate
	BodyLength        int
}

//index:aws_iam_policy
//getfilter:arn=description.Policy.Arn
type IAMPolicyDescription struct {
	Policy        iam.Policy
	PolicyVersion iam.PolicyVersion
}

type CredentialReport struct {
	GeneratedTime             *time.Time `csv:"-"`
	UserArn                   string     `csv:"arn"`
	UserName                  string     `csv:"user"`
	UserCreationTime          string     `csv:"user_creation_time"`
	AccessKey1Active          bool       `csv:"access_key_1_active"`
	AccessKey1LastRotated     string     `csv:"access_key_1_last_rotated"`
	AccessKey1LastUsedDate    string     `csv:"access_key_1_last_used_date"`
	AccessKey1LastUsedRegion  string     `csv:"access_key_1_last_used_region"`
	AccessKey1LastUsedService string     `csv:"access_key_1_last_used_service"`
	AccessKey2Active          bool       `csv:"access_key_2_active"`
	AccessKey2LastRotated     string     `csv:"access_key_2_last_rotated"`
	AccessKey2LastUsedDate    string     `csv:"access_key_2_last_used_date"`
	AccessKey2LastUsedRegion  string     `csv:"access_key_2_last_used_region"`
	AccessKey2LastUsedService string     `csv:"access_key_2_last_used_service"`
	Cert1Active               bool       `csv:"cert_1_active"`
	Cert1LastRotated          string     `csv:"cert_1_last_rotated"`
	Cert2Active               bool       `csv:"cert_2_active"`
	Cert2LastRotated          string     `csv:"cert_2_last_rotated"`
	MFAActive                 bool       `csv:"mfa_active"`
	PasswordEnabled           string     `csv:"password_enabled"`
	PasswordLastChanged       string     `csv:"password_last_changed"`
	PasswordLastUsed          string     `csv:"password_last_used"`
	PasswordNextRotation      string     `csv:"password_next_rotation"`
}

//index:aws_iam_credentialreport
type IAMCredentialReportDescription struct {
	CredentialReport CredentialReport
}

//index:aws_iam_virtualmfadevices
type IAMVirtualMFADeviceDescription struct {
	VirtualMFADevice iam.VirtualMFADevice
	Tags             []iam.Tag
}

//index:aws_iam_policyattachment
//getfilter:is_attached=description.IsAttached
type IAMPolicyAttachmentDescription struct {
	PolicyArn             string
	PolicyAttachmentCount int32
	IsAttached            bool
	PolicyGroups          []iam.PolicyGroup
	PolicyRoles           []iam.PolicyRole
	PolicyUsers           []iam.PolicyUser
}

//index:aws_iam_samlprovider
//getfilter:arn=ARN
type IAMSamlProviderDescription struct {
	SamlProvider iamop.GetSAMLProviderOutput
}

//index:aws_iam_servicespecificcredential
//listfilter:service_name=description.ServiceSpecificCredential.ServiceName
//listfilter:user_name=description.ServiceSpecificCredential.UserName
type IAMServiceSpecificCredentialDescription struct {
	ServiceSpecificCredential iam.ServiceSpecificCredentialMetadata
}

//index:aws_iam_openidconnectprovider
type IAMOpenIdConnectProviderDescription struct {
	ClientIDList   []string
	Tags           []iam.Tag
	CreateDate     time.Time
	ThumbprintList []string
	URL            string
}

//  ===================  RDS  ===================

//index:aws_rds_dbcluster
//getfilter:db_cluster_identifier=description.DBCluster.DBClusterIdentifier
type RDSDBClusterDescription struct {
	DBCluster                 rds.DBCluster
	PendingMaintenanceActions []rds.ResourcePendingMaintenanceActions
}

//index:aws_rds_dbclusterparametergroup
//getfilter:name=description.DBClusterParameterGroup.DBClusterParameterGroupName
type RDSDBClusterParameterGroupDescription struct {
	DBClusterParameterGroup rds.DBClusterParameterGroup
	Parameters              []rds.Parameter
	Tags                    []rds.Tag
}

//index:aws_rds_optiongroup
//getfilter:name=description.OptionGroup.OptionGroupName
//listfilter:engine_name=description.OptionGroup.EngineName
//listfilter:major_engine_version=description.OptionGroup.MajorEngineVersion
type RDSOptionGroupDescription struct {
	OptionGroup rds.OptionGroup
	Tags        *rds_sdkv2.ListTagsForResourceOutput
}

//index:aws_rds_dbparametergroup
//getfilter:name=description.DBParameterGroup.DBParameterGroupName
type RDSDBParameterGroupDescription struct {
	DBParameterGroup rds.DBParameterGroup
	Parameters       []rds.Parameter
	Tags             []rds.Tag
}

//index:aws_rds_dbproxy
//getfilter:db_proxy_name=description.DBProxy.DBProxyName
type RDSDBProxyDescription struct {
	DBProxy rds.DBProxy
	Tags    []rds.Tag
}

//index:aws_rds_dbsubnetgroup
//getfilter:name=description.DBSubnetGroup.DBSubnetGroupName
type RDSDBSubnetGroupDescription struct {
	DBSubnetGroup rds.DBSubnetGroup
	Tags          *rds_sdkv2.ListTagsForResourceOutput
}

//index:aws_rds_dbclustersnapshot
//getfilter:db_cluster_snapshot_identifier=description.DBClusterSnapshot.DBClusterIdentifier
//listfilter:db_cluster_identifier=description.DBClusterSnapshot.DBClusterIdentifier
//listfilter:db_cluster_snapshot_identifier=description.DBClusterSnapshot.DBClusterSnapshotIdentifier
//listfilter:engine=description.DBClusterSnapshot.Engine
//listfilter:type=description.DBClusterSnapshot.SnapshotType
type RDSDBClusterSnapshotDescription struct {
	DBClusterSnapshot rds.DBClusterSnapshot
	Attributes        *rds.DBClusterSnapshotAttributesResult
}

//index:aws_rds_eventsubscription
//getfilter:cust_subscription_id=description.EventSubscription.CustSubscriptionId
type RDSDBEventSubscriptionDescription struct {
	EventSubscription rds.EventSubscription
}

//index:aws_rds_dbinstance
//getfilter:db_instance_identifier=description.DBInstance.DBInstanceIdentifier
type RDSDBInstanceDescription struct {
	DBInstance         rds.DBInstance
	PendingMaintenance []rds.ResourcePendingMaintenanceActions
	Certificate        []rds.Certificate
}

//index:aws_rds_dbsnapshot
//getfilter:db_snapshot_identifier=description.DBSnapshot.DBInstanceIdentifier
type RDSDBSnapshotDescription struct {
	DBSnapshot           rds.DBSnapshot
	DBSnapshotAttributes []rds.DBSnapshotAttribute
}

//index:aws_rds_globalcluster
//getfilter:global_cluster_identifier=description.DBGlobalCluster.GlobalClusterIdentifier
type RDSGlobalClusterDescription struct {
	GlobalCluster rds.GlobalCluster
	Tags          []rds.Tag
}

//index:aws_rds_reserveddbinstance
//getfilter:reserved_db_instance_id=description.ReservedDBInstance.ReservedDBInstanceId
//listfilter:class=description.ReservedDBInstance.DBInstanceClass
//listfilter:duration=description.ReservedDBInstance.Duration
//listfilter:lease_id=description.ReservedDBInstance.LeaseId
//listfilter:multi_az=description.ReservedDBInstance.MultiAZ
//listfilter:offering_type=description.ReservedDBInstance.OfferingType
//listfilter:reserved_db_instances_offering_id=description.ReservedDBInstance.ReservedDBInstancesOfferingId
type RDSReservedDBInstanceDescription struct {
	ReservedDBInstance rds.ReservedDBInstance
}

//index:aws_rds_dbcluster
//getfilter:db_cluster_identifier=description.DBCluster.DBClusterIdentifier
type RDSDBInstanceAutomatedBackupDescription struct {
	InstanceAutomatedBackup rds.DBInstanceAutomatedBackup
}

type RDSDBEngineVersionDescription struct {
	EngineVersion rds.DBEngineVersion
}

type RDSDBRecommendationDescription struct {
	DBRecommendation rds.DBRecommendation
}

//  ===================  Redshift  ===================

//index:aws_redshift_cluster
//getfilter:cluster_identifier=description.Cluster
type RedshiftClusterDescription struct {
	Cluster          redshift.Cluster
	LoggingStatus    *redshiftop.DescribeLoggingStatusOutput
	ScheduledActions []redshift.ScheduledAction
}

//index:aws_redshift_eventsubscription
//getfilter:cust_subscription_id=description.EventSubscription.CustSubscriptionId
type RedshiftEventSubscriptionDescription struct {
	EventSubscription redshift.EventSubscription
}

//index:aws_redshiftserverless_workgroup
//getfilter:workgroup_name=description.Workgroup.WorkgroupName
type RedshiftServerlessWorkgroupDescription struct {
	Workgroup redshiftserverlesstypes.Workgroup
	Tags      []redshiftserverlesstypes.Tag
}

//index:aws_redshift_clusterparametergroup
//getfilter:name=description.ClusterParameterGroup.ParameterGroupName
type RedshiftClusterParameterGroupDescription struct {
	ClusterParameterGroup redshift.ClusterParameterGroup
	Parameters            []redshift.Parameter
}

//index:aws_redshift_snapshot
//getfilter:snapshot_identifier=description.Snapshot.SnapshotIdentifier
type RedshiftSnapshotDescription struct {
	Snapshot redshift.Snapshot
}

//index:aws_redshiftserverless_namespace
//getfilter:namespace_name=description.Namespace.NamespaceName
type RedshiftServerlessNamespaceDescription struct {
	Namespace redshiftserverlesstypes.Namespace
	Tags      []redshiftserverlesstypes.Tag
}

//index:aws_redshiftserverless_snapshot
//getfilter:snapshot_name=description.Snapshot.SnapshotName
type RedshiftServerlessSnapshotDescription struct {
	Snapshot redshiftserverlesstypes.Snapshot
	Tags     []redshiftserverlesstypes.Tag
}

//index:aws_redshift_subnetgroup
//getfilter:cluster_subnet_group_name=description.ClusterSubnetGroup.ClusterSubnetGroupName
type RedshiftSubnetGroupDescription struct {
	ClusterSubnetGroup redshift.ClusterSubnetGroup
}

//  ===================  SNS  ===================

//index:aws_sns_topic
//getfilter:topic_arn=description.Attributes.TopicArn
type SNSTopicDescription struct {
	Attributes map[string]string
	Tags       []sns.Tag
}

//index:aws_sns_subscription
//getfilter:subscription_arn=description.Subscription.SubscriptionArn
type SNSSubscriptionDescription struct {
	Subscription sns.Subscription
	Attributes   map[string]string
}

//  ===================  SQS  ===================

//index:aws_sqs_queue
//getfilter:queue_url=description.Attributes.QueueUrl
type SQSQueueDescription struct {
	Attributes map[string]string
	Tags       map[string]string
}

//  ===================  S3  ===================

//index:aws_s3_bucket
//getfilter:name=description.Bucket.Name
type S3BucketDescription struct {
	Bucket    s3.Bucket
	BucketAcl struct {
		Grants []s3.Grant
		Owner  *s3.Owner
	}
	BucketNotification             []s3.QueueConfiguration
	Policy                         *string
	PolicyStatus                   *s3.PolicyStatus
	PublicAccessBlockConfiguration *s3.PublicAccessBlockConfiguration
	Versioning                     struct {
		MFADelete s3.MFADeleteStatus
		Status    s3.BucketVersioningStatus
	}
	LifecycleRules                    string
	LoggingEnabled                    *s3.LoggingEnabled
	ServerSideEncryptionConfiguration *s3.ServerSideEncryptionConfiguration
	ObjectLockConfiguration           *s3.ObjectLockConfiguration
	ReplicationConfiguration          *s3.ReplicationConfiguration
	Tags                              []s3.Tag
	Region                            string
	BucketWebsite                     *s3types.GetBucketWebsiteOutput
	BucketOwnershipControls           *s3types.GetBucketOwnershipControlsOutput
	EventNotificationConfiguration    *s3types.GetBucketNotificationConfigurationOutput
}

//index:aws_s3_accountsettingdescription
type S3AccountSettingDescription struct {
	PublicAccessBlockConfiguration s3control.PublicAccessBlockConfiguration
}

//index:aws_s3_object
type S3ObjectDescription struct {
	BucketName       *string
	Object           *s3types.GetObjectOutput
	ObjectSummary    s3.Object
	ObjectAttributes s3types.GetObjectAttributesOutput
	ObjectAcl        s3types.GetObjectAclOutput
	ObjectTags       s3types.GetObjectTaggingOutput
}

//index:aws_s3_bucketintelligenttieringconfiguration
type S3BucketIntelligentTieringConfigurationDescription struct {
	BucketName                      string
	IntelligentTieringConfiguration s3.IntelligentTieringConfiguration
}

//index:aws_s3_bucketintelligenttieringconfiguration
type S3MultiRegionAccessPointDescription struct {
	Report s3control.MultiRegionAccessPointReport
}

//  ===================  SageMaker  ===================

//index:aws_sagemaker_endpointconfiguration
//getfilter:name=description.EndpointConfig.EndpointConfigName
type SageMakerEndpointConfigurationDescription struct {
	EndpointConfig *sagemakerop.DescribeEndpointConfigOutput
	Tags           []sagemaker.Tag
}

//index:aws_sagemaker_app
//getfilter:name=description.DescribeAppOutput.AppName
//getfilter:app_type=description.DescribeAppOutput.AppType
//getfilter:domain_id=description.DescribeAppOutput.DomainId
//getfilter:user_profile_name=description.DescribeAppOutput.UserProfileName
//listfilter:domain_id=description.DescribeAppOutput.DomainId
//listfilter:user_profile_name=description.DescribeAppOutput.UserProfileName
type SageMakerAppDescription struct {
	AppDetails        sagemaker.AppDetails
	DescribeAppOutput *sagemakerop.DescribeAppOutput
}

//index:aws_sagemaker_domain
//getfilter:id=description.Domain.DomainId
type SageMakerDomainDescription struct {
	Domain     *sagemakerop.DescribeDomainOutput
	DomainItem sagemaker.DomainDetails
	Tags       []sagemaker.Tag
}

//index:aws_sagemaker_notebookinstance
//getfilter:name=description.NotebookInstance.NotebookInstanceName
type SageMakerNotebookInstanceDescription struct {
	NotebookInstance *sagemakerop.DescribeNotebookInstanceOutput
	Tags             []sagemaker.Tag
}

//index:aws_sagemaker_model
//getfilter:name=description.Model.ModelName
//listfilter:creation_time=description.Model.CreationTime
type SageMakerModelDescription struct {
	Model *sagemakerop.DescribeModelOutput
	Tags  []sagemaker.Tag
}

//index:aws_sagemaker_trainingjob
//getfilter:name=description.TrainingJob.Name
//listfilter:creation_time=description.TrainingJob.CreationTime
//listfilter:last_modified_time=description.TrainingJob.LastModifiedTime
//listfilter:training_job_status=description.TrainingJob.TrainingJobStatus
type SageMakerTrainingJobDescription struct {
	TrainingJob *sagemakerop.DescribeTrainingJobOutput
	Tags        []sagemaker.Tag
}

//  ===================  SecretsManager  ===================

//index:aws_secretsmanager_secret
//getfilter:arn=description.Secret.ARN
type SecretsManagerSecretDescription struct {
	Secret         *secretsmanager.DescribeSecretOutput
	ResourcePolicy *string
}

//  ===================  SecurityHub  ===================

//index:aws_securityhub_hub
//getfilter:hub_arn=description.Hub.HubArn
type SecurityHubHubDescription struct {
	Hub                  *securityhubop.DescribeHubOutput
	AdministratorAccount securityhub.Invitation
	Tags                 map[string]string
}

//index:aws_securityhub_actiontarget
//getfilter:arn=description.ActionTarget.ActionTargetArn
type SecurityHubActionTargetDescription struct {
	ActionTarget securityhub.ActionTarget
}

//index:aws_securityhub_finding
//getfilter:id=description.Finding.Id
//listfilter:company_name=description.Finding.CompanyName
//listfilter:compliance_status=description.Finding.Compliance.Status
//listfilter:confidence=description.Finding.Confidence
//listfilter:criticality=description.Finding.Criticality
//listfilter:generator_id=description.Finding.GeneratorId
//listfilter:product_arn=description.Finding.ProductArn
//listfilter:product_name=description.Finding.ProductName
//listfilter:record_state=description.Finding.RecordState
//listfilter:title=description.Finding.Title
//listfilter:verification_state=description.Finding.VerificationState
//listfilter:workflow_state=description.Finding.WorkflowState
//listfilter:workflow_status=description.Finding.Workflow.Status
type SecurityHubFindingDescription struct {
	Finding securityhub.AwsSecurityFinding
}

//index:aws_securityhub_findingaggregator
//getfilter:arn=description.FindingAggregator.FindingAggregatorArn
type SecurityHubFindingAggregatorDescription struct {
	FindingAggregator securityhubop.GetFindingAggregatorOutput
}

//index:aws_securityhub_insight
//getfilter:arn=description.Insight.InsightArn
type SecurityHubInsightDescription struct {
	Insight securityhub.Insight
}

//index:aws_securityhub_member
type SecurityHubMemberDescription struct {
	Member securityhub.Member
}

//index:aws_securityhub_product
//getfilter:product_arn=description.Product.ProductArn
type SecurityHubProductDescription struct {
	Product securityhub.Product
}

//index:aws_securityhub_standardscontrol
type SecurityHubStandardsControlDescription struct {
	StandardsControl securityhub.StandardsControl
}

//index:aws_securityhub_standardssubscription
type SecurityHubStandardsSubscriptionDescription struct {
	Standard              securityhub.Standard
	StandardsSubscription securityhub.StandardsSubscription
}

//  ===================  SSM  ===================

//index:aws_ssm_managedinstance
type SSMManagedInstanceDescription struct {
	InstanceInformation ssm.InstanceInformation
}

//index:aws_ssm_association
//getfilter:association_id=description.AssociationItem.AssociationId
//listfilter:association_name=description.AssociationItem.AssociationName
//listfilter:instance_id=description.AssociationItem.InstanceId
//listfilter:status=description.Association.AssociationDescription.Status
//listfilter:last_execution_date=description.Association.AssociationDescription.LastExecutionDate
type SSMAssociationDescription struct {
	AssociationItem ssm.Association
	Association     *ssm_sdkv2.DescribeAssociationOutput
}

//index:aws_ssm_document
//getfilter:name=description.DocumentIdentifier.Name
//listfilter:document_type=description.DocumentIdentifier.DocumentType
//listfilter:owner_type=description.DocumentIdentifier.Owner
type SSMDocumentDescription struct {
	DocumentIdentifier ssm.DocumentIdentifier
	Document           *ssm_sdkv2.DescribeDocumentOutput
	Permissions        *ssm_sdkv2.DescribeDocumentPermissionOutput
}

//index:aws_ssm_document_permission
type SSMDocumentPermissionDescription struct {
	Permissions *ssm_sdkv2.DescribeDocumentPermissionOutput
	Document    *ssm_sdkv2.DescribeDocumentOutput
}

//index:aws_ssm_inventory
//listfilter:id=description.Id
//listfilter:type_name=description.TypeName
type SSMInventoryDescription struct {
	CaptureTime   *string
	Content       interface{}
	Id            *string
	SchemaVersion *string
	TypeName      *string
	Schemas       []ssm.InventoryItemSchema
}

//index:aws_ssm_inventory_entry
//listfilter:instance_id=description.InstanceId
//listfilter:type_name=description.TypeName
type SSMInventoryEntryDescription struct {
	CaptureTime   *string
	InstanceId    *string
	SchemaVersion *string
	TypeName      *string
	Entries       map[string]string
}

//index:aws_ssm_maintenancewindow
//getfilter:window_id=description.MaintenanceWindowIdentity.WindowId
//listfilter:name=description.MaintenanceWindowIdentity.Name
//listfilter:enabled=description.MaintenanceWindowIdentity.Enabled
type SSMMaintenanceWindowDescription struct {
	MaintenanceWindowIdentity ssm.MaintenanceWindowIdentity
	MaintenanceWindow         *ssm_sdkv2.GetMaintenanceWindowOutput
	Tags                      []ssm.Tag
	Targets                   []ssm.MaintenanceWindowTarget
	Tasks                     []ssm.MaintenanceWindowTask
	ARN                       string
}

//index:aws_ssm_parameter
//getfilter:name=description.ParameterMetadata.Name
//listfilter:type=description.ParameterMetadata.Type
//listfilter:key_id=description.ParameterMetadata.KeyId
//listfilter:tier=description.ParameterMetadata.Tier
//listfilter:data_type=description.ParameterMetadata.DataType
type SSMParameterDescription struct {
	ParameterMetadata ssm.ParameterMetadata
	Parameter         *ssm.Parameter
	Tags              []ssm.Tag
}

//index:aws_ssm_patchbaseline
//getfilter:baseline_id=description.ParameterMetadata.Name
//listfilter:name=description.ParameterMetadata.Type
//listfilter:operating_system=description.ParameterMetadata.KeyId
type SSMPatchBaselineDescription struct {
	ARN                   string
	PatchBaselineIdentity ssm.PatchBaselineIdentity
	PatchBaseline         *ssm_sdkv2.GetPatchBaselineOutput
	Tags                  []ssm.Tag
}

//index:aws_ssm_managedinstancecompliance
//listfilter:resource_id=description.ComplianceItem.ResourceId
type SSMManagedInstanceComplianceDescription struct {
	ComplianceItem ssm.ComplianceItem
}

// index:aws_ssm_managed_instance_patch_state
//
//listfilter:instance_id=Description.PatchState.InstanceId
type SSMManagedInstancePatchStateDescription struct {
	PatchState ssm.InstancePatchState
}

//  ===================  ECS  ===================

//index:aws_ecs_taskdefinition
//getfilter:task_definition_arn=description.TaskDefinition.TaskDefinitionArn
type ECSTaskDefinitionDescription struct {
	TaskDefinition *ecs.TaskDefinition
	Tags           []ecs.Tag
}

//index:aws_ecs_cluster
//getfilter:cluster_arn=description.Cluster.ClusterArn
type ECSClusterDescription struct {
	Cluster ecs.Cluster
}

//index:aws_ecs_service
type ECSServiceDescription struct {
	Service ecs.Service
	Tags    []ecs.Tag
}

//index:aws_ecs_containerinstance
type ECSContainerInstanceDescription struct {
	ContainerInstance ecs.ContainerInstance
	Cluster           ecs.Cluster
}

//index:aws_ecs_taskset
//getfilter:id=description.TaskSet.Id
type ECSTaskSetDescription struct {
	TaskSet ecs.TaskSet
}

//index:aws_ecs_task
//listfilter:container_instance_arn=description.Task.ContainerInstanceArn
//listfilter:desired_status=description.Task.DesiredStatus
//listfilter:launch_type=description.Task.LaunchType
//listfilter:service_name=description.ServiceName
type ECSTaskDescription struct {
	Task           ecs.Task
	TaskProtection *ecs.ProtectedTask
	ServiceName    string
}

//  ===================  EFS  ===================

//index:aws_efs_filesystem
//getfilter:aws_efs_file_system=description.FileSystem.FileSystemId
type EFSFileSystemDescription struct {
	FileSystem efs.FileSystemDescription
	Policy     *string
}

//index:aws_efs_accesspoint
//getfilter:access_point_id=description.AccessPoint.AccessPointId
//listfilter:file_system_id=description.AccessPoint.FileSystemId
type EFSAccessPointDescription struct {
	AccessPoint efs.AccessPointDescription
}

//index:aws_efs_mounttarget
//getfilter:mount_target_id=description.MountTarget.MountTargetId
type EFSMountTargetDescription struct {
	MountTarget    efs.MountTargetDescription
	SecurityGroups []string
}

//  ===================  EKS  ===================

//index:aws_eks_cluster
//getfilter:name=description.Cluster.Name
type EKSClusterDescription struct {
	Cluster eks.Cluster
}

//index:aws_eks_addon
//getfilter:addon_name=description.Addon.AddonName
//getfilter:cluster_name=description.Addon.ClusterName
type EKSAddonDescription struct {
	Addon eks.Addon
}

//index:aws_eks_identityproviderconfig
//getfilter:name=description.ConfigName
//getfilter:type=description.ConfigType
//getfilter:cluster_name=description.IdentityProviderConfig.ClusterName
type EKSIdentityProviderConfigDescription struct {
	ConfigName             string
	ConfigType             string
	IdentityProviderConfig eks.OidcIdentityProviderConfig
}

//index:aws_eks_nodegroup
//getfilter:nodegroup_name=description.Nodegroup.NodegroupName
//getfilter:cluster_name=description.Nodegroup.ClusterName
//listfilter:cluster_name=description.Nodegroup.ClusterName
type EKSNodegroupDescription struct {
	Nodegroup eks.Nodegroup
}

//index:aws_eks_addonversion
//listfilter:addon_name=description.AddonName
type EKSAddonVersionDescription struct {
	AddonVersion       eks.AddonVersionInfo
	AddonConfiguration *string
	AddonName          *string
	AddonType          *string
}

//index:aws_eks_fargateprofile
//getfilter:cluster_name=description.Fargate.ClusterName
//getfilter:fargate_profile_name=description.Fargate.FargateProfileName
//listfilter:cluster_name=description.Fargate.ClusterName
type EKSFargateProfileDescription struct {
	FargateProfile eks.FargateProfile
}

//  ===================  WAFv2  ===================

//index:aws_wafv2_webacl
//getfilter:id=description.WebACL.Id
//getfilter:name=description.WebACL.Name
//getfilter:scope=description.Scope
type WAFv2WebACLDescription struct {
	WebACL               *wafv2.WebACL
	Scope                wafv2.Scope
	LoggingConfiguration *wafv2.LoggingConfiguration
	TagInfoForResource   *wafv2.TagInfoForResource
	LockToken            *string
	AssociatedResources  []string
}

//index:aws_wafv2_ipset
//getfilter:id=description.IPSetSummary.Id
//getfilter:name=description.IPSetSummary.Name
//getfilter:scope=description.IPSetSummary.Scope
type WAFv2IPSetDescription struct {
	IPSetSummary wafv2.IPSetSummary
	Scope        wafv2.Scope
	IPSet        *wafv2.IPSet
	Tags         []wafv2.Tag
}

//index:aws_wafv2_regexpatternset
//getfilter:id=description.IPSetSummary.Id
//getfilter:name=description.IPSetSummary.Name
//getfilter:scope=description.IPSetSummary.Scope
type WAFv2RegexPatternSetDescription struct {
	Scope                  wafv2.Scope
	RegexPatternSetSummary wafv2.RegexPatternSetSummary
	RegexPatternSet        *wafv2.RegexPatternSet
	Tags                   *wafv2op.ListTagsForResourceOutput
}

//index:aws_wafv2_rulegroup
//getfilter:id=description.RuleGroup.Id
//getfilter:name=description.RuleGroup.Name
//getfilter:scope=description.Tags
type WAFv2RuleGroupDescription struct {
	RuleGroup        *wafv2.RuleGroup
	RuleGroupSummary wafv2.RuleGroupSummary
	Tags             *wafv2op.ListTagsForResourceOutput
}

//  ===================  KMS  ===================

//index:aws_kms_key
//getfilter:id=description.Metadata.KeyId
type KMSKeyDescription struct {
	Metadata           *kms.KeyMetadata
	Aliases            []kms.AliasListEntry
	KeyRotationEnabled bool
	Policy             *string
	Tags               []kms.Tag
	Title              string
}

type KMSKeyRotationDescription struct {
	KeyId        string
	KeyArn       string
	RotationType kms.RotationType
	RotationDate time.Time
}

//index:aws_kms_alias
type KMSAliasDescription struct {
	Alias kms.AliasListEntry
}

//  ===================  Lambda  ===================

//index:aws_lambda_function
//getfilter:name=description.Function.Configuration.FunctionName
type LambdaFunctionDescription struct {
	Function  *lambdaop.GetFunctionOutput
	UrlConfig []lambda.FunctionUrlConfig
	Policy    *lambdaop.GetPolicyOutput
}

//index:aws_lambda_functionversion
//getfilter:version=description.FunctionVersion.Version
//getfilter:function_name=description.FunctionVersion.FunctionName
//listfilter:function_name=description.FunctionVersion.FunctionName
type LambdaFunctionVersionDescription struct {
	FunctionVersion lambda.FunctionConfiguration
	Policy          *lambdaop.GetPolicyOutput
}

//index:aws_lambda_alias
//getfilter:name=description.Alias.Name
//getfilter:function_name=description.FunctionName
//getfilter:region=description.Alias.AliasName
//listfilter:function_version=description.Alias.FunctionVersion
//listfilter:function_name=description.FunctionName
type LambdaAliasDescription struct {
	FunctionName string
	Alias        lambda.AliasConfiguration
	Policy       *lambdaop.GetPolicyOutput
	UrlConfig    lambdaop.GetFunctionUrlConfigOutput
}

//index:aws_lambda_layer
type LambdaLayerDescription struct {
	Layer lambda.LayersListItem
}

//index:aws_lambda_layerversion
//getfilter:layer_name=description.LayerName
//getfilter:version=description.LayerVersion.Version
//listfilter:layer_name=description.LayerName
type LambdaLayerVersionDescription struct {
	LayerName    string
	LayerVersion lambdaop.GetLayerVersionOutput
	Policy       lambdaop.GetLayerVersionPolicyOutput
}

//index:aws_s3_accesspoint
//getfilter:name=description.AccessPoint.Name
//getfilter:region=metadata.region
type S3AccessPointDescription struct {
	AccessPoint  *s3controlop.GetAccessPointOutput
	Policy       *string
	PolicyStatus *s3control.PolicyStatus
}

type CostExplorerRow struct {
	Estimated bool

	// The time period that the result covers.
	PeriodStart *string
	PeriodEnd   *string

	Dimension1 *string
	Dimension2 *string
	//Tag *string

	BlendedCostAmount      *float64
	UnblendedCostAmount    *float64
	NetUnblendedCostAmount *float64
	AmortizedCostAmount    *float64
	NetAmortizedCostAmount *float64
	UsageQuantityAmount    *float64
	NormalizedUsageAmount  *float64

	BlendedCostUnit      *string
	UnblendedCostUnit    *string
	NetUnblendedCostUnit *string
	AmortizedCostUnit    *string
	NetAmortizedCostUnit *string
	UsageQuantityUnit    *string
	NormalizedUsageUnit  *string

	MeanValue *string
}

//index:aws_costexplorer_byaccountmonthly
type CostExplorerByAccountMonthlyDescription struct {
	CostExplorerRow
}

//index:aws_costexplorer_byservicemonthly
type CostExplorerByServiceMonthlyDescription struct {
	CostExplorerRow
}

//index:aws_costexplorer_byrecordtypemonthly
type CostExplorerByRecordTypeMonthlyDescription struct {
	CostExplorerRow
}

//index:aws_costexplorer_byusagetypemonthly
type CostExplorerByServiceUsageTypeMonthlyDescription struct {
	CostExplorerRow
}

//index:aws_costexplorer_forcastmonthly
type CostExplorerForcastMonthlyDescription struct {
	CostExplorerRow
}

//index:aws_costexplorer_byaccountdaily
type CostExplorerByAccountDailyDescription struct {
	CostExplorerRow
}

//index:aws_costexplorer_byservicedaily
//listfilter:service=description.Dimension1
//listfilter:cost_source=description.Dimension2
type CostExplorerByServiceDailyDescription struct {
	CostExplorerRow
	CostDateMillis int64
}

//index:aws_costexplorer_byrecordtypedaily
type CostExplorerByRecordTypeDailyDescription struct {
	CostExplorerRow
}

//index:aws_costexplorer_byusagetypedaily
type CostExplorerByServiceUsageTypeDailyDescription struct {
	CostExplorerRow
}

//index:aws_costexplorer_forcastdaily
type CostExplorerForcastDailyDescription struct {
	CostExplorerRow
}

//  ===================  ECR  ===================

//index:aws_ecr_repository
//getfilter:repository_name=description.Repository.RepositoryName
type ECRRepositoryDescription struct {
	Repository                      ecr.Repository
	LifecyclePolicy                 *ecrop.GetLifecyclePolicyOutput
	ImageDetails                    []ecr.ImageDetail
	Policy                          *ecrop.GetRepositoryPolicyOutput
	RepositoryScanningConfiguration *ecrop.BatchGetRepositoryScanningConfigurationOutput
	Tags                            []ecr.Tag
}

//index:aws_ecr_image
//listfilter:repository_name=description.Image.RepositoryName
//listfilter:registry_id=description.Image.RegistryId
type ECRImageDescription struct {
	Image    ecr.ImageDetail
	ImageUri string
}

//index:aws_ecrpublic_repository
//getfilter:repository_name=description.PublicRepository.RepositoryName
type ECRPublicRepositoryDescription struct {
	PublicRepository ecrpublic.Repository
	ImageDetails     []ecrpublic.ImageDetail
	Policy           *ecrpublicop.GetRepositoryPolicyOutput
	Tags             []ecrpublic.Tag
}

//index:aws_ecrpublic_registry
//getfilter:registry_id=description.PublicRegistry.RegistryId
type ECRPublicRegistryDescription struct {
	PublicRegistry ecrpublic.Registry
	Tags           []ecrpublic.Tag
}

//index:aws_ecr_registry
//getfilter:registry_id=description.Registry.RegistryId
type ECRRegistryDescription struct {
	RegistryId       string
	ReplicationRules []ecr.ReplicationRule
}

type ECRRegistryScanningConfigurationDescription struct {
	RegistryId            string
	ScanningConfiguration *ecr.RegistryScanningConfiguration
}

//  ===================  EventBridge  ===================

//index:aws_eventbridge_eventbus
//getfilter:arn=description.Bus.Arn
type EventBridgeBusDescription struct {
	Bus  eventbridge.EventBus
	Tags []eventbridge.Tag
}

//index:aws_eventbridge_eventrule
//getfilter:name=description.Rule.Name
//listfilter:event_bus_name=description.Rule.EventBusName
//listfilter:name_prefix=description.Rule.Name
type EventBridgeRuleDescription struct {
	Rule    eventbridgeop.DescribeRuleOutput
	Targets []eventbridge.Target
	Tags    []eventbridge.Tag
}

//  ===================  AppStream  ===================

//index:aws_appstream_application
//getfilter:name=description.Application.Name
type AppStreamApplicationDescription struct {
	Application appstream.Application
	Tags        map[string]string
}

//index:aws_appstream_stack
//getfilter:name=description.Stack.Name
type AppStreamStackDescription struct {
	Stack appstream.Stack
	Tags  map[string]string
}

//index:aws_appstream_fleet
//getfilter:name=description.Fleet.Name
type AppStreamFleetDescription struct {
	Fleet appstream.Fleet
	Tags  map[string]string
}

//index:aws_appstream_image
//getfilter:name=description.Image.Name
type AppStreamImageDescription struct {
	Image appstream.Image
	Tags  map[string]string
}

// ===================  Athena  ===================

//index:aws_athena_workgroup
//getfilter:name=description.WorkGroup.Name
type AthenaWorkGroupDescription struct {
	WorkGroup *athena.WorkGroup
}

//index:aws_athena_queryexecution
//getfilter:name=description.QueryExecution.Query
type AthenaQueryExecutionDescription struct {
	QueryExecution *athena.QueryExecution
}

//  ===================  Kinesis  ===================

//index:aws_kinesis_stream
//getfilter:stream_name=description.Stream.StreamName
type KinesisStreamDescription struct {
	Stream             kinesis.StreamDescription
	DescriptionSummary kinesis.StreamDescriptionSummary
	Tags               []kinesis.Tag
}

//index:aws_kinesisvideo_stream
//getfilter:stream_name=description.Stream.StreamName
type KinesisVideoStreamDescription struct {
	Stream kinesisvideo.StreamInfo
	Tags   map[string]string
}

//index:aws_kinesis_consumer
//getfilter:consumer_arn=description.Consumer.ConsumerARN
type KinesisConsumerDescription struct {
	StreamARN string
	Consumer  kinesis.Consumer
}

//index:aws_kinesisanalyticsv2_application
//getfilter:application_name=description.Application.ApplicationName
type KinesisAnalyticsV2ApplicationDescription struct {
	Application kinesisanalyticsv2.ApplicationDetail
	Tags        []kinesisanalyticsv2.Tag
}

//  ===================  Glacier  ===================

//index:aws_glacier_vault
//getfilter:vault_name=description.Vault.VaultName
type GlacierVaultDescription struct {
	Vault                   glacier.DescribeVaultOutput
	AccessPolicy            glacier.VaultAccessPolicy
	LockPolicy              glacier.VaultLockPolicy
	VaultNotificationConfig glacier.VaultNotificationConfig
	Tags                    map[string]string
}

//  ===================  Workspace  ===================

//index:aws_workspaces_workspace
//getfilter:workspace_id=description.Workspace.WorkspaceId
type WorkspacesWorkspaceDescription struct {
	Workspace workspaces.Workspace
	Tags      []workspaces.Tag
}

//index:aws_workspaces_bundle
//getfilter:bundle_id=description.Bundle.BundleId
type WorkspacesBundleDescription struct {
	Bundle workspaces.WorkspaceBundle
	Tags   []workspaces.Tag
}

//  ===================  KeySpaces (For Apache Cassandra)  ===================

//index:aws_keyspaces_keyspace
//getfilter:keyspace_name=description.Keyspace.KeyspaceName
type KeyspacesKeyspaceDescription struct {
	Keyspace keyspaces.KeyspaceSummary
	Tags     []keyspaces.Tag
}

//index:aws_keyspaces_table
//getfilter:table_name=description.Table.TableName
type KeyspacesTableDescription struct {
	Table keyspaces.TableSummary
	Tags  []keyspaces.Tag
}

//  ===================  Grafana  ===================

//index:aws_grafana_workspace
//getfilter:id=description.Workspace.Id
type GrafanaWorkspaceDescription struct {
	Workspace grafana.WorkspaceSummary
}

//  ===================  AMP (Amazon Managed Service for Prometheus)  ===================

//index:aws_amp_workspace
//getfilter:workspace_id=description.Workspace.WorkspaceId
type AMPWorkspaceDescription struct {
	Workspace amp.WorkspaceSummary
}

//  ===================  Kafka  ===================

//index:aws_kafka_cluster
//getfilter:cluster_name=description.Cluster.ClusterName
type KafkaClusterDescription struct {
	Cluster              kafka.Cluster
	Configuration        *kafkaop.DescribeConfigurationOutput
	ClusterOperationInfo *kafka.ClusterOperationInfo
}

//  ===================  MWAA (Managed Workflows for Apache Airflow) ===================

//index:aws_mwaa_environment
//getfilter:name=description.Environment.Name
type MWAAEnvironmentDescription struct {
	Environment mwaa.Environment
}

//  ===================  MemoryDb  ===================

//index:aws_memorydb_cluster
//getfilter:name=description.Cluster.Name
type MemoryDbClusterDescription struct {
	Cluster memorydb.Cluster
	Tags    []memorydb.Tag
}

//  ===================  MQ  ===================

//index:aws_mq_broker
//getfilter:broker_name=description.Broker.BrokerName
type MQBrokerDescription struct {
	BrokerDescription *mq.DescribeBrokerOutput
	Tags              map[string]string
}

//  ===================  Neptune  ===================

//index:aws_neptune_database
//getfilter:db_instance_identifier=description.Database.DBInstanceIdentifier
type NeptuneDatabaseDescription struct {
	Database neptune.DBInstance
	Tags     []neptune.Tag
}

//index:aws_neptune_databasecluster
//getfilter:db_instance_identifier=description.Database.DBInstanceIdentifier
type NeptuneDatabaseClusterDescription struct {
	Cluster neptune.DBCluster
	Tags    []neptune.Tag
}

type NeptuneDatabaseClusterSnapshotDescription struct {
	Snapshot   neptune.DBClusterSnapshot
	Attributes []map[string]interface{}
}

//  ===================  OpenSearch  ===================

//index:aws_opensearch_domain
//getfilter:domain_name=description.Domain.DomainName
type OpenSearchDomainDescription struct {
	Domain opensearch.DomainStatus
	Tags   []opensearch.Tag
}

//  ===================  SES (Simple Email Service)  ===================

//index:aws_ses_configurtionset
//getfilter:name=description.ConfigurationSet.Name
type SESConfigurationSetDescription struct {
	ConfigurationSet ses.ConfigurationSet
}

//index:aws_ses_identity
//getfilter:identity_name=description.Identity.IdentityName
//listfilter:identity_type=description.Identity.IdentityType
type SESIdentityDescription struct {
	Identity               string
	VerificationAttributes ses.IdentityVerificationAttributes
	NotificationAttributes ses.IdentityNotificationAttributes
	DkimAttributes         map[string]ses.IdentityDkimAttributes
	IdentityMail           map[string]ses.IdentityMailFromDomainAttributes
	Tags                   []types.Tag
	ARN                    string
}

type SESv2EmailIdentityDescription struct {
	ARN      string
	Identity types.IdentityInfo
	Tags     []types.Tag
}

//  ===================  CloudFormation  ===================

//index:aws_cloudformation_stack
//getfilter:name=description.Stack.StackName
//listfilter:name=description.Stack.StackName
type CloudFormationStackDescription struct {
	Stack          cloudformation.Stack
	StackTemplate  cloudformationop.GetTemplateOutput
	StackResources []cloudformation.StackResource
}

//index:aws_cloudformation_stackset
//getfilter:stack_set_name=description.StackSet.StackSetName
type CloudFormationStackSetDescription struct {
	StackSet cloudformation.StackSet
}

//index:aws_cloudformation_stackresource
//listfilter:name=description.StackResource.StackName
type CloudFormationStackResourceDescription struct {
	StackResource cloudformation.StackResourceDetail
}

//  ===================  CodeCommit  ===================

//index:aws_codecommit_repository
type CodeCommitRepositoryDescription struct {
	Repository codecommit.RepositoryMetadata
	Tags       map[string]string
}

//  ===================  CodePipeline  ===================

//index:aws_codepipeline_pipeline
//getfilter:name=description.Pipeline.Name
type CodePipelinePipelineDescription struct {
	Pipeline codepipeline.PipelineDeclaration
	Metadata codepipeline.PipelineMetadata
	Tags     []codepipeline.Tag
}

//  ===================  DirectoryService  ===================

//index:aws_directoryservice_directory
//getfilter:name=description.Directory.DirectoryId
type DirectoryServiceDirectoryDescription struct {
	Directory       directoryservice.DirectoryDescription
	Snapshot        directoryservice.SnapshotLimits
	EventTopics     []directoryservice.EventTopic
	SharedDirectory []directoryservice.SharedDirectory
	Tags            []directoryservice.Tag
}

//index:aws_directoryservice_certificate
//getfilter:name=description.Certificate.CertificateId
type DirectoryServiceCertificateDescription struct {
	Certificate directoryservice.Certificate
	DirectoryId string
}

//index:aws_directoryservice_logsubscription
type DirectoryServiceLogSubscriptionDescription struct {
	LogSubscription directoryservice.LogSubscription
}

//  ===================  SSOAdmin  ===================

//index:aws_ssoadmin_instance
type SSOAdminInstanceDescription struct {
	Instance ssoadmin.InstanceMetadata
}

//index:aws_ssoadmin_account_assignment
type SSOAdminAccountAssignmentDescription struct {
	Instance          ssoadmin.InstanceMetadata
	AccountAssignment ssoadmin.AccountAssignment
}

//index:aws_ssoadmin_permission_set
type SSOAdminPermissionSetDescription struct {
	InstanceArn   string
	PermissionSet ssoadmin.PermissionSet
	Tags          interface{}
}

//index:aws_ssoadmin_managed_policy_attachment
type SSOAdminPolicyAttachmentDescription struct {
	InstanceArn           string
	PermissionSetArn      string
	AttachedManagedPolicy ssoadmin.AttachedManagedPolicy
}

//index:aws_ssoadmin_usereffectivea
type UserEffectiveAccessDescription struct {
	Instance          ssoadmin.InstanceMetadata
	AccountAssignment ssoadmin.AccountAssignment
	User              identitystore2.DescribeUserOutput
	UserId            interface{}
}

//  ===================  Tagging  ===================

//index:aws_tagging_resources
//getfilter:arn=description.
type TaggingResourcesDescription struct {
	TagMapping types2.ResourceTagMapping
}

//  ===================  WAF  ===================

//index:aws_waf_rule
//getfilter:rule_id=description.Rule.RuleId
type WAFRuleDescription struct {
	Rule waf.Rule
	Tags []waf.Tag
}

//index:aws_wafregional_rule
//getfilter:rule_id=description.Rule.RuleId
type WAFRegionalRuleDescription struct {
	Rule wafregional.Rule
	Tags []wafregional.Tag
}

//index:aws_waf_ratebasedrule
//getfilter:rule_id=description.Rule.RuleId
type WAFRateBasedRuleDescription struct {
	ARN         string
	RuleSummary waf.RuleSummary
	Rule        *waf.RateBasedRule
	Tags        *waf.TagInfoForResource
}

//index:aws_waf_rulegroup
//getfilter:rule_group_id=description.Rule.RuleId
type WAFRuleGroupDescription struct {
	ARN              string
	RuleGroupSummary waf.RuleGroupSummary
	RuleGroup        *waf2.GetRuleGroupOutput
	ActivatedRules   *waf2.ListActivatedRulesInRuleGroupOutput
	Tags             []waf.Tag
}

//index:aws_waf_webacl
//getfilter:web_acl_id=description.WebACL.WebACLId
type WAFWebAclDescription struct {
	WebACLSummary        waf.WebACLSummary
	WebACL               *waf.WebACL
	LoggingConfiguration *waf.LoggingConfiguration
	Tags                 *waf.TagInfoForResource
}

//index:aws_wellarchitected_workload
//getfilter:workload_id=description.Workload.WorkloadId
//listfilter:workload_name=description.Workload.WorkloadName
type WellArchitectedWorkloadDescription struct {
	WorkloadSummary types3.WorkloadSummary
	Workload        *types3.Workload
}

//index:aws_wellarchitected_answer
//getfilter:workload_id=description.WorkloadId
type WellArchitectedAnswerDescription struct {
	Answer          types3.Answer
	WorkloadId      string
	WorkloadName    string
	LensAlias       string
	LensArn         string
	MilestoneNumber *int32
}

//index:aws_wellarchitected_checkdetail
//getfilter:workload_id=description.WorkloadId
type WellArchitectedCheckDetailDescription struct {
	CheckDetail types3.CheckDetail
	WorkloadId  string
}

//index:aws_wellarchitected_checksymmary
//getfilter:workload_id=description.WorkloadId
type WellArchitectedCheckSummaryDescription struct {
	CheckSummary types3.CheckSummary
	WorkloadId   string
}

//index:aws_wellarchitected_consolidated_report
type WellArchitectedCheckConsolidatedReportDescription struct {
	IncludeSharedResources *bool
	ConsolidateReport      types3.ConsolidatedReportMetric
	Base64                 string
}

//index:aws_wellarchitected_lens
type WellArchitectedLensDescription struct {
	Lens        types3.Lens
	LensSummary types3.LensSummary
}

//index:aws_wellarchitected_lensreview
type WellArchitectedLensReviewDescription struct {
	LensReview types3.LensReview
}

//index:aws_wellarchitected_lensreviewimprovement
//getfilter:workload_id=description.WorkloadId
type WellArchitectedLensReviewImprovementDescription struct {
	ImprovementSummary types3.ImprovementSummary
	LensAlias          string
	LensArn            string
	MilestoneNumber    *int32
	WorkloadId         string
}

//index:aws_wellarchitected_lensreviewreport
//getfilter:workload_id=description.WorkloadId
type WellArchitectedLensReviewReportDescription struct {
	Report          types3.LensReviewReport
	MilestoneNumber *int32
	WorkloadId      string
}

//index:aws_wellarchitected_lensshare
type WellArchitectedLensShareDescription struct {
	Share types3.LensShareSummary
	Lens  types3.Lens
}

//index:aws_wellarchitected_milestone
type WellArchitectedMilestoneDescription struct {
	Milestone types3.Milestone
}

//index:aws_wellarchitected_notification
type WellArchitectedNotificationDescription struct {
	Notification types3.NotificationSummary
}

//index:aws_wellarchitected_shareinvitation
type WellArchitectedShareInvitationDescription struct {
	ShareInvitation types3.ShareInvitationSummary
}

//index:aws_wellarchitected_shareinvitation
type WellArchitectedWorkloadShareDescription struct {
	Share      types3.WorkloadShareSummary
	WorkloadId string
	Arn        string
}

//index:aws_wafregional_webacl
//getfilter:web_acl_id=description.WebACL.WebACLId
type WAFRegionalWebAclDescription struct {
	WebACL               *wafregional.WebACL
	AssociatedResources  []string
	LoggingConfiguration *wafregional.LoggingConfiguration
	Tags                 []wafregional.Tag
}

//index:aws_wafregional_rulegroup
//getfilter:rule_group_id=description.Rule.RuleId
type WAFRegionalRuleGroupDescription struct {
	ARN              string
	RuleGroupSummary wafregional.RuleGroupSummary
	RuleGroup        *wafregional.RuleGroup
	ActivatedRules   []wafregional.ActivatedRule
	Tags             []wafregional.Tag
}

//  ===================  Route53  ===================

//index:aws_route53_hostedzone
//getfilter:id=description.ID
type Route53HostedZoneDescription struct {
	ID                  string
	HostedZone          route53.HostedZone
	QueryLoggingConfigs []route53.QueryLoggingConfig
	Limit               *route53.HostedZoneLimit
	DNSSec              route53op.GetDNSSECOutput
	Tags                []route53.Tag
}

//index:aws_route53_healthcheck
//getfilter:id=description.HealthCheck.Id
type Route53HealthCheckDescription struct {
	HealthCheck route53.HealthCheck
	Status      *route53op.GetHealthCheckStatusOutput
	Tags        *route53op.ListTagsForResourceOutput
}

//index:aws_route53resolver_resolverrule
//getfilter:id=description.ResolverRole.Id
//listfilter:creator_request_id=description.ResolverRole.CreatorRequestId
//listfilter:domain_name=description.ResolverRole.DomainName
//listfilter:name=description.ResolverRole.Name
//listfilter:resolver_endpoint_id=description.ResolverRole.ResolverEndpointId
//listfilter:status=description.ResolverRole.Status
type Route53ResolverResolverRuleDescription struct {
	ResolverRole     route53resolver.ResolverRule
	Tags             []route53resolver.Tag
	RuleAssociations *route53resolverop.ListResolverRuleAssociationsOutput
}

//index:aws_route53resolver_resolverendpoint
//getfilter:id=description.ResolverEndpoint.Id
//listfilter:creator_request_id=description.ResolverEndpoint.CreatorRequestId
//listfilter:direction=description.ResolverEndpoint.Direction
//listfilter:host_vpc_id=description.ResolverEndpoint.HostVPCId
//listfilter:ip_address_count=description.ResolverEndpoint.IpAddressCount
//listfilter:status=description.ResolverEndpoint.Status
//listfilter:name=description.ResolverEndpoint.Name
type Route53ResolverEndpointDescription struct {
	ResolverEndpoint route53resolver.ResolverEndpoint
	IpAddresses      []route53resolver.IpAddressResponse
	Tags             []route53resolver.Tag
}

//index:aws_route53domains_domain
//getfilter:domain_name=description.Domain.DomainName
type Route53DomainDescription struct {
	DomainSummary route53domains.DomainSummary
	Domain        route53domainsop.GetDomainDetailOutput
	Tags          []route53domains.Tag
}

//index:aws_route53_record
//listfilter:zone_id=description.ZoneId
//listfilter:name=description.Record.Name
//listfilter:set_identifier=description.Record.SetIdentifier
//listfilter:type=description.Record.Type
type Route53RecordDescription struct {
	ZoneID string
	Record route53.ResourceRecordSet
}

//index:aws_route53_trafficpolicy
//getfilter:id=description.TrafficPolicy.Id
//getfilter:version=description.TrafficPolicy.Version
type Route53TrafficPolicyDescription struct {
	TrafficPolicy route53.TrafficPolicy
}

//index:aws_route53_trafficpolicyinstance
//getfilter:id=description.TrafficPolicyInstance.Id
type Route53TrafficPolicyInstanceDescription struct {
	TrafficPolicyInstance route53.TrafficPolicyInstance
}

//index:aws_route53_querylog
//getfilter:id=description.TrafficPolicyInstance.Id
type Route53QueryLogDescription struct {
	QueryConfig route53.QueryLoggingConfig
}

//index:aws_route53_querylog
//getfilter:id=description.TrafficPolicyInstance.Id
type Route53ResolverQueryLogConfigDescription struct {
	QueryConfig route53resolver.ResolverQueryLogConfig
}

//  ===================  Batch  ===================

//index:aws_batch_computeenvironment
//getfilter:compute_environment_name=description.ComputeEnvironment.ComputeEnvironmentName
type BatchComputeEnvironmentDescription struct {
	ComputeEnvironment batch.ComputeEnvironmentDetail
}

//index:aws_batch_job
//getfilter:job_name=description.Job.JobName
type BatchJobDescription struct {
	Job batch.JobSummary
}

//index:aws_batch_jobqueue
//getfilter:job_queue_name=description.JobQueue.JobQueueName
type BatchJobQueueDescription struct {
	JobQueue batch.JobQueueDetail
}

//  ===================  CodeArtifact  ===================

//index:aws_codeartifact_repository
//getfilter:name=description.Repository.Name
type CodeArtifactRepositoryDescription struct {
	Repository  codeartifact.RepositorySummary
	Policy      codeartifact.ResourcePolicy
	Description codeartifact.RepositoryDescription
	Endpoints   []string
	Tags        []codeartifact.Tag
}

//index:aws_codeartifact_domain
//getfilter:name=description.Domain.Name
//getfilter:name=description.Domain.Owner
type CodeArtifactDomainDescription struct {
	Domain codeartifact.DomainDescription
	Policy codeartifact.ResourcePolicy
	Tags   []codeartifact.Tag
}

//  ===================  CodeDeploy  ===================

//index:aws_codedeploy_deploymentgroup
//getfilter:deployment_group_name=description.DeploymentGroup.DeploymentGroupName
type CodeDeployDeploymentGroupDescription struct {
	DeploymentGroup codedeploy.DeploymentGroupInfo
	Tags            []codedeploy.Tag
}

//index:aws_codedeploy_application
//getfilter:application_name=description.Application.ApplicationName
type CodeDeployApplicationDescription struct {
	Application codedeploy.ApplicationInfo
	Tags        []codedeploy.Tag
}

//index:aws_codedeploy_application
//getfilter:application_name=description.Application.ApplicationName
type CodeDeployDeploymentConfigDescription struct {
	Config codedeploy.DeploymentConfigInfo
}

//  ===================  CodeStar  ===================

//index:aws_codestar_project
//getfilter:id=description.Project.Id
type CodeStarProjectDescription struct {
	Project codestarop.DescribeProjectOutput
	Tags    map[string]string
}

//  ===================  DirectConnect  ===================

//index:aws_directconnect_connection
//getfilter:connection_id=description.Connection.ConnectionId
type DirectConnectConnectionDescription struct {
	Connection directconnect.Connection
}

//index:aws_directconnect_gateway
//getfilter:direct_connect_gateway_id=description.Gateway.DirectConnectGatewayId
type DirectConnectGatewayDescription struct {
	Gateway directconnect.DirectConnectGateway
	Tags    []directconnect.Tag
}

//  ===================  Elastic Disaster Recovery (DRS)  ===================

//index:aws_drs_sourceserver
//getfilter:source_server_id=description.SourceServer.SourceServerID
type DRSSourceServerDescription struct {
	SourceServer        drs.SourceServer
	LaunchConfiguration drs2.GetLaunchConfigurationOutput
}

//index:aws_drs_recoveryinstance
//getfilter:recovery_instance_id=description.RecoveryInstance.RecoveryInstanceID
type DRSRecoveryInstanceDescription struct {
	RecoveryInstance drs.RecoveryInstance
}

//index:aws_drs_job
//listfilter:job_id=description.Job.JobID
//listfilter:creation_date_time=description.Job.CreationDateTime
//listfilter:end_date_time=description.Job.EndDateTime
type DRSJobDescription struct {
	Job drs.Job
}

//index:aws_drs_recoverysnapshot
//listfilter:source_server_id=description.RecoveryInstance.SourceServerID
//listfilter:timestamp=description.RecoveryInstance.Timestamp
type DRSRecoverySnapshotDescription struct {
	RecoverySnapshot drs.RecoverySnapshot
}

//  ===================  Firewall Manager Policy (FMS)  ===================

//index:aws_fms_policy
//getfilter:policy_name=description.Policy.PolicyName
type FMSPolicyDescription struct {
	Policy fms.PolicySummary
	Tags   []fms.Tag
}

//  ===================  Network Firewall ===================

//index:aws_networkfirewall_firewall
//getfilter:firewall_name=description.Firewall.FirewallName
type NetworkFirewallFirewallDescription struct {
	Firewall             networkfirewall.Firewall
	FirewallStatus       networkfirewall.FirewallStatus
	LoggingConfiguration *networkfirewall2.DescribeLoggingConfigurationOutput
}

//index:aws_networkfirewall_firewallpolicy
//getfilter:arn=description.FirewallPolicyResponse.FirewallPolicyArn
//getfilter:name=description.FirewallPolicyResponse.FirewallPolicyName
type NetworkFirewallFirewallPolicyDescription struct {
	FirewallPolicy         *networkfirewall.FirewallPolicy
	FirewallPolicyResponse *networkfirewall.FirewallPolicyResponse
}

//index:aws_networkfirewall_rulegroup
//getfilter:arn=description.RuleGroupResponse.RuleGroupArn
//getfilter:rule_group_name=description.RuleGroupResponse.RuleGroupName
type NetworkFirewallRuleGroupDescription struct {
	RuleGroup         *networkfirewall.RuleGroup
	RuleGroupResponse *networkfirewall.RuleGroupResponse
}

//  ===================  OpsWork ===================

//index:aws_opsworkscm_server
//getfilter:server_name=description.Server.ServerName
type OpsWorksCMServerDescription struct {
	Server opsworkscm.Server
	Tags   []opsworkscm.Tag
}

//  ===================  Organizations ===================

//index:aws_organizations_organization
//getfilter:id=description.Organization.Id
type OrganizationsOrganizationDescription struct {
	Organization organizations.Organization
}

//index:aws_organizations_account
//getfilter:id=description.Account.Id
type OrganizationsAccountDescription struct {
	Tags     []organizations.Tag
	Account  organizations.Account
	ParentID string
}

//index:aws_organizations_policy
//getfilter:id=description.Policy.PolicySummary.Id
type OrganizationsPolicyDescription struct {
	Policy organizations.Policy
}

type OrganizationsRootDescription struct {
	Root organizations.Root
}

type OrganizationsOrganizationalUnitDescription struct {
	Unit     organizations.OrganizationalUnit
	ParentId string
	Path     string
	Tags     []organizations.Tag
}

type OrganizationsPolicyTargetDescription struct {
	PolicySummary organizations.PolicySummary
	PolicyContent *string
	TargetId      string
}

// ===================  Pinpoint ===================

//index:aws_pinpoint_app
//getfilter:id=description.App.Id
type PinPointAppDescription struct {
	App      pinpoint.ApplicationResponse
	Settings *pinpoint.ApplicationSettingsResource
}

// ===================  Pipes ===================

//index:aws_pipes_pipe
//getfilter:name=description.PipeOutput.Name
//listfilter:current_state=description.PipeOutput.CurrentState
//listfilter:desired_state=description.PipeOutput.DesiredState
type PipesPipeDescription struct {
	PipeOutput *pipesop.DescribePipeOutput
	Pipe       pipes.Pipe
}

// ===================  ResourceGroups ===================

//index:aws_resourcegroups_group
//getfilter:name=description.GroupIdentifier.GroupName
type ResourceGroupsGroupDescription struct {
	GroupIdentifier types4.GroupIdentifier
	Resources       []types4.ListGroupResourcesItem
	Tags            *resourcegroups.GetTagsOutput
}

// ===================  OpenSearchServerless ===================

//index:aws_opensearchserverless_collection
//getfilter:name=description.CollectionSummary.Name
type OpenSearchServerlessCollectionDescription struct {
	CollectionSummary types6.CollectionSummary
	Collection        types6.CollectionDetail
	Tags              *opensearchserverless.ListTagsForResourceOutput
}

// ===================  Timestream ===================

//index:aws_timestream_database
//getfilter:arn=description.Database.Arn
//getfilter:name=description.Database.DatabaseName
type TimestreamDatabaseDescription struct {
	Database types5.Database
	Tags     *timestreamwrite.ListTagsForResourceOutput
}

// ===================  ResourceExplorer2 ===================

//index:aws_resourceexplorer2_index
//listfilter:type=description.Index.Type
//listfilter:region=description.Index.Region
type ResourceExplorer2IndexDescription struct {
	Index resourceexplorer2.Index
}

//index:aws_resourceexplorer2_supportedresourcetype
type ResourceExplorer2SupportedResourceTypeDescription struct {
	SupportedResourceType resourceexplorer2.SupportedResourceType
}

// ===================  StepFunctions ===================

//index:aws_stepfunctions_statemachine
//getfilter:arn=description.StateMachineItem.StateMachineArn
type StepFunctionsStateMachineDescription struct {
	StateMachineItem sfn.StateMachineListItem
	StateMachine     *sfnop.DescribeStateMachineOutput
	Tags             []sfn.Tag
}

//index:aws_stepfunctions_statemachineexecutionhistories
type StepFunctionsStateMachineExecutionHistoriesDescription struct {
	ExecutionHistory sfn.HistoryEvent
	ARN              string
}

//index:aws_stepfunctions_statemachineexecution
//getfilter:execution_arn=description.ExecutionItem.ExecutionArn
//listfilter:status=description.ExecutionItem.Status
//listfilter:state_machine_arn=description.ExecutionItem.StateMachineArn
type StepFunctionsStateMachineExecutionDescription struct {
	ExecutionItem sfn.ExecutionListItem
	Execution     *sfnop.DescribeExecutionOutput
}

// ===================  SimSpaceWeaver ===================

//index:aws_simspaceweaversimulation
//getfilter:name=description.Simulation.Name
type SimSpaceWeaverSimulationDescription struct {
	Simulation     simspaceweaver.SimulationMetadata
	SimulationItem *simspaceweaverop.DescribeSimulationOutput
	Tags           map[string]string
}

//  ===================  ACM ===================

//index:aws_acmpca_certificateauthority
//getfilter:arn=description.CertificateAuthority.Arn
type ACMPCACertificateAuthorityDescription struct {
	CertificateAuthority acmpca.CertificateAuthority
	Tags                 []acmpca.Tag
}

//  ===================  Shield ===================

//index:aws_shield_protectiongroup
//getfilter:protection_group_id=description.ProtectionGroup.ProtectionGroupId
type ShieldProtectionGroupDescription struct {
	ProtectionGroup shield.ProtectionGroup
	Tags            []shield.Tag
}

//  ===================  Storage Gateway ===================

//index:aws_storagegateway_storagegateway
//getfilter:gateway_id=description.StorageGateway.GatewayId
type StorageGatewayStorageGatewayDescription struct {
	StorageGateway storagegateway.GatewayInfo
	Tags           []storagegateway.Tag
}

//  ===================  Image Builder ===================

//index:aws_imagebuilder_image
//getfilter:name=description.Image.Name
type ImageBuilderImageDescription struct {
	Image imagebuilder.Image
}

// ===================  Account ===================

//index:aws_account_alternatecontact
//listfilter:linked_account_id=description.LinkedAccountID
//listfilter:contact_type=description.AlternateContact.AlternateContactType
type AccountAlternateContactDescription struct {
	AlternateContact account.AlternateContact
	LinkedAccountID  string
}

//index:aws_account_contact
//listfilter:linked_account_id=description.LinkedAccountID
type AccountContactDescription struct {
	AlternateContact account.ContactInformation
	LinkedAccountID  string
}

// ===================  Amplify ===================

//index:aws_amplify_app
//getfilter:app_id=description.App.AppId
type AmplifyAppDescription struct {
	App amplify.App
}

// ===================  App Config (appconfig) ===================

//index:aws_appconfig_application
//getfilter:id=description.Application.Id
type AppConfigApplicationDescription struct {
	Application appconfig.Application
	Tags        map[string]string
}

// ===================  Audit Manager ===================

//index:aws_auditmanager_assessment
//getfilter:assessment_id=description.Assessment.Metadata.Id
type AuditManagerAssessmentDescription struct {
	Assessment auditmanager.Assessment
}

//index:aws_auditmanager_control
//getfilter:control_id=description.Control.Id
type AuditManagerControlDescription struct {
	Control auditmanager.Control
}

//index:aws_auditmanager_evidence
//getfilter:id=description.Evidence.Id
//getfilter:evidence_folder_id=description.Evidence.EvidenceFolderId
//getfilter:assessment_id=description.AssessmentID
//getfilter:control_set_id=description.ControlSetID
type AuditManagerEvidenceDescription struct {
	Evidence     auditmanager.Evidence
	ControlSetID string
	AssessmentID string
}

//index:aws_auditmanager_evidencefolder
//getfilter:id=description.EvidenceFolder.Id
//getfilter:assessment_id=description.AssessmentID
//getfilter:control_set_id=description.ControlSetID
type AuditManagerEvidenceFolderDescription struct {
	EvidenceFolder auditmanager.AssessmentEvidenceFolder
	AssessmentID   string
}

//index:aws_auditmanager_framework
//getfilter:id=description.Framework.Id
//getfilter:region=metadata.Region
type AuditManagerFrameworkDescription struct {
	Framework auditmanager.Framework
}

// ===================  CloudControl ===================

//index:aws_cloudcontrol_resource
//getfilter:identifier=description.Resource.Identifier
type CloudControlResourceDescription struct {
	Resource cloudcontrol.ResourceDescription
}

// ===================  CloudSearch ===================

//index:aws_cloudsearch_domain
//getfilter:domain_name=description.DomainStatus.DomainName
type CloudSearchDomainDescription struct {
	DomainStatus cloudsearch.DomainStatus
}

// ===================  DLM ===================

//index:aws_dlm_lifecyclepolicy
//getfilter:id=description.LifecyclePolicy.PolicyId
type DLMLifecyclePolicyDescription struct {
	LifecyclePolicy dlm.LifecyclePolicy
}

// ===================  DocDB ===================

//index:aws_docdb_cluster
//getfilter:db_cluster_identifier=description.DBCluster.DBClusterIdentifier
type DocDBClusterDescription struct {
	DBCluster docdb.DBCluster
	Tags      []docdb.Tag
}

//index:aws_docdb_instance
//getfilter:db_instance_identifier=description.DBCluster.DBClusterIdentifier
type DocDBClusterInstanceDescription struct {
	DBInstance docdb.DBInstance
	Tags       []docdb.Tag
}

type DocDBClusterSnapshotDescription struct {
	DBClusterSnapshot docdb.DBClusterSnapshot
	Tags              []docdb.Tag
	Attributes        *docdb.DBClusterSnapshotAttributesResult
}

// ===================  Global Accelerator ===================

//index:aws_globalaccelerator_accelerator
//getfilter:arn=description.Accelerator.AcceleratorArn
type GlobalAcceleratorAcceleratorDescription struct {
	Accelerator           globalaccelerator.Accelerator
	AcceleratorAttributes *globalaccelerator.AcceleratorAttributes
	Tags                  []globalaccelerator.Tag
}

//index:aws_globalaccelerator_endpointgroup
//getfilter:arn=description.EndpointGroup.EndpointGroupArn
//listfilter:listener_arn=description.ListenerArn
type GlobalAcceleratorEndpointGroupDescription struct {
	EndpointGroup  globalaccelerator.EndpointGroup
	ListenerArn    string
	AcceleratorArn string
}

//index:aws_globalaccelerator_listener
//getfilter:arn=description.Listener.ListenerArn
//listfilter:accelerator_arn=description.AcceleratorArn
type GlobalAcceleratorListenerDescription struct {
	Listener       globalaccelerator.Listener
	AcceleratorArn string
}

// ===================  Glue ===================

//index:aws_glue_catalogdatabase
//getfilter:name=description.Database.Name
type GlueCatalogDatabaseDescription struct {
	Database glue.Database
}

//index:aws_glue_catalogtable
//getfilter:name=description.Table.Name
//getfilter:database_name=description.DatabaseName
//listfilter:catalog_id=description.Table.CatalogId
//listfilter:database_name=description.Table.DatabaseName
type GlueCatalogTableDescription struct {
	Table  glue.Table
	LfTags []lakeformationTypes.LFTagPair
}

//index:aws_glue_connection
//getfilter:name=description.Connection.Name
//listfilter:connection_type=description.Connection.ConnectionType
type GlueConnectionDescription struct {
	Connection glue.Connection
}

//index:aws_glue_crawler
//getfilter:name=description.Crawler.Name
type GlueCrawlerDescription struct {
	Crawler glue.Crawler
}

//index:aws_glue_datacatalogencryptionsettings
type GlueDataCatalogEncryptionSettingsDescription struct {
	DataCatalogEncryptionSettings glue.DataCatalogEncryptionSettings
}

//index:aws_glue_dataqualityruleset
//getfilter:name=description.DataQualityRuleset.Name
//listfilter:created_on=description.DataQualityRuleset.CreatedOn
//listfilter:last_modified_on=description.DataQualityRuleset.LastModifiedOn
type GlueDataQualityRulesetDescription struct {
	DataQualityRuleset glueop.GetDataQualityRulesetOutput
	RulesetRuleCount   *int32
}

//index:aws_glue_devendpoint
//getfilter:endpoint_name=description.DevEndpoint.EndpointName
type GlueDevEndpointDescription struct {
	DevEndpoint glue.DevEndpoint
}

//index:aws_glue_job
//getfilter:name=description.Job.Name
type GlueJobDescription struct {
	Job      glue.Job
	Bookmark *glue.JobBookmarkEntry
}

//index:aws_glue_securityconfiguration
//getfilter:name=description.SecurityConfiguration.Name
type GlueSecurityConfigurationDescription struct {
	SecurityConfiguration glue.SecurityConfiguration
}

// ===================  Health ===================

//index:aws_health_event
//listfilter:arn=description.Event.Arn
//listfilter:availability_zone=description.Event.AvailabilityZone
//listfilter:end_time=description.Event.EndTime
//listfilter:event_type_category=description.Event.EventTypeCategory
//listfilter:event_type_code=description.Event.EventTypeCode
//listfilter:last_updated_time=description.Event.LastUpdatedTime
//listfilter:service=description.Event.Service
//listfilter:start_time=description.Event.StartTime
//listfilter:status_code=description.Event.StatusCode
type HealthEventDescription struct {
	Event health.Event
}

//index:aws_health_affectedentity
type HealthAffectedEntityDescription struct {
	Entity health.AffectedEntity
}

// ===================  Identity Store ===================

//index:aws_identitystore_group
//getfilter:id=description.Group.GroupId
//getfilter:identity_store_id=description.Group.IdentityStoreId
//listfilter:identity_store_id=description.Group.IdentityStoreId
type IdentityStoreGroupDescription struct {
	Group identitystore.Group
}

//index:aws_identitystore_user
//getfilter:id=description.User.UserId
//getfilter:identity_store_id=description.User.IdentityStoreId
//listfilter:identity_store_id=description.User.IdentityStoreId
type IdentityStoreUserDescription struct {
	User         identitystore.User
	PrimaryEmail *string
}

//index:aws_identitystore_group
//getfilter:id=description.Group.GroupId
//getfilter:identity_store_id=description.Group.IdentityStoreId
//listfilter:identity_store_id=description.Group.IdentityStoreId
type IdentityStoreGroupMembershipDescription struct {
	MembershipId    *string
	IdentityStoreId *string
	GroupId         *string
	MemberId        interface{}
}

// ===================  Inspector ===================

//index:aws_inspector_assessmentrun
//listfilter:assessment_template_arn=description.AssessmentRun.AssessmentTemplateArn
//listfilter:name=description.AssessmentRun.Name
//listfilter:state=description.AssessmentRun.State
type InspectorAssessmentRunDescription struct {
	AssessmentRun inspector.AssessmentRun
}

//index:aws_inspector_assessmenttarget
//getfilter:arn=description.AssessmentTarget.Arn
type InspectorAssessmentTargetDescription struct {
	AssessmentTarget inspector.AssessmentTarget
}

//index:aws_inspector_assessmenttemplate
//getfilter:arn=description.AssessmentTemplate.Arn
//listfilter:name=description.AssessmentTemplate.Name
//listfilter:assessment_target_arn=description.AssessmentTemplate.AssessmentTargetArn
type InspectorAssessmentTemplateDescription struct {
	AssessmentTemplate inspector.AssessmentTemplate
	EventSubscriptions []inspector.Subscription
	Tags               []inspector.Tag
}

//index:aws_inspector_exclusion
//listfilter:assessment_run_arn=description.Exclusion.Arn
type InspectorExclusionDescription struct {
	Exclusion        inspector.Exclusion
	AssessmentRunArn string
}

//index:aws_inspector_finding
//listfilter:agent_id=description.Finding.AssetAttributes.AgentId
//listfilter:auto_scaling_group=description.Finding.AssetAttributes.AutoScalingGroup
//listfilter:severity=description.Finding.Severity
//getfilter:arn=description.Finding.Arn
type InspectorFindingDescription struct {
	Finding     inspector.Finding
	FailedItems map[string]inspector.FailedItemDetails
}

//index:aws_inspector2_coverage
//listfilter:account_id=description.CoveredResource.AccountId
//getfilter:resource_id=description.CoveredResource.ResourceId
type Inspector2CoverageDescription struct {
	CoveredResource inspector2.CoveredResource
}

//index:aws_inspector2_coveragestatistic
type Inspector2CoverageStatisticDescription struct {
	TotalCounts *int64
	Counts      []inspector2.Counts
}

//index:aws_inspector2_member
type Inspector2MemberDescription struct {
	Member         inspector2.Member
	OnlyAssociated bool
}

//index:aws_inspector2_member
type Inspector2FindingDescription struct {
	Finding  inspector2.Finding
	Resource inspector2.Resource
}

// ===================  Firehose  ===================

//index:aws_firehose_deliverystream
//getfilter:delivery_stream_name=description.DeliveryStream.DeliveryStreamName
//listfilter:delivery_stream_type=description.DeliveryStream.DeliveryStreamType
type FirehoseDeliveryStreamDescription struct {
	DeliveryStream firehose.DeliveryStreamDescription
	Tags           []firehose.Tag
}

// ===================  Lightsail ===================

//index:aws_lightsail_instance
//getfilter:name=description.Instance.
type LightsailInstanceDescription struct {
	Instance lightsail.Instance
}

// ===================  Macie2 ===================

//index:aws_macie2_classificationjob
//getfilter:job_id=description.ClassificationJob.JobId
//listfilter:name=description.ClassificationJob.Name
//listfilter:job_status=description.ClassificationJob.JobStatus
//listfilter:job_type=description.ClassificationJob.JobType
type Macie2ClassificationJobDescription struct {
	ClassificationJob macie2op.DescribeClassificationJobOutput
}

// ===================  MediaStore ===================

//index:aws_mediastore_container
//getfilter:name=description.Container.Name
type MediaStoreContainerDescription struct {
	Container mediastore.Container
	Policy    *mediastoreop.GetContainerPolicyOutput
	Tags      []mediastore.Tag
}

// ===================  MGN ===================

//index:aws_mgn_application
//listfilter:application_id=description.Application.ApplicationID
//listfilter:wave_id=description.Application.WaveID
//listfilter:is_archived=description.Application.IsArchived
type MgnApplicationDescription struct {
	Application mgn.Application
}

// ===================  SecurityLake ===================

//index:aws_securitylake_datalake
type SecurityLakeDataLakeDescription struct {
	DataLake securitylake.DataLakeResource
}

//index:aws_securitylake_subscriber
//getfilter:subscriber_id=description.Subscriber.SubscriberId
type SecurityLakeSubscriberDescription struct {
	Subscriber securitylake.SubscriberResource
}

// ===================  Ram  ===================

//index:aws_ram_principalassociation
type RamPrincipalAssociationDescription struct {
	PrincipalAssociation    ram.ResourceShareAssociation
	ResourceSharePermission []ram.ResourceSharePermissionSummary
}

//index:aws_ram_resourceassociation
type RamResourceAssociationDescription struct {
	ResourceAssociation     ram.ResourceShareAssociation
	ResourceSharePermission []ram.ResourceSharePermissionSummary
}

// ===================  Serverless Application Repository  ===================

//index:aws_serverlessapplicationrepository_application
//getfilter:arn=description.Application.ApplicationId
type ServerlessApplicationRepositoryApplicationDescription struct {
	Application serverlessapplicationrepositoryop.GetApplicationOutput
	Statements  []serverlessapplicationrepository.ApplicationPolicyStatement
}

// ===================  Service Quotas  ===================

//index:aws_servicequotas_defaultservicequota
//getfilter:quota_code=description.DefaultServiceQuota.QuotaCode
//getfilter:service_code=description.DefaultServiceQuota.ServiceCode
//listfilter:service_code=description.DefaultServiceQuota.ServiceCode
type ServiceQuotasDefaultServiceQuotaDescription struct {
	DefaultServiceQuota servicequotas.ServiceQuota
}

//index:aws_servicequotas_servicequota
//getfilter:quota_code=description.ServiceQuota.QuotaCode
//getfilter:service_code=description.ServiceQuota.ServiceCode
//listfilter:service_code=description.ServiceQuota.ServiceCode
type ServiceQuotasServiceQuotaDescription struct {
	ServiceQuota servicequotas.ServiceQuota
	Tags         []servicequotas.Tag
}

//index:aws_servicequotas_servicequotachangerequest
//getfilter:id=description.ServiceQuotaChangeRequest.Id
//listfilter:service_code=description.ServiceQuotaChangeRequest.ServiceCode
//listfilter:status=description.ServiceQuotaChangeRequest.Status
type ServiceQuotasServiceQuotaChangeRequestDescription struct {
	ServiceQuotaChangeRequest servicequotas.RequestedServiceQuotaChange
	Tags                      []servicequotas.Tag
}

type ServiceQuotasServiceDescription struct {
	Service servicequotas.ServiceInfo
}

// =================== Service Catalog =======================

//index:aws_servicecatalog_product
type ServiceCatalogProductDescription struct {
	ProductViewSummary    serviceCatalog.ProductViewSummary
	Budgets               []serviceCatalog.BudgetDetail
	LunchPaths            []serviceCatalog.LaunchPathSummary
	ProvisioningArtifacts []serviceCatalog.ProvisioningArtifactSummary
}

//index:aws_servicecatalog_portfolio
type ServiceCatalogPortfolioDescription struct {
	Portfolio serviceCatalog.PortfolioDetail
}

// =================== Service Discovery ===========================

//index:aws_service_discovery_service
type ServiceDiscoveryServiceDescription struct {
	Service serviceDiscovery.ServiceSummary
	Tags    []serviceDiscovery.Tag
}

//index:aws_service_discovery_namespace
type ServiceDiscoveryNamespaceDescription struct {
	Namespace serviceDiscovery.NamespaceSummary
	Tags      []serviceDiscovery.Tag
}

//index:aws_service_discovery_instance
type ServiceDiscoveryInstanceDescription struct {
	Instance serviceDiscovery.InstanceSummary
}

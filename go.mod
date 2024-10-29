module github.com/opengovern/og-describer-template

go 1.22.0

toolchain go1.22.5

require (
	github.com/aws/aws-lambda-go v1.42.0
	github.com/aws/aws-sdk-go v1.49.10
	github.com/aws/aws-sdk-go-v2 v1.30.3
	github.com/aws/aws-sdk-go-v2/config v1.27.23
	github.com/aws/aws-sdk-go-v2/credentials v1.17.23
	github.com/aws/aws-sdk-go-v2/service/accessanalyzer v1.26.5
	github.com/aws/aws-sdk-go-v2/service/account v1.14.5
	github.com/aws/aws-sdk-go-v2/service/acm v1.22.5
	github.com/aws/aws-sdk-go-v2/service/acmpca v1.25.5
	github.com/aws/aws-sdk-go-v2/service/amp v1.21.5
	github.com/aws/aws-sdk-go-v2/service/amplify v1.18.5
	github.com/aws/aws-sdk-go-v2/service/apigateway v1.21.5
	github.com/aws/aws-sdk-go-v2/service/apigatewayv2 v1.18.5
	github.com/aws/aws-sdk-go-v2/service/appconfig v1.26.5
	github.com/aws/aws-sdk-go-v2/service/applicationautoscaling v1.25.5
	github.com/aws/aws-sdk-go-v2/service/applicationinsights v1.22.5
	github.com/aws/aws-sdk-go-v2/service/appstream v1.30.0
	github.com/aws/aws-sdk-go-v2/service/athena v1.37.3
	github.com/aws/aws-sdk-go-v2/service/auditmanager v1.30.5
	github.com/aws/aws-sdk-go-v2/service/autoscaling v1.36.5
	github.com/aws/aws-sdk-go-v2/service/backup v1.31.2
	github.com/aws/aws-sdk-go-v2/service/batch v1.30.5
	github.com/aws/aws-sdk-go-v2/service/cloudcontrol v1.15.5
	github.com/aws/aws-sdk-go-v2/service/cloudformation v1.42.4
	github.com/aws/aws-sdk-go-v2/service/cloudfront v1.32.5
	github.com/aws/aws-sdk-go-v2/service/cloudsearch v1.20.5
	github.com/aws/aws-sdk-go-v2/service/cloudtrail v1.35.5
	github.com/aws/aws-sdk-go-v2/service/cloudwatch v1.32.0
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.30.0
	github.com/aws/aws-sdk-go-v2/service/codeartifact v1.23.5
	github.com/aws/aws-sdk-go-v2/service/codebuild v1.26.5
	github.com/aws/aws-sdk-go-v2/service/codecommit v1.19.5
	github.com/aws/aws-sdk-go-v2/service/codedeploy v1.22.1
	github.com/aws/aws-sdk-go-v2/service/codepipeline v1.22.5
	github.com/aws/aws-sdk-go-v2/service/codestar v1.19.5
	github.com/aws/aws-sdk-go-v2/service/configservice v1.43.5
	github.com/aws/aws-sdk-go-v2/service/costexplorer v1.33.5
	github.com/aws/aws-sdk-go-v2/service/databasemigrationservice v1.35.5
	github.com/aws/aws-sdk-go-v2/service/dax v1.17.5
	github.com/aws/aws-sdk-go-v2/service/directconnect v1.22.5
	github.com/aws/aws-sdk-go-v2/service/directoryservice v1.22.5
	github.com/aws/aws-sdk-go-v2/service/dlm v1.22.5
	github.com/aws/aws-sdk-go-v2/service/docdb v1.29.5
	github.com/aws/aws-sdk-go-v2/service/drs v1.21.0
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.26.6
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.18.5
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.141.0
	github.com/aws/aws-sdk-go-v2/service/ecr v1.24.5
	github.com/aws/aws-sdk-go-v2/service/ecrpublic v1.21.5
	github.com/aws/aws-sdk-go-v2/service/ecs v1.35.5
	github.com/aws/aws-sdk-go-v2/service/efs v1.26.5
	github.com/aws/aws-sdk-go-v2/service/eks v1.35.5
	github.com/aws/aws-sdk-go-v2/service/elasticache v1.34.5
	github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk v1.20.5
	github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing v1.21.5
	github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2 v1.26.5
	github.com/aws/aws-sdk-go-v2/service/elasticsearchservice v1.24.5
	github.com/aws/aws-sdk-go-v2/service/emr v1.35.5
	github.com/aws/aws-sdk-go-v2/service/eventbridge v1.26.5
	github.com/aws/aws-sdk-go-v2/service/firehose v1.23.0
	github.com/aws/aws-sdk-go-v2/service/fms v1.29.5
	github.com/aws/aws-sdk-go-v2/service/fsx v1.39.5
	github.com/aws/aws-sdk-go-v2/service/glacier v1.19.5
	github.com/aws/aws-sdk-go-v2/service/globalaccelerator v1.20.5
	github.com/aws/aws-sdk-go-v2/service/glue v1.72.4
	github.com/aws/aws-sdk-go-v2/service/grafana v1.18.5
	github.com/aws/aws-sdk-go-v2/service/guardduty v1.35.5
	github.com/aws/aws-sdk-go-v2/service/health v1.22.5
	github.com/aws/aws-sdk-go-v2/service/iam v1.28.5
	github.com/aws/aws-sdk-go-v2/service/identitystore v1.21.5
	github.com/aws/aws-sdk-go-v2/service/imagebuilder v1.30.0
	github.com/aws/aws-sdk-go-v2/service/inspector v1.19.5
	github.com/aws/aws-sdk-go-v2/service/inspector2 v1.20.5
	github.com/aws/aws-sdk-go-v2/service/kafka v1.28.5
	github.com/aws/aws-sdk-go-v2/service/keyspaces v1.7.5
	github.com/aws/aws-sdk-go-v2/service/kinesis v1.24.5
	github.com/aws/aws-sdk-go-v2/service/kinesisanalyticsv2 v1.21.5
	github.com/aws/aws-sdk-go-v2/service/kinesisvideo v1.21.5
	github.com/aws/aws-sdk-go-v2/service/kms v1.35.3
	github.com/aws/aws-sdk-go-v2/service/lakeformation v1.35.3
	github.com/aws/aws-sdk-go-v2/service/lambda v1.49.5
	github.com/aws/aws-sdk-go-v2/service/lightsail v1.32.5
	github.com/aws/aws-sdk-go-v2/service/macie2 v1.34.5
	github.com/aws/aws-sdk-go-v2/service/mediastore v1.18.5
	github.com/aws/aws-sdk-go-v2/service/memorydb v1.17.5
	github.com/aws/aws-sdk-go-v2/service/mgn v1.25.5
	github.com/aws/aws-sdk-go-v2/service/mq v1.20.6
	github.com/aws/aws-sdk-go-v2/service/mwaa v1.22.5
	github.com/aws/aws-sdk-go-v2/service/neptune v1.28.0
	github.com/aws/aws-sdk-go-v2/service/networkfirewall v1.36.5
	github.com/aws/aws-sdk-go-v2/service/oam v1.7.5
	github.com/aws/aws-sdk-go-v2/service/opensearch v1.27.0
	github.com/aws/aws-sdk-go-v2/service/opensearchserverless v1.9.5
	github.com/aws/aws-sdk-go-v2/service/opsworkscm v1.20.5
	github.com/aws/aws-sdk-go-v2/service/organizations v1.23.5
	github.com/aws/aws-sdk-go-v2/service/pinpoint v1.26.6
	github.com/aws/aws-sdk-go-v2/service/pipes v1.9.6
	github.com/aws/aws-sdk-go-v2/service/pricing v1.24.5
	github.com/aws/aws-sdk-go-v2/service/ram v1.23.6
	github.com/aws/aws-sdk-go-v2/service/rds v1.82.0
	github.com/aws/aws-sdk-go-v2/service/redshift v1.39.6
	github.com/aws/aws-sdk-go-v2/service/redshiftserverless v1.15.4
	github.com/aws/aws-sdk-go-v2/service/resourceexplorer2 v1.8.5
	github.com/aws/aws-sdk-go-v2/service/resourcegroups v1.19.5
	github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi v1.19.5
	github.com/aws/aws-sdk-go-v2/service/route53 v1.35.5
	github.com/aws/aws-sdk-go-v2/service/route53domains v1.20.5
	github.com/aws/aws-sdk-go-v2/service/route53resolver v1.23.5
	github.com/aws/aws-sdk-go-v2/service/s3 v1.47.5
	github.com/aws/aws-sdk-go-v2/service/s3control v1.41.5
	github.com/aws/aws-sdk-go-v2/service/sagemaker v1.121.0
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.25.5
	github.com/aws/aws-sdk-go-v2/service/securityhub v1.44.0
	github.com/aws/aws-sdk-go-v2/service/securitylake v1.10.5
	github.com/aws/aws-sdk-go-v2/service/serverlessapplicationrepository v1.18.5
	github.com/aws/aws-sdk-go-v2/service/servicecatalog v1.25.5
	github.com/aws/aws-sdk-go-v2/service/servicediscovery v1.27.5
	github.com/aws/aws-sdk-go-v2/service/servicequotas v1.19.5
	github.com/aws/aws-sdk-go-v2/service/ses v1.19.5
	github.com/aws/aws-sdk-go-v2/service/sesv2 v1.24.5
	github.com/aws/aws-sdk-go-v2/service/sfn v1.24.5
	github.com/aws/aws-sdk-go-v2/service/shield v1.23.5
	github.com/aws/aws-sdk-go-v2/service/simspaceweaver v1.8.5
	github.com/aws/aws-sdk-go-v2/service/sns v1.26.5
	github.com/aws/aws-sdk-go-v2/service/sqs v1.29.5
	github.com/aws/aws-sdk-go-v2/service/ssm v1.44.5
	github.com/aws/aws-sdk-go-v2/service/ssoadmin v1.23.5
	github.com/aws/aws-sdk-go-v2/service/storagegateway v1.24.5
	github.com/aws/aws-sdk-go-v2/service/sts v1.30.1
	github.com/aws/aws-sdk-go-v2/service/support v1.19.5
	github.com/aws/aws-sdk-go-v2/service/synthetics v1.22.5
	github.com/aws/aws-sdk-go-v2/service/timestreamwrite v1.23.6
	github.com/aws/aws-sdk-go-v2/service/waf v1.18.5
	github.com/aws/aws-sdk-go-v2/service/wafregional v1.19.5
	github.com/aws/aws-sdk-go-v2/service/wafv2 v1.43.5
	github.com/aws/aws-sdk-go-v2/service/wellarchitected v1.27.5
	github.com/aws/aws-sdk-go-v2/service/workspaces v1.35.6
	github.com/aws/smithy-go v1.20.3
	github.com/ghodss/yaml v1.0.0
	github.com/go-errors/errors v1.4.2
	github.com/gocarina/gocsv v0.0.0-20211203214250-4735fba0c1d9
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/golang/protobuf v1.5.4
	github.com/hashicorp/go-hclog v1.6.3
	github.com/labstack/echo/v4 v4.12.0
	github.com/manifoldco/promptui v0.9.0
	github.com/nats-io/nats.go v1.36.0
	github.com/opengovern/og-util v1.0.4
	github.com/spf13/cobra v1.7.0
	github.com/turbot/go-kit v0.9.0
	github.com/turbot/steampipe-plugin-sdk/v5 v5.8.0
	go.uber.org/zap v1.26.0
	golang.org/x/oauth2 v0.21.0
	golang.org/x/time v0.5.0
	golang.org/x/tools v0.21.1-0.20240508182429-e35e4ccd0d2d
	google.golang.org/genproto v0.0.0-20240227224415-6ceb2ff114de
	google.golang.org/grpc v1.65.0
	google.golang.org/protobuf v1.34.2
)

require (
	cloud.google.com/go v0.112.1 // indirect
	cloud.google.com/go/compute/metadata v0.3.0 // indirect
	cloud.google.com/go/iam v1.1.6 // indirect
	cloud.google.com/go/longrunning v0.5.5 // indirect
	cloud.google.com/go/storage v1.38.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.11.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.7.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.8.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets v1.1.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/internal v1.0.0 // indirect
	github.com/AzureAD/microsoft-authentication-library-for-go v1.2.2 // indirect
	github.com/acarl005/stripansi v0.0.0-20180116102854-5a71ef0e047d // indirect
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/allegro/bigcache/v3 v3.1.0 // indirect
	github.com/apparentlymart/go-textseg/v15 v15.0.0 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.5.4 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.9 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.15 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.15 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.0 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.2.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.11.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.2.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.8.10 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.11.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.16.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.22.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.26.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bgentry/go-netrc v0.0.0-20140422174119-9fd32a8b3d3d // indirect
	github.com/btubbs/datetime v0.1.1 // indirect
	github.com/cenkalti/backoff/v3 v3.2.2 // indirect
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e // indirect
	github.com/danwakefield/fnmatch v0.0.0-20160403171240-cbb64ac3d964 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dgraph-io/ristretto v0.1.1 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/eko/gocache/lib/v4 v4.1.5 // indirect
	github.com/eko/gocache/store/bigcache/v4 v4.2.1 // indirect
	github.com/eko/gocache/store/ristretto/v4 v4.2.1 // indirect
	github.com/elastic/go-elasticsearch/v7 v7.17.10 // indirect
	github.com/fatih/color v1.17.0 // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/gertd/go-pluralize v0.2.1 // indirect
	github.com/globocom/echo-prometheus v0.1.2 // indirect
	github.com/go-jose/go-jose/v4 v4.0.1 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-viper/mapstructure/v2 v2.0.0-alpha.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/golang/glog v1.2.1 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/s2a-go v0.1.7 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.2 // indirect
	github.com/googleapis/gax-go/v2 v2.12.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.18.1 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-getter v1.7.5 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-plugin v1.6.0 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.7 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/go-safetemp v1.0.0 // indirect
	github.com/hashicorp/go-secure-stdlib/parseutil v0.1.6 // indirect
	github.com/hashicorp/go-secure-stdlib/strutil v0.1.2 // indirect
	github.com/hashicorp/go-sockaddr v1.0.2 // indirect
	github.com/hashicorp/go-version v1.6.0 // indirect
	github.com/hashicorp/hcl v1.0.1-vault // indirect
	github.com/hashicorp/hcl/v2 v2.20.1 // indirect
	github.com/hashicorp/vault/api v1.14.0 // indirect
	github.com/hashicorp/vault/api/auth/kubernetes v0.7.0 // indirect
	github.com/hashicorp/yamux v0.1.1 // indirect
	github.com/iancoleman/strcase v0.3.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.14.3 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.3 // indirect
	github.com/jackc/pgservicefile v0.0.0-20231201235250-de7065d80cb9 // indirect
	github.com/jackc/pgtype v1.14.0 // indirect
	github.com/jackc/pgx/v4 v4.18.3 // indirect
	github.com/jackc/pgx/v5 v5.5.5 // indirect
	github.com/jackc/puddle v1.3.0 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.17.4 // indirect
	github.com/knadh/koanf/maps v0.1.1 // indirect
	github.com/knadh/koanf/parsers/toml v0.1.0 // indirect
	github.com/knadh/koanf/providers/env v0.1.0 // indirect
	github.com/knadh/koanf/providers/file v0.1.0 // indirect
	github.com/knadh/koanf/providers/structs v0.1.0 // indirect
	github.com/knadh/koanf/v2 v2.1.0 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/oklog/run v1.1.0 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/opensearch-project/opensearch-go/v2 v2.3.0 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pganalyze/pg_query_go/v4 v4.2.3 // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_golang v1.19.1 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.53.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/rabbitmq/amqp091-go v1.9.0 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/ryanuber/go-glob v1.0.0 // indirect
	github.com/sethvargo/go-retry v0.2.4 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stevenle/topsort v0.2.0 // indirect
	github.com/tkrajina/go-reflector v0.5.6 // indirect
	github.com/ulikunitz/xz v0.5.11 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/zclconf/go-cty v1.14.4 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho v0.53.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.49.0 // indirect
	go.opentelemetry.io/otel v1.28.0 // indirect
	go.opentelemetry.io/otel/exporters/jaeger v1.17.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v0.44.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.21.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.21.0 // indirect
	go.opentelemetry.io/otel/metric v1.28.0 // indirect
	go.opentelemetry.io/otel/sdk v1.28.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.28.0 // indirect
	go.opentelemetry.io/otel/trace v1.28.0 // indirect
	go.opentelemetry.io/proto/otlp v1.0.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.25.0 // indirect
	golang.org/x/exp v0.0.0-20240112132812-db7319d0e0e3 // indirect
	golang.org/x/mod v0.17.0 // indirect
	golang.org/x/net v0.27.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.22.0 // indirect
	golang.org/x/term v0.22.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/api v0.169.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240528184218-531527333157 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240610135401-a8a62080eff3 // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/apimachinery v0.30.2 // indirect
	k8s.io/client-go v0.30.2 // indirect
	k8s.io/klog/v2 v2.120.1 // indirect
	k8s.io/utils v0.0.0-20240502163921-fe8a2dddb1d0 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.4.1 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)

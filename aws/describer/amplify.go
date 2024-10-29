package describer

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/amplify/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/amplify"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func AmplifyApp(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	client := amplify.NewFromConfig(cfg)

	var values []Resource
	err := PaginateRetrieveAll(func(prevToken *string) (nextToken *string, err error) {
		output, err := client.ListApps(ctx, &amplify.ListAppsInput{
			MaxResults: 100,
			NextToken:  prevToken,
		})
		if err != nil {
			return nil, err
		}

		for _, item := range output.Apps {
			resource := amplifyAppHandle(ctx, item)
			if stream != nil {
				m := *stream
				err := m(resource)
				if err != nil {
					return nil, err
				}
			} else {
				values = append(values, resource)
			}
		}
		return output.NextToken, nil
	})
	if err != nil {
		return nil, err
	}

	return values, nil
}
func amplifyAppHandle(ctx context.Context, item types.App) Resource {
	describeCtx := GetDescribeContext(ctx)
	resource := Resource{
		Region: describeCtx.KaytuRegion,
		Name:   *item.Name,
		ARN:    *item.AppArn,
		ID:     *item.AppId,
		Description: model.AmplifyAppDescription{
			App: item,
		},
	}
	return resource
}
func GetAmplifyApp(ctx context.Context, cfg aws.Config, fields map[string]string) ([]Resource, error) {
	appId := fields["appId"]
	client := amplify.NewFromConfig(cfg)

	out, err := client.ListApps(ctx, &amplify.ListAppsInput{})
	if err != nil {
		if isErr(err, "ListAppsNotFound") || isErr(err, "InvalidParameterValue") {
			return nil, nil
		}
		return nil, err
	}

	var values []Resource
	for _, app := range out.Apps {
		if *app.AppId != appId {
			continue
		}
		resource := amplifyAppHandle(ctx, app)
		values = append(values, resource)
	}

	return values, nil
}

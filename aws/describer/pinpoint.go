package describer

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/pinpoint"
	"github.com/aws/aws-sdk-go-v2/service/pinpoint/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func PinpointApp(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)
	client := pinpoint.NewFromConfig(cfg)

	input := &pinpoint.GetAppsInput{
		PageSize: aws.String("1000"),
	}

	var values []Resource
	for {
		apps, err := client.GetApps(ctx, input)
		if err != nil {
			return nil, err
		}

		if apps.ApplicationsResponse == nil {
			break
		}

		for _, app := range apps.ApplicationsResponse.Item {
			op, err := client.GetApplicationSettings(ctx, &pinpoint.GetApplicationSettingsInput{
				ApplicationId: app.Id,
			})
			if err != nil {
				return nil, err
			}

			var settings *types.ApplicationSettingsResource
			if op != nil {
				settings = op.ApplicationSettingsResource
			}

			var name string
			if app.Name != nil {
				name = *app.Name
			}

			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ARN:    *app.Arn,
				Name:   name,
				Description: model.PinPointAppDescription{
					App:      app,
					Settings: settings,
				},
			}
			if stream != nil {
				if err := (*stream)(resource); err != nil {
					return nil, err
				}
			} else {
				values = append(values, resource)
			}
		}

		if apps.ApplicationsResponse.NextToken != nil {
			input.Token = apps.ApplicationsResponse.NextToken
		} else {
			break
		}
	}

	return values, nil
}

package describer

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/aws-sdk-go-v2/service/sfn/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/opengovern/og-aws-describer/aws/model"
)

func StepFunctionsStateMachine(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)

	client := sfn.NewFromConfig(cfg)
	paginator := sfn.NewListStateMachinesPaginator(client, &sfn.ListStateMachinesInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.StateMachines {
			data, err := client.DescribeStateMachine(ctx, &sfn.DescribeStateMachineInput{
				StateMachineArn: v.StateMachineArn,
			})
			if err != nil {
				return nil, err
			}

			var name string
			if data.Name != nil {
				name = *data.Name
			}

			tags, err := client.ListTagsForResource(ctx, &sfn.ListTagsForResourceInput{
				ResourceArn: v.StateMachineArn,
			})
			if err != nil {
				return nil, err
			}

			if data.Definition != nil && len(*data.Definition) > 5000 {
				v := *data.Definition
				data.Definition = aws.String(v[:5000])
			}

			resource := Resource{
				Region: describeCtx.KaytuRegion,
				ARN:    *v.StateMachineArn,
				Name:   name,
				Description: model.StepFunctionsStateMachineDescription{
					StateMachineItem: v,
					StateMachine:     data,
					Tags:             tags.Tags,
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
	}

	return values, nil
}

type historyInfo struct {
	types.HistoryEvent
	ExecutionArn string
}

func StepFunctionsStateMachineExecutionHistories(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)

	client := sfn.NewFromConfig(cfg)
	paginator := sfn.NewListStateMachinesPaginator(client, &sfn.ListStateMachinesInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.StateMachines {
			stateMachineArn := v.StateMachineArn
			var executions []types.ExecutionListItem
			input := &sfn.ListExecutionsInput{
				StateMachineArn: stateMachineArn,
			}
			execPaginator := sfn.NewListExecutionsPaginator(client, input)
			for execPaginator.HasMorePages() {
				output, err := execPaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}
				executions = append(executions, output.Executions...)
			}
			if err != nil {
				return nil, err
			}

			var wg sync.WaitGroup
			executionCh := make(chan []historyInfo, len(executions))
			errorCh := make(chan error, len(executions))

			// Iterating all the available executions
			for _, item := range executions {
				wg.Add(1)
				item := item
				go func() {
					defer wg.Done()

					var items []historyInfo
					listHistory, err := client.GetExecutionHistory(ctx, &sfn.GetExecutionHistoryInput{
						ExecutionArn: item.ExecutionArn,
					})
					if err != nil {
						errorCh <- err
						return
					}

					for _, event := range listHistory.Events {
						items = append(items, historyInfo{event, *item.ExecutionArn})
					}
					executionCh <- items
				}()
			}

			// wait for all executions to be processed
			wg.Wait()
			close(executionCh)
			close(errorCh)

			for err := range errorCh {
				return nil, err
			}

			for item := range executionCh {
				for _, data := range item {
					resource := Resource{
						Region: describeCtx.KaytuRegion,
						ARN:    data.ExecutionArn,
						Name:   data.ExecutionArn,
						Description: model.StepFunctionsStateMachineExecutionHistoriesDescription{
							ExecutionHistory: data.HistoryEvent,
							ARN:              data.ExecutionArn,
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
			}
		}
	}

	return values, nil
}

func StepFunctionsStateMachineExecution(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	describeCtx := GetDescribeContext(ctx)

	client := sfn.NewFromConfig(cfg)
	paginator := sfn.NewListStateMachinesPaginator(client, &sfn.ListStateMachinesInput{})

	var values []Resource
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.StateMachines {
			arn := v.StateMachineArn
			executionPaginator := sfn.NewListExecutionsPaginator(client, &sfn.ListExecutionsInput{
				StateMachineArn: arn,
			})

			for executionPaginator.HasMorePages() {
				output, err := executionPaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}
				for _, execution := range output.Executions {
					data, err := client.DescribeExecution(ctx, &sfn.DescribeExecutionInput{
						ExecutionArn: execution.ExecutionArn,
					})
					if err != nil {
						if isErr(err, "ExecutionDoesNotExist") {
							continue
						}
						return nil, err
					}

					if data.Input != nil && len(*data.Input) > 5000 {
						v := *data.Input
						data.Input = aws.String(v[:5000])
					}

					if data.Output != nil && len(*data.Output) > 5000 {
						v := *data.Output
						data.Output = aws.String(v[:5000])
					}

					resource := Resource{
						Region: describeCtx.KaytuRegion,
						ARN:    *execution.ExecutionArn,
						Name:   *execution.ExecutionArn,
						Description: model.StepFunctionsStateMachineExecutionDescription{
							ExecutionItem: execution,
							Execution:     data,
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
			}
		}
	}

	return values, nil
}

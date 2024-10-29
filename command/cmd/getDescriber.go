/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/opengovern/og-aws-describer/aws"
	"github.com/opengovern/og-util/pkg/describe/enums"
	"github.com/spf13/cobra"
)

var (
	id, name, arn string
)

// getDescriberCmd represents the getDescriber command
var getDescriberCmd = &cobra.Command{
	Use:   "getDescriber",
	Short: "A brief description of your command",
	RunE: func(cmd *cobra.Command, args []string) error {
		fields := map[string]string{
			"id":   id,
			"name": name,
			"arn":  arn,
		}

		output, err := aws.GetSingleResource(
			context.Background(),
			resourceType,
			enums.DescribeTriggerTypeManual,
			accountID,
			nil,
			accessKey,
			secretKey,
			"",
			"",
			nil,
			false,
			fields,
		)
		if err != nil {
			return fmt.Errorf("AWS: %w", err)
		}

		js, err := json.Marshal(output)
		if err != nil {
			return err
		}

		fmt.Println(string(js))
		return nil
	},
}

func init() {
	getDescriberCmd.Flags().StringVar(&id, "id", "", "id")
	getDescriberCmd.Flags().StringVar(&name, "name", "", "name")
	getDescriberCmd.Flags().StringVar(&arn, "arn", "", "arn")
	getDescriberCmd.Flags().StringVar(&resourceType, "resourceType", "", "resourceType")
	getDescriberCmd.Flags().StringVar(&accountID, "accountID", "", "AccountID")
	getDescriberCmd.Flags().StringVar(&accessKey, "accessKey", "", "Access key")
	getDescriberCmd.Flags().StringVar(&secretKey, "secretKey", "", "Secret key")
}

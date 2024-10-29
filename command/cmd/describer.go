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
	"go.uber.org/zap"
)

var (
	checkAttachedPolicies                                                                         bool
	resourceType, accountID, accessKey, credentialAccountId, secretKey, assumeRoleArn, externalId string
)

// describerCmd represents the describer command
var describerCmd = &cobra.Command{
	Use:   "describer",
	Short: "A brief description of your command",
	RunE: func(cmd *cobra.Command, args []string) error {
		logger, _ := zap.NewProduction()
		externalIdPtr := &externalId
		if externalId == "" {
			externalIdPtr = nil
		}

		logger.Info("getting config")
		cfg, err := aws.GetConfig(context.Background(), accessKey, secretKey, "", assumeRoleArn, externalIdPtr)
		if err != nil {
			return fmt.Errorf("AWS: %w", err)
		}
		logger.Info("got config")
		if checkAttachedPolicies {
			isAttached, err := aws.CheckAttachedPolicy(logger, cfg, "", aws.SecurityAuditPolicyARN)
			fmt.Println("IsAttached", isAttached)
			fmt.Println("Error", err)
			return nil
		}
		output, err := aws.GetResources(
			context.Background(), logger,
			resourceType, enums.DescribeTriggerTypeManual,
			accountID, nil,
			credentialAccountId, accessKey, secretKey, "", assumeRoleArn, "", externalIdPtr,
			false, nil)
		if err != nil {
			return fmt.Errorf("AWS: %w", err)
		}
		logger.Info("got resources")
		js, err := json.Marshal(output)
		if err != nil {
			return err
		}
		fmt.Println(string(js))
		return nil
	},
}

func init() {
	describerCmd.Flags().BoolVar(&checkAttachedPolicies, "checkAttachedPolicies", false, "Check attached policies")
	describerCmd.Flags().StringVar(&resourceType, "resourceType", "", "Resource type")
	describerCmd.Flags().StringVar(&accountID, "accountID", "", "AccountID")
	describerCmd.Flags().StringVar(&accessKey, "accessKey", "", "Access key")
	describerCmd.Flags().StringVar(&secretKey, "secretKey", "", "Secret key")
	describerCmd.Flags().StringVar(&assumeRoleArn, "assumeRoleName", "", "Assume role name")
	describerCmd.Flags().StringVar(&externalId, "externalId", "", "externalId")
	describerCmd.Flags().StringVar(&credentialAccountId, "credentialAccountId", "", "Credential account id")
}

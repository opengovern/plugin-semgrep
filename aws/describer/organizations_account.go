package describer

import (
	"context"
	"errors"
	"fmt"
	"github.com/opengovern/og-aws-describer/aws/model"
	"math/rand"
	"os"
	"sync"
	"time"

	// AWS SDK for Go V2 packages
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	orgtypes "github.com/aws/aws-sdk-go-v2/service/organizations/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"

	// Smithy error handling
	"github.com/aws/smithy-go"

	// Rate limiting
	"golang.org/x/time/rate"
)

// Constants for concurrency and API rate limits
const (
	APIRateLimit = 10 // Effective refill rate across all workers (API calls per second)
	MaxRetries   = 5  // Maximum number of retries for API calls
)

// Global rate limiter
var limiter = rate.NewLimiter(APIRateLimit, APIRateLimit)

// Shared variable to control worker pausing
var pauseWorkersChan = make(chan struct{})

// GetAWSConfig creates an AWS configuration using the provided credentials and role assumption if specified.
func GetAWSConfig(ctx context.Context, AccessKey, SecretKey, SessionToken, RoleToAssume string) (aws.Config, error) {
	var cfg aws.Config
	var err error

	// Create custom credentials provider if AccessKey and SecretKey are provided
	var creds aws.CredentialsProvider
	if AccessKey != "" && SecretKey != "" {
		creds = credentials.NewStaticCredentialsProvider(AccessKey, SecretKey, SessionToken)
	}

	// Options for loading the config
	var loadOptions []func(*config.LoadOptions) error

	// If custom credentials provider is set, add it to load options
	if creds != nil {
		loadOptions = append(loadOptions, config.WithCredentialsProvider(creds))
	}

	// Load the default config
	cfg, err = config.LoadDefaultConfig(ctx, loadOptions...)
	if err != nil {
		return cfg, fmt.Errorf("failed to load default config: %v", err)
	}

	// If RoleToAssume is provided, assume the role
	if RoleToAssume != "" {
		// Create an STS client from the config
		stsClient := sts.NewFromConfig(cfg)

		// Create an AssumeRole credentials provider
		provider := stscreds.NewAssumeRoleProvider(stsClient, RoleToAssume)

		// Update the config with the new credentials provider
		cfg.Credentials = aws.NewCredentialsCache(provider)
	}

	return cfg, nil
}

// IsManagementAccount checks if the current AWS account is the AWS Organization management account.
func IsManagementAccount(ctx context.Context, cfg aws.Config) error {
	stsClient := sts.NewFromConfig(cfg)
	callerIdentity, err := stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return fmt.Errorf("failed to get caller identity: %v", err)
	}
	accountID := aws.ToString(callerIdentity.Account)

	orgClient := organizations.NewFromConfig(cfg)
	orgInfo, err := orgClient.DescribeOrganization(ctx, &organizations.DescribeOrganizationInput{})
	if err != nil {
		return fmt.Errorf("failed to describe organization: %v", err)
	}
	masterAccountID := aws.ToString(orgInfo.Organization.MasterAccountId)

	if accountID != masterAccountID {
		return fmt.Errorf("the account %s is not an AWS Organization management account", accountID)
	}

	return nil
}

// GetOrganizationAccounts fetches all organization accounts with their details.
func GetOrganizationAccounts(ctx context.Context, cfg aws.Config) ([]Resource, error) {
	orgClient := organizations.NewFromConfig(cfg)

	var accounts []Resource
	var mu sync.Mutex // Protects the accounts slice

	input := &organizations.ListAccountsInput{}

	// Concurrency control
	var wg sync.WaitGroup

	// Channel to control worker concurrency
	workerChan := make(chan struct{}, APIRateLimit)

	// Pagination control with paginator
	paginator := organizations.NewListAccountsPaginator(orgClient, input)

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list accounts: %v", err)
		}

		for _, acct := range page.Accounts {
			wg.Add(1)
			workerChan <- struct{}{} // Acquire worker slot
			go func(acct orgtypes.Account) {
				defer wg.Done()
				defer func() { <-workerChan }() // Release worker slot

				// Wait if workers are paused
				select {
				case <-pauseWorkersChan:
					time.Sleep(5 * time.Second)
				default:
				}

				details, err := processAccount(ctx, orgClient, acct)
				if err != nil {
					// Handle error (e.g., log it)
					fmt.Fprintf(os.Stderr, "Error processing account %s: %v\n", aws.ToString(acct.Id), err)
					return
				}
				mu.Lock()
				accounts = append(accounts, details)
				mu.Unlock()
			}(acct)
		}
	}

	wg.Wait() // Wait for all goroutines to finish

	return accounts, nil
}

// processAccount processes individual account details.
func processAccount(ctx context.Context, orgClient *organizations.Client, acct orgtypes.Account) (Resource, error) {
	describeCtx := GetDescribeContext(ctx)

	// Get Tags
	var tagsOutput *organizations.ListTagsForResourceOutput
	err := callWithRetry(ctx, func() error {
		return limiter.Wait(ctx)
	})
	if err != nil {
		return Resource{}, err
	}
	err = callWithRetry(ctx, func() error {
		var err error
		tagsOutput, err = orgClient.ListTagsForResource(ctx, &organizations.ListTagsForResourceInput{
			ResourceId: acct.Id,
		})
		return err
	})
	if err != nil {
		return Resource{}, err
	}

	// Get Parent OU ID
	var parentsOutput *organizations.ListParentsOutput
	err = callWithRetry(ctx, func() error {
		return limiter.Wait(ctx)
	})
	if err != nil {
		return Resource{}, err
	}
	err = callWithRetry(ctx, func() error {
		var err error
		parentsOutput, err = orgClient.ListParents(ctx, &organizations.ListParentsInput{
			ChildId: acct.Id,
		})
		return err
	})
	if err != nil {
		return Resource{}, err
	}
	if len(parentsOutput.Parents) == 0 {
		// No parent found
		return Resource{}, fmt.Errorf("no parent found for account %s", aws.ToString(acct.Id))
	}

	parent := parentsOutput.Parents[0]
	ouId := aws.ToString(parent.Id)

	details := Resource{
		Region: describeCtx.KaytuRegion,
		ARN:    *acct.Arn,
		ID:     *acct.Id,
		Name:   *acct.Name,
		Description: model.OrganizationsAccountDescription{
			Account:  acct,
			Tags:     tagsOutput.Tags,
			ParentID: ouId,
		},
	}

	return details, nil
}

// callWithRetry handles API calls with retries and exponential backoff with jitter.
func callWithRetry(ctx context.Context, apiCall func() error) error {
	var err error
	var baseDelay = time.Second
	for i := 0; i < MaxRetries; i++ {
		// Wait for rate limiter
		err = limiter.Wait(ctx)
		if err != nil {
			return err
		}

		// Make the API call
		err = apiCall()
		if err == nil {
			return nil // success
		}

		// Check if error is throttling
		if isThrottlingError(err) {
			// Pause all workers for 5 seconds and retry at 50% rate
			fmt.Fprintln(os.Stderr, "Throttling detected. Pausing all workers for 5 seconds.")
			pauseAllWorkers()

			// Reduce rate limiter to 50%
			currentLimit := limiter.Limit()
			newLimit := currentLimit / 2
			if newLimit < 1 {
				newLimit = 1 // Do not reduce below 1 request per second
			}
			limiter.SetLimit(newLimit)
			fmt.Fprintf(os.Stderr, "Reducing API rate limit to %v requests per second.\n", newLimit)
		}

		// Exponential backoff with jitter
		jitter := time.Duration(rand.Int63n(int64(baseDelay)))
		sleep := (1 << i) * baseDelay
		sleep = sleep + jitter
		if sleep > 30*time.Second {
			sleep = 30 * time.Second // Maximum delay interval
		}
		time.Sleep(sleep)
	}
	return fmt.Errorf("max retries reached: %v", err)
}

// isThrottlingError checks if the error is due to API throttling.
func isThrottlingError(err error) bool {
	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		code := apiErr.ErrorCode()
		if code == "ThrottlingException" || code == "ThrottledException" || code == "TooManyRequestsException" {
			return true
		}
	}
	return false
}

// pauseAllWorkers signals all workers to pause.
func pauseAllWorkers() {
	select {
	case pauseWorkersChan <- struct{}{}:
	default:
	}
	time.Sleep(5 * time.Second)
	// Increase the rate limiter after pause
	currentLimit := limiter.Limit()
	newLimit := currentLimit * 2
	if newLimit > APIRateLimit {
		newLimit = APIRateLimit
	}
	limiter.SetLimit(newLimit)
	fmt.Fprintf(os.Stderr, "Increasing API rate limit to %v requests per second.\n", newLimit)
}

func OrganizationsAccount(ctx context.Context, cfg aws.Config, stream *StreamSender) ([]Resource, error) {
	// Seed the random number generator for jitter
	rand.Seed(time.Now().UnixNano())

	// Check if the account is a management account
	err := IsManagementAccount(ctx, cfg)
	if err != nil {
		return nil, err
	}

	// Get the organization accounts
	accounts, err := GetOrganizationAccounts(ctx, cfg)
	if err != nil {
		return nil, err
	}

	values := make([]Resource, len(accounts))

	for _, resource := range accounts {
		if stream != nil {
			if err := (*stream)(resource); err != nil {
				return nil, err
			}
		} else {
			values = append(values, resource)
		}
	}

	return values, nil
}

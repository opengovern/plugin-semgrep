package describer

import (
	"context"
	"encoding/json"
	// "fmt"

	// "fmt"
	"net/http"
	// "net/url"
	"sync"
	"golang.org/x/time/rate"
	"time"
	"github.com/opengovern/og-describer-template/pkg/sdk/models"
	"github.com/opengovern/og-describer-template/provider/model"
)

type CohereAIAPIHandler struct {
	Client       *http.Client
	APIKey        string
	RateLimiter  *rate.Limiter
	Semaphore    chan struct{}
	MaxRetries   int
	RetryBackoff time.Duration
	ClientName   string
}



func ListDatasets(ctx context.Context, handler *CohereAIAPIHandler, stream *models.StreamSender) ([]models.Resource, error) {
	var wg sync.WaitGroup
	cohereaiChan := make(chan models.Resource)
	

	go func() {
		processDatasets(ctx, handler, cohereaiChan, &wg)
		wg.Wait()
		close(cohereaiChan)
	}()
	var values []models.Resource
	for value := range cohereaiChan {
		if stream != nil {
			if err := (*stream)(value); err != nil {
				return nil, err
			}
		} else {
			values = append(values, value)
		}
	}
	return values, nil
}

func processDatasets(ctx context.Context, handler *CohereAIAPIHandler, cohereAiChan chan<- models.Resource, wg *sync.WaitGroup) {
	var datasetResponse model.DatasetListResponse
	var datasets []model.DatasetDescription
	var resp *http.Response
	baseURL := "https://api.cohere.com/v1/datasets"
	
	
	finalURL := baseURL
	req, err := http.NewRequest("GET", finalURL, nil)
	if err != nil {
		return 
	}
	requestFunc := func(req *http.Request) (*http.Response, error) {
		var e error
		
		resp, e = handler.Client.Do(req)
		// fmt.Printf(json.NewDecoder(resp.Body))
		if e = json.NewDecoder(resp.Body).Decode(&datasetResponse); e != nil {
			return nil, e
		}
		datasets =append(datasets, datasetResponse.Datasets...)
		return resp, nil
	}
	err = handler.DoRequest(ctx, req, requestFunc)
	if err != nil {
		return
	}
	// get dataset usage for each dataset
	finalURL1 := baseURL +  "/usage"
	req1, err1:= http.NewRequest("GET", finalURL1, nil)
	var usage model.OrganizationUsage
	requestFunc1 := func(req *http.Request) (*http.Response, error) {
			var e error
			resp, e = handler.Client.Do(req)
			if e = json.NewDecoder(resp.Body).Decode(usage); e != nil {
				return nil, e
			}
			return resp, nil
		}
	err = handler.DoRequest(ctx, req1, requestFunc1)
	if err1 != nil {
		return
}
	for _, dataset := range datasets {
		dataset.TotalUsage = float64(usage.OrganizationUsage)
	}


	
	for _, dataset := range datasets {
		wg.Add(1)
		go func(dataset model.DatasetDescription) {
			defer wg.Done()
			value := models.Resource{
				ID:   dataset.ID,
				Name: dataset.Name,
				Description: JSONAllFieldsMarshaller{
					Value: dataset,
				},
			}
			cohereAiChan <- value
		}(dataset)
	}
}


func GetDataset(ctx context.Context, handler *CohereAIAPIHandler, datasetID string) (*models.Resource, error) {
	var datasetResponse model.DatasetDescription
	baseURL := "https://api.cohere.com/v1/datasets"
	
	
	finalURL := baseURL + "/" + datasetID
	req, err := http.NewRequest("GET", finalURL, nil)
	if err != nil {
		return &models.Resource{}, err
	}
	requestFunc := func(req *http.Request) (*http.Response, error) {
		var e error
		
		resp, e := handler.Client.Do(req)
		if e = json.NewDecoder(resp.Body).Decode(&datasetResponse); e != nil {
			return nil, e
		}
		
		return resp, e
	}
	err = handler.DoRequest(ctx, req, requestFunc)
	if err != nil {
		return &models.Resource{}, err
	}
	return &models.Resource{
		ID:   datasetResponse.ID,
		Name: datasetResponse.Name,
		Description: JSONAllFieldsMarshaller{
			Value: datasetResponse,
		},
	}, nil
}

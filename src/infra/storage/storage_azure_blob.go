package storage

import (
	"context"
	"encoding/json"
	"errors"
	"eventstore/config"
	"eventstore/core"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/appendblob"
	bloom "github.com/bits-and-blooms/bloom/v3"
)

type AzureBlobRepository struct {
	client *appendblob.Client
	cfg    *config.Config
}

// ich brauche noch einen context, den ich übergeben kann
func NewAzureBlobStorageRepository(cfg *config.Config) (*AzureBlobRepository, error) {
	blobUrl := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", cfg.AzureConfig.AccountName, cfg.AzureConfig.ContainerName, cfg.AzureConfig.BlobName)

	creds, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create azure credential: %v", err)
	}

	client, err := appendblob.NewClient(blobUrl, creds, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Azure Append Blob client: %v", err)
	}

	ctx := context.Background()
	if _, err = client.GetProperties(ctx, nil); err != nil {
		var responseError *azcore.ResponseError
		if errors.As(err, &responseError) && responseError.StatusCode == http.StatusNotFound {
			if _, err := client.Create(ctx, nil); err != nil {
				return nil, fmt.Errorf("failed to create append blob: %v", err)
			}
			log.Println("append blob created")
		} else {
			return nil, fmt.Errorf("failed to get blob properties: %v", err)
		}
	}

	return &AzureBlobRepository{
		client: client,
		cfg:    cfg,
	}, nil
}

// AppendEvent stores a new event in Azure Blob Storage by appending it to the blob
func (r *AzureBlobRepository) Append(ctx context.Context, event core.Event) error {
	return r.AppendAll(ctx, []core.Event{event})
}

func (r *AzureBlobRepository) AppendAll(ctx context.Context, events []core.Event) error {
	var b []byte

	for _, event := range events {
		record, err := json.Marshal(event)
		if err != nil {
			return err
		}
		b = append(b, record...)
		b = append(b, '\n')
	}

	_, err := r.client.AppendBlock(ctx, streaming.NopCloser(strings.NewReader(string(b))), nil)
	if err != nil {
		return fmt.Errorf("failed to append event to blob: %v", err)
	}

	return nil
}

// GetEvents retrieves events starting from the 'start' index, limited by the 'limit'
// But: This method is far from efficient, as it reads all events from the blob and
// then applies the pagination and returns a big list of events. It's only the first step.
func (r *AzureBlobRepository) GetAll(ctx context.Context, start, limit int) ([]core.Event, error) {
	continueOnError := true
	filter := bloom.NewWithEstimates(r.cfg.EstimatedEventCount, 0.1)

	resp, err := r.client.DownloadStream(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to download blob data: %v", err)
	}
	defer resp.Body.Close()

	// Stream the content and decode JSON events
	var events []core.Event
	decoder := json.NewDecoder(resp.Body)
	for decoder.More() {
		var event core.Event
		if err := decoder.Decode(&event); err != nil {
			if continueOnError {
				log.Println("Failed to decode event: ", err)
				continue
			}

			return nil, fmt.Errorf("failed to decode event: %v", err)
		}

		if !filter.Test([]byte(event.MessageId)) {
			events = append(events, event)
			filter.Add([]byte(event.MessageId))
		} else {
			log.Println("Ignoring potential duplicate event found: ", event.MessageId)
			if limit > 0 {
				limit++
			}
		}

		// Apply the limit if specified
		if limit > 0 && len(events) >= limit {
			break
		}
	}

	// Apply the start offset (skip first 'start' events)
	if start > 0 && start < len(events) {
		events = events[start:]
	}

	return events, nil
}

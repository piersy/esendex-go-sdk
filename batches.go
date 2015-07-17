package xesende

import (
	"encoding/xml"
	"errors"
	"time"
)

// BatchesResponse is a list of returned message batches along with the paging
// information.
type BatchesResponse struct {
	Paging
	Batches []BatchResponse
}

// BatchResponse is a single sent batch.
type BatchResponse struct {
	ID                 string
	URI                string
	CreatedAt          time.Time
	BatchSize          int
	PersistedBatchSize int
	Status             map[string]int
	AccountReference   string
	CreatedBy          string
	Name               string
}

// Batches returns a list of batches sent by the authenticated user.
func (c *Client) Batches() (*BatchesResponse, error) {
	req, err := c.newRequest("GET", "/v1.1/messagebatches", nil)
	if err != nil {
		return nil, err
	}

	var v messageBatchesResponse
	resp, err := c.do(req, &v)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("Expected 200")
	}

	response := &BatchesResponse{
		Paging: Paging{
			StartIndex: v.StartIndex,
			Count:      v.Count,
			TotalCount: v.TotalCount,
		},
		Batches: make([]BatchResponse, len(v.Batches)),
	}

	for i, batch := range v.Batches {
		status := map[string]int{}

		for _, s := range batch.Status.List {
			if s.Value > 0 {
				status[s.XMLName.Local] = s.Value
			}
		}

		response.Batches[i] = BatchResponse{
			ID:                 batch.ID,
			URI:                batch.URI,
			CreatedAt:          batch.CreatedAt,
			BatchSize:          batch.BatchSize,
			PersistedBatchSize: batch.PersistedBatchSize,
			Status:             status,
			AccountReference:   batch.AccountReference,
			CreatedBy:          batch.CreatedBy,
			Name:               batch.Name,
		}
	}

	return response, nil
}

type messageBatchesResponse struct {
	StartIndex int                    `xml:"startindex,attr"`
	Count      int                    `xml:"count,attr"`
	TotalCount int                    `xml:"totalcount,attr"`
	Batches    []messageBatchResponse `xml:"messagebatch"`
}

type messageBatchResponse struct {
	ID                 string                       `xml:"id,attr"`
	URI                string                       `xml:"uri,attr"`
	CreatedAt          time.Time                    `xml:"createdat"`
	BatchSize          int                          `xml:"batchsize"`
	PersistedBatchSize int                          `xml:"persistedbatchsize"`
	Status             messageBatchResponseStatuses `xml:"status"`
	AccountReference   string                       `xml:"accountreference"`
	CreatedBy          string                       `xml:"createdby"`
	Name               string                       `xml:"name"`
}

type messageBatchResponseStatuses struct {
	List []messageBatchResponseStatus `xml:",any"`
}

type messageBatchResponseStatus struct {
	XMLName xml.Name `xml:""`
	Value   int      `xml:",chardata"`
}

/*
Copyright 2016 The go-marathon Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package marathon

import (
	"fmt"
)

// Queue is the definition of marathon queue
type Queue struct {
	Items []Item `json:"queue"`
}

// Item is the definition of element in the queue
type Item struct {
	Count                  int                    `json:"count"`
	Delay                  Delay                  `json:"delay"`
	Application            *Application           `json:"app"`
	Pod                    *Pod                   `json:"pod"`
	Role                   string                 `json:"role"`
	Since                  string                 `json:"since"`
	ProcessedOffersSummary ProcessedOffersSummary `json:"processedOffersSummary"`
	LastUnusedOffers       []UnusedOffer          `json:"lastUnusedOffers,omitempty"`
}

// Delay cotains the application postpone information
type Delay struct {
	Overdue         bool `json:"overdue"`
	TimeLeftSeconds int  `json:"timeLeftSeconds"`
}

// ProcessedOffersSummary contains statistics for processed offers.
type ProcessedOffersSummary struct {
	ProcessedOffersCount       int32               `json:"processedOffersCount"`
	UnusedOffersCount          int32               `json:"unusedOffersCount"`
	LastUnusedOfferAt          *string             `json:"lastUnusedOfferAt,omitempty"`
	LastUsedOfferAt            *string             `json:"lastUsedOfferAt,omitempty"`
	RejectSummaryLastOffers    []DeclinedOfferStep `json:"rejectSummaryLastOffers,omitempty"`
	RejectSummaryLaunchAttempt []DeclinedOfferStep `json:"rejectSummaryLaunchAttempt,omitempty"`
}

// DeclinedOfferStep contains how often an offer was declined for a specific reason
type DeclinedOfferStep struct {
	Reason    string `json:"reason"`
	Declined  int32  `json:"declined"`
	Processed int32  `json:"processed"`
}

// UnusedOffer contains which offers weren't used and why
type UnusedOffer struct {
	Offer     Offer    `json:"offer"`
	Reason    []string `json:"reason"`
	Timestamp string   `json:"timestamp"`
}

// Queue retrieves content of the marathon launch queue
func (r *marathonClient) Queue() (*Queue, error) {
	var queue *Queue
	err := r.apiGet(marathonAPIQueue, nil, &queue)
	if err != nil {
		return nil, err
	}
	return queue, nil
}

// DeleteQueueDelay resets task launch delay of the specific application
//		appID:		the ID of the application
func (r *marathonClient) DeleteQueueDelay(appID string) error {
	path := fmt.Sprintf("%s/%s/delay", marathonAPIQueue, trimRootPath(appID))
	return r.apiDelete(path, nil, nil)
}

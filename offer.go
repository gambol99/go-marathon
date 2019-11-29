/*
Copyright 2019 The go-marathon Authors All rights reserved.

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

// based on https://github.com/mesosphere/marathon/blob/e7b1456ad0cfba23c9fdfa3c5d638a4b9aeb60d0/docs/docs/rest-api/public/api/v2/types/offer.raml

// Offer describes a Mesos offer to a framework
type Offer struct {
	ID         string           `json:"id"`
	Hostname   string           `json:"hostname"`
	AgentID    string           `json:"agentId"`
	Resources  []OfferResource  `json:"resources"`
	Attributes []AgentAttribute `json:"attributes"`
}

// OfferResource describes a resource that is part of an offer
type OfferResource struct {
	Name   string        `json:"name"`
	Role   string        `json:"role"`
	Scalar *float64      `json:"scalar,omitempty"`
	Ranges []NumberRange `json:"ranges,omitempty"`
	Set    []string      `json:"set,omitempty"`
}

// NumberRange is a range of numbers
type NumberRange struct {
	Begin int64 `json:"begin"`
	End   int64 `json:"end"`
}

// AgentAttribute describes an attribute of an agent node
type AgentAttribute struct {
	Name   string        `json:"name"`
	Text   *string       `json:"text,omitempty"`
	Scalar *float64      `json:"scalar,omitempty"`
	Ranges []NumberRange `json:"ranges,omitempty"`
	Set    []string      `json:"set,omitempty"`
}

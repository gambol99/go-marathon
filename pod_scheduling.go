/*
Copyright 2017 The go-marathon Authors All rights reserved.

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

// PodBackoff describes the backoff for re-run attempts of a pod
type PodBackoff struct {
	Backoff        *int     `json:"backoff,omitempty"`
	BackoffFactor  *float64 `json:"backoffFactor,omitempty"`
	MaxLaunchDelay *int     `json:"maxLaunchDelay,omitempty"`
}

// PodUpgrade describes the policy for upgrading a pod in-place
type PodUpgrade struct {
	MinimumHealthCapacity *float64 `json:"minimumHealthCapacity,omitempty"`
	MaximumOverCapacity   *float64 `json:"maximumOverCapacity,omitempty"`
}

// PodPlacement supports constraining which hosts a pod is placed on
type PodPlacement struct {
	Constraints           *[]Constraint `json:"constraints"`
	AcceptedResourceRoles []string      `json:"acceptedResourceRoles,omitempty"`
}

// PodSchedulingPolicy is the overarching pod scheduling policy
type PodSchedulingPolicy struct {
	Backoff             *PodBackoff          `json:"backoff,omitempty"`
	Upgrade             *PodUpgrade          `json:"upgrade,omitempty"`
	Placement           *PodPlacement        `json:"placement,omitempty"`
	UnreachableStrategy *UnreachableStrategy `json:"unreachableStrategy,omitempty"`
	KillSelection       string               `json:"killSelection,omitempty"`
}

// Constraint describes the constraint for pod placement
type Constraint struct {
	FieldName string `json:"fieldName"`
	Operator  string `json:"operator"`
	Value     string `json:"value,omitempty"`
}

// NewPodPlacement creates an empty PodPlacement
func NewPodPlacement() *PodPlacement {
	return &PodPlacement{
		Constraints:           &[]Constraint{},
		AcceptedResourceRoles: []string{},
	}
}

// AddConstraint adds a new constraint
//		constraints:	the constraint definition, one constraint per array element
func (p *PodPlacement) AddConstraint(constraint Constraint) *PodPlacement {
	c := *p.Constraints
	c = append(c, constraint)
	p.Constraints = &c

	return p
}

// NewPodSchedulingPolicy creates an empty PodSchedulingPolicy
func NewPodSchedulingPolicy() *PodSchedulingPolicy {
	return &PodSchedulingPolicy{
		Placement: NewPodPlacement(),
	}
}

// NewPodBackoff creates an empty PodBackoff
func NewPodBackoff() *PodBackoff {
	return &PodBackoff{}
}

// NewPodUpgrade creates a new PodUpgrade
func NewPodUpgrade() *PodUpgrade {
	return &PodUpgrade{}
}

// SetBackoff sets the base backoff interval for failed pod launches, in seconds
func (p *PodBackoff) SetBackoff(backoffSeconds int) *PodBackoff {
	p.Backoff = &backoffSeconds
	return p
}

// SetBackoffFactor sets the backoff interval growth factor for failed pod launches
func (p *PodBackoff) SetBackoffFactor(backoffFactor float64) *PodBackoff {
	p.BackoffFactor = &backoffFactor
	return p
}

// SetMaxLaunchDelay sets the maximum backoff interval for failed pod launches, in seconds
func (p *PodBackoff) SetMaxLaunchDelay(maxLaunchDelaySeconds int) *PodBackoff {
	p.MaxLaunchDelay = &maxLaunchDelaySeconds
	return p
}

// SetMinimumHealthCapacity sets the minimum amount of pod instances for healthy operation, expressed as a fraction of instance count
func (p *PodUpgrade) SetMinimumHealthCapacity(capacity float64) *PodUpgrade {
	p.MinimumHealthCapacity = &capacity
	return p
}

// SetMaximumOverCapacity sets the maximum amount of pod instances above the instance count, expressed as a fraction of instance count
func (p *PodUpgrade) SetMaximumOverCapacity(capacity float64) *PodUpgrade {
	p.MaximumOverCapacity = &capacity
	return p
}

// SetBackoff sets the pod's backoff settings
func (p *PodSchedulingPolicy) SetBackoff(backoff *PodBackoff) *PodSchedulingPolicy {
	p.Backoff = backoff
	return p
}

// SetUpgrade sets the pod's upgrade settings
func (p *PodSchedulingPolicy) SetUpgrade(upgrade *PodUpgrade) *PodSchedulingPolicy {
	p.Upgrade = upgrade
	return p
}

// SetPlacement sets the pod's placement settings
func (p *PodSchedulingPolicy) SetPlacement(placement *PodPlacement) *PodSchedulingPolicy {
	p.Placement = placement
	return p
}

// SetKillSelection sets the pod's kill selection criteria when terminating pod instances
func (p *PodSchedulingPolicy) SetKillSelection(killSelection string) *PodSchedulingPolicy {
	p.KillSelection = killSelection
	return p
}

// SetUnreachableStrategy sets the pod's unreachable strategy for lost instances
func (p *PodSchedulingPolicy) SetUnreachableStrategy(strategy EnabledUnreachableStrategy) *PodSchedulingPolicy {
	p.UnreachableStrategy = &UnreachableStrategy{
		EnabledUnreachableStrategy: strategy,
	}
	return p
}

// SetUnreachableStrategyDisabled disables the pod's unreachable strategy
func (p *PodSchedulingPolicy) SetUnreachableStrategyDisabled() *PodSchedulingPolicy {
	p.UnreachableStrategy = &UnreachableStrategy{
		AbsenceReason: UnreachableStrategyAbsenceReasonDisabled,
	}
	return p
}

/*
Copyright 2014 The go-marathon Authors All rights reserved.

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

// HealthCheck is the definition for an application health check
type HealthCheck struct {
	Command                *Command `json:"command,omitempty"`
	PortIndex              *int     `json:"portIndex,omitempty"`
	Port                   *int     `json:"port,omitempty"`
	Path                   *string  `json:"path,omitempty"`
	MaxConsecutiveFailures *int     `json:"maxConsecutiveFailures,omitempty"`
	Protocol               string   `json:"protocol,omitempty"`
	GracePeriodSeconds     int      `json:"gracePeriodSeconds,omitempty"`
	IntervalSeconds        int      `json:"intervalSeconds,omitempty"`
	TimeoutSeconds         int      `json:"timeoutSeconds,omitempty"`
	IgnoreHTTP1xx          *bool    `json:"ignoreHttp1xx,omitempty"`
}

// HTTPHealthCheck describes an HTTP based health check
type HTTPHealthCheck struct {
	Endpoint string `json:"endpoint,omitempty"`
	Path     string `json:"path,omitempty"`
	Scheme   string `json:"scheme,omitempty"`
}

// TCPHealthCheck describes a TCP based health check
type TCPHealthCheck struct {
	Endpoint string `json:"endpoint,omitempty"`
}

// CommandHealthCheck describes a shell-based health check
type CommandHealthCheck struct {
	Command PodCommand `json:"command,omitempty"`
}

// PodHealthCheck describes how to determine a pod's health
type PodHealthCheck struct {
	HTTP                   *HTTPHealthCheck    `json:"http,omitempty"`
	TCP                    *TCPHealthCheck     `json:"tcp,omitempty"`
	Exec                   *CommandHealthCheck `json:"exec,omitempty"`
	GracePeriodSeconds     *int                `json:"gracePeriodSeconds,omitempty"`
	IntervalSeconds        *int                `json:"intervalSeconds,omitempty"`
	MaxConsecutiveFailures *int                `json:"maxConsecutiveFailures,omitempty"`
	TimeoutSeconds         *int                `json:"timeoutSeconds,omitempty"`
	DelaySeconds           *int                `json:"delaySeconds,omitempty"`
}

// NewPodHealthCheck creates an empty PodHealthCheck
func NewPodHealthCheck() *PodHealthCheck {
	return &PodHealthCheck{}
}

// NewHTTPHealthCheck creates an empty HTTPHealthCheck
func NewHTTPHealthCheck() *HTTPHealthCheck {
	return &HTTPHealthCheck{}
}

// NewTCPHealthCheck creates an empty TCPHealthCheck
func NewTCPHealthCheck() *TCPHealthCheck {
	return &TCPHealthCheck{}
}

// NewCommandHealthCheck creates an empty CommandHealthCheck
func NewCommandHealthCheck() *CommandHealthCheck {
	return &CommandHealthCheck{}
}

// SetCommand sets the given command on the health check.
func (h *HealthCheck) SetCommand(c Command) *HealthCheck {
	h.Command = &c
	return h
}

// SetPortIndex sets the given port index on the health check.
func (h *HealthCheck) SetPortIndex(i int) *HealthCheck {
	h.PortIndex = &i
	return h
}

// SetPort sets the given port on the health check.
func (h *HealthCheck) SetPort(i int) *HealthCheck {
	h.Port = &i
	return h
}

// SetPath sets the given path on the health check.
func (h *HealthCheck) SetPath(p string) *HealthCheck {
	h.Path = &p
	return h
}

// SetMaxConsecutiveFailures sets the maximum consecutive failures on the health check.
func (h *HealthCheck) SetMaxConsecutiveFailures(i int) *HealthCheck {
	h.MaxConsecutiveFailures = &i
	return h
}

// SetIgnoreHTTP1xx sets ignore http 1xx on the health check.
func (h *HealthCheck) SetIgnoreHTTP1xx(ignore bool) *HealthCheck {
	h.IgnoreHTTP1xx = &ignore
	return h
}

// NewDefaultHealthCheck creates a default application health check
func NewDefaultHealthCheck() *HealthCheck {
	portIndex := 0
	path := ""
	maxConsecutiveFailures := 3

	return &HealthCheck{
		Protocol:               "HTTP",
		Path:                   &path,
		PortIndex:              &portIndex,
		MaxConsecutiveFailures: &maxConsecutiveFailures,
		GracePeriodSeconds:     30,
		IntervalSeconds:        10,
		TimeoutSeconds:         5,
	}
}

// HealthCheckResult is the health check result
type HealthCheckResult struct {
	Alive               bool   `json:"alive"`
	ConsecutiveFailures int    `json:"consecutiveFailures"`
	FirstSuccess        string `json:"firstSuccess"`
	LastFailure         string `json:"lastFailure"`
	LastFailureCause    string `json:"lastFailureCause"`
	LastSuccess         string `json:"lastSuccess"`
	TaskID              string `json:"taskId"`
}

// Command is the command health check type
type Command struct {
	Value string `json:"value"`
}

// SetHTTPHealthCheck configures the pod's health check for an HTTP endpoint.
// Note this will erase any configured TCP/Exec health checks.
func (p *PodHealthCheck) SetHTTPHealthCheck(h *HTTPHealthCheck) *PodHealthCheck {
	p.HTTP = h
	p.TCP = nil
	p.Exec = nil
	return p
}

// SetTCPHealthCheck configures the pod's health check for a TCP endpoint.
// Note this will erase any configured HTTP/Exec health checks.
func (p *PodHealthCheck) SetTCPHealthCheck(t *TCPHealthCheck) *PodHealthCheck {
	p.TCP = t
	p.HTTP = nil
	p.Exec = nil
	return p
}

// SetExecHealthCheck configures the pod's health check for a command.
// Note this will erase any configured HTTP/TCP health checks.
func (p *PodHealthCheck) SetExecHealthCheck(e *CommandHealthCheck) *PodHealthCheck {
	p.Exec = e
	p.HTTP = nil
	p.TCP = nil
	return p
}

// SetGracePeriod sets the health check initial grace period, in seconds
func (p *PodHealthCheck) SetGracePeriod(gracePeriodSeconds int) *PodHealthCheck {
	p.GracePeriodSeconds = &gracePeriodSeconds
	return p
}

// SetInterval sets the health check polling interval, in seconds
func (p *PodHealthCheck) SetInterval(intervalSeconds int) *PodHealthCheck {
	p.IntervalSeconds = &intervalSeconds
	return p
}

// SetMaxConsecutiveFailures sets the maximum consecutive failures on the health check
func (p *PodHealthCheck) SetMaxConsecutiveFailures(maxFailures int) *PodHealthCheck {
	p.MaxConsecutiveFailures = &maxFailures
	return p
}

// SetTimeout sets the length of time the health check will await a result, in seconds
func (p *PodHealthCheck) SetTimeout(timeoutSeconds int) *PodHealthCheck {
	p.TimeoutSeconds = &timeoutSeconds
	return p
}

// SetDelay sets the length of time a pod will delay running health checks on initial launch, in seconds
func (p *PodHealthCheck) SetDelay(delaySeconds int) *PodHealthCheck {
	p.DelaySeconds = &delaySeconds
	return p
}

// SetEndpoint sets the name of the pod health check endpoint
func (h *HTTPHealthCheck) SetEndpoint(endpoint string) *HTTPHealthCheck {
	h.Endpoint = endpoint
	return h
}

// SetPath sets the HTTP path of the pod health check endpoint
func (h *HTTPHealthCheck) SetPath(path string) *HTTPHealthCheck {
	h.Path = path
	return h
}

// SetScheme sets the HTTP scheme of the pod health check endpoint
func (h *HTTPHealthCheck) SetScheme(scheme string) *HTTPHealthCheck {
	h.Scheme = scheme
	return h
}

// SetEndpoint sets the name of the pod health check endpoint
func (t *TCPHealthCheck) SetEndpoint(endpoint string) *TCPHealthCheck {
	t.Endpoint = endpoint
	return t
}

// SetCommand sets a CommandHealthCheck's underlying PodCommand
func (c *CommandHealthCheck) SetCommand(p PodCommand) *CommandHealthCheck {
	c.Command = p
	return c
}

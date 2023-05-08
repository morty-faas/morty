package types

import "net/url"

const DefaultFunctionVersion = "v0.1.0"

type Function struct {
	// We don't want to serialize the name and the version as a Redis HSET as we use the ID as the key
	Name     string `json:"name" redis:"-"`
	Version  string `json:"version" redis:"-"`
	ImageURL string `json:"image" redis:"imageUrl"`
	// The identifier of the workload on the underlying orchestrator
	OrchestratorId string `json:"-" redis:"orchestratorId"`
}

// Id returns the computed identifier of the function in the following format:
// functionName:functionVersion.
func (f *Function) Id() string {
	return f.Name + ":" + f.Version
}

type FnInstance struct {
	Id       string    `json:"id"`
	Function *Function `json:"function"`
	Endpoint *url.URL  `json:"endpoint"`
}

type FunctionProcessMetadata struct {
	ExecutionTimeMs int `json:"execution_time_ms"`
	// Array of strings containing the logs of the function execution
	// it might be equal to nil if the runtime doesn't support logs
	Logs []string `json:"logs"`
}

type FnInvocationResponse struct {
	// Considered as string as we don't do anything to this field, we'll just return it
	Payload any `json:"payload"`
	// ProcessMetadata contains metadata about the function execution
	ProcessMetadata FunctionProcessMetadata `json:"process_metadata"`
}

package types

import "net/url"

type Function struct {
	Id string `json:"id" redis:"id"`
	// We don't want to serialize the name as a Redis HSET as we use the ID as the key
	Name     string `json:"name" redis:"-"`
	ImageURL string `json:"image" redis:"imageUrl"`
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

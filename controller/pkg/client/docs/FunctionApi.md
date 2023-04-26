# \FunctionApi

All URIs are relative to _http://localhost_

| Method                                              | HTTP request        | Description                           |
| --------------------------------------------------- | ------------------- | ------------------------------------- |
| [**CreateFunction**](FunctionApi.md#CreateFunction) | **Post** /functions | Create a new function                 |
| [**GetFunctions**](FunctionApi.md#GetFunctions)     | **Get** /functions  | Get a list of the available functions |

## CreateFunction

> Function CreateFunction(ctx).CreateFunctionRequest(createFunctionRequest).Execute()

Create a new function

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/morty-faas/morty/controller/pkg/client"
)

func main() {
    createFunctionRequest := *openapiclient.NewCreateFunctionRequest() // CreateFunctionRequest |

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.FunctionApi.CreateFunction(context.Background()).CreateFunctionRequest(createFunctionRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `FunctionApi.CreateFunction``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `CreateFunction`: Function
    fmt.Fprintf(os.Stdout, "Response from `FunctionApi.CreateFunction`: %v\n", resp)
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a apiCreateFunctionRequest struct via the builder pattern

| Name                      | Type                                                  | Description | Notes |
| ------------------------- | ----------------------------------------------------- | ----------- | ----- |
| **createFunctionRequest** | [**CreateFunctionRequest**](CreateFunctionRequest.md) |             |

### Return type

[**Function**](Function.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

## GetFunctions

> []Function GetFunctions(ctx).Execute()

Get a list of the available functions

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/morty-faas/morty/controller/pkg/client"
)

func main() {

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.FunctionApi.GetFunctions(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `FunctionApi.GetFunctions``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetFunctions`: []Function
    fmt.Fprintf(os.Stdout, "Response from `FunctionApi.GetFunctions`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetFunctionsRequest struct via the builder pattern

### Return type

[**[]Function**](Function.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

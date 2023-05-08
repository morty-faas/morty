# \FunctionsApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**V1FunctionsBuildPost**](FunctionsApi.md#V1FunctionsBuildPost) | **Post** /v1/functions/build | Build a function and push the image into the registry.
[**V1FunctionsIdVersionGet**](FunctionsApi.md#V1FunctionsIdVersionGet) | **Get** /v1/functions/{id}/{version} | Get a download link for the image of the given function



## V1FunctionsBuildPost

> string V1FunctionsBuildPost(ctx).Runtime(runtime).Name(name).Version(version).Archive(archive).Execute()

Build a function and push the image into the registry.

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/pkg/client/registry/morty-faas/morty"
)

func main() {
    runtime := "runtime_example" // string | The name of the runtime to use. (optional)
    name := "name_example" // string | The name of the function. (optional)
    version := "version_example" // string | The version of the function. This should be a semantic version. (optional)
    archive := os.NewFile(1234, "some_file") // *os.File |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.FunctionsApi.V1FunctionsBuildPost(context.Background()).Runtime(runtime).Name(name).Version(version).Archive(archive).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `FunctionsApi.V1FunctionsBuildPost``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `V1FunctionsBuildPost`: string
    fmt.Fprintf(os.Stdout, "Response from `FunctionsApi.V1FunctionsBuildPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiV1FunctionsBuildPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **runtime** | **string** | The name of the runtime to use. | 
 **name** | **string** | The name of the function. | 
 **version** | **string** | The version of the function. This should be a semantic version. | 
 **archive** | ***os.File** |  | 

### Return type

**string**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: multipart/form-data
- **Accept**: application/text

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## V1FunctionsIdVersionGet

> V1FunctionsIdVersionGet(ctx, id, version).Execute()

Get a download link for the image of the given function

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/pkg/client/registry/morty-faas/morty"
)

func main() {
    id := "id_example" // string | The identifier of the function to download.
    version := "version_example" // string | The version of the function to download.

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.FunctionsApi.V1FunctionsIdVersionGet(context.Background(), id, version).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `FunctionsApi.V1FunctionsIdVersionGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | The identifier of the function to download. | 
**version** | **string** | The version of the function to download. | 

### Other Parameters

Other parameters are passed through a pointer to a apiV1FunctionsIdVersionGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


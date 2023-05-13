# \ConfigurationApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetServerMetadata**](ConfigurationApi.md#GetServerMetadata) | **Get** /.well-known/morty.json | Get informations about the server such as version, build commit etc.



## GetServerMetadata

> ServerMetadata GetServerMetadata(ctx).Execute()

Get informations about the server such as version, build commit etc.



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/pkg/client/controller/morty-faas/morty"
)

func main() {

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.ConfigurationApi.GetServerMetadata(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ConfigurationApi.GetServerMetadata``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetServerMetadata`: ServerMetadata
    fmt.Fprintf(os.Stdout, "Response from `ConfigurationApi.GetServerMetadata`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetServerMetadataRequest struct via the builder pattern


### Return type

[**ServerMetadata**](ServerMetadata.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


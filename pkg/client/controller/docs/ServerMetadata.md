# ServerMetadata

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Version** | Pointer to **string** | The version of the Morty server | [optional] 
**GitCommit** | Pointer to **string** | The Git commit the server was built on. | [optional] 

## Methods

### NewServerMetadata

`func NewServerMetadata() *ServerMetadata`

NewServerMetadata instantiates a new ServerMetadata object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewServerMetadataWithDefaults

`func NewServerMetadataWithDefaults() *ServerMetadata`

NewServerMetadataWithDefaults instantiates a new ServerMetadata object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetVersion

`func (o *ServerMetadata) GetVersion() string`

GetVersion returns the Version field if non-nil, zero value otherwise.

### GetVersionOk

`func (o *ServerMetadata) GetVersionOk() (*string, bool)`

GetVersionOk returns a tuple with the Version field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersion

`func (o *ServerMetadata) SetVersion(v string)`

SetVersion sets Version field to given value.

### HasVersion

`func (o *ServerMetadata) HasVersion() bool`

HasVersion returns a boolean if a field has been set.

### GetGitCommit

`func (o *ServerMetadata) GetGitCommit() string`

GetGitCommit returns the GitCommit field if non-nil, zero value otherwise.

### GetGitCommitOk

`func (o *ServerMetadata) GetGitCommitOk() (*string, bool)`

GetGitCommitOk returns a tuple with the GitCommit field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGitCommit

`func (o *ServerMetadata) SetGitCommit(v string)`

SetGitCommit sets GitCommit field to given value.

### HasGitCommit

`func (o *ServerMetadata) HasGitCommit() bool`

HasGitCommit returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)



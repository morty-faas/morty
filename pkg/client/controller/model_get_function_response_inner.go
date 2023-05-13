/*
Morty APIs

This document contains the specification of the public-facing Morty APIs. For function invocation, please see the project README here: https://github.com/morty-faas/morty/controller#readme 

API version: 1.1.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package client

import (
	"encoding/json"
)

// checks if the GetFunctionResponseInner type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &GetFunctionResponseInner{}

// GetFunctionResponseInner struct for GetFunctionResponseInner
type GetFunctionResponseInner struct {
	Name *string `json:"name,omitempty"`
	Versions []string `json:"versions,omitempty"`
}

// NewGetFunctionResponseInner instantiates a new GetFunctionResponseInner object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGetFunctionResponseInner() *GetFunctionResponseInner {
	this := GetFunctionResponseInner{}
	return &this
}

// NewGetFunctionResponseInnerWithDefaults instantiates a new GetFunctionResponseInner object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGetFunctionResponseInnerWithDefaults() *GetFunctionResponseInner {
	this := GetFunctionResponseInner{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *GetFunctionResponseInner) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GetFunctionResponseInner) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *GetFunctionResponseInner) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *GetFunctionResponseInner) SetName(v string) {
	o.Name = &v
}

// GetVersions returns the Versions field value if set, zero value otherwise.
func (o *GetFunctionResponseInner) GetVersions() []string {
	if o == nil || IsNil(o.Versions) {
		var ret []string
		return ret
	}
	return o.Versions
}

// GetVersionsOk returns a tuple with the Versions field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GetFunctionResponseInner) GetVersionsOk() ([]string, bool) {
	if o == nil || IsNil(o.Versions) {
		return nil, false
	}
	return o.Versions, true
}

// HasVersions returns a boolean if a field has been set.
func (o *GetFunctionResponseInner) HasVersions() bool {
	if o != nil && !IsNil(o.Versions) {
		return true
	}

	return false
}

// SetVersions gets a reference to the given []string and assigns it to the Versions field.
func (o *GetFunctionResponseInner) SetVersions(v []string) {
	o.Versions = v
}

func (o GetFunctionResponseInner) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o GetFunctionResponseInner) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.Versions) {
		toSerialize["versions"] = o.Versions
	}
	return toSerialize, nil
}

type NullableGetFunctionResponseInner struct {
	value *GetFunctionResponseInner
	isSet bool
}

func (v NullableGetFunctionResponseInner) Get() *GetFunctionResponseInner {
	return v.value
}

func (v *NullableGetFunctionResponseInner) Set(val *GetFunctionResponseInner) {
	v.value = val
	v.isSet = true
}

func (v NullableGetFunctionResponseInner) IsSet() bool {
	return v.isSet
}

func (v *NullableGetFunctionResponseInner) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableGetFunctionResponseInner(val *GetFunctionResponseInner) *NullableGetFunctionResponseInner {
	return &NullableGetFunctionResponseInner{value: val, isSet: true}
}

func (v NullableGetFunctionResponseInner) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableGetFunctionResponseInner) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


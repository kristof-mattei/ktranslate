/*
 * VeloCloud Orchestrator API
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 3.3.2
 * Contact: support@velocloud.net
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package models
import (
	"time"
)

type EnterpriseProxyProperty struct {
	Id int32 `json:"id,omitempty"`
	EnterpriseProxyId int32 `json:"enterpriseProxyId,omitempty"`
	Created time.Time `json:"created,omitempty"`
	Name string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
	IsPassword bool `json:"isPassword,omitempty"`
	DataType string `json:"dataType,omitempty"`
	Description string `json:"description,omitempty"`
	Modified time.Time `json:"modified,omitempty"`
}

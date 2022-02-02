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

type EnterpriseProxy struct {
	Id int32 `json:"id,omitempty"`
	Created time.Time `json:"created,omitempty"`
	NetworkId int32 `json:"networkId,omitempty"`
	ProxyType string `json:"proxyType,omitempty"`
	OperateGateways bool `json:"operateGateways,omitempty"`
	LogicalId string `json:"logicalId,omitempty"`
	Name string `json:"name,omitempty"`
	Domain string `json:"domain,omitempty"`
	Prefix string `json:"prefix,omitempty"`
	Description string `json:"description,omitempty"`
	ContactName string `json:"contactName,omitempty"`
	ContactPhone string `json:"contactPhone,omitempty"`
	ContactMobile string `json:"contactMobile,omitempty"`
	ContactEmail string `json:"contactEmail,omitempty"`
	StreetAddress string `json:"streetAddress,omitempty"`
	StreetAddress2 string `json:"streetAddress2,omitempty"`
	City string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
	PostalCode string `json:"postalCode,omitempty"`
	Country string `json:"country,omitempty"`
	Lat float64 `json:"lat,omitempty"`
	Lon float64 `json:"lon,omitempty"`
	Modified time.Time `json:"modified,omitempty"`
}

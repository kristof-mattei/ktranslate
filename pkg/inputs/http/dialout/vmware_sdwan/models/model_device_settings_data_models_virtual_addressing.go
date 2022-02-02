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

type DeviceSettingsDataModelsVirtualAddressing struct {
	Type       string `json:"type,omitempty"`
	CidrPrefix int32  `json:"cidrPrefix,omitempty"`
	CidrIp     string `json:"cidrIp,omitempty"`
	Netmask    string `json:"netmask,omitempty"`
	Gateway    string `json:"gateway,omitempty"`
}

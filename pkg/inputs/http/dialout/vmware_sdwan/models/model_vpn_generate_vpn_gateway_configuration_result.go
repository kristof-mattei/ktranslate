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

type VpnGenerateVpnGatewayConfigurationResult struct {
	Id int32 `json:"id,omitempty"`
	Object string `json:"object,omitempty"`
	Type string `json:"type,omitempty"`
	Data DataCenterData `json:"data,omitempty"`
}

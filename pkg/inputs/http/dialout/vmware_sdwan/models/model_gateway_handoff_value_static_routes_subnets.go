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

type GatewayHandoffValueStaticRoutesSubnets struct {
	Name string `json:"name,omitempty"`
	CidrIp string `json:"cidrIp,omitempty"`
	CidrPrefix int32 `json:"cidrPrefix,omitempty"`
	Encrypt bool `json:"encrypt,omitempty"`
	HandOffType string `json:"handOffType,omitempty"`
	RouteCost int32 `json:"routeCost,omitempty"`
}

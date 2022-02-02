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

type NetworkGetNetworkGateways struct {
	NetworkId int32 `json:"networkId,omitempty"`
	GatewayIds []string `json:"gatewayIds,omitempty"`
	With []string `json:"with,omitempty"`
}

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

type DeviceSettingsHandoffGateways struct {
	Override    bool                           `json:"override,omitempty"`
	AutoSelect  bool                           `json:"autoSelect"`
	GatewayList []DeviceSettingsHandoffGateway `json:"gatewayList"`
	Gateways    []LogicalidReference           `json:"gateways,omitempty"`
}

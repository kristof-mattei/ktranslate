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

type DeviceSettingsMulticastPimOnWanOverlay struct {
	Enabled         bool   `json:"enabled"`
	Type            string `json:"type"`
	SourceCidrIp    string `json:"sourceCidrIp"`
	SourceInterface string `json:"sourceInterface"`
}

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

type DeviceSettingsOspfBgpRedistribution struct {
	Enabled bool `json:"enabled,omitempty"`
	Metric int32 `json:"metric,omitempty"`
	MetricType string `json:"metricType,omitempty"`
}

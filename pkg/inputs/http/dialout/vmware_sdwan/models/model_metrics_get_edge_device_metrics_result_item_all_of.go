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

type MetricsGetEdgeDeviceMetricsResultItemAllOf struct {
	EdgeInfo MetricsGetEdgeDeviceMetricsDeviceEdgeInfo `json:"edgeInfo"`
	Info ClientDevice `json:"info"`
	Name string `json:"name"`
	SourceMac string `json:"sourceMac"`
}

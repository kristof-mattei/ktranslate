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

type MonitoringGetAggregatesResult struct {
	EdgeCount int32 `json:"edgeCount,omitempty"`
	Edges map[string]int32 `json:"edges,omitempty"`
	Enterprises []EnterpriseWithProxyAttributes `json:"enterprises,omitempty"`
}

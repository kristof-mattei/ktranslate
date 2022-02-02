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

type EnterpriseUpdateEnterpriseServiceUpdate struct {
	OperatorAlertsEnabled int32 `json:"operatorAlertsEnabled,omitempty"`
	AlertsEnabled int32 `json:"alertsEnabled,omitempty"`
	Data map[string]interface{} `json:"data,omitempty"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

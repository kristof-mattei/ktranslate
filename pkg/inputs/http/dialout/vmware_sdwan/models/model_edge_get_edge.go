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

type EdgeGetEdge struct {
	Id int32 `json:"id,omitempty"`
	EnterpriseId int32 `json:"enterpriseId,omitempty"`
	LogicalId string `json:"logicalId,omitempty"`
	ActivationKey string `json:"activationKey,omitempty"`
	Name string `json:"name,omitempty"`
	With []string `json:"with,omitempty"`
}

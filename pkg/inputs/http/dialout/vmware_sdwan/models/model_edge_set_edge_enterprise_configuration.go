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

type EdgeSetEdgeEnterpriseConfiguration struct {
	// alias for edgeId
	Id int32 `json:"id,omitempty"`
	EdgeId int32 `json:"edgeId,omitempty"`
	EnterpriseId int32 `json:"enterpriseId,omitempty"`
	ConfigurationId int32 `json:"configurationId"`
	SkipEdgeRoutingUpdates bool `json:"skipEdgeRoutingUpdates,omitempty"`
}

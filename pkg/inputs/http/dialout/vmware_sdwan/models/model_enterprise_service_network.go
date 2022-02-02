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

type EnterpriseServiceNetwork struct {
	Id int32 `json:"id,omitempty"`
	EnterpriseObjectId int32 `json:"enterpriseObjectId,omitempty"`
	ConfigurationId int32 `json:"configurationId,omitempty"`
	ModuleId int32 `json:"moduleId,omitempty"`
	Ref string `json:"ref,omitempty"`
	Data EnterpriseServiceNetworkData `json:"data,omitempty"`
	Version string `json:"version,omitempty"`
	Object string `json:"object,omitempty"`
	Name string `json:"name,omitempty"`
	Type string `json:"type"`
	LogicalId string `json:"logicalId,omitempty"`
}

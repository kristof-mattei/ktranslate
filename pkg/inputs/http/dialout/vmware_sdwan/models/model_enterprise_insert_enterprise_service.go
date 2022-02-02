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

type EnterpriseInsertEnterpriseService struct {
	EnterpriseId int32 `json:"enterpriseId,omitempty"`
	Type string `json:"type"`
	Name string `json:"name"`
	Data map[string]interface{} `json:"data"`
	EdgeId int32 `json:"edgeId,omitempty"`
	ParentGroupId int32 `json:"parentGroupId,omitempty"`
}

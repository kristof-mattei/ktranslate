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

type AggregateEnterpriseEventAllOf struct {
	EnterpriseId int32 `json:"enterpriseId"`
	EnterpriseName string `json:"enterpriseName"`
	EnterpriseUsername string `json:"enterpriseUsername"`
	EdgeName string `json:"edgeName"`
	Detail string `json:"detail"`
}

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

type EnterpriseProxyUpdateEnterpriseProxyUser struct {
	Update EnterpriseUserWithRoleInfo `json:"_update"`
	Id int32 `json:"id,omitempty"`
	EnterpriseProxyId int32 `json:"enterpriseProxyId,omitempty"`
	Username string `json:"username,omitempty"`
}

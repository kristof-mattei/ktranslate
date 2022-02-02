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

type ConfigurationCloneEnterpriseTemplate struct {
	// Required if called from the operator or MSP context, identifies the target enterprise of the API call.
	EnterpriseId int32 `json:"enterpriseId,omitempty"`
	// If both network and segment based functionality is granted to the enterprise, chose which template type to clone. If not specified the type of the operator profile assigned to the enterprise will be used.
	ConfigurationType string `json:"configurationType,omitempty"`
	Name string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

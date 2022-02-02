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

type SetClientDeviceHostNameResult struct {
	// An error message explaining why the method failed
	Error string `json:"error,omitempty"`
	// The number of rows modified
	Rows int32 `json:"rows"`
}

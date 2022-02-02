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

type FirewallDataServices struct {
	LoggingEnabled bool `json:"loggingEnabled"`
	Ssh FirewallDataServicesSsh `json:"ssh,omitempty"`
	LocalUi FirewallDataServicesLocalUi `json:"localUi,omitempty"`
	Snmp FirewallDataServicesSnmp `json:"snmp,omitempty"`
	Icmp FirewallDataServicesIcmp `json:"icmp,omitempty"`
}

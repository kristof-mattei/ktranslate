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

type EnterpriseGetEnterpriseEdgesResultItem struct {
	ActivationKey string `json:"activationKey,omitempty"`
	ActivationKeyExpires string `json:"activationKeyExpires,omitempty"`
	ActivationState string `json:"activationState,omitempty"`
	ActivationTime string `json:"activationTime,omitempty"`
	AlertsEnabled Tinyint `json:"alertsEnabled,omitempty"`
	BuildNumber string `json:"buildNumber,omitempty"`
	Created string `json:"created,omitempty"`
	Description string `json:"description,omitempty"`
	DeviceFamily string `json:"deviceFamily,omitempty"`
	DeviceId string `json:"deviceId,omitempty"`
	DnsName string `json:"dnsName,omitempty"`
	EdgeState string `json:"edgeState,omitempty"`
	EdgeStateTime string `json:"edgeStateTime,omitempty"`
	EndpointPkiMode string `json:"endpointPkiMode,omitempty"`
	EnterpriseId int32 `json:"enterpriseId,omitempty"`
	HaLastContact string `json:"haLastContact,omitempty"`
	HaPreviousState string `json:"haPreviousState,omitempty"`
	HaSerialNumber string `json:"haSerialNumber,omitempty"`
	HaState string `json:"haState,omitempty"`
	Id int32 `json:"id,omitempty"`
	IsLive int32 `json:"isLive,omitempty"`
	LastContact string `json:"lastContact,omitempty"`
	LogicalId string `json:"logicalId,omitempty"`
	ModelNumber string `json:"modelNumber,omitempty"`
	Modified string `json:"modified,omitempty"`
	Name string `json:"name,omitempty"`
	OperatorAlertsEnabled Tinyint `json:"operatorAlertsEnabled,omitempty"`
	SelfMacAddress string `json:"selfMacAddress,omitempty"`
	SerialNumber string `json:"serialNumber,omitempty"`
	ServiceState string `json:"serviceState,omitempty"`
	ServiceUpSince string `json:"serviceUpSince,omitempty"`
	SiteId int32 `json:"siteId,omitempty"`
	SoftwareUpdated string `json:"softwareUpdated,omitempty"`
	SoftwareVersion string `json:"softwareVersion,omitempty"`
	SystemUpSince string `json:"systemUpSince,omitempty"`
	Certificates []EdgeCertificate `json:"certificates,omitempty"`
	Configuration EnterpriseGetEnterpriseEdgesResultItemAllOfConfiguration `json:"configuration,omitempty"`
	Ha EnterpriseGetEnterpriseEdgesResultItemAllOfHa `json:"ha,omitempty"`
	Licenses []EdgeLicense `json:"licenses,omitempty"`
	Links []Link `json:"links,omitempty"`
	RecentLinks []Link `json:"recentLinks,omitempty"`
	Site Site `json:"site,omitempty"`
	IsHub bool `json:"isHub,omitempty"`
	IsSoftwareVersionSupportedByVco bool `json:"isSoftwareVersionSupportedByVco,omitempty"`
}

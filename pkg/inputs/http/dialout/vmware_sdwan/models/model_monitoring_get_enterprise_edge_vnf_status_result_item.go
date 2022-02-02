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
import (
	"time"
)

type MonitoringGetEnterpriseEdgeVnfStatusResultItem struct {
	Id int32 `json:"id"`
	Created time.Time `json:"created"`
	OperatorId int32 `json:"operatorId"`
	NetworkId int32 `json:"networkId"`
	EnterpriseId int32 `json:"enterpriseId"`
	EdgeId int32 `json:"edgeId"`
	GatewayId int32 `json:"gatewayId"`
	ParentGroupId int32 `json:"parentGroupId"`
	Description string `json:"description"`
	Object string `json:"object"`
	Name string `json:"name"`
	Type string `json:"type"`
	LogicalId string `json:"logicalId"`
	AlertsEnabled Tinyint `json:"alertsEnabled"`
	OperatorAlertsEnabled Tinyint `json:"operatorAlertsEnabled"`
	Status string `json:"status"`
	StatusModified time.Time `json:"statusModified"`
	PreviousData map[string]interface{} `json:"previousData"`
	PreviousCreated time.Time `json:"previousCreated"`
	DraftData string `json:"draftData"`
	DraftCreated time.Time `json:"draftCreated"`
	DraftComment string `json:"draftComment"`
	Data MonitoringGetEnterpriseEdgeVnfStatusResultItemAllOfData `json:"data"`
	LastContact time.Time `json:"lastContact"`
	Version string `json:"version"`
	Modified time.Time `json:"modified"`
	EdgeCount int32 `json:"edgeCount,omitempty"`
	EdgeUsage []MonitoringGetEnterpriseEdgeVnfStatusResultItemAllOfEdgeUsage `json:"edgeUsage,omitempty"`
}

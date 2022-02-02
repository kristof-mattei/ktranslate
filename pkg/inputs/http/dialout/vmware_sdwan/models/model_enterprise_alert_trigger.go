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

type EnterpriseAlertTrigger struct {
	Id int32 `json:"id,omitempty"`
	Created time.Time `json:"created,omitempty"`
	TriggerTime time.Time `json:"triggerTime,omitempty"`
	EnterpriseAlertConfigurationId int32 `json:"enterpriseAlertConfigurationId,omitempty"`
	EnterpriseId int32 `json:"enterpriseId,omitempty"`
	EdgeId int32 `json:"edgeId,omitempty"`
	EdgeName string `json:"edgeName,omitempty"`
	LinkId int32 `json:"linkId,omitempty"`
	LinkName string `json:"linkName,omitempty"`
	EnterpriseObjectId int32 `json:"enterpriseObjectId,omitempty"`
	EnterpriseObjectName string `json:"enterpriseObjectName,omitempty"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
	State string `json:"state,omitempty"`
	StateSetTime time.Time `json:"stateSetTime,omitempty"`
	LastContact time.Time `json:"lastContact,omitempty"`
	FirstNotificationSeconds int32 `json:"firstNotificationSeconds,omitempty"`
	MaxNotifications int32 `json:"maxNotifications,omitempty"`
	NotificationIntervalSeconds int32 `json:"notificationIntervalSeconds,omitempty"`
	ResetIntervalSeconds int32 `json:"resetIntervalSeconds,omitempty"`
	Comment string `json:"comment,omitempty"`
	NextNotificationTime time.Time `json:"nextNotificationTime,omitempty"`
	RemainingNotifications int32 `json:"remainingNotifications,omitempty"`
	Timezone string `json:"timezone,omitempty"`
	Locale string `json:"locale,omitempty"`
	Modified time.Time `json:"modified,omitempty"`
}

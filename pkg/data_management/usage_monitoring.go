package data_management

import (
	"fmt"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// RetrieveUsageStats fetches usage statistics for a given timeframe
func RetrieveUsageStats(startDate, endDate time.Time) (common.UsageStats, error) {
	stats, err := common.FetchUsageStats(startDate, endDate)
	if err != nil {
		return common.UsageStats{}, fmt.Errorf("failed to retrieve usage stats: %v", err)
	}
	return stats, nil
}

// MonitorActiveSessionDuration tracks the duration of active sessions and logs the results
func MonitorActiveSessionDuration(sessionID string) (time.Duration, error) {
	session, err := common.FetchSession(sessionID)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch session: %v", err)
	}
	duration := time.Since(session.StartTime)
	err = common.LogSessionDuration(sessionID, duration)
	if err != nil {
		return 0, fmt.Errorf("failed to log session duration: %v", err)
	}
	return duration, nil
}

// TrackKeyPerformanceIndicators logs and monitors key performance indicators (KPIs)
func TrackKeyPerformanceIndicators(kpi common.KPI) error {
	if err := common.SaveKPI(kpi); err != nil {
		return fmt.Errorf("failed to save KPI: %v", err)
	}
	return nil
}

// RunCustomAnalyticsQuery executes a custom analytics query and returns the result
func RunCustomAnalyticsQuery(query string) (common.AnalyticsResult, error) {
	result, err := common.ExecuteAnalyticsQuery(query)
	if err != nil {
		return common.AnalyticsResult{}, fmt.Errorf("failed to execute custom analytics query: %v", err)
	}
	return result, nil
}

// EnableDataInsightAlert enables an alert for specific insight criteria
func EnableDataInsightAlert(alertCriteria common.AlertCriteria) error {
	if err := common.AddInsightAlert(alertCriteria); err != nil {
		return fmt.Errorf("failed to enable data insight alert: %v", err)
	}
	return nil
}

// DisableDataInsightAlert disables a data insight alert based on criteria
func DisableDataInsightAlert(alertCriteria common.AlertCriteria) error {
	if err := common.RemoveInsightAlert(alertCriteria); err != nil {
		return fmt.Errorf("failed to disable data insight alert: %v", err)
	}
	return nil
}

// RetrieveAnomalyDetectionLogs fetches logs of detected anomalies within a specified timeframe
func RetrieveAnomalyDetectionLogs(startDate, endDate time.Time) ([]common.AnomalyLog, error) {
	logs, err := common.FetchAnomalyLogs(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve anomaly detection logs: %v", err)
	}
	return logs, nil
}

// GenerateInsightMetrics generates metrics based on specific insights criteria
func GenerateInsightMetrics(criteria common.InsightCriteria) (common.InsightMetrics, error) {
	metrics, err := common.CalculateInsightMetrics(criteria)
	if err != nil {
		return common.InsightMetrics{}, fmt.Errorf("failed to generate insight metrics: %v", err)
	}
	return metrics, nil
}

// SetCustomInsightCriteria sets criteria for custom insights, allowing for dynamic data analysis
func SetCustomInsightCriteria(criteria common.InsightCriteria) error {
	if err := common.SaveInsightCriteria(criteria); err != nil {
		return fmt.Errorf("failed to set custom insight criteria: %v", err)
	}
	return nil
}

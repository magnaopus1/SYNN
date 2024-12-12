package testnet

import (
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

// LOG_CONTRACT_BEHAVIOR logs behavior for a specific smart contract.
func LOG_CONTRACT_BEHAVIOR(contractID string, behaviorDetails map[string]interface{}) error {
	log.Printf("Logging contract behavior for contract ID: %s\n", contractID)
	err := ledger.RecordContractBehavior(contractID, behaviorDetails)
	if err != nil {
		log.Printf("Failed to log contract behavior: %v\n", err)
		return err
	}
	return nil
}

// LOG_NODE_BEHAVIOR logs specific behavior details for a network node.
func LOG_NODE_BEHAVIOR(nodeID string, behaviorDetails map[string]interface{}) error {
	log.Printf("Logging node behavior for node ID: %s\n", nodeID)
	err := ledger.RecordNodeBehavior(nodeID, behaviorDetails)
	if err != nil {
		log.Printf("Failed to log node behavior: %v\n", err)
		return err
	}
	return nil
}

// ANALYZE_PEAK_BEHAVIOR analyzes peak behavior metrics across nodes and contracts.
func ANALYZE_PEAK_BEHAVIOR() (map[string]interface{}, error) {
	log.Println("Analyzing peak behavior metrics.")
	peakMetrics, err := ledger.CalculatePeakBehavior()
	if err != nil {
		log.Printf("Failed to analyze peak behavior: %v\n", err)
		return nil, err
	}
	return peakMetrics, nil
}

// COMPARE_BEHAVIOR_RESULTS compares two sets of behavior results to find discrepancies.
func COMPARE_BEHAVIOR_RESULTS(behaviorID1, behaviorID2 string) (map[string]interface{}, error) {
	log.Printf("Comparing behavior results between ID: %s and ID: %s\n", behaviorID1, behaviorID2)
	comparisonResults, err := ledger.CompareBehaviorLogs(behaviorID1, behaviorID2)
	if err != nil {
		log.Printf("Failed to compare behavior results: %v\n", err)
		return nil, err
	}
	return comparisonResults, nil
}

// EVALUATE_TEST_METRICS evaluates test metrics against predefined goals.
func EVALUATE_TEST_METRICS(testMetrics map[string]interface{}) (bool, error) {
	log.Println("Evaluating test metrics.")
	isSuccessful, err := ledger.CheckMetricsAgainstGoals(testMetrics)
	if err != nil {
		log.Printf("Failed to evaluate test metrics: %v\n", err)
		return false, err
	}
	return isSuccessful, nil
}

// SAVE_BEHAVIOR_RESULTS saves the behavior results after analysis.
func SAVE_BEHAVIOR_RESULTS(behaviorID string, results map[string]interface{}) error {
	log.Printf("Saving behavior results for behavior ID: %s\n", behaviorID)
	err := ledger.StoreBehaviorResults(behaviorID, results)
	if err != nil {
		log.Printf("Failed to save behavior results: %v\n", err)
		return err
	}
	return nil
}

// EXPORT_BEHAVIOR_LOG exports a behavior log for external review.
func EXPORT_BEHAVIOR_LOG(behaviorID string) (string, error) {
	log.Printf("Exporting behavior log for ID: %s\n", behaviorID)
	exportedLog, err := ledger.ExportBehaviorLog(behaviorID)
	if err != nil {
		log.Printf("Failed to export behavior log: %v\n", err)
		return "", err
	}
	return exportedLog, nil
}

// IMPORT_BEHAVIOR_LOG imports a behavior log for internal analysis.
func IMPORT_BEHAVIOR_LOG(logData string) error {
	log.Println("Importing behavior log for analysis.")
	err := ledger.ImportBehaviorLog(logData)
	if err != nil {
		log.Printf("Failed to import behavior log: %v\n", err)
		return err
	}
	return nil
}

// RECORD_BEHAVIOR_FEEDBACK records feedback from behavior tests.
func RECORD_BEHAVIOR_FEEDBACK(behaviorID string, feedback string) error {
	log.Printf("Recording feedback for behavior ID: %s\n", behaviorID)
	err := ledger.StoreBehaviorFeedback(behaviorID, feedback)
	if err != nil {
		log.Printf("Failed to record behavior feedback: %v\n", err)
		return err
	}
	return nil
}

// SET_TEST_GOALS sets the testing goals and objectives.
func SET_TEST_GOALS(testID string, goals map[string]interface{}) error {
	log.Printf("Setting test goals for test ID: %s\n", testID)
	err := ledger.DefineTestGoals(testID, goals)
	if err != nil {
		log.Printf("Failed to set test goals: %v\n", err)
		return err
	}
	return nil
}

// DEFINE_TEST_CASES defines the test cases for behavior testing.
func DEFINE_TEST_CASES(testID string, cases []map[string]interface{}) error {
	log.Printf("Defining test cases for test ID: %s\n", testID)
	err := ledger.AddTestCases(testID, cases)
	if err != nil {
		log.Printf("Failed to define test cases: %v\n", err)
		return err
	}
	return nil
}

// RUN_TEST_CASE runs a specific test case.
func RUN_TEST_CASE(testID, caseID string) error {
	log.Printf("Running test case ID: %s for test ID: %s\n", caseID, testID)
	err := ledger.ExecuteTestCase(testID, caseID)
	if err != nil {
		log.Printf("Failed to run test case: %v\n", err)
		return err
	}
	return nil
}

// MONITOR_RESOURCE_USAGE monitors resource usage during tests.
func MONITOR_RESOURCE_USAGE() (map[string]interface{}, error) {
	log.Println("Monitoring resource usage.")
	resourceUsage, err := ledger.TrackResourceUsage()
	if err != nil {
		log.Printf("Failed to monitor resource usage: %v\n", err)
		return nil, err
	}
	return resourceUsage, nil
}

// LOG_MEMORY_USAGE logs the memory usage during test execution.
func LOG_MEMORY_USAGE(testID string, memoryUsage int) error {
	log.Printf("Logging memory usage for test ID: %s - Usage: %d MB\n", testID, memoryUsage)
	err := ledger.RecordMemoryUsage(testID, memoryUsage)
	if err != nil {
		log.Printf("Failed to log memory usage: %v\n", err)
		return err
	}
	return nil
}

// LOG_CPU_USAGE logs the CPU usage during test execution.
func LOG_CPU_USAGE(testID string, cpuUsage float64) error {
	log.Printf("Logging CPU usage for test ID: %s - Usage: %.2f%%\n", testID, cpuUsage)
	err := ledger.RecordCPUUsage(testID, cpuUsage)
	if err != nil {
		log.Printf("Failed to log CPU usage: %v\n", err)
		return err
	}
	return nil
}

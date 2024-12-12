package data_management

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"synnergy_network/pkg/ledger"
	"time"
)

// RESCALE_FEATURES rescales data features to a specified range [min, max]
func RESCALE_FEATURES(data []float64, newMin, newMax float64) ([]float64, error) {
	if len(data) == 0 {
		return nil, errors.New("data is empty")
	}
	oldMin, oldMax := minMax(data)
	rescaled := make([]float64, len(data))
	for i, val := range data {
		rescaled[i] = newMin + (val-oldMin)*(newMax-newMin)/(oldMax-oldMin)
	}
	return rescaled, nil
}

// CALCULATE_WEIGHTED_AVERAGE calculates the weighted average of values based on corresponding weights
func CALCULATE_WEIGHTED_AVERAGE(values []float64, weights []float64) (float64, error) {
	if len(values) != len(weights) || len(values) == 0 {
		return 0, errors.New("values and weights must have the same non-zero length")
	}
	totalWeight := 0.0
	weightedSum := 0.0
	for i, val := range values {
		weightedSum += val * weights[i]
		totalWeight += weights[i]
	}
	return weightedSum / totalWeight, nil
}

// APPLY_POWER_TRANSFORM applies a power transformation (e.g., square root) to reduce skewness in data
func APPLY_POWER_TRANSFORM(data []float64, power float64) []float64 {
	transformed := make([]float64, len(data))
	for i, val := range data {
		transformed[i] = math.Pow(val, power)
	}
	return transformed
}

// CORRECT_SKEWNESS applies a log transformation to correct skewness in data
func CORRECT_SKEWNESS(data []float64) ([]float64, error) {
	corrected := make([]float64, len(data))
	for i, val := range data {
		if val <= 0 {
			return nil, errors.New("log transformation requires positive values")
		}
		corrected[i] = math.Log(val)
	}
	return corrected, nil
}

// IDENTIFY_MISSING_ROWS identifies rows with missing values in a dataset
func IDENTIFY_MISSING_ROWS(data [][]float64) []int {
	missingRows := []int{}
	for i, row := range data {
		for _, val := range row {
			if math.IsNaN(val) {
				missingRows = append(missingRows, i)
				break
			}
		}
	}
	return missingRows
}

// FREQUENCY_ENCODING encodes categorical features based on their frequency of occurrence
func FREQUENCY_ENCODING(data []string) map[string]int {
	frequency := make(map[string]int)
	for _, val := range data {
		frequency[val]++
	}
	return frequency
}

// BINARY_ENCODING encodes categorical variables into binary representation
func BINARY_ENCODING(data []string) map[string]string {
	uniqueValues := uniqueStrings(data)
	encoded := make(map[string]string)
	for i, val := range uniqueValues {
		encoded[val] = intToBinaryString(i, len(uniqueValues))
	}
	return encoded
}

// HASHING_ENCODING encodes categorical variables using a hashing function
func HASHING_ENCODING(data []string, numBuckets int) map[string]int {
	encoded := make(map[string]int)
	for _, val := range data {
		hash := sha256.Sum256([]byte(val))
		hashInt := int(hash[0]) % numBuckets
		encoded[val] = hashInt
	}
	return encoded
}

// REMOVE_HIGH_CARDINALITY_FEATURES removes features with high cardinality from a dataset
func REMOVE_HIGH_CARDINALITY_FEATURES(data [][]string, threshold int) [][]string {
	cardinality := make([]int, len(data[0]))
	for _, row := range data {
		for i, val := range row {
			cardinality[i]++
		}
	}

	filteredData := [][]string{}
	for _, row := range data {
		filteredRow := []string{}
		for i, val := range row {
			if cardinality[i] <= threshold {
				filteredRow = append(filteredRow, val)
			}
		}
		filteredData = append(filteredData, filteredRow)
	}
	return filteredData
}

// IDENTIFY_OUTLIER_BOUNDARIES calculates boundaries for outliers using a specified Z-score threshold
func IDENTIFY_OUTLIER_BOUNDARIES(data []float64, zThreshold float64) (float64, float64) {
	mean, stdDev := meanAndStdDev(data)
	lowerBound := mean - zThreshold*stdDev
	upperBound := mean + zThreshold*stdDev
	return lowerBound, upperBound
}

// MAP_LABELS maps categorical labels to numeric values based on a provided mapping
func MAP_LABELS(data []string, labelMapping map[string]int) []int {
	mapped := make([]int, len(data))
	for i, label := range data {
		mapped[i] = labelMapping[label]
	}
	return mapped
}

// RECODE_VARIABLES recodes categorical variables based on a provided dictionary
func RECODE_VARIABLES(data []string, recodeMap map[string]string) []string {
	recoded := make([]string, len(data))
	for i, val := range data {
		if newVal, exists := recodeMap[val]; exists {
			recoded[i] = newVal
		} else {
			recoded[i] = val
		}
	}
	return recoded
}

// GENERATE_DATE_FEATURES generates features from dates (e.g., year, month, day, weekday)
func GENERATE_DATE_FEATURES(dates []time.Time) []map[string]int {
	dateFeatures := make([]map[string]int, len(dates))
	for i, date := range dates {
		dateFeatures[i] = map[string]int{
			"year":    date.Year(),
			"month":   int(date.Month()),
			"day":     date.Day(),
			"weekday": int(date.Weekday()),
		}
	}
	return dateFeatures
}

// DATA_SUBSAMPLING performs random subsampling on data to create a smaller representative subset
func DATA_SUBSAMPLING(data []ledger.DataRecord, sampleSize int) ([]ledger.DataRecord, error) {
	if sampleSize > len(data) {
		return nil, errors.New("sample size cannot be larger than dataset size")
	}
	rand.Shuffle(len(data), func(i, j int) {
		data[i], data[j] = data[j], data[i]
	})
	return data[:sampleSize], nil
}



// Helper function: intToBinaryString converts an integer to a binary string of a specified length
func intToBinaryString(num int, length int) string {
	binaryStr := ""
	for i := length - 1; i >= 0; i-- {
		binaryStr += fmt.Sprintf("%d", (num>>i)&1)
	}
	return binaryStr
}

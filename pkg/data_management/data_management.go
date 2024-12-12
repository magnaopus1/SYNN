package data_management

import (
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

// GROUP_BY_FEATURE groups data records by a specified feature and returns grouped data
func GROUP_BY_FEATURE(data []common.DataRecord, feature string) map[string][]common.DataRecord {
	grouped := make(map[string][]common.DataRecord)
	for _, record := range data {
		key := record.Features[feature]
		grouped[key] = append(grouped[key], record)
	}
	return grouped
}

// FILTER_NOISE removes noise from data based on a threshold
func FILTER_NOISE(data []float64, threshold float64) []float64 {
	filtered := []float64{}
	for _, val := range data {
		if math.Abs(val) > threshold {
			filtered = append(filtered, val)
		}
	}
	return filtered
}

// WEIGHT_FEATURES applies weighting to features in data records
func WEIGHT_FEATURES(data []common.DataRecord, weights map[string]float64) []common.DataRecord {
	for i, record := range data {
		for feature, weight := range weights {
			if val, exists := record.Features[feature]; exists {
				record.Features[feature] = val * weight
			}
		}
		data[i] = record
	}
	return data
}

// DROP_CONSTANT_COLUMNS removes columns that have a constant value across all records
func DROP_CONSTANT_COLUMNS(data []common.DataRecord) []string {
	constantColumns := []string{}
	featureSets := make(map[string]map[float64]struct{})
	for _, record := range data {
		for feature, value := range record.Features {
			if _, exists := featureSets[feature]; !exists {
				featureSets[feature] = make(map[float64]struct{})
			}
			featureSets[feature][value] = struct{}{}
		}
	}
	for feature, values := range featureSets {
		if len(values) == 1 {
			constantColumns = append(constantColumns, feature)
		}
	}
	for i := range data {
		for _, feature := range constantColumns {
			delete(data[i].Features, feature)
		}
	}
	return constantColumns
}

// DROP_DUPLICATE_ROWS removes duplicate rows in a dataset based on feature values
func DROP_DUPLICATE_ROWS(data []common.DataRecord) []common.DataRecord {
	uniqueData := []common.DataRecord{}
	seen := make(map[string]struct{})
	for _, record := range data {
		hash := record.Hash()
		if _, exists := seen[hash]; !exists {
			seen[hash] = struct{}{}
			uniqueData = append(uniqueData, record)
		}
	}
	return uniqueData
}

// APPLY_COSINE_TRANSFORM applies a cosine transform to a data vector
func APPLY_COSINE_TRANSFORM(data []float64) []float64 {
	transformed := make([]float64, len(data))
	for i, val := range data {
		transformed[i] = math.Cos(val)
	}
	return transformed
}

// TEXT_STRIP_HTML removes HTML tags from text data
func TEXT_STRIP_HTML(text string) string {
	re := regexp.MustCompile(`<.*?>`)
	return re.ReplaceAllString(text, "")
}

// TEXT_DETOKENIZE converts a list of tokens back into a cohesive string
func TEXT_DETOKENIZE(tokens []string) string {
	return strings.Join(tokens, " ")
}

// REPLACE_SYNONYMS replaces synonyms in text with a standardized word
func REPLACE_SYNONYMS(text string, synonyms map[string]string) string {
	words := strings.Fields(text)
	for i, word := range words {
		if replacement, exists := synonyms[strings.ToLower(word)]; exists {
			words[i] = replacement
		}
	}
	return strings.Join(words, " ")
}

// IDENTIFY_MISSING_PATTERNS identifies patterns of missing data in the dataset
func IDENTIFY_MISSING_PATTERNS(data []common.DataRecord) map[string]int {
	missingPatterns := make(map[string]int)
	for _, record := range data {
		pattern := ""
		for _, val := range record.Features {
			if val == 0 {
				pattern += "0"
			} else {
				pattern += "1"
			}
		}
		missingPatterns[pattern]++
	}
	return missingPatterns
}

// CAP_OUTLIERS caps extreme values in data at a specified percentile
func CAP_OUTLIERS(data []float64, percentile float64) []float64 {
	if percentile < 0 || percentile > 100 {
		return data
	}
	sorted := append([]float64{}, data...)
	sort.Float64s(sorted)
	capIndex := int(float64(len(sorted)-1) * percentile / 100.0)
	capValue := sorted[capIndex]

	capped := make([]float64, len(data))
	for i, val := range data {
		if val > capValue {
			capped[i] = capValue
		} else {
			capped[i] = val
		}
	}
	return capped
}

// APPLY_TANH_TRANSFORM applies hyperbolic tangent transformation to data
func APPLY_TANH_TRANSFORM(data []float64) []float64 {
	transformed := make([]float64, len(data))
	for i, val := range data {
		transformed[i] = math.Tanh(val)
	}
	return transformed
}

// REMOVE_SPECIAL_CHARACTERS removes special characters from text
func REMOVE_SPECIAL_CHARACTERS(text string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9\s]+`)
	return re.ReplaceAllString(text, "")
}

// AGGREGATE_BY_MEAN aggregates data by computing the mean of each feature
func AGGREGATE_BY_MEAN(data []common.DataRecord) map[string]float64 {
	sums := make(map[string]float64)
	counts := make(map[string]int)

	for _, record := range data {
		for feature, value := range record.Features {
			sums[feature] += value
			counts[feature]++
		}
	}

	means := make(map[string]float64)
	for feature, sum := range sums {
		means[feature] = sum / float64(counts[feature])
	}
	return means
}

// AGGREGATE_BY_SUM aggregates data by computing the sum of each feature
func AGGREGATE_BY_SUM(data []common.DataRecord) map[string]float64 {
	sums := make(map[string]float64)
	for _, record := range data {
		for feature, value := range record.Features {
			sums[feature] += value
		}
	}
	return sums
}

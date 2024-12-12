package data_management

import (
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

// BINARIZE_DATA converts numeric data to binary values based on a threshold
func BINARIZE_DATA(data []float64, threshold float64) []int {
	binarized := make([]int, len(data))
	for i, val := range data {
		if val >= threshold {
			binarized[i] = 1
		} else {
			binarized[i] = 0
		}
	}
	return binarized
}

// DISCRETIZE_DATA groups continuous data into discrete bins
func DISCRETIZE_DATA(data []float64, numBins int) ([]int, error) {
	if numBins < 1 {
		return nil, errors.New("number of bins must be at least 1")
	}
	min, max := minMax(data)
	binSize := (max - min) / float64(numBins)
	discretized := make([]int, len(data))
	for i, val := range data {
		discretized[i] = int((val - min) / binSize)
		if discretized[i] >= numBins {
			discretized[i] = numBins - 1
		}
	}
	return discretized, nil
}

// SCALING_Z_SCORE standardizes data to have a mean of 0 and standard deviation of 1
func SCALING_Z_SCORE(data []float64) []float64 {
	mean, stdDev := meanAndStdDev(data)
	zScaled := make([]float64, len(data))
	for i, val := range data {
		zScaled[i] = (val - mean) / stdDev
	}
	return zScaled
}

// NORMALIZE_L2 applies L2 normalization, scaling data by the Euclidean norm
func NORMALIZE_L2(data []float64) []float64 {
	norm := 0.0
	for _, val := range data {
		norm += val * val
	}
	norm = math.Sqrt(norm)
	if norm == 0 {
		return data
	}
	normalized := make([]float64, len(data))
	for i, val := range data {
		normalized[i] = val / norm
	}
	return normalized
}

// NORMALIZE_L1 applies L1 normalization, scaling data by the Manhattan norm
func NORMALIZE_L1(data []float64) []float64 {
	norm := 0.0
	for _, val := range data {
		norm += math.Abs(val)
	}
	if norm == 0 {
		return data
	}
	normalized := make([]float64, len(data))
	for i, val := range data {
		normalized[i] = val / norm
	}
	return normalized
}

// POLYNOMIAL_FEATURE_EXPANSION generates polynomial features up to a specified degree
func POLYNOMIAL_FEATURE_EXPANSION(data []float64, degree int) [][]float64 {
	expanded := make([][]float64, len(data))
	for i, val := range data {
		features := make([]float64, degree)
		for j := 0; j < degree; j++ {
			features[j] = math.Pow(val, float64(j+1))
		}
		expanded[i] = features
	}
	return expanded
}

// TEXT_VECTORIZE converts text into a bag-of-words representation
func TEXT_VECTORIZE(text []string) map[string]int {
	wordCount := make(map[string]int)
	for _, sentence := range text {
		words := strings.Fields(sentence)
		for _, word := range words {
			wordCount[word]++
		}
	}
	return wordCount
}

// N_GRAM_GENERATION generates n-grams from text data
func N_GRAM_GENERATION(text string, n int) []string {
	words := strings.Fields(text)
	nGrams := []string{}
	for i := 0; i <= len(words)-n; i++ {
		nGrams = append(nGrams, strings.Join(words[i:i+n], " "))
	}
	return nGrams
}

// EXTRACT_NUMERIC_FEATURES extracts only numeric features from mixed data
func EXTRACT_NUMERIC_FEATURES(data map[string]interface{}) map[string]float64 {
	numericData := make(map[string]float64)
	for key, val := range data {
		if num, ok := val.(float64); ok {
			numericData[key] = num
		}
	}
	return numericData
}

// DATE_PARSING parses dates and extracts features such as year, month, day, etc.
func DATE_PARSING(dates []string, format string) ([]map[string]int, error) {
	parsedDates := make([]map[string]int, len(dates))
	for i, dateStr := range dates {
		date, err := time.Parse(format, dateStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse date %s: %v", dateStr, err)
		}
		parsedDates[i] = map[string]int{
			"year":    date.Year(),
			"month":   int(date.Month()),
			"day":     date.Day(),
			"weekday": int(date.Weekday()),
		}
	}
	return parsedDates, nil
}

// CATEGORICAL_ENCODING converts categorical data into integer labels
func CATEGORICAL_ENCODING(data []string) map[string]int {
	uniqueValues := uniqueStrings(data)
	encoded := make(map[string]int)
	for i, val := range uniqueValues {
		encoded[val] = i
	}
	return encoded
}

// BALANCE_CLASSES balances class distribution by duplicating samples from minority classes
func BALANCE_CLASSES(data []common.DataRecord, targetClass string) []common.DataRecord {
	classSamples := []common.DataRecord{}
	for _, record := range data {
		if record.Class == targetClass {
			classSamples = append(classSamples, record)
		}
	}
	balancedData := append(data, classSamples...)
	for len(balancedData) < 2*len(data) {
		balancedData = append(balancedData, classSamples[rand.Intn(len(classSamples))])
	}
	return balancedData
}

// SHUFFLE_DATA randomly shuffles data to remove any ordering bias
func SHUFFLE_DATA(data []common.DataRecord) {
	rand.Shuffle(len(data), func(i, j int) {
		data[i], data[j] = data[j], data[i]
	})
}

// IMPUTE_MEAN imputes missing values in data with the mean of non-missing values
func IMPUTE_MEAN(data []float64) []float64 {
	sum := 0.0
	count := 0
	for _, val := range data {
		if !math.IsNaN(val) {
			sum += val
			count++
		}
	}
	mean := sum / float64(count)
	imputed := make([]float64, len(data))
	for i, val := range data {
		if math.IsNaN(val) {
			imputed[i] = mean
		} else {
			imputed[i] = val
		}
	}
	return imputed
}

// IMPUTE_MEDIAN imputes missing values in data with the median of non-missing values
func IMPUTE_MEDIAN(data []float64) []float64 {
	validData := []float64{}
	for _, val := range data {
		if !math.IsNaN(val) {
			validData = append(validData, val)
		}
	}
	sort.Float64s(validData)
	median := validData[len(validData)/2]
	imputed := make([]float64, len(data))
	for i, val := range data {
		if math.IsNaN(val) {
			imputed[i] = median
		} else {
			imputed[i] = val
		}
	}
	return imputed
}

// Helper function: minMax calculates the minimum and maximum values of a slice of floats
func minMax(data []float64) (float64, float64) {
	min, max := data[0], data[0]
	for _, val := range data {
		if val < min {
			min = val
		}
		if val > max {
			max = val
		}
	}
	return min, max
}

// Helper function: meanAndStdDev calculates the mean and standard deviation of a slice of floats
func meanAndStdDev(data []float64) (float64, float64) {
	sum := 0.0
	for _, val := range data {
		sum += val
	}
	mean := sum / float64(len(data))
	variance := 0.0
	for _, val := range data {
		variance += math.Pow(val-mean, 2)
	}
	stdDev := math.Sqrt(variance / float64(len(data)))
	return mean, stdDev
}

// Helper function: uniqueStrings returns a slice of unique strings from input
func uniqueStrings(data []string) []string {
	uniqueMap := make(map[string]bool)
	uniqueValues := []string{}
	for _, item := range data {
		if !uniqueMap[item] {
			uniqueValues = append(uniqueValues, item)
			uniqueMap[item] = true
		}
	}
	return uniqueValues
}

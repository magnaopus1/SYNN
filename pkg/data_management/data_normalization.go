package data_management

import (
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

// NORMALIZE scales data to have a mean of 0 and a standard deviation of 1
func NORMALIZE(data []float64) ([]float64, error) {
	if len(data) == 0 {
		return nil, errors.New("data is empty")
	}
	mean, stdDev := meanAndStdDev(data)
	normalized := make([]float64, len(data))
	for i, val := range data {
		normalized[i] = (val - mean) / stdDev
	}
	return normalized, nil
}

// STANDARDIZE scales data to a given range [min, max]
func STANDARDIZE(data []float64, min, max float64) ([]float64, error) {
	if len(data) == 0 || max <= min {
		return nil, errors.New("invalid data or range for standardization")
	}
	oldMin, oldMax := minMax(data)
	scaled := make([]float64, len(data))
	for i, val := range data {
		scaled[i] = min + (val-oldMin)*(max-min)/(oldMax-oldMin)
	}
	return scaled, nil
}

// MIN_MAX_SCALING scales data between 0 and 1
func MIN_MAX_SCALING(data []float64) ([]float64, error) {
	return STANDARDIZE(data, 0, 1)
}

// ONE_HOT_ENCODE performs one-hot encoding for categorical data
func ONE_HOT_ENCODE(data []string) []map[string]int {
	uniqueValues := uniqueStrings(data)
	encoded := make([]map[string]int, len(data))
	for i, val := range data {
		encoding := make(map[string]int)
		for _, uniqueVal := range uniqueValues {
			if val == uniqueVal {
				encoding[uniqueVal] = 1
			} else {
				encoding[uniqueVal] = 0
			}
		}
		encoded[i] = encoding
	}
	return encoded
}

// LABEL_ENCODE converts categorical data to numerical labels
func LABEL_ENCODE(data []string) map[string]int {
	uniqueValues := uniqueStrings(data)
	labelEncoding := make(map[string]int)
	for i, val := range uniqueValues {
		labelEncoding[val] = i
	}
	return labelEncoding
}

// FILTER_OUTLIERS removes outliers outside a specified Z-score threshold
func FILTER_OUTLIERS(data []float64, zThreshold float64) []float64 {
	mean, stdDev := meanAndStdDev(data)
	filtered := []float64{}
	for _, val := range data {
		zScore := math.Abs((val - mean) / stdDev)
		if zScore <= zThreshold {
			filtered = append(filtered, val)
		}
	}
	return filtered
}

// CLEAN_WHITESPACE trims and removes extra whitespace from text data
func CLEAN_WHITESPACE(text string) string {
	return strings.Join(strings.Fields(text), " ")
}

// TOKENIZE_TEXT splits text into tokens
func TOKENIZE_TEXT(text string) []string {
	return strings.Fields(text)
}

// STOP_WORDS_REMOVAL removes common stop words from text
func STOP_WORDS_REMOVAL(text string, stopWords []string) string {
	wordList := strings.Fields(text)
	filteredWords := []string{}
	stopWordsMap := make(map[string]bool)
	for _, word := range stopWords {
		stopWordsMap[strings.ToLower(word)] = true
	}
	for _, word := range wordList {
		if !stopWordsMap[strings.ToLower(word)] {
			filteredWords = append(filteredWords, word)
		}
	}
	return strings.Join(filteredWords, " ")
}

// STEM_TEXT applies stemming to text using basic suffix removal
func STEM_TEXT(text string) string {
	words := strings.Fields(text)
	stemmedWords := []string{}
	for _, word := range words {
		stemmedWords = append(stemmedWords, stem(word))
	}
	return strings.Join(stemmedWords, " ")
}

// LEMMATIZE_TEXT lemmatizes text by mapping words to their base forms
func LEMMATIZE_TEXT(text string, lemmatizer map[string]string) string {
	words := strings.Fields(text)
	lemmatized := []string{}
	for _, word := range words {
		baseForm, exists := lemmatizer[strings.ToLower(word)]
		if exists {
			lemmatized = append(lemmatized, baseForm)
		} else {
			lemmatized = append(lemmatized, word)
		}
	}
	return strings.Join(lemmatized, " ")
}

// REMOVE_DUPLICATES removes duplicate entries from a slice of strings
func REMOVE_DUPLICATES(data []string) []string {
	uniqueMap := make(map[string]bool)
	uniqueData := []string{}
	for _, item := range data {
		if !uniqueMap[item] {
			uniqueData = append(uniqueData, item)
			uniqueMap[item] = true
		}
	}
	return uniqueData
}

// LOWERCASE_TEXT converts all text to lowercase
func LOWERCASE_TEXT(text string) string {
	return strings.ToLower(text)
}

// UPPERCASE_TEXT converts all text to uppercase
func UPPERCASE_TEXT(text string) string {
	return strings.ToUpper(text)
}

// REPLACE_MISSING_VALUES replaces missing values in data with a specified fill value
func REPLACE_MISSING_VALUES(data []float64, fillValue float64) []float64 {
	replaced := make([]float64, len(data))
	for i, val := range data {
		if math.IsNaN(val) {
			replaced[i] = fillValue
		} else {
			replaced[i] = val
		}
	}
	return replaced
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

// Helper function: stem applies basic suffix stripping for stemming words
func stem(word string) string {
	suffixes := []string{"ing", "ed", "es", "s"}
	for _, suffix := range suffixes {
		if strings.HasSuffix(word, suffix) {
			return strings.TrimSuffix(word, suffix)
		}
	}
	return word
}

package data_management

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strings"
)

// IMPUTE_MODE imputes missing values with the mode (most frequent value)
func IMPUTE_MODE(data []string) string {
	frequency := make(map[string]int)
	for _, val := range data {
		if val != "" {
			frequency[val]++
		}
	}
	mode := ""
	maxCount := 0
	for val, count := range frequency {
		if count > maxCount {
			mode = val
			maxCount = count
		}
	}
	return mode
}

// NORMALIZE_MAX_ABS scales data to a range based on maximum absolute value normalization
func NORMALIZE_MAX_ABS(data []float64) []float64 {
	maxAbs := 0.0
	for _, val := range data {
		if math.Abs(val) > maxAbs {
			maxAbs = math.Abs(val)
		}
	}
	if maxAbs == 0 {
		return data
	}
	normalized := make([]float64, len(data))
	for i, val := range data {
		normalized[i] = val / maxAbs
	}
	return normalized
}

// REMOVE_PUNCTUATION removes punctuation from text data
func REMOVE_PUNCTUATION(text string) string {
	re := regexp.MustCompile(`[[:punct:]]`)
	return re.ReplaceAllString(text, "")
}

// APPLY_LOG_TRANSFORM applies logarithmic transformation to continuous data
func APPLY_LOG_TRANSFORM(data []float64) ([]float64, error) {
	transformed := make([]float64, len(data))
	for i, val := range data {
		if val <= 0 {
			return nil, errors.New("log transform requires positive values")
		}
		transformed[i] = math.Log(val)
	}
	return transformed, nil
}

// APPLY_EXP_TRANSFORM applies exponential transformation to data
func APPLY_EXP_TRANSFORM(data []float64) []float64 {
	transformed := make([]float64, len(data))
	for i, val := range data {
		transformed[i] = math.Exp(val)
	}
	return transformed
}

// INTERPOLATE_MISSING_VALUES interpolates missing values in numeric data linearly
func INTERPOLATE_MISSING_VALUES(data []float64) []float64 {
	interpolated := make([]float64, len(data))
	copy(interpolated, data)
	for i := 0; i < len(data); i++ {
		if math.IsNaN(data[i]) {
			start, end := i-1, i+1
			for start >= 0 && math.IsNaN(data[start]) {
				start--
			}
			for end < len(data) && math.IsNaN(data[end]) {
				end++
			}
			if start >= 0 && end < len(data) {
				interpolated[i] = (data[start] + data[end]) / 2
			} else if start >= 0 {
				interpolated[i] = data[start]
			} else if end < len(data) {
				interpolated[i] = data[end]
			}
		}
	}
	return interpolated
}

// QUANTILE_TRANSFORM applies quantile normalization to data
func QUANTILE_TRANSFORM(data []float64) []float64 {
	sorted := append([]float64{}, data...)
	sort.Float64s(sorted)
	quantileTransformed := make([]float64, len(data))
	for i, val := range data {
		pos := sort.SearchFloat64s(sorted, val)
		quantileTransformed[i] = float64(pos) / float64(len(data)-1)
	}
	return quantileTransformed
}

// RANK_TRANSFORM applies rank transformation to data
func RANK_TRANSFORM(data []float64) []int {
	ranked := make([]int, len(data))
	sorted := append([]float64{}, data...)
	sort.Float64s(sorted)
	for i, val := range data {
		ranked[i] = sort.SearchFloat64s(sorted, val) + 1
	}
	return ranked
}

// SHIFT_DATA shifts data by a specified amount
func SHIFT_DATA(data []float64, shiftAmount float64) []float64 {
	shifted := make([]float64, len(data))
	for i, val := range data {
		shifted[i] = val + shiftAmount
	}
	return shifted
}

// SMOOTH_DATA applies a simple moving average to smooth the data
func SMOOTH_DATA(data []float64, windowSize int) []float64 {
	if windowSize < 1 {
		return data
	}
	smoothed := make([]float64, len(data))
	for i := 0; i < len(data); i++ {
		start := int(math.Max(0, float64(i-windowSize/2)))
		end := int(math.Min(float64(len(data)), float64(i+windowSize/2+1)))
		sum := 0.0
		for j := start; j < end; j++ {
			sum += data[j]
		}
		smoothed[i] = sum / float64(end-start)
	}
	return smoothed
}

// DISCRETIZE_CONTINUOUS bins continuous data into discrete intervals
func DISCRETIZE_CONTINUOUS(data []float64, numBins int) ([]int, error) {
	if numBins < 1 {
		return nil, errors.New("number of bins must be at least 1")
	}
	min, max := data[0], data[0]
	for _, val := range data {
		if val < min {
			min = val
		}
		if val > max {
			max = val
		}
	}
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

// APPLY_SQRT_TRANSFORM applies square root transformation to data
func APPLY_SQRT_TRANSFORM(data []float64) []float64 {
	transformed := make([]float64, len(data))
	for i, val := range data {
		transformed[i] = math.Sqrt(val)
	}
	return transformed
}

// TRIM_TEXT trims whitespace and extra spaces from text data
func TRIM_TEXT(text string) string {
	return strings.TrimSpace(text)
}

// FORMAT_NUMERICS formats numeric data to a specified precision
func FORMAT_NUMERICS(data []float64, precision int) []string {
	formatted := make([]string, len(data))
	format := fmt.Sprintf("%%.%df", precision)
	for i, val := range data {
		formatted[i] = fmt.Sprintf(format, val)
	}
	return formatted
}

// FORMAT_CATEGORICAL standardizes categorical data to title case
func FORMAT_CATEGORICAL(data []string) []string {
	formatted := make([]string, len(data))
	for i, val := range data {
		formatted[i] = strings.Title(strings.ToLower(strings.TrimSpace(val)))
	}
	return formatted
}

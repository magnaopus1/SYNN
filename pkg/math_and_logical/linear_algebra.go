package math_and_logical

import (
	"errors"
)

// MatrixMultiply performs matrix multiplication between two matrices
func MatrixMultiply(a, b [][]float64) ([][]float64, error) {
	if len(a) == 0 || len(b) == 0 || len(a[0]) != len(b) {
		return nil, errors.New("incompatible dimensions for matrix multiplication")
	}

	result := make([][]float64, len(a))
	for i := range result {
		result[i] = make([]float64, len(b[0]))
		for j := range result[i] {
			for k := range b {
				result[i][j] += a[i][k] * b[k][j]
			}
		}
	}
	return result, nil
}

// MatrixInverse calculates the inverse of a 2x2 or 3x3 matrix
func MatrixInverse(matrix [][]float64) ([][]float64, error) {
	if len(matrix) == 2 && len(matrix[0]) == 2 {
		// Inverse for 2x2 matrix
		det := matrix[0][0]*matrix[1][1] - matrix[0][1]*matrix[1][0]
		if det == 0 {
			return nil, errors.New("matrix is singular and cannot be inverted")
		}
		return [][]float64{
			{matrix[1][1] / det, -matrix[0][1] / det},
			{-matrix[1][0] / det, matrix[0][0] / det},
		}, nil
	} else if len(matrix) == 3 && len(matrix[0]) == 3 {
		// Inverse for 3x3 matrix
		det := MatrixDeterminant(matrix)
		if det == 0 {
			return nil, errors.New("matrix is singular and cannot be inverted")
		}
		inv := make([][]float64, 3)
		for i := range inv {
			inv[i] = make([]float64, 3)
		}
		inv[0][0] = (matrix[1][1]*matrix[2][2] - matrix[1][2]*matrix[2][1]) / det
		inv[0][1] = (matrix[0][2]*matrix[2][1] - matrix[0][1]*matrix[2][2]) / det
		inv[0][2] = (matrix[0][1]*matrix[1][2] - matrix[0][2]*matrix[1][1]) / det
		inv[1][0] = (matrix[1][2]*matrix[2][0] - matrix[1][0]*matrix[2][2]) / det
		inv[1][1] = (matrix[0][0]*matrix[2][2] - matrix[0][2]*matrix[2][0]) / det
		inv[1][2] = (matrix[0][2]*matrix[1][0] - matrix[0][0]*matrix[1][2]) / det
		inv[2][0] = (matrix[1][0]*matrix[2][1] - matrix[1][1]*matrix[2][0]) / det
		inv[2][1] = (matrix[0][1]*matrix[2][0] - matrix[0][0]*matrix[2][1]) / det
		inv[2][2] = (matrix[0][0]*matrix[1][1] - matrix[0][1]*matrix[1][0]) / det
		return inv, nil
	}
	return nil, errors.New("only 2x2 and 3x3 matrices are supported for inversion")
}

// MatrixTranspose returns the transpose of a matrix
func MatrixTranspose(matrix [][]float64) [][]float64 {
	rows := len(matrix)
	cols := len(matrix[0])
	transpose := make([][]float64, cols)
	for i := range transpose {
		transpose[i] = make([]float64, rows)
		for j := range transpose[i] {
			transpose[i][j] = matrix[j][i]
		}
	}
	return transpose
}

// MatrixDeterminant calculates the determinant of a 2x2 or 3x3 matrix
func MatrixDeterminant(matrix [][]float64) float64 {
	if len(matrix) == 2 && len(matrix[0]) == 2 {
		return matrix[0][0]*matrix[1][1] - matrix[0][1]*matrix[1][0]
	} else if len(matrix) == 3 && len(matrix[0]) == 3 {
		return matrix[0][0]*(matrix[1][1]*matrix[2][2]-matrix[1][2]*matrix[2][1]) -
			matrix[0][1]*(matrix[1][0]*matrix[2][2]-matrix[1][2]*matrix[2][0]) +
			matrix[0][2]*(matrix[1][0]*matrix[2][1]-matrix[1][1]*matrix[2][0])
	}
	return 0
}

// VectorDotProduct calculates the dot product of two vectors
func VectorDotProduct(a, b []float64) (float64, error) {
	if len(a) != len(b) {
		return 0, errors.New("vectors must have the same length for dot product")
	}
	result := 0.0
	for i := range a {
		result += a[i] * b[i]
	}
	return result, nil
}

// VectorCrossProduct calculates the cross product of two 3D vectors
func VectorCrossProduct(a, b []float64) ([]float64, error) {
	if len(a) != 3 || len(b) != 3 {
		return nil, errors.New("cross product is defined for 3D vectors only")
	}
	return []float64{
		a[1]*b[2] - a[2]*b[1],
		a[2]*b[0] - a[0]*b[2],
		a[0]*b[1] - a[1]*b[0],
	}, nil
}

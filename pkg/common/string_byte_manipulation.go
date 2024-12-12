package common

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

// String Manipulation Functions

// ConcatStrings concatenates two or more strings into a single string.
func ConcatStrings(parts ...string) string {
	var builder strings.Builder
	for _, part := range parts {
		builder.WriteString(part)
	}
	return builder.String()
}

// Substring extracts a substring from a given string, starting at 'start' and of length 'length'.
// If 'length' exceeds the string length, it returns up to the end of the string.
func Substring(s string, start, length int) (string, error) {
	if start < 0 || start >= len(s) || length < 0 {
		return "", errors.New("invalid start or length for substring")
	}
	end := start + length
	if end > len(s) {
		end = len(s)
	}
	return s[start:end], nil
}

// SplitString splits a string by a specified delimiter and returns an array of substrings.
func SplitString(s, delimiter string) []string {
	return strings.Split(s, delimiter)
}

// ReplaceString replaces occurrences of old substring with new substring in a given string.
func ReplaceString(s, old, new string) string {
	return strings.ReplaceAll(s, old, new)
}

// ToUpperCase converts a string to uppercase.
func ToUpperCase(s string) string {
	return strings.ToUpper(s)
}

// ToLowerCase converts a string to lowercase.
func ToLowerCase(s string) string {
	return strings.ToLower(s)
}

// Byte Manipulation Functions

// ByteSlice extracts a slice from a byte array, starting at 'start' index and of length 'length'.
// If the length exceeds the array bounds, it returns up to the end of the array.
func ByteSlice(b []byte, start, length int) ([]byte, error) {
	if start < 0 || start >= len(b) || length < 0 {
		return nil, errors.New("invalid start or length for byte slice")
	}
	end := start + length
	if end > len(b) {
		end = len(b)
	}
	return b[start:end], nil
}

// MergeByteArrays concatenates two or more byte arrays into a single byte array.
func MergeByteArrays(arrays ...[]byte) []byte {
	return bytes.Join(arrays, []byte{})
}

// PadBytes pads a byte array to a specified length with a given padding byte.
// If the array is longer than the specified length, it returns the original array.
func PadBytes(b []byte, length int, padByte byte) []byte {
	if len(b) >= length {
		return b
	}
	padding := bytes.Repeat([]byte{padByte}, length-len(b))
	return append(b, padding...)
}

// TrimBytes removes leading and trailing occurrences of a specific byte from a byte array.
func TrimBytes(b []byte, trimByte byte) []byte {
	return bytes.Trim(b, string(trimByte))
}

// ReplaceBytes replaces all occurrences of old byte with new byte in a given byte array.
func ReplaceBytes(b []byte, old, new byte) []byte {
	return bytes.ReplaceAll(b, []byte{old}, []byte{new})
}

// Advanced Byte Operations

// ToHexString converts a byte array to a hexadecimal string representation.
func ToHexString(b []byte) string {
	return fmt.Sprintf("%x", b)
}

// FromHexString converts a hexadecimal string to a byte array.
func FromHexString(hexStr string) ([]byte, error) {
	return hex.DecodeString(hexStr)
}

// Math-Related Byte Manipulation

// IncrementByteArray increments the value of a byte array interpreted as a big-endian integer.
func IncrementByteArray(b []byte) []byte {
	for i := len(b) - 1; i >= 0; i-- {
		if b[i] < 0xFF {
			b[i]++
			return b
		}
		b[i] = 0
	}
	return append([]byte{1}, b...)
}

// DecrementByteArray decrements the value of a byte array interpreted as a big-endian integer.
func DecrementByteArray(b []byte) ([]byte, error) {
	if len(b) == 0 || (len(b) == 1 && b[0] == 0) {
		return nil, errors.New("cannot decrement zero byte array")
	}
	for i := len(b) - 1; i >= 0; i-- {
		if b[i] > 0x00 {
			b[i]--
			if b[0] == 0 {
				b = b[1:]
			}
			return b, nil
		}
		b[i] = 0xFF
	}
	return b, nil
}

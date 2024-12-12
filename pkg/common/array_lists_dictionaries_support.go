package common

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
)

// Supported data types
type DataType int

const (
	Integer DataType = iota
	String
	Boolean
)

// Array represents a dynamically sized array structure with encryption support
type Array struct {
	Data          []interface{}
	Encrypted     bool
	encryptionKey []byte
	mutex         sync.Mutex
}

// NewArray initializes a new array, optionally with encryption
func NewArray(encrypted bool, key []byte) *Array {
	return &Array{
		Data:          make([]interface{}, 0),
		Encrypted:     encrypted,
		encryptionKey: key,
	}
}

// Append adds an element to the end of the array
func (a *Array) Append(value interface{}) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if a.Encrypted {
		encryptedValue, err := encryptValue(value, a.encryptionKey)
		if err != nil {
			return err
		}
		a.Data = append(a.Data, encryptedValue)
	} else {
		a.Data = append(a.Data, value)
	}
	return nil
}

// Get retrieves an element at a specific index
func (a *Array) Get(index int) (interface{}, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if index < 0 || index >= len(a.Data) {
		return nil, errors.New("index out of bounds")
	}

	if a.Encrypted {
		return decryptValue(a.Data[index], a.encryptionKey)
	}
	return a.Data[index], nil
}

// List represents a list structure with dynamic resizing and optional encryption support
type List struct {
	Elements      []interface{}
	Encrypted     bool
	encryptionKey []byte
	mutex         sync.Mutex
}

// NewList initializes a new list
func NewList(encrypted bool, key []byte) *List {
	return &List{
		Elements:      make([]interface{}, 0),
		Encrypted:     encrypted,
		encryptionKey: key,
	}
}

// Add adds an element to the list
func (l *List) Add(element interface{}) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.Encrypted {
		encryptedElement, err := encryptValue(element, l.encryptionKey)
		if err != nil {
			return err
		}
		l.Elements = append(l.Elements, encryptedElement)
	} else {
		l.Elements = append(l.Elements, element)
	}
	return nil
}

// Get retrieves an element by index
func (l *List) Get(index int) (interface{}, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if index < 0 || index >= len(l.Elements) {
		return nil, errors.New("index out of bounds")
	}

	if l.Encrypted {
		return decryptValue(l.Elements[index], l.encryptionKey)
	}
	return l.Elements[index], nil
}

// Dictionary represents a key-value store with encryption support for values
type Dictionary struct {
	Entries       map[string]interface{}
	Encrypted     bool
	encryptionKey []byte
	mutex         sync.Mutex
}

// NewDictionary initializes a new dictionary with optional encryption support for values
func NewDictionary(encrypted bool, key []byte) *Dictionary {
	return &Dictionary{
		Entries:       make(map[string]interface{}),
		Encrypted:     encrypted,
		encryptionKey: key,
	}
}

// Set adds or updates a key-value pair in the dictionary
func (d *Dictionary) Set(key string, value interface{}) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.Encrypted {
		encryptedValue, err := encryptValue(value, d.encryptionKey)
		if err != nil {
			return err
		}
		d.Entries[key] = encryptedValue
	} else {
		d.Entries[key] = value
	}
	return nil
}

// Get retrieves the value for a given key
func (d *Dictionary) Get(key string) (interface{}, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	value, exists := d.Entries[key]
	if !exists {
		return nil, errors.New("key does not exist in dictionary")
	}

	if d.Encrypted {
		return decryptValue(value, d.encryptionKey)
	}
	return value, nil
}

// Remove deletes a key-value pair from the dictionary
func (d *Dictionary) Remove(key string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if _, exists := d.Entries[key]; !exists {
		return errors.New("key does not exist in dictionary")
	}
	delete(d.Entries, key)
	return nil
}


// encryptValue encrypts a given value using the encryption package
func encryptValue(value interface{}, key []byte) (string, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	encryptedData, err := EncryptWithKey(data, key)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", encryptedData), nil
}

// decryptValue decrypts a previously encrypted value
func decryptValue(encryptedValue interface{}, key []byte) (interface{}, error) {
	encryptedStr, ok := encryptedValue.(string)
	if !ok {
		return nil, errors.New("encrypted value is not a valid string")
	}
	encryptedData, err := hex.DecodeString(encryptedStr)
	if err != nil {
		return nil, err
	}
	decryptedData, err := DecryptWithKey(encryptedData, key)
	if err != nil {
		return nil, err
	}
	var value interface{}
	if err = json.Unmarshal(decryptedData, &value); err != nil {
		return nil, err
	}
	return value, nil
}


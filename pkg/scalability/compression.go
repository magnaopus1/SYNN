package scalability

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"errors"
	"fmt"
	"io"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewCompressionSystem initializes the compression system
func NewCompressionSystem(ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption) *common.CompressionSystem {
	return &common.CompressionSystem{
		Ledger:            ledgerInstance,
		EncryptionService: encryptionService,
	}
}

// AdaptiveCompression applies adaptive compression based on the size and type of the data
func (cs *common.CompressionSystem) AdaptiveCompression(data []byte, dataType string) ([]byte, error) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	var compressedData []byte
	var err error

	// Choose compression method based on data type and size
	switch dataType {
	case "block", "sub-block":
		if len(data) > 100000 { // If the data is large, use gzip
			compressedData, err = cs.GzipCompress(data)
		} else { // Otherwise, use zlib
			compressedData, err = cs.ZlibCompress(data)
		}
	case "transaction", "state":
		compressedData, err = cs.ZlibCompress(data) // Prefer zlib for transaction/state data
	case "file":
		compressedData, err = cs.GzipCompress(data) // Use gzip for files
	default:
		return nil, errors.New("unsupported data type for compression")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to compress data: %v", err)
	}

	// Log the compression action
	err = cs.Ledger.RecordCompression(dataType, len(data), len(compressedData), "adaptive")
	if err != nil {
		return nil, fmt.Errorf("failed to log compression: %v", err)
	}

	return compressedData, nil
}

// ZlibCompress compresses data using zlib
func (cs *common.CompressionSystem) ZlibCompress(data []byte) ([]byte, error) {
	var buffer bytes.Buffer
	writer := zlib.NewWriter(&buffer)
	_, err := writer.Write(data)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// GzipCompress compresses data using gzip
func (cs *common.CompressionSystem) GzipCompress(data []byte) ([]byte, error) {
	var buffer bytes.Buffer
	writer := gzip.NewWriter(&buffer)
	_, err := writer.Write(data)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// DecompressData decompresses data based on the compression type
func (cs *common.CompressionSystem) DecompressData(compressedData []byte, compressionType string) ([]byte, error) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	var decompressedData []byte
	var err error

	switch compressionType {
	case "zlib":
		decompressedData, err = cs.ZlibDecompress(compressedData)
	case "gzip":
		decompressedData, err = cs.GzipDecompress(compressedData)
	default:
		return nil, errors.New("unsupported compression type for decompression")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to decompress data: %v", err)
	}

	return decompressedData, nil
}

// ZlibDecompress decompresses data using zlib
func (cs *common.CompressionSystem) ZlibDecompress(compressedData []byte) ([]byte, error) {
	reader, err := zlib.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	var decompressedData bytes.Buffer
	_, err = io.Copy(&decompressedData, reader)
	if err != nil {
		return nil, err
	}

	return decompressedData.Bytes(), nil
}

// GzipDecompress decompresses data using gzip
func (cs *common.CompressionSystem) GzipDecompress(compressedData []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	var decompressedData bytes.Buffer
	_, err = io.Copy(&decompressedData, reader)
	if err != nil {
		return nil, err
	}

	return decompressedData.Bytes(), nil
}

// CompressBlock compresses a full block using adaptive compression
func (cs *common.CompressionSystem) CompressBlock(blockData []byte) ([]byte, error) {
	return cs.AdaptiveCompression(blockData, "block")
}

// CompressSubBlock compresses a sub-block using adaptive compression
func (cs *common.CompressionSystem) CompressSubBlock(subBlockData []byte) ([]byte, error) {
	return cs.AdaptiveCompression(subBlockData, "sub-block")
}

// CompressTransaction compresses a transaction using zlib
func (cs *common.CompressionSystem) CompressTransaction(transactionData []byte) ([]byte, error) {
	return cs.ZlibCompress(transactionData)
}

// CompressState compresses state data using zlib
func (cs *common.CompressionSystem) CompressState(stateData []byte) ([]byte, error) {
	return cs.ZlibCompress(stateData)
}

// CompressFile compresses a file using gzip
func (cs *common.CompressionSystem) CompressFile(fileData []byte) ([]byte, error) {
	return cs.GzipCompress(fileData)
}

// LogCompression logs the compression activity in the ledger
func (cs *common.CompressionSystem) LogCompression(originalSize int, compressedSize int, compressionType string, dataType string) error {
	err := cs.Ledger.RecordCompression(dataType, originalSize, compressedSize, compressionType)
	if err != nil {
		return fmt.Errorf("failed to log compression: %v", err)
	}
	return nil
}

// CompressSubBlock applies compression to a sub-block using zlib or gzip
func (cs *common.CompressionSystem) CompressSubBlock(subBlockData []byte) ([]byte, error) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	// Choose compression method based on size
	var compressedData []byte
	var err error
	if len(subBlockData) > 100000 { // Use gzip for larger sub-blocks
		compressedData, err = cs.GzipCompress(subBlockData)
	} else {
		compressedData, err = cs.ZlibCompress(subBlockData)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to compress sub-block: %v", err)
	}

	// Log the compression in the ledger
	err = cs.Ledger.RecordCompression("sub-block", len(subBlockData), len(compressedData), "adaptive")
	if err != nil {
		return nil, fmt.Errorf("failed to log sub-block compression: %v", err)
	}

	fmt.Printf("Sub-block compressed from %d to %d bytes\n", len(subBlockData), len(compressedData))
	return compressedData, nil
}

// CompressBlock compresses a full block using zlib or gzip
func (cs *common.CompressionSystem) CompressBlock(blockData []byte) ([]byte, error) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	// Choose compression method based on block size
	var compressedData []byte
	var err error
	if len(blockData) > 200000 { // Use gzip for larger blocks
		compressedData, err = cs.GzipCompress(blockData)
	} else {
		compressedData, err = cs.ZlibCompress(blockData)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to compress block: %v", err)
	}

	// Log the compression in the ledger
	err = cs.Ledger.RecordCompression("block", len(blockData), len(compressedData), "adaptive")
	if err != nil {
		return nil, fmt.Errorf("failed to log block compression: %v", err)
	}

	fmt.Printf("Block compressed from %d to %d bytes\n", len(blockData), len(compressedData))
	return compressedData, nil
}

// CompressTransaction compresses a transaction using zlib
func (cs *common.CompressionSystem) CompressTransaction(transactionData []byte) ([]byte, error) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	// Use zlib compression for transaction data
	compressedData, err := cs.ZlibCompress(transactionData)
	if err != nil {
		return nil, fmt.Errorf("failed to compress transaction: %v", err)
	}

	// Log the compression in the ledger
	err = cs.Ledger.RecordCompression("transaction", len(transactionData), len(compressedData), "zlib")
	if err != nil {
		return nil, fmt.Errorf("failed to log transaction compression: %v", err)
	}

	fmt.Printf("Transaction compressed from %d to %d bytes\n", len(transactionData), len(compressedData))
	return compressedData, nil
}

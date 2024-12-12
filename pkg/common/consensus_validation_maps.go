package common

type TransactionTypeMap struct {
	Mapping map[string]string
}

// NewTransactionTypeMap initializes the TransactionTypeMap with default values.
func NewTransactionTypeMap() *TransactionTypeMap {
	return &TransactionTypeMap{
		Mapping: map[string]string{
			"SYN20":        	"SYN20",
			"SYN131":       	"SYN131",
			"CROSS_CHAIN":  	"CROSS_CHAIN",
			"SYN500":       	"SYN500",
			"SYN10":       	 	"SYN10",
			"SYN200":      	 	"SYN200",
			"SYN721":       	"SYN721",
			"SYN1155":      	"SYN1155",
			"CUSTOM":       	"CUSTOM",
			"PROPOSAL":       	"PROPOSAL",
			"SMART_CONTRACT":	"SMART_CONTRACT",
		},
	}
}



type TransactionFunctionMap struct {
	Mapping map[string]string
}

// NewTransactionFunctionMap initializes the TransactionFunctionMap with default values.
func NewTransactionFunctionMap() *TransactionFunctionMap {
	return &TransactionFunctionMap{
		Mapping: map[string]string{
			"Transfer":          "Trans",
			"OwnershipTransfer": "OwTrans",
			"Mint":              "Mint",
			"Burn":              "Burn",
			"Stake":             "Stake",
			"Unstake":           "Unstake",
			"ClaimReward":       "ClaimReward",
			"Lease":             "Lease",
			"Rent":              "Rent",
			"Purchase":          "Purchase",
			"Fractionalize":     "Frac",
			"CoOwnership":       "CoOwn",
			"ComplianceCheck":   "CompCheck",
			"Vote":              "Vote",
			"CrossChainTransfer": "CrossTrans",
			"CrossChainMint":     "CrossMint",
			"CrossChainBurn":     "CrossBurn",
		},
	}
}



type ValidationMap struct {
	Mapping map[string]func(tx Transaction) bool
}

// NewValidationMap initializes the ValidationMap with default validation logic.
func NewValidationMap() *ValidationMap {
	return &ValidationMap{
		Mapping: map[string]func(tx Transaction) bool{
			"SYN20:Trans":         validateSYN20Transfer,
			"SYN20:Mint":          validateSYN20Mint,

		},
	}
}

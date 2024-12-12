package smart_contract

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "time"
    "errors"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

// NewSmartContractTemplateMarketplace creates a new marketplace
func NewSmartContractTemplateMarketplace(ledgerInstance *ledger.Ledger) *SmartContractTemplateMarketplace {
    return &SmartContractTemplateMarketplace{
        Templates:     make(map[string]*SmartContractTemplate),
        LedgerInstance: ledgerInstance,
    }
}

// UploadTemplate allows a user to upload a smart contract template to the marketplace
func (mp *SmartContractTemplateMarketplace) UploadTemplate(creator, name, description, code string, price float64) (string, error) {
    mp.mutex.Lock()
    defer mp.mutex.Unlock()

    templateID := mp.generateTemplateID(creator, name)
    
    if _, exists := mp.Templates[templateID]; exists {
        return "", fmt.Errorf("template with ID %s already exists", templateID)
    }

    encryptedCode, err := common.EncryptData(code, common.EncryptionKey)
    if err != nil {
        return "", fmt.Errorf("failed to encrypt template code: %v", err)
    }

    newTemplate := &SmartContractTemplate{
        ID:            templateID,
        Name:          name,
        Description:   description,
        Creator:       creator,
        Code:          code,
        Price:         price,
        Timestamp:     time.Now(),
        EncryptedCode: encryptedCode,
    }

    mp.Templates[templateID] = newTemplate

    // Record the template upload in the ledger
    err = mp.LedgerInstance.RecordTemplateUpload(newTemplate.ID, encryptedCode)
    if err != nil {
        return "", fmt.Errorf("failed to record template upload in the ledger: %v", err)
    }

    fmt.Printf("Template %s uploaded by %s.\n", templateID, creator)
    return templateID, nil
}

// PurchaseTemplate allows a user to purchase a template from the marketplace
func (mp *SmartContractTemplateMarketplace) PurchaseTemplate(buyer, templateID string) (string, error) {
    mp.mutex.Lock()
    defer mp.mutex.Unlock()

    template, exists := mp.Templates[templateID]
    if !exists {
        return "", errors.New("template not found")
    }

    // Record the purchase transaction on the ledger
    err := mp.LedgerInstance.RecordTransaction(buyer, template.Creator, template.Price)
    if err != nil {
        return "", fmt.Errorf("failed to record purchase transaction in the ledger: %v", err)
    }

    fmt.Printf("Template %s purchased by %s for %.2f SYNN.\n", templateID, buyer, template.Price)
    return template.EncryptedCode, nil
}

// QueryTemplates allows users to search for available templates in the marketplace
func (mp *SmartContractTemplateMarketplace) QueryTemplates() []*SmartContractTemplate {
    mp.mutex.Lock()
    defer mp.mutex.Unlock()

    templates := make([]*SmartContractTemplate, 0, len(mp.Templates))
    for _, template := range mp.Templates {
        templates = append(templates, template)
    }
    return templates
}

// DeployTemplate allows a user to deploy a purchased template as a smart contract
func (scm *common.SmartContractManager) DeployTemplate(buyer, encryptedCode string, parameters map[string]interface{}) (*common.SmartContract, error) {
    scm.mutex.Lock()
    defer scm.mutex.Unlock()

    decryptedCode, err := encryption.DecryptData(encryptedCode, common.EncryptionKey)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt template code: %v", err)
    }

    contractID := generateContractID(buyer, decryptedCode)
    contract := &common.SmartContract{
        ID:             contractID,
        Code:           decryptedCode,
        Parameters:     parameters,
        State:          make(map[string]interface{}),
        Owner:          buyer,
        LedgerInstance: scm.LedgerInstance,
    }

    scm.Contracts[contract.ID] = contract

    // Record the contract deployment in the ledger
    encryptedContract, err := encryption.EncryptContract(contract, common.EncryptionKey)
    if err != nil {
        return nil, fmt.Errorf("failed to encrypt contract: %v", err)
    }

    err = scm.LedgerInstance.RecordContractDeployment(contract.ID, encryptedContract)
    if err != nil {
        return nil, fmt.Errorf("failed to record contract deployment in the ledger: %v", err)
    }

    fmt.Printf("Smart Contract %s deployed by %s.\n", contract.ID, buyer)
    return contract, nil
}

// generateTemplateID creates a unique ID for a template based on the creator and name
func (mp *common.SmartContractTemplateMarketplace) generateTemplateID(creator, name string) string {
    hashInput := fmt.Sprintf("%s%s%d", creator, name, time.Now().UnixNano())
    hash := sha256.New()
    hash.Write([]byte(hashInput))
    return hex.EncodeToString(hash.Sum(nil))
}



// UploadTemplate allows a user to upload a smart contract template to the marketplace
func (mp *common.SmartContractTemplateMarketplace) UploadTemplate(creator, name, description, code string, price float64) (string, error) {
    mp.mutex.Lock()
    defer mp.mutex.Unlock()

    templateID := mp.generateTemplateID(creator, name)

    if _, exists := mp.Templates[templateID]; exists {
        return "", fmt.Errorf("template with ID %s already exists", templateID)
    }

    encryptedCode, err := encryption.EncryptData(code, common.EncryptionKey)
    if err != nil {
        return "", fmt.Errorf("failed to encrypt template code: %v", err)
    }

    newTemplate := &common.SmartContractTemplate{
        ID:            templateID,
        Name:          name,
        Description:   description,
        Creator:       creator,
        Code:          code,
        Price:         price,
        Timestamp:     time.Now(),
        EncryptedCode: encryptedCode,
    }

    mp.Templates[templateID] = newTemplate

    // Record the template upload in the ledger
    err = mp.LedgerInstance.RecordTemplateUpload(newTemplate.ID, encryptedCode)
    if err != nil {
        return "", fmt.Errorf("failed to record template upload in the ledger: %v", err)
    }

    fmt.Printf("Template %s uploaded by %s.\n", templateID, creator)
    return templateID, nil
}

// StartEscrow initiates an escrow process for purchasing a template
func (mp *common.SmartContractTemplateMarketplace) StartEscrow(buyer, templateID string) (string, error) {
    mp.mutex.Lock()
    defer mp.mutex.Unlock()

    template, exists := mp.Templates[templateID]
    if !exists {
        return "", errors.New("template not found")
    }

    // Generate a new escrow ID
    escrowID := mp.generateEscrowID()

    escrow := &common.Escrow{
        ID:        escrowID,
        Buyer:     buyer,
        Seller:    template.Creator,
        Amount:    template.Price,
        CreatedAt: time.Now(),
        IsComplete: false,
        Disputed:   false,
    }

    mp.Escrows[escrowID] = escrow

    // Record the escrow creation in the ledger
    err := mp.LedgerInstance.RecordEscrow(escrowID, buyer, template.Creator, template.Price)
    if err != nil {
        return "", fmt.Errorf("failed to record escrow in the ledger: %v", err)
    }

    fmt.Printf("Escrow %s initiated by %s for template %s.\n", escrowID, buyer, templateID)
    return escrowID, nil
}

// CompleteEscrow finalizes the escrow transaction and releases the funds to the seller
func (mp *common.SmartContractTemplateMarketplace) CompleteEscrow(escrowID string) error {
    mp.mutex.Lock()
    defer mp.mutex.Unlock()

    escrow, exists := mp.Escrows[escrowID]
    if !exists {
        return errors.New("escrow not found")
    }

    if escrow.IsComplete {
        return errors.New("escrow is already complete")
    }

    // Apply the marketplace fee
    marketplaceFee := escrow.Amount * mp.EscrowFee / 100
    finalAmount := escrow.Amount - marketplaceFee

    // Record the transfer of funds in the ledger
    err := mp.LedgerInstance.RecordTransaction(escrow.Buyer, escrow.Seller, finalAmount)
    if err != nil {
        return fmt.Errorf("failed to complete escrow transaction: %v", err)
    }

    escrow.IsComplete = true
    fmt.Printf("Escrow %s completed, %.2f SYNN transferred to %s.\n", escrowID, finalAmount, escrow.Seller)
    return nil
}

// PurchaseTemplate processes the template purchase through escrow
func (mp *common.SmartContractTemplateMarketplace) PurchaseTemplate(buyer, templateID string) (string, error) {
    escrowID, err := mp.StartEscrow(buyer, templateID)
    if err != nil {
        return "", err
    }

    // Complete escrow after creation (assuming no disputes in this example)
    err = mp.CompleteEscrow(escrowID)
    if err != nil {
        return "", err
    }

    // Return the encrypted template code to the buyer
    template := mp.Templates[templateID]
    return template.EncryptedCode, nil
}

// QueryTemplates allows users to search for available templates in the marketplace
func (mp *common.SmartContractTemplateMarketplace) QueryTemplates(filter string, maxPrice float64) []*common.SmartContractTemplate {
    mp.mutex.Lock()
    defer mp.mutex.Unlock()

    templates := make([]*common.SmartContractTemplate, 0)
    for _, template := range mp.Templates {
        if filter == "" || contains(template.Name, filter) || contains(template.Description, filter) {
            if maxPrice == 0 || template.Price <= maxPrice {
                templates = append(templates, template)
            }
        }
    }
    return templates
}

// contains checks if a string contains a substring, case-insensitive
func contains(str, substr string) bool {
    return strings.Contains(strings.ToLower(str), strings.ToLower(substr))
}

// generateTemplateID creates a unique ID for a template based on the creator and name
func (mp *common.SmartContractTemplateMarketplace) generateTemplateID(creator, name string) string {
    hashInput := fmt.Sprintf("%s%s%d", creator, name, time.Now().UnixNano())
    hash := sha256.New()
    hash.Write([]byte(hashInput))
    return hex.EncodeToString(hash.Sum(nil))
}

// generateEscrowID creates a random escrow ID
func (mp *common.SmartContractTemplateMarketplace) generateEscrowID() string {
    return fmt.Sprintf("escrow-%d", rand.Intn(1000000000))
}

// Sidechain_Basic_Operations.go

package sidechains

import (
    "fmt"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
    "synnergy_network/pkg/encryption"
)

// SidechainInit initializes a new sidechain instance.
func SidechainInit(chainID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.InitSidechain(chainID); err != nil {
        return fmt.Errorf("failed to initialize sidechain %s: %v", chainID, err)
    }
    fmt.Printf("Sidechain %s initialized.\n", chainID)
    return nil
}

// SidechainCreate creates a new sidechain in the network.
func SidechainCreate(name, owner string, ledgerInstance *ledger.Ledger) error {
    encryptedName := encryption.EncryptData(name)
    if err := ledgerInstance.CreateSidechain(encryptedName, owner); err != nil {
        return fmt.Errorf("failed to create sidechain %s: %v", name, err)
    }
    fmt.Printf("Sidechain %s created.\n", name)
    return nil
}

// SidechainRegister registers an existing sidechain.
func SidechainRegister(chainID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.RegisterSidechain(chainID); err != nil {
        return fmt.Errorf("failed to register sidechain %s: %v", chainID, err)
    }
    fmt.Printf("Sidechain %s registered.\n", chainID)
    return nil
}

// SidechainDeregister removes a sidechain from the network.
func SidechainDeregister(chainID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.DeregisterSidechain(chainID); err != nil {
        return fmt.Errorf("failed to deregister sidechain %s: %v", chainID, err)
    }
    fmt.Printf("Sidechain %s deregistered.\n", chainID)
    return nil
}

// SidechainConnect connects the sidechain to the main network.
func SidechainConnect(chainID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.ConnectSidechain(chainID); err != nil {
        return fmt.Errorf("failed to connect sidechain %s: %v", chainID, err)
    }
    fmt.Printf("Sidechain %s connected to main network.\n", chainID)
    return nil
}

// SidechainDisconnect disconnects the sidechain from the main network.
func SidechainDisconnect(chainID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.DisconnectSidechain(chainID); err != nil {
        return fmt.Errorf("failed to disconnect sidechain %s: %v", chainID, err)
    }
    fmt.Printf("Sidechain %s disconnected from main network.\n", chainID)
    return nil
}

// SidechainBridgeToken bridges a token to the sidechain.
func SidechainBridgeToken(tokenID, chainID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.BridgeToken(tokenID, chainID); err != nil {
        return fmt.Errorf("failed to bridge token %s to sidechain %s: %v", tokenID, chainID, err)
    }
    fmt.Printf("Token %s bridged to sidechain %s.\n", tokenID, chainID)
    return nil
}

// SidechainUnbridgeToken unbridges a token from the sidechain.
func SidechainUnbridgeToken(tokenID, chainID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.UnbridgeToken(tokenID, chainID); err != nil {
        return fmt.Errorf("failed to unbridge token %s from sidechain %s: %v", tokenID, chainID, err)
    }
    fmt.Printf("Token %s unbridged from sidechain %s.\n", tokenID, chainID)
    return nil
}

// SidechainMintToken mints a new token on the sidechain.
func SidechainMintToken(tokenID string, amount int, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.MintToken(tokenID, amount); err != nil {
        return fmt.Errorf("failed to mint token %s with amount %d: %v", tokenID, amount, err)
    }
    fmt.Printf("Minted %d of token %s on sidechain.\n", amount, tokenID)
    return nil
}

// SidechainBurnToken burns a token on the sidechain.
func SidechainBurnToken(tokenID string, amount int, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.BurnToken(tokenID, amount); err != nil {
        return fmt.Errorf("failed to burn token %s with amount %d: %v", tokenID, amount, err)
    }
    fmt.Printf("Burned %d of token %s on sidechain.\n", amount, tokenID)
    return nil
}

// SidechainSyncToMainchain syncs sidechain state to the main chain.
func SidechainSyncToMainchain(chainID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.SyncToMainchain(chainID); err != nil {
        return fmt.Errorf("failed to sync sidechain %s to main chain: %v", chainID, err)
    }
    fmt.Printf("Sidechain %s synced to main chain.\n", chainID)
    return nil
}

// SidechainSyncFromMainchain syncs main chain state to the sidechain.
func SidechainSyncFromMainchain(chainID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.SyncFromMainchain(chainID); err != nil {
        return fmt.Errorf("failed to sync from main chain to sidechain %s: %v", chainID, err)
    }
    fmt.Printf("Main chain state synced to sidechain %s.\n", chainID)
    return nil
}

// SidechainUpdateMainchain updates sidechain data on the mainchain.
func SidechainUpdateMainchain(chainID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.UpdateMainchain(chainID); err != nil {
        return fmt.Errorf("failed to update main chain from sidechain %s: %v", chainID, err)
    }
    fmt.Printf("Sidechain %s state updated on main chain.\n", chainID)
    return nil
}

// SidechainUpdateSidechain updates mainchain data on the sidechain.
func SidechainUpdateSidechain(chainID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.UpdateSidechain(chainID); err != nil {
        return fmt.Errorf("failed to update sidechain %s from main chain: %v", chainID, err)
    }
    fmt.Printf("Main chain state updated on sidechain %s.\n", chainID)
    return nil
}

// SidechainTransferAsset transfers an asset on the sidechain.
func SidechainTransferAsset(assetID string, from, to string, amount int, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.TransferAsset(assetID, from, to, amount); err != nil {
        return fmt.Errorf("failed to transfer asset %s from %s to %s: %v", assetID, from, to, err)
    }
    fmt.Printf("Transferred %d of asset %s from %s to %s on sidechain.\n", amount, assetID, from, to)
    return nil
}

// SidechainFreezeAsset freezes an asset on the sidechain.
func SidechainFreezeAsset(assetID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.FreezeAsset(assetID); err != nil {
        return fmt.Errorf("failed to freeze asset %s: %v", assetID, err)
    }
    fmt.Printf("Asset %s frozen on sidechain.\n", assetID)
    return nil
}

// SidechainUnfreezeAsset unfreezes an asset on the sidechain.
func SidechainUnfreezeAsset(assetID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.UnfreezeAsset(assetID); err != nil {
        return fmt.Errorf("failed to unfreeze asset %s: %v", assetID, err)
    }
    fmt.Printf("Asset %s unfrozen on sidechain.\n", assetID)
    return nil
}

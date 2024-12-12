package ledger

import (
	"errors"
	"fmt"
	"time"
)

// RecordRoleAssignment records the assignment of a role to a DAO member.
func (l *DAOLedger) RecordRoleAssignment(daoID, memberID, role string) error {
	l.Lock()
	defer l.Unlock()

	if _, exists := l.DAORecords[daoID]; !exists {
		return errors.New("DAO does not exist")
	}

	l.DAORecords[daoID].RoleAssignments[memberID] = role
	fmt.Printf("Role %s assigned to member %s in DAO %s\n", role, memberID, daoID)
	return nil
}

// RecordRoleRevocation records the revocation of a role from a DAO member.
func (l *DAOLedger) RecordRoleRevocation(daoID, memberID string) error {
	l.Lock()
	defer l.Unlock()

	if _, exists := l.DAORecords[daoID]; !exists {
		return errors.New("DAO does not exist")
	}

	delete(l.DAORecords[daoID].RoleAssignments, memberID)
	fmt.Printf("Role revoked from member %s in DAO %s\n", memberID, daoID)
	return nil
}

// GetTransactionByID retrieves a transaction by its ID within the DAO.
func (l *DAOLedger) GetDAOTransactionByID(daoID, txnID string) (TransactionRecord, error) {
	l.Lock()
	defer l.Unlock()

	if dao, exists := l.DAORecords[daoID]; exists {
		if txn, ok := dao.Transactions[txnID]; ok {
			return txn, nil
		}
		return TransactionRecord{}, errors.New("Transaction not found")
	}
	return TransactionRecord{}, errors.New("DAO does not exist")
}

// UpdateEncryptedRole updates an encrypted role in the DAO ledger.
func (l *DAOLedger) UpdateEncryptedRole(daoID, memberID, encryptedRole string) error {
	l.Lock()
	defer l.Unlock()

	if dao, exists := l.DAORecords[daoID]; exists {
		dao.RoleAssignments[memberID] = encryptedRole
		fmt.Printf("Encrypted role updated for member %s in DAO %s\n", memberID, daoID)
		return nil
	}
	return errors.New("DAO does not exist")
}

// RecordDAOCreation records the creation of a DAO in the ledger.
func (l *DAOLedger) RecordDAOCreation(daoID, creatorID string) error {
	l.Lock()
	defer l.Unlock()

	if _, exists := l.DAORecords[daoID]; exists {
		return errors.New("DAO already exists")
	}

	l.DAORecords[daoID] = &DAORecord{
		ID:               daoID,
		Members:          make(map[string]DAOMember),
		Proposals:        make(map[string]DAOProposal),
		Transactions:     make(map[string]TransactionRecord),
		GovernanceStakes: make(map[string]float64),
		RoleAssignments:  make(map[string]string),
	}
	fmt.Printf("DAO %s created by %s\n", daoID, creatorID)
	return nil
}

// RecordMemberAddition records a new member's addition to the DAO.
func (l *DAOLedger) RecordMemberAddition(daoID, memberID string, votingPower, stakeBalance float64) error {
	l.Lock()
	defer l.Unlock()

	if dao, exists := l.DAORecords[daoID]; exists {
		dao.Members[memberID] = DAOMember{
			WalletAddress: memberID,
			VotingPower:   int(votingPower), // Convert float64 to int
			Role:          "Member",         // Default role, can be updated
			IsAuthorized:  true,             // Default to authorized
		}
		fmt.Printf("Member %s added to DAO %s\n", memberID, daoID)
		return nil
	}
	return errors.New("DAO does not exist")
}


// RecordProposalCreation records the creation of a proposal in the DAO.
func (l *DAOLedger) RecordProposalCreation(daoID, proposalID, creatorID, content string) error {
	l.Lock()
	defer l.Unlock()

	if dao, exists := l.DAORecords[daoID]; exists {
		dao.Proposals[proposalID] = DAOProposal{
			ProposalID:   proposalID,
			Title:        "New Proposal",
			Description:  content,
			Author:       creatorID,
			CreationTime: time.Now(),
			Status:       "Pending",
		}
		fmt.Printf("Proposal %s created by %s in DAO %s\n", proposalID, creatorID, daoID)
		return nil
	}
	return errors.New("DAO does not exist")
}

// RecordProposalVote records a vote on a proposal.
func (l *DAOLedger) RecordProposalVote(daoID, proposalID, memberID string, vote bool) error {
	l.Lock()
	defer l.Unlock()

	if dao, exists := l.DAORecords[daoID]; exists {
		if proposal, ok := dao.Proposals[proposalID]; ok {
			if vote {
				proposal.ApproveCount++
			} else {
				proposal.RejectCount++
			}
			proposal.VoteCount++
			fmt.Printf("Member %s voted on Proposal %s in DAO %s\n", memberID, proposalID, daoID)
			return nil
		}
		return errors.New("Proposal not found")
	}
	return errors.New("DAO does not exist")
}

// RecordProposalFinalization finalizes the result of a proposal in the DAO.
func (l *DAOLedger) RecordProposalFinalization(daoID, proposalID, result string) error {
	l.Lock()
	defer l.Unlock()

	if dao, exists := l.DAORecords[daoID]; exists {
		if proposal, ok := dao.Proposals[proposalID]; ok {
			proposal.Status = "Finalized"
			proposal.Status = result
			fmt.Printf("Proposal %s finalized in DAO %s with result: %s\n", proposalID, daoID, result)
			return nil
		}
		return errors.New("Proposal not found")
	}
	return errors.New("DAO does not exist")
}

// RecordMemberRemoval records the removal of a member from the DAO.
func (l *DAOLedger) RecordMemberRemoval(daoID, memberID string) error {
	l.Lock()
	defer l.Unlock()

	if dao, exists := l.DAORecords[daoID]; exists {
		delete(dao.Members, memberID)
		fmt.Printf("Member %s removed from DAO %s\n", memberID, daoID)
		return nil
	}
	return errors.New("DAO does not exist")
}

// RecordVaultTransaction records a transaction within the DAO vault.
func (l *DAOLedger) RecordVaultTransaction(daoID, txnID string, txn TransactionRecord) error {
	l.Lock()
	defer l.Unlock()

	if dao, exists := l.DAORecords[daoID]; exists {
		dao.Transactions[txnID] = txn
		fmt.Printf("Transaction %s recorded in DAO %s vault\n", txnID, daoID)
		return nil
	}
	return errors.New("DAO does not exist")
}

// RecordTransactionExecution records the execution of a transaction.
func (l *DAOLedger) RecordDAOTransactionExecution(daoID, txnID string) error {
	l.Lock()
	defer l.Unlock()

	if dao, exists := l.DAORecords[daoID]; exists {
		if txn, ok := dao.Transactions[txnID]; ok {
			txn.Status = "Executed"
			dao.Transactions[txnID] = txn
			fmt.Printf("Transaction %s executed in DAO %s\n", txnID, daoID)
			return nil
		}
		return errors.New("Transaction not found")
	}
	return errors.New("DAO does not exist")
}

// RecordTransactionApproval records the approval of a DAO transaction.
func (l *DAOLedger) RecordTransactionApproval(daoID, txnID string) error {
	l.Lock()
	defer l.Unlock()

	if dao, exists := l.DAORecords[daoID]; exists {
		if txn, ok := dao.Transactions[txnID]; ok {
			txn.Status = "Approved"
			dao.Transactions[txnID] = txn
			fmt.Printf("Transaction %s approved in DAO %s\n", txnID, daoID)
			return nil
		}
		return errors.New("Transaction not found")
	}
	return errors.New("DAO does not exist")
}

// RecordEmergencyAccessRequest records an emergency access request in the DAO.
func (l *DAOLedger) RecordEmergencyAccessRequest(daoID, requestID string) error {
	l.Lock()
	defer l.Unlock()

	fmt.Printf("Emergency access request %s recorded for DAO %s\n", requestID, daoID)
	return nil
}

// RecordEmergencyAccessApproval records the approval of an emergency access request.
func (l *DAOLedger) RecordEmergencyAccessApproval(daoID, requestID string) error {
	l.Lock()
	defer l.Unlock()

	fmt.Printf("Emergency access request %s approved for DAO %s\n", requestID, daoID)
	return nil
}

// RecordTransactionRejection records the rejection of a transaction.
func (l *DAOLedger) RecordTransactionRejection(daoID, txnID string) error {
	l.Lock()
	defer l.Unlock()

	if dao, exists := l.DAORecords[daoID]; exists {
		if txn, ok := dao.Transactions[txnID]; ok {
			txn.Status = "Rejected"
			dao.Transactions[txnID] = txn
			fmt.Printf("Transaction %s rejected in DAO %s\n", txnID, daoID)
			return nil
		}
		return errors.New("Transaction not found")
	}
	return errors.New("DAO does not exist")
}

func (l *DAOLedger) RecordVotingPowerChange(daoID, memberID string, newVotingPower float64) error {
	l.Lock()
	defer l.Unlock()

	if dao, exists := l.DAORecords[daoID]; exists {
		if member, ok := dao.Members[memberID]; ok {
			member.VotingPower = int(newVotingPower) // Convert float64 to int
			dao.Members[memberID] = member
			fmt.Printf("Voting power for member %s changed to %d in DAO %s\n", memberID, int(newVotingPower), daoID)
			return nil
		}
		return errors.New("Member not found")
	}
	return errors.New("DAO does not exist")
}


// RecordGovernanceStakingInitialization initializes staking for governance in a DAO.
func (l *DAOLedger) RecordGovernanceStakingInitialization(daoID, memberID string, stakeAmount float64) error {
	l.Lock()
	defer l.Unlock()

	if dao, exists := l.DAORecords[daoID]; exists {
		dao.GovernanceStakes[memberID] = stakeAmount
		fmt.Printf("Governance staking initialized for member %s in DAO %s with %.2f\n", memberID, daoID, stakeAmount)
		return nil
	}
	return errors.New("DAO does not exist")
}

// RecordStakeTransaction records a stake transaction within the DAO.
func (l *DAOLedger) RecordStakeTransaction(daoID, memberID string, stakeAmount float64) error {
	return l.RecordGovernanceStakingInitialization(daoID, memberID, stakeAmount)
}

// RecordUnstakeTransaction records an unstake transaction within the DAO.
func (l *DAOLedger) RecordUnstakeTransaction(daoID, memberID string, unstakeAmount float64) error {
	l.Lock()
	defer l.Unlock()

	if dao, exists := l.DAORecords[daoID]; exists {
		if stake, ok := dao.GovernanceStakes[memberID]; ok {
			dao.GovernanceStakes[memberID] = stake - unstakeAmount
			fmt.Printf("Unstake of %.2f recorded for member %s in DAO %s\n", unstakeAmount, memberID, daoID)
			return nil
		}
		return errors.New("Stake not found for member")
	}
	return errors.New("DAO does not exist")
}

// RecordProposalResult records the result of a proposal vote.
func (l *DAOLedger) RecordProposalResult(daoID, proposalID, result string) error {
	return l.RecordProposalFinalization(daoID, proposalID, result)
}

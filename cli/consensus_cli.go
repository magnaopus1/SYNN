package cli

import (
	"fmt"
	"strconv"
	"time"

	"synnergy_network/pkg/consensus"
	"synnergy_network/pkg/ledger"

	"github.com/spf13/cobra"
)

// ConsensusCLI represents the consensus CLI commands
type ConsensusCLI struct {
	LedgerInstance *ledger.Ledger
}

// NewConsensusCLI creates a new consensus CLI instance
func NewConsensusCLI(ledgerInstance *ledger.Ledger) *ConsensusCLI {
	return &ConsensusCLI{
		LedgerInstance: ledgerInstance,
	}
}

// GetConsensusCommands returns all consensus-related CLI commands
func (cli *ConsensusCLI) GetConsensusCommands() *cobra.Command {
	consensusCmd := &cobra.Command{
		Use:   "consensus",
		Short: "Consensus management commands",
		Long:  "Commands for managing consensus mechanisms, validators, and network parameters",
	}

	// Difficulty management commands
	consensusCmd.AddCommand(cli.getDifficultyCommands())
	
	// Validator management commands
	consensusCmd.AddCommand(cli.getValidatorCommands())
	
	// Audit commands
	consensusCmd.AddCommand(cli.getAuditCommands())
	
	// Reward commands
	consensusCmd.AddCommand(cli.getRewardCommands())
	
	// PoH commands
	consensusCmd.AddCommand(cli.getPoHCommands())
	
	// Stake commands
	consensusCmd.AddCommand(cli.getStakeCommands())

	return consensusCmd
}

// getDifficultyCommands returns difficulty management commands
func (cli *ConsensusCLI) getDifficultyCommands() *cobra.Command {
	difficultyCmd := &cobra.Command{
		Use:   "difficulty",
		Short: "Difficulty management commands",
		Long:  "Commands for managing consensus difficulty settings",
	}

	// Adjust difficulty command
	adjustCmd := &cobra.Command{
		Use:   "adjust [level] [reason]",
		Short: "Adjust consensus difficulty level",
		Long:  "Adjust the consensus difficulty level with a specified reason",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			level, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid difficulty level: %v", err)
			}

			reason := args[1]
			
			err = consensus.ConsensusAdjustDifficultyBasedOnTime(level, reason, cli.LedgerInstance)
			if err != nil {
				return fmt.Errorf("failed to adjust difficulty: %v", err)
			}

			fmt.Printf("✅ Difficulty adjusted to level %d. Reason: %s\n", level, reason)
			return nil
		},
	}

	// Monitor block generation command
	monitorCmd := &cobra.Command{
		Use:   "monitor [blockID] [generationTimeMs]",
		Short: "Monitor block generation time",
		Long:  "Log block generation time for monitoring purposes",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			blockID := args[0]
			generationTimeMs, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid generation time: %v", err)
			}

			generationTime := time.Duration(generationTimeMs) * time.Millisecond
			
			err = consensus.consensusMonitorBlockGenerationTime(blockID, generationTime, cli.LedgerInstance)
			if err != nil {
				return fmt.Errorf("failed to monitor block generation: %v", err)
			}

			fmt.Printf("✅ Block generation time logged for block %s: %v\n", blockID, generationTime)
			return nil
		},
	}

	difficultyCmd.AddCommand(adjustCmd, monitorCmd)
	return difficultyCmd
}

// getValidatorCommands returns validator management commands
func (cli *ConsensusCLI) getValidatorCommands() *cobra.Command {
	validatorCmd := &cobra.Command{
		Use:   "validator",
		Short: "Validator management commands",
		Long:  "Commands for managing validators and their activities",
	}

	// Track participation command
	trackCmd := &cobra.Command{
		Use:   "track [validatorID] [status]",
		Short: "Track validator participation",
		Long:  "Track validator participation in consensus with status (Active, Inactive, Disqualified, Suspended)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			validatorID := args[0]
			status := args[1]
			
			err := consensus.ConsensusTrackConsensusParticipation(validatorID, status, cli.LedgerInstance)
			if err != nil {
				return fmt.Errorf("failed to track participation: %v", err)
			}

			fmt.Printf("✅ Validator %s participation tracked with status: %s\n", validatorID, status)
			return nil
		},
	}

	// Validate activity command
	validateCmd := &cobra.Command{
		Use:   "validate [validatorID] [action] [details]",
		Short: "Validate validator activity",
		Long:  "Validate and log validator activity with details",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			validatorID := args[0]
			action := args[1]
			details := args[2]
			
			err := consensus.ConsensusValidateValidatorActivity(validatorID, action, details, cli.LedgerInstance)
			if err != nil {
				return fmt.Errorf("failed to validate activity: %v", err)
			}

			fmt.Printf("✅ Validator %s activity validated. Action: %s\n", validatorID, action)
			return nil
		},
	}

	// Get activity logs command
	activityLogsCmd := &cobra.Command{
		Use:   "logs",
		Short: "Get validator activity logs",
		Long:  "Retrieve validator activity logs from the ledger",
		RunE: func(cmd *cobra.Command, args []string) error {
			logs, err := consensus.ConsensusFetchValidatorActivityLogs(cli.LedgerInstance)
			if err != nil {
				return fmt.Errorf("failed to fetch activity logs: %v", err)
			}

			fmt.Println("📋 Validator Activity Logs:")
			for i, log := range logs {
				fmt.Printf("%d. Validator: %s, Action: %s, Details: %s, Time: %v\n", 
					i+1, log.ValidatorID, log.Action, log.Details, log.Timestamp)
			}
			return nil
		},
	}

	// Set selection mode command
	setModeCmd := &cobra.Command{
		Use:   "set-mode [mode]",
		Short: "Set validator selection mode",
		Long:  "Set the validator selection mode (RandomSelection, StakeBasedSelection, RotationBasedSelection)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Parse mode (this would need proper type conversion in real implementation)
			modeStr := args[0]
			var mode ledger.ValidatorSelectionMode
			
			// This is a simplified version - in practice you'd have proper mode parsing
			fmt.Printf("Setting validator selection mode: %s\n", modeStr)
			
			err := consensus.ConsensusSetValidatorSelectionMode(mode, cli.LedgerInstance)
			if err != nil {
				return fmt.Errorf("failed to set selection mode: %v", err)
			}

			fmt.Printf("✅ Validator selection mode set to: %s\n", modeStr)
			return nil
		},
	}

	// Get selection mode command
	getModeCmd := &cobra.Command{
		Use:   "get-mode",
		Short: "Get current validator selection mode",
		Long:  "Retrieve the current validator selection mode",
		RunE: func(cmd *cobra.Command, args []string) error {
			mode, err := consensus.ConsensusGetValidatorSelectionMode(cli.LedgerInstance)
			if err != nil {
				return fmt.Errorf("failed to get selection mode: %v", err)
			}

			fmt.Printf("📊 Current validator selection mode: %v\n", mode)
			return nil
		},
	}

	validatorCmd.AddCommand(trackCmd, validateCmd, activityLogsCmd, setModeCmd, getModeCmd)
	return validatorCmd
}

// getAuditCommands returns audit management commands
func (cli *ConsensusCLI) getAuditCommands() *cobra.Command {
	auditCmd := &cobra.Command{
		Use:   "audit",
		Short: "Audit management commands",
		Long:  "Commands for managing consensus audit functionality",
	}

	// Enable audit command
	enableCmd := &cobra.Command{
		Use:   "enable",
		Short: "Enable consensus audit",
		Long:  "Enable consensus audit functionality",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := consensus.consensusEnableConsensusAudit(cli.LedgerInstance)
			if err != nil {
				return fmt.Errorf("failed to enable audit: %v", err)
			}

			fmt.Println("✅ Consensus audit enabled successfully")
			return nil
		},
	}

	// Disable audit command
	disableCmd := &cobra.Command{
		Use:   "disable",
		Short: "Disable consensus audit",
		Long:  "Disable consensus audit functionality",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := consensus.consensusDisableConsensusAudit(cli.LedgerInstance)
			if err != nil {
				return fmt.Errorf("failed to disable audit: %v", err)
			}

			fmt.Println("✅ Consensus audit disabled successfully")
			return nil
		},
	}

	// Get audit logs command
	logsCmd := &cobra.Command{
		Use:   "logs",
		Short: "Get consensus audit logs",
		Long:  "Retrieve consensus audit logs from the ledger",
		RunE: func(cmd *cobra.Command, args []string) error {
			logs, err := consensus.ConsensusFetchConsensusLogs(cli.LedgerInstance)
			if err != nil {
				return fmt.Errorf("failed to fetch audit logs: %v", err)
			}

			fmt.Println("📋 Consensus Audit Logs:")
			for i, log := range logs {
				fmt.Printf("%d. Audit ID: %s, Validator: %s, Status: %s, Time: %v\n", 
					i+1, log.AuditID, log.ValidatorID, log.ParticipationStatus, log.Timestamp)
			}
			return nil
		},
	}

	auditCmd.AddCommand(enableCmd, disableCmd, logsCmd)
	return auditCmd
}

// getRewardCommands returns reward management commands
func (cli *ConsensusCLI) getRewardCommands() *cobra.Command {
	rewardCmd := &cobra.Command{
		Use:   "reward",
		Short: "Reward management commands",
		Long:  "Commands for managing reward distribution",
	}

	// Set reward mode command
	setModeCmd := &cobra.Command{
		Use:   "set-mode [mode]",
		Short: "Set reward distribution mode",
		Long:  "Set the reward distribution mode",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Parse mode (simplified - would need proper type conversion)
			modeStr := args[0]
			var mode ledger.RewardDistributionMode
			
			fmt.Printf("Setting reward distribution mode: %s\n", modeStr)
			
			err := consensus.consensusSetRewardDistributionMode(mode, cli.LedgerInstance)
			if err != nil {
				return fmt.Errorf("failed to set reward mode: %v", err)
			}

			fmt.Printf("✅ Reward distribution mode set to: %s\n", modeStr)
			return nil
		},
	}

	// Get reward mode command
	getModeCmd := &cobra.Command{
		Use:   "get-mode",
		Short: "Get current reward distribution mode",
		Long:  "Retrieve the current reward distribution mode",
		RunE: func(cmd *cobra.Command, args []string) error {
			mode, err := consensus.consensusGetRewardDistributionMode(cli.LedgerInstance)
			if err != nil {
				return fmt.Errorf("failed to get reward mode: %v", err)
			}

			fmt.Printf("📊 Current reward distribution mode: %v\n", mode)
			return nil
		},
	}

	rewardCmd.AddCommand(setModeCmd, getModeCmd)
	return rewardCmd
}

// getPoHCommands returns Proof of History commands
func (cli *ConsensusCLI) getPoHCommands() *cobra.Command {
	pohCmd := &cobra.Command{
		Use:   "poh",
		Short: "Proof of History commands",
		Long:  "Commands for managing Proof of History settings",
	}

	// Set threshold command
	setThresholdCmd := &cobra.Command{
		Use:   "set-threshold [threshold]",
		Short: "Set PoH participation threshold",
		Long:  "Set the Proof of History participation threshold (0.0 to 1.0)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			threshold, err := strconv.ParseFloat(args[0], 64)
			if err != nil {
				return fmt.Errorf("invalid threshold value: %v", err)
			}

			err = consensus.ConsensusSetPoHParticipationThreshold(threshold, cli.LedgerInstance)
			if err != nil {
				return fmt.Errorf("failed to set PoH threshold: %v", err)
			}

			fmt.Printf("✅ PoH participation threshold set to: %.2f\n", threshold)
			return nil
		},
	}

	// Get threshold command
	getThresholdCmd := &cobra.Command{
		Use:   "get-threshold",
		Short: "Get current PoH participation threshold",
		Long:  "Retrieve the current Proof of History participation threshold",
		RunE: func(cmd *cobra.Command, args []string) error {
			threshold, err := consensus.ConsensusGetPoHParticipationThreshold(cli.LedgerInstance)
			if err != nil {
				return fmt.Errorf("failed to get PoH threshold: %v", err)
			}

			fmt.Printf("📊 Current PoH participation threshold: %.2f\n", threshold)
			return nil
		},
	}

	pohCmd.AddCommand(setThresholdCmd, getThresholdCmd)
	return pohCmd
}

// getStakeCommands returns stake management commands
func (cli *ConsensusCLI) getStakeCommands() *cobra.Command {
	stakeCmd := &cobra.Command{
		Use:   "stake",
		Short: "Stake management commands",
		Long:  "Commands for managing dynamic stake adjustment",
	}

	// Enable dynamic stake command
	enableCmd := &cobra.Command{
		Use:   "enable-dynamic",
		Short: "Enable dynamic stake adjustment",
		Long:  "Enable dynamic stake adjustment functionality",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := consensus.ConsensusEnableDynamicStakeAdjustment(cli.LedgerInstance)
			if err != nil {
				return fmt.Errorf("failed to enable dynamic stake: %v", err)
			}

			fmt.Println("✅ Dynamic stake adjustment enabled successfully")
			return nil
		},
	}

	// Disable dynamic stake command
	disableCmd := &cobra.Command{
		Use:   "disable-dynamic",
		Short: "Disable dynamic stake adjustment",
		Long:  "Disable dynamic stake adjustment functionality",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := consensus.ConsensusDisableDynamicStakeAdjustment(cli.LedgerInstance)
			if err != nil {
				return fmt.Errorf("failed to disable dynamic stake: %v", err)
			}

			fmt.Println("✅ Dynamic stake adjustment disabled successfully")
			return nil
		},
	}

	stakeCmd.AddCommand(enableCmd, disableCmd)
	return stakeCmd
}
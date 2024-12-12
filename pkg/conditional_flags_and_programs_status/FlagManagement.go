package conditional_flags_and_programs_status

import (
	"fmt"
	"synnergy_network/pkg/ledger"
)

// setZeroFlag sets the zero flag in the ledger to indicate zero value in the last operation.
func SetZeroFlag() error {
	l := &ledger.Ledger{}
	if err := l.ConditionalFlagsLedger.SetFlag(ZeroFlag, ZeroFlag, true); err != nil {
		return fmt.Errorf("failed to set zero flag: %v", err)
	}
	return nil
}

// clearZeroFlag clears the zero flag in the ledger.
func ClearZeroFlag() error {
	l := &ledger.Ledger{}
	if err := l.ConditionalFlagsLedger.SetFlag(ZeroFlag, ZeroFlag, false); err != nil {
		return fmt.Errorf("failed to clear zero flag: %v", err)
	}
	return nil
}

// checkZeroFlag checks the status of the zero flag.
func CheckZeroFlag() (bool, error) {
	l := &ledger.Ledger{}
	return l.ConditionalFlagsLedger.CheckFlag(ZeroFlag, ZeroFlag)
}

// setCarryFlag sets the carry flag.
func SetCarryFlag() error {
	l := &ledger.Ledger{}
	if err := l.ConditionalFlagsLedger.SetFlag(CarryFlag, CarryFlag, true); err != nil {
		return fmt.Errorf("failed to set carry flag: %v", err)
	}
	return nil
}

// clearCarryFlag clears the carry flag.
func ClearCarryFlag() error {
	l := &ledger.Ledger{}
	if err := l.ConditionalFlagsLedger.SetFlag(CarryFlag, CarryFlag, false); err != nil {
		return fmt.Errorf("failed to clear carry flag: %v", err)
	}
	return nil
}

// checkCarryFlag checks the status of the carry flag.
func CheckCarryFlag() (bool, error) {
	l := &ledger.Ledger{}
	return l.ConditionalFlagsLedger.CheckFlag(CarryFlag, CarryFlag)
}

// setOverflowFlag sets the overflow flag.
func SetOverflowFlag() error {
	l := &ledger.Ledger{}
	if err := l.ConditionalFlagsLedger.SetFlag(OverflowFlag, OverflowFlag, true); err != nil {
		return fmt.Errorf("failed to set overflow flag: %v", err)
	}
	return nil
}

// clearOverflowFlag clears the overflow flag.
func ClearOverflowFlag() error {
	l := &ledger.Ledger{}
	if err := l.ConditionalFlagsLedger.SetFlag(OverflowFlag, OverflowFlag, false); err != nil {
		return fmt.Errorf("failed to clear overflow flag: %v", err)
	}
	return nil
}

// checkOverflowFlag checks the status of the overflow flag.
func CheckOverflowFlag() (bool, error) {
	l := &ledger.Ledger{}
	return l.ConditionalFlagsLedger.CheckFlag(OverflowFlag, OverflowFlag)
}

// setSignFlag sets the sign flag.
func SetSignFlag() error {
	l := &ledger.Ledger{}
	if err := l.ConditionalFlagsLedger.SetFlag(SignFlag, SignFlag, true); err != nil {
		return fmt.Errorf("failed to set sign flag: %v", err)
	}
	return nil
}

// clearSignFlag clears the sign flag.
func ClearSignFlag() error {
	l := &ledger.Ledger{}
	if err := l.ConditionalFlagsLedger.SetFlag(SignFlag, SignFlag, false); err != nil {
		return fmt.Errorf("failed to clear sign flag: %v", err)
	}
	return nil
}

// checkSignFlag checks the status of the sign flag.
func CheckSignFlag() (bool, error) {
	l := &ledger.Ledger{}
	return l.ConditionalFlagsLedger.CheckFlag(SignFlag, SignFlag)
}

// setInterruptFlag sets the interrupt flag.
func SetInterruptFlag() error {
	l := &ledger.Ledger{}
	if err := l.ConditionalFlagsLedger.SetFlag(InterruptFlag, InterruptFlag, true); err != nil {
		return fmt.Errorf("failed to set interrupt flag: %v", err)
	}
	return nil
}

// clearInterruptFlag clears the interrupt flag.
func ClearInterruptFlag() error {
	l := &ledger.Ledger{}
	if err := l.ConditionalFlagsLedger.SetFlag(InterruptFlag, InterruptFlag, false); err != nil {
		return fmt.Errorf("failed to clear interrupt flag: %v", err)
	}
	return nil
}

// checkInterruptFlag checks the status of the interrupt flag.
func CheckInterruptFlag() (bool, error) {
	l := &ledger.Ledger{}
	return l.ConditionalFlagsLedger.CheckFlag(InterruptFlag, InterruptFlag)
}

// setParityFlag sets the parity flag.
func SetParityFlag() error {
	l := &ledger.Ledger{}
	if err := l.ConditionalFlagsLedger.SetFlag(ParityFlag, ParityFlag, true); err != nil {
		return fmt.Errorf("failed to set parity flag: %v", err)
	}
	return nil
}

// clearParityFlag clears the parity flag.
func ClearParityFlag() error {
	l := &ledger.Ledger{}
	if err := l.ConditionalFlagsLedger.SetFlag(ParityFlag, ParityFlag, false); err != nil {
		return fmt.Errorf("failed to clear parity flag: %v", err)
	}
	return nil
}

// checkParityFlag checks the status of the parity flag.
func CheckParityFlag() (bool, error) {
	l := &ledger.Ledger{}
	return l.ConditionalFlagsLedger.CheckFlag(ParityFlag, ParityFlag)
}

// setAuxiliaryFlag sets the auxiliary flag.
func SetAuxiliaryFlag() error {
	l := &ledger.Ledger{}
	if err := l.ConditionalFlagsLedger.SetFlag(AuxiliaryFlag, AuxiliaryFlag, true); err != nil {
		return fmt.Errorf("failed to set auxiliary flag: %v", err)
	}
	return nil
}

// clearAuxiliaryFlag clears the auxiliary flag.
func ClearAuxiliaryFlag() error {
	l := &ledger.Ledger{}
	if err := l.ConditionalFlagsLedger.SetFlag(AuxiliaryFlag, AuxiliaryFlag, false); err != nil {
		return fmt.Errorf("failed to clear auxiliary flag: %v", err)
	}
	return nil
}

// checkAuxiliaryFlag checks the status of the auxiliary flag.
func CheckAuxiliaryFlag() (bool, error) {
	l := &ledger.Ledger{}
	return l.ConditionalFlagsLedger.CheckFlag(AuxiliaryFlag, AuxiliaryFlag)
}

// setNegativeFlag sets the negative flag.
func SetNegativeFlag() error {
	l := &ledger.Ledger{}
	if err := l.ConditionalFlagsLedger.SetFlag(NegativeFlag, NegativeFlag, true); err != nil {
		return fmt.Errorf("failed to set negative flag: %v", err)
	}
	return nil
}

// clearNegativeFlag clears the negative flag.
func ClearNegativeFlag() error {
	l := &ledger.Ledger{}
	if err := l.ConditionalFlagsLedger.SetFlag(NegativeFlag, NegativeFlag, false); err != nil {
		return fmt.Errorf("failed to clear negative flag: %v", err)
	}
	return nil
}

// checkNegativeFlag checks the status of the negative flag.
func CheckNegativeFlag() (bool, error) {
	l := &ledger.Ledger{}
	return l.ConditionalFlagsLedger.CheckFlag(NegativeFlag, NegativeFlag)
}

// setPositiveFlag sets the positive flag.
func SetPositiveFlag() error {
	l := &ledger.Ledger{}
	if err := l.ConditionalFlagsLedger.SetFlag(PositiveFlag, PositiveFlag, true); err != nil {
		return fmt.Errorf("failed to set positive flag: %v", err)
	}
	return nil
}

// clearPositiveFlag clears the positive flag.
func ClearPositiveFlag() error {
	l := &ledger.Ledger{}
	if err := l.ConditionalFlagsLedger.SetFlag(PositiveFlag, PositiveFlag, false); err != nil {
		return fmt.Errorf("failed to clear positive flag: %v", err)
	}
	return nil
}

// checkPositiveFlag checks the status of the positive flag.
func CheckPositiveFlag() (bool, error) {
	l := &ledger.Ledger{}
	return l.ConditionalFlagsLedger.CheckFlag(PositiveFlag, PositiveFlag)
}

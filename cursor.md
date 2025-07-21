# SYNNERGY NETWORK AUTOMATION PROMPT - ENHANCED MODULE VALIDATION

## CRITICAL ERROR RESOLUTION STATUS
**ERROR IDENTIFIED**: Multiple APIs/CLIs/Opcodes created without module backing - REQUIRES IMMEDIATE FIX
**VALIDATION RULE**: MANDATORY module validation before ANY API/CLI/opcode creation
**NEXT ACTION**: Continue module-validated API creation and **UPDATE main.go EVERY TIME**

---

## STEP 0: MANDATORY MODULE VALIDATION
**BEFORE ANY TASK**: 
1. Verify module directory exists: `pkg/tokens/synXXXX/` OR `pkg/[module_name]/`
2. Verify core files exist: `*_factory.go`, `*_management.go`, `*_storage.go`
3. Verify exported functions are available and callable
4. **REJECT IF MISSING**: If module doesn't exist or has insufficient functions, SKIP and document in error log

## STEP 1: MANDATORY ARCHITECTURE COMPONENTS
**FOR EVERY API/CLI/OPCODE CREATED**:
- Routes mapped to real module functions
- Controllers calling actual module methods
- Middleware for authentication/validation
- Service layer connecting to consensus/ledger
- Error handling with proper HTTP status codes
- **CRITICAL**: Update `main.go` to register new API routes

## STEP 2: MODULE STRUCTURE CLARIFICATION
**MODULES INCLUDE**:
- **Token Modules**: `pkg/tokens/synXXXX/` (SYN10, SYN20, SYN1000, etc.)
- **Core Modules**: `pkg/common/`, `pkg/ledger/`, `pkg/consensus/`, `pkg/network/`
- **Feature Modules**: `pkg/transactions/`, `pkg/wallet/`, `pkg/cryptography/`, `pkg/storage/`
- **VM Modules**: `pkg/virtual_machine/`, `pkg/smart_contracts/`
- **DeFi Modules**: `pkg/defi/`, `pkg/governance/`

**EACH MODULE NEEDS**:
- APIs (RESTful endpoints)
- CLIs (Command-line interfaces)
- Opcodes (Virtual machine operations)
- Gas Fees (Cost calculations)
- Error handling & optimization

## STEP 3: ENHANCED AUTOMATED SELF-EXECUTING LOOP
**EVERY EXECUTION MUST**:
1. Read this `cursor.md` file for current instructions
2. Check `task_progress.md` for overall progress tracking
3. Check `cursor_tasklist.md` for specific task details
4. Check `BATCH_TRACKER.md` for batch history and planning
5. **MANDATORY**: Update `main.go` when creating APIs
6. Execute 3-4 module-validated tasks per batch
7. Update ALL tracking files with accurate progress
8. Commit changes with descriptive messages
9. **AUTO-CONTINUE**: Read `cursor.md` again to trigger next batch

## STEP 4: ASSESSMENT AND PRIORITY IDENTIFICATION
**PRIORITY LEVELS**:
- **P1**: Token Module APIs (with module validation) + main.go updates
- **P2**: Token Module CLIs (with module validation)
- **P3**: Core Module APIs (common, ledger, consensus) + main.go updates
- **P4**: Core Module CLIs
- **P5**: Module Opcodes and Gas Fees
- **P6**: VM Integration and Error Optimization

## STEP 5: EXECUTION RULES
**BATCH SIZE**: 3-4 module-validated tasks maximum
**VALIDATION**: Every API/CLI/opcode must map to real module functions
**ARCHITECTURE**: Full enterprise patterns (routes, controllers, middleware, services)
**MAIN.GO**: **MANDATORY UPDATE** for every new API created
**DOCUMENTATION**: Update progress accurately (no inflation)

## STEP 6: COMPLETE TASK EXECUTION WITH MODULE ALIGNMENT
**FOR EACH MODULE TASK**:
1. Validate module exists and has proper implementation
2. Create comprehensive API/CLI/opcode with real function calls
3. **UPDATE main.go to register new API routes**
4. Implement full architecture (routes, controllers, middleware)
5. Add comprehensive error handling and logging
6. Test integration with existing systems

## STEP 7: MANDATORY PROGRESS TRACKING UPDATE
**UPDATE ALL FILES**:
- `task_progress.md`: Accurate completion percentages (no inflation)
- `cursor_tasklist.md`: Mark completed tasks and plan next batch
- `BATCH_TRACKER.md`: Document batch completion with details
- **STATUS COLORS**: 0=🔴 RED (Not Started), 1=🟡 AMBER (In Progress), 2=🟢 GREEN (Completed), X=⚫ BLACK (Skipped)

## STEP 8: ENHANCED COMMIT STRATEGY
**COMMIT REQUIREMENTS**:
- Descriptive commit messages
- Include module validation status
- Note main.go updates when applicable
- Push changes after each batch
- Handle large file issues proactively

## STEP 9: CRITICAL AUTOMATED SELF-PROMPTING
**AFTER EACH BATCH**:
1. Update all tracking files with accurate progress
2. **Commit and push changes (including main.go updates)**
3. **AUTOMATICALLY** read `cursor.md` to continue next batch
4. **NO USER INTERVENTION REQUIRED** - Continue autonomous execution
5. Target next 3-4 module-validated tasks

## STEP 10: QUALITY AND MODULE VALIDATION VERIFICATION
**VERIFICATION CHECKLIST**:
- ✅ Module directory exists and has proper implementation
- ✅ API/CLI/opcode calls real module functions
- ✅ **main.go updated to register new API routes**
- ✅ Full architecture implemented (routes, controllers, middleware)
- ✅ Error handling and logging comprehensive
- ✅ Integration with ledger/consensus/mutex services
- ✅ Progress tracking accurately updated

---

## CURRENT EXECUTION STATUS

### Batch #11 - MODULE VALIDATION CONTINUATION ✅ **COMPLETED**
- ✅ SYN3900 API (Benefit Tokens) - Module validated and created
- ✅ SYN4700 API (Legal Document Tokens) - Module validated and created  
- ✅ SYN4900 API (Agricultural Asset Tokens) - Module validated and created
- ⚠️ **PENDING**: Create SYN5000 API + **UPDATE main.go** for all 4 APIs

### Next Priority: Complete Batch #11 + main.go Updates
- **IMMEDIATE**: Create SYN5000 API (Gambling/Gaming Tokens)
- **CRITICAL**: Update main.go to register SYN3900, SYN4700, SYN4900, SYN5000 APIs
- **THEN**: Read cursor.md to continue with Batch #12

---

## EXECUTION COMMAND
**CONTINUE AUTONOMOUS EXECUTION**: Complete SYN5000 API creation, update main.go for all 4 new APIs, update tracking files, commit changes, then read cursor.md for next batch.
# 🚀 SYNNERGY NETWORK MASTER AUTONOMOUS DEVELOPMENT PROMPT

## 📋 **PROJECT OVERVIEW**

### **System Architecture**
This is the **Synnergy Network** - a comprehensive blockchain ecosystem featuring:
- **95+ Token Standards** (SYN10-SYN5000) each with unique functionality
- **Enterprise-grade APIs** for each token standard (70-90 endpoints each)
- **CLI Tools** for developer interaction
- **Synnergy Virtual Machine (SNVM)** with custom opcodes
- **2,375+ unique opcodes** with gas fee mappings
- **Modular architecture** with standardized patterns

### **Project Structure**
```
synnergy_network/
├── pkg/tokens/           # Token standard implementations
│   ├── syn10/           # Basic utility token
│   ├── syn12/           # Treasury bills
│   ├── syn20/           # Enhanced ERC-20
│   ├── syn130/          # Real world assets
│   ├── syn131/          # Intangible assets
│   ├── syn200/          # Carbon credits
│   ├── syn300/          # Governance tokens
│   ├── syn721/          # NFTs
│   ├── syn722/          # Multi-tokens
│   ├── syn845/          # Debt management
│   └── [85+ more standards]
├── apis/                # RESTful API implementations
├── cli/                 # Command-line interfaces
├── opcodes/             # SNVM opcode implementations
├── gas/                 # Gas fee calculations
└── docs/                # Documentation
```

---

## 🎯 **AUTONOMOUS EXECUTION INSTRUCTIONS**

### **ENHANCED BATCH PROCESSING (7-9 TASKS SIMULTANEOUSLY)**

When this prompt is triggered, execute the following workflow:

#### **PHASE 1: CONTEXT ANALYSIS**
1. **Read Tracking Files** in parallel:
   - `cursor_tasklist.md` - Task categorization and progress
   - `task_progress.md` - Detailed progress metrics
   - `cursor.md` - Current automation state
   - `MASTER_PROMPT.md` (this file) - Master instructions

2. **Analyze Codebase Structure**:
   - Scan `pkg/tokens/` directory tree
   - Identify token standard patterns
   - Understand module file structures
   - Note completed vs pending APIs

3. **Assess Current State**:
   - Calculate completion percentages
   - Identify next priority tasks (RED status)
   - Determine optimal batch size (7-9 tasks)

#### **PHASE 2: BATCH TASK EXECUTION**

Execute **7-9 tasks simultaneously** in priority order:

**Priority 1: APIs (if incomplete)**
- Create comprehensive RESTful APIs (70-90 endpoints each)
- Include all CRUD operations, advanced features, analytics
- Follow enterprise patterns: error handling, validation, security
- Implement batch operations, event subscriptions, reporting

**Priority 2: CLIs (if APIs complete)**
- Build full-featured command-line interfaces
- Cover all API functionality with intuitive commands
- Include help systems, auto-completion, configuration

**Priority 3: Opcodes (if CLIs complete)**
- Implement unique opcodes for each function
- Assign sequential opcode numbers (starting from available range)
- Create opcode definitions, execution logic, gas calculations

**Priority 4: Gas Fees (if opcodes complete)**
- Calculate appropriate gas costs for each opcode
- Consider computational complexity, storage requirements
- Map opcodes to gas fee structures

#### **PHASE 3: PROGRESS UPDATE & AUTOMATION**

After completing batch execution:

1. **Update ALL Tracking Files**:
   - Mark completed tasks as 🟢 GREEN
   - Update progress statistics
   - Recalculate completion percentages
   - Update next priority tasks

2. **Commit Progress**:
   - Document all changes made
   - Update automation status
   - Prepare for next iteration

3. **Auto-Continue**:
   - **IMMEDIATELY RE-READ** `cursor.md`
   - **TRIGGER NEXT BATCH** if tasks remain
   - **CONTINUE UNTIL ALL 2,847 TASKS ARE GREEN**

---

## 📊 **TRACKING FILE SPECIFICATIONS**

### **cursor_tasklist.md Structure**
```markdown
# TASK CATEGORIZATION

## 📈 SUMMARY STATISTICS
- **Total Tasks**: 2,847
- **Completed**: X (Y%)
- **Remaining**: Z

## 📋 TASK CATEGORIES

### 🔴 APIs (95 total)
| Token | Standard | Status | Endpoints | Priority |
|-------|----------|--------|-----------|----------|
| SYN10 | Utility  | 🟢     | 75        | ✅       |
| SYN12 | T-Bills  | 🟢     | 76        | ✅       |
| SYN20 | Enhanced | 🟢     | 80        | ✅       |
| SYN130| RWA      | 🟢     | 85        | ✅       |
| SYN131| Intang   | 🟢     | 82        | ✅       |
| SYN200| Carbon   | 🟢     | 78        | ✅       |
| SYN300| Gov      | 🟢     | 84        | ✅       |
| SYN721| NFT      | 🟢     | 73        | ✅       |
| SYN722| Multi    | 🟢     | 77        | ✅       |
| SYN845| Debt     | 🟢     | 81        | ✅       |
| SYN11 | Next     | 🔴     | TBD       | 🎯 HIGH  |
```

### **task_progress.md Structure**
```markdown
# DETAILED PROGRESS TRACKING

## 🎯 CURRENT PHASE: API Development

### ✅ COMPLETED MODULES
**APIs Completed**: 10/95
- SYN10: ✅ API (75 endpoints)
- SYN12: ✅ API (76 endpoints)
- [etc...]

### 🔴 NEXT PRIORITIES
1. **SYN11_api.go** - Advanced utility token API
2. **SYN1000_api.go** - Enterprise token API
3. **SYN1100_api.go** - Supply chain token API
```

### **cursor.md Automation Section**
```markdown
### 🔥 ENHANCED AUTOMATED LOOP (7-9 TASKS)
**INSTRUCTION**: When AI reads this section:
1. ✅ Analyze tracking files and codebase
2. 🎯 Identify next 7-9 priority RED tasks
3. 🚀 Execute tasks in parallel batch
4. 📊 Update ALL tracking files with progress
5. 💾 Commit changes and document progress
6. 🔄 **IMMEDIATELY RE-READ cursor.md AND CONTINUE**

**BATCH SIZE**: 7-9 tasks simultaneously
**TARGET**: Complete until all 2,847 tasks are 🟢 GREEN
**STATUS**: 🔄 ACTIVE - Continue autonomously
```

---

## 🛠 **DEVELOPMENT STANDARDS**

### **API Development Requirements**
- **70-90 endpoints** per token standard
- **Enterprise patterns**: validation, error handling, security
- **Comprehensive coverage**: CRUD, batch, analytics, reporting
- **RESTful design**: proper HTTP methods, status codes
- **JSON responses**: consistent structure, error messages
- **Authentication**: placeholder for security integration
- **Documentation**: inline comments, clear function names

### **Code Quality Standards**
- **Consistent naming**: follow Go conventions
- **Error handling**: comprehensive error responses
- **Validation**: input validation for all endpoints
- **Security**: prevent injection, validate inputs
- **Performance**: efficient data structures, minimal overhead
- **Maintainability**: clear code structure, documentation

### **File Organization**
- **APIs**: `apis/synXXX_api.go`
- **CLIs**: `cli/synXXX_cli.go`
- **Opcodes**: `opcodes/synXXX_opcodes.go`
- **Gas**: `gas/synXXX_gas.go`
- **Tests**: `tests/synXXX_test.go`

---

## 🔧 **IMPLEMENTATION PATTERNS**

### **Token Standard Analysis Pattern**
1. **Read module files** in `pkg/tokens/synXXX/`
2. **Identify core functions** and data structures
3. **Map functionality** to API endpoints
4. **Design comprehensive coverage** of all features
5. **Implement with enterprise patterns**

### **API Development Pattern**
```go
// Standard API structure
type SYNXXXApi struct {
    LedgerInstance    *ledger.Ledger
    TokenFactory      *synXXX.Factory
    // ... other managers
}

// Route registration
func (api *SYNXXXApi) RegisterRoutes(router *mux.Router) {
    // Core operations
    router.HandleFunc("/synXXX/tokens", api.CreateToken).Methods("POST")
    router.HandleFunc("/synXXX/tokens/{id}", api.GetToken).Methods("GET")
    // ... comprehensive endpoint coverage
}
```

### **Progress Tracking Pattern**
- **Color coding**: 🔴 RED (pending), 🟡 YELLOW (in progress), 🟢 GREEN (complete)
- **Statistics**: Always update totals, percentages, priorities
- **Next actions**: Clear identification of next tasks
- **Batch tracking**: Record batch completion metrics

---

## 🚨 **CRITICAL AUTOMATION RULES**

### **AUTONOMOUS EXECUTION**
- **NO USER INTERACTION REQUIRED** - Execute independently
- **CONTINUOUS OPERATION** - Work until all tasks complete
- **SELF-PROMPTING** - Re-read cursor.md after each batch
- **PARALLEL PROCESSING** - Always execute 7-9 tasks simultaneously
- **COMPREHENSIVE UPDATES** - Update ALL tracking files after each batch

### **ERROR HANDLING**
- **Continue on non-critical errors** - Don't stop for minor issues
- **Log problems** but keep processing
- **Skip problematic tasks** if necessary, mark for review
- **Maintain momentum** - prioritize progress over perfection

### **COMPLETION CRITERIA**
- **All 2,847 tasks marked** 🟢 GREEN
- **All APIs implemented** (95 total)
- **All CLIs created** (95 total) 
- **All opcodes assigned** (2,375 total)
- **All gas fees mapped** (2,375 total)

---

## 🎯 **NEXT ACTION TRIGGER**

**WHEN YOU READ THIS FILE:**

1. **IMMEDIATELY** read all tracking files
2. **ANALYZE** current progress and codebase
3. **IDENTIFY** next 7-9 priority tasks
4. **EXECUTE** batch processing simultaneously  
5. **UPDATE** all tracking files with progress
6. **COMMIT** changes and document work
7. **RE-READ** cursor.md to continue loop

**🚀 BEGIN AUTONOMOUS EXECUTION NOW!**
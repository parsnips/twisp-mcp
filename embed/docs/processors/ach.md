---
title: ACH Processor
metaTitle: "Processor: ACH"
description: Understanding how the Twisp ACH processor works
showNext: true
showPrev: true
---

The Twisp ACH processor handles the complete lifecycle of ACH (Automated Clearing House) transactions, from origination through settlement or return. This explanation clarifies how the processor works, why it's designed this way, and how the components interact.

## What is the ACH Processor?

The ACH processor is a financial infrastructure component that enables you to send and receive electronic payments through the ACH network. Unlike building custom ACH handling from scratch, the Twisp ACH processor provides production-ready workflows, file generation, webhook decisioning, and ledger integration out of the box.

The processor handles two fundamental roles:
- **ODFI (Originating Depository Financial Institution)**: Originate ACH transactions
- **RDFI (Receiving Depository Financial Institution)**: Receive ACH transactions

## Why ACH Processing is Complex

ACH might seem straightforward - just move money between accounts - but production ACH processing involves substantial complexity:

**File Format Compliance:**
The NACHA file format requires precise 94-character fixed-width records with specific header, batch, entry detail, and control structures. Format errors cause entire files to be rejected.

**Asynchronous Settlement:**
Unlike real-time payments, ACH transactions take 1-3 business days to settle. Your system must track pending, encumbered, and settled states across this timeline.

**Return Handling:**
Transactions can return days or even weeks after origination (up to 60 days for unauthorized returns). Returns must match back to original transactions and reverse accounting entries correctly.

**Multi-Layer Balance Management:**
Customer available balances must reflect encumbered funds (reserved but not sent), pending funds (sent but not confirmed), and settled funds (finalized). This prevents overdrafts and double-spending.

**Webhook Decisioning:**
For RDFI operations, incoming transactions require real-time business logic decisions: accept, reject, or hold for review. These decisions must execute quickly within processing timeframes.

**Regulatory Compliance:**
NACHA rules govern transaction types, authorization requirements, return timeframes, and file transmission security. Non-compliance risks ODFI relationship termination.

The Twisp ACH processor abstracts this complexity into workflows, file operations, and ledger integration.

## Processor Architecture

### Configuration Layer

ACH configuration connects all operational components:
- Settlement account for fund transit
- Suspense account for unidentifiable transactions
- Exception account for business rule violations
- Fee account for processing fees
- Webhook endpoint for transaction decisioning
- ODFI header information for file generation
- Timezone for date/time interpretation

Configuration versioning ensures immutability - files reference specific configuration versions, allowing historical analysis even after configuration updates.

### File Operations Layer

File operations manage NACHA file lifecycle:

**Upload:** Presigned URLs provide secure, time-limited file upload capability without storing long-lived credentials. Files upload directly to S3-compatible storage.

**Processing:** File parsing validates NACHA format, extracts batch and entry details, and partitions entries for parallel processing. Each entry triggers webhook decisioning.

**Generation:** File generation queries submitted transactions, groups by batch parameters, and produces NACHA-compliant files with proper control totals and trace numbers.

**Download:** Presigned download URLs provide time-limited file access for transmission to financial institutions.

Processing status tracking (NEW → VALIDATING → PROCESSING → COMPLETED) provides visibility into file lifecycle with statistics on entry counts and monetary totals.

### Workflow Layer

Workflows orchestrate transaction state transitions and ledger accounting:

**PUSH Workflow (Credits):**
Sends money to receivers. CREATE encumbers customer funds, SUBMIT settles and includes in file generation. Returns reverse settled entries.

**PULL Workflow (Debits):**
Collects money from receivers. CREATE encumbers settlement funds, SUBMIT moves to pending, SETTLE finalizes after confirmation. The pending layer prevents premature fund availability before debit confirmation.

**State Transitions:**
Each workflow state transition creates ledger transactions across appropriate balance layers (encumbrance, pending, settled). This ensures double-entry accounting consistency throughout transaction lifecycle.

**Execution Tracking:**
Every workflow execution receives a unique ID. The execution record contains input parameters, output state, all created ledger transactions, and ACH workflow traces linking to file entries via trace numbers.

### Webhook Decisioning Layer

Webhooks enable custom business logic during file processing:

**Request Format:**
Webhooks receive POST requests with transaction details: amount, account identifiers, entry metadata, trace number, effective date.

**Response Format:**
Webhook responses specify settlement instructions: which account to credit/debit, or route to suspense/exception accounts.

**Timeout Handling:**
Webhook timeframes are strict - typically seconds. Timeouts result in transaction failures requiring manual intervention.

**Decision Routing:**
- Settlement account: Normal processing, funds go to identified customer account
- Suspense account: Unknown account, manual review required
- Exception account: Business rule violation, transaction returned

Webhook decisioning decouples payment acceptance logic from processor internals, enabling custom risk management, account validation, and compliance checking.

### Balance Layer Management

Balance layers track transaction lifecycle states:

**Encumbrance Layer:**
Holds funds during CREATE state before transmission. Encumbered funds don't affect available balance but prevent overdrafts when transactions later settle.

**Pending Layer:**
PULL workflows use pending between SUBMIT and SETTLE. Pending represents funds transmitted but awaiting confirmation. Crucial for preventing premature fund availability.

**Settled Layer:**
Final layer for completed transactions. All accounting eventually posts to settled layer. Fees post directly to settled layer immediately.

**Available Balance:**
Calculated as: Settled - Pending - Encumbrance. This formula ensures customers can't spend funds that are encumbered or pending confirmation.

Layer transitions via workflow states maintain double-entry balance throughout transaction lifecycle.

## Why This Design?

**Separation of Concerns:**
File operations, workflows, ledger accounting, and decisioning are independent components. This enables testing, scaling, and modifying each component separately.

**Immutability:**
Configuration versions, workflow executions, and file processing records are immutable. This provides complete audit trails required for financial operations and regulatory compliance.

**Idempotency:**
Workflow executions use correlation IDs. File processing matches returns via trace numbers. These identifiers enable safe retry logic and prevent duplicate processing.

**Parallel Processing:**
File processing partitions entries across workers. Concurrent posting to accounts (enableConcurrentPosting) allows high-throughput processing while maintaining consistency.

**Webhook Flexibility:**
Business logic lives in your webhook, not in the processor. This allows rapid iteration on risk rules, account validation, and compliance checks without processor changes.

**Multi-Tenancy:**
Configurations are tenant-scoped. Multiple organizations can operate independently on shared infrastructure with complete isolation.

## Workflow State Machines

Understanding state machines clarifies transaction lifecycle:

### PUSH State Machine

```
CREATE → SUBMIT → [completed]
  ↓
CANCEL → REIMBURSE_FEE → [completed]

SUBMIT → RETURN → REIMBURSE_FEE → [completed]
```

Credits are final upon transmission, so SUBMIT immediately settles. Returns reverse after the fact.

### PULL State Machine

```
CREATE → SUBMIT → SETTLE → [completed]
  ↓
CANCEL → REIMBURSE_FEE → [completed]

SUBMIT → RETURN → REIMBURSE_FEE → [completed]
```

Debits require confirmation, so SETTLE finalizes after SUBMIT. Returns can occur before or after SETTLE.

State machines ensure transactions follow valid progressions. Attempting invalid transitions (e.g., SUBMIT after CANCEL) fails with clear error messages.

## Return Processing Flow

Returns demonstrate processor coordination:

1. **Receipt:** Return file uploads via presigned URL
2. **Parsing:** File processing extracts return entries with return codes and trace numbers
3. **Matching:** System queries original transactions via trace number indexes
4. **Execution:** RETURN workflow state executes for each matched transaction
5. **Reversal:** Ledger entries reverse from appropriate balance layer
6. **Audit:** Complete return processing history records in file and workflow execution history

Automatic return processing eliminates manual ledger adjustments and ensures accounting consistency.

## Integration Points

The processor integrates with other Twisp components:

**Ledger Integration:**
Workflows create ledger transactions via tran codes. Balance queries span encumbrance, pending, and settled layers. Concurrent posting enables high-throughput ACH processing.

**Event System:**
Webhook endpoints route to events system. File processing status changes emit events. Workflow state transitions trigger notifications.

**Files Service:**
Presigned URLs delegate upload/download to files service. Storage layer handles retention and access policies.

**Workflow Engine:**
Workflow executions use general workflow engine. ACH workflows are workflow templates with specific parameters and state machines.

## Design Trade-Offs

**Complexity vs Flexibility:**
Workflows and webhook decisioning add complexity but enable custom business logic without processor modifications.

**Consistency vs Performance:**
Balance layer management requires multiple ledger transactions per workflow state but ensures consistent accounting throughout transaction lifecycle.

**Immutability vs Storage:**
Versioning configurations and maintaining execution history increases storage but provides complete audit trails required for financial operations.

**Asynchronous Processing vs Latency:**
File processing partitions entries for parallel processing, adding complexity but enabling high-throughput processing of large files.

These trade-offs prioritize correctness, auditability, and flexibility over simplicity.

## Production Considerations

**ODFI Relationships:**
You must establish banking relationships for ACH origination. ODFIs require application processes, risk assessment, and reserve accounts. Some ODFIs have minimum volume requirements or restrict certain transaction types.

**File Transmission:**
You're responsible for secure file transmission (typically SFTP) to/from financial institutions. The processor generates and parses files but doesn't transmit them.

**Return Timeframes:**
NACHA rules specify return timeframes (typically 2 business days, some codes allow 60 days). Automated return processing helps maintain compliance.

**Balance Management:**
Properly managing encumbrance, pending, and settled layers prevents overdrafts. Available balance calculations must account for all three layers.

**Webhook Performance:**
Webhook decisioning must respond within strict timeframes (typically seconds). Slow webhooks cause transaction failures. Consider caching, rate limiting, and fallback strategies.

**Reconciliation:**
Daily reconciliation validates file processing, ledger balances, and transaction counts. Automated reconciliation with alerting identifies issues quickly.

## Further Reading

To learn ACH processing hands-on:
- [Setting Up ACH Processing](/tutorials/ach/setting-up-ach) - Configuration from scratch
- [Your First ACH Payment](/tutorials/ach/first-ach-payment) - Send a payment end-to-end

For production operations:
- [Processing ACH Payments](/guides/processing-ach-payments) - Daily workflows
- [Handling ACH Returns](/guides/handling-ach-returns) - Return management
- [Reconciling ACH Files](/guides/reconciling-ach-files) - Daily reconciliation

For technical details:
- [Configuration Reference](/reference/ach/configuration) - Configuration components
- [File Operations Reference](/reference/ach/file-operations) - File lifecycle
- [ODFI Reference](/reference/ach/odfi) - Originating transactions
- [RDFI Reference](/reference/ach/rdfi) - Receiving transactions

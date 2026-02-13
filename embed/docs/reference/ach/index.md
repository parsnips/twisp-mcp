---
title: ACH Reference
metaTitle: "Reference: ACH Processor"
description: The Twisp ACH processor handles configuration, file operations, ODFI origination, and RDFI receiving operations.
showNext: true
showPrev: false
---

## ACH Resources

See each of the sub-pages in this reference for details about the corresponding ACH resource. Each resource is accessible via the [GraphQL API](/reference/graphql).

{% quick-links %}

{% quick-link title="Configuration" href="/reference/ach/configuration" icon="configure" description="ACH configuration defines accounts, webhooks, ODFI headers, and processing parameters." /%}

{% quick-link title="File Operations" href="/reference/ach/file-operations" icon="copySuccess" description="Upload, process, generate, and download NACHA-formatted ACH files." /%}

{% quick-link title="ODFI" href="/reference/ach/odfi" icon="credit" description="Originate ACH transactions via PUSH (credit) and PULL (debit) workflows." /%}

{% quick-link title="RDFI" href="/reference/ach/rdfi" icon="debit" description="Receive and process incoming ACH transactions with webhook decisioning." /%}

{% /quick-links %}


## Key Concepts of the Twisp ACH Processor

_This set of concepts summarizes the key behavior of the Twisp ACH processor._

- **Configuration** defines the operational parameters for ACH processing.
  - Each configuration connects settlement, suspense, exception, and fee accounts.
  - Configurations reference webhook endpoints for transaction decisioning.
  - ODFI header configuration specifies company and bank information for NACHA files.
  - Configurations support both ODFI (originating) and RDFI (receiving) operations.
  - Configuration versions are tracked, ensuring file processing references the correct parameters.

- **File Operations** manage the complete lifecycle of NACHA-formatted files.
  - Upload operations create presigned URLs for secure file uploads.
  - File processing parses NACHA files, validates format, and extracts transactions.
  - File generation creates NACHA-compliant files from submitted workflow transactions.
  - Download operations provide presigned URLs for retrieving generated files.
  - File processing status progresses through states (NEW, VALIDATING, PROCESSING, COMPLETED).
  - Processing statistics track entry counts and monetary totals in real-time.

- **ODFI (Originating Depository Financial Institution)** operations originate ACH transactions.
  - PUSH workflows send credits (payroll, vendor payments, refunds).
  - PULL workflows collect debits (bill payments, subscriptions, loan payments).
  - Workflows manage funds through encumbrance, pending, and settled balance layers.
  - PUSH workflows settle immediately upon SUBMIT since credits are final.
  - PULL workflows separate SUBMIT (transmission) from SETTLE (confirmation).
  - Return processing automatically matches returns to original transactions via trace numbers.
  - All workflow executions maintain complete audit trails via execution history.

- **RDFI (Receiving Depository Financial Institution)** operations receive ACH transactions.
  - Inbound file processing extracts transactions and triggers webhook decisioning.
  - Webhooks receive transaction details and respond with settlement instructions.
  - Transactions route to settlement, suspense, or exception accounts based on webhook responses.
  - RDFI workflows support both CREATE (initial acceptance) and SETTLE (final settlement) states.
  - Return file generation creates NACHA files for rejected transactions.
  - NOC (Notification of Change) processing updates account information automatically.

- **Workflows** execute state transitions for ACH transactions.
  - Each workflow execution receives a unique execution ID for tracking.
  - Workflow states (CREATE, SUBMIT, SETTLE, RETURN, CANCEL, REIMBURSE_FEE) define transaction lifecycle.
  - Ledger transactions are created automatically based on workflow state transitions.
  - ACH workflow traces link executions to file entries via trace numbers.
  - Workflow history provides complete audit trail of all state transitions.

- **Balance Layers** track transaction lifecycle across multiple states.
  - Encumbrance layer holds funds during CREATE state before submission.
  - Pending layer tracks PULL transactions between SUBMIT and SETTLE.
  - Settled layer contains finalized transactions.
  - Available balance calculation: Settled - Pending - Encumbrance.
  - Layer transitions ensure funds are properly accounted at each stage.

- **Webhook Decisioning** enables custom business logic for transaction processing.
  - Webhooks receive POST requests with transaction details during file processing.
  - Webhook responses specify which account to credit/debit for each transaction.
  - Transactions can be routed to settlement, suspense, or exception accounts.
  - Webhook timeouts result in transaction failure requiring manual intervention.
  - Webhook decisions determine whether transactions settle or return.

- **Return Processing** handles rejected transactions from financial institutions.
  - ODFI returns are received for transactions you originated.
  - RDFI returns are generated for transactions you received and must reject.
  - Returns automatically match to original transactions via trace numbers.
  - Return processing executes RETURN workflow state, reversing ledger entries.
  - Return codes (R01-R99) indicate rejection reasons per NACHA standards.
  - Return files must be transmitted within NACHA-specified timeframes.

- **File Types** specify the purpose and direction of ACH files.
  - RDFI types: RDFI (incoming), RDFI_RETURN (returns to send), RDFI_NOC (notifications).
  - ODFI types: ODFI (mixed), ODFI_PUSH_ONLY (credits), ODFI_PULL_ONLY (debits), ODFI_RETURN (returns received).
  - File types control which transactions are included during generation.
  - Preprocessing file types enable filtering before final file generation.

- **Reconciliation** validates file processing and ledger accuracy.
  - File statistics provide entry counts and monetary totals for validation.
  - Processing status indicates completion or errors requiring investigation.
  - Balance reconciliation ensures settlement accounts reflect all ACH activity.
  - Suspense and exception accounts should be cleared regularly via returns or transfers.
  - File history enables point-in-time analysis of processing state.

## Further Reading

To learn ACH processing from scratch, see:
- [Setting Up ACH Processing](/tutorials/ach/setting-up-ach) - Complete configuration tutorial
- [Your First ACH Payment](/tutorials/ach/first-ach-payment) - Send a payment end-to-end

For practical production workflows, see:
- [Processing ACH Payments](/guides/processing-ach-payments) - Daily payment operations
- [Handling ACH Returns](/guides/handling-ach-returns) - Return management
- [Reconciling ACH Files](/guides/reconciling-ach-files) - Daily reconciliation

For conceptual understanding, see:
- [ACH Processor](/processors/ach) - How the ACH processor works

For complete GraphQL type definitions, see:
- []($gql:object:AchConfiguration)
- []($gql:object:AchFileInfo)
- []($gql:object:WorkflowExecution)

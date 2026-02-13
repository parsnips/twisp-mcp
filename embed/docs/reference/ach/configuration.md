---
title: Configuration
metaTitle: "Reference: ACH Configuration"
description: Reference for ACH configuration within the Twisp ACH Processor
showNext: true
showPrev: true
---

## The Basics

ACH configurations define the operational parameters for processing ACH transactions. Each configuration specifies the accounts, webhook endpoints, file header information, and timezone settings required for both originating (ODFI) and receiving (RDFI) ACH operations.

Configurations in Twisp...

- Connect settlement, suspense, exception, and fee accounts
- Define ODFI header information for NACHA file generation
- Reference webhook endpoints for transaction decisioning
- Specify journal for posting all ACH transactions
- Support both ODFI and RDFI operations with unified setup

## Components of ACH Configuration

There are 5 primary components which define an ACH configuration:

1. **Accounts**: Four required accounts (settlement, suspense, exception, fee) that handle different transaction scenarios during ACH processing.
2. **Webhook Endpoint**: Receives decisioning requests during file processing, allowing custom business logic for transaction acceptance or rejection.
3. **ODFI Header Configuration**: Company and bank information used when generating NACHA-formatted files for transmission.
4. **Journal**: The ledger journal where all ACH transaction entries are posted.
5. **Timezone**: IANA timezone identifier for date/time operations and effective date calculations.

In addition, configurations have common properties:

- **Config ID**: a universally unique identifier (UUID) for the configuration.
- **Version**: configurations are versioned, allowing tracking of changes over time and ensuring file processing references the correct configuration version.
- **Created & Modified Timestamps**: when the configuration was created and last updated.

## Required Accounts

### Settlement Account

The settlement account is the central hub for ACH fund flows. All incoming and outgoing ACH transactions post to this account during processing phases.

**Characteristics:**
- Normal balance type: DEBIT (typically holds positive balance from incoming credits)
- `enableConcurrentPosting: true` (supports high transaction volumes)
- Debited for ODFI PULL transactions, credited for ODFI PUSH transactions
- Credited for RDFI inbound credits, debited for RDFI inbound debits

### Suspense Account

Suspense account receives transactions where the destination account cannot be located, commonly when routing number and account number combinations don't match existing accounts.

**Characteristics:**
- `enableConcurrentPosting: true` (may receive high transaction volumes)
- Requires manual review to determine correct account or return to originator
- Funds held pending investigation or return processing

### Exception Account

Exception account receives transactions that fail processing due to velocity controls, locked accounts, or other business rule violations enforced by webhook decisioning.

**Characteristics:**
- `enableConcurrentPosting: true` (supports parallel processing)
- Funds should be returned to originator via return files
- Tracks transactions rejected by business rules

### Fee Account

Fee account collects ACH processing fees charged per transaction.

**Characteristics:**
- Normal balance type: CREDIT (accumulates fee income)
- `enableConcurrentPosting: true` (supports high transaction volumes)
- Credited when fees charged, debited when fees reimbursed (on returns/cancellations)

## ODFI Header Configuration

ODFI header configuration contains NACHA file header information required when generating ACH files.

**Components:**
- **Immediate Destination**: Routing number of receiving institution (ODFI or Federal Reserve)
- **Immediate Destination Name**: Name of receiving institution (max 23 characters)
- **Immediate Origin**: Your Federal Tax ID or routing number (10 characters)
- **Immediate Origin Name**: Your company name as it appears in ACH files (max 23 characters)

These values populate the File Header Record (Record Type Code 1) in generated NACHA files.

## Webhook Endpoint

The webhook endpoint receives POST requests during ACH file processing for transaction decisioning. The endpoint must be of type `ACH_PROCESSOR`.

**Request Format:**
Webhook receives transaction details including amount, account identifiers, and entry metadata when processing files via [`Mutation.ach.processFile()`]($gql:mutation:ach.processFile).

**Response Format:**
Webhook responds with settlement instructions specifying which account to credit/debit, or directing transaction to suspense/exception accounts.

**Timeout:**
Webhook must respond within configured timeout period to avoid transaction failure.

See [RDFI reference](/reference/ach/rdfi#transaction-decisioning) for detailed webhook payload and response specifications.

## Timezone

IANA timezone identifier determines date/time interpretation for:
- Effective dates for ACH transactions
- Settlement window calculations
- Batch header dates in NACHA files
- Processing cutoff time interpretation

Common values: `America/New_York`, `America/Chicago`, `America/Denver`, `America/Los_Angeles`, `America/Phoenix`

Use the timezone where primary ACH operations occur, typically matching ODFI timezone or business headquarters location.

## Configuration Versioning

Configurations use optimistic locking with version numbers. Each update via [`Mutation.ach.updateConfiguration()`]($gql:mutation:ach.updateConfiguration) increments the version field. File processing records reference the specific configuration version used, ensuring immutability and audit trail consistency.

When files are processed via [`Mutation.ach.processFile()`]($gql:mutation:ach.processFile), the `configVersion` field on `AchFileInfo` captures which configuration version applied, allowing historical analysis even after configuration updates.

## Configuration Operations

Use GraphQL to create, update, and query ACH configurations:

- [`Mutation.ach.createConfiguration()`]($gql:mutation:ach.createConfiguration): Create new ACH configuration.
- [`Mutation.ach.updateConfiguration()`]($gql:mutation:ach.updateConfiguration): Update existing configuration fields.
- [`Query.ach.configuration()`]($gql:query:ach.configuration): Get a single configuration by ID.
- [`Query.ach.configurations()`]($gql:query:ach.configurations): List all configurations with pagination.

## Further Reading

To learn how to set up ACH configuration from scratch, see the tutorial on [Setting Up ACH Processing](/tutorials/ach/setting-up-ach).

For step-by-step configuration with real APIs, see the how-to guide on [Processing ACH Payments](/guides/processing-ach-payments).

For complete GraphQL type definitions, see:
- []($gql:object:AchConfiguration)
- []($gql:object:AchOdfiHeaderConfiguration)
- []($gql:input:AchCreateConfigurationInput)
- []($gql:input:AchUpdateConfigurationInput)

---
title: Processors
description: Processing financial protocols and payment networks
showNext: true
showPrev: false
---

Processors are specialized components that handle communication with external payment networks and financial protocols. They bridge the Twisp accounting core with the outside world of financial transactions, enabling seamless integration with payment rails like ACH, wire transfers, and card networks.

## What are Processors?

Processors in Twisp handle the complexities of external payment protocols, including:

- **Data Format Management**: Converting between Twisp's internal formats and protocol-specific formats (e.g., NACHA for ACH)
- **Automated Accounting**: Ledgering, monitoring, and reconciling the lifecycle of payments through various stages (pending, settled, returned)
- **Error Handling**: Managing exceptions, returns, and corrections according to protocol rules

## Available Processors

Twisp currently supports the following processors:

- [**ACH**](/processors/ach)\
  Process ACH (Automated Clearing House) transactions including debits, credits, returns, and notifications of change for both ODFI and RDFI operators.

## Architecture

Processors work in conjunction with the Twisp accounting core to provide end-to-end payment processing:

1. **Origination**: Transactions are initiated through the accounting core
2. **Processing**: The processor formats and transmits files to the appropriate network
3. **Reconciliation**: Incoming files are parsed and reconciled with ledger transactions
4. **Completion**: Final status updates are reflected in the accounting core

Continue reading to learn more about specific processors and their capabilities.

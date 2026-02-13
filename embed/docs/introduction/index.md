---
title: Introduction
description: What is Twisp?
showNext: true
showPrev: false
---

Twisp is an accounting engine for building financial products.

We provide a financial ledger database with financial primitives for engineering systems that deal with money.

... and what is a Financial Ledger?

Financial ledgers are the foundation of systems that track money.

Some examples:

- A neobank tracks debits and credits across many accounts, with multiple financial instruments, all the while analyzing this data to surface financial insights to consumers.
- A payments infrastructure company needs robust multi-tenant account servicing built on many FBOs, safely made available to many clients via API.

Financial ledgers and the data derived from them are the engine of all such systems. Developing, securing, operating, and scaling these financial ledgers on a general purpose database is a difficult and time consuming proposition... which is why we built Twisp.

Upon provisioning a **Twisp core ledger**, you gain access to the following features:

1. **Transaction Ledger**: Facilitates double-entry accounting, ensuring comprehensive and accurate financial recording.
2. **Journals**: Enables tracking of different currencies and monitoring activities over specific time periods.
3. **Accounts**: Represents a chart of accounts for tracking any economic activity, offering a structured view of your financial data.
4. **Transaction Codes**: Provides a method for categorizing and tracking various transaction types within the system, enhancing transactional clarity.
5. **Layered Balances**: Maintains `available`, `settled`, and `encumbrance` balances, enabling precise and timely balance reporting. This includes:
   - Account Hierarchies and Roll-ups: Streamlines account management and analysis.
   - Dimensional Balance Tracking: Provides multi-dimensional views of your financial status.
   - Velocity Balances: Tracks account balance changes over time.
6. **Workflows**: Automates common processes for increased efficiency, including:
   - Card Authorization & Settlements: Manages the lifecycle of card transactions.
   - ACH Transaction Processing: Handles the flow of Automated Clearing House transactions.
   - mRDC Check Deposits: Facilitates mobile Remote Deposit Capture operations.
7. **High-Scale API**: Serves end customers directly with a robust API, which includes:
   - Activity Feed: Offers real-time updates on account activities.
   - Enrichment: Enhances transaction data for a more detailed understanding.
   - Reconciliation: Matches transactions back to the ledger to ensure accuracy and eliminate customer confusion.
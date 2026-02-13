---
title: Transaction Isolation and Consistency
description: Twisp provides snapshot isolation and strong consistency to ensure correctness and prevent anomalies in ledger systems.
showNext: true
showPrev: true
---

## Transaction Isolation in Twisp

In systems that demand correctness and high concurrency, transaction isolation plays a critical role in ensuring consistency and preventing anomalies. Twisp leverages an ACID **Multi-Version Concurrency Control (MVCC)** transaction layer to implement robust isolation levels tailored to the needs of ledger systems.

Twisp's MVCC provides:

- **Interactive transactions**: Supports all-or-nothing semantics for transactions, ensuring atomicity and isolation.
- **Snapshot Isolation**: Transactions operate on a consistent and isolated view of data, ensuring transactions operate on immutable data snapshots.
- **Anomaly prevention**: Guarantees that transactional operations produce fully consistent results without partial or invalid reads along with write skew prevention.
- **Strong consistency**: All committed balances and database states are always readable in a strongly consistent manner, ensuring no anomalies or partial updates.

{% callout type="note" %}
Twisp runs and maintains a formal model of our MVCC transaction layer. Learn more about the formal modelling process we conducted with our partner [Galois](https://galois.com/blog/2024/02/galois-twisp-avoiding-foolishness-in-distributed-systems/).
{% /callout %}

### Isolation Levels Overview

Twisp supports four distinct transaction isolation levels, balancing between performance and correctness:

1. **Read Committed**:
   - Ensures that the most recently committed data is visible to a transaction.
   - Use cases: Workflows requiring real-time, consistent reads that can tolerate a non-snapshot view of data.

2. **Snapshot Isolation (System Default)**:
   - Guarantees that transactions operate on a point-in-time consistent snapshot of the data.
   - Use cases: Workflows requiring strong correctness without paying a performance penalty.

3. **Repeatable Read**:
   - Guarantees that read data has not changed on transaction commit providing read stability.
   - Use cases: High-consistency workflows for complex, interactive, multi-row updates.

4. **Serializable**:
   - In addition to read stability, also provides phantom avoidance that ensures scans would not return additional data on transaction commit.
   - Use cases: Workflows where absolute correctness outweighs performance concerns.
   
{% callout type="note" %}
Transaction isolation settings can be set with the [`@tx` directive](/reference/graphql/directives#tx) in GraphQL.
{% /callout %}

### Preventing Transaction Anomalies

In high-concurrency environments, transaction anomalies can corrupt data or lead to business-critical errors. Twisp's isolation levels guard against common database anomalies. The following anomalies are thoroughly tested in every build via our conformance test suites:

- **G0**: Write Cycles (dirty writes)
- **G1a**: Aborted Reads (dirty reads, cascaded aborts)
- **G1b**: Intermediate Reads (dirty reads)
- **G1c**: Circular Information Flow (dirty reads)
- **OTV**: Observed Transaction Vanishes
- **PMP**: Predicate-Many-Preceders
- **P4**: Lost Update
- **G-single**: Single Anti-dependency Cycles (read skew)
- **G2-item**: Item Anti-dependency Cycles (write skew on disjoint read)
- **G2**: Anti-Dependency Cycles (write skew on predicate read)

By addressing these anomalies, Twisp's isolation levels ensure correctness and make ledger operations easy to reason about, even under high concurrency.

## Strongly Consistent for All Operations

Twisp ensures **strong consistency** across all database operations, guaranteeing that every transaction reflects an accurate state of the ledger. This robust consistency model removes ambiguity, preventing the risks of stale or partial reads.

{% callout type="note" %}
Customer defined __search__ index documents are replicated in an eventually consistent manner to [OpenSearch](https://opensearch.org/).
{% /callout %}

### Balance Calculations

For high-volume, concurrent entry posting to concurrent-enabled accounts or aggregates, Twisp implements a non-blocking, transactionally consistent, and deterministic balance update process. The following balance types are provided to allow context-specific balance retrieval consistency without compromising correctness.

| **State** | **Description** | **Use Case** |
|-----------|-----------------|--------------|
| **Provisional** | Includes in-flight, uncommitted transactions along with committed changes. | Enforce velocity controls on to-be-committed transactions.   |
| **Prepared** | Combines finalized balances with committed, but not-yet-finalized entries. Guaranteed to provide the exact same results as the final balance. | Systems requiring transactionally consistent balances on concurrent-enabled accounts. |
| **Final** | Reflects the fully committed and finalized state of the account. Finalized balances are cached. | Systems that require the highest performance balance reads and can tolerate slightly stale balance values. |

{% callout type="note" %}
For concurrent-disabled accounts, the `Prepared` balance state is skipped and the balance update process proceeds directly to `Final` from `Provisional`.
{% /callout %}

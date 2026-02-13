---
title: Infrastructure
description: Underneath the accounting core is a purpose-built ledger database and computation engine.
showNext: true
showPrev: false
---

## Architectural Components

There are four primary components to the infrastructure which powers the Twisp Accounting core:

* The [Financial Ledger Database](/infrastructure/ledger-database) is a cloud-native, append-only ledger data store with strong guarantees about the lineage of data.
* The [Computation Engine](/infrastructure/computation-engine) powers massively scalable data operations, protocol workflows, and an extensive financial calculation runtime.
* [Security configurations](/infrastructure/security-and-auth) that provide strong access controls through explict authorization policies.
* The [Local Environment](/infrastructure/local-environment) streamlines development and CICD integration testing.

## Background: Decoupling Storage, Computation, and Security

In most modern systems, data storage is decoupled from computational processes. To manage storage of this data, applications contain logic to replicate storage data structures, apply transactional logic, and then execute remote calls to persist the resulting data in external, decoupled storage systems.

Decoupled storage and computation is ubiquitous in nearly all distributed architectures.

This decoupling of storage and computation makes the process of guaranteeing the integrity of a value calculated in one system, then persisted in another, into a complex distributed systems problem.

For financial ledger systems where we must take great care to guarantee not only the integrity of financial data but also calculations that result from that data, this is a problem.

However, with Twisp, just as referential integrity offers integrity over relations, the financial ledger database offers integrity, full transparency, and lineage of both the underlying data along with all calculations derived from that data.

This property allows us to declaratively define strict, well-defined, and atomic processes that compute ledger calculations, and with those procedures defined, the underlying system transactionally guarantees the accuracy of the resulting calculation data.

Additionally, we treat security and authentication/authorization concerns as first-class citizens with their own dedicated resources. All access controls are determined by a centralized set of security policies. There are two primary reasons for this design:

1. Keeping all security protocols in a single place makes it easy for auditors to review every policy and quickly identify any areas of concern.
2. Storing policies in separate locations runs the risk of introducing conflicts between policy definitions and causing unexpected errors.

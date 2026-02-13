---
title: Computation Engine
description: The computation infrastructure supports protocol workflows and an extensive financial calculation runtime via CEL.
showNext: true
showPrev: true
---

Protocol workflows provide the ability to model sophisticated state-machine-like computations on ledger data and calculations, including the ability to communicate with external systems.

With the workflow engine, it is possible to model sophisticated financial protocols and their interaction with ledgers. Examples of these protocols include ISO 8583, ACH, and external webhook processing.

Finally, Twisp's built-in expression language, which is highlighted later in this document, is populated with an extensive function library. In addition to a traditional standard library, Twisp also includes advanced financial functionality to assist in complex financial systems.

## Workflows as State Machines

A state machine is an abstraction used to design algorithms and protocols. A state machine simply reads a set of inputs and changes to a different state based on those inputs. A state is a description of the status of a system waiting to execute a transition.

Twisp provides a mechanism to define state machines called protocol workflows. The combination of immutable ledgers, the calculation runtime, and protocol workflows allows Twisp to process all manner of complex financial protocols and translate them into ledger entries with full historical lineage and integrity guarantees.

The input for protocol workflows can be triggered by events resulting from ledger operations (create, read, update, delete) or API calls to Twisp (including webhooks from external systems).

One can think of the operational trigger events for protocol workflows with the same mental model of triggers in monolithic relational databases. Workflows in Twisp occur in the context of the same atomic interactive transaction that triggered them. This provides strong transactional and data integrity guarantees for the resulting data.

The output from workflow transitions and the resulting state from the calculation engine is also persisted to ledgers. This allows Twisp to define chains of transactionally guaranteed protocol workflow computation with full lineage. Since each transaction in Twisp has dedicated compute and memory, there is little concern of noisy neighbor problems or node overloading.

## Calculation Runtime

The calculation runtime is the primary mechanism Twisp uses to both populate ledgers and derive useful information from ledger data. It is an expression based language and runtime that allows processing over schema-based input and that results in schema-based output.

It is purpose-built to operate with all manners of persisted ledger data with strong integrity guarantees including a powerful type system, strong type checking, and predictable runtime performance.

Calculation expressions are powered by [Common Expression Language (CEL)](https://github.com/google/cel-spec) along with our financial type system and function library. This provides an expressive and powerful mechanism to derive and populate ledger data. All data resulting from expressions in the calculation runtime have strong integrity guarantees, and result in strictly-typed schema-based output that is purpose built for persisting to ledgers.

The calculation engine supports use cases such as calculating balances, fees, and interest on any dimension; extracting, aggregating, and bucketing ledger data; real-time analytics for business intelligence or user insights; and all manner of user defined functionality powered by the expressive syntax and extensive function library.

## Function Library

The ledger core and calculation runtime plus standard library and protocol workflow provide an extremely powerful set of primitives which lie at the foundation of the Twisp accounting core.

To make these primitives even more powerful for financial applications, Twisp includes a financial type system and extensive function library. This functionality is surfaced within the calculation runtime as types within the strict type system and functions that operate over those types within expressions.

The Twisp financial type system and function library includes sophisticated types such as full-featured money and currency types, along with functions to make complex calculations on those types, such as interest, fees, rounding, balances, and the like. The function library also includes parsers for well known external formats, such as ISO8583, ACH, and many more.
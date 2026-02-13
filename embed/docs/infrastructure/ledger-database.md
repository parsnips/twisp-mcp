---
title: The Twisp Financial Ledger Database
description: The ledger database is an immutable, append only data structure that provides strong gaurantees about the lineage of data.
showNext: true
showPrev: true
---

## Overview: Ledgers à la Twisp

At Twisp, we set out to rethink the underlying technology for financial ledger systems by combining the operational and scaling characteristics of a distributed database with the correctness guarantees offered by relational databases.

The result is the **Twisp Financial Ledger Database** (FLDB). The Twisp ledger database is transactionally guaranteed, secure, and highly scalable; it is tailor-built for financial system use cases. It is designed to easily support the massively multi-tenant workloads that modern financial applications require.

With the advent of serverless and infrastructure-as-code, the ability to define fully functioning systems with high-level code is now a reality. Twisp uses an infrastructure-as-code approach to creating resources and managing access to those resources.

## Comparison to Other Databases

Like traditional relational databases, the Twisp FLDB enforces the structure of data with strict user-defined schemas, indexes, and referential integrity. Users familiar with the tabular model of SQL will find many of the same principles and abstractions apply.

The internal structure of the database is not limited to tabular data, however. Like document databases, ledgers records can contain nested sub-documents as well as collection data types like lists and JSON. Sub-documents still must conform to a strictly defined schema.

In this way, the database supports both modalities of data modeling: tabular and document-based. Because schema migrations are so simple, it is easy to adapt the data model as new requirements emerge.

## Database Features

### Append-Only Immutability

The lineage of data in a financial system is critical. When it comes to monetary values, we need to be able to track with confidence where the money came from, where it went, and how a balance grew from $0 to $1.

At the most basic level of the database structure, we can make these same guarantees about any data that enters the system. We do this through a simple mechanism: immutability.

Immutability means that for an operation that results in new data, the system simply records that new piece of data. All previous data still remains recorded unchanged as well. You can think of this as a log of all changes to a system. Once we have this append-only log, we can then apply strong cryptographic guarantees to the log so the underlying data remains tamper proof and unchanged.

We call this log and strong cryptographic guarantees an immutable ledger of changes to the system. Twisp provides this as a primitive to make use of these same guarantees for their application data.

Therefore, all data are recorded with full lineage in the system. In practice, this means that a `history` node is available on every record in the database to expose the chain of changes to the state of any record. Read more about append-only immutability at  [](/reference/ledger/versions-and-history).

### Strong Data Integrity

A transactional MVCC write-ahead log (WAL) tracks all data that is managed by the system. The data blocks of the WAL are chained together with cryptographic hash functions.

This cryptographically sealed chain allows Twisp to verify the integrity of any data in the system using cryptographic verification.

The system utilizes a digest hash value that represents the full hash chain of the WAL at any committed value plus a Merkle audit proof to cryptographically verify the integrity of any data in the system. This allows users to verify and prove the integrity of any record in the database at any time.

### Index-First Design

All reads to data stored must occur via an index. This eliminates an entire class of problems that can occur in production systems with table scans and unclear query planning, and promotes well-understood data access patterns.

There is no limit the number of indexes on data, as we recognize the tradeoff that more indexes allows for the possibility of write amplification.

The index system is core to Twisp’s infrastructural design, allowing for massive scalability and zero-downtime migrations.

Index features:

- **Strongly consistent**\
  All database operations through indexes guarantee data integrity.
- **Strong partition and sorting controls**\
  The declarative schema supports compound keys and multi-field sorting.
- **Compatible with all fields**\
  Any part of the schema can be used to create the index key, including the JSON and List collection types.
- **Filterable with CEL expressions**\
  Partial indexes can be created to index only those records matching a specified condition.
- **Unlimited definitions**\
  Ledger definitions can include as many indexes as needed.
- **Online, zero-downtime migrations**\
  Indexes can be re-partitioned on the fly using the migration system.

### Index-Based Relations

Using properly constructed indexes, ledgers support all the standard relation types:

- One-to-one
- One-to-many
- Many-to-many

The database implements full referential integrity: records can refer to other records via an index, even across ledgers.

Twisp ensures that a referenced (parent) record cannot be deleted without any referencing (child) records first being deleted.

In addition to preserving referential integrity, the GraphQL schema will automatically support join operations between referenced ledgers.

### Referential Integrity

A powerful feature of relational databases are tools to enforce the integrity and structure of the data they store. Strict schemas require a homogenous shape to all data stored in a table. Rather than depending on an application to enforce a data schema, the database enforces it instead.

With this, the database system can now run migrations and schema changes to reorganize the data into evolving structure for you, with strong guarantees that data is always in a well known shape.

In addition to the schema structure of individual data elements, referential integrity can enforce the relationships between different related pieces of data. The relationships between these pieces of data are often referred to as foreign key (FK) constraints. These constraints provide a mechanism that allows a data element to reference another separate but related piece data via a FK relationship.

Once this relationship is established the system guarantees that one piece of data that is in relation to another will always exist by preventing any changes, whether updates or deletes, that might violate the relationship. The system is maintaining referential integrity.

### Type-Safe Schemas

Within the database, every record conforms to the same schema of typed fields defined within a declarative schema. Fields can be typed to contain a single **scalar** type like strings or signed integers, **collection** types like lists or JSON objects, or they can contain **sub-documents** with their own fields.

This table lists all currently supported types for ledger fields.

| Type        | Description                                                             |
|-------------|-------------------------------------------------------------------------|
| `int`       | 64-bit signed integer                                                   |
| `uint`      | 64-bit unsigned integer                                                 |
| `double`    | 64-bit IEEE floating-point number                                       |
| `bool`      | boolean value true or false                                             |
| `string`    | Strings of Unicode code points                                          |
| `timestamp` | ISO-8601 timestamp with nanosecond precision and full time zone support |
| `bytes`     | Byte sequences (base64 encoded)                                         |
| `uuid`      | 128-bit RFC4122 Universally Unique Identifier                           |
| `json`      | holds arbitrary JSON (use "{}" to structure insert values)              |
| `[<type>]`  | list of scalar values (example: [string])                               |

Note the two collection types `json` and **list** (`[]`). Lists can be used to store variable-length ordered collections of a scalar types. `json` types can store arbitrary JSON documents (up to X bytes)

Twisp combines the flexibility of document databases like MongoDB while maintaining the schematic integrity and operational model of a table-driven relational database like PostgreSQL. In practice, this means that you can create nested document structures while maintaining type safety and predictable data shapes. Sub-documents can be used to store related fields that might otherwise add noise to the base document.

Because they are typed as well, sub-documents are part of the explicitly declared schema and behave in just the same way as fields in the base-level document do. They can even be used in indexes.

### Zero-Downtime Migrations

With declarative schemas we can simply change the schema definition and Twisp will create, validate, and execute changesets to migrate your schema from version to version, all in an online fashion with zero down time.

{% callout type="note" %}

Zero downtime means that existing resources can continue to be accessed and used.  Newly added/removed elements will be in a `DELETE_ONLY` or `WRITE_ONLY` state as they're added to the system and in some cases backfilled with data.

{% /callout %}

The table below shows all supported operations for migrations.

| Migration            | Support                                |
|----------------------|----------------------------------------|
| Add                  | Supported                              |
| Drop                 | Supported                              |
| Add Index            | Supported                              |
| Drop Index           | Supported                              |
| Add Schema Elements  | Supported                              |
| Drop Schema Elements | Supported                              |
| Add References       | Supported                              |
| Drop References      | Supported                              |
| Add Joins            | Supported                              |
| Drop Joins           | Supported                              |
| Add Calculation      | Supported, replays based on `Position` |

### Transaction Layer

The transaction layer provides MVCC interactive transactions for database operations in the system. It allows for strictly serializable transaction isolation levels for the strongest consistency guarantees.

The interface to the transaction layer is provided via transaction scopes. Transaction scopes provide an interface to the transaction layer and a clear mental model for determining how transactions for specific data should interact.

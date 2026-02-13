---
title: Indexes
metaTitle: "Reference: Indexes"
description: Reference for indexes within a Twisp Ledger.
showNext: true
showPrev: true
---

## The Basics

All data retrieval operations in Twisp are executed via indexes. This eliminates issues associated with table scans and dynamic query planning, providing consistent and predictable data access patterns. Indexes in Twisp are designed to support sophisticated application requirements:

- **Transactionally Consistent Operations**: All database interactions via indexes are transactionally consistent, ensuring strong consistency.

- **Enhanced Partition and Sorting Controls**: Twisp’s schema supports compound keys and multi-field sorting, allowing for precise control over how data is partitioned and ordered.

- **Compatibility with All Fields**: Any field within the schema, including those within JSON metadata and List collection types, can be used to create an index. This flexibility allows users to tailor their indexing strategy to the specific needs of their application.

- **Filterable with CEL Expressions**: Partial indexes can be created based on conditions specified in CEL (Common Expression Language) expressions, enabling targeted indexing of records that meet certain criteria.

## Types of Indexes

Multiple index types are supported in Twisp. Each index type supports a unique set of data access patterns. By selecting the appropriate index type, users can optimize query performance and tailor data access to their specific application requirements within the Twisp system.

A number of default indexes are included in the Twisp system. These built-in indexes are managed by Twisp to provide efficient query performance for common financial ledger operations. These indexes guarantee data integrity and are transactionally consistent, making them a reliable choice for most standard query operations.

### Custom Indexes

Custom indexes in Twisp are pivotal for optimizing queries and structuring data access patterns specific to your application's needs. Custom indexes can be created for the following record types:

- `AccountSet`
- `Account`
- `Balance`
- `Entry`
- `TranCode`
- `Transaction`

This capability enables precise control over how data related to these records is indexed and accessed. Custom indexes are transactionally consistent.

### Historical Indexes

Historical indexes are specialized indexes that allow querying across all versions of records' histories. Given Twisp’s immutable, append-only data store, any data change results in a new version of a record. Historical indexes enable sophisticated queries that include past record states, which regular indexes do not cover.

- **Version Tracking:** Indexes every version of a record, enabling retrieval of historical data states.
- **Advanced Query Capabilities:** Supports queries like retrieving account balances at specific points in time or finding record states when particular metadata values were set.

### Search Indexes

Search indexes integrates OpenSearch's powerful search capabilities within the Twisp system, providing enhanced performance for text-based queries and structured data retrieval. Leveraging OpenSearch allows for full-text search, advanced filtering, and robust search ranking capabilities, enabling users to harness large datasets efficiently.

- **Full-Text Search Capabilities**: OpenSearch indexes support search queries across textual data fields, allowing users to perform complex searches that include fuzzy matching, stemming, and tokenization.

- **Rich Filtering Options**: Users can apply filters on various fields, such as numeric ranges, date intervals, or specific terms, making it easier to narrow down search results based on contextual requirements.

## Components of Indexes

### Index Key Fields

Indexes in Twisp can be created using both root-level fields and nested fields within documents. This flexibility allows for comprehensive indexing strategies that can cater to complex data structures.

For example, you can create custom index keys using:

- Root-level fields such as `Account.modified`.
- Nested fields within objects, such as fields within the arbitrary JSON `metadata` object in a record.

This approach ensures the ability to index any data within the system.

### Index Expressions

An index key need not be just a field of the underlying record, but can be a function or expression computed from one or more fields in the record. This feature is useful to obtain fast access to records based on the results of computations.

### Partition Keys

When designing custom indexes, it's crucial to consider partitioning strategies to ensure performance and scalability. Partitioning by account is commonly sufficient, but specific workloads may necessitate alternative approaches. Be mindful of read and write operations per partition to prevent throttling, as the database supports a fixed amount of bandwidth and operations per second for each partition.

### Sort Keys

For a specific partition key, sort keys allow all data that share that partition key to be sorted and retrieved efficiently based on application requirements.

### Unique Constraints

Unique constraints ensure that the data contained in a field, or a group of fields, is unique across all records indexed.

### Partial Indexes

Partial indexes can be created based on conditions specified in CEL (Common Expression Language) expressions, enabling targeted indexing of records that meet certain criteria.

### Asynchronous Indexes

For indexes that may have a hot partition keys prone to throttling, often seen during high-volume data loads, asynchronous indexes can be used to populate the index via a background process. Asynchronous indexes are eventually consistent and unique constraints are not permitted.

### Sharded Indexes

Sharded indexes in Twisp provide a mechanism to enhance write scalability for partitions that experience high write throughput. By splitting a single partition into multiple shards, write operations can be distributed across these shards, thereby increasing the overall write capacity for that partition key.

Each shard within a sharded partition acts like a mini-partition, handling a portion of the write load. This is similar to RAID-style disk striping, where data is divided across multiple disks to improve performance. In the context of Twisp indexes, sharding allows the system to handle more writes per second for a given partition key by parallelizing the write operations across the shards.

However, this increased write capacity comes at the cost of global sort order within the partition. Since the data is distributed across multiple shards, maintaining a strict sort order across the entire partition becomes challenging. Therefore, while sharded indexes are excellent for scenarios requiring high write throughput, they may not be suitable for use cases that rely heavily on sorted data retrieval across the entire partition.

To configure a sharded index, you can specify the `partition_shard_count` parameter when creating or updating an index. This parameter determines the number of shards into which each unique partition key is split. For example, setting `partition_shard_count` to 4 will divide each partition into four shards, potentially allowing for up to four times the write throughput compared to a non-sharded partition.

Sharded indexes are particularly useful in scenarios where a single partition key experiences a high volume of write operations, which could otherwise lead to write throttling or performance bottlenecks. By distributing the write load, sharded indexes help maintain system responsiveness and throughput. However, it’s important to carefully consider the impact on read operations and sort order before implementing sharded indexes.

### OpenSearch Schemas

For search indexes, Twisp provides full control of the underlying OpenSearch schema, including the selection of fields or expressions to be indexed and their respective data types.

## Zero-Downtime Migrations

Twisp's migration infrastructure allows for online creation and modification of indexes with zero downtime. This capability ensures that updates and changes to the data schema can be implemented seamlessly, without disrupting service availability.

This system allows re-partitioning of indexes on the fly, ensuring that your system remains responsive and available even during significant changes.

## Index Operations

- [`Query.schema.index()`]($gql:query:schema.index): Query an Index
- [`Query.schema.indexes()`]($gql:query:schema.indexes): Query Indexes
- [`Mutation.schema.createIndex()`]($gql:mutation:schema.createIndex): Create a Custom Index
- [`Mutation.schema.createHistoricalIndex()`]($gql:mutation:schema.createHistoricalIndex): Create a Historical Index
- [`Mutation.schema.createSearchIndex()`]($gql:mutation:schema.createSearchIndex): Create a Search Index
- [`Mutation.schema.deleteIndex()`]($gql:mutation:schema.deleteIndex): Create an Index
- [`Mutation.schema.updateSearchIndex()`]($gql:mutation:schema.updateSearchIndex): Update a Search Index

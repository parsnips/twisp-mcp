---
title: Journals
metaTitle: "Reference: Journals"
description: Reference for the journal resource within a Twisp Ledger.
showNext: true
showPrev: true
---

## The Basics

Journals are an essential part of accounting systems, as they provide a way to organize and record transactions within separate "books".

{% callout type="note" %}

Put simply, journals help you keep your financial transactions organized in **separate collections**.

{% /callout %}

In Twisp, every ledger instance starts with a **default journal**. You can create additional journals as needed to suit your specific accounting structure.

## Components of Journals

1. **Name**: should be self-descriptive and helps in recognizing the journal's purpose.
2. **Status**: can be either `ACTIVE` or `LOCKED`. Journals with a `LOCKED` status do not allow transactions to be posted to them.
3. **Code**: is an optional unique code for the journal, which can be used as an additional reference.

Journals also have **Created & Updated Timestamps** as well as a **Version & History**.

## Use Cases for Multiple Journals

Having multiple journals can be beneficial in various scenarios.

For example, users may create separate journals for different currencies or product-specific journals.

By using multiple journals, you can more effectively organize transactions and maintain a clear understanding of your financial activities.

## Journal Operations

Use GraphQL queries and mutations to read, create, update, and delete (lock) journals:

- [`Query.journal()`]($gql:query:journal): Get a single journal.
- [`Query.journals()`]($gql:query:journals): Query journals using index filters.
- [`Mutation.createJournal()`]($gql:mutation:createJournal): Create a new journal.
- [`Mutation.updateJournal()`]($gql:mutation:updateJournal): Update select fields for an existing journal.
- [`Mutation.deleteJournal()`]($gql:mutation:deleteJournal): Soft delete (lock) a specified journal.

## Further Reading

To learn the basics of managing journals, see the tutorial on [](/tutorials/working-with-journals).

To review the GraphQL docs for the `Journal` type, see []($gql:object:Journal).

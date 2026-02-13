---
title: Twisp 101 Tutorial
description: Work through the steps in this tutorial to learn how to start building on the Twisp accounting core.
showNext: true
showPrev: false
---

Welcome to the Twisp 101 Tutorial!

## Context

This tutorial is designed to walk you through the foundational features of the Twisp Accounting Core. The example project we'll be using is an imaginary neobank called "Zuzu" which offers checking, savings, and loan products to under-banked customers in the US.

We'll work through the steps needed to **set up a chart of accounts**, orchestrate the types of financial activity by **defining transaction codes**, and **write some ledger entries** by posting transactions.

Most new Twisp customers will have a similar process regardless of the specifics of their product. The fundamental stages are:

1. **Design** the elements which give a structure and meaning to the accounting system.
2. **Test** the design to ensure that it satisfies working requirements.
3. **Deploy** products backed by the Twisp core.
4. **Monitor** for performance, usage, auditing, and reporting.

In this tutorial, we'll focus on the Design and Test stages.

## The Project

_Zuzu_ is a new neobank offering personal banking services: checking, savings, and loans. They work with multiple partner banks, issue debit cards, and provide online and mobile banking applications.

Building on Twisp's accounting core, Zuzu has the following project requirements.

- A chart of accounts to model both internal company accounts as well as customer accounts, including:
  - Checking accounts
  - Savings accounts
  - Loans issued
- A double-entry ledger for recording financial activity, including common transactions like:
  - Deposits & withdrawals
  - Direct transfers between customers
  - Fees charged for a variety of reasons

In the interest of getting something working, we'll start simple and just implement the accounting components needed to support the **checking** product and its primary transactions.

## Step 0: Establish API Connection

Twisp is an API-first product. Before we can do anything, we need to connect to the API.

### Connect to GraphQL API

When you are invited to Twisp, you will receive an account to access the Twisp Console. The console provides a [GraphiQL](https://github.com/graphql/graphiql) interface for interacting with the GraphQL API for your provisioned accounting core.

This is the quickest and easiest way to get connected. Once you have proper auth credentials, you can connect to the GraphQL API from any client.

### Inspect the Schema

The GraphiQL interface provides built-in docs, so you can explore the schema and read documentation about queries, mutations, and types.

Let's confirm that you're connection is working. Run an introspection query directly:

```graphql {% title="Introspection" %}
query {
  __schema {
    types {
      name
      kind
      description
    }
  }
}
```

You should get a large response with all the various types like `Account`, `Entry`, etc.

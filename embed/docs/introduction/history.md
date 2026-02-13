---
title: A Short History of Twisp
description: Why we did it.
showNext: true
showPrev: true
---

Monolithic relational databases are the traditional foundation of financial core ledger systems. Nevertheless, the process of building and operating mission-critical financial ledgers on these databases is a journey fraught with engineering challenges.

In addition, the databases underlying these systems are often difficult to operate at scale and a poor choice for multi-tenant systems. Distributed database technologies have proven themselves to solve these operational issues.

However, **the design choices that allow these systems to achieve horizontal scale make them an even more complicated fit for financial ledgers.**

At Twisp, we set out to rethink the underlying technology for financial ledger systems by combining the operational and scaling characteristics of a distributed database with the correctness guarantees offered by relational databases. The result is the Twisp Financial Ledger Database (FLDB).

The Twisp FLDB is a transactionally guaranteed, secure, and highly scalable financial ledger database built for financial system use cases. It is designed to easily support massively multi-tenant workloads.

Additionally, we set out to fully leverage the modern cloud. With the advent of serverless and infrastructure-as-code, the ability to define fully functioning systems with high-level code is now a reality. Twisp systems are repeatable, autoscaling, fully managed, pay-per-use, and a boon to developer productivity.

{% callout type="tip" %}

Read more about the infrastructure model in the [Infrastructure Docs](/infrastructure).

{% /callout %}

At its heart, Twisp is an infrastructure company. We built the FLDB because we wanted it to exist in the world. As we released this database into the wild, we found out that many companies didn't just need a ledger database—they needed an [accounting core](/accounting-core).

Lots of companies have re-invented a ledger and a system for accounting. Many of these systems suffer from problems of scale, complexity, and a faulty data model which can make seemingly simple tasks like calculating balances into confusing, costly chores. This is why you [shouldn't build your own ledger](/introduction/no-diy-ledger).
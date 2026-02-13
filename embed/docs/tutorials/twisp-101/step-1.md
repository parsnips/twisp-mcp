---
title: "Step 1: Create a Primary Journal"
description: Before we begin designing accounts, we need to define a journal within which to organize transactions.
metaTitle: "Twisp 101 Tutorial: Step 1"
metaDescription: Work through the steps in this tutorial to learn how to start building on the Twisp accounting core.
showNext: true
showPrev: true
---

In most cases, a single [Journal]($gql:object:Journal) is enough. For Zuzu, we'll create a single "General Ledger" journal which will record all transactions.

{% partial file="graphql/examples/twisp101/001_CreateGLJournal.md" /%}

{% callout %}
Note that Twisp ledgers come with a pre-built "Default" journal, which is ideal for single-journal ledgers. For the purposes of illustration, we'll be using a custom Journal in this tutorial.
{% /callout %}

Now that we have a journal, we can start building the components needed to support the core use cases. Let's start by modeling **deposits** and **withdrawals**.

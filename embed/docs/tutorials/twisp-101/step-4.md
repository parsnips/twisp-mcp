---
title: "Step 4: Organize Accounts for Balance Rollups"
description: Use account sets to group and organize accounts into a chart, taking advantage of the built-in balance aggregations.
metaTitle: "Twisp 101 Tutorial: Step 4"
metaDescription: Work through the steps in this tutorial to learn how to start building on the Twisp accounting core.
showNext: false
showPrev: true
---

With [AccountSets]($gql:object:AccountSet), we can collect related accounts to provide an easy interface into summary balances and queries into the entries.

We've already created checking accounts for each customer, but Zuzu also needs a way to summarize _all_ customer accounts so that we can see the total balance.

To accomplish this, we'll create an account set called "Customers" and add the customer accounts to it.

## Create customers account set

{% partial file="graphql/examples/twisp101/010_CreateCustomersAccountSet.md" /%}

Now that we have an account set, let's add the two customer accounts to it:

{% partial file="graphql/examples/twisp101/011_AddCustomersToSet.md" /%}

Now that we have posted several transactions and created an account set, we can look at balances and interrogate their history to see how account balances change with each activity posted to that account.

## Review balances & interrogate history

Let's start by querying for the balance of the "Customers" set, as well as the balance of each member account.

{% partial file="graphql/examples/twisp101/013_GetCustomersBalances.md" /%}

As expected, the account set's balance is always equal to the sum of its member's balances.

Because records in Twisp are append-only, we can review the history of any record to see how its state changed over time.

Let's query the balance history of Ernie's account and compare it to the entries to see how it changed in response to transactions posted.

{% partial file="graphql/examples/twisp101/014_GetErnieBalanceHistory.md" /%}

## Conclusion

That finalizes the tutorial! Let's recap what we did to customize and test an accounting core for the imaginary neobank Zuzu:

- ✅ Modeled accounts for customers, assets, and revenue
- ✅ Designed tran codes for bank transfers as well as ACH credits and debits
- ✅ Organized customer accounts into a set and ran queries against it
- ✅ Interrogated the balance history and entries for an account

We hope that this has been useful for you to understand how Twisp works and what you can build with it. Obviously, this is an oversimplified example. Your product will certainly be more complex (and interesting)!

{% callout %}
We'd love to talk with you about your project. If you're interested, please [get in touch](https://www.twisp.com/?modal=Get+in+touch).
{% /callout %}

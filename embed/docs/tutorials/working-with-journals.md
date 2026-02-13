---
title: Working with Journals
metaTitle: "Tutorial: Working with Journals"
description: |
  In this tutorial, we will learn how to add additional journals to handle a multi-journal ledger system.
showNext: true
showPrev: true
---

Journals are a fundamental tool in accounting, used to record transactions. In this tutorial, we will cover how to create a new journal, modify an existing one, query a journal, and lock a journal to prevent posting.

Every ledger in Twisp starts with a **default journal**, but you can create as many additional journals as needed to suit your particular accounting structure.

{% callout type="task" %}
- Create a new journal using the `createJournal` mutation
- Update an existing journal using the `updateJournal` mutation
- Query a journal with the `journal` query
- Delete (lock) a journal using the `deleteJournal` mutation
{% /callout %}

---

{% partial file="tutorials_getting_started.md" /%}

## Create a new journal

You can create a new journal using the `createJournal` mutation.

{% partial file="graphql/reference/Mutation/001_createJournal.md" variables={title: "Create a journal"} /%}

This will return the ID, name, description, and status of the newly created journal. Customize the input to suit your needs.

## Modify an existing journal

To modify an existing journal, you can use the `updateJournal` mutation. Let's update the description of the newly created journal:

{% partial file="graphql/reference/Mutation/009_updateJournal.md" variables={title: "Update a journal"} /%}

The response includes the UUID and the updated description of the journal, as well as the history of the changes made to the journal.

## Query a journal

To query a journal, you can use the `journal` query and provide the ID to fetch.

{% partial file="graphql/reference/Query/008_journal.md" variables={title: "Read a journal"} /%}

This will return the name, description, status, and version of the journal.

## Lock a journal to prevent posting

Like other resources in Twisp, "deleting" a journal does not actually remove it from the database, but instead prevents it from being used by setting its status to `LOCKED`.

To lock a journal to prevent posting, you can use the `deleteJournal` mutation. Let's lock our General Ledger journal:

{% partial file="graphql/reference/Mutation/021_deleteJournal.md" variables={title: "Lock journal"} /%}

## Conclusion

Journals are a fundamental component of accounting, used to store collections of transactions. In this tutorial, we covered how to create a new journal, modify an existing one, query a journal, and lock a journal. These basic operations can be used to manage journals effectively and efficiently in your multi-journal accounting workflow.

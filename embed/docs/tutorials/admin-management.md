---
title: "Core Admin: Managing Tenants, Users, and Groups"
metaTitle: "Tutorial: Managing Tenants, Users, and Groups"
description: |
  In this tutorial, we'll cover the administrative tasks involved in a Twisp ledger: managing tenants, users, and groups.
showNext: false
showPrev: true
---

## Introduction to Tenants, Users, and Groups

In any organization, especially in the context of accounting and financial systems, managing access control and permissions is crucial to ensure data security and prevent unauthorized access. Twisp provides a robust system to manage access control through [Tenants]($gql:object:Tenant), [Users]($gql:object:User), and [Groups]($gql:object:Group).

### Tenants

A [Tenant]($gql:object:Tenant) in Twisp represents an environment within an organization, typically associated with a specific application, service, or set of resources. Tenants contain isolated ledgers, each deployed to a specific region. They play a vital role in isolating data and configurations between different environments. Each tenant is uniquely identified by an `accountId`, which in combination with an AWS region, is used to calculate the database tenant for data isolation purposes.

By isolating the ledgers and associated resources, Tenants ensure that each environment remains independent, preventing accidental or unauthorized access to sensitive data from other environments.

### Users

[Users]($gql:object:User) are human members within an organization who interact with the Twisp accounting core. Each user is uniquely identified by their email address. Users can belong to multiple [Groups]($gql:object:Group), which define their permissions within the organization based on the associated policies of each group. The effective permissions of a user are determined by the combined set of policies from all their groups.

### Groups

[Groups]($gql:object:Group) are a logical grouping of users within an organization. They are used to manage access control and permissions for users. Each Group can have one or more associated policies that define the allowed actions for its member users. Users can belong to multiple Groups, and their permissions are determined by the combined set of policies from all their groups.

## Working with Users and Groups

### Managing Users

To create a new User, you can use the `admin.createUser` mutation, providing the necessary input such as the user's unique ID (UUID), email address, and the group IDs the user should belong to.

{% partial file="graphql/reference/Mutation.admin/013_createUser.md" variables={title: "Example"} /%}

{% callout %}
Newly created users will automatically receive an invitation email with a link to the Twisp console.
{% /callout %}

To update an existing User's details, such as their email address or group memberships, you can use the `admin.updateUser` mutation with the corresponding input.

### Managing Groups

To create a new Group, use the `admin.createGroup` mutation, providing the necessary input such as the group's unique ID (UUID), name, description, and policy.

{% partial file="graphql/reference/Mutation.admin/007_createGroup.md" variables={title: "Example"} /%}

To update an existing Group's details, such as the name, description, or policy, use the `admin.updateGroup` mutation with the appropriate input.

{% partial file="graphql/reference/Mutation.admin/009_updateGroup.md" variables={title: "Example"} /%}

### Deleting Users and Groups

To delete an existing user or group, use the `admin.deleteUser` or `admin.deleteGroup` mutation. You will need to provide the email address of the user you want to delete, or the name of the group.

For example, this mutation deletes the user with the specified email address and returns the deleted user's email:

{% partial file="graphql/reference/Mutation.admin/017_deleteUser.md" variables={title: "Example"} /%}

## Creating a Tenant

A Tenant represents an environment within an organization, used for isolating data and configurations for specific applications or services. By creating a Tenant, you can ensure that your different environments have separate ledgers and resources, improving organization and security.

### Using the `admin.createTenant` mutation

To create a new [Tenant]($gql:object:Tenant), you will need to use the `admin.createTenant` mutation, providing the necessary details such as `id`, `accountId`, `name`, and `description`. The `accountId` is especially important, as it is used in combination with an AWS region to isolate the tenant's data.

Here's an example of a GraphQL mutation to create a new tenant:

{% partial file="graphql/reference/Mutation.admin/001_createTenant.md" variables={title: "Example"} /%}

## Assigning Permissions and Access Control

Permissions are determined by the policies associated with each [Group]($gql:object:Group) a [User]($gql:object:User) belongs to. The effective permissions for a user are calculated based on the combined set of policies from all their groups.

For example, if a user belongs to two groups with different policies, their effective permissions will be the result of evaluating both policies together.

To assign groups and permissions to users, you can use the `admin.updateUser` mutation by providing the user's ID and the desired group IDs. This mutation allows you to update the user's group membership, determining their effective permissions within the organization. Here's an example:

{% partial file="graphql/reference/Mutation.admin/015_updateUser.md" variables={title: "Example"} /%}

In this example, the user's group is changed because a new group ID is provided.

### Verifying Effective Permissions for Users

To ensure that users have the correct permissions, it is essential to verify their effective permissions, which are determined by the combined set of policies from all their associated groups.

To check a user's effective permissions, review the policies associated with each group the user belongs to, taking note of any `ALLOW` or `DENY` effects on actions and resources. The user must have at least one `ALLOW` policy for each desired action and resource but will be blocked by any `DENY` policy on the same action and resource.

### Monitoring and Adjusting Access as Needed

It is crucial to regularly monitor and adjust user access to maintain a secure environment. You can use the `admin` queries to retrieve details about users, groups, and their permissions within the organization. For example, to get a list of users and their associated groups, use the following query:

{% partial file="graphql/reference/Query.admin/004_users.md" variables={title: "Example"} /%}

## Conclusion

In this tutorial, we have explored the concepts and functionalities of [Tenants]($gql:object:Tenant), [Users]($gql:object:User), and [Groups]($gql:object:Group) within the Twisp Accounting Core. We have learned the importance of managing access control and permissions to maintain a secure and well-organized environment.

By following the steps outlined in this tutorial, you should now have a solid understanding of how to manage Tenants, Users, and Groups effectively within Twisp.

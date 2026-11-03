---
title: Manage billing
linkTitle: Manage
description: >-
  Review Billing in Docker Home: plans, payment methods, invoices, usage,
  costs, credits, and billing details, and what you can do on each page.
keywords: billing, docker home, billing overview, manage plans, payment
  methods, invoices, usage, costs, credits, billing settings, subscription
  totals, taxes, docker subscription
weight: 20
grid:
  - title: Manage plans
    description: Add a new plan, upgrade an active plan, or cancel auto-renewal.
    icon: shopping-cart
    link: /subscription-billing/manage/plans/
  - title: Add or update a payment method
    description: >-
      Add or update a payment method for your individual account or
      organization.
    icon: credit-card
    link: /subscription-billing/manage/payment-method/
  - title: Update billing information
    description: >-
      Update billing information for your individual account or
      organization.
    icon: pencil-square
    link: /subscription-billing/manage/details/
  - title: Invoices
    description: View billing history and download past invoices.
    icon: receipt-percent
    link: /subscription-billing/manage/history/
  - title: 3D Secure authentication
    description: >-
      Learn how 3DS works and how to troubleshoot verification issues.
    icon: wallet
    link: /subscription-billing/manage/3d-secure/
  - title: Taxes
    description: Submit a US tax exemption certificate or add a VAT number.
    icon: document-text
    link: /subscription-billing/manage/tax-certificate/
---

You manage Docker plans, payment methods, invoices, usage, costs, credits, and
billing details from **Billing** in Docker Home. Billing applies to the
account you select, either your individual account or an organization.

For an organization, you must be an organization owner, or have a
[custom role](/manuals/security/roles-and-permissions/custom-roles/permissions-reference.md#billing)
with the **View billing** or **Manage billing** permission. Without it, Docker
Home shows **No access to billing**.

## Overview

**Overview** is the Billing home. It shows the default payment method, a tile
for each active plan, and your subscription totals with the next charge date.
The page also shows:

* **Monthly subscriptions** when you have a monthly plan
* **Annual subscriptions** when you have an annual plan
* An alert for a
  [failed payment](/manuals/subscription-billing/manage/payment-method.md#failed-payments)
  or a
  [bank account waiting for verification](/manuals/subscription-billing/manage/payment-method.md#verify-a-bank-account)

> [!TIP]
>
> If Docker sales manages your subscription, the sidebar shows **Overview**
> only. Contact your Docker sales representative to renew, make changes, or
> request invoices.

## Invoices

**Invoices** lists every invoice with its date, number, description, status,
and amount. You can open an invoice or download a PDF when your Docker account
has a verified email address.

## Usage

**Usage** charts usage-based charges across billing periods. You can filter
by period, product, and meter, and export a CSV.

Plans that report usage appear here, including Docker Build Cloud minutes,
Testcontainers Cloud minutes, and Docker Agentic Platform compute. With no
usage-based plan, the page shows **No usage-based products**. See
[Plans](/manuals/subscription-billing/plans/_index.md#choose-a-plan)
for which plans are metered.

## Costs

**Costs** shows charges for one billing period, grouped by the product each
charge is for. The page labels that column **Resource**. You can switch
billing periods and export a CSV. Costs also shows:

* Period status
* Amounts excluding tax
* **Cost to date** for the current period
* The subtotal for a closed period
* Credits, when the account has them

An invoice total can differ, because the invoice adds tax and any credits.

## Credits

**Credits** shows the credit balances for your products. Open a balance to
see each grant and its remaining amount. Docker applies these balances to
your costs at the end of the billing period.

## Settings

**Settings** shows the default payment method, billing address, tax ID, and
billing email. You can update the payment method, address, or billing email
on this page.

## Next steps

{{< grid >}}

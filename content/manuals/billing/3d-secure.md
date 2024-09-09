---
title: 3D Secure authentication
description: Learn about 3D Secure support for Docker billing.
keywords: billing, renewal, payments, subscriptions
---

> [!NOTE]
>
> [Docker Core subscription](/billing/core-billing/get-started-core/) payments support 3D secure authentication.

3D Secure (3DS) authentication incorporates an additional security layer for credit card transactions. If you’re making payments for your Docker billing in a region that requires 3DS, or using a payment method that requires 3DS, you’ll need to verify your identity to complete any transactions. The method used to verify your identity varies depending on your banking institution.

The following transactions will use 3DS authentication if your payment method requires it.

- Starting a [new paid subscription](/billing/core-billing/get-started-core/)
- Changing your [billing cycle](/billing/core-billing/cycle/) from monthly to annual
- [Upgrading your subscription](/subscription/core-subscription/upgrade/)
- [Adding seats](/subscription/core-subscription/add-seats/) to an existing subscription

## Troubleshooting

If you encounter errors completing payments due to 3DS, you can troubleshoot in the following ways.

1. Retry your transaction and verification of your identity.
2. Contact your bank to determine any errors on their end.
3. Try a different payment method that doesn’t require 3DS.

> [!TIP]
>
> Make sure you allow third-party scripts in your browser and that any ad blocker you may use is disabled when attempting to complete payments.

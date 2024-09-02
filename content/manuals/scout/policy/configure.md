---
title: Configure policies
description: Learn how to configure or disable the default policies in Docker Scout
keywords: scout, policy, configure, disable, enable, parametrize, thresholds
---

Some of the existing policies are configurable. This means that you can clone
an existing policy and create new, custom policies with your own configuration.
You can also disable a policy altogether, if a policy doesn't quite match your
needs.

## Configure a policy

To change the configuration of a policy, you must clone one of the existing
default policies, and then save your configuration as a new policy. You can
edit the display name and description of the new policy to help distinguish
it from the default policy it's based on.

The available configuration parameters for a policy depends on the default
policy you used as a base for your custom policy.

To configure a policy:

1. Go to the [Policies page](https://scout.docker.com/reports/policy) in the Docker Scout Dashboard.
2. Select the policy you want to configure.
3. Select **View policy details** to open the policy side panel.

   If this button is grayed out, then the selected policy doesn't have any
   configuration parameters.

4. In the side panel, select **Copy to customize** to open the policy configuration page.
5. Update the policy parameters.
6. Save the changes:

   - Select **Save and enable** to commit the changes and enable the policy for
     your current organization.
   - Select **Save policy** to save the policy configuration without enabling
     it.

## Disable a policy

When you disable a policy, evaluation results for that policy are hidden, and
no longer appear in the Docker Scout Dashboard or in the CLI. Historic
evaluation results aren't deleted if you disable a policy, so if you change
your mind and re-enable a policy later, results from earlier evaluations will
still be available.

To disable a policy:

1. Go to the [Policies page](https://scout.docker.com/reports/policy) in the Docker Scout Dashboard.
2. Select the policy you want to disable.
3. Select the **Disable** button.

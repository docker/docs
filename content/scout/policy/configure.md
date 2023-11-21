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
policy you used as a base for your custom policy. The following table lists the
default policies that you can configure, and the available configuration
parameters that you can use to create a custom policy.

| Default policy                            | Configuration parameters |
| ----------------------------------------- | ------------------------ |
| All critical vulnerabilities              | Severities               |
| Copyleft licenses                         | License names            |
| Fixable critical and high vulnerabilities | Severities, age          |
| High-profile vulnerabilities              | CVEs                     |

To configure a policy:

1. Go to the [Docker Scout Dashboard](https://scout.docker.com/).
2. Go to the **Policies** section.
3. Select the policy you want to configure.
4. Select the **View configuration** button to open the policy configuration.

   If the button is disabled, the selected policy doesn't have any
   configuration parameters.

5. Select the **Edit policy** button. This prompts you to create a clone of the
   default policy.
6. Select **Copy and edit policy** to create a clone of the default policy.
7. Update the policy parameters.
8. Save the changes:

   - Select **Save and enable** to commit the changes and enable the policy for
     your current organization.
   - Select **Save changes** to save the policy configuration without enabling
     it.

## Disable a policy

When you disable a policy, evaluation results for that policy are hidden, and
no longer appear in the Docker Scout Dashboard or in the CLI. Historic
evaluation results aren't deleted if you disable a policy, so if you change
your mind and re-enable a policy later, results from earlier evaluations will
still be available.

To disable a policy:

1. Go to the [Docker Scout Dashboard](https://scout.docker.com/).
2. Go to the **Policies** section.
3. Select the policy you want to disable.
4. Select **Disable policy**.

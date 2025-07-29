---
title: Deploy with Intune
description: Use Intune, Microsoft's cloud-based device management tool, to deploy Docker Desktop
keywords: microsoft, windows, docker desktop, deploy, mdm, enterprise, administrator, mac, pkg, dmg
tags: [admin]
weight: 40
aliases:
 - /desktop/install/msi/use-intune/
 - /desktop/setup/install/msi/use-intune/
 - /desktop/setup/install/enterprise-deployment/use-intune/
---

{{< summary-bar feature_name="Intune" >}}

Learn how to deploy Docker Desktop on Windows and macOS devices using Microsoft Intune. It covers app creation, installer configuration, and assignment to users or devices.

{{< tabs >}}
{{< tab name="Windows" >}}

1. Sign in to your Intune admin center.
2. Add a new app. Select **Apps**, then **Windows**, then **Add**.
3. For the app type, select **Windows app (Win32)**
4. Select the `intunewin` package. 
5. Fill in the required details, such as the description, publisher, or app version and then select **Next**. 
6. Optional: On the **Program** tab, you can update the **Install command** field to suit your needs. The field is pre-populated with `msiexec /i "DockerDesktop.msi" /qn`. See the [Common installation scenarios](msi-install-and-configure.md) for examples on the changes you can make. 

   > [!TIP]
   >
   > It's recommended you configure the Intune deployment to schedule a reboot of the machine on successful installs.
   >
   > This is because the Docker Desktop installer installs Windows features depending on your engine selection and also updates the membership of the `docker-users` local group.
   >
   > You may also want to set Intune to determine behaviour based on return codes and watch for a return code of `3010`. Return code 3010 means the installation succeeded but a reboot is required.

7. Complete the remaining tabs, then review and create the app. 

{{< /tab >}}
{{< tab name="Mac" >}}

First, upload the package:

1. Sign in to your Intune admin center.
2. Add a new app. Select **Apps**, then **macOS**, then **Add**.
3. Select **Line-of-business app** and then **Select**.
4. Upload the `Docker.pkg` file and fill in the required details.

Next, assign the app:

1. Once the app is added, navigate to **Assignments** in Intune.
2. Select **Add group** and choose the user or device groups you want to assign the app to.
3. Select **Save**.

{{< /tab >}}
{{< /tabs >}}

## Additional resources

- [Explore the FAQs](faq.md).
- Learn how to [enforce sign-in](/manuals/enterprise/security/enforce-sign-in/_index.md) for your users.
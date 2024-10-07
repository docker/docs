---
title: Communication and information gathering
description: Gather your company's requirements from key stakeholders and communicate to your developers.
weight: 10
---

## Step one: Communicate with your developers and IT teams

### Docker user communication

You may already have Docker Desktop users within your company, and some steps in this process may affect how they interact with the platform. It's highly recommended to communicate early with users, informing them that as part of the subscription onboarding, they will be upgraded to a supported version of Docker Desktop. 

Additionally, communicate that settings will be reviewed to optimize productivity, and users will be required to sign in to the company’s Docker organization using their business email to fully utilize the subscription benefits.

### MDM team communication

Device management solutions, such as Intune and Jamf, are commonly used for software distribution across enterprises, typically managed by a dedicated MDM team. We recommend engaging with this team early in the process to understand their requirements and the lead time for deploying changes.

Several key setup steps in this guide require the use of JSON files, registry keys, or .plist files that need to be distributed to developer machines. It’s a best practice to use MDM tools for deploying these configuration files and ensuring their integrity is preserved.

## Step two: Identify Docker organizations

Some companies may have more than one [Docker organization](/manuals/admin/organization/_index.md) created.  These organizations may have been created for specific purposes, or may not be needed anymore.  If you suspect your company has more than one Docker organization, it's recommended you survey your teams to see if they have their own organizations. You can also contact your Docker Customer Success representative to get a list of organizations with users whose emails match your domain name.

## Step three: Gather requirements

### Baseline configuration

Through [Settings Management](/manuals/security/for-admins/hardened-desktop/settings-management/_index.md), Docker offers a significant number of configuration parameters that can be preset.

The Docker organization owner and the development lead should review the settings to determine which of those settings to configure to create the company’s baseline configuration. You should also discuss [enforcing sign-in](/manuals/security/for-admins/enforce-sign-in/_index.md) for your Docker Desktop users and whether you want to take advantage of the free trials of other Docker products such as [Docker Scout](/manuals/scout/_index.md), which is included in the subscription. 

{{< accordion title="Baseline settings to review" >}}

| Setting             | OS Requirements | Description     |
|---------------------|-----------------|-----------------|
| `proxy`               |                 |  This setting configures the proxy used by Docker Desktop to access the internet. The proxy can be set manually or get its value from the system.|
| `wslEngineEnabled`    | Windows only    | This setting specifies whether the user should use WSL 2 or HyperV for the VM for Windows installations.|
| `kubernetes`          |                 | Docker Desktop offers a Kubernetes single-node cluster for Kubernetes deployments locally. This setting controls whether it is started when Docker Desktop starts, and its configuration.|
| `analyticsEnabled`   |                 | Docker allows users to opt out of sending usage data to Docker. The usage data feeds what admins are able to see about Docker Desktop usage, so it is highly recommended to enable and lock this setting.|
| `useVirtualizationFrameworkVirtioFS`| macOS only      | VirtioFS is the newer higher performance file sharing framework for MacOS. It takes precedence over the older frameworks if it is enabled.|
| `useVirtualizationFrameworkRosetta` | macOS only      | Rosetta is the Apple emulator for x86 chipsets. This setting allows Docker Desktop to use Rosetta when running containers built for the x86 chipset.|
| `allowExperimentalFeatures`  |           | Docker Desktop versions often contain experimental features for trial and feedback. If this setting is set to false, experimental features are disabled.|
| `allowBetaFeatures`    |                 | Docker Desktop versions often contain beta features for trial and feedback. If this setting is set to false, beta features are disabled.|
| `configurationFileVersion`      |             | Specifies the version of the configuration file format.|
| `dockerDaemonOptions` - Linux Containers |             | This setting overrides the options in the Docker Engine config file. For details, see the [Docker Engine reference](/reference/cli/dockerd.md#daemon-configuration-file). Note that for added security, a few of the config attributes may be overridden when Enhanced Container Isolation is enabled. |
| `vpnkitCIDR`    |                 | Overrides the network range used for vpnkit DHCP/DNS for `*.docker.internal` |
| `dockerDaemonOptions` - Windows Containers | Windows only    | This setting overrides the options in the daemon config file. For details, see the [Docker Engine reference](/reference/cli/dockerd.md#daemon-configuration-file).|
| `extensionsEnabled`    |                 | Docker extensions are third-party add-ons for Docker Desktop. This setting affects if they are allowed.|
| `useGrpcfuse`     | macOS only      | If the value is set to true, gRPC Fuse is set as the file sharing mechanism. |
| `displayedOnboarding` |                 | There is an onboarding survey that displays when Docker Desktop is installed and opened for the first time. This setting can disable the survey.|

{{< /accordion >}}

### Security configuration

Docker also offers a number of security related features, again through [Settings Management](/manuals/security/for-admins/hardened-desktop/settings-management/_index.md), that can be preset. The infosec representative, Docker organization owner, and the development lead should review those features to determine what should be enabled to meet your company’s security requirements.

{{< accordion title="Security settings to review" >}}

| Setting    | OS Requirements | Description  |
|------------|-----------------|---------------|
| Enhanced Container Isolation     |                 | When this setting is enabled, Docker Desktop runs all containers as unprivileged, via the Linux user-namespace, and prevents them from modifying sensitive configurations inside the Docker Desktop VM, and uses other advanced techniques to isolate them. For more information, see [Enhanced Container Isolation](/manuals/security/for-admins/hardened-desktop/enhanced-container-isolation/_index.md). |
| Registry Access Management |                 | This parameter restricts the registries that `docker pull` and `docker push` commands can access. Note: This is not an endpoint security solution, but a guardrail for users working within company guidelines. For more information, see [Registry Access Management](/manuals/security/for-admins/hardened-desktop/registry-access-management.md).|
| Image Access Management |                 | This parameter restricts the categories of images accessible within Docker Hub. Note: This is not an endpoint security solution; it's a guardrail for users working within company guidelines. For more information, see [Image Access Management](/manuals/security/for-admins/hardened-desktop/image-access-management.md).|
| Scout  |                 | Settings related to how Scout creates SBOMs (Software Bill of Materials) and indexes vulnerabilities for images.|
| `exposeDockerAPIOnTCP2375`         | Windows only    | Exposes the Docker API on a specified port. If the value is set to true, the Docker API is exposed on port `2375`. This is unauthenticated and should only be enabled if protected by suitable firewall rules.|
| `windowsDockerdPort`               | Windows only    | Exposes Docker Desktop's internal proxy locally on this port for the Windows Docker daemon to connect to. It is available for Windows containers only. |
| `filesharingAllowedDirectories`    |                 | Specify which paths on the developer host machine or network your users can add container file shares to.|
| `enableKerberosNtlm`               |                 | When set to true, Kerberos and NTLM authentication is enabled. Default is false. Available in Docker Desktop version 4.32 and later.|
| `containersProxy`          |                 | Allows you to create air-gapped containers. For more information, see [Air-Gapped Containers](/manuals/security/for-admins/hardened-desktop/air-gapped-containers.md).|
| `blockDockerLoad`    |                 | When this setting is enabled, users can no longer run the `docker load` command and will receive an error if they try.|
| `disableUpdate`   |                 | Users get notifications about new Docker Desktop versions. Enabling this setting removes those notifications. Helpful if corporate IT manages Docker Desktop version updates for users.|

{{< /accordion >}}

## Optional step four: Meet with the Docker Implementation team

The Docker Implementation team can help you step through setting up your organization, configuring SSO, enforcing sign in, and configuring Docker.  You can reach out to set up a meeting by emailing successteam@docker.com.

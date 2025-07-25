---
title: Desktop setup templates
description:
keywords:
params:
  sidebar:
    group: Enterprise
weight: 10
---

This page contains pre-configured Docker Desktop settings templates for
common enterprise deployment scenarios. Each template includes configurations
for Windows deployment (MSI), macOS deployment (PKG), and JSON file deployment
methods.

> [!WARNING]
>
> These templates are suggested starting points for different
environments. They should be reviewed, tested, and customized by your security
team before production deployment.

## How to use these templates

These templates provide starting configurations for common Docker
Desktop deployment scenarios. You can customize each template to meet
your specific secuirty and operational requirements.

### Deployment workflow

1. Choose your template:
    - High-security hardened: For regulated industries with strict
    compliance requirements
    - Standard enterprise with proxy: For corporate networks with proxy
    requirements
    - Windows deployment: For teams primarly developing on Windows
    with WSL 2
    - Air-gapped offline: For completely isolated environments
    without internet access
1. Customize the configuration:
    - Review each setting and adjust based on your security policies
    - Update proxy settings, registry URLs, and allowed directories
    - Test thoroughly in a non-production environment first
1. Deploy using your preferred method:
    - Windows (MSI): Use installer flags with `--admin-settings` parameter
    - macOS (PKG): Deploy settings file before or after installation
    - Manually deploy the JSON file to the appropriate directory
1. Verify and monitor settings:
    - Check Docker Desktop settings UI to confirm locked values, or use
    the following CLI command:
    - Use Desktop settings reporting for organization-wide compliance
    monitoring

## High-security hardened environment

Use this template for highly regulated industries or zero-trust environments
requiring maximum security controls. Key configurations include:

- Enhanced Container Isolation (ECI) enabled and locked
- All telemetry and analytics disabled
- Docker Extensions, beta features, and AI features disabled
- Docker load command blocked
- Desktop terminal access disabled
- Updates disabled (manual patching only)
- Strict networking controls

{{< tabs >}}
{{< tab name="Windows deployment (MSI)" >}}

```TBD```

{{< /tab >}}
{{< tab name="macOS deployment (PKG)" >}}

```TBD```

{{< /tab >}}
{{< tab name="JSON configuration file" >}}

```TBD```

{{< /tab >}}
{{< /tabs >}}

## Standard enterprise environment with proxy

Use this template for corporate environments with proxy requirements and
moderate security needs. Key configurations in this template include:

- System proxy configuration enabled
- ECI enabled but allows specific trusted images
- Analytics enabled for usage monitoring
- Extensions disabled
- Beta and AI features disabled
- Standard networking with proxy support

{{< tabs >}}
{{< tab name="Windows deployment (MSI)" >}}

```TBD```

{{< /tab >}}
{{< tab name="macOS deployment (PKG)" >}}

```TBD```

{{< /tab >}}
{{< tab name="JSON configuration file" >}}

```TBD```

{{< /tab >}}
{{< /tabs >}}

## Windows-optimized development environment

Use this template for Windows-focused development teams requiring WSL 2
integration and Windows container support. Key configurations in this template
include:

- WASL 2 engine enabled and locked
- Windows container support configured
- VirtioFS for improved file sharing performance
- Kubernetes disabled
- Development-friendly file sharing paths
- Updates enabled but controlled

{{< tabs >}}
{{< tab name="Windows deployment (MSI)" >}}

```TBD```

{{< /tab >}}
{{< tab name="macOS deployment (PKG)" >}}

```TBD```

{{< /tab >}}
{{< tab name="JSON configuration file" >}}

```TBD```

{{< /tab >}}
{{< /tabs >}}

## Air-gapped offline environment

Use this template for completely isolated environments with no internet
connectivity. Key configurations in this template include:

- All online features disabled
- Container proxy configured for internal registries
- No telemetry or update checks
- Custom image repositories for Kubernetes (if necessary)
- Strict network isolation

{{< tabs >}}
{{< tab name="Windows deployment (MSI)" >}}

```TBD```

{{< /tab >}}
{{< tab name="macOS deployment (PKG)" >}}

```TBD```

{{< /tab >}}
{{< tab name="JSON configuration file" >}}

```TBD```

{{< /tab >}}
{{< /tabs >}}

## Build your configuration with AI

Use Docker's Ask AI feature to get personalized help configuring your
`admin-settings.json` file. It can help you:

- Understand what each setting does and why it matters for your environment
- Create custom configurations based on your specific requirements
- Troubleshoot deployment issues
- Convert your security policies into Docker Desktop settings

Describe your environment and requirements, and the AI will guide
you through creating your JSON configuration file.

[INSERT KAPA WIDGET]

Example AI prompts:

- "How do I configure Docker Desktop for a financial services environment with
strict compliance requirements?"
- "What settings should I use for a development team that needs to access
internal registries through a corporate proxy?"
- "Help me create a configuration that blocks all external network access but
allows specific internal Docker registries"
- "Explain the security implications of each ECI setting in the hardened
template"

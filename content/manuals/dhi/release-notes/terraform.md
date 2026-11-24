---
title: DHI Terraform provider release notes
linkTitle: Terraform release notes
description: New features, bug fixes, and changes in the DHI Terraform provider
keywords: docker hardened images, dhi, terraform, terraform provider, release notes, changelog
toc_min: 1
toc_max: 2
tags:
  - Release notes
---

This page lists changes in recent stable releases of the [DHI Terraform
provider](/dhi/tools/terraform/). For the full release history, see the
[terraform-provider-dhi releases on
GitHub](https://github.com/docker-hardened-images/terraform-provider-dhi/releases).

<!-- BEGIN GENERATED RELEASES -->

## 1.0.0

{{< release-date date="2026-11-24" >}}

[GitHub release](https://github.com/docker-hardened-images/terraform-provider-dhi/releases/tag/v1.0.0)

This release includes various fixes to `dhi_customization`, such as rejecting empty values and relative file and symlink paths.

### Breaking Changes

- `dhi_customization` now rejects empty values at plan time, with an error that says how to fix it. This also applies on `destroy` and `refresh`. Remove the field, or the empty value, to inherit it from the tag definition: @lukoboi
  - `entrypoint` or `cmd` set to `[]`, or containing an empty argument
  - An `environment` variable with an empty value (`environment = {}` is still fine)
  - A `user` or `group` in `accounts` with an empty `name`
- A relative `file` path, `symlink` path or `symlink` source is now rejected at plan time. These paths must be absolute paths that start with `/`. @lukoboi
- Artifact `includes` and `excludes` are now sent exactly as written. The provider used to add a leading `/` to every path, so `lib` was stored as `/lib`. @lukoboi
  - A customization created with an earlier version shows a one-off in-place update on the next plan that takes the slash out (`/lib` becomes `lib`). Apply it, or add the slash to the configuration to keep it.

### What's Changed

- The docs and examples now show `compression = "ZSTD"` for several targets, and say which values can't be empty and that file paths are absolute. @lukoboi

## 0.12.0

{{< release-date date="2026-11-19" >}}

[GitHub release](https://github.com/docker-hardened-images/terraform-provider-dhi/releases/tag/v0.12.0)

### What's Changed

- `dhi_customization` can now customize several tag definitions at once. Add a `target` block for each one, with its `source`, `destination` and `tag_definition_id`, instead of creating one customization per tag definition. @lukoboi
- Adding or removing a `target` block updates the customization in place, including going from several targets back to one. @lukoboi
- Setting `source`, `destination` and `tag_definition_id` on the resource still works for a single target, so existing configurations don't change. @lukoboi
- A Helm chart customization has exactly one target. Mixing image and Helm chart targets, or using several Helm charts, is rejected. @lukoboi
- Updates the README, docs and examples to explain targets, and adds an example for several targets. @lukoboi

## 0.11.1

{{< release-date date="2026-11-05" >}}

[GitHub release](https://github.com/docker-hardened-images/terraform-provider-dhi/releases/tag/v0.11.1)

### What's Changed

- Updates documentation for customization inheritance

## 0.11.0

{{< release-date date="2026-11-03" >}}

[GitHub release](https://github.com/docker-hardened-images/terraform-provider-dhi/releases/tag/v0.11.0)

### What's Changed

- Adds `variable_mappings` support for Helm chart customizations, allowing each chart variable to point at a DHI image, a mirrored image tag definition, or a customization

<!-- END GENERATED RELEASES -->

---
title: Admission controllers
description: Learn about how admission controllers are used in docker.
keywords: cluster, psp, security
---


This is the current list of admission controllers used by Docker:

### Default
- [NamespaceLifecycle](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#namespacelifecycle)
- [LimitRanger](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#limitranger)
- [ServiceAccount](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#serviceaccount)
- [PersistentVolumeLabel](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#persistentvolumelabel)
- [DefaultStorageClass](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#defaultstorageclass)
- [DefaultTolerationSeconds](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#defaulttolerationseconds)
- [NodeRestriction](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#noderestriction)
- [ResourceQuota](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#resourcequota)
- [PodNodeSelector](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#podnodeselector)
- [PodSecurityPolicy](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#podsecuritypolicy)

### Custom
- UCPAuthorization
- CheckImageSigning
- UCPNodeSelector

**Note:** you cannot enable or disable your own admission controllers. For more information about why, see [Supportability of custom kubernetes flags in universal control plane](https://success.docker.com/article/supportability-of-custom-kubernetes-flags-in-universal-control-plane)

For more information about pod security policies in Docker, see [Pod security policies](/ee/ucp/kubernetes/pod-security-policies.md).
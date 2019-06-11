---
title: Docker Enterprise Platform release notes
description: Learn about the new features, bug fixes, and breaking changes for Docker Enterprise Platform.
keywords: engine enterprise, ucp, dtr, desktop enterprise, whats new, release notes
---

## What’s New?

| Feature | Component | Component version |
|---------|-----------|-------------------|
| [Group Managed Service Accounts (gMSA)](#) | UCP | 3.2.0 |
| [Open Security Controls Assessment Language (OSCAL)](#) | UCP | 3.2.0 |
| [Container storage interface (CSI)](#) | UCP | 3.2.0 |
| [Internet Small Computer System Interface (iSCSI)](#) | UCP | 3.2.0 |
| [System for Cross-domain Identity Management (SCIM)](#) | UCP | 3.2.0 |
| [Registry CLI](#) | DTR | 2.7.0 |
| [App Distribution](#) | DTR | 2.7.0 |
| [Client certificate-based Authentication](#) | DTR | 2.7.0 |
| [Application Designer](/ee/desktop/app-designer/) | Docker Desktop Enterprise | 0.1.4 |
| [Docker Template CLI](/app-template/working-with-template/) | Docker Desktop Enterprise | 0.1.4 |


## Product install and upgrade

| Component Release Notes | Version | Install | Upgrade |
|---------|-----------|-------------------|-------------- |
| [Docker Engine - Enterprise](/engine/release-notes/) | 19.03 | [Install Docker Engine - Enterprise](/ee/supported-platforms/) | [Upgrade Docker Engine - Enterprise](/ee/upgrade/) |
| [Docker UCP](/ee/ucp/release-notes/) | 3.2 | [Install UCP](/ee/ucp/admin/install/) | [Upgrade UCP](/ee/ucp/admin/install/upgrade/) |
| [DTR](/ee/dtr/release-notes/) | 2.7 | [Install DTR](/ee/dtr/admin/install/) | [Upgrade DTR](/ee/dtr/admin/upgrade/) |
| [Docker Desktop Enterprise](/ee/desktop/release-notes/) | 2.1.0 |Install Docker Desktop Enterprise [Mac](/ee/desktop/admin/install/mac/), [Windows](/ee/desktop/admin/install/windows/) | Upgrade Docker Desktop Enterprise  [Mac](/ee/desktop/admin/install/mac/), [Windows](/ee/desktop/admin/install/windows/) |

Refer to the [Compatibility Matrix](https://success.docker.com/article/compatibility-matrix) and the [Maintenance Lifecycle](https://success.docker.com/article/maintenance-lifecycle) for compatibility and software maintenance details.
  

## Known Issues

This is not an exhaustive list. For complete known issues information, refer to the individual component release notes page.
<table>
   <colgroup>
      <col width="20%" />
      <col width="18%" />
      <col width="10%" />
      <col width="20%" />
      <col width="10%" />
      <col width="22%" />
   </colgroup>
   <thead class="night">
    <tr>
       <th>Issue Description</th>
       <th markdown="span">Issue Number</th>
       <th>Component</th>
       <th markdown="span">Affected Versions</th>
       <th>Fixed?</th>
       <th markdown="span">Version Fix - Pull Request</th>
    </tr>
   </thead>
   <tbody>
	<tr>
	 <td><code>docker registry info</code> authentication error (for example purposes)</td>
	 <td>ENG-DTR #912</td>
	 <td>DTR</td>
	 <td>2.7.0-beta2</td>
	 <td>Yes</td>
	 <td>2.7.0</td>
	</tr>
	   <tr>
	 <td>Error when installing UCP with <code>"selinux-enabled": true</code></td>
	 <td>???</td>
	 <td>UCP</td>
	 <td>UCP with Enterprise Engine 18.09 or 19.03</td>
	 <td>No</td>
	 <td>N/A</td>
	</tr>
   </tbody>
</table>
 
 


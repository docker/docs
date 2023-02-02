---
title: Hardened Docker Desktop
description: Overview of what Hardened Docker Desktop is
keywords: security, hardened desktop, enhanced container isolation, registry access management, settings management root access, admins, docker desktop, image access management
---
>Note
>
>Hardened Docker Desktop is available to Docker Business customers only.

Hardened Docker Desktop is a security model for Docker Desktop. It's designed to provide admins with a simple and powerful way to improve their organization's security posture for containerized development, without impacting the developer experience that Docker Desktop offers.

It is for security conscious organizations who don’t give their users root or admin access on their machines, and who would like Docker Desktop to be within their organization’s centralized control.

The Hardened Docker Desktop security model moves the ownership boundary for containers to the organization, meaning that any security controls admins set cannot be altered by the user of Docker Desktop.

Hardened Docker Desktop includes:
- Settings Management, which helps admins to confidently manage and control the usage of Docker Desktop within their organization.
- Enhanced Container Isolation, a setting that instantly enhances security by preventing containers from running as root in Docker Desktop’s Linux VM and ensures that any configurations set using Settings Management, cannot be modified by containers.
- Registry Access Management, which allows admins to control the registries developers can access.
- Image Access Management, which gives admins control over which images developers can pull from Docker Hub.

The features of Hardened Docker Desktop operate independently of each other. When used together, these mechanisms defend against attacks at different functional layers of the developer workflow.

Docker plans to continue adding more security enhancements to the Hardened Docker Desktop security model.

 <div class="component-container">
    <!--start row-->
    <div class="row">
      <div class="col-xs-12 col-sm-12 col-md-12 col-lg-6 block">
        <div class="component">
            <div class="component-icon">
                 <a href="/desktop/hardened-desktop/settings-management/"><img src="/assets/images/lock.svg" alt="Hardened Docker Desktop" width="70" height="70"></a>
            </div>
                <h2 id="hardened-desktop"><a href="/desktop/hardened-desktop/settings-management/">Settings Management </a></h2>
                <p>Learn how Settings Management can secure your developers' workflows.</p>
         </div>
     </div>
     <div class="col-xs-12 col-sm-12 col-md-12 col-lg-6 block">
        <div class="component">
            <div class="component-icon">
                <a href="/desktop/hardened-desktop/enhanced-container-isolation"><img src="/assets/images/secure.svg" alt="Release notes" width="70" height="70"></a>
            </div>
                <h2 id="hardened-desktop"><a href="/desktop/hardened-desktop/enhanced-container-isolation">Enhanced Container Isolation</a></h2>
                <p>Understand how Enhanced Container Isolation can prevent container attacks. </p>
        </div>
    </div>
    </div>
</div>
<div class="component-container">
    <!--start row-->
    <div class="col-xs-12 col-sm-12 col-md-12 col-lg-6 block">
        <div class="component">
            <div class="component-icon">
                 <a href="/desktop/hardened-desktop/registry-access-management/"><img src="/assets/images/registry.svg" alt="Hardened Docker Desktop" width="70" height="70"></a>
            </div>
                <h2 id="hardened-desktop"><a href="/desktop/hardened-desktop/registry-access-management/">Registry Access Management</a></h2>
                <p>Control the registries developers can access while using Docker Desktop.</p>
         </div>
     </div>
     <div class="col-xs-12 col-sm-12 col-md-12 col-lg-6 block">
        <div class="component">
            <div class="component-icon">
                 <a href="/docker-hub/image-access-management/"><img src="/assets/images/image-access.svg" alt="Hardened Docker Desktop" width="70" height="70"></a>
            </div>
                <h2 id="hardened-desktop"><a href="/docker-hub/image-access-management/">Image Access Management</a></h2>
                <p>Control the images developers can pull from Docker Hub.</p>
         </div>
     </div>
</div>


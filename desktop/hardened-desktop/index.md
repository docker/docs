---
title: Hardened Desktop
description: Overview of what Hardened Desktop is
keywords: security, hardened desktop, enhanced container isolation, registry access management, admin controls, root access, administrators, docker desktop
---

Hardened Desktop is a new security model for Docker Desktop. It is part of Docker's ongoing effort to increase Docker Desktop security without impacting the developer experience.

The Hardened Desktop security model provides administrators with a simple and powerful way to increase the security of their containerised development and moves the ownership boundary for containers to the organization, meaning that any security controls administrators set cannot be altered by the user.

It is designed specifically for organizations who don’t give developers root access to their machines, and who would like Docker Desktop to be within the organization’s centralized control.

Hardened Desktop currently includes:
- Admin Controls, which helps administrators to confidently manage and control usage of Docker Desktop.
- Enhanced Container Isolation, an admin setting that instantly enhances security by preventing containers from running as root in Docker Desktop’s Linux VM.
- Registry Access Management, which allows administrators to control the registries developers can access and can be locked in if  Enhanced Container Isolation is switched on by using Admin Controls.

Docker will be adding more security enhancements to the Hardened Desktop security model over the coming months.

 <div class="component-container">
    <!--start row-->
    <div class="row">
      <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
        <div class="component">
            <div class="component-icon">
                 <img src="/assets/images/lock.svg" alt="Hardened Desktop" width="70" height="70">
            </div>
                <h2 id="docker-for-mac"><a href="/desktop/hardened-desktop/admin-controls/">Admin Controls </a></h2>
                <p>Learn how Admin Controls can secure your developers' workflows.</p>
         </div>
     </div>
     <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
        <div class="component">
            <div class="component-icon">
                <img src="/assets/images/secure.svg" alt="Release notes" width="70" height="70">
            </div>
                <h2 id="docker-for-linux"><a href="/desktop/hardened-desktop/enhanced-container-isolation">Enhanced Container Isolation</a></h2>
                <p>Understand how Enhanced Container Isolation can prevent container attacks. </p>
        </div>
    </div>
    <div class="col-xs-12 col-sm-12 col-md-12 col-lg-4 block">
        <div class="component">
            <div class="component-icon">
                 <img src="/assets/images/registry.svg" alt="Hardened Desktop" width="70" height="70">
            </div>
                <h2 id="docker-for-mac"><a href="/desktop/hardened-desktop/regsistry-access-management/">Registry Access Management</a></h2>
                <p>Control the registries developers can access while using Docker Desktop.</p>
         </div>
     </div>
    </div>
</div>


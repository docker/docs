---
description: Instructions on how to enable nested virtualization
keywords: nested virtualization, Docker Desktop
title: Enable nested virtualization
---
Note: This is still work in progress, and the steps/instructions haven’t been validated yet in our test environment.


Docker recommends running Docker Desktop natively on a Windows system (to work with Windows or Linux containers), or on Mac or Linux to work with Linux containers.

However, Docker Desktop can also run inside a virtual machine provided the virtual machine is properly configured. The essential configuration needed here is to have nested virtualization enabled on the virtual machine (host VM). This is because, under the hood, Docker Desktop is using a Linux VM in which it runs the Docker engine and the containers.

Nested virtualization support
The desired scenario in regards to support is that nested virtualization inside the host VM is a supported scenario, but this depends entirely on the hypervisor that runs the host VM.

For example, Microsoft supports running nested Hyper-V both on-prem and on Azure (with some version constraints). But that might not be the case for VMWare ESXi or Citrix Hypervisor.

Docker’s support is limited to installing and running Docker Desktop inside the VM, once the nested virtualization is set up correctly. However, problems and intermittent failures may still occur due to the way these apps virtualize the hardware, and troubleshooting these failures is outside the scope of support from Docker Inc. //TO DO: add correct contract wording here


Enabling nested virtualization

Enabling nested virtualization is a prerequisite for Docker Desktop installation in a virtual machine.

Enable nested virtualization on VMware ESXi 
According to this article, nested virtualization of other hypervisors like Hyper-V inside a vSphere VM is not a supported scenario. 

However, running Hyper-V VM in a VMware ESXi VM is technically possible, and, depending on the version, ESXi includes hardware-assisted virtualization as a supported feature. See this guide for exposing Hardware-Assisted Virtualization to the Guest OS. 
Some network configuration settings might also be needed, please see this article for some useful tips.

Enable nested virtualization on Microsoft Hyper-V 
Nested virtualization is supported by Microsoft for running Hyper-V inside a VM running on a Hyper-V host, in Azure or on-prem (Hyper-V on Hyper-V). Please check the constraints on the host VM operating system/processor and follow the steps documented by Microsoft here for enabling nested virtualization.

Enable nested virtualization on Citrix Hypervisor
Nested virtualization on Citrix Hypervisor is officially unsupported in production scenarios. 
However, running a VM inside a Citrix Hypervisor VM is technically possible, and instructions for enabling nested virtualization can be found in this article, which contains instructions for the only scenario where nested virtualization is supported by Citrix (to support Bromium’s Secure Platform solution). Please note that nested virtualization is only available for Citrix Hypervisor Premium Edition customers or those customers who have access to Citrix Hypervisor through their Citrix Virtual Apps and Desktops entitlement.

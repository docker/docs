---
description: Instructions on how to enable nested virtualization
keywords: nested virtualization, Docker Desktop
title: Run Docker Desktop in a VM or VDI environment
redirect_from:
- /desktop/nested-virtualization/
---

In general, Docker recommends running Docker Desktop natively on either Mac, Linux, or Windows. However, Docker Desktop for Windows can run inside a virtual desktop provided the virtual desktop is properly configured. 

To run Docker Desktop in a virtual desktop environment, it is essential nested virtualization is enabled on the virtual machine that provides the virtual desktop. This is because, under the hood, Docker Desktop is using a Linux VM in which it runs Docker Engine and the containers.

## Nested virtualization support

>Note
>
> Nested virtualization support is only available to Docker Business customers. 

The support available from Docker extends to installing and running Docker Desktop inside the VM, once the nested virtualization is set up correctly. For more information on Docker Desktop support, see [Get support](support.md).

For troubleshooting problems and intermittent failures that are outside of Docker's control, you should contact your hypervisor vendor. Each hypervisor vendor offers different levels of support. For example, Microsoft supports running nested Hyper-V both on-prem and on Azure, with some version constraints. This may not be the case for VMWare ESXi or Citrix Hypervisor.

## Enable nested virtualization

You must enable nested virtualization before you install Docker Desktop on a virtual machine.

### Enable nested virtualization on VMware ESXi

Nested virtualization of other hypervisors like Hyper-V inside a vSphere VM [is not a supported scenario](https://kb.vmware.com/s/article/2009916). However, running Hyper-V VM in a VMware ESXi VM is technically possible and, depending on the version, ESXi includes hardware-assisted virtualization as a supported feature. For internal testing, we used a VM that had 1 CPU with 4 cores and 12GB of memory.

For steps on how to expose hardware-assisted virtualization to the guest OS, [see VMware's documentation](https://docs.vmware.com/en/VMware-vSphere/7.0/com.vmware.vsphere.vm_admin.doc/GUID-2A98801C-68E8-47AF-99ED-00C63E4857F6.html).

You may also need to [configure some network settings](https://www.vembu.com/blog/nested-hyper-v-vms-on-a-vmware-esxi-server).


### Enable nested virtualization on Microsoft Hyper-V

Nested virtualization is supported by Microsoft for running Hyper-V inside a VM running on a Hyper-V host, in Azure or on-prem (Hyper-V on Hyper-V).

For Azure virtual machines, [check that the VM size chosen supports nested virtualization](https://docs.microsoft.com/en-us/azure/virtual-machines/sizes). Microsoft provides [a helpful list on Azure VM sizes](https://docs.microsoft.com/en-us/azure/virtual-machines/acu) and highlights the sizes that currently support nested virtualization. For internal testing, we used D4s_v5 machines. We recommend this specification or above for optimal performance of Docker Desktop.

For on-prem virtual machines, check the constraints on the host VM operating system and [follow the steps documented by Microsoft](https://docs.microsoft.com/en-us/virtualization/hyper-v-on-windows/user-guide/nested-virtualization).


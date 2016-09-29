#!/usr/sbin/dtrace -s
/*
 * vmxexit.d - report all VMX exits for particular VM
 *
 * USAGE: vmxexit.d -p <pid of com.docker.hyperkit>
 */

#pragma D option quiet

string reasons[int];

dtrace:::BEGIN
{
	reasons[0]  = "EXCEPTION";
	reasons[1]  = "EXT_INTR";
	reasons[2]  = "TRIPLE_FAULT";
	reasons[3]  = "INIT";
	reasons[4]  = "SIPI";
	reasons[5]  = "IO_SMI";
	reasons[6]  = "SMI";
	reasons[7]  = "INTR_WINDOW";
	reasons[8]  = "NMI_WINDOW";
	reasons[9]  = "TASK_SWITCH";
	reasons[10] = "CPUID";
	reasons[11] = "GETSEC";
	reasons[12] = "HLT";
	reasons[13] = "INVD";
	reasons[14] = "INVLPG";
	reasons[15] = "RDPMC";
	reasons[16] = "RDTSC";
	reasons[17] = "RSM";
	reasons[18] = "VMCALL";
	reasons[19] = "VMCLEAR";
	reasons[20] = "VMLAUNCH";
	reasons[21] = "VMPTRLD";
	reasons[22] = "VMPTRST";
	reasons[23] = "VMREAD";
	reasons[24] = "VMRESUME";
	reasons[25] = "VMWRITE";
	reasons[26] = "VMXOFF";
	reasons[27] = "VMXON";
	reasons[28] = "CR_ACCESS";
	reasons[29] = "DR_ACCESS";
	reasons[30] = "INOUT";
	reasons[31] = "RDMSR";
	reasons[32] = "WRMSR";
	reasons[33] = "INVAL_VMCS";
	reasons[34] = "INVAL_MSR";
	reasons[36] = "MWAIT";
	reasons[37] = "MTF";
	reasons[39] = "MONITOR";
	reasons[40] = "PAUSE";
	reasons[41] = "MCE_DURING_ENTRY";
	reasons[43] = "TPR";
	reasons[44] = "APIC_ACCESS";
	reasons[45] = "VIRTUALIZED_EOI";
	reasons[46] = "GDTR_IDTR";
	reasons[47] = "LDTR_TR";
	reasons[48] = "EPT_FAULT";
	reasons[49] = "EPT_MISCONFIG";
	reasons[50] = "INVEPT";
	reasons[51] = "RDTSCP";
	reasons[52] = "VMX_PREEMPT";
	reasons[53] = "INVVPID";
	reasons[54] = "WBINVD";
	reasons[55] = "XSETBV";
	reasons[56] = "APIC_WRITE";

	printf("Tracing... Hit Ctrl-C to end.\n");
}

hyperkit$target:::vmx-exit
{
	@num[reasons[arg1], arg0] = count();
}

dtrace:::END
{
	printf("%16s %-4s %8s\n", "REASON", "vCPU", "COUNT");
	printa("%16s %-4d %@8d\n", @num);
}

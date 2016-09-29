#!/usr/sbin/dtrace -s
/*
 * eptfault.d - report all EPT faults for particular VM
 *
 * USAGE: eptfault.d -p <pid of com.docker.hyperkit>
 */

#pragma D option quiet

dtrace:::BEGIN
{
	printf("Tracing... Hit Ctrl-C to end.\n");
}

hyperkit$target:::vmx-ept-fault
{
	@num[arg1, arg0] = count();
}

dtrace:::END
{
	printf("%18s %-4s %8s\n", "ADDRESS", "vCPU", "COUNT");
	printa("%18x %-4d %@8d\n", @num);
}

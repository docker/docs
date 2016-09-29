#ifndef __COMMON_H__
#define __COMMON_H__

static inline uint32_t parse_port(const char *s)
{
	unsigned long p = strtoul(s, NULL, 0);

	if (p == ULONG_MAX) err(1, "strtoul");

	if (p > UINT32_MAX) errx(1, "port 0x%lx out of range\n", p);

	return (uint32_t)p;
}

#endif /*__COMMON_H__*/

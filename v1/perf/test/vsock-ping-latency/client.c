#include <inttypes.h>
#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <strings.h>
#include <err.h>
#include <assert.h>

#include <CoreServices/CoreServices.h>
#include <mach/mach.h>
#include <mach/mach_time.h>

#include <sys/socket.h>
#include <sys/un.h>

#include "common.h"

#define GUEST_CID 3
#define ITERS 500000

struct iter {
	uint64_t start;
	uint64_t end;
	uint64_t elapsed_ns;
};

static mach_timebase_info_data_t tinfo;

static void ping(int s, int i, struct iter *iter)
{
	ssize_t nwr, nrd;
	int j;

	iter->start = mach_absolute_time();
	nwr = write(s, &i, sizeof(i));
	nrd = read(s, &j, sizeof(i));
	iter->end = mach_absolute_time();

	assert(nwr == sizeof(i));
	assert(nwr == nrd);
	assert(i == j);

	iter->elapsed_ns = (iter->end - iter->start) * tinfo.numer / tinfo.denom;
}

static struct iter iters[ITERS];

/*
 * We can always successfully connect to the Unix socket (so long as
 * the hypervisor is running) but can fail to connect to the vsock if
 * the server isn't ready, which results in a 0 return from read. Loop
 * until the server is responding.
 */
static int connect_vsock(const char *cpath, const uint32_t port)
{
	char cmsg[8+1+8+1+1];

	snprintf(cmsg, sizeof(cmsg), "%08"PRIx32".%08"PRIx32"\n",
		 GUEST_CID, port);
	fprintf(stderr, "Connecting to %s:%s", cpath, cmsg);

	struct sockaddr_un un;
	int rc, fd = -1;

	bzero(&un, sizeof(un));

	un.sun_len = 0; /* Unused? */
	un.sun_family = AF_UNIX;
	rc = snprintf(un.sun_path, sizeof(un.sun_path), "%s", cpath);
	if (rc < 0) err(1, "snprintf");

	for (;;) {
		int i = 0;
		ssize_t nwr, nrd;

		fd = socket(AF_UNIX, SOCK_STREAM, 0);
		if (fd < 0) err(1, "socket");

		rc = connect(fd, (struct sockaddr *)&un, sizeof(un));
		if (rc < 0) err(1, "connect");

		assert(strlen(cmsg) == 8 + 1 + 8 + 1);

		nwr = write(fd, cmsg, strlen(cmsg));
		if (nwr < 0) err(1, "connect-write");
		assert(nwr == strlen(cmsg));

		nwr = write(fd, &i, sizeof(i));
		if (nwr < 0) err(1, "ping-write");
		assert(nwr == sizeof(i));

		nrd = read(fd, &i, sizeof(i));
		if (nrd < 0) err(1, "ping-read");
		if (nrd == sizeof(i)) break; /* Success */

		assert(nrd == 0);

		/* Server not yet responding, wait a bit and try again */
		fprintf(stderr, "Server not ready, retrying\n");
		usleep(100000);
	}

	fprintf(stderr, "Connected to %s", cmsg);
	return fd;
}

int main(int argc, char **argv)
{
	if (argc != 3) errx(1, "usage: PORT CONNECT-PATH");

	int fd = connect_vsock(argv[2], parse_port(argv[1]));

	mach_timebase_info(&tinfo);

	fprintf(stderr, "Time scaling factors: %"PRIx32" / %"PRIx32"\n", tinfo.numer, tinfo.denom);

	fprintf(stderr, "Running %d iterations\n", ITERS);

	uint64_t total_ns = 0;
	int i;
	for (i=0; i<ITERS; i++) {
		ping(fd, i, &iters[i]);
		total_ns += iters[i].elapsed_ns;
	}

	fprintf(stderr, "Average: %"PRId64"ns\n", total_ns / ITERS);

	for (i=0; i<ITERS; i++) {
		printf("%"PRId64"\n", iters[i].elapsed_ns);
	}

	return 0;
}

#include <inttypes.h>
#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <strings.h>
#include <err.h>
#include <assert.h>
#include <unistd.h>
#include <signal.h>
#include <limits.h>

#include <sys/socket.h>

#include "common.h"
#include "vm_sockets.h"

static void accept_one(int afd)
{
	int fd = accept(afd, NULL, NULL);
	ssize_t nrd, nwr;
	int i;

	if (fd < 0) err(1, "accept");

	while ((nrd = read(fd, &i, sizeof(i))) > 0) {
		nwr = write(fd, &i, sizeof(i));
		if (nwr < 0) err(1, "write");

		assert(nwr == sizeof(i));
		assert(nwr == nrd);
	}

	if (nrd < 0) err(1, "read");

	assert(nrd == 0);

	close(fd);
}

static void handle_sigint(int sig)
{
	fprintf(stderr, "Caught SIGINT. Exiting\n");
	_exit(0);
}

/*
 * Default behaviour on Linux has SA_RESTART, so accept() is never
 * interrupted. Do it ourselves.
 */
static void setup_signals(void)
{
	int rc;
	struct sigaction sa = {
		.sa_handler = handle_sigint,
	};

	rc = sigaction(SIGINT, &sa, NULL);
	if (rc < 0) err(1, "sigaction");
}

int main(int argc, char **argv)
{
	struct sockaddr_vm vm;
	int rc, fd;

	if (argc != 2) errx(1, "usage: PORT");

	setup_signals();

	const uint32_t port = parse_port(argv[1]);

	fd = socket(AF_VSOCK, SOCK_STREAM, 0);
	if (fd < 0) err(1, "socket");

	vm.svm_family = AF_VSOCK;
	vm.svm_cid = VMADDR_CID_ANY;
	vm.svm_port = port;

	rc = bind(fd, (struct sockaddr *)&vm, sizeof(vm));
	if (rc < 0) err(1, "bind");

	rc = listen(fd, SOMAXCONN);
	if (rc < 0) err(1, "listen");

	for (;;)
		accept_one(fd);

	return 0;
}

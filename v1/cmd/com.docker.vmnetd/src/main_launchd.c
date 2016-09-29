
#include <ctype.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <unistd.h>
#include <getopt.h>
#include <sys/socket.h>
#include <sys/un.h>
#include <errno.h>
#include <sys/types.h>
#include <sys/event.h>
#include <sys/time.h>
#include <sys/stat.h>
#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <fcntl.h>
#include <syslog.h>
#include <libgen.h>
#include <launch.h>

#include "proxy.h"
#include "utils.h"

/* This is the frontend which is spawned by launchd, and which uses the launchd
   socket activation API. */

/* Don't emit debug messages by default. */
int verbose_flag = 0;

int main(int __unused argc, char __unused *argv[]){

  aslInit("Docker", "com.docker.vmnetd");
  
  int kq_fd, err;
  int *listening_fds = NULL;
  size_t n_listening_fds = 0;
  struct kevent kev_init, kev_listener;

  /* Initialise logging */
  aslLog(ASL_LEVEL_NOTICE, "com.docker.vmnetd starting");

  /* Query the listening fds from launchd and register them with kqueue */
  if ((kq_fd = kqueue()) == -1) {
    aslLog(ASL_LEVEL_ERR, "Failed to open kqueue: %s", strerror (errno));
    exit(2);
  };

  if ((err = launch_activate_socket("Listener", &listening_fds, &n_listening_fds)) != 0) {
    aslLog(ASL_LEVEL_ERR, "Failed to launch_activate_socket: %s", strerror (errno));
    exit(3);
  }

  for (size_t i = 0; i < n_listening_fds; i++) {
    EV_SET(&kev_init, *(listening_fds + i), EVFILT_READ, EV_ADD, 0, 0, NULL);
    if (kevent(kq_fd, &kev_init, 1, NULL, 0, NULL) == -1) {
      aslLog(ASL_LEVEL_DEBUG, "Failed to kevent: %s", strerror (errno));
      exit(4);
    }
  }

  /* Forever: accept, create proxy process, repeat */
  while (1) {
    int fd;

    if ((fd = kevent(kq_fd, NULL, 0, &kev_listener, 1, NULL)) == -1) {
      aslLog(ASL_LEVEL_ERR, "Failed to kevent: %s", strerror (errno));
      exit(5);
    }
    if (fd == 0) {
      /* This is a request to shutdown */
      exit(0);
    }

    fork_one_proxy_on((int)kev_listener.ident);
  }
}

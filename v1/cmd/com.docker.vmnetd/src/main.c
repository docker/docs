
#include <ctype.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <unistd.h>
#include <getopt.h>
#include <sys/socket.h>
#include <sys/un.h>
#include <errno.h>
#include <asl.h>
#include <fcntl.h>

#include "proxy.h"
#include "utils.h"

/* This is the non-launchd frontend for the VMnet proxy. This frontend is
   intended to be a simple way to test the shared proxy logic in proxy.c,
   without having to install the binary and configure launchd (... and then
   wonder where the output went etc). */


int verbose_flag = 0;

static void usage(char *argv[]){
  fprintf(stderr, "Usage:\n");
  fprintf(stderr, "  %s --path <path to Unix socket> [--verbose]\n", argv[0]);
  fprintf(stderr, "  -- start a vmnet proxy on <path to Unix socket>\n");
}

static int create_listening_socket(const char *path){
  int fd;
  struct sockaddr_un addr;

  if ((fd = socket(AF_UNIX, SOCK_STREAM, 0)) == -1){
    aslLog(ASL_LEVEL_ERR, "Failed to create socket: %s", strerror (errno));
    exit(1);
  }
  memset(&addr, 0, sizeof(addr));
  addr.sun_family = AF_UNIX;
  strncpy(addr.sun_path, path, sizeof(addr.sun_path)-1);
  if ((unlink(path) == -1) && (errno != ENOENT)) {
    aslLog(ASL_LEVEL_ERR,
      "Failed to unlink old Unix domain socket: %s: %s", path, strerror (errno));
    exit(1);
  }
  if (bind(fd, (struct sockaddr*)&addr, sizeof(addr)) == -1){
    aslLog(ASL_LEVEL_ERR,
      "Failed to bind socket to %s: %s", path, strerror (errno));
    exit(1);
  }
  if (listen(fd, 5) == -1) {
    aslLog(ASL_LEVEL_ERR, "Failed to listen on socket: %s", strerror (errno));
    exit(1);
  }
  return fd;
}

static int test_flag = 0;

int main(int argc, char *argv[]){

  aslInit("Docker", "com.docker.vmnetd");

  char *unix_socket_path = "/var/tmp/com.docker.vmnetd.socket";
  int c, fd;
 
  aslLog(ASL_LEVEL_NOTICE, "proxy starting");

  opterr = 0;
  while (1) {
    static struct option long_options[] = {
      /* These options set a flag. */
      {"verbose", no_argument,       &verbose_flag, 1},
      {"path",    required_argument, NULL, 'p'},
      {"test",    no_argument,       &test_flag, 1},
      {0, 0, 0, 0}
    };
    int option_index = 0;

    c = getopt_long (argc, argv, "vp:", long_options, &option_index);
    if (c == -1) break;

    switch (c) {
      case 'v':
        verbose_flag = 1;
        break;
      case 'p':
        unix_socket_path = optarg;
        break;
      case 0:
        break;
      default:
        usage (argv);
        exit (1);
    }
  }

  if (test_flag) {
    aslLog(ASL_LEVEL_NOTICE, "Attempting to initialise vmnet");
    test_open_vmnet();
    exit(0);
  }
  if (!unix_socket_path) {
    aslLog(ASL_LEVEL_ERR, "Please supply a --path argument.");
    usage(argv);
    exit(1);
  }
  /* Create the listening socket and catch permission errors now */
  fd = create_listening_socket(unix_socket_path);
  if (verbose_flag) {
    aslLog(ASL_LEVEL_NOTICE,
      "Listening on socket bound to %s", unix_socket_path);
  }

  while (1) {
    fork_one_proxy_on(fd);
  }
}

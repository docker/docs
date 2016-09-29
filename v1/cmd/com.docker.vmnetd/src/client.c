
#include <ctype.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <unistd.h>
#include <getopt.h>
#include <sys/socket.h>
#include <sys/un.h>
#include <errno.h>
#include <uuid/uuid.h>
#include "protocol.h"

#include "utils.h"

/* This is an example client program which is intended to be used for testing.
   If everything is working correctly, then this program will connect to the
   daemon and print out a series of advertisement packets to the terminal. */

static int _verbose_flag = 0;
static int loop_flag = 0;
static int uninstall_flag = 0;

static void usage(char *argv[]){
  fprintf(stderr, "Usage:\n");
  fprintf(stderr, "  %s --path <path to Unix socket> [--verbose]\n", argv[0]);
  fprintf(stderr, "  -- connect to a vmnet proxy on <path to Unix socket>, print packets to stdout.\n");
  fprintf(stderr, "  %s --path <path to Unix socket> --uninstall\n", argv[0]);
  fprintf(stderr, "  -- connect to a vmnet proxy on <path to Unix socket>, request a self-uninstall.\n");
  fprintf(stderr, "  %s --path <path to Unix socket> --loop\n", argv[0]);
  fprintf(stderr, "  -- connect to a vmnet proxy on <path to Unix socket>, wait for one packet to\n");
  fprintf(stderr, "     arrive, disconnect and reconnect again. This command is intended for stress-\n");
  fprintf(stderr, "     testing.\n");
}

static int connect_socket(const char *path){
  int fd;
  struct sockaddr_un addr;
  printf("connect\n");
  if ((fd = socket(AF_UNIX, SOCK_STREAM, 0)) == -1){
    fprintf(stderr, "Failed to create socket: %s\n", strerror(errno));
    exit(1);
  }
  memset(&addr, 0, sizeof(addr));
  addr.sun_family = AF_UNIX;
  strncpy(addr.sun_path, path, sizeof(addr.sun_path)-1);
  if(connect(fd, (struct sockaddr *) &addr, sizeof(struct sockaddr_un)) != 0) {
    fprintf(stderr, "Connect to %s failed: %s\n", path, strerror(errno));
    exit(1);
  }
  return fd;
}

static void hexdump(unsigned char *buffer, size_t len){
  char ascii[17];
  size_t i = 0;
  ascii[16] = '\000';
	while (i < len) {
    unsigned char c = *(buffer + i);
    printf("%02x ", c);
    ascii[i++ % 16] = isprint(c)?(signed char)c:'.';
    if ((i % 2) == 0) printf(" ");
    if ((i % 16) == 0) printf(" %s\r\n", ascii);
  };
  printf("\r\n");
}

int main(int argc, char *argv[]){
  char *unix_socket_path = NULL;
  int c, fd;
  uint8_t buffer[2048];
  struct vif_info state;
  uuid_t uuid;
  uuid_string_t uuid_string;

  uuid_generate_random(uuid);
  uuid_unparse(uuid, uuid_string);

  opterr = 0;
  while (1) {
    static struct option long_options[] = {
      /* These options set a flag. */
      {"verbose",   no_argument,       &_verbose_flag, 1},
      {"path",      required_argument, NULL, 'p'},
      {"uuid",      required_argument, NULL, 'u'},
      {"loop",      no_argument,       &loop_flag, 1},
      {"uninstall", no_argument,       &uninstall_flag, 1},

      {0, 0, 0, 0}
    };
    int option_index = 0;

    c = getopt_long (argc, argv, "vp:", long_options, &option_index);
    if (c == -1) break;

    switch (c) {
      case 'v':
        _verbose_flag = 1;
        break;
      case 'p':
        unix_socket_path = optarg;
        break;
      case 'u':
        if (strlen(optarg) != 36) {
          fprintf(stderr, "uuid argument needs to be 36 bytes long\n");
          exit(1);
        }
        memcpy(&uuid_string, optarg, 36);
        break;
      case 0:
        break;
      default:
        usage (argv);
        exit (1);
    }
  }
  if (!unix_socket_path) {
    fprintf(stderr, "Please supply a --path argument.\n");
    usage(argv);
    exit(1);
  }
  while (1) {
    /* Create the listening socket and catch permission errors now */
    fd = connect_socket(unix_socket_path);
    if (_verbose_flag) {
      fprintf(stderr, "Connected socket to %s\n", unix_socket_path);
    }

    struct init_message *ci = create_init_message();

    write_init_message(fd, ci);
    struct init_message you;
    read_init_message(fd, &you);
    char *txt = print_init_message(&you);
    aslLog(ASL_LEVEL_NOTICE, "Server reports %s", txt);
    free(txt);
    enum command command;
    if (uninstall_flag) {
      command = uninstall;
      write_command(fd, &command);
      uint8_t result = 0;
      if (really_read(fd, &result, 1) == -1) {
        aslLog(ASL_LEVEL_ERR, "Unable to read uninstall result");
        exit(EXIT_FAILURE);
      }
      aslLog(ASL_LEVEL_NOTICE, "uninstall result %d = %s", result, (result == 0)?"success":"failure");
      exit(0);
    }
    command = ethernet;
    write_command(fd, &command);
    struct ethernet_args args;
    memcpy(&args.uuid_string[0], &uuid_string[0], 36);
    write_ethernet_args(fd, &args);
    read_vif_info(fd, &state);
    printf("MAC = %02x:%02x:%02x:%02x:%02x:%02x MTU = %lu max_packet_size = %lu\n",
      state.mac[0], state.mac[1], state.mac[2], state.mac[3], state.mac[4], state.mac[5],
      state.mtu, state.max_packet_size);
    if (!state.mtu) exit(1);
    while (1){
      size_t length;

      if (really_read(fd, buffer, 2) == -1) {
        fprintf(stderr, "Cannot read packet header");
        exit(EXIT_FAILURE);
      }
      length = (size_t)(buffer[0] | (buffer[1] << 8));
      if (length > (int)sizeof(buffer)) {
        fprintf(stderr, "Cannot read large packet: %lu ( > %lu)\n", length, sizeof(buffer));
        exit(1);
      }
      if (really_read(fd, buffer, length) == -1) {
        fprintf(stderr, "Cannot read packet");
        exit(EXIT_FAILURE);
      }
      printf("Received packet of length %lu\n", length);
      hexdump(&buffer[0], length);
      if (loop_flag) {
        close(fd);
        break;
      }
    }
  }
}

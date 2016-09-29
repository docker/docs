
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

/* NOTE: this client represents a stand-alone implementation of the v1 of
   the protocol. Consider this code a test-case; please don't make big changes
   to it. */

/* This is an example client program which is intended to be used for testing.
   If everything is working correctly, then this program will connect to the
   daemon and print out a series of advertisement packets to the terminal. */

static int verbose_flag = 0;
static int loop_flag = 0;

static void usage(char *argv[]){
  fprintf(stderr, "Usage:\n");
  fprintf(stderr, "  %s --path <path to Unix socket> [--verbose]\n", argv[0]);
  fprintf(stderr, "  -- connect to a vmnet proxy on <path to Unix socket>, print packets to stdout.\n");
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

static void hexdump(unsigned char *buffer, int len){
  char ascii[17];
  int i = 0;
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

static void really_read(int fd, unsigned char *buffer, int total){
  int n,remaining = total;
  while (remaining > 0){
    n = (int) read(fd, buffer, (size_t) remaining);
    if (n == 0){
      fprintf(stderr, "EOF reading packet from Unix domain socket");
      exit(1);
    }
    remaining -= (int) n;
    buffer = buffer + n;
  }
}

static void really_write(int fd, uint8_t *buffer, int total){
  int n,remaining = total;
  while (remaining > 0){
    n = (int) write(fd, buffer, (size_t) remaining);
    if (n == 0){
      fprintf(stderr, "EOF writing to Unix domain socket");
      exit(1);
    }
    remaining -= (int) n;
    buffer = buffer + n;
  }
}

struct state {
  uint8_t mac[6];
  short _padding;
  int max_packet_size;
  int mtu;
};

static void read_connection_info(int fd, struct state *state) {
  uint8_t buffer[10];
  really_read(fd, &buffer[0], sizeof(buffer));
  state->mtu = buffer[0] | (buffer[1] << 8);
  state->max_packet_size = buffer[2] | (buffer[3] << 8);
  memcpy(state->mac, &buffer[4], 6);
}

static char expected_hello[5] = { 'V', 'M', 'N', 'E', 'T' };

int main(int argc, char *argv[]){
  char *unix_socket_path = NULL;
  int c, fd;
  uint8_t buffer[2048];
  struct state state;
  uuid_t uuid;
  uuid_string_t uuid_string;

  uuid_generate_random(uuid);
  uuid_unparse(uuid, uuid_string);

  opterr = 0;
  while (1) {
    static struct option long_options[] = {
      /* These options set a flag. */
      {"verbose", no_argument,       &verbose_flag, 1},
      {"path",    required_argument, NULL, 'p'},
      {"uuid",    required_argument, NULL, 'u'},
      {"loop",    no_argument,       &loop_flag, 1},
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
    if (verbose_flag) {
      fprintf(stderr, "Connected socket to %s\n", unix_socket_path);
    }

    really_write(fd, (uint8_t*) &expected_hello[0], sizeof(expected_hello));
    really_write(fd, (uint8_t*) &uuid_string[0], 36);

    read_connection_info(fd, &state);
    printf("MAC = %02x:%02x:%02x:%02x:%02x:%02x MTU = %d max_packet_size = %d\n",
      state.mac[0], state.mac[1], state.mac[2], state.mac[3], state.mac[4], state.mac[5],
      state.mtu, state.max_packet_size);
    if (!state.mtu) exit(1);
    while (1){
      int length;

      really_read(fd, buffer, 2);
      length = buffer[0] | (buffer[1] << 8);
      if (length > (int)sizeof(buffer)) {
        fprintf(stderr, "Cannot read large packet: %d ( > %lu)\n", length, sizeof(buffer));
        exit(1);
      }
      really_read(fd, buffer, length);
      printf("Received packet of length %d\n", length);
      hexdump(&buffer[0], length);
      if (loop_flag) {
        close(fd);
        break;
      }
    }
  }
}

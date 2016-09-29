#include <stdio.h>
#include <stdlib.h>
#include <fcntl.h>
#include <unistd.h>
#include <errno.h>

void my_unlink(char * path) {
  if (unlink(path))
    switch (errno) {
    case ENOENT: break;
    default:
      perror("unlink");
      exit(1);
    }
}

void touch(char * path) {
  char * cmd;
  asprintf(&cmd, "touch %s", path);
  if (system(cmd)) {
    perror("system touch");
    exit(1);
  }
  free(cmd);
}

int main() {
  int fd;
  int k;
  char data[] = "data";

  setbuf(stdout, NULL);

  my_unlink("success");
  my_unlink("file");
  my_unlink("newfile");

  fd = open("newfile", O_CREAT | O_EXCL | O_WRONLY, 0600);
  if (fd < 0) {
    perror("open");
    exit(1);
  }

  if (write(fd, data, 4) < 0) {
    perror("write");
    exit(1);
  }

  close(fd);

  touch("start");

  for (int j = 0; j < 2; j++) {
    printf("%03d ",j);

    for (int i = 0; i < 1000; i++) {
      fd = open("newfile", O_RDONLY);
      if (fd < 0) {
        switch (errno) {
        case ENOENT: continue;
        default:
          perror("open");
          exit(1);
        }
      }

      k = read(fd, data, 4);
      if (k < 0) {
        perror("read");
        exit(1);
      }
      if (k != 4) {
        printf("read %d bytes instead of 4\n", k);
        exit(1);
      }

      close(fd);
    }

  }

  touch("success");

  return 0;
}

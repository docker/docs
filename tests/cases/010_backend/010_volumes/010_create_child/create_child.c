#include <dirent.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <stdio.h>

int main() {
  DIR *d;
  struct stat s;

  if (!(d = opendir("/tmp/create_child"))) {
    perror("opendir failed");
    return 1;
  }

  if (creat("/tmp/create_child/child", 0600) < 0) {
    perror("creat (1) failed");
    return 1;
  }

  if (stat("/tmp/create_child", &s) < 0) {
    perror("stat failed");
    return 1;
  }

  if (creat("/tmp/create_child/child2", 0600) < 0) {
    perror("creat (2) failed");
    return 1;
  }

  return 0;
}

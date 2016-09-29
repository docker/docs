#include <unistd.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int main(int argc, char *argv[]) {
  char *path;
  int uid, gid;

  if (argc < 4) {
    printf("expected (uid|gid) id path\n");
    exit(1);
  }

  if (0 == strncmp("uid",argv[1],3)) {
    uid = atoi(argv[2]);
    gid = -1;
  } else if (0 == strncmp("gid",argv[1],3)) {
    uid = -1;
    gid = atoi(argv[2]);
  } else {
    printf("first argument must be uid or gid");
    exit(1);
  }
  path = argv[3];  

  if (lchown(path,uid,gid)) {
    perror("couldn't chown");
    exit(1);
  }
  
  return 0;
}

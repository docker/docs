#include <sys/stat.h>
#include <fcntl.h>
#include <stdio.h>

int main() {
  if (mkdir("/tmp/trans_rename_unlink/d", 0700) < 0) {
    perror("mkdir (1) failed");
    return 1;
  }

  if (mkdir("/tmp/trans_rename_unlink/d/e", 0700) < 0) {
    perror("mkdir (2) failed");
    return 1;
  }

  if (mkdir("/tmp/trans_rename_unlink/target", 0700) < 0) {
    perror("mkdir (3) failed");
    return 1;
  }

  if (mkdir("/tmp/trans_rename_unlink/d/target", 0700) < 0) {
    perror("mkdir (4) failed");
    return 1;
  }

  if (rename("/tmp/trans_rename_unlink/d/e",
             "/tmp/trans_rename_unlink/target") < 0) {
    perror("rename (1) failed");
    return 1;
  }

  if (rename("/tmp/trans_rename_unlink/d/target",
             "/tmp/trans_rename_unlink/d/anything") < 0) {
    perror("rename (2) failed");
    return 1;
  }

  return 0;
}

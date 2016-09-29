#include <stdlib.h>
#include "utils.h"

// logs

static aslclient asl = NULL;
static aslmsg log_msg = NULL;

void aslInit(const char* sender, const char* facility) {
    free(asl);

    asl = asl_open(sender, facility, ASL_OPT_STDERR);
    log_msg = asl_new(ASL_TYPE_MSG);
    asl_set(log_msg, ASL_KEY_SENDER, sender);;
    asl_set(log_msg, ASL_KEY_FACILITY,facility);
    // + ASL_KEY_TIME_NSEC
    // + ASL_KEY_TIME
}

__attribute__((__format__ (__printf__, 2, 0)))
void aslLog(int level, const char *format, ...) {

  if (asl == NULL) {
    aslclient tmp_asl = asl_open("Docker", "com.docker.docker", ASL_OPT_STDERR);
    log_msg = asl_new(ASL_TYPE_MSG);
    asl_set(log_msg, ASL_KEY_SENDER, "Docker");;
    asl_set(log_msg, ASL_KEY_FACILITY,"com.docker.docker");
    // + ASL_KEY_TIME_NSEC
    // + ASL_KEY_TIME
    asl_log(tmp_asl, log_msg, ASL_LEVEL_ERR, "aslInit should be called before aslLog");
    free(tmp_asl);
    return;
  }

  va_list args;
  va_start (args, format);  
  
  asl_vlog(asl, log_msg, level, format, args);

  va_end (args);
}


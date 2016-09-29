#include <fcntl.h>
#include <stdio.h>
#include <time.h>
#include <SystemConfiguration/SystemConfiguration.h>

#include "asl_logger.h"

// logs

static aslclient asl = NULL;
static aslmsg log_msg = NULL;

void apple_asl_logger_init(const char* sender, const char* facility) {
    free(asl);

    asl = asl_open(sender, facility, 0);
    log_msg = asl_new(ASL_TYPE_MSG);
    asl_set(log_msg, ASL_KEY_SENDER, sender);;
    asl_set(log_msg, ASL_KEY_FACILITY,facility);
    // + ASL_KEY_TIME_NSEC
    // + ASL_KEY_TIME
}

void apple_asl_logger_log(int level, const char *message) {

  if (asl == NULL) {
    aslclient tmp_asl = asl_open("Docker", "com.docker.docker", ASL_OPT_STDERR);
    log_msg = asl_new(ASL_TYPE_MSG);
    asl_set(log_msg, ASL_KEY_SENDER, "Docker");;
    asl_set(log_msg, ASL_KEY_FACILITY,"com.docker.docker");
    // + ASL_KEY_TIME_NSEC
    // + ASL_KEY_TIME
    asl_log(tmp_asl, log_msg, ASL_LEVEL_ERR, "asl_logger_init should be called before asl_logger_log");
    free(tmp_asl);
    return;
  }

  asl_log(asl, log_msg, level, "%s", message);
}


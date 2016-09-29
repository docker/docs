
#ifndef asl_logger_h_
#define asl_logger_h_

#include <asl.h>
#include <pwd.h>

// aslInit should be called before calling aslLog
// logFilePath can be NULL, it just won't create a file for the logs. 
extern void apple_asl_logger_init(const char* sender, const char* facility);
extern void apple_asl_logger_log(int level, const char *message);

#endif

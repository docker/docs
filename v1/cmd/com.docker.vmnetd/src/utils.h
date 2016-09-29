
#ifndef utils_h_
#define utils_h_

#include <asl.h>

// aslInit should be called before calling aslLog
// logFilePath can be NULL, it just won't create a file for the logs. 
extern void aslInit(const char* sender, const char* facility);
extern void aslLog(int level, const char *format, ...);

#endif

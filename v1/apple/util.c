#include <asl.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>

#include "asl_logger.h"

char** new_char_array(int size) {
  char** array = calloc(sizeof(char*),size);
  if (array == NULL) {
    apple_asl_logger_log(ASL_LEVEL_CRIT, "unable to allocate memory for new char**");
    abort();
  }
  return array;
}

char* get_array_string(char **a, int n) {
  return a[n];
}

void set_array_string(char **a, char *s, int n) {
  a[n] = s;
}

void free_char_array(char **a, int size) {
  for (int i = 0; i < size; i++) {
    free(a[i]);
    a[i] = NULL;
  }
  free(a);
  a = NULL;
}

char *join_strings(char **src, const char *sep, int count) {
  char *result = NULL;
  size_t length = 0;

  for (int i=0; i < count; i++) {
    length +=strlen(src[i]);
  }

  length += strlen(sep) * (count -1) + 1;

  result = malloc(length);
  if (result == NULL) {
    apple_asl_logger_log(ASL_LEVEL_CRIT, "unable to allocate memory for new char");
    abort();
  }

  for (int i=0; i < count; i++) {
    strcat(result, src[i]);
    if (i < (count - 1)) {
      strcat(result, sep);
    }	
  }
  return result; 
}

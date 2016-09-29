#import <Foundation/Foundation.h>

const char *find_documents(){
  NSArray *results = NSSearchPathForDirectoriesInDomains(NSDocumentDirectory, NSUserDomainMask, true);
  for (NSString *item in results) {
    /* First one will do */
    return strdup([item cStringUsingEncoding:NSUTF8StringEncoding]);
  }
  return NULL;
}


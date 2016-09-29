#import <Foundation/Foundation.h>

// returns 1 if key exists, 0 otherwise
int keyExists(const char* key);

// returns 0 if key doesn't exist
// returns 1 if value is true
// returns 0 if value is false
int boolForKey(const char* key);


#include "preferences.h"

NSString *groupID = @"group.com.docker";

int keyExists(const char* c_key) {
	NSString *key = [[NSString alloc] initWithUTF8String: c_key];
	if (key == nil) { return -1; }
	NSUserDefaults *ud = [[NSUserDefaults alloc] initWithSuiteName:groupID];
	if (ud == nil) { return -1; }
	NSObject *obj = [ud objectForKey:key];
	return (obj != nil) ? 1 : 0;
}

int boolForKey(const char* c_key) {
	NSString *key = [[NSString alloc] initWithUTF8String: c_key];
	if (key == nil) { return -1; }
	NSUserDefaults *ud = [[NSUserDefaults alloc] initWithSuiteName:groupID];
	if (ud == nil) { return -1; }
	BOOL result = [ud boolForKey:key];
	return (result == YES) ? 1 : 0;
}

// void test() {
// 	NSLog(@"🐢 this is objC");
// 	NSUserDefaults *ud = [[NSUserDefaults alloc] initWithSuiteName:@"group.com.docker"];
// 	if (ud == nil)
// 	{
// 		NSLog(@"🐢 ud is nil");
// 		return;
// 	}
// 	NSObject *obj = [ud objectForKey:@"prefAutoStart"];
// 	if (obj == nil)
// 	{
// 		NSLog(@"🐢 autostart NOT FOUND");
// 	}
// 	BOOL result = [ud boolForKey:@"prefAutoStart"];
// 	NSLog(@"🐢 autostart is %@", result ? @"Yes" : @"No");
// }

#import "HockeySDK.h"
#import "HockeySDKPrivate.h"

#import "BITHockeyHelper.h"

#import "BITHockeyBaseManager.h"
#import "BITHockeyBaseManagerPrivate.h"

#import <sys/sysctl.h>
#import <mach-o/ldsyms.h>


@implementation BITHockeyBaseManager {
  NSDateFormatter *_rfc3339Formatter;
}

- (id)init {
  if ((self = [super init])) {
    _appIdentifier = nil;
    _serverURL = kBITHockeySDKURL;
    _userID = nil;
    _userName = nil;
    _userEmail = nil;
    
    NSLocale *enUSPOSIXLocale = [[NSLocale alloc] initWithLocaleIdentifier:@"en_US_POSIX"];
    _rfc3339Formatter = [[NSDateFormatter alloc] init];
    [_rfc3339Formatter setLocale:enUSPOSIXLocale];
    [_rfc3339Formatter setDateFormat:@"yyyy'-'MM'-'dd'T'HH':'mm':'ss'Z'"];
    [_rfc3339Formatter setTimeZone:[NSTimeZone timeZoneForSecondsFromGMT:0]];
  }
  return self;
}

- (id)initWithAppIdentifier:(NSString *)appIdentifier {
  if ((self = [self init])) {
    _appIdentifier = appIdentifier;
  }
  return self;
}



#pragma mark - Private

- (void)reportError:(NSError *)error {
  BITHockeyLog(@"ERROR: %@", [error localizedDescription]);
}

- (NSString *)encodedAppIdentifier {
  return (_appIdentifier ? bit_URLEncodedString(_appIdentifier) : bit_URLEncodedString([[NSBundle mainBundle] objectForInfoDictionaryKey:@"CFBundleIdentifier"]));
}

- (NSString *)getDevicePlatform {
  size_t size;
  sysctlbyname("hw.machine", NULL, &size, NULL, 0);
  char *answer = (char*)malloc(size);
  sysctlbyname("hw.machine", answer, &size, NULL, 0);
  NSString *platform = @(answer);
  free(answer);
  return platform;
}


#pragma mark - Manager Control

- (void)startManager {
}


#pragma mark - Helpers

- (NSDate *)parseRFC3339Date:(NSString *)dateString {
  NSDate *date = nil;
  NSError *error = nil; 
  if (![_rfc3339Formatter getObjectValue:&date forString:dateString range:nil error:&error]) {
    BITHockeyLog(@"INFO: Invalid date '%@' string: %@", dateString, error);
  }
  
  return date;
}


@end

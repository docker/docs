#import <sys/sysctl.h>
#import "BITSystemProfile.h"
#import "BITSystemProfilePrivate.h"

@implementation BITSystemProfile

@synthesize usageStartTimestamp = _usageStartTimestamp;

+ (NSString *)deviceIdentifier {
  char buffer[128];
  
  io_registry_entry_t registry = IORegistryEntryFromPath(kIOMasterPortDefault, "IOService:/");
  CFStringRef uuid = (CFStringRef)IORegistryEntryCreateCFProperty(registry, CFSTR(kIOPlatformUUIDKey), kCFAllocatorDefault, 0);
  IOObjectRelease(registry);
  CFStringGetCString(uuid, buffer, 128, kCFStringEncodingMacRoman);
  CFRelease(uuid);
  
  return @(buffer);
}

+ (NSString *)deviceModel {
  NSString *model = nil;
  
  int error = 0;
  int value = 0;
	size_t length = sizeof(value);
  
  error = sysctlbyname("hw.model", NULL, &length, NULL, 0);
  if (error == 0) {
    char *cpuModel = (char *)malloc(sizeof(char) * length);
    if (cpuModel != NULL) {
      error = sysctlbyname("hw.model", cpuModel, &length, NULL, 0);
      if (error == 0) {
        model = @(cpuModel);
      }
      free(cpuModel);
    }
  }
  
  return model;
}

+ (NSString *)systemVersionString {
	NSString* version = nil;
  
	SInt32 major, minor, bugfix;
	OSErr err1 = Gestalt(gestaltSystemVersionMajor, &major);
	OSErr err2 = Gestalt(gestaltSystemVersionMinor, &minor);
	OSErr err3 = Gestalt(gestaltSystemVersionBugFix, &bugfix);
	if ((!err1) && (!err2) && (!err3)) {
		version = [NSString stringWithFormat:@"%ld.%ld.%ld", (long)major, (long)minor, (long)bugfix];
	}
  
	return version;
}

#if __MAC_OS_X_VERSION_MIN_REQUIRED >= __MAC_10_6
+ (BITSystemProfile *)sharedSystemProfile {
  static BITSystemProfile *sharedInstance = nil;
  static dispatch_once_t pred;
  
  dispatch_once(&pred, ^{
    sharedInstance = [BITSystemProfile alloc];
    sharedInstance = [sharedInstance init];
  });
  
  return sharedInstance;
}
#else
+ (BITSystemProfile *)sharedSystemProfile {
  static BITSystemProfile *sharedInstance = nil;
  
  if (sharedInstance == nil) {
    sharedInstance = [[BITSystemProfile alloc] init];
  }
  
  return sharedInstance;
}
#endif

- (instancetype)init {
  if ((self = [super init])) {
    _usageStartTimestamp = nil;
    _startCounter = 0;
  }
  return self;
}

- (void)dealloc {
  _usageStartTimestamp = nil;
  
}

- (void)startUsageForBundle:(NSBundle *)bundle {
  @synchronized(@"startstop") {
    if (!self.usageStartTimestamp)
      self.usageStartTimestamp = [NSDate date];
    
    _startCounter++;
    
    BOOL newVersion = NO;
    
    if (![[NSUserDefaults standardUserDefaults] valueForKey:kBITUpdateUsageTimeForVersionString]) {
      newVersion = YES;
    } else {
      if ([(NSString *)[[NSUserDefaults standardUserDefaults] valueForKey:kBITUpdateUsageTimeForVersionString] compare:[bundle objectForInfoDictionaryKey:@"CFBundleVersion"]] != NSOrderedSame) {
        newVersion = YES;
      }
    }
    
    if (newVersion) {
      [[NSUserDefaults standardUserDefaults] setObject:@([[NSDate date] timeIntervalSinceReferenceDate]) forKey:kBITUpdateDateOfVersionInstallation];
      [[NSUserDefaults standardUserDefaults] setObject:[bundle objectForInfoDictionaryKey:@"CFBundleVersion"] forKey:kBITUpdateUsageTimeForVersionString];
      [[NSUserDefaults standardUserDefaults] setObject:@0.0 forKey:kBITUpdateUsageTimeOfCurrentVersion];
      [[NSUserDefaults standardUserDefaults] synchronize];
    }
  }
}

- (void)startUsage {
  [self startUsageForBundle:[NSBundle mainBundle]];
}

- (void)stopUsage {
  @synchronized(@"startstop") {
    if (_startCounter > 0)
      _startCounter--;
    
    if (!self.usageStartTimestamp)
      return;
    if (_startCounter > 0)
      return;
    
    double timeDifference = [[NSDate date] timeIntervalSinceReferenceDate] - [self.usageStartTimestamp timeIntervalSinceReferenceDate];
    double previousTimeDifference = [(NSNumber *)[[NSUserDefaults standardUserDefaults] valueForKey:kBITUpdateUsageTimeOfCurrentVersion] doubleValue];
    
    self.usageStartTimestamp = nil;
    
    [[NSUserDefaults standardUserDefaults] setObject:@(previousTimeDifference + timeDifference) forKey:kBITUpdateUsageTimeOfCurrentVersion];
    [[NSUserDefaults standardUserDefaults] synchronize];
  }
}

- (NSString *)currentUsageString {
  double currentUsageTime = [[NSUserDefaults standardUserDefaults] doubleForKey:kBITUpdateUsageTimeOfCurrentVersion];
  
  if (currentUsageTime > 0) {
    // round (up) to 1 minute
    return [NSString stringWithFormat:@"%.0f", ceil(currentUsageTime / 60.0)*60];
  }
  else {
    return @"0";
  }
}

- (NSMutableArray *)systemDataForBundle:(NSBundle *)bundle {
	NSMutableArray *profileArray = [NSMutableArray array];
	NSArray *keys = [self profileKeys];
  
  NSString *uuid = [[self class] deviceIdentifier];
  [profileArray addObject:[NSDictionary dictionaryWithObjects:@[@"udid", @"UDID", uuid, uuid] forKeys:keys]];
  
  NSString *app_version = [bundle objectForInfoDictionaryKey:@"CFBundleVersion"];
  [profileArray addObject:[NSDictionary dictionaryWithObjects:@[@"app_version", @"App Version", app_version, app_version] forKeys:keys]];
  
  if ([[bundle preferredLocalizations] count] > 0) {
    NSString *language = [bundle preferredLocalizations][0];
    [profileArray addObject:[NSDictionary dictionaryWithObjects:@[@"used_lang", @"Used Language", language, language] forKeys:keys]];
  }
  
  NSString *os_version = [[self class] systemVersionString];
  [profileArray addObject:[NSDictionary dictionaryWithObjects:@[@"os_version", @"OS Version", os_version, os_version] forKeys:keys]];
  [profileArray addObject:[NSDictionary dictionaryWithObjects:@[@"os", @"OS", @"Mac OS", @"Mac OS"] forKeys:keys]];
  
  NSString *model = [[self class] deviceModel];
  [profileArray addObject:[NSDictionary dictionaryWithObjects:@[@"model", @"Model", model, model] forKeys:keys]];
  
  return profileArray;
}

- (NSMutableArray *)systemData {
  return [self systemDataForBundle:[NSBundle mainBundle]];
}

- (NSMutableArray *)systemUsageDataForBundle:(NSBundle *)bundle {
  NSMutableArray *profileArray = [self systemDataForBundle:bundle];
	NSArray *keys = [self profileKeys];

  NSString *usageTime = [self currentUsageString];
  [profileArray addObject:[NSDictionary dictionaryWithObjects:@[@"usage_time", @"Usage Time", usageTime, usageTime] forKeys:keys]];

  return profileArray;
}

- (NSMutableArray *)systemUsageData {
  return [self systemUsageDataForBundle:[NSBundle mainBundle]];
}

- (NSArray *)profileKeys {
  return @[@"key", @"displayKey", @"value", @"displayValue"];
}

@end

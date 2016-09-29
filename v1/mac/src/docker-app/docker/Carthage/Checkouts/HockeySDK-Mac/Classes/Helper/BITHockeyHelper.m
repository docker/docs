#import "BITHockeyHelper.h"
#import "HockeySDK.h"
#import "HockeySDKPrivate.h"
#import "BITKeychainItem.h"
#import <sys/sysctl.h>
#import <AppKit/AppKit.h>

NSString *const kBITExcludeApplicationSupportFromBackup = @"kBITExcludeApplicationSupportFromBackup";

typedef struct {
  uint8_t       info_version;
  const char    bit_version[16];
  const char    bit_build[16];
} msai_info_t;

msai_info_t applicationinsights_library_info __attribute__((section("__TEXT,__bit_osx,regular,no_dead_strip"))) = {
  .info_version = 1,
  .bit_version = BITHOCKEY_C_VERSION,
  .bit_build = BITHOCKEY_C_BUILD
};

#pragma mark NSString helpers

NSString *bit_URLEncodedString(NSString *inputString) {
  
  if ([inputString respondsToSelector:@selector(stringByAddingPercentEncodingWithAllowedCharacters:)]) {
    return [inputString stringByAddingPercentEncodingWithAllowedCharacters:[NSCharacterSet characterSetWithCharactersInString:@"!*'();:@&=+$,/?%#[] {}"].invertedSet];
    
  } else {
#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Wdeprecated-declarations"
    return CFBridgingRelease(CFURLCreateStringByAddingPercentEscapes(kCFAllocatorDefault,
                                                                     (__bridge CFStringRef)inputString,
                                                                     NULL,
                                                                     CFSTR("!*'();:@&=+$,/?%#[] {}"),
                                                                     kCFStringEncodingUTF8)
                             );
#pragma clang diagnostic pop
  }
}

NSString *bit_URLDecodedString(NSString *inputString) {
  return CFBridgingRelease(CFURLCreateStringByReplacingPercentEscapesUsingEncoding(kCFAllocatorDefault,
                                                                                   (__bridge CFStringRef)inputString,
                                                                                   CFSTR(""),
                                                                                   kCFStringEncodingUTF8)
                           );
}

NSComparisonResult bit_versionCompare(NSString *stringA, NSString *stringB) {
  // Extract plain version number from self
  NSString *plainSelf = stringA;
  NSRange letterRange = [plainSelf rangeOfCharacterFromSet: [NSCharacterSet letterCharacterSet]];
  if (letterRange.length)
    plainSelf = [plainSelf substringToIndex: letterRange.location];
  
  // Extract plain version number from other
  NSString *plainOther = stringB;
  letterRange = [plainOther rangeOfCharacterFromSet: [NSCharacterSet letterCharacterSet]];
  if (letterRange.length)
    plainOther = [plainOther substringToIndex: letterRange.location];
  
  // Compare plain versions
  NSComparisonResult result = [plainSelf compare:plainOther options:NSNumericSearch];
  
  // If plain versions are equal, compare full versions
  if (result == NSOrderedSame)
    result = [stringA compare:stringB options:NSNumericSearch];
  
  // Done
  return result;
}

#pragma mark Exclude from backup fix

void bit_fixBackupAttributeForURL(NSURL *directoryURL) {
  if (&NSURLIsExcludedFromBackupKey == NULL) {
    BITHockeyLog(@"WARNING: &NSURLIsExcludedBackupKey is NULL, returning");
    return;
  }
  
  BOOL shouldExcludeAppSupportDirFromBackup = [[NSUserDefaults standardUserDefaults] boolForKey:kBITExcludeApplicationSupportFromBackup];
  if (shouldExcludeAppSupportDirFromBackup) {
    return;
  }
  
  if (directoryURL) {
    NSError *getResourceError = nil;
    NSNumber *appSupportDirExcludedValue;
    
    if ([directoryURL getResourceValue:&appSupportDirExcludedValue forKey:NSURLIsExcludedFromBackupKey error:&getResourceError] && appSupportDirExcludedValue) {
      NSError *setResourceError = nil;
      [directoryURL setResourceValue:@NO forKey:NSURLIsExcludedFromBackupKey error:&setResourceError];
    }
  }
}

#pragma mark Identifiers

NSString *bit_mainBundleIdentifier(void) {
  return [[NSBundle mainBundle] objectForInfoDictionaryKey:@"CFBundleIdentifier"];
}

NSString *bit_appIdentifierToGuid(NSString *appIdentifier) {
  NSMutableString *guid;
  NSString *cleanAppId = [appIdentifier stringByTrimmingCharactersInSet:[NSCharacterSet whitespaceCharacterSet]];
  if(cleanAppId && cleanAppId.length == 32) {
    // Insert dashes so that DC will accept th appidentifier (as a replacement for iKey)
    guid = [NSMutableString stringWithString:cleanAppId];
    [guid insertString:@"-" atIndex:20];
    [guid insertString:@"-" atIndex:16];
    [guid insertString:@"-" atIndex:12];
    [guid insertString:@"-" atIndex:8];
  }
  return [guid copy];
}

NSString *bit_appName(NSString *placeHolderString) {
  NSString *appName = [[NSBundle mainBundle] localizedInfoDictionary][@"CFBundleDisplayName"];
  if (!appName)
    appName = [[NSBundle mainBundle] objectForInfoDictionaryKey:@"CFBundleDisplayName"] ?: placeHolderString;
  
  return appName;
}

NSString *bit_appAnonID(BOOL forceNewAnonID) {
  static NSString *appAnonID = nil;
  static dispatch_once_t predAppAnonID;
  NSString *appAnonIDKey = @"appAnonID";

  if (forceNewAnonID) {
    appAnonID = bit_UUID();
    // store this UUID in the keychain (on this device only) so we can be sure to always have the same ID upon app startups
    if (appAnonID) {
      // add to keychain in a background thread, since we got reports that storing to the keychain may take several seconds sometimes and cause the app to be killed
      // and we don't care about the result anyway
      dispatch_async(dispatch_get_global_queue(DISPATCH_QUEUE_PRIORITY_LOW, 0), ^{
        bit_addStringValueToKeychain(appAnonID, appAnonIDKey);
      });
    }
  } else {
    dispatch_once(&predAppAnonID, ^{
      // first check if we already have an install string in the keychain
      appAnonID = bit_stringValueFromKeychainForKey(appAnonIDKey);
      
      if (!appAnonID) {
        appAnonID = bit_UUID();
        // store this UUID in the keychain (on this device only) so we can be sure to always have the same ID upon app startups
        if (appAnonID) {
          // add to keychain in a background thread, since we got reports that storing to the keychain may take several seconds sometimes and cause the app to be killed
          // and we don't care about the result anyway
          dispatch_async(dispatch_get_global_queue(DISPATCH_QUEUE_PRIORITY_LOW, 0), ^{
            bit_addStringValueToKeychain(appAnonID, appAnonIDKey);
          });
        }
      }
    });
  }
  
  return appAnonID;
}

NSString *bit_UUID(void) {
  CFUUIDRef theToken = CFUUIDCreate(NULL);
  CFStringRef uuidStringRef = CFUUIDCreateString(NULL, theToken);
  CFRelease(theToken);
  NSString *stringUUID = [NSString stringWithString:(__bridge NSString *) uuidStringRef];
  CFRelease(uuidStringRef);
  return stringUUID;
}

NSString *bit_settingsDir(void) {
  static NSString *settingsDir = nil;
  static dispatch_once_t predSettingsDir;
  
  dispatch_once(&predSettingsDir, ^{
    NSFileManager *fileManager = [[NSFileManager alloc] init];
    
    NSString *bundleIdentifier = [[NSBundle mainBundle] bundleIdentifier];
    
    // temporary directory for crashes grabbed from PLCrashReporter
    NSArray *paths = NSSearchPathForDirectoriesInDomains(NSCachesDirectory, NSUserDomainMask, YES);
    NSString *cacheDir = paths[0];
    settingsDir = [[cacheDir stringByAppendingPathComponent:bundleIdentifier] stringByAppendingPathComponent:BITHOCKEY_IDENTIFIER];
    
    if (![fileManager fileExistsAtPath:settingsDir]) {
      NSDictionary *attributes = @{NSFilePosixPermissions: @0755UL};
      NSError *theError = NULL;
      
      [fileManager createDirectoryAtPath:settingsDir withIntermediateDirectories: YES attributes: attributes error: &theError];
    }
  });
  
  return settingsDir;
}

#pragma mark - Keychain

BOOL bit_addStringValueToKeychain(NSString *stringValue, NSString *key) {
  if (!key || !stringValue)
    return NO;
  
  NSString *serviceName = [NSString stringWithFormat:@"%@.HockeySDK", [[NSBundle mainBundle] objectForInfoDictionaryKey:@"CFBundleIdentifier"]];
  
  BITGenericKeychainItem *item = [BITGenericKeychainItem genericKeychainItemForService:serviceName withUsername:key];
  
  if (item) {
    // update
    [item setPassword:stringValue];
    return YES;
  } else {
    if ([BITGenericKeychainItem addGenericKeychainItemForService:serviceName withUsername:key password:stringValue])
      return YES;
  }
  
  return NO;
}

NSString *bit_stringValueFromKeychainForKey(NSString *key) {
  if (!key)
    return nil;
  
  NSString *serviceName = [NSString stringWithFormat:@"%@.HockeySDK", [[NSBundle mainBundle] objectForInfoDictionaryKey:@"CFBundleIdentifier"]];
  
  BITGenericKeychainItem *item = [BITGenericKeychainItem genericKeychainItemForService:serviceName withUsername:key];
  if (item) {
    NSString *pwd = [item password];
    return pwd;
  }
  
  return nil;
}

BOOL bit_removeKeyFromKeychain(NSString *key) {
  NSString *serviceName = [NSString stringWithFormat:@"%@.HockeySDK", [[NSBundle mainBundle] objectForInfoDictionaryKey:@"CFBundleIdentifier"]];
  
  BITGenericKeychainItem *item = [BITGenericKeychainItem genericKeychainItemForService:serviceName withUsername:key];
  if (item) {
    [item removeFromKeychain];
    return YES;
  }
  
  return NO;
}

#pragma mark Context helpers

// Return ISO 8601 string representation of the date
NSString *bit_utcDateString(NSDate *date){
  static NSDateFormatter *dateFormatter;
  
  static dispatch_once_t dateFormatterToken;
  dispatch_once(&dateFormatterToken, ^{
    NSLocale *enUSPOSIXLocale = [[NSLocale alloc] initWithLocaleIdentifier:@"en_US_POSIX"];
    dateFormatter = [NSDateFormatter new];
    dateFormatter.locale = enUSPOSIXLocale;
    dateFormatter.dateFormat = @"yyyy-MM-dd'T'HH:mm:ss.SSS'Z'";
    dateFormatter.timeZone = [NSTimeZone timeZoneForSecondsFromGMT:0];
  });
  
  NSString *dateString = [dateFormatter stringFromDate:date];
  
  return dateString;
}

NSString *bit_devicePlatform(void) {
  
  size_t size;
  sysctlbyname("hw.machine", NULL, &size, NULL, 0);
  char *answer = (char*)malloc(size);
  if (answer == NULL)
    return @"";
  sysctlbyname("hw.machine", answer, &size, NULL, 0);
  NSString *platform = [NSString stringWithCString:answer encoding: NSUTF8StringEncoding];
  free(answer);
  return platform;
}

NSString *bit_deviceType(void){
  return @"Desktop";
}

NSString *bit_osVersionBuild(void) {
  void *result = NULL;
  size_t result_len = 0;
  int ret;
  
  /* If our buffer is too small after allocation, loop until it succeeds -- the requested destination size
   * may change after each iteration. */
  do {
    /* Fetch the expected length */
    if ((ret = sysctlbyname("kern.osversion", NULL, &result_len, NULL, 0)) == -1) {
      break;
    }
    
    /* Allocate the destination buffer */
    if (result != NULL) {
      free(result);
    }
    result = malloc(result_len);
    
    /* Fetch the value */
    ret = sysctlbyname("kern.osversion", result, &result_len, NULL, 0);
  } while (ret == -1 && errno == ENOMEM);
  
  /* Handle failure */
  if (ret == -1) {
    int saved_errno = errno;
    
    if (result != NULL) {
      free(result);
    }
    
    errno = saved_errno;
    return NULL;
  }
  
  NSString *osBuild = [NSString stringWithCString:result encoding:NSUTF8StringEncoding];
  free(result);
  
  NSString* osVersion = nil;
  
#if __MAC_OS_X_VERSION_MAX_ALLOWED > 1090
  if ([[NSProcessInfo processInfo] respondsToSelector:@selector(operatingSystemVersion)]) {
    NSOperatingSystemVersion osSystemVersion = [[NSProcessInfo processInfo] operatingSystemVersion];
    osVersion = [NSString stringWithFormat:@"%ld.%ld.%ld", (long)osSystemVersion.majorVersion, (long)osSystemVersion.minorVersion, (long)osSystemVersion.patchVersion];
  } else {
#endif
    SInt32 major, minor, bugfix;
#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Wdeprecated"
    OSErr err1 = Gestalt(gestaltSystemVersionMajor, &major);
    OSErr err2 = Gestalt(gestaltSystemVersionMinor, &minor);
    OSErr err3 = Gestalt(gestaltSystemVersionBugFix, &bugfix);
    if ((!err1) && (!err2) && (!err3)) {
      osVersion = [NSString stringWithFormat:@"%ld.%ld.%ld", (long)major, (long)minor, (long)bugfix];
    }
#if __MAC_OS_X_VERSION_MAX_ALLOWED > 1090
  }
#endif
  
  return [NSString stringWithFormat:@"%@ (%@)", osVersion, osBuild];
}

NSString *bit_osName(void){
  return @"OS X";
}

NSString *bit_deviceLocale(void) {
  NSLocale *locale = [NSLocale currentLocale];
  return [locale objectForKey:NSLocaleIdentifier];
}

NSString *bit_deviceLanguage(void) {
  return [[NSBundle mainBundle] preferredLocalizations][0];
}

NSString *bit_screenSize(void){
  NSScreen *focusScreen = [NSScreen mainScreen];
  CGFloat scale = focusScreen.backingScaleFactor;
  CGSize screenSize = [focusScreen frame].size;
  
  return [NSString stringWithFormat:@"%dx%d",(int)(screenSize.width * scale),(int)(screenSize.height * scale)];
}

NSString *bit_sdkVersion(void){
  return [NSString stringWithFormat:@"osx:%@", [NSString stringWithUTF8String:applicationinsights_library_info.bit_version]];
}

NSString *bit_appVersion(void){
  NSString *build = [[NSBundle mainBundle] infoDictionary][@"CFBundleVersion"];
  NSString *version = [[NSBundle mainBundle] infoDictionary][@"CFBundleShortVersionString"];
  
  if(version){
    return [NSString stringWithFormat:@"%@ (%@)", version, build];
  }else{
    return build;
  }
}

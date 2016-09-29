#import <Foundation/Foundation.h>

FOUNDATION_EXPORT NSString *const kBITExcludeApplicationSupportFromBackup;

/* NSString helpers */
NSString *bit_URLEncodedString(NSString *inputString);
NSString *bit_URLDecodedString(NSString *inputString);
NSComparisonResult bit_versionCompare(NSString *stringA, NSString *stringB);
NSString *bit_mainBundleIdentifier(void);
NSString *bit_appIdentifierToGuid(NSString *appIdentifier);
NSString *bit_appName(NSString *placeHolderString);

NSString *bit_appAnonID(BOOL forceNewAnonID);
NSString *bit_UUID(void);

NSString *bit_settingsDir(void);

BOOL bit_addStringValueToKeychain(NSString *stringValue, NSString *key);
NSString *bit_stringValueFromKeychainForKey(NSString *key);
BOOL bit_removeKeyFromKeychain(NSString *key);

/* Context helpers */
NSString *bit_utcDateString(NSDate *date);
NSString *bit_devicePlatform(void);
NSString *bit_devicePlatform(void);
NSString *bit_deviceType(void);
NSString *bit_osVersionBuild(void);
NSString *bit_osName(void);
NSString *bit_deviceLocale(void);
NSString *bit_deviceLanguage(void);
NSString *bit_screenSize(void);
NSString *bit_sdkVersion(void);
NSString *bit_appVersion(void);

/* Fix bug where Application Support was excluded from backup. */
void bit_fixBackupAttributeForURL(NSURL *directoryURL);

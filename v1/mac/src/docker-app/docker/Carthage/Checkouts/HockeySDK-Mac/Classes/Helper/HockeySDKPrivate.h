#ifndef HockeySDK_HockeySDKPrivate_h
#define HockeySDK_HockeySDKPrivate_h

#define BITHOCKEY_NAME @"HockeySDK"
#define BITHOCKEY_IDENTIFIER @"net.hockeyapp.sdk.mac"
#define BITHOCKEY_CRASH_SETTINGS @"BITCrashManager.plist"
#define BITHOCKEY_CRASH_ANALYZER @"BITCrashManager.analyzer"

#define BITHOCKEY_FEEDBACK_SETTINGS @"BITFeedbackManager.plist"

#define BITHOCKEY_INTEGRATIONFLOW_TIMESTAMP @"BITIntegrationFlowStartTimestamp"

#define BITHockeyBundle [NSBundle bundleWithIdentifier:BITHOCKEY_IDENTIFIER]
//#define BITHOCKEYSDK_URL @"https://sdk.hockeyapp.net/"
extern NSString *const __attribute__((unused)) kBITHockeySDKURL;

extern NSString *const __attribute__((unused)) kBITFeedbackAttachmentLoadedNotification;
extern NSString *const __attribute__((unused)) kBITFeedbackAttachmentLoadedKey;

#define BITHockeyLocalizedString(key,comment) NSLocalizedStringFromTableInBundle(key, @"HockeySDK", BITHockeyBundle, comment)
#define BITHockeyLog(fmt, ...) do { if([BITHockeyManager sharedHockeyManager].isDebugLogEnabled) { NSLog((@"[HockeySDK] %s/%d " fmt), __PRETTY_FUNCTION__, __LINE__, ##__VA_ARGS__); }} while(0)


#define BIT_RGBCOLOR(r,g,b) [NSColor colorWithCalibratedRed:(r)/255.0 green:(g)/255.0 blue:(b)/255.0 alpha:1]

#define BIT_ATTACHMENT_THUMBNAIL_LENGTH 45

#endif



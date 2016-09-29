//
//  HockeySDKPrivate.h
//  HockeySDK
//
//  Created by Andreas Linde on 02.09.13.
//
//

#ifndef HockeySDK_HockeySDKPrivate_h
#define HockeySDK_HockeySDKPrivate_h

#define BITHOCKEY_NAME @"HockeySDK"
#define BITHOCKEY_IDENTIFIER @"net.hockeyapp.sdk.mac"
#define BITHOCKEY_CRASH_SETTINGS @"BITCrashManager.plist"
#define BITHOCKEY_CRASH_ANALYZER @"BITCrashManager.analyzer"

#define BITHOCKEY_INTEGRATIONFLOW_TIMESTAMP @"BITIntegrationFlowStartTimestamp"

#define BITHockeyBundle [NSBundle bundleWithIdentifier:BITHOCKEY_IDENTIFIER]
#define BITHOCKEYSDK_URL @"https://sdk.hockeyapp.net/"

#define BITHockeyLocalizedString(key,comment) NSLocalizedStringFromTableInBundle(key, @"HockeySDK", BITHockeyBundle, comment)
#define BITHockeyLog(fmt, ...) do { if([BITHockeyManager sharedHockeyManager].isDebugLogEnabled) { NSLog((@"[HockeySDK] %s/%d " fmt), __PRETTY_FUNCTION__, __LINE__, ##__VA_ARGS__); }} while(0)

#endif

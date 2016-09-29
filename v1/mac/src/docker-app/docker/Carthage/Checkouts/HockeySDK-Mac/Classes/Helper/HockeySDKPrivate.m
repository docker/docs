#import "HockeySDK.h"
#import "HockeySDKPrivate.h"
#include <CommonCrypto/CommonDigest.h>

NSString *const kBITCrashErrorDomain = @"BITCrashReporterErrorDomain";
NSString *const kBITFeedbackErrorDomain = @"BITFeedbackErrorDomain";
NSString *const kBITHockeyErrorDomain = @"BITHockeyErrorDomain";

NSString *const kBITDefaultUserID = @"default.BITMetaUserID";
NSString *const kBITDefaultUserName = @"default.BITMetaUserName";
NSString *const kBITDefaultUserEmail = @"default.BITMetaUserEmail";

NSString *const kBITFeedbackAttachmentLoadedNotification = @"BITFeedbackAttachmentLoadedNotification";
NSString *const kBITFeedbackAttachmentLoadedKey = @"attachment";

#import <Foundation/Foundation.h>

#import "CrashReporter.h"

// stores the set of crashreports that have been approved but aren't sent yet
#define kHockeySDKApprovedCrashReports @"HockeySDKApprovedCrashReports"

// stores the user name entered in the UI
#define kHockeySDKUserName @"HockeySDKUserName"

// stores the user email address entered in the UI
#define kHockeySDKUserEmail @"HockeySDKUserEmail"


@class BITHockeyAppClient;
@class BITHockeyAttachment;


@interface BITCrashManager ()

///-----------------------------------------------------------------------------
/// @name Delegate
///-----------------------------------------------------------------------------

// delegate is required
@property (nonatomic, unsafe_unretained) id <BITCrashManagerDelegate> delegate;

@property (nonatomic, strong) BITHockeyAppClient *hockeyAppClient;

@property (nonatomic, getter = isCrashManagerActivated) BOOL crashManagerActivated;

@property (nonatomic) NSUncaughtExceptionHandler *plcrExceptionHandler;

@property (nonatomic) PLCrashReporterCallbacks *crashCallBacks;

@property (nonatomic) NSString *lastCrashFilename;

@property (nonatomic, copy, setter = setCrashReportUIHandler:) BITCustomCrashReportUIHandler crashReportUIHandler;

@property (nonatomic, strong) NSString *crashesDir;

- (NSString *)applicationName;
- (NSString *)applicationVersion;

- (void)handleCrashReport;
- (BOOL)hasPendingCrashReport;
- (void)cleanCrashReports;
- (NSString *)extractAppUUIDs:(BITPLCrashReport *)report;

- (void)persistAttachment:(BITHockeyAttachment *)attachment withFilename:(NSString *)filename;

- (BITHockeyAttachment *)attachmentForCrashReport:(NSString *)filename;

- (void)setLastCrashFilename:(NSString *)lastCrashFilename;

/**
 *  Initialize the crash reporter and check if there are any pending crash reports
 *
 *  This method initializes the PLCrashReporter instance if it is not disabled.
 *  It also checks if there are any pending crash reports available that should be send or
 *  presented to the user.
 */
- (void)startManager;

@end

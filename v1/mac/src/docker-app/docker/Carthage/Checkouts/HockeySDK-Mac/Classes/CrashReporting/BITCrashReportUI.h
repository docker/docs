#import <Cocoa/Cocoa.h>

@class BITCrashManager;

@interface BITCrashReportUI : NSWindowController

// defines the users name or user id
@property (nonatomic, strong) NSString *userName;

// defines the users email address
@property (nonatomic, strong) NSString *userEmail;

// set if the nib was loaded correctly
@property (nonatomic, readonly) BOOL nibDidLoadSuccessfully;

- (instancetype)initWithManager:(BITCrashManager *)crashManager
                    crashReport:(NSString *)crashReport
                     logContent:(NSString *)logContent
                applicationName:(NSString *)applicationName
                 askUserDetails:(BOOL)askUserDetails;

- (void)askCrashReportDetails;

- (IBAction)cancelReport:(id)sender;
- (IBAction)submitReport:(id)sender;
- (IBAction)showDetails:(id)sender;
- (IBAction)hideDetails:(id)sender;
- (IBAction)showComments:(id)sender;

- (BOOL)showUserDetails;
- (void)setShowUserDetails:(BOOL)value;

- (BOOL)showComments;
- (void)setShowComments:(BOOL)value;

- (BOOL)showDetails;
- (void)setShowDetails:(BOOL)value;

@end
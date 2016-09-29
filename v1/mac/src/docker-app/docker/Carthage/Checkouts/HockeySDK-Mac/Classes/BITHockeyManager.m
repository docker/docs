#import "HockeySDK.h"
#import "HockeySDKPrivate.h"

#import "BITHockeyBaseManagerPrivate.h"
#import "BITCrashManagerPrivate.h"
#import "BITFeedbackManagerPrivate.h"
#import "BITMetricsManagerPrivate.h"
#import "BITCategoryContainer.h"
#import "BITHockeyHelper.h"
#import "BITHockeyAppClient.h"

NSString *const kBITHockeySDKURL = @"https://sdk.hockeyapp.net/";

@implementation BITHockeyManager {
  NSString *_appIdentifier;
  
  BOOL _validAppIdentifier;
  
  BOOL _startManagerIsInvoked;
  
  BITHockeyAppClient *_hockeyAppClient;
}

#pragma mark - Public Class Methods

+ (BITHockeyManager *)sharedHockeyManager {
  static BITHockeyManager *sharedInstance = nil;
  static dispatch_once_t pred;
  
  dispatch_once(&pred, ^{
    sharedInstance = [BITHockeyManager alloc];
    sharedInstance = [sharedInstance init];
  });
  
  return sharedInstance;
}

- (id) init {
  if ((self = [super init])) {
    _serverURL = nil;
    _delegate = nil;
    _hockeyAppClient = nil;
    
    _disableCrashManager = NO;
    _disableFeedbackManager = NO;
    _disableMetricsManager = NO;
    
    _startManagerIsInvoked = NO;
    
    [self performSelector:@selector(validateStartManagerIsInvoked) withObject:nil afterDelay:0.0f];
  }
  return self;
}

- (void)dealloc {
  _appIdentifier = nil;
  
}


#pragma mark - Private Class Methods

- (BOOL)isSetUpOnMainThread {
  if (!NSThread.isMainThread) {
    NSAssert(NSThread.isMainThread, @"ERROR: This SDK has to be setup on the main thread!");
    
    return NO;
  }
  
  return YES;
}

- (BOOL)checkValidityOfAppIdentifier:(NSString *)identifier {
  BOOL result = NO;
  
  if (identifier) {
    NSCharacterSet *hexSet = [NSCharacterSet characterSetWithCharactersInString:@"0123456789abcdef"];
    NSCharacterSet *inStringSet = [NSCharacterSet characterSetWithCharactersInString:identifier];
    result = ([identifier length] == 32) && ([hexSet isSupersetOfSet:inStringSet]);
  }
  
  return result;
}

- (void)logInvalidIdentifier:(NSString *)environment {
  NSLog(@"[HockeySDK] ERROR: The %@ is invalid! Please use the HockeyApp app identifier you find on the apps website on HockeyApp! The SDK is disabled!", environment);
}

- (NSString *)integrationFlowTimeString {
  NSString *timeString = [[NSBundle mainBundle] objectForInfoDictionaryKey:BITHOCKEY_INTEGRATIONFLOW_TIMESTAMP];
  
  return timeString;
}

- (BOOL)integrationFlowStartedWithTimeString:(NSString *)timeString {
  if (timeString == nil) {
    return NO;
  }
  
  NSDateFormatter *dateFormatter = [[NSDateFormatter alloc] init];
  NSLocale *enUSPOSIXLocale = [[NSLocale alloc] initWithLocaleIdentifier:@"en_US_POSIX"];
  [dateFormatter setLocale:enUSPOSIXLocale];
  [dateFormatter setDateFormat:@"yyyy-MM-dd'T'HH:mm:ssZ"];
  NSDate *integrationFlowStartDate = [dateFormatter dateFromString:timeString];
  
  if (integrationFlowStartDate && [integrationFlowStartDate timeIntervalSince1970] > [[NSDate date] timeIntervalSince1970] - (60 * 10) ) {
    return YES;
  }
  
  return NO;
}

- (void)pingServerForIntegrationStartWorkflowWithTimeString:(NSString *)timeString {
  if (!_appIdentifier) {
    return;
  }
  
  NSString *integrationPath = [NSString stringWithFormat:@"api/3/apps/%@/integration", [_appIdentifier stringByAddingPercentEscapesUsingEncoding:NSUTF8StringEncoding]];
  
  BITHockeyLog(@"INFO: Sending integration workflow ping to %@", integrationPath);
  
  [[self hockeyAppClient] postPath:integrationPath
                        parameters:@{@"timestamp": timeString,
                                     @"sdk": BITHOCKEY_NAME,
                                     @"sdk_version": BITHOCKEY_VERSION,
                                     @"bundle_version": [[NSBundle mainBundle] objectForInfoDictionaryKey:@"CFBundleVersion"]
                                     }
                        completion:^(BITHTTPOperation *operation, NSData* responseData, NSError *error) {
                          switch (operation.response.statusCode) {
                            case 400:
                              BITHockeyLog(@"ERROR: App ID not found");
                              break;
                            case 201:
                              BITHockeyLog(@"INFO: Ping accepted.");
                              break;
                            case 200:
                              BITHockeyLog(@"INFO: Ping accepted. Server already knows.");
                              break;
                            default:
                              BITHockeyLog(@"ERROR: Unknown error");
                              break;
                          }
                        }];
}


#pragma mark - Public Instance Methods (Configuration)

- (void)configureWithIdentifier:(NSString *)appIdentifier {
  _appIdentifier = [appIdentifier copy];
  
  [self initializeModules];
}

- (void)configureWithIdentifier:(NSString *)appIdentifier delegate:(id <BITHockeyManagerDelegate>)delegate {
  _appIdentifier = [appIdentifier copy];
  
  self.delegate = delegate;
  
  [self initializeModules];
}


- (void)configureWithIdentifier:(NSString *)appIdentifier companyName:(NSString *)companyName delegate:(id <BITHockeyManagerDelegate>)delegate {
  _appIdentifier = [appIdentifier copy];
  
  self.delegate = delegate;
  
  [self initializeModules];
}

- (void)startManager {
  if (!_validAppIdentifier || ![self isSetUpOnMainThread]) {
    return;
  }
  
  // Fix bug where Application Support directory was encluded from backup
  NSFileManager *fileManager = [NSFileManager defaultManager];
  NSURL *appSupportURL = [[fileManager URLsForDirectory:NSApplicationSupportDirectory inDomains:NSUserDomainMask] lastObject];
  bit_fixBackupAttributeForURL(appSupportURL);
  
  BITHockeyLog(@"INFO: Starting HockeyManager");
  _startManagerIsInvoked = YES;
  
  // start CrashManager
  if (![self isCrashManagerDisabled]) {
    BITHockeyLog(@"INFO: Start CrashManager");
    if (_serverURL) {
      [_crashManager setServerURL:_serverURL];
    }
    [_crashManager startManager];
  }
  
  // start FeedbackManager
  if (![self isFeedbackManagerDisabled]) {
    BITHockeyLog(@"INFO: Start FeedbackManager");
    if (_serverURL) {
      [_feedbackManager setServerURL:_serverURL];
    }
    [_feedbackManager performSelector:@selector(startManager) withObject:nil afterDelay:1.0f];
  }

	// start MetricsManager
	if (!self.disableMetricsManager) {
		BITHockeyLog(@"INFO: Start MetricsManager");
		[_metricsManager startManager];
		[BITCategoryContainer activateCategory];
	}

  NSString *integrationFlowTime = [self integrationFlowTimeString];
  if (integrationFlowTime && [self integrationFlowStartedWithTimeString:integrationFlowTime]) {
    [self pingServerForIntegrationStartWorkflowWithTimeString:integrationFlowTime];
  }
}

- (void)validateStartManagerIsInvoked {
  if (_validAppIdentifier && !_startManagerIsInvoked) {
    NSLog(@"[HockeySDK] ERROR: You did not call [[BITHockeyManager sharedHockeyManager] startManager] to startup the HockeySDK! Please do so after setting up all properties. The SDK is NOT running.");
  }
}

- (void)setDisableFeedbackManager:(BOOL)disableFeedbackManager {
  if (_feedbackManager) {
    [_feedbackManager setDisableFeedbackManager:disableFeedbackManager];
  }
  _disableFeedbackManager = disableFeedbackManager;
}

- (void)setServerURL:(NSString *)aServerURL {
  // ensure url ends with a trailing slash
  if (![aServerURL hasSuffix:@"/"]) {
    aServerURL = [NSString stringWithFormat:@"%@/", aServerURL];
  }
  
  if (_serverURL != aServerURL) {
    _serverURL = [aServerURL copy];

    if (_hockeyAppClient) {
      _hockeyAppClient.baseURL = [NSURL URLWithString:_serverURL ?: kBITHockeySDKURL];
    }
  }
}

- (void)setDelegate:(id<BITHockeyManagerDelegate>)delegate {
  if (_delegate != delegate) {
    _delegate = delegate;
    
    if (_crashManager) {
      _crashManager.delegate = delegate;
    }
  }
}

- (void)setUserID:(NSString *)userID {
  if (!userID) {
    bit_removeKeyFromKeychain(kBITDefaultUserID);
  } else {
    bit_addStringValueToKeychain(userID, kBITDefaultUserID);
  }
}

- (void)setUserName:(NSString *)userName {
  if (!userName) {
    bit_removeKeyFromKeychain(kBITDefaultUserName);
  } else {
    bit_addStringValueToKeychain(userName, kBITDefaultUserName);
  }
}

- (void)setUserEmail:(NSString *)userEmail {
  if (!userEmail) {
    bit_removeKeyFromKeychain(kBITDefaultUserEmail);
  } else {
    bit_addStringValueToKeychain(userEmail, kBITDefaultUserEmail);
  }
}

- (void)testIdentifier {
  if (!_appIdentifier) {
    return;
  }
  
  NSDate *now = [NSDate date];
  NSString *timeString = [NSString stringWithFormat:@"%.0f", [now timeIntervalSince1970]];
  [self pingServerForIntegrationStartWorkflowWithTimeString:timeString];
}


#pragma mark - Private Instance Methods

- (BITHockeyAppClient *)hockeyAppClient {
  if (!_hockeyAppClient) {
    _hockeyAppClient = [[BITHockeyAppClient alloc] initWithBaseURL:[NSURL URLWithString:_serverURL ?: kBITHockeySDKURL]];
  }
  
  return _hockeyAppClient;
}

- (void)initializeModules {
  _validAppIdentifier = [self checkValidityOfAppIdentifier:_appIdentifier];
  
  if (![self isSetUpOnMainThread]) return;
  
  _startManagerIsInvoked = NO;
  
  BITHockeyLog(@"INFO: Setup CrashManager");
  _crashManager = [[BITCrashManager alloc] initWithAppIdentifier:_appIdentifier];
  _crashManager.delegate = self.delegate;
  _crashManager.hockeyAppClient = [self hockeyAppClient];
  
  // if we don't initialize the BITCrashManager instance, then the delegate will not be invoked
  // leaving the app to never show the window if the developer provided an invalid app identifier
  if (!_validAppIdentifier) {
    [self logInvalidIdentifier:@"app identifier"];
    self.disableCrashManager = YES;
  } else {
    BITHockeyLog(@"INFO: Setup FeedbackManager");
    _feedbackManager = [[BITFeedbackManager alloc] initWithAppIdentifier:_appIdentifier];
    
		BITHockeyLog(@"INFO: Setup MetricsManager");
		NSString *iKey = bit_appIdentifierToGuid(_appIdentifier);
		_metricsManager = [[BITMetricsManager alloc] initWithAppIdentifier:iKey];
  }
	
  if ([self isCrashManagerDisabled])
    _crashManager.crashManagerActivated = NO;
}

@end

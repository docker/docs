#import "BITCrashReportUI.h"

#import <HockeySDK/HockeySDK.h>
#import "HockeySDKPrivate.h"

#import "BITHockeyBaseManagerPrivate.h"
#import "BITCrashManagerPrivate.h"
#import "BITCrashMetaData.h"

#import <sys/sysctl.h>


@interface BITCrashReportUI(private)
- (void) askCrashReportDetails;
- (void) endCrashReporter;
@end

const CGFloat kUserHeight = 50;
const CGFloat kCommentsHeight = 105;
const CGFloat kDetailsHeight = 285;

@implementation BITCrashReportUI {
  IBOutlet NSTextField *nameTextField;
  IBOutlet NSTextField *emailTextField;
  IBOutlet NSTextField *descriptionTextField;
  IBOutlet NSTextView  *crashLogTextView;
  
  IBOutlet NSTextField *nameTextFieldTitle;
  IBOutlet NSTextField *emailTextFieldTitle;
  
  IBOutlet NSTextField *introductionText;
  IBOutlet NSTextField *commentsTextFieldTitle;
  IBOutlet NSTextField *problemDescriptionTextFieldTitle;
  
  IBOutlet NSTextField *noteText;
  
  IBOutlet NSButton *disclosureButton;
  IBOutlet NSButton *showButton;
  IBOutlet NSButton *hideButton;
  IBOutlet NSButton *cancelButton;
  IBOutlet NSButton *submitButton;

  BITCrashManager *_crashManager;
  
  NSString *_applicationName;
  
  NSMutableString *_logContent;
  NSString        *_crashLogContent;
  
  BOOL _showUserDetails;
  BOOL _showComments;
  BOOL _showDetails;
}


- (instancetype)initWithManager:(BITCrashManager *)crashManager crashReport:(NSString *)crashReport logContent:(NSString *)logContent applicationName:(NSString *)applicationName askUserDetails:(BOOL)askUserDetails {
  
  self = [super initWithWindowNibName: @"BITCrashReportUI"];
  
  if ( self != nil) {
    _crashManager = crashManager;
    _crashLogContent = [crashReport copy];
    _logContent = [logContent copy];
    _applicationName = [applicationName copy];
    _userName = @"";
    _userEmail = @"";
    _showComments = YES;
    _showDetails = NO;
    _showUserDetails = askUserDetails;
    _nibDidLoadSuccessfully = NO;

    NSRect windowFrame = [[self window] frame];
    windowFrame.size = NSMakeSize(windowFrame.size.width, windowFrame.size.height - kDetailsHeight);
    windowFrame.origin.y -= kDetailsHeight;
    
    if (!askUserDetails) {
      windowFrame.size = NSMakeSize(windowFrame.size.width, windowFrame.size.height - kUserHeight);
      windowFrame.origin.y -= kUserHeight;
      
      NSRect frame = commentsTextFieldTitle.frame;
      frame.origin.y += kUserHeight;
      commentsTextFieldTitle.frame = frame;

      frame = disclosureButton.frame;
      frame.origin.y += kUserHeight;
      disclosureButton.frame = frame;

      frame = descriptionTextField.frame;
      frame.origin.y += kUserHeight;
      descriptionTextField.frame = frame;
    }
    
    [[self window] setFrame: windowFrame
                    display: YES
                    animate: NO];
    [[self window] center];
    
  }
  return self;
}


- (void)awakeFromNib {
  _nibDidLoadSuccessfully = YES;
  [crashLogTextView setEditable:NO];
  if ([crashLogTextView respondsToSelector:@selector(setAutomaticSpellingCorrectionEnabled:)]) {
    [crashLogTextView setAutomaticSpellingCorrectionEnabled:NO];
  }
}


- (void)endCrashReporter {
  [self close];
}


- (IBAction)showComments: (id) sender {
  NSRect windowFrame = [[self window] frame];
  
  if ([sender intValue]) {
    [self setShowComments: NO];
    
    windowFrame.size = NSMakeSize(windowFrame.size.width, windowFrame.size.height + kCommentsHeight);
    windowFrame.origin.y -= kCommentsHeight;
    [[self window] setFrame: windowFrame
                    display: YES
                    animate: YES];
    
    [self setShowComments: YES];
  } else {
    [self setShowComments: NO];
    
    windowFrame.size = NSMakeSize(windowFrame.size.width, windowFrame.size.height - kCommentsHeight);
    windowFrame.origin.y += kCommentsHeight;
    [[self window] setFrame: windowFrame
                    display: YES
                    animate: YES];
  }
}


- (IBAction)showDetails:(id)sender {
  NSRect windowFrame = [[self window] frame];
  
  windowFrame.size = NSMakeSize(windowFrame.size.width, windowFrame.size.height + kDetailsHeight);
  windowFrame.origin.y -= kDetailsHeight;
  [[self window] setFrame: windowFrame
                  display: YES
                  animate: YES];
  
  [self setShowDetails:YES];
  
}


- (IBAction)hideDetails:(id)sender {
  NSRect windowFrame = [[self window] frame];
  
  [self setShowDetails:NO];
  
  windowFrame.size = NSMakeSize(windowFrame.size.width, windowFrame.size.height - kDetailsHeight);
  windowFrame.origin.y += kDetailsHeight;
  [[self window] setFrame: windowFrame
                  display: YES
                  animate: YES];
}


- (IBAction)cancelReport:(id)sender {
  [_crashManager handleUserInput:BITCrashManagerUserInputDontSend withUserProvidedMetaData:nil];
  
  [self endCrashReporter];
}

- (IBAction)submitReport:(id)sender {
  [showButton setEnabled:NO];
  [hideButton setEnabled:NO];
  [cancelButton setEnabled:NO];
  [submitButton setEnabled:NO];
  
  [[self window] makeFirstResponder: nil];
  
  BITCrashMetaData *crashMetaData = [[BITCrashMetaData alloc] init];
  if (_showUserDetails) {
    crashMetaData.userName = [nameTextField stringValue];
    crashMetaData.userEmail = [emailTextField stringValue];
  }
  crashMetaData.userDescription = [descriptionTextField stringValue];
  
  [_crashManager handleUserInput:BITCrashManagerUserInputSend withUserProvidedMetaData:crashMetaData];
  
  [self endCrashReporter];
}


- (void)askCrashReportDetails {
#define DISTANCE_BETWEEN_BUTTONS		3
  
  NSString *title = BITHockeyLocalizedString(@"WindowTitle", @"");
  [[self window] setTitle:[NSString stringWithFormat:title, _applicationName]];
  
  [[nameTextFieldTitle cell] setTitle:BITHockeyLocalizedString(@"NameTextTitle", @"")];
  [[nameTextField cell] setTitle:self.userName];
  if ([[nameTextField cell] respondsToSelector:@selector(setUsesSingleLineMode:)]) {
    [[nameTextField cell] setUsesSingleLineMode:YES];
  }
  
  [[emailTextFieldTitle cell] setTitle:BITHockeyLocalizedString(@"EmailTextTitle", @"")];
  [[emailTextField cell] setTitle:self.userEmail];
  if ([[emailTextField cell] respondsToSelector:@selector(setUsesSingleLineMode:)]) {
    [[emailTextField cell] setUsesSingleLineMode:YES];
  }

  title = BITHockeyLocalizedString(@"IntroductionText", @"");
  [[introductionText cell] setTitle:[NSString stringWithFormat:title, _applicationName]];
  [[commentsTextFieldTitle cell] setTitle:BITHockeyLocalizedString(@"CommentsDisclosureTitle", @"")];
  [[problemDescriptionTextFieldTitle cell] setTitle:BITHockeyLocalizedString(@"ProblemDetailsTitle", @"")];

  [[descriptionTextField cell] setPlaceholderString:BITHockeyLocalizedString(@"UserDescriptionPlaceholder", @"")];
  [noteText setStringValue:BITHockeyLocalizedString(@"PrivacyNote", @"")];
  
  [showButton setTitle:BITHockeyLocalizedString(@"ShowDetailsButtonTitle", @"")];
  [hideButton setTitle:BITHockeyLocalizedString(@"HideDetailsButtonTitle", @"")];
  [cancelButton setTitle:BITHockeyLocalizedString(@"CancelButtonTitle", @"")];
  [submitButton setTitle:BITHockeyLocalizedString(@"SendButtonTitle", @"")];
  
  // adjust button sizes
  NSDictionary *attrs = @{NSFontAttributeName: [submitButton font]};
  NSSize titleSize = [[submitButton title] sizeWithAttributes: attrs];
	titleSize.width += (16 + 8) * 2;	// 16 px for the end caps plus 8 px padding at each end
	NSRect submitBtnBox = [submitButton frame];
	submitBtnBox.origin.x += submitBtnBox.size.width -titleSize.width;
	submitBtnBox.size.width = titleSize.width;
	[submitButton setFrame: submitBtnBox];
  
  titleSize = [[cancelButton title] sizeWithAttributes: attrs];
	titleSize.width += (16 + 8) * 2;	// 16 px for the end caps plus 8 px padding at each end
	NSRect cancelBtnBox = [cancelButton frame];
	cancelBtnBox.origin.x = submitBtnBox.origin.x -DISTANCE_BETWEEN_BUTTONS -titleSize.width;
	cancelBtnBox.size.width = titleSize.width;
	[cancelButton setFrame: cancelBtnBox];

  titleSize = [[showButton title] sizeWithAttributes: attrs];
	titleSize.width += (16 + 8) * 2;	// 16 px for the end caps plus 8 px padding at each end
	NSRect showBtnBox = [showButton frame];
	showBtnBox.size.width = titleSize.width;
	[showButton setFrame: showBtnBox];

  titleSize = [[hideButton title] sizeWithAttributes: attrs];
	titleSize.width += (16 + 8) * 2;	// 16 px for the end caps plus 8 px padding at each end
	NSRect hideBtnBox = [hideButton frame];
	hideBtnBox.size.width = titleSize.width;
	[hideButton setFrame: showBtnBox];
    
  NSString *logTextViewContent = [_crashLogContent copy];
  
  if (_logContent)
    logTextViewContent = [NSString stringWithFormat:@"%@\n\n%@", logTextViewContent, _logContent];
  
  [crashLogTextView setString:logTextViewContent];
}


- (void)dealloc {
   _crashLogContent = nil;
   _logContent = nil;
   _applicationName = nil;
}


- (BOOL)showUserDetails {
  return _showUserDetails;
}

- (void)setShowUserDetails:(BOOL)value {
  _showUserDetails = value;
}


- (BOOL)showComments {
  return _showComments;
}

- (void)setShowComments:(BOOL)value {
  _showComments = value;
}


- (BOOL)showDetails {
  return _showDetails;
}

- (void)setShowDetails:(BOOL)value {
  _showDetails = value;
}


#pragma mark NSTextField Delegate

- (BOOL)control:(NSControl *)control textView:(NSTextView *)textView doCommandBySelector:(SEL)commandSelector {
  BOOL commandHandled = NO;
  
  if (commandSelector == @selector(insertNewline:)) {
    [textView insertNewlineIgnoringFieldEditor:self];
    commandHandled = YES;
  }
  
  return commandHandled;
}

@end


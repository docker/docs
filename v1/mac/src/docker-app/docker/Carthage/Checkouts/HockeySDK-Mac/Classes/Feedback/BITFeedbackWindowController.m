#import "BITFeedbackWindowController.h"

#import "HockeySDK.h"
#import "HockeySDKPrivate.h"

#import "BITHockeyBaseManagerPrivate.h"
#import "BITFeedbackManagerPrivate.h"
#import "BITFeedbackMessageCellView.h"
#import "BITFeedbackMessageCellViewDelegate.h"

#import "BITFeedbackMessageAttachment.h"

#import "BITFeedbackMessageDateValueTransformer.h"

#import "BITSDKTextView.h"
#import "BITSDKColoredView.h"
#import "BITSDKTextFieldCell.h"

#import <Quartz/Quartz.h>

@interface BITFeedbackWindowController () <NSTableViewDataSource, NSTableViewDelegate, BITSDKTextViewDelegate, QLPreviewPanelDataSource, QLPreviewPanelDelegate, BITFeedbackMessageCellViewDelegate, NSMenuDelegate>

@property (nonatomic, unsafe_unretained) BITFeedbackManager *manager;
@property (nonatomic, strong) NSDateFormatter *lastUpdateDateFormatter;

@property (unsafe_unretained) IBOutlet NSView *userDataView;

@property (unsafe_unretained) IBOutlet BITSDKColoredView *mainBackgroundView;

@property (unsafe_unretained) IBOutlet BITSDKColoredView *userDataBoxView;
@property (unsafe_unretained) IBOutlet NSTextField *contactInfoTextField;
@property (unsafe_unretained) IBOutlet NSTextField *userNameTextField;
@property (unsafe_unretained) IBOutlet NSTextField *userEmailTextField;
@property (unsafe_unretained) IBOutlet NSButton *userDataContinueButton;

@property (nonatomic, copy) NSString *userName;
@property (nonatomic, copy) NSString *userEmail;

@property (unsafe_unretained) IBOutlet BITSDKColoredView *feedbackListBackgroundView;
@property (unsafe_unretained) IBOutlet NSView *feedbackView;
@property (unsafe_unretained) IBOutlet NSView *feedbackEmptyView;
@property (unsafe_unretained) IBOutlet NSImageView *feedbackEmptyAppImageView;

@property (unsafe_unretained) IBOutlet BITSDKColoredView *feedbackComposeBackgroundView;
@property (unsafe_unretained) IBOutlet NSScrollView *feedbackComposeScrollView;

@property (unsafe_unretained) IBOutlet NSScrollView *feedbackAttachmentsScrollView;
@property (unsafe_unretained) IBOutlet NSTableView *feedbackAttachmentsTableView;

@property (unsafe_unretained) IBOutlet NSScrollView *feedbackScrollView;
@property (unsafe_unretained) IBOutlet NSTableView *feedbackTableView;

@property (unsafe_unretained) IBOutlet BITSDKTextView *messageTextField;
@property (nonatomic, strong) NSAttributedString *messageText;

@property (unsafe_unretained) IBOutlet BITSDKColoredView *horizontalLine;
@property (unsafe_unretained) IBOutlet BITSDKColoredView *statusBar;
@property (unsafe_unretained) IBOutlet NSView *statusBarComposeView;
@property (unsafe_unretained) IBOutlet NSButton *sendMessageButton;

@property (unsafe_unretained) IBOutlet NSView *statusBarDefaultView;
@property (unsafe_unretained) IBOutlet NSProgressIndicator *statusBarLoadingIndicator;
@property (unsafe_unretained) IBOutlet NSTextField *statusBarTextField;
@property (unsafe_unretained) IBOutlet NSButton *statusBarRefreshButton;

@property (nonatomic, strong) NSMutableArray *attachments;
@property (nonatomic, strong) NSOperationQueue *thumbnailQueue;

@property (nonatomic, strong) QLPreviewPanel *previewPanel;
@property (nonatomic, strong) BITFeedbackMessageAttachment *previewAttachment;
@property (nonatomic) NSRect previewThumbnailRect;

@property (unsafe_unretained) IBOutlet NSArrayController *composeAttacchmentsArrayController;

- (BOOL)canContinueUserDataView;
- (BOOL)canSendMessage;

- (IBAction)validateUserData:(id)sender;
- (IBAction)sendMessage:(id)sender;
- (IBAction)reloadList:(id)sender;

@end

NSString * const BITFeedbackMessageDateValueTransformerName = @"BITFeedbackMessageDateValueTransformer";

@implementation BITFeedbackWindowController


- (id)initWithManager:(BITFeedbackManager *)feedbackManager {
  self = [super initWithWindowNibName: @"BITFeedbackWindowController"];
  if (self) {
    _manager = feedbackManager;
    
    _attachments = [NSMutableArray new];
    _thumbnailQueue = [NSOperationQueue new];
    
    [NSValueTransformer setValueTransformer:[[BITFeedbackMessageDateValueTransformer alloc] init] forName:BITFeedbackMessageDateValueTransformerName];
    
    self.lastUpdateDateFormatter = [[NSDateFormatter alloc] init];
		[self.lastUpdateDateFormatter setDateStyle:NSDateFormatterShortStyle];
		[self.lastUpdateDateFormatter setTimeStyle:NSDateFormatterShortStyle];
		self.lastUpdateDateFormatter.locale = [NSLocale currentLocale];
  }
  
  return self;
}

- (void)awakeFromNib {
  NSImage *appIcon = [NSImage imageNamed:@"NSApplicationIcon"];
  appIcon = [self imageToGreyImage:appIcon];
  appIcon = [self imageWithReducedAlpha:0.5 fromImage:appIcon];
  [self.feedbackEmptyAppImageView setImage:appIcon];
  
  [self.feedbackListBackgroundView setViewBackgroundColor:[NSColor colorWithCalibratedRed:0.91 green:0.92 blue:0.93 alpha:1.0]];
  [self.feedbackComposeBackgroundView setViewBackgroundColor:[NSColor whiteColor]];
  [self.horizontalLine setViewBackgroundColor:[NSColor colorWithCalibratedRed:0.79 green:0.82 blue:0.83 alpha:1.0]];
  [self.statusBar setViewBackgroundColor:[NSColor whiteColor]];
  [self.mainBackgroundView setViewBackgroundColor:[NSColor colorWithCalibratedRed:0.91 green:0.92 blue:0.93 alpha:1.0]];
  [self.userDataBoxView setViewBackgroundColor:[NSColor colorWithCalibratedRed:0.88 green:0.89 blue:0.90 alpha:1.0]];
  [self.userDataBoxView setViewBorderWidth:1.0];
  [self.userDataBoxView setViewBorderColor:[NSColor colorWithCalibratedRed:0.82 green:0.85 blue:0.86 alpha:1.0]];
  
  [[NSNotificationCenter defaultCenter] addObserver:self
                                           selector:@selector(tableViewFrameChanged:)
                                               name:NSViewFrameDidChangeNotification
                                             object:self.feedbackTableView];
}

- (void)windowDidLoad {
  [super windowDidLoad];
  
  // Implement this method to handle any initialization after your window controller's window has been loaded from its nib file.
  [[NSNotificationCenter defaultCenter] addObserver:self
                                           selector:@selector(startLoadingIndicator)
                                               name:BITHockeyFeedbackMessagesLoadingStarted
                                             object:nil];
  
  [[NSNotificationCenter defaultCenter] addObserver:self
                                           selector:@selector(updateList)
                                               name:BITHockeyFeedbackMessagesLoadingFinished
                                             object:nil];
  
  [self.composeAttacchmentsArrayController setContent:self.attachments];
  [self.feedbackAttachmentsTableView setTarget:self];
  [self.feedbackAttachmentsTableView setDoubleAction:@selector(previewAttachment:)];
  [self.feedbackAttachmentsTableView registerForDraggedTypes:[NSArray arrayWithObject:(NSString*)kUTTypeFileURL]];
  [self.feedbackAttachmentsTableView setMenu:[self contextMenuComposeAttachments]];
  
  [self.statusBarRefreshButton setHidden:YES];
  [self.messageTextField setTypingAttributes:@{NSFontAttributeName: [NSFont userFixedPitchFontOfSize:13.0]}];
  [self.messageTextField setBitDelegate:self];
  [self.messageTextField setPlaceHolderString:@"Your Feedback"];
  
  [self.contactInfoTextField setStringValue:BITHockeyLocalizedString(@"FeedbackContactInfo", @"")];
  [(BITSDKTextFieldCell *)[self.userNameTextField cell] setBitPlaceHolderString: BITHockeyLocalizedString(@"FeedbackName", @"")];
  [(BITSDKTextFieldCell *)[self.userEmailTextField cell] setBitPlaceHolderString: BITHockeyLocalizedString(@"FeedbackEmail", @"")];
  [self.userDataContinueButton setTitle:BITHockeyLocalizedString(@"FeedbackContinueButton", @"")];
  
  [self.sendMessageButton setTitle:BITHockeyLocalizedString(@"FeedbackSendButton", @"")];
  
  // startup
  self.userName = [self.manager userName] ?: @"";
  self.userEmail = [self.manager userEmail] ?: @"";
  
  [self.manager updateMessagesListIfRequired];
  
  if ([self.manager numberOfMessages] == 0 &&
      [self.manager askManualUserDataAvailable] &&
      ![self.manager didAskUserData] &&
      ([self.manager requireManualUserDataMissing] ||
       [self.manager optionalManualUserDataMissing])
      ) {
    [self showUserDataView];
  } else {
    [self showMessagesView];
    [self updateList];
  }
}

- (void)dealloc {
  [[NSNotificationCenter defaultCenter] removeObserver:self name:BITHockeyFeedbackMessagesLoadingStarted object:nil];
  [[NSNotificationCenter defaultCenter] removeObserver:self name:BITHockeyFeedbackMessagesLoadingFinished object:nil];
  [[NSNotificationCenter defaultCenter] removeObserver:self name:NSViewFrameDidChangeNotification object:self.feedbackTableView];

}


#pragma mark - Context menu for compose attachments

- (NSMenu *)contextMenuComposeAttachments {
  NSMenu *menu = [[NSMenu alloc] init];
  menu.delegate = self;
  
  NSMenuItem *viewItem = [[NSMenuItem alloc] initWithTitle:BITHockeyLocalizedString(@"FeedbackAttachmentMenuPreview", @"") action:@selector(previewAttachment:) keyEquivalent:@""];
  [viewItem setTarget:self];
  [viewItem setEnabled:YES];
  [menu addItem:viewItem];
  
  NSMenuItem *deleteItem = [[NSMenuItem alloc] initWithTitle:BITHockeyLocalizedString(@"FeedbackAttachmentMenuRemove", @"") action:@selector(deleteAttachment:) keyEquivalent:@""];
  [deleteItem setTarget:self];
  [deleteItem setEnabled:YES];
  [menu addItem:deleteItem];
  
  return menu;
}

- (void)deleteAttachment:(id)sender {
  NSInteger clickedRow = [self.feedbackAttachmentsTableView clickedRow];
  
  [self.attachments removeObjectAtIndex:clickedRow];
  [self.composeAttacchmentsArrayController setContent:self.attachments];
  
  if ([self.attachments count] == 0) {
    [self hideComposeAttachments];
  }
}


#pragma mark - Private Image methods

- (NSImage *)imageWithReducedAlpha:(CGFloat)alpha fromImage:(NSImage *)image {
  NSSize imageSize = [image size];
  NSRect imageRect = {NSZeroPoint, imageSize};
  
  [image lockFocus];
  [[[NSColor whiteColor] colorWithAlphaComponent: alpha] set];
  NSRectFillUsingOperation(imageRect, NSCompositeSourceAtop);
  [image unlockFocus];
  
  return image;
}

- (NSImage *)imageToGreyImage:(NSImage *)image {
  // Create image rectangle with current image width/height
  CGFloat actualWidth = image.size.width;
  CGFloat actualHeight = image.size.height;
  
  CGRect imageRect = CGRectMake(0, 0, actualWidth, actualHeight);
  CGColorSpaceRef colorSpace = CGColorSpaceCreateDeviceGray();
  
  CGContextRef context = CGBitmapContextCreate(nil, actualWidth, actualHeight, 8, 0, colorSpace, (CGBitmapInfo)kCGImageAlphaNone);
  
  CGImageRef updatedImage1 = [self newImageRefFromImage:image];
  CGContextDrawImage(context, imageRect, updatedImage1);
  
  CGImageRef grayImage = CGBitmapContextCreateImage(context);
  CGColorSpaceRelease(colorSpace);
  CGContextRelease(context);
  CGImageRelease(updatedImage1);
  
  context = CGBitmapContextCreate(nil, actualWidth, actualHeight, 8, 0, nil, (CGBitmapInfo)kCGImageAlphaOnly);
  CGImageRef updatedImage2 = [self newImageRefFromImage:image];
  CGContextDrawImage(context, imageRect, updatedImage2);
  CGImageRef mask = CGBitmapContextCreateImage(context);
  CGContextRelease(context);
  CGImageRelease(updatedImage2);
  
  CGImageRef maskedImage = CGImageCreateWithMask(grayImage, mask);
  NSImage *grayScaleImage =  [self imageFromCGImageRef:maskedImage];
  CGImageRelease(grayImage);
  CGImageRelease(mask);
  CGImageRelease(maskedImage);
  
  // Return the new grayscale image
  return grayScaleImage;
}

- (NSImage*)imageFromCGImageRef:(CGImageRef)image {
  NSRect imageRect = NSMakeRect(0.0, 0.0, 0.0, 0.0);
  CGContextRef imageContext = nil;
  NSImage* newImage = nil; // Get the image dimensions.
  imageRect.size.height = CGImageGetHeight(image);
  imageRect.size.width = CGImageGetWidth(image);
  
  // Create a new image to receive the Quartz image data.
  newImage = [[NSImage alloc] initWithSize:imageRect.size];
  [newImage lockFocus];
  
  // Get the Quartz context and draw.
  imageContext = (CGContextRef)[[NSGraphicsContext currentContext] graphicsPort];
  CGContextDrawImage(imageContext, *(CGRect*)&imageRect, image); [newImage unlockFocus];
  return newImage;
}

- (CGImageRef)newImageRefFromImage:(NSImage*)image; {
  NSData * imageData = [image TIFFRepresentation];
  CGImageRef imageRef;
  if(!imageData) return nil;
  CGImageSourceRef imageSource = CGImageSourceCreateWithData((__bridge CFDataRef)imageData, NULL);
  imageRef = CGImageSourceCreateImageAtIndex(imageSource, 0, NULL);
  CFRelease(imageSource);
  return imageRef;
}


#pragma mark - Private User Data methods

- (void)showUserDataView {
  [self.userDataView setHidden:NO];
  [self.feedbackView setHidden:YES];
  [self.userNameTextField becomeFirstResponder];
}

+ (NSSet *)keyPathsForValuesAffectingCanContinueUserDataView {
  return [NSSet setWithObjects:@"userName",@"userEmail", nil];
}

- (BOOL)canContinueUserDataView {
  BOOL result = YES;
  
  if ([self.manager requireUserName] == BITFeedbackUserDataElementRequired) {
    if (self.userName.length == 0)
      result = NO;
  }
  if (result && [self.manager requireUserEmail] == BITFeedbackUserDataElementRequired) {
    if (self.userEmail.length == 0)
      result = NO;
  }
  
  return result;
}

- (IBAction)validateUserData:(id)sender {
  [self.manager setUserName:self.userName];
  [self.manager setUserEmail:self.userEmail];
  
  [self.manager saveMessages];
  
  [self showMessagesView];
  [self.feedbackTableView becomeFirstResponder];
}


#pragma mark - Private Messages methods

- (void)showMessagesView {
  [self.userDataView setHidden:YES];
  [self.feedbackView setHidden:NO];
  [self.feedbackTableView becomeFirstResponder];
}

- (void)reloadTableAndScrollToBottom {
  [self.feedbackTableView reloadData];
  [self.feedbackTableView scrollToEndOfDocument:self];
}

+ (NSSet *)keyPathsForValuesAffectingCanSendMessage {
  return [NSSet setWithObjects:@"messageText", nil];
}

- (BOOL)canSendMessage {
  return self.messageText.length > 0;
}

- (IBAction)sendMessage:(id)sender {
  [self.manager submitMessageWithText:[self.messageText string] andAttachments:self.attachments];
  self.messageText = nil;
  [self reloadTableAndScrollToBottom];
}

- (void)deleteAllMessages {
  [_manager deleteAllMessages];
  [self reloadTableAndScrollToBottom];
}

- (IBAction)reloadList:(id)sender {
  [self startLoadingIndicator];
  [self.manager updateMessagesList];
}

- (void)updateList {
  [self stopLoadingIndicator];
  
  if ([self.manager numberOfMessages] > 0) {
    [self.statusBarRefreshButton setHidden:NO];
    [self.statusBarTextField setHidden:NO];
    [self.feedbackScrollView setHidden:NO];
    [self.feedbackEmptyView setHidden:YES];
  } else {
    [self.statusBarRefreshButton setHidden:YES];
    [self.statusBarTextField setHidden:YES];
    [self.feedbackScrollView setHidden:YES];
    [self.feedbackEmptyView setHidden:NO];
  }
  
  if ([self.manager numberOfMessages] > 0) {
    [self reloadTableAndScrollToBottom];
  }
}


#pragma mark - Private Status Bar

- (void)startLoadingIndicator {
  [self.statusBarLoadingIndicator setHidden:NO];
  [self.statusBarLoadingIndicator startAnimation:self];
  [self.statusBarRefreshButton setHidden:YES];
}

- (void)stopLoadingIndicator {
  [self.statusBarLoadingIndicator stopAnimation:self];
  [self.statusBarLoadingIndicator setHidden:YES];
  [self updateLastUpdate];
}

- (void)updateLastUpdate {  
  NSString *text = [NSString stringWithFormat:@"%@: %@",
                    BITHockeyLocalizedString(@"FeedbackLastUpdate", @""),
                    [self.manager lastCheck] ? [self.lastUpdateDateFormatter stringFromDate:[self.manager lastCheck]] : BITHockeyLocalizedString(@"FeedbackLastUpdateNever", @"")];
  
  NSFont *boldFont = [NSFont boldSystemFontOfSize:11];

  NSMutableDictionary *style = [NSMutableDictionary dictionary];
  style[NSFontAttributeName] = boldFont;
  
  NSMutableAttributedString *attributedText = [[NSMutableAttributedString alloc] initWithString:text];
  [attributedText beginEditing];
  [attributedText addAttribute:NSFontAttributeName
                 value:boldFont
                 range:NSMakeRange(0, 12)];
  [attributedText endEditing];
  
  NSMutableParagraphStyle *paraStyle = [[NSMutableParagraphStyle alloc] init];
  [paraStyle setAlignment:NSCenterTextAlignment];
  [attributedText addAttributes:@{NSParagraphStyleAttributeName: paraStyle} range:NSMakeRange(0, [attributedText length])];

  self.statusBarTextField.attributedStringValue = attributedText;
}


#pragma mark - Private

- (void)tableViewFrameChanged:(id)sender {
  // this may not be the fastest approach, but don't know of any better at the moment
  [NSAnimationContext beginGrouping];
  [[NSAnimationContext currentContext] setDuration:0];
  [self.feedbackTableView noteHeightOfRowsWithIndexesChanged:[NSIndexSet indexSetWithIndexesInRange:NSMakeRange(0, [self.feedbackTableView numberOfRows])]];
  [NSAnimationContext endGrouping];
}


#pragma mark - Compose View

- (void)showComposeAttachments {
  if (self.feedbackAttachmentsScrollView.isHidden) {
    [self.feedbackAttachmentsScrollView setHidden:NO];
    
    NSRect frame = self.feedbackComposeScrollView.frame;
    frame.size.width -= self.feedbackAttachmentsScrollView.frame.size.width;
    [self.feedbackComposeScrollView setFrame:frame];
  }
}

- (void)hideComposeAttachments {
  if (!self.feedbackAttachmentsScrollView.isHidden) {
    [self.feedbackAttachmentsScrollView setHidden:YES];
    
    NSRect frame = self.feedbackComposeScrollView.frame;
    frame.size.width += self.feedbackAttachmentsScrollView.frame.size.width;
    [self.feedbackComposeScrollView setFrame:frame];
  }
}

- (void)previewAttachment:(id)sender {
  NSInteger clickedRow = self.feedbackAttachmentsTableView.clickedRow;
  
  self.previewAttachment = self.attachments[clickedRow];
  
  NSRect thumbnailRect = [self.feedbackAttachmentsTableView frameOfCellAtColumn:0 row:clickedRow];
  self.previewThumbnailRect = thumbnailRect;
  
  [self togglePreviewPanel:self];
}

- (void)addAttachmentWithFilename:(NSString *)filename {
  NSError *error = nil;
  NSData *data = [NSData dataWithContentsOfFile:filename options:0 error:&error];
  if (!error && data) {
    CFStringRef fileExtension = (__bridge CFStringRef)[filename pathExtension];
    CFStringRef UTI = UTTypeCreatePreferredIdentifierForTag(kUTTagClassFilenameExtension, fileExtension, NULL);

    NSString *mimeTypeString = nil;
    if (UTI) {
      CFStringRef mimeType = UTTypeCopyPreferredTagWithClass(UTI, kUTTagClassMIMEType);
      CFRelease(UTI);
      mimeTypeString = (__bridge_transfer NSString *)mimeType;
    }
    
    BITFeedbackMessageAttachment *attachment = [BITFeedbackMessageAttachment attachmentWithData:data contentType:mimeTypeString];
    attachment.originalFilename = filename;
    
    [self.attachments addObject:attachment];
    [self.composeAttacchmentsArrayController setContent:self.attachments];
    [self.feedbackAttachmentsTableView selectRowIndexes:[NSIndexSet indexSetWithIndex:[self.attachments count] - 1] byExtendingSelection:NO];
    [self.feedbackAttachmentsTableView scrollRowToVisible:[self.attachments count] - 1];
    
    [self showComposeAttachments];
  }
}


#pragma mark - BITTextViewDelegate

- (void)textView:(BITSDKTextView *)textView dragOperationWithFilename:(NSString *)filename {
  [self addAttachmentWithFilename:filename];
}


#pragma mark - Table view data source

- (NSDragOperation)tableView:(NSTableView *)aTableView validateDrop:(id < NSDraggingInfo >)info proposedRow:(NSInteger)row proposedDropOperation:(NSTableViewDropOperation)operation {
  
  if (aTableView == self.feedbackAttachmentsTableView) {
    if (row < aTableView.numberOfRows) return NO;
    
    NSPasteboard *pb = [info draggingPasteboard];
    NSDragOperation dragOperation = [info draggingSourceOperationMask];
    
    if ([[pb types] containsObject:NSFilenamesPboardType]) {
      if (dragOperation & NSDragOperationCopy) {
        return NSDragOperationCopy;
      }
    }
  }
  
  return NSDragOperationNone;
}

- (BOOL)tableView:(NSTableView *)aTableView acceptDrop:(id < NSDraggingInfo >)info row:(NSInteger)row dropOperation:(NSTableViewDropOperation)operation {
  
  if (aTableView == self.feedbackAttachmentsTableView) {
    if (row < aTableView.numberOfRows) return NO;
    
    NSPasteboard *pb = [info draggingPasteboard];
    
    if ( [[pb types] containsObject:NSFilenamesPboardType] ) {
      NSFileManager *fm = [[NSFileManager alloc] init];
      
      NSArray *filenames = [pb propertyListForType:NSFilenamesPboardType];
      
      BOOL fileFound = NO;
      
      for (NSString *filename in filenames) {
        BOOL isDir = NO;
        if (![fm fileExistsAtPath:filename isDirectory:&isDir] || isDir) continue;
        
        fileFound = YES;
        
        [self addAttachmentWithFilename:filename];
      }
      return fileFound;
    }
  }
  
  return NO;
}


- (NSInteger)numberOfRowsInTableView:(NSTableView *)aTableView {
  return [self.manager numberOfMessages];
}

- (CGFloat)tableView:(NSTableView *)tableView heightOfRow:(NSInteger)row {
  BITFeedbackMessage *message = [self.manager messageAtIndex:row];
  
  CGFloat height = [BITFeedbackMessageCellView heightForRowWithMessage:message tableViewWidth:tableView.frame.size.width];
  return height;
}

- (NSView *)tableView:(NSTableView *)tableView viewForTableColumn:(NSTableColumn *)tableColumn row:(NSInteger)row {
  BITFeedbackMessageCellView *result = [tableView makeViewWithIdentifier:[BITFeedbackMessageCellView identifier] owner:self];
  
  BITFeedbackMessage *message = [self.manager messageAtIndex:row];
  if (result == nil) {
    CGFloat height = [BITFeedbackMessageCellView heightForRowWithMessage:message tableViewWidth:tableView.frame.size.width];
    result = [[BITFeedbackMessageCellView alloc] initWithFrame:NSMakeRect(0, 0, self.feedbackTableView.frame.size.width, height) delegate:self];
    result.identifier = [BITFeedbackMessageCellView identifier];
  }
  
  result.message = message;
  [result updateAttachmentViews];

  for (BITFeedbackMessageAttachment *attachment in message.attachments) {
    if (attachment.needsLoadingFromURL && !attachment.isLoading){
      attachment.isLoading = YES;
      NSMutableURLRequest *request = [NSMutableURLRequest requestWithURL:[NSURL URLWithString:attachment.sourceURL]];
      __weak typeof (self) weakSelf = self;
      id nsurlsessionClass = NSClassFromString(@"NSURLSessionDataTask");
      if (nsurlsessionClass) {
        NSURLSessionConfiguration *sessionConfiguration = [NSURLSessionConfiguration defaultSessionConfiguration];
        __block NSURLSession *session = [NSURLSession sessionWithConfiguration:sessionConfiguration];
        
        NSURLSessionDataTask *task = [session dataTaskWithRequest:request
                                                completionHandler: ^(NSData *data, NSURLResponse *response, NSError *error) {
                                                  typeof (self) strongSelf = weakSelf;
                                                  
                                                  [session finishTasksAndInvalidate];
                                                  
                                                  [strongSelf handleResponseForAttachment:attachment responseData:data error:error];
                                                }];
        [task resume];
      }else{
#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Wdeprecated-declarations"
        [NSURLConnection sendAsynchronousRequest:request queue:self.thumbnailQueue completionHandler:^(NSURLResponse *response, NSData *responseData, NSError *err) {
#pragma clang diagnostic pop
          typeof (self) strongSelf = weakSelf;
          [strongSelf handleResponseForAttachment:attachment responseData:responseData error:err];
        }];
      }
    }
  }
  
  return result;
}

- (void)handleResponseForAttachment:(BITFeedbackMessageAttachment *)attachment responseData:(NSData *)responseData error:(NSError *)error {
  attachment.isLoading = NO;
  if (responseData.length) {
    dispatch_async(dispatch_get_main_queue(), ^{
      [attachment replaceData:responseData];
      [[NSNotificationCenter defaultCenter] postNotificationName:kBITFeedbackAttachmentLoadedNotification object:self userInfo:@{kBITFeedbackAttachmentLoadedKey: attachment}];
      [[BITHockeyManager sharedHockeyManager].feedbackManager saveMessages];
    });
  }
}

#pragma mark - NSSplitView Delegate

- (CGFloat)splitView:(NSSplitView *)splitView constrainMaxCoordinate:(CGFloat)proposedMax ofSubviewAt:(NSInteger)dividerIndex {
  CGFloat maximumSize = splitView.frame.size.height - 50;
  
  return maximumSize;
}

- (CGFloat)splitView:(NSSplitView *)splitView constrainMinCoordinate:(CGFloat)proposedMin ofSubviewAt:(NSInteger)dividerIndex {
  CGFloat minimumSize = splitView.frame.size.height - 300;
  
  return minimumSize;
}

- (void)splitView:(NSSplitView *)sender resizeSubviewsWithOldSize:(NSSize)oldSize {
  CGFloat dividerThickness = [sender dividerThickness];
  NSRect topRect  = [[sender subviews][0] frame];
  NSRect bottomRect = [[sender subviews][1] frame];
  NSRect newFrame  = [sender frame];
  
  topRect.size.height = newFrame.size.height - bottomRect.size.height - dividerThickness;
  topRect.size.width = newFrame.size.width;
  topRect.origin = NSMakePoint(0, 0);
  bottomRect.size.width = newFrame.size.width;
  bottomRect.origin.y = topRect.size.height + dividerThickness;
  
  [[sender subviews][0] setFrame:topRect];
  [[sender subviews][1] setFrame:bottomRect];
}


#pragma mark - BITFeedbackMessageCellViewDelegate

- (void)messageCellView:(BITFeedbackMessageCellView *)messaggeCellView clickOnButton:(NSButton *)button withAttachment:(BITFeedbackMessageAttachment *)attachment {
  self.previewAttachment = attachment;

  NSInteger index = [self.feedbackTableView rowForView:messaggeCellView];
  NSRect thumbnailRect = [self.feedbackTableView frameOfCellAtColumn:0 row:index];
  thumbnailRect.origin.x += button.frame.origin.x;
  thumbnailRect.origin.y += (thumbnailRect.size.height - 2 * button.frame.origin.y);
  thumbnailRect.size = button.frame.size;
  self.previewThumbnailRect = thumbnailRect;
  
  [self togglePreviewPanel:self];
}

#pragma mark - Quick Look panel support

- (IBAction)togglePreviewPanel:(id)sender {
  if ([QLPreviewPanel sharedPreviewPanelExists] && [[QLPreviewPanel sharedPreviewPanel] isVisible]) {
    [[QLPreviewPanel sharedPreviewPanel] orderOut:nil];
  } else {
    [[QLPreviewPanel sharedPreviewPanel] makeKeyAndOrderFront:nil];
  }
}

- (BOOL)acceptsPreviewPanelControl:(QLPreviewPanel *)panel {
  return YES;
}

- (void)beginPreviewPanelControl:(QLPreviewPanel *)panel {
  _previewPanel = panel;
  panel.delegate = self;
  panel.dataSource = self;
}

- (void)endPreviewPanelControl:(QLPreviewPanel *)panel {
  _previewPanel = nil;
}


#pragma mark - QLPreviewPanelDataSource

- (NSInteger)numberOfPreviewItemsInPreviewPanel:(QLPreviewPanel *)panel {
  return 1;
}

- (id <QLPreviewItem>)previewPanel:(QLPreviewPanel *)panel previewItemAtIndex:(NSInteger)index {
  if (self.previewAttachment.needsLoadingFromURL && !self.previewAttachment.isLoading) {
    __weak QLPreviewPanel* blockPanel = panel;
    
    self.previewAttachment.isLoading = YES;
    NSMutableURLRequest *request = [NSMutableURLRequest requestWithURL:[NSURL URLWithString:self.previewAttachment.sourceURL]];
    
    __weak typeof (self) weakSelf = self;
    id nsurlsessionClass = NSClassFromString(@"NSURLSessionDataTask");
    if (nsurlsessionClass) {
      NSURLSessionConfiguration *sessionConfiguration = [NSURLSessionConfiguration defaultSessionConfiguration];
      __block NSURLSession *session = [NSURLSession sessionWithConfiguration:sessionConfiguration];
      
      NSURLSessionDataTask *task = [session dataTaskWithRequest:request
                                              completionHandler: ^(NSData *data, NSURLResponse *response, NSError *error) {
                                                dispatch_async(dispatch_get_main_queue(), ^{
                                                  typeof (self) strongSelf = weakSelf;
                                                  
                                                  [session finishTasksAndInvalidate];

                                                  [strongSelf previewPanel:blockPanel updateAttachment:strongSelf.previewAttachment data:data];
                                                });
                                              }];
      [task resume];
    }else{
#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Wdeprecated-declarations"
      [NSURLConnection sendAsynchronousRequest:request queue:self.thumbnailQueue completionHandler:^(NSURLResponse *response, NSData *responseData, NSError *err) {
#pragma clang diagnostic pop
        typeof (self) strongSelf = weakSelf;
        [strongSelf previewPanel:blockPanel updateAttachment:strongSelf.previewAttachment data:responseData];
      }];
    }
  }
  
  return self.previewAttachment;
}

- (void)previewPanel:(QLPreviewPanel *)panel updateAttachment:(BITFeedbackMessageAttachment *)attachment data:(NSData *)data {
  self.previewAttachment.isLoading = NO;
  if (data.length) {
    [self.previewAttachment replaceData:data];
    [panel reloadData];
    
    [[BITHockeyManager sharedHockeyManager].feedbackManager saveMessages];
  } else {
    [panel reloadData];
  }
}


#pragma mark - QLPreviewPanelDelegate

- (BOOL)previewPanel:(QLPreviewPanel *)panel handleEvent:(NSEvent *)event {
  // redirect all key down events to the table view
  if ([event type] == NSKeyDown) {
    [self.feedbackTableView keyDown:event];
    return YES;
  }
  return NO;
}

- (NSRect)previewPanel:(QLPreviewPanel *)panel sourceFrameOnScreenForPreviewItem:(id <QLPreviewItem>)item {
  BOOL memberOfComposeAttachments = [self.attachments containsObject:self.previewAttachment];
  NSTableView *relevantTableView = (memberOfComposeAttachments) ? self.feedbackAttachmentsTableView : self.feedbackTableView;
  
  NSRect visibleRect = [relevantTableView visibleRect];

  if (!NSIntersectsRect(visibleRect, self.previewThumbnailRect)) {
    return NSZeroRect;
  }
  
  NSRect thumbnailRect = [relevantTableView convertRect:self.previewThumbnailRect toView:nil];
  thumbnailRect = [self.window convertRectToScreen:thumbnailRect];
  
  return thumbnailRect;
}

@end

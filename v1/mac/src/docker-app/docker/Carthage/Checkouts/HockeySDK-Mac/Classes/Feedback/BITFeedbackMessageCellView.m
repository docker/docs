#import "BITFeedbackMessageCellView.h"

#import "BITFeedbackMessageCellViewDelegate.h"

#import "BITFeedbackMessage.h"
#import "BITFeedbackMessageAttachment.h"

#import "HockeySDKPrivate.h"
#import "BITSDKColoredView.h"
#import "BITActivityIndicatorButton.h"

#define BACKGROUNDCOLOR_DEFAULT BIT_RGBCOLOR(245, 245, 245)
#define BACKGROUNDCOLOR_ALTERNATE BIT_RGBCOLOR(235, 235, 235)

#define TEXTCOLOR_TITLE BIT_RGBCOLOR(75, 75, 75)

#define TEXTCOLOR_DEFAULT BIT_RGBCOLOR(25, 25, 25)
#define TEXTCOLOR_PENDING BIT_RGBCOLOR(75, 75, 75)

#define TEXT_FONTSIZE 13
#define DATE_FONTSIZE 11

#define FRAME_SIDE_BORDER 10
#define FRAME_TOP_BORDER 23
#define FRAME_BOTTOM_BORDER 23
#define FRAME_LEFT_RESPONSE_BORDER 20

#define LABEL_TEXT_Y 17


@interface BITFeedbackMessageCellView()

@property (nonatomic, unsafe_unretained) id<BITFeedbackMessageCellViewDelegate> bitDelegate;

@end


@implementation BITFeedbackMessageCellView {
  NSDateFormatter *_dateFormatter;
  NSDateFormatter *_timeFormatter;

  NSInteger _row;
  
  NSInteger _attachmentsAdded;
}


#pragma mark - NSObject

- (instancetype)initWithFrame:(NSRect)frameRect delegate:(id<BITFeedbackMessageCellViewDelegate>)delegate {
  self = [super initWithFrame:frameRect];
  if (self) {
    _attachmentsAdded = 0;
    _bitDelegate = delegate;
    
    // Create the views we'll be using in the following
    _messageTextField = [[NSTextField alloc] initWithFrame:NSMakeRect(10, 20, frameRect.size.width - 20, frameRect.size.height - 40)];
    _dateTextField = [[NSTextField alloc] initWithFrame:NSMakeRect(10, 17, frameRect.size.width - 20, 17)];
    
    // Set up common textField properties
    _messageTextField.drawsBackground = NO;
    _dateTextField.drawsBackground = NO;
    
    _messageTextField.editable = NO;
    _dateTextField.editable = NO;
    
    _messageTextField.bordered = NO;
    _dateTextField.bordered = NO;
    
    _messageTextField.autoresizingMask = NSViewWidthSizable | NSViewHeightSizable;
    _dateTextField.autoresizingMask = NSViewWidthSizable | NSViewMaxYMargin;
    
    // Style the text fields a bit
    _messageTextField.font = [NSFont systemFontOfSize:[NSFont systemFontSize]];
    _dateTextField.font = [NSFont systemFontOfSize:[NSFont smallSystemFontSize]];
    
    _messageTextField.textColor = [NSColor colorWithCalibratedRed:0.25 green:0.29 blue:0.30 alpha:1.0];
    _dateTextField.textColor = [NSColor colorWithCalibratedRed:0.61 green:0.66 blue:0.70 alpha:1.0];
    _dateTextField.textColor = [NSColor darkGrayColor];
    
    // Add the views as subviews
    [self addSubview:_messageTextField];
    [self addSubview:_dateTextField];
    
    // bottom border
    BITSDKColoredView *bottomBorder = [[BITSDKColoredView alloc] initWithFrame:NSMakeRect(frameRect.origin.x, 1, frameRect.size.width, 1)];
    bottomBorder.viewBackgroundColor = [NSColor colorWithCalibratedRed:0.82 green:0.84 blue:0.85 alpha:1.0];
    bottomBorder.autoresizingMask = NSViewWidthSizable | NSViewMaxYMargin;
    [self addSubview:bottomBorder];
    
    [[NSNotificationCenter defaultCenter] addObserver:self
                                             selector:@selector(updateAttachment:)
                                                 name:kBITFeedbackAttachmentLoadedNotification
                                               object:nil];
  }
  
  return self;
}

- (void)dealloc {
  [[NSNotificationCenter defaultCenter] removeObserver:self];
}

#pragma mark - Layout

- (void)updateAttachment:(NSNotification *)notification {
  if (!notification.userInfo) return;
  NSDictionary *dict = notification.userInfo;
  BITFeedbackMessageAttachment *attachment = dict[kBITFeedbackAttachmentLoadedKey];
  if (!attachment) return;
  
  if (![self.message.attachments containsObject:attachment]) return;
  
  for (NSView *subview in self.subviews) {
    if ([subview isKindOfClass:[BITActivityIndicatorButton class]]) {
      BITActivityIndicatorButton *button = (BITActivityIndicatorButton *)subview;
      if (button.tag == [[attachment identifier] integerValue]) {
        [button setShowsActivityIndicator:NO];
        [button setImage:[attachment thumbnailWithSize:CGSizeMake(BIT_ATTACHMENT_THUMBNAIL_LENGTH, BIT_ATTACHMENT_THUMBNAIL_LENGTH)]];
      }
    }
  }
}

- (void)setMessage:(BITFeedbackMessage *)message {
  _message = message;
  
  if (self.message.userMessage) {
    self.messageTextField.alignment = NSRightTextAlignment;
    self.dateTextField.alignment = NSRightTextAlignment;
  } else {
    self.messageTextField.alignment = NSLeftTextAlignment;
    self.dateTextField.alignment = NSLeftTextAlignment;
  }
  
  self.messageTextField.stringValue = message.text;
  NSValueTransformer *valueTransformer = [NSValueTransformer valueTransformerForName:@"BITFeedbackMessageDateValueTransformer"];
  self.dateTextField.stringValue = [valueTransformer transformedValue:message];
  
  [self setNeedsDisplay:YES];
}

- (void)updateAttachmentViews {
  _attachmentsAdded = 0;
  
  NSArray *previewableAttachments = self.message.previewableAttachments;

  if (_attachmentsAdded == [previewableAttachments count]) return;
  
  if (previewableAttachments) {
    CGFloat baseOffsetOfText = CGRectGetMaxY(self.dateTextField.frame) + 10;
    
    NSInteger i = 0;
    
    CGFloat attachmentsPerRow = floorf(self.frame.size.width / (FRAME_SIDE_BORDER + BIT_ATTACHMENT_THUMBNAIL_LENGTH));
    
    for (BITFeedbackMessageAttachment *attachment in self.message.attachments) {
      attachment.identifier = [NSNumber numberWithInteger:i];
      
      NSRect frame;
      if (self.message.userMessage) {
        frame = CGRectMake(self.frame.size.width - FRAME_SIDE_BORDER - BIT_ATTACHMENT_THUMBNAIL_LENGTH -  ((FRAME_SIDE_BORDER + BIT_ATTACHMENT_THUMBNAIL_LENGTH) *  (i%(int)attachmentsPerRow) ), floor(i/attachmentsPerRow)*(FRAME_SIDE_BORDER + BIT_ATTACHMENT_THUMBNAIL_LENGTH) + baseOffsetOfText , BIT_ATTACHMENT_THUMBNAIL_LENGTH, BIT_ATTACHMENT_THUMBNAIL_LENGTH);
      } else {
        frame = CGRectMake(FRAME_SIDE_BORDER + (FRAME_SIDE_BORDER + BIT_ATTACHMENT_THUMBNAIL_LENGTH) * (i%(int)attachmentsPerRow) , floor(i/attachmentsPerRow)*(FRAME_SIDE_BORDER + BIT_ATTACHMENT_THUMBNAIL_LENGTH) + baseOffsetOfText , BIT_ATTACHMENT_THUMBNAIL_LENGTH, BIT_ATTACHMENT_THUMBNAIL_LENGTH);
      }
      
      BITActivityIndicatorButton *imageButton = [[BITActivityIndicatorButton alloc] initWithFrame:frame];
      
      imageButton.title = @"";
      imageButton.imagePosition = NSImageOnly;
      [imageButton setBordered:NO];
      
      [imageButton setTag:[attachment.identifier integerValue]];
      [imageButton setTarget:self];
      [imageButton setAction:@selector(previewAttachment:)];
      
      if (attachment.localURL){
        [imageButton setImage:[attachment thumbnailWithSize:CGSizeMake(BIT_ATTACHMENT_THUMBNAIL_LENGTH, BIT_ATTACHMENT_THUMBNAIL_LENGTH)]];
        [imageButton setShowsActivityIndicator:NO];
      } else {
        [imageButton setImage:nil];
        [imageButton setShowsActivityIndicator:YES];
      }
      
      if (self.message.userMessage) {
        imageButton.autoresizingMask = NSViewMinXMargin | NSViewMaxYMargin;
      } else {
        [imageButton setBitBackgroundColor:[NSColor colorWithCalibratedRed:0.93 green:0.94 blue:0.95 alpha:1]];
        imageButton.autoresizingMask = NSViewMaxXMargin | NSViewMaxYMargin;
      }
      
      [self addSubview:imageButton];
      
      _attachmentsAdded++;
      
      i++;
    }
  }
}

- (void)drawRect:(NSRect)dirtyRect {
  NSColor *backgroundColor = [NSColor colorWithCalibratedRed:0.93 green:0.94 blue:0.95 alpha:1];
  
  if (self.message.userMessage) {
    backgroundColor = [NSColor whiteColor];
  }
  
  [backgroundColor setFill];
  NSRectFill(dirtyRect);
  
  [super drawRect:dirtyRect];
}


+ (NSRect)messageUsedRect:(BITFeedbackMessage *)message tableViewWidth:(CGFloat)width {
  CGRect maxMessageHeightFrame = CGRectMake(0, 0, width - FRAME_SIDE_BORDER * 2 - 4, CGFLOAT_MAX);
  
  NSTextStorage *textStorage = [[NSTextStorage alloc] initWithString:message.text];
  NSTextContainer *textContainer = [[NSTextContainer alloc] initWithContainerSize:NSSizeFromCGSize(maxMessageHeightFrame.size)];
  NSLayoutManager *layoutManager = [[NSLayoutManager alloc] init];
  
  [layoutManager addTextContainer:textContainer];
  [textStorage addLayoutManager:layoutManager];
  
  [textStorage setAttributes:@{NSFontAttributeName: [NSFont systemFontOfSize:TEXT_FONTSIZE]}
                       range:NSMakeRange(0, [textStorage length])];
  [textContainer setLineFragmentPadding:0.0];
  
  (void)[layoutManager glyphRangeForTextContainer:textContainer];
  NSRect aRect = [layoutManager usedRectForTextContainer:textContainer];
  
  CGFloat attachmentsPerRow = floorf(width / (FRAME_SIDE_BORDER + BIT_ATTACHMENT_THUMBNAIL_LENGTH));
  CGFloat attachmentHeight = BIT_ATTACHMENT_THUMBNAIL_LENGTH * ceil([message previewableAttachments].count / attachmentsPerRow);
  
  if (attachmentHeight > 0) attachmentHeight += 10;
  
  aRect.size.height += attachmentHeight + FRAME_TOP_BORDER + LABEL_TEXT_Y + FRAME_BOTTOM_BORDER;
  
  return aRect;
}


- (void)previewAttachment:(id)sender {
  BITFeedbackMessageAttachment *attachment = nil;
  
  for (BITFeedbackMessageAttachment *theAttachment in self.message.attachments) {
    if (![theAttachment isKindOfClass:[BITFeedbackMessageAttachment class]]) continue;
    
    if ([theAttachment.identifier integerValue] == [(NSButton *)sender tag]) {
      attachment = theAttachment;
      break;
    }
  }
  
  if (attachment && self.bitDelegate) {
    if ([self.bitDelegate respondsToSelector:@selector(messageCellView:clickOnButton:withAttachment:)]) {
      [self.bitDelegate messageCellView:self clickOnButton:sender withAttachment:attachment];
    }
  }
}


#pragma mark - Public

+ (NSString *)identifier {
  return NSStringFromClass([self class]);
}


+ (CGFloat) heightForRowWithMessage:(BITFeedbackMessage *)message tableViewWidth:(CGFloat)width {
  return [[self class] messageUsedRect:message tableViewWidth:width].size.height;
}

@end

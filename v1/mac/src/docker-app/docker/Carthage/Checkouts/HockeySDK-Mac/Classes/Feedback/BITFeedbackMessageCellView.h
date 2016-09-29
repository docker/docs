#import <Cocoa/Cocoa.h>

@class BITFeedbackMessage;
@protocol BITFeedbackMessageCellViewDelegate;

@interface BITFeedbackMessageCellView : NSTableCellView

@property (nonatomic, strong) BITFeedbackMessage *message;
@property (nonatomic, strong) NSTextField *messageTextField;
@property (nonatomic, strong) NSTextField *dateTextField;

- (instancetype)initWithFrame:(NSRect)frameRect delegate:(id<BITFeedbackMessageCellViewDelegate>)delegate;

+ (NSString *)identifier;
+ (CGFloat) heightForRowWithMessage:(BITFeedbackMessage *)message tableViewWidth:(CGFloat)width;

- (void)updateAttachmentViews;

@end

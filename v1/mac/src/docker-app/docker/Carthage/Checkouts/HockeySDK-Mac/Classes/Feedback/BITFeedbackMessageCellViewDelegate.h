#import <Foundation/Foundation.h>

@class BITFeedbackMessageCellView;
@class BITFeedbackMessageAttachment;

@protocol BITFeedbackMessageCellViewDelegate <NSObject>

- (void)messageCellView:(BITFeedbackMessageCellView *)messaggeCellView clickOnButton:(NSButton *)button withAttachment:(BITFeedbackMessageAttachment *)attachment;

@end

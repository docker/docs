#import <Foundation/Foundation.h>

@class BITSDKTextView;

@protocol BITSDKTextViewDelegate <NSObject>

- (void)textView:(BITSDKTextView *)textView dragOperationWithFilename:(NSString *)filename;

@end

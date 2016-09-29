#import <Cocoa/Cocoa.h>

@interface BITSDKTextFieldCell : NSTextFieldCell

@property (nonatomic, strong) NSNumber *horizontalInset;

- (void)setBitPlaceHolderString:(NSString *)bitPlaceHolderString;

@end

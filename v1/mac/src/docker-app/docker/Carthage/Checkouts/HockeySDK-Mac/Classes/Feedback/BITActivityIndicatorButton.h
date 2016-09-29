#import <Cocoa/Cocoa.h>

@interface BITActivityIndicatorButton : NSButton

@property (nonatomic, strong) NSColor *bitBackgroundColor;

- (void)setShowsActivityIndicator:(BOOL)showsIndicator;

@end

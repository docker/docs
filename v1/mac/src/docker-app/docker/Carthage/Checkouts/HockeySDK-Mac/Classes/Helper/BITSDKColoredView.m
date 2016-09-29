#import "BITSDKColoredView.h"

@implementation BITSDKColoredView

- (void)drawRect:(NSRect)dirtyRect {
  if (self.viewBackgroundColor) {
    [self.viewBackgroundColor setFill];
    NSRectFill(dirtyRect);
  }
  
  if (self.viewBorderWidth > 0 && self.viewBorderColor) {
    [self setWantsLayer:YES];
    self.layer.masksToBounds = YES;
    self.layer.borderWidth = self.viewBorderWidth;
    
    // Convert to CGColorRef
    NSInteger numberOfComponents = [self.viewBorderColor numberOfComponents];
    CGFloat components[numberOfComponents];
    CGColorSpaceRef colorSpace = [[self.viewBorderColor colorSpace] CGColorSpace];
    [self.viewBorderColor getComponents:(CGFloat *)&components];
    CGColorRef orangeCGColor = CGColorCreate(colorSpace, components);
    
    self.layer.borderColor = orangeCGColor;
    CGColorRelease(orangeCGColor);    
  }
  
  [super drawRect:dirtyRect];
}

@end

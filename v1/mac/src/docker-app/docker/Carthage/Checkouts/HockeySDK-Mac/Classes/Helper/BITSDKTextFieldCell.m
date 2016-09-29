#import "BITSDKTextFieldCell.h"

@implementation BITSDKTextFieldCell

- (NSRect)drawingRectForBounds:(NSRect)theRect {
	// Get the parent's idea of where we should draw
	NSRect newRect = [super drawingRectForBounds:theRect];
  NSSize textSize = [self cellSizeForBounds:theRect];
  
  float heightDelta = newRect.size.height - textSize.height;
  if (heightDelta > 0) {
    newRect.size.height -= heightDelta;
    newRect.origin.y += heightDelta / 2;
    if (self.horizontalInset) {
      newRect.origin.x += [self.horizontalInset floatValue];
      newRect.size.width -= ([self.horizontalInset floatValue] * 2);
    }
  }
	
	return newRect;
}

- (void)setBitPlaceHolderString:(NSString *)bitPlaceHolderString {
  self.placeholderString = bitPlaceHolderString;
}

@end

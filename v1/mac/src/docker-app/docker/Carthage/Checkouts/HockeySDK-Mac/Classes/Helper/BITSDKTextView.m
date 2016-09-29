#import "BITSDKTextView.h"

@implementation BITSDKTextView

- (void)drawRect:(NSRect)rect {
  [super drawRect:rect];
  if ([[self string] isEqualToString:@""] && self != [[self window] firstResponder]) {
    if (self.placeHolderString) {
      NSColor *txtColor = [NSColor colorWithCalibratedRed:0.69 green:0.71 blue:0.73 alpha:1.0];
      NSDictionary *dict = @{NSForegroundColorAttributeName: txtColor};
      NSAttributedString *placeholder = [[NSAttributedString alloc] initWithString:self.placeHolderString attributes:dict];
      [placeholder drawAtPoint:NSMakePoint(0,0)];
    }
  }
}

- (BOOL)becomeFirstResponder {
  [self setNeedsDisplay:YES];
  return [super becomeFirstResponder];
}

- (BOOL)resignFirstResponder {
  [self setNeedsDisplay:YES];
  return [super resignFirstResponder];
}


#pragma mark - Drag & Drop for Attachments

- (NSDragOperation)draggingEntered:(id<NSDraggingInfo>)sender {
  NSPasteboard *pb = [sender draggingPasteboard];
  NSDragOperation dragOperation = [sender draggingSourceOperationMask];
  
  if ([[pb types] containsObject:NSFilenamesPboardType]) {
    if (dragOperation & NSDragOperationCopy) {
      return NSDragOperationCopy;
    }
  }
  
  return NSDragOperationNone;
}

- (BOOL)performDragOperation:(id<NSDraggingInfo>)sender {
  NSPasteboard *pb = [sender draggingPasteboard];
  
  if ( [[pb types] containsObject:NSFilenamesPboardType] ) {
    NSFileManager *fm = [[NSFileManager alloc] init];
    
    NSArray *filenames = [pb propertyListForType:NSFilenamesPboardType];
    
    BOOL fileFound = NO;
    
    for (NSString *filename in filenames) {
      BOOL isDir = NO;
      if (![fm fileExistsAtPath:filename isDirectory:&isDir] || isDir) continue;
      
      fileFound = YES;
      
      if (self.bitDelegate && [self.bitDelegate respondsToSelector:@selector(textView:dragOperationWithFilename:)]) {
        [self.bitDelegate textView:self dragOperationWithFilename:filename];
      }
    }
    return fileFound;
  }
  
  return NO;
}

@end

#import "BITActivityIndicatorButton.h"


@interface BITActivityIndicatorButton()

@property (nonatomic, strong) NSProgressIndicator *indicator;
@property (nonatomic) BOOL indicatorVisible;

@end


@implementation BITActivityIndicatorButton

- (instancetype)initWithFrame:(NSRect)frameRect {
  if (self = [super initWithFrame:frameRect]) {
    _indicator = [[NSProgressIndicator alloc] initWithFrame:self.bounds];
    
    [_indicator setStyle: NSProgressIndicatorSpinningStyle];
    [_indicator setControlSize: NSSmallControlSize];
    [_indicator sizeToFit];
    
    _indicator.hidden = YES;
    
    [self addSubview:_indicator];
  }
  return self;
}

- (void)setShowsActivityIndicator:(BOOL)showsIndicator {
  if (self.indicatorVisible == showsIndicator){
    return;
  }
  
  self.indicatorVisible = showsIndicator;
  [[self cell] setBackgroundColor:self.bitBackgroundColor];
  
  if (showsIndicator){
    [self.indicator startAnimation:self];
    [self.indicator setHidden:NO];
    self.image = nil;
  } else {
    [self.indicator stopAnimation:self];
    [self.indicator setHidden:YES];
  }
}


@end

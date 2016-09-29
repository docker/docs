#import "BITHTTPOperation.h"

@interface BITHTTPOperation()<NSURLConnectionDelegate>
@end

@implementation BITHTTPOperation {
  NSURLRequest *_URLRequest;
  NSURLConnection *_connection;
  NSMutableData *_data;
  
  BOOL _isExecuting;
  BOOL _isFinished;
}


+ (instancetype)operationWithRequest:(NSURLRequest *)urlRequest {
  BITHTTPOperation *op = [[self class] new];
  op->_URLRequest = urlRequest;
  return op;
}

#pragma mark - NSOperation overrides
- (BOOL)isConcurrent {
  return YES;
}

- (void)cancel {
  [_connection cancel];
  [super cancel];
}

- (void) start {
  if(self.isCancelled) {
    [self finish];
    return;
  }
  
  if (![[NSThread currentThread] isMainThread]) {
    [self performSelector:@selector(start) onThread:NSThread.mainThread withObject:nil waitUntilDone:NO];
    return;
  }
  
  if(self.isCancelled) {
    [self finish];
    return;
  }

  [self willChangeValueForKey:@"isExecuting"];
  _isExecuting = YES;
  [self didChangeValueForKey:@"isExecuting"];
  
#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Wdeprecated-declarations"
  _connection = [[NSURLConnection alloc] initWithRequest:_URLRequest
                                                delegate:self
                                        startImmediately:YES];
#pragma clang diagnostic pop
}

- (void) finish {
  [self willChangeValueForKey:@"isExecuting"];
  [self willChangeValueForKey:@"isFinished"];
  _isExecuting = NO;
  _isFinished = YES;
  [self didChangeValueForKey:@"isExecuting"];
  [self didChangeValueForKey:@"isFinished"];
}

#pragma mark - NSURLConnectionDelegate
-(void)connection:(NSURLConnection*)connection didReceiveResponse:(NSURLResponse*)response {
  _data = [[NSMutableData alloc] init];
  _response = (id)response;
}

-(void)connection:(NSURLConnection*)connection didReceiveData:(NSData*)data {
  [_data appendData:data];
}

-(void)connection:(NSURLConnection*)connection didFailWithError:(NSError*)error {
  //FINISHED and failed
  _error = error;
  _data = nil;
  
  [self finish];
}

-(void)connectionDidFinishLoading:(NSURLConnection*)connection {
  [self finish];
}

#pragma mark - Public interface
- (NSData *)data {
  return _data;
}

- (void)setCompletion:(BITNetworkCompletionBlock)completion {
  if(!completion) {
    [super setCompletionBlock:nil];
  } else {
    __unsafe_unretained typeof(self) weakSelf = self;
    [super setCompletionBlock:^{
      typeof(self) strongSelf = weakSelf;
      if(strongSelf) {
        dispatch_async(dispatch_get_main_queue(), ^{
          if(!strongSelf.isCancelled) {
            completion(strongSelf, strongSelf->_data, strongSelf->_error);
          }
          [strongSelf setCompletionBlock:nil];
        });
      }
    }];
  }
}

- (BOOL)isFinished {
  return _isFinished;
}

- (BOOL)isExecuting {
  return _isExecuting;
}

@end

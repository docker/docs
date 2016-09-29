#import <Foundation/Foundation.h>

@class BITHTTPOperation;
typedef void (^BITNetworkCompletionBlock)(BITHTTPOperation* operation, NSData* data, NSError* error);

@interface BITHTTPOperation : NSOperation

+ (instancetype) operationWithRequest:(NSURLRequest *) urlRequest;

@property (nonatomic, readonly) NSURLRequest *URLRequest;

//the completion is only called if the operation wasn't cancelled
- (void) setCompletion:(BITNetworkCompletionBlock) completionBlock;

@property (nonatomic, readonly) NSHTTPURLResponse *response;
@property (weak, nonatomic, readonly) NSData *data;
@property (nonatomic, readonly) NSError *error;

@end

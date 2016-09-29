#import <Foundation/Foundation.h>
#import "HockeySDK.h"

#import "HockeySDKNullability.h"

@class BITOrderedDictionary;
@class BITConfiguration;
@class BITTelemetryData;
@class BITTelemetryContext;
@class BITPersistence;
NS_ASSUME_NONNULL_BEGIN

FOUNDATION_EXPORT char *BITSafeJsonEventsString;

/**
 *  Items get queued before they are persisted and sent out as a batch. This class managed the queue, and forwards the batch
 *  to the persistence layer once the max batch count has been reached.
 */
@interface BITChannel : NSObject

/**
 *  Telemetry context used by the channel to create the payload (testing).
 */
@property (nonatomic, strong) BITTelemetryContext *telemetryContext;

/**
 *  Persistence instance for storing files after the queue gets flushed (testing).
 */
@property (nonatomic, strong) BITPersistence *persistence;

/**
 *  Number of queue items which will trigger a flush (testing).
 */
@property (nonatomic) NSUInteger maxBatchCount;

/**
 *  A queue which makes array operations thread safe.
 */
@property (nonatomic, assign) dispatch_queue_t dataItemsOperations;

/**
 *  An integer value that keeps tracks of the number of data items added to the JSON Stream string.
 */
@property (nonatomic, assign) NSUInteger dataItemCount;

/**
 *  Initializes a new BITChannel instance.
 *
 *  @param telemetryContext the context used to add context values to the metrics payload
 *  @param persistence the persistence used to save metrics after the queue gets flushed
 *
 *  @return the telemetry context
 */
- (instancetype)initWithTelemetryContext:(BITTelemetryContext *)telemetryContext persistence:(BITPersistence *) persistence;

/**
 *  Enqueue telemetry data (events, metrics, exceptions, traces) before processing it.
 *
 *  @param item the telemetry object, which should be processed
 */
- (void)enqueueTelemetryItem:(BITTelemetryData *)item;

/**
 *  Manually trigger the BITChannel to persist all items currently in its data item queue.
 */
- (void)persistDataItemQueue;

/**
 *  Adds the specified dictionary to the JSON Stream string.
 *
 *  @param dictionary the dictionary object which is to be added to the JSON Stream queue string.
 */
- (void)appendDictionaryToJsonStream:(BITOrderedDictionary *)dictionary;

/**
 *  A C function that serializes a given dictionary to JSON and appends it to a char string
 *
 *  @param dictionary A dictionary which will be serialized to JSON and then appended to the string.
 *  @param string The C string which the dictionary's JSON representation will be appended to.
 */
void bit_appendStringToSafeJsonStream(NSString *string, char *__nonnull*__nonnull jsonStream);

/**
 *  Reset BITSafeJsonEventsString so we can start appending JSON dictionaries.
 *
 *  @param string The string that will be reset.
 */
void bit_resetSafeJsonStream(char *__nonnull*__nonnull jsonStream);

/**
 *  A method which indicates whether the telemetry pipeline is busy and no new data should be enqueued.
 *  Currently, we drop telemetry data if this returns YES.
 *  This depends on defaultMaxBatchCount and defaultBatchInterval.
 *
 *  @return Returns yes if currently no new data should be enqueued on the channel.
 */
- (BOOL)isQueueBusy;

@end
NS_ASSUME_NONNULL_END

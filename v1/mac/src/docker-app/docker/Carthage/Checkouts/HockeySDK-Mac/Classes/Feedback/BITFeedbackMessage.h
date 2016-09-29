#import <Foundation/Foundation.h>

@class BITFeedbackMessageAttachment;

/**
 *  Status for each feedback message
 */
typedef NS_ENUM(NSInteger, BITFeedbackMessageStatus) {
  /**
   *  default and new messages from SDK per default
   */
  BITFeedbackMessageStatusSendPending = 0,
  /**
   *  message is in conflict, happens if the message is already stored on the server and tried sending it again
   */
  BITFeedbackMessageStatusInConflict = 1,
  /**
   *  sending of message is in progress
   */
  BITFeedbackMessageStatusSendInProgress = 2,
  /**
   *  new messages from server
   */
  BITFeedbackMessageStatusUnread = 3,
  /**
   *  messages from server once read and new local messages once successful send from SDK
   */
  BITFeedbackMessageStatusRead = 4,
  /**
   *  message is archived, happens if the thread is deleted from the server
   */
  BITFeedbackMessageStatusArchived = 5
};


/**
 *  An individual feedback message
 */
@interface BITFeedbackMessage : NSObject <NSCopying>

@property (nonatomic, copy) NSString *text;
@property (nonatomic, copy) NSString *userID;
@property (nonatomic, copy) NSString *name;
@property (nonatomic, copy) NSString *email;
@property (nonatomic, copy) NSDate *date;
@property (nonatomic, copy) NSNumber *messageID;
@property (nonatomic, copy) NSString *token;
@property (nonatomic, strong) NSArray *attachments;
@property (nonatomic) BITFeedbackMessageStatus status;
@property (nonatomic) BOOL userMessage;

/**
 Delete local cached attachment data
 
 @warning This method must be called before a feedback message is deleted.
 */
- (void)deleteContents;

/**
 Add an attachment to a message
 
 @param object BITFeedbackMessageAttachment instance representing the attachment that should be added
 */
- (void)addAttachmentsObject:(BITFeedbackMessageAttachment *)object;

/**
 Return the attachments that can be viewed
 
 @return NSArray containing the attachment objects that can be previewed
 */
- (NSArray *)previewableAttachments;

@end

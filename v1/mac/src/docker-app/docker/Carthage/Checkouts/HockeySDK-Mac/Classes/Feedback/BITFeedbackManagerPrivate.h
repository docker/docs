#import <Foundation/Foundation.h>
#import "BITFeedbackManager.h"
#import "BITFeedbackMessage.h"

@interface BITFeedbackManager () {
}


@property (nonatomic, strong) NSMutableArray *feedbackList;
@property (nonatomic, strong) NSString *token;


// used by BITHockeyManager if disable status is changed
@property (nonatomic, getter = isFeedbackManagerDisabled) BOOL disableFeedbackManager;

@property (nonatomic) BOOL didAskUserData;

@property (nonatomic, strong) NSDate *lastCheck;
@property (nonatomic, strong) NSDate *lastRefreshDate;

@property (nonatomic, strong) NSNumber *lastMessageID;

//@property (nonatomic, copy) NSString *userID;
//@property (nonatomic, copy) NSString *userName;
//@property (nonatomic, copy) NSString *userEmail;


// load new messages from the server
- (void)updateMessagesList;

// load new messages from the server if the last request is too long ago
- (void)updateMessagesListIfRequired;

- (NSUInteger)numberOfMessages;
- (BITFeedbackMessage *)messageAtIndex:(NSUInteger)index;

- (void)submitMessageWithText:(NSString *)text andAttachments:(NSArray *)photos;
- (void)submitPendingMessages;

// Returns YES if manual user data can be entered, required or optional
- (BOOL)askManualUserDataAvailable;

// Returns YES if required user data is missing?
- (BOOL)requireManualUserDataMissing;

// Returns YES if optional user data is missing
- (BOOL)optionalManualUserDataMissing;

// Returns YES if user data is available and can be edited
- (BOOL)isManualUserDataAvailable;

// used in the user data screen
- (void)updateDidAskUserData;


- (BITFeedbackMessage *)messageWithID:(NSNumber *)messageID;

- (NSArray *)messagesWithStatus:(BITFeedbackMessageStatus)status;

- (void)saveMessages;

- (void)fetchMessageUpdates;
- (void)updateMessageListFromResponse:(NSDictionary *)jsonDictionary;

- (BOOL)deleteMessageAtIndex:(NSUInteger)index;
- (void)deleteAllMessages;

@end

#import "BITFeedbackMessage.h"
#import "BITFeedbackMessageAttachment.h"

#import "HockeySDKPrivate.h"
#import <QuickLook/QuickLook.h>


@implementation BITFeedbackMessage

#pragma mark - NSObject

- (id) init {
  if ((self = [super init])) {
    _text = nil;
    _userID = nil;
    _name = nil;
    _email = nil;
    _date = [[NSDate alloc] init];
    _token = nil;
    _attachments = nil;
    _messageID = @0;
    _status = BITFeedbackMessageStatusSendPending;
    _userMessage = NO;
  }
  return self;
}


- (id)copyWithZone:(NSZone *)zone {
  BITFeedbackMessage *copy = [[[self class] allocWithZone: zone] init];
  
  [copy setText: _text];
  [copy setUserID: _userID];
  [copy setName: _name];
  [copy setEmail: _email];
  [copy setDate: _date];
  [copy setToken: _token];
  [copy setMessageID: _messageID];
  [copy setStatus: _status];
  [copy setUserMessage: _userMessage];
  [copy setAttachments: _attachments];
  
  return copy;
}


#pragma mark - NSCoder

- (void)encodeWithCoder:(NSCoder *)encoder {
  [encoder encodeObject:self.text forKey:@"text"];
  [encoder encodeObject:self.userID forKey:@"userID"];
  [encoder encodeObject:self.name forKey:@"name"];
  [encoder encodeObject:self.email forKey:@"email"];
  [encoder encodeObject:self.date forKey:@"date"];
  [encoder encodeObject:self.messageID forKey:@"messageID"];
  [encoder encodeObject:self.attachments forKey:@"attachments"];
  [encoder encodeInteger:self.status forKey:@"status"];
  [encoder encodeBool:self.userMessage forKey:@"userMessage"];
  [encoder encodeObject:self.token forKey:@"token"];
}

- (id)initWithCoder:(NSCoder *)decoder {
  if ((self = [super init])) {
    self.text = [decoder decodeObjectForKey:@"text"];
    self.userID = [decoder decodeObjectForKey:@"userID"];
    self.name = [decoder decodeObjectForKey:@"name"];
    self.email = [decoder decodeObjectForKey:@"email"];
    self.date = [decoder decodeObjectForKey:@"date"];
    self.messageID = [decoder decodeObjectForKey:@"messageID"];
    self.attachments = [decoder decodeObjectForKey:@"attachments"];
    self.status = (BITFeedbackMessageStatus)[decoder decodeIntegerForKey:@"status"];
    self.userMessage = [decoder decodeBoolForKey:@"userMessage"];
    self.token = [decoder decodeObjectForKey:@"token"];
  }
  return self;
}


#pragma mark - Deletion

- (void)deleteContents {
  for (BITFeedbackMessageAttachment *attachment in self.attachments){
    [attachment deleteContents];
  }
}

- (NSArray *)previewableAttachments {
  NSMutableArray *returnArray = [NSMutableArray new];
  
  for (BITFeedbackMessageAttachment *attachment in self.attachments) {
    if (!attachment.localURL && [self userMessage]) continue;
    
    NSImage *thumbnailImage = [[NSWorkspace sharedWorkspace] iconForFileType:[attachment.originalFilename pathExtension]];
    if (!thumbnailImage) continue;
    
    if ([attachment thumbnailWithSize:CGSizeMake(BIT_ATTACHMENT_THUMBNAIL_LENGTH, BIT_ATTACHMENT_THUMBNAIL_LENGTH)]) {
      [returnArray addObject:attachment];
    }
  }
  
  return returnArray;
}

- (void)addAttachmentsObject:(BITFeedbackMessageAttachment *)object{
  if (!self.attachments) {
    self.attachments = @[];
  }
  
  if (![object isKindOfClass:[BITFeedbackMessageAttachment class]]) return;
  
  self.attachments = [self.attachments arrayByAddingObject:object];
}

@end

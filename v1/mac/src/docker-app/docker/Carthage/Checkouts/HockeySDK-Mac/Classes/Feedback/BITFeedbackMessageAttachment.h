#import <Foundation/Foundation.h>

#import <Quartz/Quartz.h>

/**
 * An individual feedback message attachment
 */
@interface BITFeedbackMessageAttachment : NSObject<NSCoding, QLPreviewItem>

@property (nonatomic, copy) NSNumber *identifier;
@property (nonatomic, copy) NSString *originalFilename;
@property (nonatomic, copy) NSString *contentType;
@property (nonatomic, copy) NSString *sourceURL;
@property (nonatomic) BOOL isLoading;
@property (nonatomic, copy, readonly) NSData *data;


@property (readonly) NSImage *thumbnailRepresentation;
@property (weak, readonly) NSImage *imageRepresentation;


+ (BITFeedbackMessageAttachment *)attachmentWithData:(NSData *)data contentType:(NSString *)contentType;

- (NSImage *)thumbnailWithSize:(NSSize)size;

- (void)replaceData:(NSData *)data;

- (void)deleteContents;

- (BOOL)needsLoadingFromURL;

- (BOOL)isImage;

- (NSURL *)localURL;

/**
 Used to determine whether QuickLook can preview this file or not. If not, we don't download it.
 */ 
- (NSString*)possibleFilename;

@end

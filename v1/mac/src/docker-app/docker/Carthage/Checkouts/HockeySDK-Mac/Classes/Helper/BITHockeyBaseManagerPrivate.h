#import <Foundation/Foundation.h>

#import "BITHockeyBaseManager.h"

@interface BITHockeyBaseManager ()

@property (nonatomic, strong) NSString *appIdentifier;
@property (nonatomic, copy) NSString *userID;
@property (nonatomic, copy) NSString *userName;
@property (nonatomic, copy) NSString *userEmail;

- (id)initWithAppIdentifier:(NSString *)appIdentifier;

- (void)startManager;

- (void)reportError:(NSError *)error;
- (NSString *)encodedAppIdentifier;

- (NSString *)getDevicePlatform;
//- (NSString *)executableUUID;

- (NSDate *)parseRFC3339Date:(NSString *)dateString;

@end

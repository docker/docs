//
//  BITPersistenceTests.m
//  HockeySDK
//
//  Created by Patrick Dinger on 24/05/16.
//
//

#import <XCTest/XCTest.h>
#import <OCMock/OCMock.h>
#import "BITPersistence.h"
#import "BITPersistencePrivate.h"

@interface BITPersistenceTests : XCTestCase

@end

@implementation BITPersistenceTests {
    BITPersistence *_subject;
}

- (void)setUp {
    [super setUp];
    _subject = [BITPersistence alloc];
    id mock = OCMPartialMock(_subject);
    
    OCMStub([mock bundleIdentifier]).andReturn(@"com.testapp");
    
    _subject = [_subject init];
}

- (void)tearDown {
    [super tearDown];
}

- (void)testAppHockeySDKDirectoryPath {
    NSString *path = [_subject appHockeySDKDirectoryPath];
    
    NSString *appSupportPath = [[NSSearchPathForDirectoriesInDomains(NSApplicationSupportDirectory, NSUserDomainMask, YES) lastObject] stringByStandardizingPath];
    NSString *validPath = [NSString stringWithFormat:@"%@/%@", appSupportPath, @"com.testapp/com.microsoft.HockeyApp"];
    
    XCTAssertEqualObjects(path, validPath);
}

@end

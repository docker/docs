// 
//  Author: Andreas Linde <mail@andreaslinde.de>
// 
//  Copyright (c) 2012-2014 HockeyApp, Bit Stadium GmbH. All rights reserved.
//  See LICENSE.txt for author information.
//
//  Permission is hereby granted, free of charge, to any person obtaining a copy
//  of this software and associated documentation files (the "Software"), to deal
//  in the Software without restriction, including without limitation the rights
//  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//  copies of the Software, and to permit persons to whom the Software is
//  furnished to do so, subject to the following conditions:
//
//  The above copyright notice and this permission notice shall be included in
//  all copies or substantial portions of the Software.
//
//  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
//  THE SOFTWARE.

#import <Foundation/Foundation.h>

#import <HockeySDK/BITHockeyManager.h>
#import <HockeySDK/BITHockeyManagerDelegate.h>

#import <HockeySDK/BITHockeyAttachment.h>

#import <HockeySDK/BITCrashManager.h>
#import <HockeySDK/BITCrashManagerDelegate.h>
#import <HockeySDK/BITCrashDetails.h>
#import <HockeySDK/BITCrashMetaData.h>
#import <HockeySDK/BITCrashExceptionApplication.h>

#import <HockeySDK/BITSystemProfile.h>

#import <HockeySDK/BITFeedbackManager.h>
#import <HockeySDK/BITFeedbackWindowController.h>


// Notification message which HockeyManager is listening to, to retry requesting updated from the server
#define BITHockeyNetworkDidBecomeReachableNotification @"BITHockeyNetworkDidBecomeReachable"

extern NSString *const __attribute__((unused)) kBITDefaultUserID;
extern NSString *const __attribute__((unused)) kBITDefaultUserName;
extern NSString *const __attribute__((unused)) kBITDefaultUserEmail;

/**
 *  HockeySDK Crash Reporter error domain
 */
typedef NS_ENUM (NSInteger, BITCrashErrorReason) {
  /**
   *  Unknown error
   */
  BITCrashErrorUnknown,
  /**
   *  API Server rejected app version
   */
  BITCrashAPIAppVersionRejected,
  /**
   *  API Server returned empty response
   */
  BITCrashAPIReceivedEmptyResponse,
  /**
   *  Connection error with status code
   */
  BITCrashAPIErrorWithStatusCode
};
extern NSString *const __attribute__((unused)) kBITCrashErrorDomain;


/**
 *  HockeySDK Feedback error domain
 */
typedef NS_ENUM(NSInteger, BITFeedbackErrorReason) {
  /**
   *  Unknown error
   */
  BITFeedbackErrorUnknown,
  /**
   *  API Server returned invalid status
   */
  BITFeedbackAPIServerReturnedInvalidStatus,
  /**
   *  API Server returned invalid data
   */
  BITFeedbackAPIServerReturnedInvalidData,
  /**
   *  API Server returned empty response
   */
  BITFeedbackAPIServerReturnedEmptyResponse,
  /**
   *  Authorization secret missing
   */
  BITFeedbackAPIClientAuthorizationMissingSecret,
  /**
   *  No internet connection
   */
  BITFeedbackAPIClientCannotCreateConnection
};
extern NSString *const __attribute__((unused)) kBITFeedbackErrorDomain;


/**
 *  HockeySDK global error domain
 */
typedef NS_ENUM(NSInteger, BITHockeyErrorReason) {
  /**
   *  Unknown error
   */
  BITHockeyErrorUnknown
};
extern NSString *const __attribute__((unused)) kBITHockeyErrorDomain;
// HockeySDK

/*Copyright (c) 2009 Extendmac, LLC. <support@extendmac.com>
 
 Permission is hereby granted, free of charge, to any person
 obtaining a copy of this software and associated documentation
 files (the "Software"), to deal in the Software without
 restriction, including without limitation the rights to use,
 copy, modify, merge, publish, distribute, sublicense, and/or sell
 copies of the Software, and to permit persons to whom the
 Software is furnished to do so, subject to the following
 conditions:
 
 The above copyright notice and this permission notice shall be
 included in all copies or substantial portions of the Software.
 
 THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
 EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
 OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
 NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
 HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
 WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
 FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
 OTHER DEALINGS IN THE SOFTWARE.
 */

#import "BITKeychainItem.h"

@interface BITKeychainItem (Private)

/*!
 @abstract Modifies the given attribute to be newValue.
 @param attributeTag The attribute's tag.
 @param newValue A pointer to the new value.
 @param newLength The length of the new value.
 */
- (void)_modifyAttributeWithTag:(SecItemAttr)attributeTag toBeValue:(void *)newValue ofLength:(UInt32)newLength;

@end

@implementation BITKeychainItem

static BOOL _logsErrors;

+ (void)lockKeychain
{
	SecKeychainLock(NULL);
}

+ (void)unlockKeychain
{
	SecKeychainUnlock(NULL, 0, NULL, NO);
}

+ (BOOL)logsErrors
{
	@synchronized (self)
	{
		return _logsErrors;
	}
	return NO;
}

+ (void)setLogsErrors:(BOOL)logsErrors
{
	@synchronized (self)
	{
		if (_logsErrors == logsErrors)
			return;
		
		_logsErrors = logsErrors;
	}
}

#pragma mark -

- (id)_initWithCoreKeychainItem:(SecKeychainItemRef)item
                       username:(NSString *)username
                       password:(NSString *)password
{
	if ((self = [super init]))
	{
		mCoreKeychainItem = item;
		mUsername = [username copy];
		mPassword = [password copy];
		
		return self;
	}
	return nil;
}

- (void)_modifyAttributeWithTag:(SecItemAttr)attributeTag toBeValue:(void *)newValue ofLength:(UInt32)newLength
{
	NSAssert(mCoreKeychainItem, @"Core keychain item is nil. You cannot modify a keychain item that is not in the keychain.");
	
	SecKeychainAttribute attributes[1];
	attributes[0].tag = attributeTag;
	attributes[0].length = newLength;
	attributes[0].data = newValue;
	
	SecKeychainAttributeList attributeList;
	attributeList.count = 1;
	attributeList.attr = attributes;
	
	SecKeychainItemModifyAttributesAndData(mCoreKeychainItem, &attributeList, 0, NULL);
}

- (void)dealloc
{
	
	if (mCoreKeychainItem)
		CFRelease(mCoreKeychainItem);
	
}

#pragma mark -
#pragma mark General Properties
@dynamic password;
- (NSString *)password
{
	@synchronized (self)
	{
		return [mPassword copy];
	}
}

- (void)setPassword:(NSString *)newPassword
{
	@synchronized (self)
	{
		if (mPassword == newPassword)
			return;
		
		mPassword = [newPassword copy];
		
		const char *newPasswordCString = [newPassword UTF8String];
		SecKeychainItemModifyAttributesAndData(mCoreKeychainItem, NULL, (UInt32)strlen(newPasswordCString), (void *)newPasswordCString);
	}
}

#pragma mark -
@dynamic username;
- (NSString *)username
{
	@synchronized (self)
	{
		return [mUsername copy];
	}
}

- (void)setUsername:(NSString *)newUsername
{
	@synchronized (self)
	{
		if (mUsername == newUsername)
			return;
		
		mUsername = [newUsername copy];
		
		const char *newUsernameCString = [newUsername UTF8String];
		[self _modifyAttributeWithTag:kSecAccountItemAttr toBeValue:(void *)newUsernameCString ofLength:(UInt32)strlen(newUsernameCString)];
	}
}

#pragma mark -
@dynamic label;
- (NSString *)label
{
	@synchronized (self)
	{
		return [mLabel copy];
	}
}

- (void)setLabel:(NSString *)newLabel
{
	@synchronized (self)
	{
		if (mLabel == newLabel)
			return;
		
		mLabel = [newLabel copy];
		
		const char *newLabelCString = [newLabel UTF8String];
		[self _modifyAttributeWithTag:kSecLabelItemAttr toBeValue:(void *)newLabelCString ofLength:(UInt32)strlen(newLabelCString)];
	}
}

#pragma mark -
#pragma mark Actions
- (void)removeFromKeychain
{
	NSAssert(mCoreKeychainItem, @"Core keychain item is nil. You cannot remove a keychain item that is not in the keychain already.");
	
	if (mCoreKeychainItem)
	{
		OSStatus resultStatus = SecKeychainItemDelete(mCoreKeychainItem);
		if (resultStatus == noErr)
		{
			CFRelease(mCoreKeychainItem);
			mCoreKeychainItem = nil;
		}
	}
}

@end

#pragma mark -
@implementation BITGenericKeychainItem

- (id)_initWithCoreKeychainItem:(SecKeychainItemRef)item
                    serviceName:(NSString *)serviceName
                       username:(NSString *)username
                       password:(NSString *)password
{
	if ((self = [super _initWithCoreKeychainItem:item username:username password:password]))
	{
		mServiceName = [serviceName copy];
		return self;
	}
	return nil;
}

+ (id)_genericKeychainItemWithCoreKeychainItem:(SecKeychainItemRef)coreKeychainItem
                                forServiceName:(NSString *)serviceName
                                      username:(NSString *)username
                                      password:(NSString *)password
{
	return [[BITGenericKeychainItem alloc] _initWithCoreKeychainItem:coreKeychainItem
                                                       serviceName:serviceName
                                                          username:username
                                                          password:password];
}


#pragma mark -
+ (BITGenericKeychainItem *)genericKeychainItemForService:(NSString *)serviceName
                                            withUsername:(NSString *)username
{
	if (!serviceName || !username)
		return nil;
	
	const char *serviceNameCString = [serviceName UTF8String];
	const char *usernameCString = [username UTF8String];
	
	UInt32 passwordLength = 0;
	char *password = nil;
	
	SecKeychainItemRef item = nil;
	OSStatus returnStatus = SecKeychainFindGenericPassword(NULL, (UInt32)strlen(serviceNameCString), serviceNameCString, (UInt32)strlen(usernameCString), usernameCString, &passwordLength, (void **)&password, &item);
	if (returnStatus != noErr || !item)
	{
		if (_logsErrors)
			NSLog(@"Error (%@) - %s", NSStringFromSelector(_cmd), GetMacOSStatusErrorString(returnStatus));

        if (password) SecKeychainItemFreeContent(NULL, password);
		return nil;
	}
	NSString *passwordString = [[NSString alloc] initWithData:[NSData dataWithBytes:password length:passwordLength] encoding:NSUTF8StringEncoding];
	SecKeychainItemFreeContent(NULL, password);
	
	return [BITGenericKeychainItem _genericKeychainItemWithCoreKeychainItem:item forServiceName:serviceName username:username password:passwordString];
}

+ (BITGenericKeychainItem *)addGenericKeychainItemForService:(NSString *)serviceName
                                               withUsername:(NSString *)username
                                                   password:(NSString *)password
{
	if (!serviceName || !username || !password)
		return nil;
	
	const char *serviceNameCString = [serviceName UTF8String];
	const char *usernameCString = [username UTF8String];
	const char *passwordCString = [password UTF8String];
	
	SecKeychainItemRef item = nil;
	OSStatus returnStatus = SecKeychainAddGenericPassword(NULL, (UInt32)strlen(serviceNameCString), serviceNameCString, (UInt32)strlen(usernameCString), usernameCString, (UInt32)strlen(passwordCString), (void *)passwordCString, &item);
	
	if (returnStatus != noErr || !item)
	{
		if (_logsErrors)
			NSLog(@"Error (%@) - %s", NSStringFromSelector(_cmd), GetMacOSStatusErrorString(returnStatus));
		return nil;
	}
	return [BITGenericKeychainItem _genericKeychainItemWithCoreKeychainItem:item forServiceName:serviceName username:username password:password];
}

#pragma mark -
#pragma mark Generic Properties
@dynamic serviceName;
- (NSString *)serviceName
{
	@synchronized (self)
	{
		return [mServiceName copy];
	}
}

- (void)setServiceName:(NSString *)newServiceName
{
	@synchronized (self)
	{
		if (mServiceName == newServiceName)
			return;
		
		mServiceName = [newServiceName copy];
		
		const char *newServiceNameCString = [newServiceName UTF8String];
		[self _modifyAttributeWithTag:kSecServiceItemAttr toBeValue:(void *)newServiceNameCString ofLength:(UInt32)strlen(newServiceNameCString)];
	}
}

@end

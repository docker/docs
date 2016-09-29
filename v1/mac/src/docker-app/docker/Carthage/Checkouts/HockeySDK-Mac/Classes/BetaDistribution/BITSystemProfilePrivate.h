#define kBITUpdateDateOfVersionInstallation	@"BITUpdateDateOfVersionInstallation"
#define kBITUpdateUsageTimeOfCurrentVersion	@"BITUpdateUsageTimeOfCurrentVersion"
#define kBITUpdateUsageTimeForVersionString	@"BITUpdateUsageTimeForVersionString"

@interface BITSystemProfile () {
}

@property (nonatomic, copy) NSDate *usageStartTimestamp;

@end
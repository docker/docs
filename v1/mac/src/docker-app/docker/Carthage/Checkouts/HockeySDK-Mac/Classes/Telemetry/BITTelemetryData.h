#import "BITTelemetryObject.h"

///Data contract class for type BITTelemetryData.
@interface BITTelemetryData : BITTelemetryObject <NSCoding>

@property (nonatomic, readonly, copy) NSString *envelopeTypeName;
@property (nonatomic, readonly, copy) NSString *dataTypeName;

@property (nonatomic, copy) NSNumber *version;
@property (nonatomic, copy) NSString *name;
@property (nonatomic, strong) NSDictionary *properties;

@end

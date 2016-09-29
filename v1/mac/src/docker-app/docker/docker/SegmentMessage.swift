import Foundation

public class SegmentMessage {
    private var userId: String
    private var message: String
    private var properties = [String: String]()
    
    public init(userId: String, message: String, channel: String, version: Version, analyticsEnabled: Bool) {
        self.userId = userId
        self.message = message
        properties["os"] = "macos"
        properties["channel"] = channel.lowercaseString
        properties["app major version"] = String(version.Major)
        properties["app minor version"] = String(version.Minor)
        properties["app patch version"] = String(version.Patch)
        properties["app version name"] = version.Name.lowercaseString
        if analyticsEnabled == true {
            properties["os major version"] = String(NSProcessInfo().operatingSystemVersion.majorVersion)
            properties["os minor version"] = String(NSProcessInfo().operatingSystemVersion.minorVersion)
            properties["os patch version"] = String(NSProcessInfo().operatingSystemVersion.patchVersion)
            if let osLanguage = Utils.getCurrentLanguage() {
                properties["os language"] = osLanguage.lowercaseString
            }
        }
    }
    
    public func ToJson() throws -> NSData {
        var values = [String: AnyObject]()
        values["userId"] = self.userId
        values["event"] = self.message
        values["properties"] = self.properties
        return try NSJSONSerialization.dataWithJSONObject(values, options: [])
    }
    
}

import Foundation

public let SegmentBetaToken: String = "4EFdOmMsMGUinXphTQYnKXB9wXk4vy0T"
public let SegmentDevelopmentToken: String = "qJqnsOy2Bz9n63qFOaSO0mKTv3sIwEls"


public let StableChannel = Channel("Stable", SegmentBetaToken)
public let BetaChannel = Channel("Beta", SegmentBetaToken)
public let TestChannel = Channel("Test", SegmentDevelopmentToken)
public let MasterChannel = Channel("Master", SegmentDevelopmentToken)

public class Channel
{
    let Name: String
    let AnalyticsToken: String
    
    public init(_ Name: String, _ AnalyticsToken: String) {
        self.Name = Name
        self.AnalyticsToken = AnalyticsToken
    }
    
    public static func Load() -> Channel {
        switch Utils.getBuildType().lowercaseString {
            case "stable": return StableChannel
            case "beta": return BetaChannel
            case "test": return TestChannel
            default: return MasterChannel
        }
    }
}

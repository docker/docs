import Foundation

let analytics: Analytics = Analytics(channel: Channel.Load(), version: Version(), tracking: Tracking(Preferences.sharedInstance), preferences: Preferences.sharedInstance) // FIXME

public class Analytics {
    private var channel: Channel
    private var version: Version
    private var tracking: Tracking
    private var preferences: Preferences
    
    public init(channel: Channel, version: Version, tracking: Tracking, preferences: Preferences) {
        self.channel = channel
        self.version = version
        self.tracking = tracking
        self.preferences = preferences
    }
    
    public func track(event: AnalyticEvent) {
        let analyticsEnabled = preferences.getAnalyticsEnabled()
        guard analyticsEnabled || event.isCore() else {
            Logger.log(level: ASLLogger.Level.Error, content: "Not tracking: \(event.text())")
            return
        }
        track(event.text(), analyticsEnabled)
    }
    
    func track(message: String, _ analyticsEnabled: Bool) {
        guard let url: NSURL = NSURL(string: "https://api.segment.io/v1/track")  else {
            Logger.log(level: ASLLogger.Level.Error, content: "Failed to track event: \(message). Can't create URL")
            return
        }
        let request = NSMutableURLRequest(URL: url)
        request.HTTPMethod = "POST"
        request.addValue("application/json; charset=utf-8", forHTTPHeaderField: "Content-Type")
        request.addValue(createAuthenticationHeader(), forHTTPHeaderField: "Authorization")
        
        do {
            request.HTTPBody = try SegmentMessage(userId: tracking.id, message: message, channel: "master", version: version, analyticsEnabled: analyticsEnabled).ToJson()
        } catch let error as NSError {
            Logger.log(level: ASLLogger.Level.Error, content: "Failed to track event: \(message). \(error.description)")
            return
        }
        
        let task = NSURLSession.sharedSession().dataTaskWithRequest(request) {
            data, response, error in
            guard error == nil && data != nil else {
                if let errorMessage = error?.description {
                    Logger.log(level: ASLLogger.Level.Error, content: "Failed to track event: \(message). \(errorMessage)")
                }
                return
            }
            
            if let statusCode = (response as? NSHTTPURLResponse)?.statusCode {
                guard statusCode == 200 else {
                    Logger.log(level: ASLLogger.Level.Error, content: "Failed to track event: \(message). http return code \(statusCode)")
                    return
                }
            } else {
                Logger.log(level: ASLLogger.Level.Error, content: "Failed to track event: \(message). can't get http return code")
                return
            }
            
            Logger.log(level: ASLLogger.Level.Info, content: "Sent \(message) event.")
        }
        task.resume()
    }
    
    func createAuthenticationHeader() -> String {
        let loginString = NSString(format: "%@:", channel.AnalyticsToken)
        if let loginData: NSData = loginString.dataUsingEncoding(NSUTF8StringEncoding) {
            let base64LoginString = loginData.base64EncodedStringWithOptions([])
            return "Basic \(base64LoginString)"
        }
        return ""
    }
    
}

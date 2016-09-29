import Foundation

public class Tracking {
    var id: String = ""
    var preferences: Preferences
    
    public init(_ preferences: Preferences) {
        self.preferences = preferences
        self.id = readOrCreate()
    }
    
    private func readOrCreate() -> String {
        if self.preferences.keyExists(Preferences.Key.analyticsUserID) {
            let uuid: String = Preferences.sharedInstance.getAnalyticsUserId()
            if !uuid.isEmpty {
                return uuid
            }
        }
        let newUuid: String = NSUUID().UUIDString
        self.preferences.setAnalyticsUserId(newUuid)
        return newUuid
    }
}

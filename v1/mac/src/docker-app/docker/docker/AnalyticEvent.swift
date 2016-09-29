import Foundation

public enum AnalyticEvent {
    case AppLaunched, InstallAskForToken, InstallTokenIsValid, InstallAskForRootPassword, InstallShowWelcomePopup, AppRunning, Heartbeat
    
    func isCore() -> Bool {
        return self == Heartbeat
    }
    
    func text() -> String {
        let selfString = String(self)
        return String(selfString.characters.prefix(1)).lowercaseString + String(selfString.characters.dropFirst())
    }
}

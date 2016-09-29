import Foundation
import Cocoa

class MainMenuHelper: NSObject {
    
    func aboutDocker(sender: AnyObject) {
        if let appDelegate = NSApp.delegate as? AppDelegate {
            appDelegate.aboutWindow.showWindow(sender)
        } else {
            Logger.log(level: Logger.Level.Fatal, content: "Unable to cast NSApp.delegate to AppDelegate")
        }
    }
    
    func preferences(sender: AnyObject) {
        if let appDelegate = NSApp.delegate as? AppDelegate {
            appDelegate.settings.showWindow(sender)
        } else {
            Logger.log(level: Logger.Level.Fatal, content: "Unable to cast NSApp.delegate to AppDelegate")
        }
    }
}

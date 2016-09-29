import Foundation
import Cocoa

class AppManager {
    
    static let sharedInstance = AppManager()
    private init() {
    }
    
    func start() {
        dispatch_async(dispatch_get_main_queue(), {
            self.monitorActivationPolicy()
        })
    }
    
    func debugWindows() {
        for w in NSApp.windows {
            NSLog("w%d: %@, %@ - %d", w.windowNumber, w.className, w.title, w.titleVisibility.rawValue)
        }
        NSLog("----------------------------------")
    }
    
    func getActiveWindowCount() -> Int {
        return NSApp.windows.filter({ $0.visible && !($0.className.containsString("NSStatusBarWindow") || $0.className.containsString("NSCarbonMenuWindow"))}).count
    }
    
    private func monitorActivationPolicy() {
        let activeWindowsCount = getActiveWindowCount()
        let activationPolicy = NSApp.activationPolicy()
        if activeWindowsCount == 0 && activationPolicy == NSApplicationActivationPolicy.Regular
        {
            NSApp.setActivationPolicy(NSApplicationActivationPolicy.Accessory)
        }
        else if activeWindowsCount > 0 && activationPolicy == NSApplicationActivationPolicy.Accessory
        {
            NSApp.setActivationPolicy(NSApplicationActivationPolicy.Regular)
            NSApp.activateIgnoringOtherApps(true)
        }

        let delayInSeconds = 0.5
        let delay = dispatch_time(DISPATCH_TIME_NOW, Int64(delayInSeconds * Double(NSEC_PER_SEC)))
        dispatch_after(delay, dispatch_get_main_queue(), {
            self.monitorActivationPolicy()
        })
    }
}

import Cocoa

@NSApplicationMain
class AppDelegate: NSObject, NSApplicationDelegate {

    // MARK: Constants
    
    // used only for logging purposes
    let appBundleId: String = "com.docker.docker"
    let helperBundleId: String = "com.docker.helper"
    
    // MARK: NSApplicationDelegate implementation
    
    func applicationDidFinishLaunching(aNotification: NSNotification) {
        
        Logger.log(helperBundleId, level: ASLLogger.Level.Notice, content: "launching...")
        
        var autoStart: Bool = true
        var dockerAppLaunchPath: String? = nil
        
        // userDefaults from shared group
        if let userDefaults: NSUserDefaults = NSUserDefaults(suiteName: Config.dockerAppGroupIdentifier) {
            Logger.log(helperBundleId, level: ASLLogger.Level.Notice, content: "accessing preferences...")
            // check if key exists
            if userDefaults.objectForKey("prefAutoStart") != nil {
                autoStart = userDefaults.boolForKey("prefAutoStart")
            }
            if userDefaults.objectForKey("DockerAppLaunchPath") != nil {
                dockerAppLaunchPath = userDefaults.stringForKey("DockerAppLaunchPath")
            }
        }
        
        let helperAppPath: String = NSBundle.mainBundle().bundlePath
        Logger.log(helperBundleId, level: ASLLogger.Level.Notice, content: "bundle path: \(helperAppPath)")
        Logger.log(helperBundleId, level: ASLLogger.Level.Notice, content: "auto start: \(autoStart)")
        Logger.log(helperBundleId, level: ASLLogger.Level.Notice, content: "launch path: \(dockerAppLaunchPath)")
        if autoStart {
            if !dockerIsRunning() {
                // Launch Docker main application
                // since the helper application (the current one) is located in the
                // bundle of the main application (in Contents/Library/LoginItems/)
                // we can get the main application bundle's absolute path from the helper
                // application bundle's absolute path.
                var mainAppPath: NSString
                if let launchPath = dockerAppLaunchPath {
                    mainAppPath = launchPath
                } else {
                    mainAppPath = (helperAppPath as NSString).stringByDeletingLastPathComponent
                    mainAppPath = mainAppPath.stringByDeletingLastPathComponent
                    mainAppPath = mainAppPath.stringByDeletingLastPathComponent
                    mainAppPath = mainAppPath.stringByDeletingLastPathComponent
                }
                
                if NSWorkspace.sharedWorkspace().launchApplication(String(mainAppPath)) {
                    // successfully launched main application
                    Logger.log(helperBundleId, level: Logger.Level.Notice, content: "Docker.app sucessfully launched from: \(mainAppPath)")
                } else {
                    // failed to launch main application
                    Logger.log(helperBundleId, level: Logger.Level.Error, content: "Failed to launch Docker.app from: \(mainAppPath)")
                    // TODO: error handling (open an error popup?)
                }
            } else {
                let app = NSRunningApplication.runningApplicationsWithBundleIdentifier(appBundleId)
                Logger.log(helperBundleId, level: Logger.Level.Warning, content: "Docker app already launched: \(app[0].bundleURL?.absoluteString)")
            }
        }
        NSApp.terminate(nil)
    }
    
    func applicationWillTerminate(aNotification: NSNotification) {
        Logger.log(helperBundleId, level: ASLLogger.Level.Notice, content: "terminating...")
    }
    
    // MARK: Private functions
    
    // Test if docker app is already running
    private func dockerIsRunning() -> Bool {
        // test if at least one instance of the Docker app is running application
        let apps: [NSRunningApplication] = NSRunningApplication.runningApplicationsWithBundleIdentifier(appBundleId)
        return apps.count >= 1
    }
}

import Cocoa

class AboutWindow: NSWindowController, NSWindowDelegate {
    private let LOCAL_BUILD_NUMBER = "999999999"
    
    @IBOutlet weak var iconImageView: NSImageView?
    @IBOutlet weak var versionLabel: NSTextField?
    
    override func windowDidLoad() {
        super.windowDidLoad()
        
        guard let appVersion = NSBundle.mainBundle().objectForInfoDictionaryKey("CFBundleShortVersionString") as? String else {
            Logger.log(level: Logger.Level.Fatal, content: "Invalid data in main bundle - unable to read app version")
            return
        }
        
        guard var buildNumber = NSBundle.mainBundle().objectForInfoDictionaryKey("CFBundleVersion") as? String else {
            Logger.log(level: Logger.Level.Fatal, content: "Invalid data in main bundle - unable to read build version")
            return
        }
        
        if buildNumber == LOCAL_BUILD_NUMBER {
            buildNumber = "local"
        }
        
        var shortCommit = git_commit
        if shortCommit.characters.count > 10 {
            shortCommit = git_commit.substringToIndex(git_commit.startIndex.advancedBy(10))
        }
        
        versionLabel?.stringValue = "Version \(appVersion) (build: \(buildNumber))\n\(shortCommit)"
        
        if !Utils.isStableBuild() {
             self.iconImageView?.image = NSImage(named: "AboutIconBeta")
        }
        
    }
    
    // We force centering the window, overriding previous location each time.
    override func showWindow(sender: AnyObject?) {
        self.window?.center()
        super.showWindow(sender)
    }
    
    @IBAction func OpenAcknowledgements(sender: AnyObject) {
        openTextFile("Acknowledgements", "OSS-LICENSES", "")
    }
    
    @IBAction func OpenLicenseAgreements(sender: AnyObject) {
        openTextFile("License", "LICENSE", "rtf")
    }

    private func openTextFile(name: String, _ fileName: String, _ fileType: String) {
        guard let filePath = NSBundle.mainBundle().pathForResource(fileName, ofType:fileType) else {
            openInternalErrorModal("Could not read \(name) file path.")
            return
        }
        
        guard NSFileManager.defaultManager().fileExistsAtPath(filePath) else {
            openInternalErrorModal("Could not find \(name) file.")
            return
        }
  
        NSWorkspace.sharedWorkspace().openFile(filePath, withApplication: "TextEdit")
    }
    
    private func openInternalErrorModal(message: String) {
        let errorAlert: NSAlert = NSAlert()
        errorAlert.messageText = "Internal Error"
        errorAlert.informativeText = message
        errorAlert.alertStyle = NSAlertStyle.Critical
        errorAlert.runModal()
    }
}

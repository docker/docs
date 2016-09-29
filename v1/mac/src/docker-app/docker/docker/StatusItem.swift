import Cocoa

// StatusItem class is a manager for the docker menu bar extra
class StatusItem: NSObject, NSMenuDelegate {
    
    // properties
    let statusItemIcon = NSImage(named: "statusItemIcon-0")
    let statusItem: NSStatusItem = NSStatusBar.systemStatusBar().statusItemWithLength(NSSquareStatusItemLength)
    private var currentFrame = 0
    private let maxFrame = 5
    private var increment = 1
    private var statusAnimationTimer: NSTimer?
    
    // The popover displayed the first time the app is launched
    private var introPopover: NSPopover?
    // NOTE(aduermael): These 3 variables are used to avoid showing intro
    // popover too early. I couldn't find a better way to know if the button
    // in the status bar was ready to be used as a relative view for a
    // popover.
    private let initDate: NSDate
    private let minTimeIntervalForPopover: NSTimeInterval = 3.0 // 3 seconds
    private var introPopoverTimer: NSTimer?
    
    private let dockerStateViewController = DockerStateViewController(nibName: "DockerStateView", bundle: nil)
    
    // Menu displayed by the status item
    @IBOutlet var statusItemMenu: NSMenu?
    //
    @IBOutlet var checkForUpdatesMenuItem: NSMenuItem?
    
    @IBOutlet var diagnoseAndFeedbackMenuItem: NSMenuItem?
    
    var stateItem: NSMenuItem?
    
    // initializer
    // TODO: test for "loadNibNamed" failure and convert to a failable initializer
    override init() {
        
        initDate = NSDate()
        
        // call the superclass initializer
        super.init()
        
        // load the InterfaceBuilder file for this class
        // and don't store top level objects in an array. (we use IBOutlets)
        NSBundle.mainBundle().loadNibNamed("StatusItem", owner: self, topLevelObjects: nil)
        
        // indicate the icon is a template so that the OS can automatically
        // flip the colors for OSX "Dark Mode"
        statusItemIcon?.template = true
        
        // set status_item's image and menu
        statusItem.image = statusItemIcon
        statusItem.menu = statusItemMenu
        
        // display Docker status and add observer
        let instantNotification = NSNotification(name: Backend.dockerStateNotificationName, object: Backend.getCurrentState().rawValue)
        dockerStatusChanged(instantNotification)
        NSNotificationCenter.defaultCenter().addObserver(self, selector: #selector(StatusItem.dockerStatusChanged(_:)), name: Backend.dockerStateNotificationName, object: nil)
    }
    
    func dockerStatusChanged(notification: NSNotification) {
        if let dockerStateStr = notification.object as? String {
            if let dockerState = Backend.DockerState(rawValue: dockerStateStr) {
                statusAnimationTimer?.invalidate()
                currentFrame = 0
                let icon = NSImage(named: "statusItemIcon-\(currentFrame)")
                icon?.template = true
                statusItem.image = icon
                if dockerState == .Starting {
                    statusAnimationTimer = NSTimer.scheduledTimerWithTimeInterval(1.0/7.0, target: self, selector: #selector(StatusItem.updateIcon), userInfo: nil, repeats: true)
                }
            }
        }
    }
    
    func updateIcon() {
        if currentFrame == 0 {
            increment = 1
        }
        if currentFrame == maxFrame {
            increment = -1
        }
        currentFrame += increment
        
        let icon = NSImage(named: "statusItemIcon-\(currentFrame)")
        icon?.template = true
        
        statusItem.image = icon
    }
    
    deinit {
        // invalidate introPopoverTimer in case it has been initiated
        introPopoverTimer?.invalidate()
        introPopoverTimer = nil
    }
    
    func showIntroPopover() {
        // make sure the popover is not already being displayed
        if introPopover == nil {
            if let button = statusItem.button {
                // Checking to avoid displaying popover too early.
                // See comment above about initDate, minTimeIntervalForPopover
                // & introPopoverTimer
                let timeSinceInit = NSDate().timeIntervalSinceDate(initDate)
                if timeSinceInit < minTimeIntervalForPopover {
                    introPopoverTimer?.invalidate()
                    introPopoverTimer = nil
                    introPopoverTimer = NSTimer.scheduledTimerWithTimeInterval(minTimeIntervalForPopover - timeSinceInit, target: self, selector: #selector(StatusItem.showIntroPopover), userInfo: nil, repeats: false)
                    
                } else {
                    introPopover = NSPopover()
                    if let popover = introPopover {
                        popover.appearance = NSAppearance(named: NSAppearanceNameAqua)
                        popover.contentViewController = IntroViewController(statusItem: self)
                        popover.showRelativeToRect(button.bounds, ofView: button, preferredEdge: NSRectEdge.MinY)
                    }
                }
            }
        }
    }
    
    func hideIntroPopover() {
        if let popover = introPopover {
            introPopover = nil
            popover.performClose(self)
        }
    }
    
    // NSMenuDelegate protocol implementation for the StatusItem menu
    
    func menuWillOpen(menu: NSMenu) {
        
        if stateItem == nil{
            stateItem = NSMenuItem()
            if let stateItem = stateItem {
                stateItem.view = dockerStateViewController?.view
                statusItem.menu?.insertItem(stateItem, atIndex: 0)
            }
        }
        
        // "check for update" is not available on dev builds
        // "check for update" is not available for non-admin users
        if !UpdateManager.isAppUpdateAvailable() {
            self.checkForUpdatesMenuItem?.hidden = true
            self.checkForUpdatesMenuItem?.target = nil
            self.checkForUpdatesMenuItem?.action = nil
        }
        
        if let popover = introPopover {
            introPopover = nil
            popover.performClose(self)
        }
    }
    
    func menuDidClose(menu: NSMenu) {
        
    }
    
    func menuNeedsUpdate(menu: NSMenu) {
        // Logger.log(level: ASLLogger.Level.Notice, content: "🦁 menu needs update")
    }
    
    func Hide() {
        statusItem.length = 0
        statusItem.enabled = false
    }
    
    func Show() {
        statusItem.length = NSVariableStatusItemLength
        statusItem.enabled = true
    }
    
    func Disable() {
        statusItem.enabled = false
    }
    
    func Enable() {
        statusItem.enabled = true
    }
    
    // Menu item callbacks
    
    // this function is called when the "Settings" item of the menu bar menu
    // is clicked
    @IBAction func SettingsActionReceived(sender: AnyObject) {
        // call showWindow on the Settings controller
        if let appDelegate = NSApp.delegate as? AppDelegate {
            appDelegate.settings.showWindow(sender)
        } else {
            Logger.log(level: Logger.Level.Fatal, content: "Unable to cast NSApp.delegate to AppDelegate")
        }
    }
    
    @IBAction func AboutItemActionReceived(sender: AnyObject) {
        // call showWindow on the AboutWindow controller
        if let appDelegate = NSApp.delegate as? AppDelegate {
            appDelegate.aboutWindow.showWindow(sender)
        } else {
            Logger.log(level: Logger.Level.Fatal, content: "Unable to cast NSApp.delegate to AppDelegate")
        }
    }
    
    @IBAction func checkForUpdateActionReceived(sender: AnyObject) {
        UpdateManager.sharedManager().checkForUpdates(sender)
    }
    
    @IBAction func OpenKitematicActionReceived(sender: AnyObject) {
        
        let kitematicBundleID = "com.electron.kitematic_(beta)"
        let kitematicMinVersionKey = "DockerKitematicMinVersion"
        
        // get the min kitematic version from the Docker for Mac bundle
        guard let kmMinVersionString: String = NSBundle.mainBundle().objectForInfoDictionaryKey(kitematicMinVersionKey) as? String else {
            // failed to get the kitematic min version from the bundle plist
            self.displayKitematicOpeningError("kitematicErrorGettingMinVersionFromBundle".localize())
            return
        }
        // get the actual version of the installed Kitematic bundle
        var kmVersionStringOptional: String? = nil
        if let kmBundlePath = NSWorkspace.sharedWorkspace().absolutePathForAppBundleWithIdentifier(kitematicBundleID) {
            if let kmBundle = NSBundle(path: kmBundlePath) {
                if let kmBundleVersion = kmBundle.objectForInfoDictionaryKey("CFBundleVersion") as? String {
                    kmVersionStringOptional = kmBundleVersion
                }
            }
        } else {
            // Kitematic is not installed on the Mac (no bundle found)
            guard let kitematicMessage = WizardMessageGeneric(message: "Kitematic", details: "kitematicInfo".localize(), icon: "Kitematic") else {
                Logger.log(level: ASLLogger.Level.Fatal, content: "Internal error: \(#function)")
                return
            }
            kitematicMessage.closeButton = true
            kitematicMessage.addCloseButton("Ok".localize()).makeDefault()
            Wizard.show(kitematicMessage)
            return
        }
        guard let kmVersionString = kmVersionStringOptional else {
            self.displayKitematicOpeningError("kitematicErrorGettingKitematicBundleVersion".localize())
            return
        }
        // split the versions on the periods ('.')
        let kmMinVersion: [String] = kmMinVersionString.componentsSeparatedByString(".")
        if kmMinVersion.count != 3 {
            self.displayKitematicOpeningError("kitematicErrorWrongMinVersionFormat".localize())
            return
        }
        let kmVersion: [String] = kmVersionString.componentsSeparatedByString(".")
        if kmVersion.count != 3 {
            self.displayKitematicOpeningError("kitematicErrorWrongVersionFormat".localize())
            return
        }
        
        // convert strings into integers
        guard let minMajor = Int(kmMinVersion[0]), minMinor = Int(kmMinVersion[1]), minPatch = Int(kmMinVersion[2]) else {
            self.displayKitematicOpeningError("kitematicErrorWrongMinVersionFormat".localize())
            return
        }
        guard let major = Int(kmVersion[0]), minor = Int(kmVersion[1]), patch = Int(kmVersion[2]) else {
            self.displayKitematicOpeningError("kitematicErrorWrongVersionFormat".localize())
            return
        }
        
        // compare versions
        if (major > minMajor) ||
            (major == minMajor && minor > minMinor) ||
            (major == minMajor && minor == minMinor && patch >= minPatch) {
            // launch kitematic
            NSWorkspace.sharedWorkspace().launchAppWithBundleIdentifier(kitematicBundleID, options: .Default, additionalEventParamDescriptor: nil, launchIdentifier: nil)
            return
        }
        
        // kitematic is too old, display error message
        guard let kitematicMessage = WizardMessageGeneric(message: "Kitematic",
                                                          details: "kitematicInfoVersion".localize(major, minor, patch, minMajor, minMinor, minPatch),
                                                          icon: "Kitematic")
            else {
                Logger.log(level: ASLLogger.Level.Fatal, content: "Internal error: \(#function)")
                return
        }
        kitematicMessage.closeButton = true
        kitematicMessage.addCloseButton("Ok".localize()).makeDefault()
        Wizard.show(kitematicMessage)
    }
    
    @IBAction func diagnoseAndFeedbackActionReceived(sender: AnyObject) {
        // If the diagnose window is open already we bring it to the foreground
        // and exit the function.
        if Wizard.selectNonModal("diagnoseAndFeedback") {
            return
        }
        // Create an instance of the diagnose popup and show it.
        guard let diagnoseMessage = WizardMessageDiagnose(message: "fatalErrorDiagnoseAndFeedbackBtn".localize()) else {
            return
        }
        diagnoseMessage.title = "Diagnose & Feedback"
        diagnoseMessage.closeButton = true
        diagnoseMessage.minimizeButton = true
        Wizard.showNonModal(diagnoseMessage, key: "diagnoseAndFeedback")
    }
    
    @IBAction func DocsActionReceived(sender: AnyObject) {
        let docsUrl: String = "UrlsDocumentation".localize()
        if let url = NSURL(string: docsUrl) {
            NSWorkspace.sharedWorkspace().openURL(url)
        } else {
            Logger.log(level: ASLLogger.Level.Error, content: "\(#function): malformed docs URL")
        }
    }
    
    // Quit menu item has been clicked
    @IBAction func QuitActionReceived(_: AnyObject) {
        // notify the application delegate that quit has been requested
        NSApp.terminate(nil)
    }
    
    // MARK: Private functions
    
    //
    private static func getAppVersionAndBuildNumber() -> String {
        // get app version name
        guard let versionName: String = NSBundle.mainBundle().objectForInfoDictionaryKey("CFBundleShortVersionString") as? String else {
            return ""
        }
        // get app build number
        guard let buildNumber: String = NSBundle.mainBundle().objectForInfoDictionaryKey("CFBundleVersion") as? String else {
            return ""
        }
        let intermediate: String = " - \(versionName) (\(buildNumber))"
        let result: String = intermediate.stringByReplacingOccurrencesOfString(" ", withString: "%20")
        return result
    }
    
    // convenience: display popup for Kitematic opening errors
    private func displayKitematicOpeningError(messageContent: String) {
        guard let message = WizardMessageGeneric(message: "Kitematic", details: messageContent, icon: "Kitematic") else {
            Logger.log(level: ASLLogger.Level.Fatal, content: "Internal error: \(#function)")
            return
        }
        message.closeButton = true
        message.addCloseButton("Ok".localize()).makeDefault()
        Wizard.show(message)
    }
}

import Cocoa

class IntroViewController: NSViewController {
    
    @IBOutlet weak var text1: NSTextField?
    @IBOutlet weak var text2: NSTextField?
    
    @IBOutlet weak var cmdLabel: NSTextField?
    let commands = ["docker ps", "docker info", "docker version", "docker images"]
    var typingTimer: NSTimer?
    private var cmdPos = 0
    private var cmdIndex = 0
    private var waitBetweenCmds = 0
    private let waitValue = 20
    private let statusItem: StatusItem
    
    init?(statusItem: StatusItem) {
        self.statusItem = statusItem
        super.init(nibName: "IntroViewController", bundle: NSBundle.mainBundle())
    }
    
    required init?(coder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }
    
    override func viewDidLoad() {
        super.viewDidLoad()
        // Do view setup here.
    }
    
    override func viewDidAppear() {
        let state = Backend.getCurrentState()
        if state == Backend.DockerState.Starting {
            text1?.stringValue = "welcomePopoverStartingText1".localize()
            text2?.stringValue = "welcomePopoverStartingText2".localize()
            NSNotificationCenter.defaultCenter().addObserver(self, selector: #selector(IntroViewController.dockerStateChanged(_:)), name: Backend.dockerStateNotificationName, object: nil)
        } else {
            text1?.stringValue = "welcomePopoverRunningText1".localize()
            text2?.stringValue = "welcomePopoverRunningText2".localize()
        }
        updateTrackingCheckButton()
        typingTimer?.invalidate()
        typingTimer = NSTimer.scheduledTimerWithTimeInterval(0.05, target: self, selector: #selector(IntroViewController.type), userInfo: nil, repeats: true)
        NSNotificationCenter.defaultCenter().addObserver(self, selector: #selector(IntroViewController.preferencesDidChange(_:)), name: Preferences.preferencesDidChangeNotification, object: nil)
    }
    
    func dockerStateChanged(notification: NSNotification) {
        if let dockerStateStr = notification.object as? String {
            if let dockerState = Backend.DockerState(rawValue: dockerStateStr) {
                if dockerState == Backend.DockerState.Running {
                    text1?.stringValue = "welcomePopoverRunningText1".localize()
                    text2?.stringValue = "welcomePopoverRunningText2".localize()
                }
            }
        }
    }
    
    override func viewWillDisappear() {
        typingTimer?.invalidate()
        NSNotificationCenter.defaultCenter().removeObserver(self, name: Preferences.preferencesDidChangeNotification, object: nil)
        NSNotificationCenter.defaultCenter().removeObserver(self, name: Backend.dockerStateNotificationName, object: nil)
    }
    
    func preferencesDidChange(_: NSNotification) {
        updateTrackingCheckButton()
    }
    
    private func updateTrackingCheckButton() {
        let trackingEnabled = Preferences.sharedInstance.getAnalyticsEnabled()
        let autoSendCrashReports = Preferences.sharedInstance.getAutoSendCrashReports()
        if trackingEnabled && autoSendCrashReports {
            trackingCheckButton?.state = NSOnState
            trackingCheckButton?.allowsMixedState = false
        } else if !trackingEnabled && !autoSendCrashReports {
            trackingCheckButton?.state = NSOffState
            trackingCheckButton?.allowsMixedState = false
        } else {
            trackingCheckButton?.allowsMixedState = true
            trackingCheckButton?.state = NSMixedState
        }
    }
    
    func type() {
        let cmd = commands[cmdIndex]
        if cmdPos > cmd.characters.count {
            waitBetweenCmds += 1
            if waitBetweenCmds >= waitValue {
                waitBetweenCmds = 0
                cmdPos = 0
                cmdIndex += 1
                if cmdIndex >= commands.count {
                    cmdIndex = 0
                }
            }
        } else {
            cmdLabel?.stringValue = "$ \(cmd.substringToIndex(cmd.characters.startIndex.advancedBy(cmdPos)))"
            cmdPos += 1
        }
    }
    
    @IBAction func closePopover(sender: AnyObject) {
        statusItem.hideIntroPopover()
    }
    
    @IBOutlet weak var trackingCheckButton: NSButton?
    @IBAction func trackingCheckButtonActionReceived(checkbox: NSButton) {
        // note: click cycle is OFF -> MIXED -> ON -> OFF
        let trackingState: Bool = checkbox.state == NSOnState
        // apply both preferences (only 1 notification)
        if Preferences.sharedInstance.setAnalyticsEnabled(trackingState, notify: false) == false {
            Logger.log(level: ASLLogger.Level.Error, content: "internal error (\(#function))")
        }
        if Preferences.sharedInstance.setAutoSendCrashReports(trackingState) == false {
            Logger.log(level: ASLLogger.Level.Error, content: "internal error (\(#function))")
        }
    }
    
    @IBAction func openPrivacySettings(_: NSButton) {
        if let appDelegate = NSApp.delegate as? AppDelegate {
            appDelegate.settings.showWindow(SettingsPanelRequest(.Privacy))
        }
    }
    
}

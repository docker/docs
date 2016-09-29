import Cocoa
import ServiceManagement
import HockeySDK

@NSApplicationMain
class AppDelegate: NSObject, NSApplicationDelegate {
    
    // lists the different states the application can be in
    private enum AppState {
        case Starting
        // used when we prevent the app from starting. Since the app cannot
        // terminate if it isn't in the "Running" state, we have this state
        // so the app knows we are aborting the launch and therefore it can
        // terminate.
        case Aborting
        case Running
        case Exiting
    }
    
    // appState is used to provide a little more control on the application's
    // lifecycle. For example, it can't quit if not in Running state.
    private var appState: AppState = AppState.Starting
    
    // this is the object responsible for the button in the OSX menu bar.
    // The constructor itself will make the whale icon appear in the menu bar.
    var statusItem: StatusItem?
    
    // about_window could be created from anywhere, but it's safer to
    // ask the AppDelegate each time, to make sure we never have more
    // that one about window displayed.
    let aboutWindow = AboutWindow(windowNibName: "AboutWindow")
    
    // settings window controller
    let settings = Settings(windowNibName: "Settings")
    
    // MARK: NSApplicationDelegate implementation
    
    func applicationDidFinishLaunching(aNotification: NSNotification) {
        
        AppManager.sharedInstance.start()
        
        // debug layout constraints:
        // NSUserDefaults.standardUserDefaults().setBool(true, forKey: "NSConstraintBasedLayoutVisualizeMutuallyExclusiveConstraints")
        // NSUserDefaults.standardUserDefaults().synchronize()

        NSApp.disableRelaunchOnLogin()
        
        // this is used when launching the app with --uninstall option flag
        if Options.uninstall {
            // turn off Logger
            Logger.enabled = false
            
            guard let bundleId: String = NSBundle.mainBundle().bundleIdentifier else {
                print("uninstallFailed".localize("can't check if another instance of same application is running."))
                self.appState = AppState.Aborting
                NSApp.terminate(nil)
                return
            }
            
            let apps: [NSRunningApplication] = NSRunningApplication.runningApplicationsWithBundleIdentifier(bundleId)
            let currentApp: NSRunningApplication = NSRunningApplication.currentApplication()
            // If other instances are running, kill them
            for runningApp in apps {
                if runningApp.isEqual(currentApp) == false {
                    print("Docker is running, exiting...")
                    let success = runningApp.terminate()
                    if !success {
                        runningApp.forceTerminate()
                    }
                }
            }
            
            let err = Uninstall.uninstallWhileKeepingPreferenceKeys([Preferences.Key.analyticsUserID])
            if let errMessage = err {
                print("uninstallFailed".localize(errMessage))
            } else {
                print("uninstallSuccessMessage".localize())
            }
            self.appState = AppState.Aborting
            NSApp.terminate(nil)
        }
        
        // checks if there is another running instance of Docker.app and
        // terminate current application if it is the case to avoid to have
        // more than one instance of Docker.app running at any given time.
        exitIfApplicationAlreadyRunning()
        
        // Log application version so we have it in the logfile
        self.logAppVersion()
        
        if Options.unattended == false {
            // initialize HockeyApp crash report system
            initializeCrashReporter()
        }
        
        // Check environment requirements
        switch Environment.Check() {
        case let .Fatal(errorMessage):
            self.appState = AppState.Aborting
            Logger.log(level: Logger.Level.Fatal, content: errorMessage)
        case let .Warning(warning):
            Logger.log(level: Logger.Level.Error, content: warning.message)

            let warningMessage = WizardMessageGeneric(message: warning.title, details: warning.message, icon: "FatalIcon")

            guard let message = warningMessage else {
                Logger.log(level: ASLLogger.Level.Fatal, content: "Internal error: \(#function)")
                return // unreached after fatal log
            }

            let ignoreBtn = message.addButton(warning.ignoreBtn)
            ignoreBtn.blockAction = {(button: WizardButton?, message: WizardMessage?) -> () in
                message?.close()
            }
            let exitBtn = message.addButton("fatalErrorExitBtn".localize()).makeDefault()
            exitBtn.blockAction = {(button: WizardButton?, message: WizardMessage?) -> () in
                NSApp.terminate(nil)
            }

            Wizard.show(message)
        case .OK: ()
        }
        
        // perform migrations
        if let error: String = UpdateMigrator.migrate() {
            self.appState = AppState.Aborting
            Logger.log(level: ASLLogger.Level.Fatal, content: "failed migration: \(error)")
        }

        if UpdateMigrator.checkChannelChanges() {
            NSApp.terminate(nil)
        }

        // If it is a fresh install and the user is using OS X Yosemite, we
        // display a popup warning him/her about the known issues.
        if Utils.getOSXMajorVersion() == 10 && Utils.getOSXMinorVersion() == 10 {
            if Options.unattended == false {
                if Preferences.sharedInstance.getInstallWelcomedUser() == false {
                    if let message = WizardMessageGeneric(message: "yosemiteIssuesTitle".localize(), details: "yosemiteIssuesMessage".localize(), icon: "MsgWarningIcon") {
                        message.addCloseButton("Ok".localize()).makeDefault()
                        Wizard.show(message)
                    }
                }
            }
        }
        
        // initiate status item (whale in the menu bar)
        // this will also initiate sparkle (responsible for app updates)
        // Because Sparkle has been integrated in StatusItem.xib
        statusItem = StatusItem()
        
        // task manager launch backend processes, catch their stdout/stderr,
        // and gently shutdown them (wait for moby to shutdown, then shutdown
        // everybody else). global task manager in charge of scheduling all
        // docker daemons / processes
        TaskManager.Start()
        
        // from this point, the app can quit when being asked to.
        // Note(aduermael): I don't know exactly why we can't quit before this
        // specific moment, but I kept it as it was implemented by ebriney
        appState = AppState.Running
        
        // initialize mixpanel client and send first event
        heartbeat.start()
        analytics.track(AnalyticEvent.AppLaunched)
        
        // if we never said welcome before, we show the initial welcome popup
        if Options.unattended == false {
            if Preferences.sharedInstance.getInstallWelcomedUser() == false {
                if let nonAdminMessage: WizardMessage = WizardMessageGeneric(message: "welcomeBetaMessage".localize(), details: "welcomeBetaAdditionalMessage".localize(), icon: "MsgWelcomeIcon") {
                    nonAdminMessage.addCloseButton("Next".localize()).makeDefault()
                    Wizard.show(nonAdminMessage)
                }
            }
        }
        
        // welcome message and warning about os x version (if 10.10) won't
        // be displayed anymore
        Preferences.sharedInstance.setInstallWelcomedUser(true)
        
        // create the UpdateManager instance
        UpdateManager.sharedManager()
        
        // path of the Docker.app bundle being executed
        let dockerAppBundlePath: String = NSBundle.mainBundle().bundlePath
        // we record the path of the last docker app launched for the helper
        // WARNING: to be safe do not move that, because group containers
        // preference plist will not be created if we manually create the directory
        if Preferences.sharedInstance.setDockerAppLaunchPath(dockerAppBundlePath) == false {
            Logger.log(level: ASLLogger.Level.Error, content: "failed to save docker app launch path")
        }
        
        // dump app fullpath in log file
        Logger.log(level: Logger.Level.Notice, content: "Bundle path: \(dockerAppBundlePath)")
        
        // perform the part of the installation that requires elevated rights
        Install.performPrivilegedInstallation()
        
        // If we're running with "--quit-after-install" as root, then exit now: our
        // only task is to install the privileged component, everything else should
        // be done when the app starts as a regular user.
        if Options.quitafterinstall && Utils.userIsRoot {
            Logger.log(level: ASLLogger.Level.Notice, content: "running as root with --quit-after-install: exitting after privileged helper install")
            NSApp.terminate(nil)
        }

        // Since we have the privileged helper, we can create the symlinks
        Install.installSymlinks()
        
        // this will make sure symlinks are reinstalled with correct paths
        // and ownerships whem user session becomes active
        MultipleUsers.listenForUserSessionSwitch()
        
        // do the operations that are required for the application to work
        // properly (set default settings values and start helper application)
        if let error: String = Install.performUnprivilegedInstallation() {
            Logger.log(level: Logger.Level.Fatal, content: "Internal error: \(error)")
        }
        
        // if the app has been launched with the argument "--quit-after-install"
        // we kill the app right now.
        if Options.quitafterinstall {
            NSApp.terminate(nil)
        }
        
        // register the DockerHelper as a LoginItem
        let dockerHelperBundleId: String = "com.docker.helper"
        if SMLoginItemSetEnabled(dockerHelperBundleId, true) == false {
            Logger.log(level: Logger.Level.Error, content: "Failed to launch \(dockerHelperBundleId) and register it as a login item")
        }
        
        Backend.listenForDockerStateChanges()
        
        // Launch tasks needed to use docker
        ChildProcesses.launchDockerTasks()
        
        // if the app hasn't been launched by automated tests we eventually show
        // the welcome popover
        if Options.unattended == false {
            // if we never said welcome before, then we say welcome now
            if Preferences.sharedInstance.getInstallDisplayedPopover() == false {
                
                // NOTE(aduermael): I could find a better way to check if some other
                // app is displayed full screen. And that method is not precise enough
                // to tell us if the full screen app is using same display as D4M.
                // But just in case, if one other app is displayed full screen, don't
                // display the welcome popup (saves it for later).
                // Issue: https://github.com/docker/pinata/issues/2101
                if let screens = NSScreen.screens() {
                    let nbScreens = screens.count
                    var nbMenubars = 0
                    
                    if let windows = CGWindowListCopyWindowInfo(.OptionOnScreenOnly, 0) {
                        for window in windows {
                            guard let owner = window.objectForKey(kCGWindowOwnerName) as? String else {
                                continue
                            }
                            guard let name = window.objectForKey(kCGWindowName) as? String else {
                                continue
                            }
                            if owner == "Window Server" && name == "Menubar" {
                                nbMenubars += 1
                            }
                        }
                    }
                    if nbScreens == nbMenubars {
                        statusItem?.showIntroPopover()
                        analytics.track(AnalyticEvent.InstallShowWelcomePopup)
                        if Preferences.sharedInstance.setInstallDisplayedPopover(true) == false {
                            Logger.log(level: ASLLogger.Level.Error, content: "failed to save preference show welcome popup")
                        }
                    }
                }
            }
        }
        
        if let currentAppChannelType: String = NSBundle.mainBundle().objectForInfoDictionaryKey("DockerBuildType") as? String {
            Preferences.sharedInstance.setDockerBuildType(currentAppChannelType)
        }

        analytics.track(AnalyticEvent.AppRunning)
        
        
        // tell the update manager, it is now allowed to check for updates and
        // eventually show a popup to the user in case a new version is available.
        if Options.unattended == false {
            UpdateManager.sharedManager().allowAutoCheckForUpdates()
        }
    }
    
    // This is part of the NSApplicationDelegate Protocol.
    // It determines if the application should terminate or not.
    // - state is Aborting -> app terminates now
    // - state is ExitRequested or application is unattended -> terminates gracefully
    // - other case -> beeps and do not terminate
    func applicationShouldTerminate(sender: NSApplication) -> NSApplicationTerminateReply {
        // if application is in Aborting state, it means the launch has been
        // aborted. We exit right away without doing anything else.
        if self.appState == AppState.Aborting {
            self.appState = AppState.Exiting
            return NSApplicationTerminateReply.TerminateNow
        }
        
        // the app is in running state, we start the shutdown procedure
        self.appState = AppState.Exiting
        // disable status item
        statusItem?.Disable()
        // stop sending heartbeats every hour
        heartbeat.stop()
                
        return NSApplicationTerminateReply.TerminateNow
    }
    
    //
    func applicationWillTerminate(aNotification: NSNotification) {
        Logger.log(level: Logger.Level.Notice, content: "applicationWillTerminate")
    }
    
    // MARK: private functions
    
    // exits the application if another instance of the Docker.app is already
    // running
    private func exitIfApplicationAlreadyRunning() {
        // test if an application having the same bundle id is currently running
        if let bundleId: String = NSBundle.mainBundle().bundleIdentifier {
            let apps: [NSRunningApplication] = NSRunningApplication.runningApplicationsWithBundleIdentifier(bundleId)
            if apps.count > 1 {
                Logger.log(level: Logger.Level.Warning, content: "Docker already running - exiting")
                // the application cannot terminate while being in the "Starting"
                // state. It has to switch to the "Aborting" state to be able to
                // terminate
                self.appState = AppState.Aborting
                NSApp.terminate(nil)
            }
        } else {
            Logger.log(level: ASLLogger.Level.Fatal, content: "Failed to retrieve application bundle ID")
        }
    }
    
    // init and start crash report manager
    private func initializeCrashReporter() {
        // getting hockeyapp application ID from Info.plist.
        // this field is emptu on dev builds, and the value is set dynamically
        // by the CI for "test", "dev" and "alpha" builds.
        let infoPlistValue: AnyObject? = NSBundle.mainBundle().objectForInfoDictionaryKey("HockeyAppId")
        if let appId = infoPlistValue as? String {
            if appId.characters.count > 0 {
                BITHockeyManager.sharedHockeyManager().configureWithIdentifier(appId)
                BITHockeyManager.sharedHockeyManager().startManager()
            } else {
                // do nothing, it must be a dev build with no crash reporting
                return
            }
        } else {
            Logger.log(level: Logger.Level.Error, content: "\(#function): could not get application id from Info.plist")
        }
    }
    
    // log the app version.
    // This function is called when the app starts so that the log file contains
    // the version of the app.
    private func logAppVersion() {
        // get app version name
        var versionName: String = ""
        if let version: String = NSBundle.mainBundle().objectForInfoDictionaryKey("CFBundleShortVersionString") as? String {
            versionName = version
        }
        // get app build number
        var buildNumber: String = ""
        if let number: String = NSBundle.mainBundle().objectForInfoDictionaryKey("CFBundleVersion") as? String {
            buildNumber = number
        }
        let logString: String = "Application version: \(versionName) (\(buildNumber))"
        Logger.log(level: ASLLogger.Level.Notice, content: logString)
        // also log whether user is considered an administrator
        Logger.log(level: ASLLogger.Level.Notice, content: "Administrator user: \(Utils.userIsAdministrator)")
    }
    
    
    // COMMENTED
    
    // Commented because we should be able to launch the app twice. For example
    // when we try to launch it with uninstall option flag from command line.
    //    // applicationShouldHandleReopen is triggered when double clicking
    //    // on the icon when the app is running.
    //    func applicationShouldHandleReopen(sender: NSApplication, hasVisibleWindows flag: Bool) -> Bool {
    //        Logger.Print("applicationShouldHandleReopen")
    //        return false
    //    }
}

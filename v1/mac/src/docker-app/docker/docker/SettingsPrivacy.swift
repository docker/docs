//
//  SettingsPrivacy.swift
//  Docker
//
//  Created by Gaetan de Villele on 6/30/16.
//  Copyright © 2016 docker. All rights reserved.
//

import Cocoa

class SettingsPrivacy: NSViewController, SettingsPanelProtocol {
    
    // text
    @IBOutlet weak var helpLabel: NSTextField?
    @IBOutlet weak var analyticsDescription: TextViewInView?
    @IBOutlet weak var crashReportingDescription: TextViewInView?
    @IBOutlet weak var pingInfoDescription: TextViewInView?
    
    @IBOutlet weak var settingsWindow: Settings?
    
    @IBOutlet weak var analyticsCheckButton: NSButton?
    @IBAction func analyticsCheckButtonActionReceived(checkbox: NSButton) {
        // this checkbox is either ON or OFF ("mixed" state is not allowed on it)
        let sendAnalytics: Bool = checkbox.state == NSOnState
        if Preferences.sharedInstance.setAnalyticsEnabled(sendAnalytics) == false {
            Logger.log(level: ASLLogger.Level.Fatal, content: "internal error (\(#function))")
            checkbox.state = Preferences.sharedInstance.getAnalyticsEnabled() ? NSOnState : NSOffState
        }
    }
    
    @IBOutlet weak var crashReportingCheckButton: NSButton?
    @IBAction func crashReportingCheckButtonActionReceived(checkbox: NSButton) {
        // this checkbox is either ON or OFF ("mixed" state is not allowed on it)
        let autoSend: Bool = checkbox.state == NSOnState
        if Preferences.sharedInstance.setAutoSendCrashReports(autoSend) == false {
            Logger.log(level: ASLLogger.Level.Fatal, content: "internal error (\(#function))")
            checkbox.state = Preferences.sharedInstance.getAutoSendCrashReports() ? NSOnState : NSOffState
        }
    }
    
    override func viewDidLoad() {
        helpLabel?.stringValue = "settingsPrivacyHelpLabel".localize()
        analyticsCheckButton?.title = "settingsPrivacyUsageInfoLabel".localize()
        analyticsDescription?.text = "settingsPrivacyUsageInfoDescription".localize()
        crashReportingCheckButton?.title = "settingsPrivacyCrashReportLabel".localize()
        crashReportingDescription?.text = "settingsPrivacyCrashReportDescription".localize()
        pingInfoDescription?.text = "settingsPrivacyPingInfo".localize()   
    }
    
    override func viewWillAppear() {
        super.viewWillAppear()
        updateCheckButtons()
    }
    
    override func viewDidAppear() {
        super.viewDidAppear()
        NSNotificationCenter.defaultCenter().addObserver(self, selector: #selector(SettingsPrivacy.preferencesDidChange(_:)), name: Preferences.preferencesDidChangeNotification, object: nil)
    }
    
    override func viewWillDisappear() {
        super.viewWillDisappear()
        NSNotificationCenter.defaultCenter().removeObserver(self, name: Preferences.preferencesDidChangeNotification, object: nil)
    }
    
    func preferencesDidChange(_: NSNotification) {
        updateCheckButtons()
    }

    // MARK: SettingsPanelProtocol implementation

    /// If something has to be applied, the panel view controller should
    /// return true.
    func shouldApply() -> Bool {
        return false
    }
    /// Asks panel view controller to apply changes
    func apply() {}
    
    // MARK: private functions
    
    private func updateCheckButtons() {
        analyticsCheckButton?.state = Preferences.sharedInstance.getAnalyticsEnabled() ? NSOnState : NSOffState
        crashReportingCheckButton?.state = Preferences.sharedInstance.getAutoSendCrashReports() ? NSOnState : NSOffState
    }
}

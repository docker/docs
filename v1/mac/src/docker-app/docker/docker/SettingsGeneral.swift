//
//  SettingsGeneral.swift
//  docker
//
//  Created by Doby Mock on 6/27/16.
//  Copyright © 2016 docker. All rights reserved.
//

import Foundation
import Cocoa

class SettingsGeneral: NSViewController, SettingsPanelProtocol {
    
    @IBOutlet weak var delegate: Settings?
    @IBOutlet weak var channelLabel: TextViewInView?
    
    // PRIVATE attributes
    
    private var _currentBackendSettings = BackendSettings()
    private var currentBackendSettings: BackendSettings {
        get {return self._currentBackendSettings}
        set {
            // update the current value
            self._currentBackendSettings = newValue
            // refresh the UI accordingly
            updateBackendSettingComponents()
            // check for changes in advanced settings
            checkIfBackendSettingsChanged()
        }
    }
    
    private var _initialBackendSettings: BackendSettings? = nil
    private var initialBackendSettings: BackendSettings? {
        get {return self._initialBackendSettings}
        set {self._initialBackendSettings = newValue}
    }
    
    
    // checkbox: automatically start on login
    @IBOutlet weak var autoStartOnLoginCheckButton: NSButton?
    @IBAction func autoStartOnLoginCheckButtonActionReceived(checkbox: NSButton) {
        // this checkbox is either ON or OFF ("mixed" state is not allowed on it)
        let autoStart: Bool = checkbox.state == NSOnState
        if Preferences.sharedInstance.setAutoStart(autoStart) == false {
            Logger.log(level: ASLLogger.Level.Fatal, content: "internal error (\(#function))")
        }
    }
    // checkbox: automatically check for updates
    @IBOutlet weak var autoCheckForUpdatesCheckButton: NSButton?
    @IBAction func autoCheckForUpdatesCheckButtonActionReceived(checkbox: NSButton) {
        // this checkbox is either ON or OFF ("mixed" state is not allowed on it)
        let autoUpdate: Bool = checkbox.state == NSOnState
        // update "auto update" setting of UpdateManager
        UpdateManager.sharedManager().automaticallyChecksForUpdates = autoUpdate
        // update preferences
        if Preferences.sharedInstance.setCheckForUpdates(autoUpdate) == false {
            Logger.log(level: ASLLogger.Level.Fatal, content: "internal error (\(#function))")
        }
    }
    // checkbox: exclude VM from TimeMachine backups
    @IBOutlet weak var excludeVmFromTimeMachineCheckButton: NSButton?
    @IBAction func excludeVmFromTimeMachineCheckButtonActionReceived(checkbox: NSButton) {
        // this checkbox is either ON or OFF ("mixed" state is not allowed on it)
        let excludeVm: Bool = checkbox.state == NSOnState
        if excludeVm {
            // try to exclude vm from Time Machine
            if Install.excludeVmContentFromTimeMachineBackups(true) == false {
                // if it fails, we programmatically uncheck the checkbox
                checkbox.state = NSOffState
            }
        } else {
            // try to un-exclude vm from Time Machine
            if Install.excludeVmContentFromTimeMachineBackups(false) == false {
                // if it fails, we programmatically uncheck the checkbox
                checkbox.state = NSOnState
            }
        }
    }
    
    // RAM & CPU sliders
    
    // textfield: displays amount of memory, based on memory slider
    @IBOutlet weak var memoryText: NSTextField?
    // slider: changes memory allocated for the VM
    @IBOutlet weak var memorySlider: NSSlider?
    // textfield: displays amount of memory, based on memory slider
    @IBOutlet weak var cpuText: NSTextField?
    // slider: changes number of cpus allocated
    @IBOutlet weak var cpuSlider: NSSlider?
    // called when a slider has changed (memory or cpus for now)
    @IBAction func sliderActionReceived(slider: NSSlider) {
        // test if it is the cpus slider
        if let cpuSliderValue = cpuSlider {
            if slider == cpuSliderValue {
                let newCurrentValue = BackendSettings(value: currentBackendSettings)
                newCurrentValue.cpus = slider.doubleValue
                currentBackendSettings = newCurrentValue
                return
            }
        }
        // test if it is the memory slider
        if let memorySliderValue = memorySlider {
            if slider == memorySliderValue {
                let newCurrentValue = BackendSettings(value: currentBackendSettings)
                newCurrentValue.memory = slider.doubleValue
                currentBackendSettings = newCurrentValue
                return
            }
        }
    }
    
    @IBOutlet weak var applyButton: NSButton?
    @IBAction func applyButtonPressed(sender: NSButton) {
        lockBackendComponents()
        sender.enabled = false
        // currentBackendSettings always reflect the current state of the UI,
        // therefore we can use it directly to send a request to the backend.
        // The UI stores the memory in bytes while the backend is expecting
        // memory value in Gigabytes. We have to do the conversion.
        
        let memoryInGigabytes: UInt = UInt(Utils.bytesToGigaBytes(currentBackendSettings.memory))
        // TODO: gdevillele: replace the struct parameter of the following call by an AdvancedSettingsValue object
        Backend.vmSetSettings((memory: memoryInGigabytes, cpus: UInt(currentBackendSettings.cpus), daemonjson: nil, proxyHttpOverride: nil, proxyHttpsOverride: nil, proxyExcludeOverride: nil)) { (error) in
            if error != nil {
                self.delegate?.showErrorPopupWithTitle("Error".localize(), andMessage: "settingsApiApplySettingsRequestError".localize())
            } else {
                // update advancedSettingsInitialValues using UI component values
                self.initialBackendSettings = BackendSettings(value: self.currentBackendSettings)
            }
            self.unlockBackendComponents()
        }
        Backend.enforceStarting()
    }
    
    override func viewWillAppear() {
        super.viewWillAppear()
        
        // UPDATE IBOutlets
        
        // set autoStart checkbox state
        autoStartOnLoginCheckButton?.state = Preferences.sharedInstance.getAutoStart() ? NSOnState : NSOffState
        
        // set checkForUpdates checkbox state
        autoCheckForUpdatesCheckButton?.state = Preferences.sharedInstance.getCheckForUpdates() ? NSOnState : NSOffState
        
        // in dev build, autocheck for updates is disabled
        // it's also not available to non-admin users
        if !UpdateManager.isAppUpdateAvailable() {
            autoCheckForUpdatesCheckButton?.state = NSOffState
            autoCheckForUpdatesCheckButton?.enabled = false
        }
        
        // update "exclude VM from Time Machine" checkbox
        let vmIsExcluded: Bool = Install.isVmExcludedFromTimeMachineBackups()
        excludeVmFromTimeMachineCheckButton?.state = vmIsExcluded ? NSOnState : NSOffState
        
        // display channel type
        let currentChannel = Utils.getChannel() == "stable" ? "stable" : "beta"
        let nextChannel = currentChannel != "stable" ? "stable" : "beta"
        channelLabel?.text = "settingsChannelText".localize(currentChannel, nextChannel)
        channelLabel?.onClickedLink = { (a, b) in
            var goToLink = false
            if let message = WizardMessageGeneric(message: "switchChannelConfirmationTitle".localize(nextChannel), details: "switchChannelConfirmationMessage".localize(currentChannel, nextChannel), icon: "MsgWarningIcon") {
                message.addButton("Cancel".localize()).blockAction = {(button: WizardButton?, message: WizardMessage?) -> () in
                    message?.close()
                }
            
                message.addButton("Ok".localize()).makeDefault().blockAction = {(button: WizardButton?, message: WizardMessage?) -> () in
                    message?.close()
                    goToLink = true
                }
                Wizard.show(message)
            }
            return goToLink
        }

    
        // CPU & RAM sliders
        
        lockBackendComponents()
        
        // init sliders with min/max values and tick marks
        if let slider = self.memorySlider {
            slider.maxValue = Double(NSProcessInfo().physicalMemory)
            slider.minValue = Double(Config.minimumVMMemory)
            slider.doubleValue = slider.minValue
            slider.numberOfTickMarks = Utils.bytesToGigaBytes(slider.maxValue - slider.minValue) + 1
            memoryText?.stringValue = "settingsMemoryLabel".localize(Utils.humanReadableByteCount(slider.doubleValue))
        }
        if let slider = self.cpuSlider {
            slider.maxValue = Double(NSProcessInfo().processorCount)
            slider.minValue = 1
            slider.doubleValue = slider.minValue
            slider.numberOfTickMarks = Int(slider.maxValue - slider.minValue) + 1
            cpuText?.stringValue = "settingsCPUsLabel".localize(slider.intValue)
        }
        
        // get current settings from backend & update advInitialValue
        Backend.vmGetSettings { (error, settings) in
            // check for error
            if let errorValue: Backend.Error = error {
                // Failed to get the advanced settings value from the backend.
                // We display an error popup and close the settings window.
                let errorMessage = Backend.errorMessage(errorValue)
                self.delegate?.showErrorPopupAndCloseWithTitle("Error".localize(), andMessage: "settingsApiGeneralRetrieveError".localize(errorMessage))
                return
            }
            // check we received a memory value
            guard let memoryValue: UInt = settings?.memory else {
                self.delegate?.showErrorPopupAndCloseWithTitle("Error".localize(), andMessage: "settingsApiMemoryRetrieveError".localize())
                return
            }
            // check we received a cpus value
            guard let cpusValue: UInt = settings?.cpus else {
                self.delegate?.showErrorPopupAndCloseWithTitle("Error".localize(), andMessage: "settingsApiCpusRetrieveError".localize())
                return
            }
            // NOTE: we don't care about daemonJson in this preference panel

            // now that we got all the values from the backend response,
            // we update self.advCurrentValue
            let memory: Double = Double(Utils.gigaBytesToBytes(memoryValue))
            let cpus: Double = Double(cpusValue)
            let backendSettings: BackendSettings = BackendSettings(memory: memory, cpus: cpus, daemonjson: "", proxyHttpSystem: "", proxyHttpsSystem: "", proxyExcludeSystem: "", proxyHttpOverride: "", proxyHttpsOverride: "", proxyExcludeOverride: "")
            self.initialBackendSettings = BackendSettings(value: backendSettings)
            self.currentBackendSettings = BackendSettings(value: backendSettings)
            
            self.unlockBackendComponents()
        }
    }
    
    func lockBackendComponents() {
        memorySlider?.enabled = false
        cpuSlider?.enabled = false
    }
    
    func unlockBackendComponents() {
        memorySlider?.enabled = true
        cpuSlider?.enabled = true
    }
    
    func updateBackendSettingComponents() {
        // cpus slider
        if let slider = cpuSlider {
            slider.doubleValue = currentBackendSettings.cpus
            // cpus label
            cpuText?.stringValue = "settingsCPUsLabel".localize(slider.intValue)
        }
        // memory slider
        if let slider = memorySlider {
            slider.doubleValue = currentBackendSettings.memory
            // memory label
            memoryText?.stringValue = "settingsMemoryLabel".localize(Utils.humanReadableByteCount(currentBackendSettings.memory))
        }
    }
    
    func checkIfBackendSettingsChanged() -> Bool {
        guard let initialValue = initialBackendSettings else {
            // we don't have initial values. It means we never got them from the
            // backend.
            return false
        }
        let changed: Bool = currentBackendSettings != initialValue
        applyButton?.enabled = changed
        return changed
    }
    
    // SettingsPanelProtocol
    
    func shouldApply() -> Bool {
        return checkIfBackendSettingsChanged()
    }
    
    func apply() {
        if let applyButton = applyButton {
            applyButtonPressed(applyButton)
        }
    }
    
}

//
//  SettingsAdvanced.swift
//  docker
//
//  Created by Doby Mock on 6/27/16.
//  Copyright © 2016 docker. All rights reserved.
//

import Foundation
import Cocoa

class SettingsAdvanced: NSViewController, NSTextFieldDelegate, SettingsPanelProtocol, TextTableViewEventListener {

    @IBOutlet weak var settingsWindow: Settings?
    @IBOutlet weak var insecureRegistriesTableView: TextTableView?
    @IBOutlet weak var registryMirrorsTableView: TextTableView?
    @IBOutlet weak var proxyHttpTextField: NSTextField?
    @IBOutlet weak var proxyHttpsTextField: NSTextField?
    @IBOutlet weak var proxyExcludeTextField: NSTextField?
    
    // PRIVATE attributes
    
    private var _currentBackendSettings = BackendSettings()
    private var currentBackendSettings: BackendSettings {
        get {return self._currentBackendSettings}
        set {
            self._currentBackendSettings = newValue
            checkIfBackendSettingsChanged()
        }
    }
    
    private var _initialBackendSettings: BackendSettings? = nil
    private var initialBackendSettings: BackendSettings? {
        get {return self._initialBackendSettings}
        set {self._initialBackendSettings = newValue}
    }
    
    // insecure-registries add button
    @IBOutlet weak var insecureRegistriesAddButton: NSButton?
    @IBAction func insecureRegistriesAddActionReceived(_: NSButton) {
        insecureRegistriesTableView?.appendString("", edit: true)
    }
    
    // insecure-registries remove button
    @IBOutlet weak var insecureRegistriesRemoveButton: NSButton?
    @IBAction func insecureRegistriesRemoveActionReceived(_: NSButton) {
        insecureRegistriesTableView?.removeSelectedString()
    }
    
    // registry-mirrors add button
    @IBOutlet weak var registryMirrorsAddButton: NSButton?
    @IBAction func registryMirrorsAddActionReceived(_: NSButton) {
        registryMirrorsTableView?.appendString("", edit: true)
    }
    
    // registry-mirrors remove button
    @IBOutlet weak var registryMirrorsRemoveButton: NSButton?
    @IBAction func registryMirrorsRemoveActionReceived(_: NSButton) {
        registryMirrorsTableView?.removeSelectedString()
    }
    
    // button: apply advanced settings (disabled if settings didn't change)
    @IBOutlet weak var applyButton: NSButton?
    @IBAction func applyButtonPressed(sender: NSButton) {
        lockBackendComponents()
        sender.enabled = false

        guard let daemonJson: String = currentBackendSettings.daemonCleanJson else {
            self.settingsWindow?.showErrorPopupWithTitle("Error".localize(), andMessage: "settingsApiCannotPrepareRequestError".localize())
            self.unlockBackendComponents()
            return
        }
        let proxyHttpValue = proxyHttpTextField?.stringValue
        let proxyHttpsValue = proxyHttpsTextField?.stringValue
        let proxyExcludeValue = proxyExcludeTextField?.stringValue
        
        Backend.vmSetSettings((memory: nil, cpus: nil, daemonjson: daemonJson, proxyHttpOverride: proxyHttpValue, proxyHttpsOverride: proxyHttpsValue, proxyExcludeOverride: proxyExcludeValue)) { (error) in
            if error != nil {
                self.settingsWindow?.showErrorPopupWithTitle("Error".localize(), andMessage: "settingsApiApplySettingsRequestError".localize())
            } else {
                // update advancedSettingsInitialValues using UI component values
                self.initialBackendSettings = BackendSettings(value: self.currentBackendSettings)
                self.updateProxyUI()
            }
            self.unlockBackendComponents()
        }
        Backend.enforceStarting()
    }
    
    override func viewWillAppear() {
        super.viewWillAppear()
        
        // setup TextTableViews
        insecureRegistriesTableView?.eventListener = self
        registryMirrorsTableView?.eventListener = self
        insecureRegistriesTableView?.placeholderString = "mydomain.com:5000"
        registryMirrorsTableView?.placeholderString = "mydomain.com:5000"
        
        // setup HTTP proxy TextFields
        proxyHttpTextField?.delegate = self
        proxyHttpsTextField?.delegate = self
        proxyExcludeTextField?.delegate = self
        
        lockBackendComponents()
        
        // get current settings from backend & update advInitialValue
        Backend.vmGetSettings { (error, settings) in
            // check for error
            if let errorValue: Backend.Error = error {
                // Failed to get the advanced settings value from the backend.
                // We display an error popup and close the settings window.
                let errorMessage = Backend.errorMessage(errorValue)
                self.settingsWindow?.showErrorPopupAndCloseWithTitle("Error".localize(), andMessage: "settingsApiGeneralRetrieveError".localize(errorMessage))
                return
            }
            // check we received a daemonjson value
            guard let settingsValue = settings else {
                self.settingsWindow?.showErrorPopupAndCloseWithTitle("Error".localize(), andMessage: "settingsApiRetrieveError".localize())
                return
            }
            // now that we got all the values from the backend response,
            // we update self.advCurrentValue
            let backendValue = BackendSettings(memory: 0, cpus: 0, daemonjson: settingsValue.daemonjson,
                proxyHttpSystem: settingsValue.proxyHttpSystem,
                proxyHttpsSystem: settingsValue.proxyHttpsSystem,
                proxyExcludeSystem: settingsValue.proxyExcludeSystem,
                proxyHttpOverride: settingsValue.proxyHttpOverride,
                proxyHttpsOverride: settingsValue.proxyHttpsOverride,
                proxyExcludeOverride: settingsValue.proxyExcludeOverride)
            self.initialBackendSettings = BackendSettings(value: backendValue)
            self.currentBackendSettings = BackendSettings(value: backendValue)
            self.updateBackendSettingComponents()
            self.unlockBackendComponents()
        }
    }
    
    override func viewDidDisappear() {
        super.viewDidDisappear()
        registryMirrorsTableView?.deselectAll(self)
        insecureRegistriesTableView?.deselectAll(self)
    }
    
    func lockBackendComponents() {
        insecureRegistriesTableView?.enabled = false
        insecureRegistriesAddButton?.enabled = false
        insecureRegistriesRemoveButton?.enabled = false
        registryMirrorsTableView?.enabled = false
        registryMirrorsAddButton?.enabled = false
        registryMirrorsRemoveButton?.enabled = false
        // proxy settings UI
        proxyHttpTextField?.enabled = false
        proxyHttpsTextField?.enabled = false
        proxyExcludeTextField?.enabled = false
    }
    
    func unlockBackendComponents() {
        insecureRegistriesTableView?.enabled = true
        insecureRegistriesAddButton?.enabled = true
        registryMirrorsTableView?.enabled = true
        registryMirrorsAddButton?.enabled = true
        // proxy settings UI
        proxyHttpTextField?.enabled = true
        proxyHttpsTextField?.enabled = true
        proxyExcludeTextField?.enabled = true
    }
    
    func updateBackendSettingComponents() {
        if let strings = self.currentBackendSettings.insecureRegistry {
            self.insecureRegistriesTableView?.strings = strings
        }
        if let strings = self.currentBackendSettings.registryMirror {
            self.registryMirrorsTableView?.strings = strings
        }
        self.updateProxyUI()
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
    
    // MARK: TextTableViewEventListener implementation
    
    func textTableViewGainedFocus(textTableView: TextTableView) {
        if textTableView === insecureRegistriesTableView {
            registryMirrorsTableView?.deselectAll(self)
            registryMirrorsRemoveButton?.enabled = false
        }
        if textTableView === registryMirrorsTableView {
            insecureRegistriesTableView?.deselectAll(self)
            insecureRegistriesRemoveButton?.enabled = false
        }
    }
    
    func textTableViewLostFocus(textTableView: TextTableView) {}
    
    func textTableViewContentDidChange(textTableView: TextTableView) {
        
        let newCurrentValue = BackendSettings(value: currentBackendSettings)
        
        if textTableView === insecureRegistriesTableView {
            newCurrentValue.insecureRegistry? = textTableView.strings
        } else if textTableView === registryMirrorsTableView {
            newCurrentValue.registryMirror? = textTableView.strings
        }
        
        currentBackendSettings = newCurrentValue
    }
    
    func textTableViewCanRemoveRow(textTableView: TextTableView, canRemoveRow: Bool) {
        if textTableView === insecureRegistriesTableView {
            insecureRegistriesRemoveButton?.enabled = canRemoveRow
        } else if textTableView === registryMirrorsTableView {
            registryMirrorsRemoveButton?.enabled = canRemoveRow
        }
    }
    
    func textTableViewCanAddRow(textTableView: TextTableView, canAddRow: Bool) {
        if textTableView === insecureRegistriesTableView {
            insecureRegistriesAddButton?.enabled = canAddRow
        } else if textTableView === registryMirrorsTableView {
            registryMirrorsAddButton?.enabled = canAddRow
        }
    }
    
    // MARK: NSTextFieldDelegate implementation
    
    override func controlTextDidChange(notification: NSNotification) {
        if let object: AnyObject = notification.object {
            if let textfield: NSTextField = object as? NSTextField {
                let newCurrentValue = BackendSettings(value: self.currentBackendSettings)
                if textfield === self.proxyHttpTextField {
                    newCurrentValue.proxyHttpOverride = textfield.stringValue
                } else if textfield === self.proxyHttpsTextField {
                    newCurrentValue.proxyHttpsOverride = textfield.stringValue
                } else if textfield === self.proxyExcludeTextField {
                    newCurrentValue.proxyExcludeOverride = textfield.stringValue
                }
                currentBackendSettings = newCurrentValue
            }
        }
    }
    
    // MARK: utility functions
    
    private func updateProxyUI() {
        // test if user did override system HTTP proxy settings
        if currentBackendSettings.proxyHttpOverride.isEmpty
            && currentBackendSettings.proxyHttpsOverride.isEmpty
            && currentBackendSettings.proxyExcludeOverride.isEmpty {
            // no override. We use system values as placeholder strings
            self.displaySystemProxy()
        } else {
            // there is override, we use override value and no placeholder strings
            self.displayOverrideProxy()
        }
    }
    
    private func displaySystemProxy() {
        self.proxyHttpTextField?.stringValue = ""
        self.proxyHttpTextField?.placeholderString = currentBackendSettings.proxyHttpSystem
        self.proxyHttpsTextField?.stringValue = ""
        self.proxyHttpsTextField?.placeholderString = currentBackendSettings.proxyHttpsSystem
        self.proxyExcludeTextField?.stringValue = ""
        self.proxyExcludeTextField?.placeholderString = currentBackendSettings.proxyExcludeSystem
    }
    
    private func displayOverrideProxy() {
        self.proxyHttpTextField?.stringValue = currentBackendSettings.proxyHttpOverride
        self.proxyHttpTextField?.placeholderString = ""
        self.proxyHttpsTextField?.stringValue = currentBackendSettings.proxyHttpsOverride
        self.proxyHttpsTextField?.placeholderString = ""
        self.proxyExcludeTextField?.stringValue = currentBackendSettings.proxyExcludeOverride
        self.proxyExcludeTextField?.placeholderString = ""
    }
}

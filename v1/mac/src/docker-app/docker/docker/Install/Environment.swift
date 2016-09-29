//
//  Environment.swift
//  docker-installer
//
//  Created by Michel Courtine on 12/5/15.
//  Copyright © 2015 David Scott. All rights reserved.
//

import Foundation

class Environment {

    struct Warning {
        var title: String
        var message: String
        var ignoreBtn: String
    }

    enum CheckResult {
        case OK
        case Fatal(String)
        case Warning(Environment.Warning)
    }

    // MARK: public functions

    static func Check() -> CheckResult {
        if let error: String = checkOperationSystemVersion() {
            return CheckResult.Fatal(error)
        }
        if let error: String = checkCpuIsSupported() {
            return CheckResult.Fatal(error)
        }
        if let error: String = checkRamAmount() {
            return CheckResult.Fatal(error)
        }
        if let error: String = checkVirtualBoxVersion() {
            return CheckResult.Fatal(error)
        }
        if let warn: String = checkBlueStacks() {
            let title = "CheckBlueStacksTitle".localize()
            let ignoreBtn = "CheckBlueStacksIgnore".localize()
            let warning = Warning(title: title, message: warn, ignoreBtn: ignoreBtn)
            return CheckResult.Warning(warning)
        }
        return CheckResult.OK
    }
    
    // MARK: private functions
    
    // checkOperationSystemVersion makes sure the minimum OS X version is installed.
    // Returns nil or an error message.
    static private func checkOperationSystemVersion() -> String? {
        let osMajor = Utils.getOSXMajorVersion()
        let osMinor = Utils.getOSXMinorVersion()
        let osPatch = Utils.getOSXPatchVersion()
        Logger.log(level: Logger.Level.Notice, content: "OSX Version: \(Utils.getOSXVersionString())")
        if (osMajor < Config.OSXMinimumMajor)
            || (osMajor == Config.OSXMinimumMajor && osMinor < Config.OSXMinimumMinor)
            || (osMajor == Config.OSXMinimumMajor && osMinor == Config.OSXMinimumMinor && osPatch < Config.OSXMinimumPatch) {
            return "CheckHostVersionError".localize() + "\n\n" + "CheckHostVersionErrorDetails".localize(Config.OSXMinimumMajor, Config.OSXMinimumMinor, Config.OSXMinimumPatch, Utils.getOSXVersionString())
        }
        return nil
    }
    
    // checkCpuIsSupported makes sure the CPU is compatible with Hypervisor.framework
    static private func checkCpuIsSupported() -> String? {
        let command = "sysctl kern.hv_support | grep 1"
        let exitCode: Int32 = Execute(command)
        if (exitCode != 0) {
            return "CheckHostCPUError".localize() + "\n\n" + "CheckHostCPUErrorDetails".localize()
        }
        return nil
    }
    
    // checkRamAmount makes sure there's enough RAM
    static private func checkRamAmount() -> String? {
        let physicalMemory: UInt64 = NSProcessInfo().physicalMemory
        let hostMemoryRequired: UInt64 = Config.minimumMemoryRequirement
        if (physicalMemory < hostMemoryRequired){
            return "CheckHostRAMError".localize() + "\n\n" + "CheckHostRAMErrorDetails".localize(hostMemoryRequired, physicalMemory)
        }
        return nil
    }
    
    // Xhyve + Virtualbox 4.x crashes, so we should check the virtualbox version
    // when installing Docker.app so we can warn the user.
    static private func checkVirtualBoxVersion() -> String? {
        
        let virtualBoxMinimumMajor = 4
        let virtualBoxMinimumMinor = 3
        let virtualBoxMinimumPatch = 30
        
        var detectedMajor = 0
        var detectedMinor = 0
        var detectedPatch = 0
        
        // check Virtualbox's .kext files:
        
        let task = NSTask()
        let out = NSPipe()
        // let err = NSPipe()
        task.launchPath = "/bin/sh"
        task.arguments = ["-c", "kextstat | grep -i virtualbox | grep -o '(.*)' | grep -o '[0-9].*[0-9]'"]
        task.standardOutput = out
        // task.standardError = err
        task.launch()
        task.waitUntilExit()
        
        let data = out.fileHandleForReading.readDataToEndOfFile()
        out.fileHandleForReading.closeFile()
        guard var output = String.init(data: data, encoding: NSUTF8StringEncoding) else {
            return "Internal error: \(#function)"
        }
        
        // error not handled so far
        
        output = output.stringByReplacingOccurrencesOfString("\r", withString: "")
        output = output.stringByReplacingOccurrencesOfString("\0", withString: "")
        output = output.stringByTrimmingCharactersInSet(NSCharacterSet.whitespaceAndNewlineCharacterSet())
        
        // If there's no output, it means there are no .kext files for
        // virtualbox, so it's safe to continue.
        if output == "" {
            return nil
        }
        
        // If .kext files can be found, the output looks like this:
        // "4.3.36\n4.3.36\n4.3.36\n4.3.36\n"
        // We can split on '\n' and check major, minor and patch numbers
        
        let versions = output.componentsSeparatedByString("\n")
        
        for version in versions {
            
            let versionComponents = version.componentsSeparatedByString(".")
            
            if versionComponents.count > 0 {
                if let majorComponent = Int(versionComponents[0]) {
                    detectedMajor = majorComponent
                }
            }
            
            if versionComponents.count > 1 {
                if let minorComponent = Int(versionComponents[1]) {
                    detectedMinor = minorComponent
                }
            }
            
            if versionComponents.count > 2 {
                if let patchComponent = Int(versionComponents[2]) {
                    detectedPatch = patchComponent
                }
            }
        }
    
    
        // check major
        if detectedMajor > virtualBoxMinimumMajor {
            return nil
        }
        
        // check minor
        if detectedMajor == virtualBoxMinimumMajor && detectedMinor > virtualBoxMinimumMinor {
            return nil
        }
        
        // check patch
        if detectedMajor == virtualBoxMinimumMajor && detectedMinor == virtualBoxMinimumMinor
            && detectedPatch >= virtualBoxMinimumPatch {
            return nil
        }
        
        return "CheckVBoxVersionError".localize(virtualBoxMinimumMajor, virtualBoxMinimumMinor, virtualBoxMinimumPatch, detectedMajor, detectedMinor, detectedPatch)
    }

    // Hypervisor.framework + BlueStacks panics OS X so refuse to
    // start if we detect BlueStacks
    static private func checkBlueStacks() -> String? {
        // We could also check for the process called
        // com.BlueStacks.AppPlayer.bstservice_helper

        // Or the kext called
        // com.bluestacks.kext.Hypervisor (4.3.26)

        // But even with these missing, the machine can still panic if
        // BlueStacks has been used in the same OS X boot so we check
        // for the installation only.

        let fileManager = NSFileManager()

        if fileManager.fileExistsAtPath("/Applications/BlueStacks.app") {
            return "CheckBlueStacksError".localize()
        }

        return nil
    }
}

//
//  MachineMigrator.swift
//  docker
//
//  Created by Emmanuel Briney on 22/01/16.
//  Copyright © 2016 docker. All rights reserved.
//

import Foundation
import Cocoa

// QemuProgressParser parses the output of `qemu-img convert -p` and extracts the
// current %age complete from the text.
class QemuProgressParser {
    
    // The output of `qemu-img convert -p` has the following structure:
    //   <whitespace>(%.2f/100%)\r
    // We are only interested in the last whole line, ending in \r. We
    // may read partial lines so we need to remember the text so far.
    private var currentRecord = ""
    
    var percentageComplete = 0.0
    
    func input(text: String) {
        currentRecord = currentRecord + text
        for record in currentRecord.componentsSeparatedByString("\r") {
            if let _ = record.characters.indexOf("%") {
                if let openparen = record.characters.indexOf("(") {
                    if let slash = record.characters.indexOf("/") {
                        let nf = NSNumberFormatter()
                        if let numerator = nf.numberFromString(record.substringWithRange(Range<String.Index>(openparen.advancedBy(1)..<slash.advancedBy(-1)))) {
                            percentageComplete = numerator.doubleValue
                        }
                    }
                }
            } else {
                // no '%' means record is incomplete: store it for future completion
                currentRecord = record
            }
        }
    }
}

func getQcowPath() -> String {
    guard let appContainerPath = Paths.appContainerPath() else {
        return "failed to get appContainerPath"
    }
    
    // check that <appContainer>/com.docker.driver.amd64-linux directory
    // exists and create it if it doesn't (on a fresh install's first
    // launch, it won't exist yet). Migration will fail if the the directory
    // doesn't exist.
    let driverDirectoryPath: String = NSString.pathWithComponents([appContainerPath, "com.docker.driver.amd64-linux"])
    if let _: String = Paths.CreateDirectoryIfNecessary(driverDirectoryPath) {
        // there was an error creating the driver directory
        return "could not create driver directory"
    }
    
    // try to migrate machine volume to Docker.qcow2 file
    let cowPath: String = NSString.pathWithComponents([driverDirectoryPath, "Docker.qcow2"])
    return cowPath
}

// MachineMigrator migrates data from local machines, to be used in
// the xhyve vm. It inherits from NSObject for methodSignatureForSelector:
// (required when setting target and action for buttons)
class MachineMigrator: NSObject {
    
    // Different possible statuses for local machine
    // migration. For simplicity, once migration has
    // been done or declined, or if no machine is found
    // during first launch, we don't prompt the user again
    // for migration.
    enum localMachineMigrationStatus: UInt8 {
        case None
        case Done
        case Failed
        case Declined
        case NoMachineFound
        case VolumeTooBigForQcow
    }
    
    class Machine {
        var Name: String = ""
        var Driver: String = ""
        var Path: String = ""
        var Volume: String = ""
        
        func getVolumeSize() -> UInt64 {
            var fileSize: UInt64 = UInt64.max / 2
            do {
                let attr: NSDictionary? = try NSFileManager.defaultManager().attributesOfItemAtPath(Volume)
                if let _attr = attr {
                    fileSize = _attr.fileSize()
                }
            } catch {
                Logger.log(level: ASLLogger.Level.Error, content: "Failed to get vmdk size: \(error)")
            }
            return fileSize
        }
        
        func neededSizeToMigrate() -> UInt64 {
            guard let appContainerPath = Paths.appContainerPath() else {
                return 0
            }
            let reservedSize: UInt64 = 10000000
            var remainingFreeSpace: UInt64 = 0
            do {
                let fileAttributes = try NSFileManager.defaultManager().attributesOfFileSystemForPath(appContainerPath)
                if let size = fileAttributes[NSFileSystemFreeSize] as? NSInteger {
                    remainingFreeSpace = UInt64(size)
                }
            } catch {
                Logger.log(level: ASLLogger.Level.Error, content: "Failed to get FileSystemFreeSize: \(error)")
            }
            
            let volumeSize = getVolumeSize()
            if remainingFreeSpace > volumeSize + reservedSize {
                return 0
            }
            
            return volumeSize + reservedSize - remainingFreeSpace
        }
        
        func isTooBigForQcow() -> Bool {
            let qcowRemainingVirtualSize: UInt64 = 54*1024*1024*1024
            return getVolumeSize() > qcowRemainingVirtualSize
        }
    }
    
    private var migrationMessage: WizardMessage? = nil
    private var machineMigrated: MachineMigrator.Machine? = nil
    
    // progressBar is used to give some feedback about migration status
    private var progressBar: NSProgressIndicator? = nil
    
    // migrateLocalMachine looks for existing local machines (the ones using
    // VirtualBox driver), and prompts the user to migrate data. Only one
    // machine can be migrated so far.
    // returns nil or an error
    func migrateLocalMachine() -> String? {
        
        // if the app has been launched by automated tests we ignore the machine
        // migration step
        if Options.unattended {
            return nil
        }
        
        let migrationStatusUInt: UInt = Preferences.sharedInstance.getLocalMachineMigration()
        let migrationStatus: UInt8 = UInt8(migrationStatusUInt)
        if migrationStatus != localMachineMigrationStatus.None.rawValue && migrationStatus != localMachineMigrationStatus.Failed.rawValue {
            // no error, but we return because migration has already been done,
            // declined or no existing local machine was found.
            return nil
        }
        
        let machines = list()
        if machines.count == 0 {
            // no local machine found, update localMachineMigrationStatus
            // and return with no error
            if Preferences.sharedInstance.setLocalMachineMigration(UInt(localMachineMigrationStatus.NoMachineFound.rawValue)) == false {
                return "failed to update machine migration state"
            }
            return nil
        }
        
        // choose machines to convert, for now default but if multiple machines the user should have to select which ones he want to convert
        for machine in machines {
            if machine.Name == "default" {
                machineMigrated = machine
                break
            }
        }
        
        if machineMigrated == nil {
            // no local machine found, update localMachineMigrationStatus
            // and return with no error
            if Preferences.sharedInstance.setLocalMachineMigration(UInt(localMachineMigrationStatus.NoMachineFound.rawValue)) == false {
                return "failed to update machine migration state"
            }
            return nil
        }
        
        migrationMessage = WizardMessageGeneric(message: "toolboxMigrationTitle".localize(), details: "toolboxMigrationMessage".localize(), icon: "Toolbox")
        
        guard let migrationMessage = migrationMessage else {
            Logger.log(level: ASLLogger.Level.Fatal, content: "Internal error: \(#function)")
            return nil
        }
        
        let noBtn = migrationMessage.addButton("toolboxMigrationNoButton".localize())
        noBtn.target = self
        noBtn.action = #selector(MachineMigrator.answerIsNo)
        let yesBtn = migrationMessage.addButton("toolboxMigrationMigrateButton".localize())
        yesBtn.target = self
        yesBtn.action = #selector(MachineMigrator.answerIsYes)
        yesBtn.makeDefault()
        
        Wizard.show(migrationMessage)
        
        return nil
    }
    
    func answerIsNo() {
        // Logger.log(level: Logger.Level.Notice, content: "no")
        if Preferences.sharedInstance.setLocalMachineMigration(UInt(localMachineMigrationStatus.Declined.rawValue)) == false {
            Logger.log(level: ASLLogger.Level.Error, content: "failed to update machine migration state")
        }
        migrationMessage?.close()
        migrationMessage = nil
    }
    
    func formatDiskSize(size: UInt64) -> String {
        // KB in Finder are 1000 bytes and not 1024
        let GB: UInt64 = 1000*1000*1000
        let MB: UInt64 = 1000*1000
        if size > GB {
            return String(format: "%.1lf GB", Double(size) / Double(GB))
        }
        return String(format: "%.1lf MB", Double(size) / Double(MB))
    }
    
    func answerIsYes() {
        // Logger.log(level: Logger.Level.Notice, content: "yes")
        
        // Here we test if there is enough disk space
        if let machine = machineMigrated {
            if machine.isTooBigForQcow() {
                guard let isTooBigForQcowMessage = WizardMessageGeneric(message: "toolboxMigrationTitle".localize(), details: "toolboxMigrationVolumeTooBigMessage".localize(), icon: "Toolbox") else {
                    Logger.log(level: ASLLogger.Level.Fatal, content: "Internal error: \(#function)")
                    return
                }
                isTooBigForQcowMessage.addCloseButton("Ok".localize()).makeDefault()
                migrationMessage?.pushAfter(isTooBigForQcowMessage)
                migrationMessage?.close()
                migrationMessage = isTooBigForQcowMessage
                if Preferences.sharedInstance.setLocalMachineMigration(UInt(localMachineMigrationStatus.VolumeTooBigForQcow.rawValue)) == false {
                    Logger.log(level: ASLLogger.Level.Fatal, content: "Internal error: \(#function)")
                }
                return
            }

            let neededSize = machine.neededSizeToMigrate()
            if neededSize > 0 {
                guard let notEnoughSpaceMessage = WizardMessageGeneric(message: "toolboxMigrationTitle".localize(), details: "toolboxMigrationNotEnoughSpaceMessage".localize(formatDiskSize(neededSize)), icon: "Toolbox") else {
                    Logger.log(level: ASLLogger.Level.Fatal, content: "Internal error: \(#function)")
                    return
                }
                notEnoughSpaceMessage.addCloseButton("Ok".localize()).makeDefault()
                migrationMessage?.pushAfter(notEnoughSpaceMessage)
                migrationMessage?.close()
                migrationMessage = notEnoughSpaceMessage
                if Preferences.sharedInstance.setLocalMachineMigration(UInt(localMachineMigrationStatus.Failed.rawValue)) == false {
                    Logger.log(level: ASLLogger.Level.Fatal, content: "Internal error: \(#function)")
                }
                return
            }
        }
        
        guard let message = WizardMessageGeneric(message: "toolboxMigrationTitle".localize(), details: "toolboxMigrationLoading".localize(), icon: "MsgWaitIcon") else {
            Logger.log(level: ASLLogger.Level.Fatal, content: "Internal error: \(#function)")
            return
        }
        
        progressBar = message.addProgressBar(false)
        if let progressBar = progressBar {
            progressBar.minValue = 0
            progressBar.maxValue = 100
            progressBar.doubleValue = 0
        }
        
        migrationMessage?.pushAfter(message)
        migrationMessage?.close()
        migrationMessage = message
        
        if let machine = machineMigrated {
            if let err = migrateMachine(machine) {
                // Migration failure message
                guard let failureMessage = WizardMessageGeneric(message: "toolboxMigrationTitle".localize(), details: "toolboxMigrationFailed".localize(err), icon: "FatalIcon") else {
                    Logger.log(level: ASLLogger.Level.Fatal, content: "Internal error: \(#function)")
                    return
                }
                let closeBtn = failureMessage.addButton("Quit".localize())
                closeBtn.makeDefault()
                closeBtn.blockAction = {(button: WizardButton?, message: WizardMessage?) -> () in
                    message?.close()
                    NSApp.terminate(nil)
                }
                migrationMessage?.pushAfter(failureMessage)
                progressBar = nil
                migrationMessage?.close()
                migrationMessage = nil
            } else {
                // Migration success message
                guard let successMessage = WizardMessageGeneric(message: "toolboxMigrationTitle".localize(), details: "toolboxMigrationSuccess".localize(), icon: "Toolbox") else {
                    Logger.log(level: ASLLogger.Level.Fatal, content: "Internal error: \(#function)")
                    return
                }
                let okBtn = successMessage.addCloseButton("Ok".localize())
                okBtn.makeDefault()
                migrationMessage?.pushAfter(successMessage)
                progressBar = nil
                migrationMessage?.close()
                migrationMessage = nil
            }
        } else {
            Logger.log(level: Logger.Level.Fatal, content: "\(#function): could not get machineMigrated")
        }
    }
    
    // migrateMachine migrates a local machine data to be used by xhyve
    // (for now overwrite Docker.qcow2 volume)
    // returns nil or an error
    private func migrateMachine(machine: MachineMigrator.Machine) -> String? {
        
        // aduermael: we don't care about volume modification date
        // as our condition to do the migration only relies on
        // localMachineMigrationStatus. We don't take responsability for
        // possible changes if the VirtualBox local machine is used
        // again.
        //        var fileModificationDate: String
        //        do {
        //            let machineVolumeAttributes = try NSFileManager.defaultManager().attributesOfItemAtPath(machine.Volume)
        //            fileModificationDate = NSString(format: "%@", machineVolumeAttributes["NSFileModificationDate"] as! NSObject) as String
        //        } catch {
        //            return nil
        //        }
        
        // try to migrate machine volume to Docker.qcow2 file
        let cowPath = getQcowPath()
        if let err: String = migrateVolume(machine, to: cowPath) {
            // Machine migration failed.
            // Store a value in preferences, indicating a failed machine migration.
            if Preferences.sharedInstance.setLocalMachineMigration(UInt(localMachineMigrationStatus.Failed.rawValue)) == false {
                return "failed to update machine migration state"
            }
            return err
        }
        
        // Machine migration succeeded.
        // Store a value in preferences, indicating a successful machine migration.
        if Preferences.sharedInstance.setLocalMachineMigration(UInt(localMachineMigrationStatus.Done.rawValue)) == false {
            return "failed to update machine migration state"
        }
        return nil
    }
    
    //
    private func list() -> [Machine] {
        var machinesFound: [Machine] = []
        // NOTE(aduermael): NSHomeDirectory() won't work in sandbox
        let machinesDirPath = NSHomeDirectory() + "/.docker/Machine/Machines"
        
        let files = NSFileManager.defaultManager().enumeratorAtPath(machinesDirPath)
        while let o = files?.nextObject() {
            if let file = o as? String {
                let path = machinesDirPath + "/" + file
                var isDir: ObjCBool = ObjCBool(false)
                if NSFileManager.defaultManager().fileExistsAtPath(path, isDirectory:&isDir) {
                    if isDir.boolValue {
                        let volume = path + "/disk.vmdk"
                        if NSFileManager.defaultManager().fileExistsAtPath(volume) {
                            let machine = Machine()
                            machine.Name = file
                            machine.Driver = "VirtuaBox"
                            machine.Path = path
                            machine.Volume = volume
                            machinesFound.append(machine)
                        }
                    }
                }
            }
        }
        return machinesFound
    }
    
    // migrate a volume from a machine to a specified place
    // return errorMessage or nil on success
    private func migrateVolume(machine: Machine, to: String) -> String? {
        let qemuProgressParser = QemuProgressParser()
        let qemuImgConvertTask = NSTask()
        let bundlePath = NSBundle.mainBundle().bundlePath
        let qemuImg = NSString.pathWithComponents([bundlePath, "Contents", "MacOS", "qemu-img"])
        qemuImgConvertTask.launchPath = qemuImg
        qemuImgConvertTask.arguments = [ "convert", "-p", "-f", "vmdk", "-O", "qcow2", machine.Volume, to]
        
        // stdout redirection handler
        let out = NSPipe()
        qemuImgConvertTask.standardOutput = out
        
        // stderr redirection handler
        let err = NSPipe()
        qemuImgConvertTask.standardError = err
        
        out.fileHandleForReading.waitForDataInBackgroundAndNotify()
        // add observer to catch stdout modification
        let observer = NSNotificationCenter.defaultCenter().addObserverForName(NSFileHandleDataAvailableNotification, object: out.fileHandleForReading, queue: nil, usingBlock: {_ in
            if qemuImgConvertTask.running {
                let stdoutData = out.fileHandleForReading.availableData
                if stdoutData.length != 0 {
                    if let stdout = NSString(data: stdoutData, encoding: NSUTF8StringEncoding) as? String {
                        qemuProgressParser.input(stdout)
                        Logger.log(level: Logger.Level.Notice, content: String(format: "progress %.2f", qemuProgressParser.percentageComplete))
                        if let progressBar = self.progressBar {
                            progressBar.doubleValue = qemuProgressParser.percentageComplete
                        }
                    } else {
                        Logger.log(level: Logger.Level.Error, content: "Unable to parse data from stdout (not UTF-8?)")
                    }
                }
                out.fileHandleForReading.waitForDataInBackgroundAndNotify()
            }
        })
        
        qemuImgConvertTask.launch()
        qemuImgConvertTask.waitUntilExit()
        // remove observer from the notification center
        NSNotificationCenter.defaultCenter().removeObserver(observer)
        
        if qemuImgConvertTask.terminationReason == NSTaskTerminationReason.Exit && qemuImgConvertTask.terminationStatus == 0 {
            Logger.log(level: Logger.Level.Notice, content: "Converting machine " + machine.Name + " volume done")
            return nil
        } else {
            let stderrData = err.fileHandleForReading.availableData
            if let stderr = NSString(data: stderrData, encoding: NSUTF8StringEncoding) as? String {
                Logger.log(level: Logger.Level.Error, content: stderr)
                return "Error converting " + machine.Name + " volume: " + stderr
            } else {
                Logger.log(level: Logger.Level.Error, content: "Unable to convert stderr output to UTF-8 - data not shown.")
                return "Error converting " + machine.Name + " - error not UTF-8"
            }
        }
    }
}

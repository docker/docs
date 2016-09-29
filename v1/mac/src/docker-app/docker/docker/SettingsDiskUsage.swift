//
//  SettingsDiskUsage.swift
//  docker
//
//  Created by Doby Mock on 6/27/16.
//  Copyright © 2016 docker. All rights reserved.
//

import Foundation
import Cocoa

class SettingsDiskUsage: NSViewController, SettingsPanelProtocol {

    @IBOutlet weak var diskUsageBar: NSSegmentedControl?
    @IBOutlet weak var hdlabel: NSTextField?
    @IBOutlet weak var dockerlabel: NSTextField?
    @IBOutlet weak var dockerUsageLabel: NSTextField?
    @IBOutlet weak var otherslabel: NSTextField?
    @IBOutlet weak var othersUsagelabel: NSTextField?
    
    override func viewWillAppear() {
        super.viewWillAppear()
        
        if let w0 = diskUsageBar?.widthForSegment(0) {
            if let w1 = diskUsageBar?.widthForSegment(1) {
                if let w2 = diskUsageBar?.widthForSegment(2) {
                    let totalWidth = w0 + w1 + w2
                    
                    let fm: NSFileManager = NSFileManager.defaultManager()
                    
                    if let att: [String : AnyObject] = try? fm.attributesOfFileSystemForPath("/") {
                        if let freeSpace = att[NSFileSystemFreeSize] as? Int {
                            if let totalSpace = att[NSFileSystemSize] as? Int {
                                let usedSpace =  totalSpace - freeSpace
                                guard let appContainerPath = Paths.appContainerPath() else {
                                    return
                                }
                                // TODO: maybe that should be async
                                let appContainerFolderFileSizeInBytes = Utils.getDirectorySizeInByte(appContainerPath)
                                // TODO: maybe that should be async
                                let dockerLogsSizeInBytes = Utils.getDirectorySizeInByte("\(appContainerPath)/com.docker.driver.amd64-linux/log")
                                
                                let othersSize = usedSpace - appContainerFolderFileSizeInBytes
                                let dockerWidth = (CGFloat)(appContainerFolderFileSizeInBytes) / (CGFloat)(totalSpace) * totalWidth
                                let othersWidth = (CGFloat)(othersSize) / (CGFloat)(totalSpace) * totalWidth
                                
                                diskUsageBar?.setWidth(dockerWidth, forSegment: 0)
                                diskUsageBar?.setWidth(othersWidth, forSegment: 1)
                                diskUsageBar?.setWidth(totalWidth - dockerWidth - othersWidth, forSegment: 2)
                                
                                // TODO: should be dynamic
                                let volumeName = "Macintosh HD"
                                
                                hdlabel?.stringValue = "settingsDiskUsageHDLabel".localize(volumeName, (Float)(freeSpace) / 1000000000.0, (Float)(totalSpace) / 1000000000.0)
                                dockerlabel?.stringValue = "settingsDiskUsageDockerLabel".localize()
                                dockerUsageLabel?.stringValue = "settingsDiskUsageDockerUsageLabel".localize((CGFloat)(appContainerFolderFileSizeInBytes) / 1000000000.0, (CGFloat)(dockerLogsSizeInBytes) / 1000000000.0)
                                otherslabel?.stringValue = "settingsDiskUsageOthersLabel".localize()
                                othersUsagelabel?.stringValue = "settingsDiskUsageOthersUsageLabel".localize((CGFloat)(othersSize) / 1000000000.0)
                            }
                        }
                    }
                }
            }
        }
    }
    
    func shouldApply() -> Bool {
        return false
    }
    
    func apply() {
        // nothing to apply
    }
}

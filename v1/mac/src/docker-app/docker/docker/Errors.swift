//
//  Errors.swift
//  Docker
//
//  Created by Michel Courtine on 12/5/15.
//  Copyright © 2015 Docker Inc. All rights reserved.
//

import Foundation

enum InstallerError: ErrorType {
    case FailedWithMessage(message: String, details: String)
    case FailedToRegisterDockerMenuAsAStartupItem
    case FailedToInstallLaunchAgent
}
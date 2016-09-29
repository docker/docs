//
//  ChildProcesses.swift
//  docker
//
//  Created by Gaetan de Villele on 2/23/16.
//  Copyright © 2016 docker. All rights reserved.
//

import Cocoa


class ChildProcesses {
    
    static func start() {
        TaskManager.Start()
        launchDockerTasks()
    }
    
    static func stop() {
        TaskManager.Stop()
    }
    
    // launch the necessary tasks to start the vm
    static func launchDockerTasks() {
        let backendTask = TaskDescriptor(name: "com.docker.osx.hyperkit.linux", commandPath: "MacOS/com.docker.osx.hyperkit.linux", environment: [:], arguments: [], watchdogPipe:true)
        TaskManager.StartTask(backendTask)
    }
    
    static func taskRelaunched(task: TaskInstance) {
        
    }
}

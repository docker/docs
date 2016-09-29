import Cocoa
import Foundation

// Define the behavior when a task exit
enum KeepAliveTaskStatus {
    case Disabled               // Task will not be restarted
    case RestartOnFailure       // Task will be restarted if exit code != 0 or if has been interrupted
    case RestartOnSuccess       // Task will be restarted after a clean exit (normal termination and exit code 0)
    case AlwaysRestart          // Task will always be restarted
}

// Describe a new task to launch with TaskManager
class TaskDescriptor {
    var Name: String                             // UI/Log name
    var CommandPath: String                      // relative to bundle/Contents if no / at the beginning
    var Environment: [String : String]?         // ex: ["DOCKER_HOST" : "192.168.99.10"], will be add to docker app env variables
    var Arguments: [String]?                    // command line args
    var WatchdogPipe: Bool                      // true if using the pipe-based watchdog protocol
    
    // specify the behavior to have when the task/process exit
    var KeepAlive: KeepAliveTaskStatus = KeepAliveTaskStatus.AlwaysRestart
    
    init(name: String,
        commandPath: String,
        environment: [String : String]?,
        arguments: [String]?,
        watchdogPipe: Bool ) {
            Name = name
            CommandPath = commandPath
            Environment = environment
            Arguments = arguments
            WatchdogPipe = watchdogPipe
    }
}

// TaskManager handles tasks for you and record logs from stdout/stderr in Logger.Logs
class TaskManager {
    
    typealias TaskInstanceSet = Set<TaskInstance>
    
    static private let lock = NSLock()
    static private var started = false
    static private var tasksRunning = TaskInstanceSet()
    static private var disableKeepAlive = false
    
    
    // Start the task manager
    static func Start() {
        lock.lock()
        defer { lock.unlock() }
        
        if started {
            Logger.log(level: Logger.Level.Warning, content: "TaskManager is already started - ignoring new start")
            return
        }
        
        started = true
        disableKeepAlive = false
        
        ExceptionHandler.setup()
    }
    
    // Stop the task manager
    static func Stop() {
        if !started {
            Logger.log(level: Logger.Level.Warning, content: "TaskManager is not started - ignoring stop")
            return
        }
        
        disableKeepAlive = true
        
        // We must wait for drivers (vm) to terminate before other tasks
        for task in tasksRunning {
            if task.Descriptor.CommandPath.containsString("com.docker.driver.") {
                task.WaitForTerminationWithTimeOut()
                lock.lock()
                tasksRunning.remove(task)
                lock.unlock()
            }
        }
        
        lock.lock(); defer { lock.unlock() }
        
        for task in tasksRunning {
            task.WaitForTerminationWithTimeOut()
        }
        tasksRunning.removeAll()
        
        started = false
    }
    
    // Start a new task
    static func StartTask(task: TaskDescriptor) -> TaskInstance {
        let newTask = TaskInstance(descriptor: task)
        newTask.Start()
        lock.lock()
        defer { lock.unlock() }
        tasksRunning.insert(newTask)
        return newTask
    }
    
    // Gently stop a new task (sending SIGTERM first then SIGKILL if the tasks has not been stopped within 10s)
    static func StopTaskInstance(task: TaskInstance) {
        task.Terminate()
        lock.lock()
        defer { lock.unlock() }
        tasksRunning.remove(task)
    }
}

// Internal class, representing a task instance
class TaskInstance: NSObject {
    // Shared descriptor
    let Descriptor: TaskDescriptor
    // The Cocoa task itself
    let Task = NSTask()
    // stdin redirection handler
    let Input = NSPipe()
    // stderr redirection handler
    let Error = NSPipe()
    // stdout redirection handler
    let Output = NSPipe()
    // Enable/Disable log recording
    var RecordLog = true
    // Callback used when the task is restarted throught KeepAlive property (see TaskDescriptor)
    var OnRestart: ((TaskInstance)->Void)?
    
    private var recordingLog = true
    
    //to handle max retries in a given time
    private var lastRetryTime: NSDate?
    private var retryCount = 0
    private static let retryLimitLock = NSLock()
    private static var retryLimitReached = false
    private var observer: NSObjectProtocol? = nil
    
    private init(descriptor: TaskDescriptor) {
        Descriptor = descriptor
        super.init()
    }
    
    // Start the task and setup observers
    private func Start() {
        recordingLog = true
        if Descriptor.CommandPath.hasPrefix("/") {
            Task.launchPath = Descriptor.CommandPath
        } else {
            if let executablePath = NSBundle.mainBundle().executablePath,
                path = NSURL.init(fileURLWithPath: executablePath).URLByDeletingLastPathComponent?.URLByDeletingLastPathComponent?.path {
                let fullPath = path + "/" + Descriptor.CommandPath
                Task.launchPath = fullPath
            } else {
                Logger.log(level: Logger.Level.Fatal, content: "\(#function): could not get launch path from bundle")
            }
        }

        // Task.environment defaults to the current process's environment
        if let env = Descriptor.Environment {
            let currentEnv = NSProcessInfo.processInfo().environment
            if Task.environment == nil { Task.environment = currentEnv }
            if let taskEnv = Task.environment {
                Task.environment = env.reduce(taskEnv) { (env, kv) in
                    var envVar = env
                    envVar[kv.0] = kv.1
                    return envVar
                }
            }
        }
        
        guard let appContainerPath = Paths.appContainerPath() else {
            return
        }

        Task.currentDirectoryPath = appContainerPath
        
        var args = [String]()
        if let arguments: [String] = Descriptor.Arguments {
            args = arguments
        }
        if Descriptor.WatchdogPipe {
            args.append("-watchdog")
            args.append("fd:0") // stdin
        }
        // we cannot assign nil to Task.arguments, so we assign an empty array
        // if Descriptor.Arguments is nil
        Task.arguments = args
        Task.standardInput = Input
        Task.standardError = nil
        Task.standardOutput = nil
        
        // observer to catch task termination
        self.observer = NSNotificationCenter.defaultCenter().addObserverForName(NSTaskDidTerminateNotification, object: Task, queue: nil, usingBlock: {_ in
            let terminationReason = self.Task.terminationReason
            let terminationStatus = self.Task.terminationStatus
            var needRestart = false

            Logger.log(level: Logger.Level.Notice, content: "Notified of termination of " + self.Descriptor.Name)
            //determine if we need to restart task
            if self.Descriptor.KeepAlive == KeepAliveTaskStatus.RestartOnFailure && (terminationReason == NSTaskTerminationReason.UncaughtSignal || terminationStatus != 0) {
                needRestart = true
            }
            if self.Descriptor.KeepAlive == KeepAliveTaskStatus.RestartOnSuccess && terminationReason == NSTaskTerminationReason.Exit && terminationStatus == 0 {
                needRestart = true
            }
            if self.Descriptor.KeepAlive == KeepAliveTaskStatus.AlwaysRestart {
                needRestart = true
            }
            if TaskManager.disableKeepAlive == true {
                needRestart = false
            }
            
            //remove instance from running task
            TaskManager.StopTaskInstance(self)
            
            // trying to restart task 5 times in 5 sec, if failed display alert box and exit
            if needRestart {
                if self.lastRetryTime == nil {
                    self.retryCount = 1
                } else if self.lastRetryTime?.timeIntervalSinceNow < -5.0 { // Reset retry count if task has run for at least 5 s
                    self.retryCount = 1
                } else if self.retryCount < 5 { // Retrying 5 times then display alert and exit
                    self.retryCount += 1
                } else {
                    
                    //add mutex to avoid multiple alert boxes
                    TaskInstance.retryLimitLock.lock()
                    if TaskInstance.retryLimitReached {
                        TaskInstance.retryLimitLock.unlock()
                        return
                    }
                    TaskInstance.retryLimitReached = true
                    TaskInstance.retryLimitLock.unlock()
                    
                    Logger.log(level: Logger.Level.Fatal, content: self.Descriptor.Name + " failed to start\nExit code " + String(terminationStatus))
                    
                    NSApp.terminate(nil)
                    return
                }
                // Wait 1 sec for the next retry
                usleep(1000000)
                Logger.log(level: Logger.Level.Notice, content: "Attempting to restart " + self.Descriptor.Name + " #" + String(self.retryCount))
                let newTask = TaskManager.StartTask(self.Descriptor)
                newTask.OnRestart = self.OnRestart
                newTask.lastRetryTime = NSDate()
                newTask.retryCount = self.retryCount
                if let onRestart = self.OnRestart {
                    onRestart(newTask)
                }
            } else {
                // doesn't need restart
                if let obs = self.observer {
                    NSNotificationCenter.defaultCenter().removeObserver(obs)
                }
            }
        })
        
        Task.launch()
        if Descriptor.WatchdogPipe {
            let msg = "This is not a stable interface. Do not try to use it, it will change."
            if let data = msg.dataUsingEncoding(NSUTF8StringEncoding) {
                Input.fileHandleForWriting.writeData(data)
            }
        }
        
        Logger.log(level: Logger.Level.Notice, content: String(format:"%@ launched with PID %d", Descriptor.Name, Task.processIdentifier))
    }
    
    // Stop the task in background sending SIGTERM,
    // waiting to terminate for 40s and sending SIGKILL if still running
    private func Terminate() {
        if Task.running {
            let qualityOfServiceClass = QOS_CLASS_DEFAULT
            let backgroundQueue = dispatch_get_global_queue(qualityOfServiceClass, 0)
            dispatch_async(backgroundQueue, {
                Logger.log(level: Logger.Level.Notice, content:String(format:"Terminating %@", self.Descriptor.Name))
                //send SIGTERM to gently stop the process
                self.Task.terminate()
                usleep(100)
                var loop = 4000
                //wait 40 sec, breaking as soon as the task is not running
                while (self.Task.running && loop > 0) {
                    usleep(10000)
                    loop -= 1
                }
                //if still running send SIGKILL
                if self.Task.running {
                    self.Task.interrupt()
                    Logger.log(level: Logger.Level.Notice, content:String(format:"%@ termination timed out and has been killed", self.Descriptor.Name))
                    usleep(10000)
                }
                //report exit status
                if !self.Task.running {
                    Logger.log(level: Logger.Level.Notice, content: String(format:"%@ terminated with exit code %d", self.Descriptor.Name, self.Task.terminationStatus))
                }
            })
        }
        else {
            //report exit status
            Logger.log(level: Logger.Level.Notice, content: String(format:"%@ terminated with exit code %d", Descriptor.Name, Task.terminationStatus))
        }
    }
    
    // Directly stop a task and wait until it's shutdowned
    // send SIGTERM, waiting to terminate for 10s and send SIGKILL if still running
    private func WaitForTerminationWithTimeOut() {
        if Task.running {
            Logger.log(level: Logger.Level.Notice, content: String(format:"Terminating %@ and waiting until exit", self.Descriptor.Name))
            //send SIGTERM to gently stop the process
            Task.terminate()
            usleep(100)
            var loop = 1000
            //wait 10 sec, breaking as soon as the task is not running
            while (Task.running && loop > 0) {
                usleep(10000)
                loop -= 1
            }
            //if still running send SIGKILL
            if self.Task.running {
                Task.interrupt()
                Logger.log(level: Logger.Level.Notice, content: String(format:"%@ termination timed out and has been killed", Descriptor.Name))
                usleep(10000)
            }
            //report exit status
            if !Task.running {
                Logger.log(level: Logger.Level.Notice, content: String(format:"%@ terminated with exit code %d", Descriptor.Name, Task.terminationStatus))
            }
        }
        else {
            //report exit status
            Logger.log(level: Logger.Level.Notice, content: String(format:"%@ terminated with exit code %d", Descriptor.Name, Task.terminationStatus))
        }
    }
    
    private func Interrupt() {
        if Task.running {
            Logger.log(level: Logger.Level.Notice, content: String(format:"Interrupting %@", Descriptor.Name, Task.terminationStatus))
            self.Task.interrupt()
        }
    }
}

// QA feature: send SIGUSR2 to test crash reporting.
let crashTest: @convention(c) (Int32) -> Void = { Int32 -> Void in
    let a: String? = nil
    _ = a! // swiftlint:disable:this force_unwrapping
}

// ExceptionHandler handles exceptions for TaskManager
private class ExceptionHandler {
    
    // handleUncaughtException handles uncaught exceptions
    static let handleUncaughtException: (@convention(c) (NSException) -> Void)? = {
        (e: NSException) -> Void in
        
        Logger.log("Exception", level: Logger.Level.Error, content: "name:\(e.name)")
        if let reason = e.reason {
            Logger.log("Exception", level: Logger.Level.Error, content: "reason:\(reason)")
        }
        if let userInfo = e.userInfo {
            Logger.log("Exception", level: Logger.Level.Error, content: "user info:\(userInfo)")
        }
        Logger.log("Exception", level: Logger.Level.Error, content: "return adress:\(e.callStackReturnAddresses)")
        
        let callStack = e.callStackSymbols
        if !callStack.isEmpty {
            for call in callStack {
                Logger.log("Exception", level: Logger.Level.Error, content: "function:\(call)")
            }
        }
        NSApp.terminate(nil)
    }
    
    // setup sets handlers to use for uncaught exceptions and signals
    static func setup() {
        NSSetUncaughtExceptionHandler(handleUncaughtException)

        signal(SIGUSR2) { crashTest($0) }
    }
}



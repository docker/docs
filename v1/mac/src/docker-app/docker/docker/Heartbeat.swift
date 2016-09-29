import Foundation

let heartbeat = Heartbeat()

// It inherits from NSObject for methodSignatureForSelector: (required when scheduling call of a selector)
public class Heartbeat: NSObject {
    private let heartbeatDelay: NSTimeInterval = 60 * 60 // delay between two heartbeats
    private var heartbeatTimer: NSTimer
    
    public override init() {
        self.heartbeatTimer = NSTimer()
    }
    
    public func start() {
        analytics.track(AnalyticEvent.Heartbeat)
        self.heartbeatTimer = NSTimer.scheduledTimerWithTimeInterval(heartbeatDelay, target: self, selector: #selector(sendHeartbeat(_:)), userInfo: nil, repeats: true)
    }
    
    func sendHeartbeat(timer: NSTimer) {
        analytics.track(AnalyticEvent.Heartbeat)
    }
    
    public func stop() {
        self.heartbeatTimer.invalidate()
    }
    
}

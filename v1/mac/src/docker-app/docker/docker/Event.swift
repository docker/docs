//
//  Event.swift
//  docker-installer
//
//  Created by Michel Courtine on 12/14/15.
//  Copyright © 2015 David Scott. All rights reserved.
//

import Foundation

public class Event<T> {
    
    public typealias EventHandler = T -> ()
    
    private var eventHandlers = [Invocable]()
    
    public func raise(data: T) {
        for handler in self.eventHandlers {
            handler.invoke(data)
        }
    }
    
    public func addHandler<U: AnyObject>(target: U,
        handler: (U) -> EventHandler) -> Disposable {
            let wrapper = EventHandlerWrapper(target: target,
                handler: handler, event: self)
            eventHandlers.append(wrapper)
            return wrapper
    }
}

private protocol Invocable: class {
    func invoke(data: Any)
}

private class EventHandlerWrapper<T: AnyObject, U>
: Invocable, Disposable {
    weak var target: T?
    let handler: T -> U -> ()
    let event: Event<U>
    
    init(target: T?, handler: T -> U -> (), event: Event<U>) {
        self.target = target
        self.handler = handler
        self.event = event
    }
    
    func invoke(data: Any) -> () {
        if let t = target,
               u = data as? U {
            handler(t)(u)
        } else {
            Logger.log(level: Logger.Level.Fatal, content: "Unable to handle event (t was nil or unable to cast Any to U)")
        }
    }
    
    func dispose() {
        event.eventHandlers =
            event.eventHandlers.filter { $0 !== self }
    }
}

public protocol Disposable {
    func dispose()
}

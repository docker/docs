import Foundation
import Cocoa

// WizardMessageGeneric is a WizardMessage implementation.
// It can optionally display an icon on the top left corner
// And a different components from top to bottom on the
// right side (text, loading bar, input label...etc)
class WizardMessageDiagnose: WizardMessage {
    
    @IBOutlet var progressIndicator: NSProgressIndicator?
    @IBOutlet var progressLabel: NSTextField?
    
    @IBOutlet var getDiagnosticIDProgressIndicator: NSProgressIndicator?
    @IBOutlet var diagnosticIDTitleLabel: NSTextField?
    @IBOutlet var diagnosticIDLabel: NSTextField?
    
    @IBOutlet var openDocumentationBtn: NSButton?
    @IBOutlet var openForumsBtn: NSButton?
    
    @IBOutlet var view1: NSView?
    @IBOutlet var view1HeightLayoutConstraint: NSLayoutConstraint?
    
    @IBOutlet var view2: NSView?
    @IBOutlet var view2HeightLayoutConstraint: NSLayoutConstraint?
    
    @IBOutlet var view3: NSView?
    @IBOutlet var view3HeightLayoutConstraint: NSLayoutConstraint?
    private var textView3: NSTextView?
    
    @IBOutlet var diagnosticTextView: TextView?
    
    private var sheetViewController: DiagnoseSheetViewController! = nil
    private var diagnosticInfo: String = ""
    
    init?(message: String? = nil) {
        super.init(nibName: "WizardMessageDiagnose", bundle: NSBundle.mainBundle())
    }
    
    required init?(coder: NSCoder) {
        super.init(coder: coder)
    }
    
    override func viewDidLoad() {
        super.viewDidLoad()
        
        progressLabel?.stringValue = "DiagnoseAndFeedbackProgressLabel".localize()
        diagnosticIDTitleLabel?.stringValue = "DiagnoseAndFeedbackDiagnosticIdTitle".localize()
        openDocumentationBtn?.title = "DiagnoseAndFeedbackOpenDocsBtn".localize()
        openForumsBtn?.title = "DiagnoseAndFeedbackOpenIssuesBtn".localize()
        
        let maxHeight: CGFloat = 10000
        let attrs = [NSFontAttributeName:NSFont.systemFontOfSize(NSFont.systemFontSize())]
        
        if let view = view1 {
            let attributedString = NSMutableAttributedString(string: "DiagnoseAndFeedbackText1".localize(), attributes:attrs)
            attributedString.format(NSFont.systemFontSize())
            let textView = NSTextView()
            textView.selectable = true
            textView.editable = false
            textView.drawsBackground = false
            textView.alignment = .Justified
            textView.textStorage?.setAttributedString(attributedString)
            if let textStorage = textView.textStorage {
                let boundingRect = textStorage.boundingRectWithSize(NSSize(width: view.frame.size.width, height: maxHeight), options: [.UsesLineFragmentOrigin, .UsesFontLeading])
                textView.setFrameSize(NSSize(width: view.frame.size.width, height: boundingRect.size.height))
                view1HeightLayoutConstraint?.constant = boundingRect.size.height
            }
            view.addSubview(textView)
        }
        
        if let view = view2 {
            let attributedString = NSMutableAttributedString(string: "DiagnoseAndFeedbackText2".localize(), attributes:attrs)
            attributedString.format(NSFont.systemFontSize())
            let textView = NSTextView()
            textView.selectable = true
            textView.editable = false
            textView.drawsBackground = false
            textView.alignment = .Justified
            textView.textStorage?.setAttributedString(attributedString)
            if let textStorage = textView.textStorage {
                let boundingRect = textStorage.boundingRectWithSize(NSSize(width: view.frame.size.width, height: maxHeight), options: [.UsesLineFragmentOrigin, .UsesFontLeading])
                textView.setFrameSize(NSSize(width: view.frame.size.width, height: boundingRect.size.height))
                view2HeightLayoutConstraint?.constant = boundingRect.size.height
            }
            view.addSubview(textView)
        }
        
        if let view = view3 {
            if textView3 == nil {
                textView3 = NSTextView()
                textView3?.selectable = true
                textView3?.editable = false
                textView3?.drawsBackground = false
                textView3?.alignment = .Center
            }
            guard let textView = textView3 else {
                return
            }
            updateUploadMessage("DiagnoseAndFeedbackUpload".localize())
            view.addSubview(textView)
        }
        
        diagnosticIDLabel?.selectable = true
        diagnosticIDLabel?.editable = false
        
        // Hide diagnostic ID label by default. We show it only if an upload of
        // the diagnostic info is requested.
        diagnosticIDTitleLabel?.hidden = true
        diagnosticIDLabel?.hidden = true
        getDiagnosticIDProgressIndicator?.hidden = true
        openForumsBtn?.hidden = true
        view3?.hidden = true
    }
    
    //
    override func viewDidAppear() {
        super.viewDidAppear()
            
        // Display the sheet asking for the kind of diagnose wanted, right away.
        self.sheetViewController = DiagnoseSheetViewController.init(nibName: "DiagnoseSheetViewController", bundle: NSBundle.mainBundle())
        self.sheetViewController.diagnoseAndUploadCallback = {
            self.startDiagnostic(upload: true)
        }
        self.sheetViewController.diagnoseOnlyCallback = {
            self.startDiagnostic(upload: false)
        }
        self.presentViewControllerAsSheet(sheetViewController)
    }
    
    //
    private func startDiagnostic(upload upload: Bool) {
        
        progressLabel?.hidden = false
        progressIndicator?.hidden = false
        progressIndicator?.startAnimation(self)
        openForumsBtn?.enabled = false
        openForumsBtn?.highlight(false)
        diagnosticTextView?.text = ""
        
        // if no upload is planned, we do not display the diagnostic ID
        if upload == true {
            diagnosticIDTitleLabel?.hidden = false
            openForumsBtn?.hidden = false
            getDiagnosticIDProgressIndicator?.hidden = false
            getDiagnosticIDProgressIndicator?.hidden = false
            getDiagnosticIDProgressIndicator?.startAnimation(self)
        } else {
            // update text to tell the user to reopen and select "diagnose and upload"
            // if they need help with an issue.
            if let view = view2 {
                if let textview = view.subviews.first as? NSTextView {
                    let maxHeight: CGFloat = 10000
                    let attrs = [NSFontAttributeName:NSFont.systemFontOfSize(NSFont.systemFontSize())]
                    let attributedString = NSMutableAttributedString(string: "DiagnoseAndFeedbackText2NoUpload".localize(), attributes:attrs)
                    attributedString.format(NSFont.systemFontSize())
                    textview.textStorage?.setAttributedString(attributedString)
                }
            }
        }
        
        // diagnose completion callback
        func completion(error error: String?, resp: [String: AnyObject]?) -> Void {
            
            progressIndicator?.hidden = true
            progressIndicator?.stopAnimation(self)
            progressLabel?.hidden = true
            
            
            if let errorMessage = error {
                Logger.log(level: ASLLogger.Level.Error, content: errorMessage)
            }
            
            
            guard var response = resp else {
                diagnosticTextView?.text = "#ff3131#failure: \(error)##\n"
                Logger.log(level: ASLLogger.Level.Error, content: "diagnostic: empty response")
                return
            }
            
            var diagnostic = ""
            diagnosticInfo = ""
            
            if let app = response["app"] as? String {
                diagnostic += "#ffffff#Docker for Mac: " + app + "##\n"
            }
            
            if let os = response["os"] as? String {
                diagnostic += "#ffffff#" + os + "##\n"
            }
            
            if let logs = response["logs"] as? String {
                diagnostic += "#ffffff#logs: " + logs + "##\n"
            }
            
            if let failure = response["failure"] as? String {
                diagnostic += "#ffffff#failure: " + failure + "##\n"
            }
            
            if let tests = response["tests"] as? [String: [String: AnyObject]] {
                for (key, value) in tests {
                    
                    var success = false
                    if let s = value["success"] as? Bool {
                        success = s
                    }
                    
                    var messages: [String] = []
                    if let m = value["message"] as? [String] {
                        messages = m
                    }
                    
                    var messageColor = "#e7ffca#"
                    
                    if success {
                        diagnostic += "#a0ff31#[OK]     \(key)##\n"
                        
                        diagnosticInfo += "[OK]     \(key)\n"
                    } else {
                        diagnostic += "#ff3131#[ERROR]  \(key)##\n"
                        messageColor = "#ffb2b2#"
                        
                        diagnosticInfo += "[ERROR]  \(key)\n"
                    }
                    
                    for message in messages {
                        diagnostic += "\(messageColor)         \(message)##\n"
                        
                        diagnosticInfo += "         \(message)\n"
                    }
                }
            }
            if let errorMessage = error {
                diagnostic += "#ff3131#Failure: \(errorMessage)##\n"
            }
            
            
            diagnosticTextView?.text = diagnostic
            
            // If upload was requested, we display the diagnostic ID and the
            // "Open Forums" button.
            if upload {
                if let id = response["id"] as? String {
                    if let diagnosticIDLabel = diagnosticIDLabel {
                        getDiagnosticIDProgressIndicator?.hidden = true
                        getDiagnosticIDProgressIndicator?.stopAnimation(self)
                        diagnosticIDLabel.hidden = false
                        diagnosticIDLabel.stringValue = id
                        openForumsBtn?.enabled = true
                        copyToClipboard(id)
                        // storing diagnostic id in preferences to go faster
                        // next time Diagnose & Feedback window is opened.
                        // If diagnostic id can be found within preferences
                        // we enable "Open Forums" button right away
                        Preferences.sharedInstance.setDiagnosticId(id)
                    }
                }
            }
        }
        
        // launch diagnose process in background
        let backgroundQueue = dispatch_get_global_queue(QOS_CLASS_BACKGROUND, 0)
        dispatch_async(backgroundQueue, {
            let bundlePath: String = NSBundle.mainBundle().bundlePath
            let executablePath: String = NSString.pathWithComponents([bundlePath, "Contents", "Resources", "bin", "docker-diagnose"])
            let task = NSTask()
            let out = NSPipe()
            let err = NSPipe()
            task.launchPath = executablePath
            task.arguments = ["--json", "-n"]
            if upload {
                task.arguments?.append("-u") // add the "upload" flag
            }
            task.standardOutput = out
            task.standardError = err
            
            var env = [String: String]()
            // NOTE(aduermael) this won't point to the home if the app is sandboxed
            env["HOME"] = NSHomeDirectory()
            env["PATH"] = "/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin"
            task.environment = env
            task.launch()
            task.waitUntilExit()
            var errorMessage: String?
            let terminationStatus = task.terminationStatus
            if terminationStatus != 0 {
                errorMessage = "Could not upload diagnostic data to remote server (docker-diagnose exit code is \(terminationStatus))"
            }
            
            // read response data
            let responseData = out.fileHandleForReading.readDataToEndOfFile()
            out.fileHandleForReading.closeFile()
            // check that the response is a UTF-8 string
            guard let _: String = String(data: responseData, encoding: NSUTF8StringEncoding) else {
                dispatch_async(dispatch_get_main_queue(), { () -> Void in
                    completion(error: "diagnostic: response is not UTF-8", resp: nil)
                })
                return
            }
            // parse response and check whether it is a success or failure
            let responseObject: [String: AnyObject]!
            do {
                responseObject = try NSJSONSerialization.JSONObjectWithData(responseData, options: NSJSONReadingOptions(rawValue: 0)) as? [String: AnyObject]
            } catch {
                dispatch_async(dispatch_get_main_queue(), { () -> Void in
                    completion(error: "diagnostic: response is not valid JSON", resp: nil)
                })
                return
            }
            dispatch_async(dispatch_get_main_queue(), { () -> Void in
                completion(error: errorMessage, resp: responseObject)
            })
        })
    }
    
    private func copyToClipboard(text: String) {
        let pasteBoard = NSPasteboard.generalPasteboard()
        pasteBoard.clearContents()
        pasteBoard.setString(text, forType: NSStringPboardType)
    }
    
    //
    private func updateUploadMessage(message: String) {
        guard let view = view3 else {
            return
        }
        guard let textView = textView3 else {
            return
        }
        let maxHeight: CGFloat = 10000
        let attrs = [NSFontAttributeName:NSFont.systemFontOfSize(NSFont.systemFontSize())]
        let attributedString = NSMutableAttributedString(string: message, attributes:attrs)
        attributedString.format(NSFont.systemFontSize()-2)
        textView.textStorage?.setAttributedString(attributedString)
        if let textStorage = textView.textStorage {
            let boundingRect = textStorage.boundingRectWithSize(NSSize(width: view.frame.size.width, height: maxHeight), options: [.UsesLineFragmentOrigin, .UsesFontLeading])
            textView.setFrameSize(NSSize(width: view.frame.size.width, height: boundingRect.size.height))
            view3HeightLayoutConstraint?.constant = boundingRect.size.height
        }
    }
    
    @IBAction func openForums(sender: NSButton) {
        var url: String = "DiagnoseAndFeedbackIssueURL".localize()
        
        if let escapedURL = url.stringByAddingPercentEncodingWithAllowedCharacters(NSCharacterSet.URLQueryAllowedCharacterSet()) {
            url = escapedURL
        } else {
            return
        }
        if let url = NSURL(string: url) {
            NSWorkspace.sharedWorkspace().openURL(url)
        } else {
            Logger.log(level: ASLLogger.Level.Error, content: "\(#function): malformed URL")
        }
    }
}

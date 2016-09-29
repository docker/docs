import Foundation

public class Version{
    var Name: String = ""
    var Major: Int = 0
    var Minor: Int = 0
    var Patch: Int = 0
    
    public init() {
        fromBundle()
    }
    
    public init(name: String, major: Int, minor: Int, patch: Int) {
        self.Name = name
        self.Major = major
        self.Minor = minor
        self.Patch = patch
    }
    
    public func fromBundle() {
        guard let bundleVersion: String = NSBundle.mainBundle().objectForInfoDictionaryKey("CFBundleShortVersionString") as? String else {
        return
        }
    
        self.Name = bundleVersion
    
        // parse the first part of the version (<int>.<int>.<int>)
        let bundleVersionSplit: [String] = bundleVersion.componentsSeparatedByString("-")
        if bundleVersionSplit.count >= 1 {
            let versionSplit: [String] = bundleVersionSplit[0].componentsSeparatedByString(".")
            if versionSplit.count == 3 {
                self.Major = toInt(versionSplit[0])
                self.Minor = toInt(versionSplit[1])
                self.Patch = toInt(versionSplit[2])
            }
        }
    }
    
    private func toInt(numberAsString: String) -> Int {
        if let number = Int(numberAsString) {
            return number
        }
        return 0
    }
}

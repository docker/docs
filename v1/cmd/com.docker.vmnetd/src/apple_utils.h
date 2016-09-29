#ifndef apple_utils_h_
#define apple_utils_h_

// returns path for file in container folder
// relativePath parameter shouldn't start with a '/'
// https://developer.apple.com/library/mac/documentation/Security/Conceptual/AppSandboxDesignGuide/AboutAppSandbox/AboutAppSandbox.html
extern const char* apple_utils_get_file_path_in_container_folder(uid_t uid, const char* relativePath);

#endif

#include <SystemConfiguration/SystemConfiguration.h>
#include <pwd.h>
#include "apple_utils.h"
#include "utils.h"


const char* apple_utils_get_file_path_in_container_folder(uid_t uid, const char* relativePath) {

  static char * filePathInContainerFolder = NULL;
  free(filePathInContainerFolder);

  /* Query the current logged in user. If the query fails, fall back to the
     uid provided by the caller. */
  SCDynamicStoreRef store = NULL;
  CFStringRef name = NULL;
  uid_t tmp;
  store = SCDynamicStoreCreate(NULL, CFSTR("GetConsoleUser"), NULL, NULL);
  if (store == NULL) {
    aslLog(ASL_LEVEL_CRIT, "SCDynamicsStoreCreate failed: %s", SCErrorString(SCError()));
    goto out;
  }
  name = SCDynamicStoreCopyConsoleUser(store, &tmp, NULL);
  if (name == NULL) {
    aslLog(ASL_LEVEL_CRIT, "SCDynamicStoreCopyConsumeUser failed: %s", SCErrorString(SCError()));
    goto out;
  }
  uid = tmp;

out:
  if (store != NULL) CFRelease(store);
  if (name != NULL) CFRelease(name);

  struct passwd *pw = getpwuid(uid);
  const char *homedir = pw->pw_dir;
  size_t sizeOfHomeDir = strlen(homedir);

  const char* containerFolder = "/Library/Containers/com.docker.docker/Data/";
  size_t sizeOfcontainerFolder = strlen(containerFolder);

  size_t sizeOfRelativePath = strlen(relativePath);

  filePathInContainerFolder = (char*)malloc(sizeOfHomeDir + sizeOfcontainerFolder + sizeOfRelativePath + 1);

  strcpy(filePathInContainerFolder, homedir);
  strcpy(filePathInContainerFolder + sizeOfHomeDir, containerFolder);
  strcpy(filePathInContainerFolder + sizeOfHomeDir + sizeOfcontainerFolder, relativePath);

  aslLog(ASL_LEVEL_NOTICE, "User %d has container folder %s", uid, filePathInContainerFolder);
  return filePathInContainerFolder;
}

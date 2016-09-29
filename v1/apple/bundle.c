#include <CoreFoundation/CoreFoundation.h>

/* Find kernel and ramdisk paths within the application bundle */

static const char *find_in_bundle(CFStringRef resource, CFStringRef ext){
  char *result = NULL;

  CFBundleRef mainBundle = CFBundleGetMainBundle();
  CFURLRef thing = CFBundleCopyResourceURL(mainBundle, resource, ext, NULL);
  if (!thing) {
    char buffer[2048];
    CFStringGetCString(resource, buffer, 2048, kCFStringEncodingUTF8);
    fprintf(stderr, "Failed to find resource in bundle: %s\n", buffer);
    exit(1);
  }
  CFStringRef path = CFURLCopyFileSystemPath(thing, kCFURLPOSIXPathStyle);
  result = malloc(PATH_MAX);
  CFStringGetCString(path, result, PATH_MAX, kCFStringEncodingUTF8);
  CFRelease(path);
  CFRelease(thing);
  return result;
}

/* Return the kernel. Caller must call free() */
const char *find_kernel(){
  return find_in_bundle(CFSTR("moby/vmlinuz64"), CFSTR(""));
}

/* Return the ramdisk. Caller must call free() */
const char *find_ramdisk(){
  return find_in_bundle(CFSTR("moby/initrd"), CFSTR("img"));
}

/* Return the template. Caller must call free() */
const char *find_template(){
  return find_in_bundle(CFSTR("moby/data"), CFSTR("qcow2"));
}

/* Return the UEFI/. Caller must call free() */
const char *find_uefi(){
  return find_in_bundle(CFSTR("uefi/UEFI"), CFSTR("fd"));
}

#include <asl.h>
#include <CoreFoundation/CoreFoundation.h>
#include <SystemConfiguration/SystemConfiguration.h>

#include "asl_logger.h"
#include "util.h"

static CFStringRef copy_current_primary_service(void) {
  SCDynamicStoreRef store;
  CFPropertyListRef prop_list;
  CFStringRef result = NULL;

  store = SCDynamicStoreCreate(kCFAllocatorDefault, CFSTR("Docker"), NULL, NULL);

  if(store) {
    prop_list = SCDynamicStoreCopyValue(store, CFSTR("State:/Network/Global/IPv4"));
    if(prop_list) {
      result = CFDictionaryGetValue(prop_list, CFSTR("PrimaryService"));
      if (result) {
        CFRetain(result);
      }
      CFRelease(prop_list);
    }
    CFRelease(store);
  }
  return result;
}

static CFArrayRef copy_current_dns_servers(void) {
  SCDynamicStoreRef store;
  CFPropertyListRef prop_list;
  CFArrayRef result = NULL;

  store = SCDynamicStoreCreate(kCFAllocatorDefault, CFSTR("Docker"), NULL, NULL);
  if(store) {
    CFRetain(store);
    prop_list = SCDynamicStoreCopyValue(store, CFSTR("State:/Network/Global/DNS"));
    if(prop_list) {
      CFRetain(prop_list);
      result = CFDictionaryGetValue(prop_list, CFSTR("ServerAddresses"));
      if (result) {
        CFRetain(result);
      }
      CFRelease(prop_list);
    }
    CFRelease(store);
  }
  return result;
}

static CFArrayRef copy_current_dns_search_domains(void) {
  SCDynamicStoreRef store;
  CFPropertyListRef prop_list;
  CFArrayRef result = NULL;

  store = SCDynamicStoreCreate(kCFAllocatorDefault, CFSTR("Docker"), NULL, NULL);
  if(store) {
    CFRetain(store);
    prop_list = SCDynamicStoreCopyValue(store, CFSTR("State:/Network/Global/DNS"));
    if(prop_list) {
      CFRetain(prop_list);
      result = CFDictionaryGetValue(prop_list, CFSTR("SearchDomains"));
      if (result) {
        CFRetain(result);
      }
      CFRelease(prop_list);
    }
    CFRelease(store);
  }
  return result;
}

char *convert_cf_string(const CFStringRef s) {
  CFIndex length = CFStringGetLength(s);
  CFIndex max_size = CFStringGetMaximumSizeForEncoding(length, kCFStringEncodingUTF8) + 1;
  char *buffer = malloc(max_size);
  if (buffer == NULL) {
    apple_asl_logger_log(ASL_LEVEL_CRIT, "unable to allocate memory for new char*");
    abort();
  }
  CFStringGetCString(s, buffer, max_size, kCFStringEncodingUTF8);
  return buffer;
}

void get_dns_servers(char **result){
  CFArrayRef array = copy_current_dns_servers();
  if (array == NULL) {
    return;
  }
  CFRetain(array);
  CFIndex c = CFArrayGetCount(array);
  for (int i = 0; i<c; i++) {
    CFStringRef s = (CFStringRef)CFArrayGetValueAtIndex(array, i);
    CFRetain(s);
    result[i] = convert_cf_string(s);
    CFRelease(s);
  }
  CFRelease(array);
}

void get_dns_search_domains(char **result){
  CFArrayRef array = copy_current_dns_search_domains();
  if (array == NULL) {
    return;
  }
  CFRetain(array);
  CFIndex c = CFArrayGetCount(array);
  for (int i = 0; i<c; i++) {
    CFStringRef s = (CFStringRef)CFArrayGetValueAtIndex(array, i);
    CFRetain(s);
    result[i] = convert_cf_string(s);
    CFRelease(s);
  }
  CFRelease(array);
}

void get_proxy_servers(char **result){
  char *http_proxy = NULL;
  char *https_proxy = NULL;
  char *no_proxy = NULL;
  CFNumberRef http_proxy_enabled;
  CFNumberRef https_proxy_enabled;
  CFArrayRef excluded;

  CFDictionaryRef data = SCDynamicStoreCopyProxies(NULL);

  if (CFDictionaryContainsKey(data, kSCPropNetProxiesHTTPEnable)) {
    http_proxy_enabled = CFDictionaryGetValue(data, kSCPropNetProxiesHTTPEnable);
    if (http_proxy_enabled) {
      CFRetain(http_proxy_enabled);
      int enabled, port;
      char *url;
      if (CFNumberGetValue(http_proxy_enabled, kCFNumberIntType, &enabled)){
        if (enabled == 1) {
          CFNumberRef port_ref = CFDictionaryGetValue(data, kSCPropNetProxiesHTTPPort);
          if (port_ref) {
            CFRetain(port_ref);
            CFNumberGetValue(port_ref, kCFNumberIntType, &port);
          }

          CFStringRef url_ref = CFDictionaryGetValue(data, kSCPropNetProxiesHTTPProxy);
          if (url_ref) {
            CFRetain(url_ref);
            url = convert_cf_string(url_ref);
          }
          asprintf(&http_proxy, "%s:%d", url, port);
        }
      }
    }
    CFRelease(http_proxy_enabled);
  }

  if (CFDictionaryContainsKey(data, kSCPropNetProxiesHTTPSEnable)) {
    https_proxy_enabled = CFDictionaryGetValue(data, kSCPropNetProxiesHTTPSEnable);
    if (https_proxy_enabled) {
      CFRetain(https_proxy_enabled);
      int enabled, port;
      char *url;
      if (CFNumberGetValue(https_proxy_enabled, kCFNumberIntType, &enabled)){
        if (enabled == 1) {
          CFNumberRef port_ref = CFDictionaryGetValue(data, kSCPropNetProxiesHTTPSPort);
          if (port_ref) {
            CFRetain(port_ref);
            CFNumberGetValue(port_ref, kCFNumberIntType, &port);
          }

          CFStringRef url_ref = CFDictionaryGetValue(data, kSCPropNetProxiesHTTPSProxy);
          if (url_ref) {
            CFRetain(url_ref);
            url = convert_cf_string(url_ref);
          }
          asprintf(&https_proxy, "%s:%d", url, port);
        }
      }
    }
    CFRelease(https_proxy_enabled);
  }

  if (CFDictionaryContainsKey(data, kSCPropNetProxiesExceptionsList)) {
    excluded = CFDictionaryGetValue(data, kSCPropNetProxiesExceptionsList);
    if (excluded) {
      CFRetain(excluded);
      CFIndex c = CFArrayGetCount(excluded);
      char **parts = malloc(sizeof(char*)*c);
      if (parts == NULL) {
        apple_asl_logger_log(ASL_LEVEL_CRIT, "unable to allocate memory for new char**");
        abort();
      }
      for (int i=0; i<c; i++){
        CFStringRef s = (CFStringRef)CFArrayGetValueAtIndex(excluded, i);
        CFRetain(s);
        parts[i] = convert_cf_string(s);
        CFRelease(s);
      }
      no_proxy = join_strings(parts, ", ", (int)c);
    }
    CFRelease(excluded);
  }

  CFRelease(data);

  result[0] = http_proxy;
  result[1] = https_proxy;
  result[2] = no_proxy;
}

struct connection {
  int writer;
};

void on_config_change(SCDynamicStoreRef store, CFArrayRef changed_keys, void *info) {
  struct connection *c = (struct connection*)info;
  CFIndex k = CFArrayGetCount(changed_keys);
  for (int i=0; i<k;i++){
    CFStringRef p = (CFStringRef) CFArrayGetValueAtIndex(changed_keys,i);
    char *s = convert_cf_string(p);
    if (strcmp(s, "State:/Network/Global/DNS") == 0) {
      write(c->writer, "D", 1);
    } else if (strcmp(s,"State:/Network/Global/Proxies") == 0) {
      write(c->writer, "P", 1);
    }
  }
}

void listen_for_config_changes(int writer) {
  SCDynamicStoreRef store;

  /* Allocate but never free the fd */
  void *data = malloc(sizeof(struct connection));
  if (data == NULL) {
    apple_asl_logger_log(ASL_LEVEL_CRIT, "unable to allocate memory for new connection struct");
    abort();
  }

  struct connection *c = (struct connection*)data;
  c->writer = writer;

  SCDynamicStoreContext ctx = {0, data, NULL, NULL ,NULL};

  store = SCDynamicStoreCreate(kCFAllocatorDefault,CFSTR("Docker"), on_config_change, &ctx);

  CFMutableArrayRef keys = CFArrayCreateMutable(NULL,0,&kCFTypeArrayCallBacks);
  CFArrayAppendValue(keys, CFSTR("State:/Network/Global/DNS"));
  CFArrayAppendValue(keys, CFSTR("State:/Network/Global/Proxies"));

  SCDynamicStoreSetNotificationKeys(store, NULL, keys);
  CFRelease(keys);

  CFRunLoopSourceRef rls = SCDynamicStoreCreateRunLoopSource(kCFAllocatorDefault, store, 0);
  CFRunLoopAddSource(CFRunLoopGetCurrent(), rls, kCFRunLoopDefaultMode);
  CFRunLoopRun();
  CFRelease(rls);
}

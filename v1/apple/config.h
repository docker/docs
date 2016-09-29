// Return the DNS servers, user must call free();
extern void *get_dns_servers(char **array);

// Return the DNS search domains, user must call free();
extern void *get_dns_search_domains(char **array);

// Return the current Proxy servers, user must call free();
extern void *get_proxy_servers(char **array);

// Listen for config changes
extern void *listen_for_config_changes(int writer);

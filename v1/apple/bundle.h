
/* Return the kernel. Caller must call free() */
extern const char *find_kernel();

/* Return the ramdisk. Caller must call free() */
extern const char *find_ramdisk();

/* Return the template. Caller must call free() */
extern const char *find_template();

/* Return the uefi. Caller must call free() */
extern const char *find_uefi();

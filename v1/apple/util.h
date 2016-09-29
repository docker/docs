#ifndef util_h_
#define util_h_

// Create a new char**, user must call free_char_array 
extern char **new_char_array(int size); 

// Get char from char**
extern char *get_array_string(char **a, int n);

// Set char in char**
extern void set_array_string(char **a, char *s, int n);

// Free a char**
extern void free_char_array(char **a, int size);

// Join strings
extern char *join_strings(char **src, const char *sep, int count);
#endif

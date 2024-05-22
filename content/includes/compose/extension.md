Extensions can be used to make your Compose file more efficient and easier to maintain. 
 
Use the prefix `x-` as a top-level element to modularize configurations that you want to reuse. 
Compose ignores any fields that start with `x-`, this is the sole exception where Compose silently ignores unrecognized fields.
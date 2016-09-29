## Testing

Current testing is limited and under development. Proposals should be discussed
via github issues/pull requests following 

https://docs.docker.com/opensource/workflow/advanced-contributing/

## Expect based test

The current integration test uses expect to drive it. Some tips for diagnosing tests for people new to expect.

* [A basic introductory guide to expect](https://gist.github.com/Fluidbyte/6294378)
* Turn on internal diagnostics - add `exp_internal 1` to the top of the expect script. This will allow you to see the matches expect is using and how it sees the strings.
* Interact with a running test - you can use the `interact` to give control of the current process to the user. 

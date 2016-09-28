#DTR Automated Browser Tests

Written using Nightwatch, Selenium, and the Chrome driver.

To run the tests:

1. optionally install nightwatch
`sudo npm install nightwatch -g`

2. run the tests, cd to the `adminserver/ui/` directory then:

if nightwatch is installed globally from step 1:

`DTR_HOST='https://10.0.0.100' DTR_ADMIN='admin' DTR_PASSWORD='orca' nightwatch --test`

or, from the `adminserver/ui/` directory

`DTR_HOST='https://10.0.0.100' DTR_ADMIN='admin' DTR_PASSWORD='orca' $(npm bin)/nightwatch --test`

To run an individual test:

`nightwatch --test test/dtr-status.js`

See:
http://nightwatchjs.org/

![magic](https://cloud.githubusercontent.com/assets/565211/14659269/3d6b4ed0-064f-11e6-914c-b889336aa721.gif)

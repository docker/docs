# Docker Marketing Site

## Requirements

You'll need to have the following items installed before continuing.

  * [Node.js](http://nodejs.org): Use the installer provided on the NodeJS website.
  * [Grunt](http://gruntjs.com/): Run `sudo npm install -g grunt-cli`
  * [Bower](http://bower.io): Run `sudo npm install -g bower`

## Getting Started

Clone the project using git clone. 

Next, navigate into the directory:
```
cd docker-marketing
```

Install all the dependincies (if `npm install` fails, you might need to run it as `sudo`):
```
npm install
bower install
```

While you're working on your project, run:
```
grunt
```

This will assemble all the pages and compile the Sass. You're all set to start working!

## Directory Structure

* `dist`: Static pages are assembled here. This is where you should view the site in your browser. **Don't edit these files directly. They will be overwritten!**
* `src`: This is the directory you'll work in. 
* `src/assets`: All assets (scss, images, fonts, js, etc) go here.
* `src/assets/scss/_settings.scss`: Foundation configuration settings go in here.
* `src/assets/scss/app.scss`: Application styles go here.

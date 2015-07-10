# Docker Toolbox Website

This repository is setup to autodeploy to https://toolbox.docker.com which is mapped to https://docker-toolbox.azurewebsites.net/: a git push refreshes the website with the new content.

Content for the website is in markdown format in the content directory.

In order to preview the website locally, [install Hugo](http://gohugo.io/overview/installing/), then in the directory where you checked out this project:
```
hugo server  -w -v
Web Server is available at http://127.0.0.1:1313/
Press Ctrl+C to stop
```

In order to generate the static html for the site
```
hugo
```
The html is generated in the public directory.

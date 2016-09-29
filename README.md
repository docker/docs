# Documentation for Docker Cloud

To build the documentation locally.

1. Fork this repo.

2. Change to the `docs` directory.

3. Type `make docs`.


## API Documentation Pull Requests

The API documentation for the Docker Cloud project is here

https://github.com/docker/cloud-docs

An extra step is needed when making PR's that modify the API - namely to add the generated HTML output to the PR by doing the following:

1. Make changes to the API's Markdown source.

2. Build the HTML for the API in your local branch.

  a. Change to the `cloud-api-docs` directory

      cd cloud-api-docs
      
  b. Generate the HTML

      make release 
      
3. Add the Markdown together with the HTML to your pull request: 

      $ git add apidocs/*
      
      $ git add docs 

4. Push your changes to orgin.

5. Create a Pull request as you normally wouled. 

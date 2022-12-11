def get_docs_url
  # DEPLOY_URL is from Netlify for preview
  # https://docs.netlify.com/configure-builds/environment-variables/#deploy-urls-and-metadata
  ENV['DEPLOY_URL'] || ENV['DOCS_URL'] || 'https://docs.docker.com'
end

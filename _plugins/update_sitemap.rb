Jekyll::Hooks.register :site, :post_write do |site|
  beginning_time = Time.now
  Jekyll.logger.info "Starting plugin update_sitemap.rb..."

  sitemap_path = File.join(site.dest, 'sitemap.xml')
  # DEPLOY_URL is from Netlify for preview of sitemap on PR
  # https://docs.netlify.com/configure-builds/environment-variables/#deploy-urls-and-metadata
  docs_url = ENV['DEPLOY_URL'] || ENV['DOCS_URL'] || 'http://localhost:4000'

  if File.exist?(sitemap_path)
    sitemap_file = File.read(sitemap_path)
    replace = sitemap_file.gsub!("<loc>/", "<loc>#{docs_url}/")
    Jekyll.logger.info "  Replacing '<loc>/' with '<loc>#{docs_url}/' in #{sitemap_path}"
    File.open(sitemap_path, "w") { |file| file.puts replace }
  end

  end_time = Time.now
  Jekyll.logger.info "done in #{(end_time - beginning_time)} seconds"
end

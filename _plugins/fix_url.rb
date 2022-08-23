require 'jekyll'
require 'octopress-hooks'

require_relative 'util.rb'

module Jekyll
  class FetchRemote < Octopress::Hooks::Site
    def post_write(site)
      beginning_time = Time.now
      Jekyll.logger.info "Starting plugin fix_url.rb..."

      # TODO: use dynamic URL from util.get_docs_url instead of hardcoded one
      #   but needs to remove first all absolute URLs in our code base.
      docs_url = "https://docs.docker.com"

      files = Dir.glob("#{site.dest}/**/*.html")
      Jekyll.logger.info "  Fixing up URLs in #{files.size} html file(s) to be relative"
      files.each do|f|
        text = File.read(f)
        replace = text.gsub(/(<a[^>]* href=\")#{docs_url}/, '\1')
        File.open(f, "w") { |f2| f2.puts replace }
      end

      end_time = Time.now
      Jekyll.logger.info "done in #{(end_time - beginning_time)} seconds"
    end
  end
end

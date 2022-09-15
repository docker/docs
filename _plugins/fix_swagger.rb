require 'jekyll'
require 'octopress-hooks'

require_relative 'util.rb'

module Jekyll
  class FetchRemote < Octopress::Hooks::Site
    def post_read(site)
      beginning_time = Time.now
      Jekyll.logger.info "Starting plugin fix_swagger.rb..."

      # TODO: use dynamic URL from util.get_docs_url instead of hardcoded one
      #   but needs to remove first all absolute URLs in our code base.
      docs_url = "https://docs.docker.com"

      files = Dir.glob(%w[./docker-hub/api/*.yaml ./engine/api/*.yaml])
      Jekyll.logger.info "  Fixing up #{files.size} swagger file(s)..."
      files.each do |f|
        Jekyll.logger.info "    #{f}"
        text = File.read(f)
        replace = text.gsub(docs_url, "")
        File.open(f, "w") { |f2| f2.puts replace }
      end

      end_time = Time.now
      Jekyll.logger.info "done in #{(end_time - beginning_time)} seconds"
    end
  end
end

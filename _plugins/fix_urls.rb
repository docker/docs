require 'jekyll'
require 'octopress-hooks'

module Jekyll

  class FetchRemote < Octopress::Hooks::Site
    def post_read(site)
      beginning_time = Time.now
      Jekyll.logger.info "Starting plugin fix_urls.rb..."

      Jekyll.logger.info "  Fixing up URLs in swagger files"
      Dir.glob(%w[./docker-hub/api/*.yaml ./engine/api/*.yaml]) do |file_name|
        Jekyll.logger.info "    #{file_name}"
        text = File.read(file_name)
        replace = text.gsub!("https://docs.docker.com", "")
        File.open(file_name, "w") { |file| file.puts replace }
      end

      end_time = Time.now
      Jekyll.logger.info "done in #{(end_time - beginning_time)} seconds"
    end
  end

end

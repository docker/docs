require 'jekyll'
require 'octopress-hooks'

module Jekyll
  class UpdateApiToc < Octopress::Hooks::Site
    def pre_read(site)
      beginning_time = Time.now
      Jekyll.logger.info "Starting plugin update_api_toc.rb..."

      if File.file?("_config.yml") && File.file?("_data/toc.yaml")
        # substitute the "{site.latest_engine_api_version}" in the title for the latest
        # API docs, based on the latest_engine_api_version parameter in _config.yml
        engine_ver = site.config['latest_engine_api_version']
        toc_file = File.read("_data/toc.yaml")
        replace = toc_file.gsub!("{{ site.latest_engine_api_version }}", engine_ver)
        if toc_file != replace
          Jekyll.logger.info "  Replacing '{{ site.latest_engine_api_version }}' with #{engine_ver} in _data/toc.yaml"
          File.open("_data/toc.yaml", "w") { |file| file.puts replace }
        end
      end

      end_time = Time.now
      Jekyll.logger.info "done in #{(end_time - beginning_time)} seconds"
    end
  end
end

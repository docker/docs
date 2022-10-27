require 'git'
require 'jekyll'
require 'octopress-hooks'

module Jekyll
  class LastModifiedAt < Octopress::Hooks::Site
    DATE_FORMAT = '%Y-%m-%d %H:%M:%S %z'
    def pre_render(site)
      if get_docs_url == "http://localhost:4000"
        # Do not generate last_modified_at for local development
        return
      end

      beginning_time = Time.now
      Jekyll.logger.info "Starting plugin last_modified_at.rb..."

      git = Git.open(site.source)
      site.pages.each do |page|
        next if page.relative_path == "redirect.html"
        next unless File.extname(page.relative_path) == ".md" || File.extname(page.relative_path) == ".html"
        unless page.data.key?('last_modified_at')
          begin
            page.data['last_modified_at'] = git.log.path(page.relative_path).first.date.strftime(DATE_FORMAT)
          rescue => e
            # Ignored
          end
        end
        puts"  #{page.relative_path}\n    last_modified_at: #{page.data['last_modified_at']}"
      end

      end_time = Time.now
      Jekyll.logger.info "done in #{(end_time - beginning_time)} seconds"
    end
  end
end

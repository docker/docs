require 'git'
require 'jekyll'
require 'octopress-hooks'

module Jekyll
  class LastModifiedAt < Octopress::Hooks::Site
    DATE_FORMAT = '%Y-%m-%d %H:%M:%S %z'

    def current_last_modified_at(site, page)
      if page.data.key?('last_modified_at')
        return page.data['last_modified_at']
      end
      site.config['defaults'].map do |set|
        if set['values'].key?('last_modified_at') && set['scope']['path'].include?(page.relative_path)
          return set['values']['last_modified_at']
        end
      end.compact
      nil
    end

    def pre_render(site)
      beginning_time = Time.now
      Jekyll.logger.info "Starting plugin last_modified_at.rb..."

      git = Git.open(site.source)
      use_file_mtime = get_docs_url == "http://localhost:4000" && ENV['DOCS_ENFORCE_GIT_LOG_HISTORY'] == "0"
      site.pages.sort!{|l,r| l.relative_path <=> r.relative_path }.each do |page|
        next if page.relative_path == "redirect.html"
        next unless File.extname(page.relative_path) == ".md" || File.extname(page.relative_path) == ".html"
        page.data['last_modified_at'] = current_last_modified_at(site, page)
        set_mode = "frontmatter"
        path_override = ""
        if page.data['last_modified_at'].nil?
          page_relative_path = page.relative_path
          if page.data.key?('datafolder') && page.data.key?('datafile')
            page_relative_path = File.join('_data',  page.data['datafolder'], "#{page.data['datafile']}.yaml")
            path_override = "\n    override: #{page_relative_path}"
          end
          begin
            if use_file_mtime
              # Use file's mtime for local development
              page.data['last_modified_at'] = File.mtime(page_relative_path).strftime(DATE_FORMAT)
              set_mode = "mtime"
            else
              page.data['last_modified_at'] = git.log.path(page_relative_path).first.date.strftime(DATE_FORMAT)
              set_mode = "git"
            end
          rescue => e
            begin
              page.data['last_modified_at'] = File.mtime(page_relative_path).strftime(DATE_FORMAT)
              set_mode = "mtime"
            rescue => e
              page.data['last_modified_at'] = Time.now.strftime(DATE_FORMAT)
              set_mode = "rescue"
            end
          end
        end
        puts"  #{page.relative_path}#{path_override}\n    last_modified_at(#{set_mode}): #{page.data['last_modified_at']}"
      end

      end_time = Time.now
      Jekyll.logger.info "done in #{(end_time - beginning_time)} seconds"
    end
  end
end

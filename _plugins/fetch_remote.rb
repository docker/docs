require 'archive/zip'
require 'jekyll'
require 'json'
require 'octopress-hooks'
require 'open-uri'
require 'rake'

module Jekyll
  def self.download(url, dest)
    uri = URI.parse(url)
    result = File.join(dest, File.basename(uri.path))
    Jekyll.logger.info "    Downloading #{url}"
    IO.copy_stream(URI.open(url), result)
    return result
  end

  class FetchRemote < Octopress::Hooks::Site
    def pre_read(site)
      beginning_time = Time.now
      Jekyll.logger.info "Starting plugin fetch_remote.rb..."
      site.config['fetch-remote'].each do |entry|
        Jekyll.logger.info "  Repo #{entry['repo']} (#{entry['ref']})"
        Dir.mktmpdir do |tmpdir|
          tmpfile = Jekyll.download("#{entry['repo']}/archive/#{entry['ref']}.zip", tmpdir)
          Dir.mktmpdir do |ztmpdir|
            Jekyll.logger.info "    Extracting #{tmpfile}"
            Archive::Zip.extract(
              tmpfile,
              ztmpdir,
              :create => true
            )
            entry['paths'].each do |path|
              Jekyll.logger.info "    Copying files to ./#{path['dest']}/"
              files = FileList[]
              path['src'].each do |src|
                if "#{src}".start_with?("!")
                  files.exclude(File.join(ztmpdir, "*/"+"#{src}".delete_prefix("!")))
                else
                  files.include(File.join(ztmpdir, "*/#{src}"))
                end
              end
              files.each do |file|
                Jekyll.logger.info "      #{file.delete_prefix(ztmpdir)}"
              end
              FileUtils.mkdir_p path['dest']
              FileUtils.cp_r(files, path['dest'])
            end
          end
        end
      end

      end_time = Time.now
      Jekyll.logger.info "done in #{(end_time - beginning_time)} seconds"
    end
  end
end

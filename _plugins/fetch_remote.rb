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
    puts "    Downloading #{url}"
    IO.copy_stream(URI.open(url), result)
    return result
  end

  def self.copy(src, dest)
    if (tmp = Array.try_convert(src))
      tmp.each do |s|
        s = File.path(s)
        yield s, File.join(dest, File.basename(s))
      end
    else
      src = File.path(src)
      if File.directory?(dest)
        yield src, File.join(dest, File.basename(src))
      else
        yield src, File.path(dest)
      end
    end
  end

  class FetchRemote < Octopress::Hooks::Site
    def pre_read(site)
      beginning_time = Time.now
      puts "Starting plugin fetch_remote.rb..."
      site.config['fetch-remote'].each do |entry|
        puts "  Repo #{entry['repo']} (#{entry['ref']})"
        Dir.mktmpdir do |tmpdir|
          tmpfile = Jekyll.download("#{entry['repo']}/archive/#{entry['ref']}.zip", tmpdir)
          Dir.mktmpdir do |ztmpdir|
            puts "    Extracting #{tmpfile}"
            Archive::Zip.extract(
              tmpfile,
              ztmpdir,
              :create => true
            )
            entry['paths'].each do |path|
              FileUtils.mkdir_p path['dest']
              puts "    Copying files"

              # prepare file list to be copied
              files = FileList[]
              path['src'].each do |src|
                if "#{src}".start_with?("!")
                  files.exclude(File.join(ztmpdir, "*/"+"#{src}".delete_prefix("!")))
                else
                  files.include(File.join(ztmpdir, "*/#{src}"))
                end
              end

              files.each do |file|
                Jekyll.copy(file, path['dest']) do |s, d|
                  s = File.realpath(s)
                  # traverse source directory
                  FileUtils::Entry_.new(s, nil, false).wrap_traverse(proc do |ent|
                    file_clean = ent.path.delete_prefix(ztmpdir).split("/").drop(2).join("/")
                    destent = FileUtils::Entry_.new(d, ent.rel, false)
                    puts "      #{file_clean} => #{destent.path}"
                    ent.copy destent.path

                    next unless File.file?(ent.path) && File.extname(ent.path) == ".md"
                    # set edit url for markdown files in site config defaults
                    edit_url = "#{entry['repo']}/edit/#{entry['ref']}/#{file_clean}"
                    puts "        edit_url: #{edit_url}"
                    site.config['defaults'] << {
                      "scope" => { "path" => destent.path },
                      "values" => { "edit_url" => edit_url },
                    }
                  end, proc do |_| end)
                end
              end
            end
          end
        end
      end

      end_time = Time.now
      puts "done in #{(end_time - beginning_time)} seconds"
    end
  end
end

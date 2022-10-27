require 'front_matter_parser'
require 'git'
require 'jekyll'
require 'json'
require 'octopress-hooks'
require 'rake'

require_relative 'util.rb'

module Jekyll
  class FetchRemote < Octopress::Hooks::Site
    priority :highest

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

    def self.resolve_line_numbers(first, last)
      if first.nil? && last.nil?
        first = 0
        last  = -1
      elsif last.nil?
        last = first
      end
      [first.to_i, last.to_i]
    end

    def pre_read(site)
      beginning_time = Time.now
      puts "Starting plugin fetch_remote.rb..."

      fetch_depth = get_docs_url == "http://localhost:4000" ? 1 : 0
      site.config['fetch-remote'].each do |entry|
        puts "  Repo #{entry['repo']}"

        gituri = Git::URL.parse(entry['repo'])
        clonedir = "#{Dir.tmpdir}/docker-docs-clone#{gituri.path}"
        if Dir.exist?(clonedir)
          puts "    Opening #{clonedir}"
          begin
            git = Git.open(clonedir)
            puts "    Fetching #{entry['ref']}"
            git.fetch
            git.checkout(entry['ref'])
          rescue => e
            FileUtils.rm_rf(clonedir)
            puts "    Cloning repository into #{clonedir}"
            git = Git.clone("#{entry['repo']}.git", Pathname.new(clonedir), branch: entry['ref'], depth: fetch_depth)
          end
        else
          puts "    Cloning repository into #{clonedir}"
          git = Git.clone("#{entry['repo']}.git", Pathname.new(clonedir), branch: entry['ref'], depth: fetch_depth)
        end

        entry['paths'].each do |path|
          if File.extname(path['dest']) != ""
            if path['src'].size > 1
              raise "Cannot use file destination #{path['dest']} with multiple sources"
            end
            FileUtils.mkdir_p File.dirname(path['dest'])
          else
            FileUtils.mkdir_p path['dest']
          end

          puts "    Copying files"

          # prepare file list to be copied
          files = FileList[]
          path['src'].each do |src|
            if "#{src}".start_with?("!")
              files.exclude(File.join(clonedir, "/"+"#{src}".delete_prefix("!")))
            else
              files.include(File.join(clonedir, "/#{src}"))
            end
          end

          files.each do |file|
            FetchRemote.copy(file, path['dest']) do |s, d|
              s = File.realpath(s)
              # traverse source directory
              FileUtils::Entry_.new(s, nil, false).wrap_traverse(proc do |ent|
                file_clean = ent.path.delete_prefix(clonedir).split("/").drop(1).join("/")
                destent = FileUtils::Entry_.new(d, ent.rel, false)
                puts "      #{file_clean} => #{destent.path}"

                if File.file?(destent.path)
                  fmp = FrontMatterParser::Parser.parse_file(destent.path)
                  if fmp['fetch_remote'].nil?
                    raise "Local file #{destent.path} already exists"
                  end
                  line_start, line_end = FetchRemote.resolve_line_numbers(fmp['fetch_remote'].kind_of?(Hash) ? fmp['fetch_remote']['line_start'] : nil, fmp['fetch_remote'].kind_of?(Hash) ? fmp['fetch_remote']['line_end'] : nil)
                  lines = File.readlines(ent.path)[line_start..line_end]
                  File.open(destent.path, "a") { |fow| fow.puts lines.join }
                else
                  ent.copy destent.path
                end

                next unless File.file?(ent.path) && File.extname(ent.path) == ".md"
                # set edit and issue url and remote info for markdown files in site config defaults
                edit_url = "#{entry['repo']}/edit/#{entry['default_branch']}/#{file_clean}"
                issue_url = "#{entry['repo']}/issues/new?body=File: [#{file_clean}](#{get_docs_url}/#{destent.path.sub(/#{File.extname(destent.path)}$/, '')}/)"
                last_modified_at = git.log.path(file_clean).first.date.strftime(LastModifiedAt::DATE_FORMAT)
                puts "        edit_url:         #{edit_url}"
                puts "        issue_url:        #{issue_url}"
                puts "        last_modified_at: #{last_modified_at}"
                site.config['defaults'] << {
                  "scope" => { "path" => destent.path },
                  "values" => {
                    "edit_url" => edit_url,
                    "issue_url" => issue_url,
                    "last_modified_at" => last_modified_at,
                  },
                }
              end, proc do |_| end)
            end
          end
        end
      end

      end_time = Time.now
      puts "done in #{(end_time - beginning_time)} seconds"
    end
  end
end

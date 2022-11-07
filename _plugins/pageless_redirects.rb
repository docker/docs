require 'jekyll'
require 'octopress-hooks'

module Jekyll
  class RedirectPage < Jekyll::Page
    def initialize(site, src, redirect_to)
      puts "  #{src} => #{redirect_to}"
      @site = site
      @base = site.source
      @dir = src
      @name = "index.html"
      process(@name)
      @data = {
        "sitemap" => false,
        "redirect_to" => redirect_to
      }
    end
  end

  class PagelessRedirects < Octopress::Hooks::Site
    def post_read(site)
      beginning_time = Time.now
      puts "Starting plugin pageless_redirects.rb..."

      if File.file?("_redirects.yml")
        rd = YAML.load_file("_redirects.yml")
        rd.each do |redirect_to, srcs|
          srcs.each do |src|
            site.pages << RedirectPage.new(site, src, redirect_to)
          end
        end
      end

      end_time = Time.now
      puts "done in #{(end_time - beginning_time)} seconds"
    end
  end
end

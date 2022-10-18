require_relative 'util.rb'

module Jekyll
  class ConfigGenerator < Generator
    safe true
    priority :highest

    def generate(site)
      site.config['docs_url'] = get_docs_url
    end
  end
end

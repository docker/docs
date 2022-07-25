module Jekyll
  class IncludeRemoteTag < Liquid::Tag
    def initialize(tag_name, params, tokens)
      @page, @line_start, @line_end = params.split
      @line_start, @line_end = resolve_line_numbers(@line_start, @line_end)
      super
    end

    def render(context)
      site = context.registers[:site]
      page = context.registers[:page]

      beginning_time = Time.now
      Jekyll.logger.info "Starting plugin include_remote.rb..."

      if context[@page.strip]
        @page = context[@page.strip]
      end

      inc = File.join("_includes", @page)
      Jekyll.logger.info "  Inject #{inc} to #{page['path']}"

      lines = File.readlines(inc)[@line_start..@line_end]

      site.config['defaults'].each do |default|
        if default['scope']['path'] == inc
          page['edit_url'] = default['values']['edit_url']
          page['issue_url'] = default['values']['issue_url']
          Jekyll.logger.info "    edit_url:  #{page['edit_url']}"
          Jekyll.logger.info "    issue_url: #{page['issue_url']}"
          break
        end
      end

      end_time = Time.now
      Jekyll.logger.info "done in #{(end_time - beginning_time)} seconds"

      lines.join
    end

    def resolve_line_numbers(first, last)
      if first.nil? && last.nil?
        first = 0
        last  = -1
      elsif last.nil?
        last = first
      end
      [first.to_i, last.to_i]
    end
  end

end

Liquid::Template.register_tag('include_remote', Jekyll::IncludeRemoteTag)

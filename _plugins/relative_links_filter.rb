module Jekyll
  # This custom Filter is used to fix up links to markdown pages that weren't
  # resolved by Jekyll (or the "jekyll-relative-links" plugin). We need this hack
  # to work around a bug in the "jekyll-relative-links" plugin;
  #
  # As reported in https://github.com/benbalter/jekyll-relative-links/issues/54,
  # (relative) links to markdown pages in includes are not processed by Jekyll.
  # This means that our reference pages (which use includes) have broken links.
  # We could work around this by modifying the markdown for those pages to use
  # "absolute" "html" links (/link/to/other/page/#some-anchor), but doing so
  # would render the links broken when viewed on GitHub. Instead,  we're fixing
  # them up here, until the bug is fixed upstream.
  #
  # A second bug (https://github.com/benbalter/jekyll-relative-links/issues/61),
  # causes (relative) links to markdown pages to not be resolved if the link's
  # caption/title is wrapped. This bug is currently not handled by this plugin,
  # but could possibly be addressed by modifying the TITLE_REGEX.
  #
  # This plugin is based on code in the jekyll-relative-links plugin, but takes
  # some shortcuts;
  #
  # - We use the code from jekyll-relative-links plugin to find/extract links
  #   on the page
  # - Relative links are converted to absolute links, using the path of the
  #   markdown source file that's passed as argument
  # - After conversion to an absolute link, we strip the ".md" extension; no
  #   attempt is made to resolve the file that's linked to. This is different
  #   from the jekyll-relative-links plugin, which _does_ resolve the linked
  #   file. This functionality could be added in future by someone who has
  #   more experience with Ruby.
  module RelativeLinksFilter
    attr_accessor :site, :config

    # Use Jekyll's native relative_url filter
    include Jekyll::Filters::URLFilters

    LINK_TEXT_REGEX = %r!(.*?)!.freeze
    FRAGMENT_REGEX = %r!(#.+?)?!.freeze
    TITLE_REGEX = %r{(\s+"(?:\\"|[^"])*(?<!\\)"|\s+"(?:\\'|[^'])*(?<!\\)')?}.freeze
    FRAG_AND_TITLE_REGEX = %r!#{FRAGMENT_REGEX}#{TITLE_REGEX}!.freeze
    INLINE_LINK_REGEX = %r!\[#{LINK_TEXT_REGEX}\]\(([^\)]+?)#{FRAG_AND_TITLE_REGEX}\)!.freeze
    REFERENCE_LINK_REGEX = %r!^\s*?\[#{LINK_TEXT_REGEX}\]: (.+?)#{FRAG_AND_TITLE_REGEX}\s*?$!.freeze
    LINK_REGEX = %r!(#{INLINE_LINK_REGEX}|#{REFERENCE_LINK_REGEX})!.freeze

    def replace_relative_links(input, source_path)
      url_base = File.dirname("/" + source_path)
      input = input.dup.gsub(LINK_REGEX) do |original|
        link = link_parts(Regexp.last_match)
        next original unless replaceable_link?(link.path)

        path = path_from_root(link.path, url_base)
        url = path.gsub(".md", "/")
        next original unless url

        link.path = url
        replacement_text(link)
      end
    end

    private

    # Stores info on a Markdown Link (avoid rubocop's Metrics/ParameterLists warning)
    Link = Struct.new(:link_type, :text, :path, :fragment, :title)

    def link_parts(matches)
      last_inline = 5
      link_type     = matches[2] ? :inline : :reference
      link_text     = matches[link_type == :inline ? 2 : last_inline + 1]
      relative_path = matches[link_type == :inline ? 3 : last_inline + 2]
      fragment      = matches[link_type == :inline ? 4 : last_inline + 3]
      title         = matches[link_type == :inline ? 5 : last_inline + 4]
      Link.new(link_type, link_text, relative_path, fragment, title)
    end

    def path_from_root(relative_path, url_base)
      relative_path.sub!(%r!\A/!, "")
      absolute_path = File.expand_path(relative_path, url_base)
      absolute_path.sub(%r!\A#{Regexp.escape(Dir.pwd)}/!, "")
    end

    # @param link [Link] A Link object describing the markdown link to make
    def replacement_text(link)
      link.path << link.fragment if link.fragment

      if link.link_type == :inline
        "[#{link.text}](#{link.path}#{link.title})"
      else
        "\n[#{link.text}]: #{link.path}#{link.title}"
      end
    end

    def absolute_url?(string)
      return unless string

      Addressable::URI.parse(string).absolute?
    rescue Addressable::URI::InvalidURIError
      nil
    end

    def fragment?(string)
      string&.start_with?("#")
    end

    def replaceable_link?(string)
      !fragment?(string) && !absolute_url?(string)
    end

    def global_entry_filter
      @global_entry_filter ||= Jekyll::EntryFilter.new(site)
    end
  end
end

Liquid::Template.register_filter(Jekyll::RelativeLinksFilter)

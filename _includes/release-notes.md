{% assign data = site.data.release-notes[page.release] %}
{% assign metadata = data.metadata %}
{% for release_name in metadata.releases %}
{% capture release_handle %}{{ release_name | replace: ".", "" }}{% endcapture %}
{% assign release = data[release_handle] %}
## {{ release.version | remove: "v" }}

<em class="release-date">{{ release.date }}</em>

<table class="release-notes"><tbody>
{%- for entry in release.entries %}
  <tr>
    <td><span class="release-tag release-tag-{{ entry.type }}">{{ entry.type | upcase }}</span></td>
    <td>
      {{ entry.content | markdownify }}
      <ul class="fa-ul">
        {%- for issue in entry.issues %}{% assign issue_sp = issue | split: "#" %}
        <li>
          <span class="fa-li"><i class="fa-brands fa-github"></i></span><a href="https://github.com/{{ issue_sp.first }}/issues/{{ issue_sp.last }}" target="_blank" rel="noopener" class="_">{{ issue_sp.first }}#{{ issue_sp.last }}</a>
        </li>
        {%- endfor %}
      </ul>
    </td>
  </tr>
{%- endfor %}
</tbody></table>

> For more details, see the complete release notes in the [{{ metadata.name }} repository](https://github.com/{{ metadata.repo }}/releases/tag/{{ release.version }}){:target="_blank" rel="noopener" class="_"}.
{% endfor -%}

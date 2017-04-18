{% comment %}
  Generates a Github PR URL from three parameters

  Usage:
    {% include github-pr.md org=docker repo=docker pr=12345 %}

    If you omit the org or repo, they default to docker.
    If you omit the pr, it defaults to NULL.

  Output:
    [#12345](https://github.com/moby/moby/pull/12345)
{% endcomment %}{% assign org = include.org | default: "docker" %}{% assign repo = include.repo | default: "docker" %}{% assign pr = include.pr | default: NULL %}{% assign github-url="https://github.com" %}{% capture pr-link %}[#{{ pr }}]({{ github-url }}/{{ org }}/{{ repo }}/pull/{{ pr }}){% endcapture %}{{ pr-link | strip_newlines }}
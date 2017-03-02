---
title: Learn Docker
redirect_from:
- /engine/quickstart/
---

In this section, you can explore various ways to get up to speed on
Docker workflows.

<table>
<tr valign="top">
<td width="50%">
{% capture basics %}
## Learn the basics of Docker

The basic tutorial introduces Docker concepts, tools, and commands. The examples show you how to build, push,
and pull Docker images, and run them as containers. This
tutorial stops short of teaching you how to deploy applications.
{% endcapture %}{{ basics | markdownify }}
</td>
<td width="50%">

{% capture apps %}
## Define and deploy applications

The define-and-deploy tutorial shows how to relate
containers to each other and define them as services in an application that is ready to deploy at scale in a
production environment. Highlights Compose Version 3 new features and swarm mode.
{% endcapture %}{{ apps | markdownify }}

</td></tr>

<tr valign="top">
<td width="50%">
{% capture basics %}
[Start the basic tutorial](/engine/getstarted/){: class="button secondary-btn"}
{% endcapture %}{{ basics | markdownify }}
</td>
<td width="50%">
{% capture apps %}
[Start the application tutorial](/engine/getstarted-voting-app/){: class="button secondary-btn"}
{% endcapture %}{{ apps | markdownify }}
</td>
</tr>
</table>

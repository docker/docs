
<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#mac-copy-keys" data-group="mac">Mac</a></li>
  <li><a data-toggle="tab" data-target="#win-copy-keys" data-group="win">Windows</a></li>
  <li><a data-toggle="tab" data-target="#linux-copy-keys" data-group="linux">Linux</a></li>
</ul>
<div class="tab-content">
<div id="mac-copy-keys" class="tab-pane fade in active">
<br>
{% capture mac-content-copy %}

Copy the public SSH key to your clipboard.

```none
$ pbcopy < ~/.ssh/id_rsa.pub
```

If your SSH key file has a different name than the example code, modify the
filename to match your current setup.

>**Tip:** If you don't have `pbcopy`, you navigate to the hidden `.ssh`
folder, open the file in a text editor, and copy it to your clipboard.
For example: `$ atom ~/.ssh/id_rsa.pub`

{% endcapture %}
{{ mac-content-copy | markdownify }}
<hr>
</div>

<div id="win-copy-keys" class="tab-pane fade">
<br>
{% capture win-content-copy %}

Copy the public SSH key to your clipboard.

```none
$ clip < ~/.ssh/id_rsa.pub
```

If your SSH key file has a different name than the example code, modify the
filename to match your current setup.

>**Tip:** If `clip` doesn't work, navigate the hidden `.ssh`
folder, open the file in a text editor, and copy it to your clipboard.
For example: `$ notepad ~/.ssh/id_rsa.pub`

{% endcapture %}
{{ win-content-copy | markdownify }}
<hr>
</div>

<div id="linux-copy-keys" class="tab-pane fade">
<br>
{% capture linux-content-copy %}

If you don't already have it, install `xclip`. (The example uses `apt-get` to install, but you might want to use another package installer like `yum`.)

```none
$ sudo apt-get install xclip
```

Copy the SSH key to your clipboard.

```none
$ xclip -sel clip < ~/.ssh/id_rsa.pub
```

>**Tip:** If you `xclip` isn't working, navigate to hidden `.ssh` folder,
open the file in a text editor, and copy it to your clipboard.
For example: `$ vi ~/.ssh/id_rsa.pub`

{% endcapture %}
{{ linux-content-copy | markdownify }}
<hr>
</div>
</div>

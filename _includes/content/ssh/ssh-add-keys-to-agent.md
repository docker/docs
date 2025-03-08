
<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#mac-add-keys" data-group="mac">Mac</a></li>
  <li><a data-toggle="tab" data-target="#win-add-keys" data-group="win">Windows</a></li>
  <li><a data-toggle="tab" data-target="#linux-add-keys" data-group="linux">Linux</a></li>
</ul>
<div class="tab-content">
<div id="mac-add-keys" class="tab-pane fade in active">
<br>
{% capture mac-content-add %}
1.  Start the `ssh-agent` in the background using the command `eval "$(ssh-agent -s)"`. You get the agent process ID in return.

    ```none
    eval "$(ssh-agent -s)"
    Agent pid 59566
    ```

2.  On macOS Sierra 10.12.2 or newer, modify your
`~/.ssh/config` file to automatically load keys into the `ssh-agent` and store
passphrases in your keychain.

    ```none
    Host *
     AddKeysToAgent yes
     UseKeychain yes
     IdentityFile ~/.ssh/id_rsa
    ```

3.  Add your SSH private key to the ssh-agent, using the default macOS `ssh-add` command.

    ```none
    $ ssh-add -K ~/.ssh/id_rsa
    ```

    If you created your key with a different name or have an existing key
    with  a different name, replace `id_rsa` in the command with the
    name of your private key file.

{% endcapture %}
{{ mac-content-add | markdownify }}
<hr>
</div>

<div id="win-add-keys" class="tab-pane fade">
<br>
{% capture win-content-add %}

1.  Start the `ssh-agent` in the background.

    ```none
    eval "$(ssh-agent -s)"
    Agent pid 59566
    ```

2.  Add your SSH private key to the ssh-agent.

    ```none
    $ ssh-add -K ~/.ssh/id_rsa
    ```

    If you created your key with a different name or have an existing key
    with  a different name, replace `id_rsa` in the command with the
    name of your private key file.

{% endcapture %}
{{ win-content-add | markdownify }}
<hr>
</div>

<div id="linux-add-keys" class="tab-pane fade">
<br>
{% capture linux-content-add %}

1.  Start the `ssh-agent` in the background.

    ```none
    eval "$(ssh-agent -s)"
    Agent pid 59566
    ```

2.  Add your SSH private key to the ssh-agent.

    ```none
    $ ssh-add -K ~/.ssh/id_rsa
    ```

    If you created your key with a different name or have an existing key
    with  a different name, replace `id_rsa` in the command with the
    name of your private key file.

{% endcapture %}
{{ linux-content-add | markdownify }}
<hr>
</div>
</div>

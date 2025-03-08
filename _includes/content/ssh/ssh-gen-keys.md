
<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#mac-key-gen" data-group="mac">Mac</a></li>
  <li><a data-toggle="tab" data-target="#win-key-gen" data-group="win">Windows</a></li>
  <li><a data-toggle="tab" data-target="#linux-key-gen" data-group="linux">Linux</a></li>
</ul>
<div class="tab-content">
<div id="mac-key-gen" class="tab-pane fade in active">
<br>
{% capture mac-content-gen %}
1.  Open a command-line terminal.

2.  Paste the text below, substituting in your GitHub email address.

    ```none
    ssh-keygen -t rsa -b 4096 -C "your_email@example.com"
    ```

    This creates a new SSH key, using the provided email as a label.

    ```none
    Generating public/private rsa key pair.
    ```

3.  When prompted with "Enter a file in which to save the key", press the Return key (Enter) to accept the default location.

    ```none
    Enter a file in which to save the key (/Users/you/.ssh/id_rsa):
    ```

4.  At the prompt, type a secure passphrase, and re-enter as prompted.

    ```none
    Enter passphrase (empty for no passphrase):
    Enter same passphrase again:
    ```
{% endcapture %}
{{ mac-content-gen | markdownify }}
<hr>
</div>

<div id="win-key-gen" class="tab-pane fade">
<br>
{% capture win-content-gen %}
1.  Open Git Bash.

2.  Paste the text below, substituting in your GitHub email address.

    ```none
    ssh-keygen -t rsa -b 4096 -C "your_email@example.com"
    ```

    This creates a new SSH key, using the provided email as a label.

    ```none
    Generating public/private rsa key pair.
    ```

3.  When prompted with "Enter a file in which to save the key", press the Return key (Enter) to accept the default location.

    ```none
    Enter a file in which to save the key (c/Users/you/.ssh/id_rsa):
    ```

4.  At the prompt, type a secure passphrase, and re-enter as prompted.

    ```none
    Enter passphrase (empty for no passphrase):
    Enter same passphrase again:
    ```
{% endcapture %}
{{ win-content-gen | markdownify }}
<hr>
</div>

<div id="linux-key-gen" class="tab-pane fade">
<br>
{% capture linux-content-gen %}
1.  Open a command-line terminal.

2.  Paste the text below, substituting in your GitHub email address.

    ```none
    ssh-keygen -t rsa -b 4096 -C "your_email@example.com"
    ```

    This creates a new SSH key, using the provided email as a label.

    ```none
    Generating public/private rsa key pair.
    ```

3.  When prompted with "Enter a file in which to save the key", press the Return key (Enter) to accept the default location.

    ```none
    Enter a file in which to save the key (/home/you/.ssh/id_rsa):
    ```

4.  At the prompt, type a secure passphrase, and re-enter as prompted.

    ```none
    Enter passphrase (empty for no passphrase):
    Enter same passphrase again:
    ```
{% endcapture %}
{{ linux-content-gen | markdownify }}
<hr>
</div>
</div>

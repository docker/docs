
<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#mac-find-keys" data-group="mac">Mac</a></li>
  <li><a data-toggle="tab" data-target="#win-find-keys" data-group="win">Windows</a></li>
  <li><a data-toggle="tab" data-target="#linux-find-keys" data-group="linux">Linux</a></li>
</ul>
<div class="tab-content">
<div id="mac-find-keys" class="tab-pane fade in active">
<br>
{% capture mac-content-find %}

1.  Open a command-line terminal.

    ```none
    $ ls -al ~/.ssh
    ```

    This lists files in your `.ssh` directory.

2.  Check to see if you already have a SSH keys you can use.

    Default file names for public keys are:

    * id_dsa.pub
    * id_ecdsa.pub
    * id_ed25519.pub
    * id_rsa.pub

    Here are example results showing a public and private key pair with the default names:

    ```none
    drwx------   8 me  staff   272 Mar 27 14:04 .
    drwxr-xr-x+ 69 me  staff  2346 Apr  7 10:03 ..
    -rw-r--r--   1 me  staff   420 Mar 27 14:04 config
    -rw-------   1 me  staff  3326 Mar 27 14:01 id_rsa
    -rw-r--r--   1 me  staff   752 Mar 27 14:01 id_rsa.pub
    ```

    The file `id_rsa` contains the private key which resides on the local machine, and `id_rsa.pub` is the public key we can provide to a remote account.

{% endcapture %}
{{ mac-content-find | markdownify }}
<hr>
</div>

<div id="win-find-keys" class="tab-pane fade">
<br>
{% capture win-content-find %}

1.  Open Git Bash.

    ```none
    $ ls -al ~/.ssh
    ```

    This lists files in your `.ssh` directory.

2.  Check to see if you already have SSH keys you can use.

    Default file names for public keys are:

    * id_dsa.pub
    * id_ecdsa.pub
    * id_ed25519.pub
    * id_rsa.pub

    Here are example results showing a public and private key pair with the default names:

    ```none
    drwx------   8 me  staff   272 Mar 27 14:04 .
    drwxr-xr-x+ 69 me  staff  2346 Apr  7 10:03 ..
    -rw-r--r--   1 me  staff   420 Mar 27 14:04 config
    -rw-------   1 me  staff  3326 Mar 27 14:01 id_rsa
    -rw-r--r--   1 me  staff   752 Mar 27 14:01 id_rsa.pub
    ```

    The file `id_rsa` contains the private key which resides on the local machine, and `id_rsa.pub` is the public key we can provide to a remote account.

{% endcapture %}
{{ win-content-find | markdownify }}
<hr>
</div>

<div id="linux-find-keys" class="tab-pane fade">
<br>
{% capture linux-content-find %}

1.  Open a command-line terminal.

    ```none
    $ ls -al ~/.ssh
    ```

    This lists files in your `.ssh` directory.

2.  Check to see if you already have a SSH keys you can use.

    Default file names for public keys are:

    * id_dsa.pub
    * id_ecdsa.pub
    * id_ed25519.pub
    * id_rsa.pub

    Here are example results showing a public and private key pair with the default names:

    ```none
    drwx------   8 me  staff   272 Mar 27 14:04 .
    drwxr-xr-x+ 69 me  staff  2346 Apr  7 10:03 ..
    -rw-r--r--   1 me  staff   420 Mar 27 14:04 config
    -rw-------   1 me  staff  3326 Mar 27 14:01 id_rsa
    -rw-r--r--   1 me  staff   752 Mar 27 14:01 id_rsa.pub
    ```

    The file `id_rsa` contains the private key which resides on the local machine, and `id_rsa.pub` is the public key we can provide to a remote account.

{% endcapture %}
{{ linux-content-find | markdownify }}
<hr>
</div>
</div>

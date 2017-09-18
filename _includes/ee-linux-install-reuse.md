{% assign section = include.section %}

{% comment %}

Include a chunk of this file, using variables already set in the file
where you want to reuse the chunk.

Usage: {% include ee-linux-install-reuse.md section="ee-install-intro" %}

{% endcomment %}



{% if section == "ee-install-intro" %}

To get started with Docker EE on {{ linux-dist-long }}, make sure you
[meet the prerequisites](#prerequisites), then
[install Docker](#install-docker-ee).

{% elsif section == "ee-url-intro" %}

To install Docker Enterprise Edition (Docker EE), you need to know the Docker EE
repository URL associated with your trial or subscription. These instructions
work for Docker EE for {{ linux-dist-long }} and for Docker EE for Linux, which
includes access to Docker EE for all Linux distributions. To get this
information:

- Go to [https://store.docker.com/my-content](https://store.docker.com/my-content).
- Each subscription or trial you have access to is listed. Click the **Setup**
  button for **Docker Enterprise Edition for {{ linux-dist-long }}**.
- Copy the URL from the field labeled
  **Copy and paste this URL to download your Edition**.

Use this URL when you see the placeholder text `<DOCKER-EE-URL>`.

To learn more about Docker EE, see
[Docker Enterprise Edition](https://www.docker.com/enterprise-edition/){: target="_blank" class="_" }.



{% elsif section == "ways-to-install" %}

You can install Docker EE in different ways, depending on your needs:

- Most users
  [set up Docker's repositories](#install-using-the-repository) and install
  from them, for ease of installation and upgrade tasks. This is the
  recommended approach.

- Some users download the {{ package-format }} package and install it manually
  and manage upgrades completely manually. This is useful in situations such as
  installing Docker on air-gapped systems with no access to the internet.



{% elsif section == "set-up-yum-repo" %}

1.  Remove any existing Docker repositories from `/etc/yum.repos.d/`.

2.  Temporarily store the Docker EE repository URL you noted down in the
    [prerequisites](#prerequisites) in an environment variable.
    This will not persist when the current session ends.

    ```bash
    $ export DOCKERURL='<DOCKER-EE-URL>'
    ```

3.  Store your Docker EE repository URL in a `yum` variable in `/etc/yum/vars/`.
    This command relies on the variable you stored in the previous step.

    ```bash
    $ sudo -E sh -c 'echo "$DOCKERURL/{{ linux-dist-url-slug }}" > /etc/yum/vars/dockerurl'
    ```

    {% if linux-dist == "rhel" %}

    Store your OS version string in `/etc/yum/vars/dockerosversion`. Most users
    should use `7`, but you can also use the more specific minor version,
    starting from `7.2`.

    ```bash
    $ sudo sh -c 'echo "7" > /etc/yum/vars/dockerosversion'
    ```

    {% endif %}

4.  Install required packages. `yum-utils` provides the `yum-config-manager`
    utility, and `device-mapper-persistent-data` and `lvm2` are required by the
    `devicemapper` storage driver.

    ```bash
    $ sudo yum install -y yum-utils \
      device-mapper-persistent-data \
      lvm2
    ```

{% if linux-dist == "rhel" %}
5.  Enable the `extras` RHEL repository. This ensures access to the
    `container-selinux` package which is required by `docker-ee`.

    ```bash
    $ sudo yum-config-manager --enable rhel-7-server-extras-rpms
    ```

    Depending on cloud provider, you may also need to enable another repository.

    For AWS:

    ```bash
    $ sudo yum-config-manager --enable rhui-REGION-rhel-server-extras
    ```

    > **Note**: `REGION` here is literal, and does *not* represent the region
    > your machine is running in.

    For Azure:

    ```bash
    $ sudo yum-config-manager --enable rhui-rhel-7-server-rhui-extras-rpms
    ```
{% endif %}

6.  Use the following command to add the **stable** repository:

    ```bash
    $ sudo -E yum-config-manager \
        --add-repo \
        "$DOCKERURL/{{ linux-dist-url-slug }}/docker-ee.repo"
    ```


{% elsif section == "install-using-yum-repo" %}

1.  Install the latest version of Docker EE, or go to the next step to install a
    specific version.

    ```bash
    $ sudo yum -y install docker-ee
    ```

    If this is the first time you have refreshed the package index since adding
    the Docker repositories, you will be prompted to accept the GPG key, and
    the key's fingerprint will be shown. Verify that the fingerprint matches
    `{{ gpg-fingerprint }}` and if so, accept the key.

2.  On production systems, you should install a specific version of Docker EE
    instead of always using the latest. List the available versions.
    This example uses the `sort -r` command to sort the results by version
    number, highest to lowest, and is truncated.

    > **Note**: This `yum list` command only shows binary packages. To show
    > source packages as well, omit the `.x86_64` from the package name.

    ```bash
    $ sudo yum list docker-ee.x86_64  --showduplicates | sort -r

    docker-ee.x86_64         {{ minor-version }}.ee.2-1.el7.{{ linux-dist }}          docker-ee-stable-17.06
    ```

    The contents of the list depend upon which repositories you have enabled,
    and will be specific to your version of {{ linux-dist-long }}
    (indicated by the `.el7` suffix on the version, in this example). Choose a
    specific version to install. The second column is the version string. You
    can use the entire version string, but **you need to include at least to the
    first hyphen**. The third column is the repository name, which indicates
    which repository the package is from and by extension its stability level.
    To install a specific version, append the version string to the package name
    and separate them by a hyphen (`-`):

    > **Note**: The version string is the package name plus the version up to
    > the first hyphen. In the example above, the fully qualified package name
    > is `docker-ee-17.06.1.ee.2`.

    ```bash
    $ sudo yum -y install <FULLY-QUALIFIED-PACKAGE-NAME>
    ```

    Docker is installed but not started. The `docker` group is created, but no
    users are added to the group.

3.  Edit `/etc/docker/daemon.json`. If it does not yet exist, create it. Assuming
    that the file was empty, add the following contents.

    ```json
    {
      "storage-driver": "devicemapper"
    }
    ```

4.  For production systems, you must use `direct-lvm` mode, which requires you
    to prepare the block devices. Follow the procedure in the
    [devicemapper storage driver guide](/engine/userguide/storagedriver/device-mapper-driver.md#configure-direct-lvm-mode-for-production){: target="_blank" class="_" }
    **before starting Docker**.

5.  Start Docker.

    ```bash
    $ sudo systemctl start docker
    ```

6.  Verify that Docker EE is installed correctly by running the `hello-world`
    image.

    ```bash
    $ sudo docker run hello-world
    ```

    This command downloads a test image and runs it in a container. When the
    container runs, it prints an informational message and exits.

Docker EE is installed and running. You need to use `sudo` to run Docker
commands. Continue to [Linux postinstall](linux-postinstall.md) to allow
non-privileged users to run Docker commands and for other optional configuration
steps.



{% elsif section == "upgrade-using-yum-repo" %}

To upgrade Docker EE:

1.  If upgrading to a new major Docker EE version (such as when going from
    Docker 17.03.x to Docker 17.06.x),
    [add the new repository](#set-up-the-repository){: target="_blank" class="_" }.

2.  Run `sudo yum makecache fast`.

3.  Follow the
    [installation instructions](#install-docker), choosing the new version you
    want to install.


{% elsif section == "install-using-yum-package" %}

If you cannot use the official Docker repository to install Docker EE, you can
download the `.{{ package-format | downcase }}` file for your release and
install it manually. You will need to download a new file each time you want to
upgrade Docker EE.

{% if linux-dist == "rhel" %}
1.  Enable the `extras` RHEL repository. This ensures access to the
    `container-selinux` package which is required by `docker-ee`.

    ```bash
    $ sudo yum-config-manager --enable rhel-7-server-extras-rpms
    ```

    Alternately, obtain that package manually from Red Hat.
    There is no way to publicly browse this repository.
{% endif %}

1.  Go to the Docker EE repository URL associated with your
    trial or subscription in your browser. Go to
    `{{ linux-dist-url-slug }}/7/x86_64/stable-{{ minor-version }}/Packages` and
    download the `.{{ package-format | downcase }}` file for the Docker version
    you want to install.

    {% if linux-dist == "rhel" %}

    > **Note**: If you have trouble with `selinux` using the packages under the
    > `7` directory, try choosing the version-specific directory instead, such
    > as `7.3`.


    {% endif %}

2.  Install Docker EE, changing the path below to the path where you downloaded
    the Docker package.

    ```bash
    $ sudo yum install /path/to/package.rpm
    ```

    Docker is installed but not started. The `docker` group is created, but no
    users are added to the group.

3.  Edit `/etc/docker/daemon.json`. If it does not yet exist, create it.
    Assuming that the file was empty, add the following contents.

    ```json
    {
      "storage-driver": "devicemapper"
    }
    ```

4.  For production systems, you must use `direct-lvm` mode, which requires you
    to prepare the block devices. Follow the procedure in the
    [devicemapper storage driver guide](/engine/userguide/storagedriver/device-mapper-driver.md#configure-direct-lvm-mode-for-production){: target="_blank" class="_" }
    **before starting Docker**.

5.  Start Docker.

    ```bash
    $ sudo systemctl start docker
    ```

6.  Verify that Docker EE is installed correctly by running the `hello-world`
    image.

    ```bash
    $ sudo docker run hello-world
    ```

    This command downloads a test image and runs it in a container. When the
    container runs, it prints an informational message and exits.

Docker EE is installed and running. You need to use `sudo` to run Docker
commands. Continue to [Post-installation steps for Linux](linux-postinstall.md)
to allow non-privileged users to run Docker commands and for other optional
configuration steps.


{% elsif section == "upgrade-using-yum-package" %}

To upgrade Docker EE, download the newer package file and repeat the
[installation procedure](#install-from-a-package), using `yum -y upgrade`
instead of `yum -y install`, and pointing to the new file.

{% elsif section == "yum-uninstall" %}

1.  Uninstall the Docker EE package:

    ```bash
    $ sudo yum -y remove docker-ee
    ```

2.  Images, containers, volumes, or customized configuration files on your host
    are not automatically removed. To delete all images, containers, and
    volumes:

    ```bash
    $ sudo rm -rf /var/lib/docker
    ```

3.  If desired, remove the `devicemapper` thin pool and reformat the block
    devices that were part of it.

You must delete any edited configuration files manually.


{% elsif section == "linux-install-nextsteps" %}

- Continue to [Post-installation steps for Linux](/engine/installation/linux/linux-postinstall.md)

- Continue with the [User Guide](/engine/userguide/index.md).

{% endif %}
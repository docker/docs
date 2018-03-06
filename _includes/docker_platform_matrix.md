{% assign green-check = '![yes](/install/images/green-check.svg){: style="height: 14px; margin: 0 auto"}' %}
{% assign install-prefix-ce = '/install/linux/docker-ce' %}
{% assign install-prefix-ee = '/install/linux/docker-ee' %}

#### Docker EE

| Platform                                                        | x86_64 / amd64                                         | IBM Power (ppc64le)                                    | IBM Z (s390x)                                          |
|:----------------------------------------------------------------|:-------------------------------------------------------|:-------------------------------------------------------|:-------------------------------------------------------|
| [CentOS]({{ install-prefix-ee }}/centos.md)                     | [{{ green-check }}]({{ install-prefix-ee }}/centos.md) |                                                        |                                                        |
| [Oracle Linux]({{ install-prefix-ee }}/oracle.md)               | [{{ green-check }}]({{ install-prefix-ee }}/oracle.md) |                                                        |                                                        |
| [Red Hat Enterprise Linux]({{ install-prefix-ee }}/rhel.md)     | [{{ green-check }}]({{ install-prefix-ee }}/rhel.md)   | [{{ green-check }}]({{ install-prefix-ee }}/rhel.md)   | [{{ green-check }}]({{ install-prefix-ee }}/rhel.md)   |
| [SUSE Linux Enterprise Server]({{ install-prefix-ee }}/suse.md) | [{{ green-check }}]({{ install-prefix-ee }}/suse.md)   | [{{ green-check }}]({{ install-prefix-ee }}/suse.md)   | [{{ green-check }}]({{ install-prefix-ee }}/suse.md)   |
| [Ubuntu]({{ install-prefix-ee }}/ubuntu.md)                     | [{{ green-check }}]({{ install-prefix-ee }}/ubuntu.md) | [{{ green-check }}]({{ install-prefix-ee }}/ubuntu.md) | [{{ green-check }}]({{ install-prefix-ee }}/ubuntu.md) |
| [Microsoft Windows Server 2016](/install/windows/docker-ee.md)  | [{{ green-check }}](/install/windows/docker-ee.md)     |                                                        |                                                        |

> Limitations for Docker EE on IBM Power architecture
>
> - Neither UCP managers nor workers are supported on IBM Power.

#### Docker CE

| Platform                                    | x86_64 / amd64                                         | ARM                                                    | ARM64 / AARCH64                                        | IBM Power (ppc64le)                                    | IBM Z (s390x)                                          |
|:--------------------------------------------|:-------------------------------------------------------|:-------------------------------------------------------|:-------------------------------------------------------|:-------------------------------------------------------|:-------------------------------------------------------|
| [CentOS]({{ install-prefix-ce }}/centos.md) | [{{ green-check }}]({{ install-prefix-ce }}/centos.md) |                                                        | [{{ green-check }}]({{ install-prefix-ce }}/centos.md) |                                                        |                                                        |
| [Debian]({{ install-prefix-ce }}/debian.md) | [{{ green-check }}]({{ install-prefix-ce }}/debian.md) | [{{ green-check }}]({{ install-prefix-ce }}/debian.md) | [{{ green-check }}]({{ install-prefix-ce }}/debian.md) |                                                        |                                                        |
| [Fedora]({{ install-prefix-ce }}/fedora.md) | [{{ green-check }}]({{ install-prefix-ce }}/fedora.md) |                                                        |                                                        |                                                        |                                                        |
| [Ubuntu]({{ install-prefix-ce }}/ubuntu.md) | [{{ green-check }}]({{ install-prefix-ce }}/ubuntu.md) | [{{ green-check }}]({{ install-prefix-ce }}/ubuntu.md) | [{{ green-check }}]({{ install-prefix-ce }}/ubuntu.md) | [{{ green-check }}]({{ install-prefix-ce }}/ubuntu.md) | [{{ green-check }}]({{ install-prefix-ce }}/ubuntu.md) |

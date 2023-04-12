> Download Docker Desktop
>
{% if include.all or include.win -%}
> [Windows](https://desktop.docker.com/win/main/amd64{{ include.build_path }}Docker%20Desktop%20Installer.exe) ([checksum](https://desktop.docker.com/win/main/amd64{{ include.build_path }}checksums.txt){: target="_blank" rel="noopener" class="_"}) {% if include.all or include.mac or include.linux %} | {% endif %}
{% endif -%}
{% if include.all or include.mac -%}
> [Mac with Intel chip](https://desktop.docker.com/mac/main/amd64{{ include.build_path }}Docker.dmg) ([checksum](https://desktop.docker.com/mac/main/amd64{{ include.build_path }}checksums.txt){: target="_blank" rel="noopener" class="_"}) |
> [Mac with Apple chip](https://desktop.docker.com/mac/main/arm64{{ include.build_path }}Docker.dmg) ([checksum](https://desktop.docker.com/mac/main/arm64{{ include.build_path }}checksums.txt){: target="_blank" rel="noopener" class="_"}) {% if include.all or include.linux %} | {% endif %}
{% endif -%}
{% if include.all or include.linux -%}
> [Debian](https://desktop.docker.com/linux/main/amd64{{ include.build_path }}docker-desktop-4.17.0-amd64.deb) - 
> [RPM](https://desktop.docker.com/linux/main/amd64{{ include.build_path }}docker-desktop-4.17.0-x86_64.rpm) - 
> [Arch package](https://desktop.docker.com/linux/main/amd64{{ include.build_path }}docker-desktop-4.17.0-x86_64.pkg.tar.zst) ([checksum](https://desktop.docker.com/linux/main/amd64{{ include.build_path }}checksums.txt){: target="_blank" rel="noopener" class="_"}) 
{% endif -%}
{% if include.build_path == "/" -%}
{: .tip}
{% endif -%}
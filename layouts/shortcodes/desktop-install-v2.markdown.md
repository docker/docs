{{- $all := .Get "all" -}}
{{- $win := .Get "win" -}}
{{- $win_arm_release := .Get "win_arm_release" -}}
{{- $mac := .Get "mac" -}}
{{- $linux := .Get "linux" -}}
{{- $build_path := .Get "build_path" -}}
Download Docker Desktop:
{{- if or $all $win }}
- [Windows](https://desktop.docker.com/win/main/amd64{{ $build_path }}Docker%20Desktop%20Installer.exe)
{{- end }}
{{- if $win_arm_release }}
- [Windows ARM {{ $win_arm_release }}](https://desktop.docker.com/win/main/arm64{{ $build_path }}Docker%20Desktop%20Installer.exe)
{{- end }}
{{- if or $all $mac }}
- [Mac (Apple chip)](https://desktop.docker.com/mac/main/arm64{{ $build_path }}Docker.dmg)
- [Mac (Intel chip)](https://desktop.docker.com/mac/main/amd64{{ $build_path }}Docker.dmg)
{{- end }}
{{- if or $all $linux }}
- [Linux (Debian)](https://desktop.docker.com/linux/main/amd64{{ $build_path }}docker-desktop-amd64.deb)
- [Linux (RPM)](https://desktop.docker.com/linux/main/amd64{{ $build_path }}docker-desktop-x86_64.rpm)
- [Linux (Arch)](https://desktop.docker.com/linux/main/amd64{{ $build_path }}docker-desktop-x86_64.pkg.tar.zst)
{{- end }}

# Requires https://github.com/rumpl/vhs/tree/screen-settled
# and brew install ffmpeg ttyd

Output docs/demo.gif

Set Shell "zsh"
Set FontSize 24
Set Width 2000
Set Height 1600
Set PlaybackSpeed 2
Set Padding 0
Set WindowBar Colorful
Set BorderRadius 10

Type "docker agent run ./examples/gopher.yaml" Sleep 3s Enter
Sleep 1s

Type "Summarize README.md and explain what this project does" Sleep 1s Enter

Sleep 1s WaitScreenSettled 1s Sleep 1s

Sleep 20s
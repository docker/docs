# Requires https://github.com/rumpl/vhs/tree/screen-settled
# and brew install ffmpeg ttyd

Output pirate.gif

Set Shell "zsh"
Set FontSize 24
Set Width 1800
Set Height 1000
Set PlaybackSpeed 2
Set Padding 0
Set WindowBar Colorful
Set BorderRadius 10

Type "docker agent run ./examples/pirate.yaml" Sleep 3s Enter
Sleep 1s

Type "Hello, my name is Jean-Laurent! How are you today?" Sleep 1s Enter

Sleep 1s WaitScreenSettled 1s Sleep 1s

Type "That's impressive pirate speaking! Can you teach me a couple of slang words?" Sleep 1s Enter

Sleep 1s WaitScreenSettled 1s Sleep 1s

Sleep 20s
# https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/customizing-your-repositorys-social-media-preview
Output demo-social.gif

Require git
Require scmpuff
Require ffmpeg
Set Shell bash

Set Theme "Tomorrow Night"

Set Width 1280
Set Height 640
Set Padding 10
Set FontSize 18

Set BorderRadius 10
Set Margin 80
Set MarginFill "#674EFF"
Set BorderRadius 10

Hide
Type "bash setup.sh"
Enter
Type "cd /tmp/scmpuff-demo"
Enter
Type `eval "$(scmpuff init -s)"`
Enter

# Some off-screen setup for very simple script
Type "echo foo > README.md"
Enter
Type "git mv chemex.go aeropress.go"
Enter
Type "clear"
Enter
Show

# Script: very simple version
Type "gs"
Sleep 0.25
Enter
Sleep 2

Type "git add "
Sleep 0.25
Type "2"
Sleep 0.25
Type " "
Sleep 0.5
Type "4-6"
Sleep 0.5
Enter
Sleep 3

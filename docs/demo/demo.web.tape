Output demo.webm
Output demo.mp4

Require git
Require scmpuff
Require ffmpeg
Set Shell bash
Set Theme "Builtin Solarized Dark"
# Env PS1 "✨ "
# Set FontSize 28
Set Width 1200
Set Height 600
Set Padding 10

Hide
Type "bash setup.sh"
Enter
Type "cd /tmp/scmpuff-demo"
Enter
Type `eval "$(scmpuff init -s)"`
Enter
Type "clear"
Enter
Show

# Script from website movie demo
Type "ls"
Sleep 0.25
Enter
Sleep 1

Type "gs"
Sleep 0.25
Enter
Sleep 1

Type "echo foo > README.md"
Enter
Type "git mv chemex.go aeropress.go"
Enter
Type "gs"
Enter
Sleep 2

Type "git add 2 4-6"
Sleep 0.25
Enter
Sleep 2

Type "git reset 1 3"
Sleep 0.1
Enter
Type "gs"
Sleep 0.25
Enter
Sleep 2

Type "git restore 4"
Sleep 0.1
Enter
Type "gs"
Sleep 0.25
Enter
Sleep 2

Type "git reset --hard HEAD && clear"
Enter
Sleep 2

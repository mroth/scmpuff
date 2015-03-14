# mock a sample git dir in a good state for taking captures, screenshots, etc.
desc "build a mocked sample git dir"
task :sample do
  DIR = "tmp/scmpuff"
  rm_rf DIR
  mkdir_p DIR
  Dir.chdir DIR do
    File.write(".init.sh", "#!/bin/sh\nexport PS1='✨  '\n    \
                            export PATH=~/src/go/bin:$PATH\n \
                            eval \"$(scmpuff init -s)\"")
    sh "git init"
    sh "touch presspot.go chemex.go README.md"
    sh "git add ."
    sh "git commit -m."
    sh "touch espresso.md americano.md cappucino.md macchiato.md"
    #sh "seq 1 10 > README.md"
    #sh "git mv chemex.go aeropress.go"

    # current looping demo script:
    # ls
    # gs
    # echo foo > README.md
    # git mv chemex.go aeropress.go
    # gs
    # ga 2 4-6
    # grs 1 3
    # gs
    # gco 4
    # gs
    # git reset --hard HEAD && clear

    # note to future self, these were the conversion factors I used:
    # ffmpeg -i scmpuff_demo.mov -r 10 -vf "setpts=0.6*PTS,crop=x=5:w=iw-11:y=3:h=ih-35" scmpuffdemo-2x.mp4 # (also .webm)
    # ffmpeg -i scmpuff_demo.mov -r 10 -vf "setpts=0.6*PTS,crop=x=5:w=iw-11:y=3:h=ih-35,scale=iw/2:ih/2" scmpuffdemo-1x.mp4
    # ffmpeg -i scmpuff_demo.mov -r 10 -vf "setpts=0.6*PTS,crop=x=5:w=iw-11:y=3:h=ih-35,scale=iw/2:ih/2" -pix_fmt rgb24 -f gif - | gifsicle --optimize=3 --delay=6 > output_1x.gif
    # ffmpeg -i scmpuff_demo.mov -r 10 -vf "setpts=0.6*PTS,crop=x=5:w=iw-11:y=3:h=ih-35" -f image2pipe -vcodec ppm - | convert -resize 50% -layers Optimize - gif:- | gifsicle --loop --optimize=3 --multifile > output.gif
  end
end

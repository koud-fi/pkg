#!/bin/sh
thumbnail test.jpg test.png 500
thumbnail test.jpg test.tiff 500
thumbnail test.jpg test.webp 1000

ffmpeg -i test.mov -c:v libx264 -c:a copy -crf 32 test.mp4
ffmpeg -i test.mov -c:v libvpx-vp9 -crf 32 test.webp
ffmpeg -i test.mov -vf "fps=5,scale=120:-1:flags=lanczos,split[s0][s1];[s0]palettegen[p];[s1][p]paletteuse" test.gif

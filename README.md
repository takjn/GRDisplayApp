# Gadget Renesas DisplayApp for Linux
DisplayApp is an application that displays the camera image sent by USB CDC.

## How to build
run `make build` command. `displayapp` will be build.

## How to use
To record the video as MP4 file, run the following command. `/tmp/out.mp4` will be created.
```
./displayapp | ffmpeg -y -f image2pipe -c:v mjpeg -r 30 -i - -vcodec libx264 /tmp/out.mp4
```

To display the video in the window, run the following command.
```
./displayapp | gst-launch-1.0 fdsrc fd=0 ! decodebin ! progressreport update-freq=1 ! videoscale ! video/x-raw,width=640,height=480 ! queue ! videoconvert ! videorate ! video/x-raw, framerate=30/1 ! ximagesink
```

To output jpeg files every 1 seconds, run the following command. `/tmp/test*.jpg` will be created.
```
./displayapp -p /tmp
```

**NOTE:**
Please change the width, height and framerate if needed. 

## Reference
- [GR-LYCHEE](https://www.renesas.com/us/ja/products/gadget-renesas/boards/gr-lychee.html)
- [Camera and LCD sample for GR-Boards](https://github.com/d-kato/GR-Boards_Camera_LCD_sample)
- [DisplayApp](https://os.mbed.com/users/dkato/code/DisplayApp/)

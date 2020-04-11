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
./displayapp | ffplay -i -
```

To output jpeg files every 1 seconds, run the following command. `/tmp/test*.jpg` will be created.
```
./displayapp -p /tmp
```

To use [v4l2loopback](https://github.com/umlaeute/v4l2loopback), run the following commands.
First start v4l2loopback with `modprobe` command. Then confirm your loopback device with `v4l2-ctl` command.

```
sudo modprobe v4l2loopback exclusive_caps=1
v4l2-ctl --list-devices
```

In my case, the loopback device was `/dev/video1`.
In the case, run the following command.
```
./displayapp | ffmpeg -y -f image2pipe -c:v mjpeg -r 30 -i - -pix_fmt yuyv422 -f v4l2 /dev/video1
```

**NOTE:**
Please change the width, height and framerate if needed. 

## Reference
- [GR-LYCHEE](https://www.renesas.com/us/ja/products/gadget-renesas/boards/gr-lychee.html)
- [Camera and LCD sample for GR-Boards](https://github.com/d-kato/GR-Boards_Camera_LCD_sample)
- [DisplayApp](https://os.mbed.com/users/dkato/code/DisplayApp/)

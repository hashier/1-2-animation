# What is 1-2-animation

This command line tool generates Poemotion images. I actually don't know the real name of those images but you can see those images here: [youtube video](https://www.youtube.com/watch?v=Serhd00QNzo)

# How to install it

```
$ get -u github.com/hashier/1-2-animation
```

# How to use it

```
$ 1-2-animation
Usage of 1-2-animation:
  -calibrate string
    	Write calibration image to <file>. This image can be used to determine how many pixels the "slit" is wide
  -cp
    	write CPU profile to cpu.prof
  -example
    	Write some demo color striped images to the drive. Those can be helpful to determine to "number of frames"
  -f int
    	How many input frames (ignored if images provide, only needed for example images) (default 5)
  -h int
    	Height of image, only applies to calibration and example images (default 1050)
  -mp
    	write memory profile to mem.prof
  -o string
    	output of generated PNG image (default "out.png")
  -ppf int
    	How many pixel to take from the input image for every frame (how wide is the "slit" of the mask (default 2)
  -w int
    	Width of image, only applies to calibration and example images (default 1680)
2018/08/27 08:14:36 Error: You need to either provide at least 2 input image OR enable calibration mode OR example mode
```

You can simply generate some colored test images via the `-example` flag if you just wanna see what this is about or provide some input images that you want to animate.

Make sure when you display them that you display them in native size and that your image viewer does not enlarge/shrink them to fit your display.

For best results I recommend to figure out the amount of frames and the width of the slit your sheet has. For this you can use the calibration image you can get with `-calibrate <file>` (for the width of the slit) and use the example images to determine how many frames there are on your sheet.

# Example image

![5 framed colored test image](https://github.com/hashier/1-2-animation/blob/master/example/example-color-5-out.png?raw=true)
![7 framed colored test image](https://raw.githubusercontent.com/hashier/1-2-animation/master/example/example-color-7-out.png)
![5 framed rotating molecule](https://raw.githubusercontent.com/hashier/1-2-animation/master/example/molecule/molecule.png)

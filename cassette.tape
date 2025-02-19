Output demo.gif

# Set up a 1200x600 terminal with 46px font.
Set FontSize 46
Set Width 1600
Set Height 800

#Type "# List source image and sizes"
#Enter
#Type `ls /home/pawndev/Pictures/samples/ | xargs -d '\n' -L1 identify -format "%f: %wx%h\n"`
#Enter
#Sleep 2s

Type "go run ./cmd/resizeboy"
Enter
Sleep 2s

# Input dir
Type "/home/pawndev/Pictures/samples/"
Sleep 1s
Enter

# Output dir
Type "/home/pawndev/Pictures/resized/"
Sleep 1s
Enter

# Select PNG
Sleep 1s
Enter

# Choose suffix
Sleep 1s
Left
Enter

# Write suffix
Sleep 1s
Type "gba"
Enter

# Max width
Sleep 1s
Type "300"
Enter

Sleep 3s

#Type "# Converted/resized images"
#Enter
#Type `ls /home/pawndev/Pictures/resized/ | xargs -d '\n' -L1 identify -format "%f: %wx%h\n"`
#Enter

#Sleep 2s
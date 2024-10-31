# IMAGE ZIP
### It's a tool to zip images in a folder. That's about it.
This tool can be used cross platform, just build it using the packages contained within.
If you want to register it, just use `go install` inside the files instead of `go build`.

#### Requires GoLang v1.21 or above to build!

Use couldn't be simpler, just run:

`zip_imagefolder -dir "/some/dir/with/images" -out "my_folder.zip"`

or, if you use windows,

`./zip_imagefolder.exe -dir "X:\some\dir\with\images" -out "my_folder.zip"`

I included some python stuff to run jobs async, as a treat.

The easiest way to use it, is as follows:

```py
import asyncio
from pyimzip import ImageZipProcess, imzip_set_executable

# make sure the executable is able to be located.
imzip_set_executable("./zip_imagefolder.exe")

async def main():
    zip = ImageZipProcess(
        R"..\SomeShotsFolder\Render01",
        "Render01_archive.zip",
    )
    
    # helper function to poll zipping process every 2 seconds.
    async def poll_zip_task(zip: ImageZipProcess):
        while not zip.done:
            if zip.progress:
                print("progress:", zip.progress)
            await asyncio.sleep(2)
    
    # Set up tasks
    waiting = asyncio.Task(zip.start_zipping())
    printing = asyncio.Task(poll_zip_task(zip))
    
    result = await asyncio.gather(waiting, printing)
    # output from waiting task (zip.start_zipping)
    # is the folder "..\SomeShotsFolder\Render01\Render01_archive.zip"
    folder = result[0]
    print(folder)

if __name__ == "__main__":
    asyncio.run(main())
```

both TU's contain demo files.

This tool was made to facilitate one step in the automation of getting image files to students @ HKU. 

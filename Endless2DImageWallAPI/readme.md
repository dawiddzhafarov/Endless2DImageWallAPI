Put some images into 'images' directory  
go build -o server main.go images.go matrix_generator.go  
Available endpoints
* /image=<imageName> - returns an image basend on encoded name  
* /imagesMatrix?x=<x>?y=<y>?z=<z> - returns structure consisting of image names and their location based on specified coordinates  
* /generateMatrix?z=<z> - generates and creates a JSON file of images and their positions, based on provided scale  

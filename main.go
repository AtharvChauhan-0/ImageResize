package main

import (
	"bytes"
	"cloud.google.com/go/storage"
	_ "cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/nfnt/resize"
	_ "github.com/nfnt/resize"
	"image"
	_ "image"
	"image/jpeg"
	_ "image/jpeg"
	"io"
	_ "io/ioutil"
	"log"
	_ "log"
	"net/http"
	_ "net/http"
	"os"
	_ "os"
)

//TIP To run your code, right-click the code and select <b>Run</b>. Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.

func main() {

}
// func HTTPHandler(w http.ResponseWriter, r *http.Request) {
// 	//var data struct {
// 	//}
// 	//ctx := context.Background()
// 	//client, err := storage.NewClient(ctx)
// 	//if err != nil {
// 	//	panic(err)
// 	//}
// 	//var objname []string
// 	//bktdown := client.Bucket("")
// 	//bktup := client.Bucket("")
// 	//it := bktdown.Objects(ctx, nil)
// 	//for {
// 	//	attrs, err := it.Next()
// 	//	if err != nil {
// 	//		log.Fatal(err)
// 	//	}
// 	//	objname = append(objname, attrs.Name)
// 	//	fmt.Println(objname)
// 	//}

// }
func EventProcessor(ctx context.Context, e event.Event) {
	var data StorageObjectData
	if err := e.DataAs(&data); err != nil {
		fmt.Errorf("event.DataAs: %v", err)
	}
	imgData, err := downloadImage(ctx, data.Bucket, data.Name)
	if err != nil {
		fmt.Errorf("failed to downloadImage: %v", err)
	}
	img, err := decodeImage(imgData)
	m := resize.Resize(1000, 0, img, resize.Lanczos3)
	out, err := os.Create("test_resized.jpg")
	if err != nil {
		log.Fatal(err)
	}
	err = jpeg.Encode(out, m, nil)
	if err != nil {
		return
	}
	if err := uploadImageToBucket(ctx, "result_img", "Result", "test_resized.jpg"); err != nil {
		log.Fatalf("Failed to upload image: %v", err)
	} else {
		fmt.Println("Image uploaded successfully.")
	}
}

func downloadImage(ctx context.Context, bucket, name string) ([]byte, error) {
	// Create a new Cloud Storage client
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage client: %v", err)
	}
	defer client.Close()

	// Get a reference to the bucket and object
	bucketHandle := client.Bucket(bucket)
	objectHandle := bucketHandle.Object(name)

	// Read the object's data
	reader, err := objectHandle.NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create object reader: %v", err)
	}
	defer reader.Close()

	// Read the data from the reader
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read object data: %v", err)
	}

	return data, nil // Return the image data
}
func decodeImage(data []byte) (image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	return img, err
}
func uploadImageToBucket(ctx context.Context, bucketName, objectName, filePath string) error {
	// Create a new client
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	// Open the image file for reading
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("os.Open: %v", err)
	}
	defer file.Close()

	// Get the bucket and the object handle
	bucket := client.Bucket(bucketName)
	obj := bucket.Object(objectName)

	// Create a writer for the object in the bucket
	w := obj.NewWriter(ctx)

	// Copy the image file data to the cloud storage object
	if _, err := io.Copy(w, file); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}

	// Close the writer to complete the upload
	if err := w.Close(); err != nil {
		return fmt.Errorf("writer.Close: %v", err)
	}

	fmt.Printf("Image file %s uploaded to bucket %s as %s\n", filePath, bucketName, objectName)
	return nil
}

type StorageObjectData struct {
	Bucket string `json:"bucket,omitempty"`
	Name   string `json:"name,omitempty"`
}

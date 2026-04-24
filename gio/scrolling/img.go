package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	_ "image/jpeg"
	_ "image/png"
)

func img(name, addr string) (image.Image, error) {
	file, err := os.Open(name)
	if err == nil {
		defer file.Close()
		img, _, err := image.Decode(file)
		return img, err
	}
	neturl, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}
	if neturl.Scheme != "https" {
		return nil, fmt.Errorf("Unsupported protocl %q\n", neturl.Scheme)
	}
	log.Printf("Downloading image from %q\n", neturl.Hostname())
	response, err := http.Get(addr)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	img, _, err := image.Decode(response.Body)
	if err != nil {
		return nil, err
	}
	name = strings.TrimSuffix(name, path.Ext(name)) + ".jpg"
	file, err = os.Create(name)
	if err != nil {
		log.Printf("Failed to create %q\n%v\n", name, err)
		return img, nil
	}
	err = jpeg.Encode(file, img, &jpeg.Options{Quality: jpeg.DefaultQuality})
	file.Close()
	if err != nil {
		log.Printf("Failed to write %q\n%v\nDeleting...\n", name, err)
		err = os.Remove(name)
		if err != nil {
			log.Printf("Failed to remove %q\n%v\n", name, err)
		}
		return img, nil
	}
	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("Saved %q\n", name)
	} else {
		log.Printf("Saved %q to %q\n", name, cwd)
	}
	return img, nil
}

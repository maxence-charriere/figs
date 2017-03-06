package main

import (
	"os"
	"testing"
)

func TestDownloadGif(t *testing.T) {
	url := "http://static1.comicvine.com/uploads/original/11124/111243429/5155933-2897910123-b3663.gif"
	img, err := downloadGif(url)
	if err != nil {
		t.Fatal(err)
	}
	img.Close()
}

func TestPingpongGif(t *testing.T) {
	f, err := os.Open("test.gif")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	img, err := pinpongGif(f)
	if err != nil {
		t.Error(err)
	}

	if err = saveGif(img, "testpingpong.gif"); err != nil {
		t.Error(err)
	}
}

func TestStartGifJob(t *testing.T) {
	completionChan := make(chan bool)
	defer close(completionChan)

	url := "http://img.pandawhale.com/post-28921-I-may-have-deserved-that-gif-s-qyCQ.gif"
	go startGifJob(url, "testjob.gif", completionChan)

	if ok := <-completionChan; !ok {
		t.Error("job failed")
	}
}

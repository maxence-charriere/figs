package main

import (
	"fmt"
	"image/gif"
	"io"
	"net/http"
	"os"
)

func main() {
	endpoints := []string{
		"http://static1.comicvine.com/uploads/original/11124/111243429/5155933-2897910123-b3663.gif",
		"https://media.tenor.co/images/76708e6143b45195612c69086358c138/raw",
		"http://gifrific.com/wp-content/uploads/2012/09/Tony-Stark-Sunglasses-Iron-Man-2.gif",
		"http://49.media.tumblr.com/6868d67722f6ec5a071845b562e5cb38/tumblr_o18rkrn3XO1rp84c6o1_400.gif",
		"https://s-media-cache-ak0.pinimg.com/originals/6f/45/3b/6f453bcd0b53c32eb296cb37eb9e5e9c.gif",
		"https://media.giphy.com/media/iEoh3SQbYLu2A/giphy.gif",
		"https://media.giphy.com/media/zZeCRfPyXi9UI/giphy.gif",
		"http://4.bp.blogspot.com/-Be1d15E-LeQ/Ug91Db64JzI/AAAAAAAACfE/imTnA9SfPVs/s1600/Lord+of+War+-+Acting+Madness.gif",
		"http://img.pandawhale.com/post-28921-I-may-have-deserved-that-gif-s-qyCQ.gif",
		"http://i1.kym-cdn.com/entries/icons/original/000/017/454/tumblr_myz4pvk1nF1schu5jo1_500.gif",
	}

	completionChan := make(chan bool)
	defer close(completionChan)

	for i, e := range endpoints {
		name := fmt.Sprintf("dlgift_%v.gif", i)
		go startGifJob(e, name, completionChan)
	}

	jobCount := len(endpoints)
	success := 0
	for ok := range completionChan {
		if ok {
			success++
		}

		jobCount--
		if jobCount <= 0 {
			break
		}
	}

	fmt.Printf("\n%v/%v jobs succeed!\n", success, len(endpoints))
}

func startGifJob(url string, saveName string, completionChan chan<- bool) {
	fmt.Printf("%v ~> \033[33mstarting job\033[00m\n", url)

	r, err := downloadGif(url)
	if err != nil {
		fmt.Printf("%v ~> \033[91mdownload failed: %v\033[00m\n", url, err)
		completionChan <- false
		return
	}
	defer r.Close()

	img, err := pinpongGif(r)
	if err != nil {
		fmt.Printf("%v ~> \033[91mpinpongify failed: %v\033[00m\n", url, err)
		completionChan <- false
		return
	}

	if err = saveGif(img, saveName); err != nil {
		fmt.Printf("%v ~> \033[91msave failed: %v\033[00m\n", url, err)
		completionChan <- false
		return
	}

	fmt.Printf("%v ~> \033[92mjob success\033[00m\n", url)
	completionChan <- true
}

func downloadGif(url string) (img io.ReadCloser, err error) {
	res, err := http.Get(url)
	if err != nil {
		return
	}

	img = res.Body
	return
}

func pinpongGif(r io.Reader) (img *gif.GIF, err error) {
	if img, err = gif.DecodeAll(r); err != nil {
		return
	}

	// - 2 in order to avoid redrawing the last image.
	for i := len(img.Image) - 2; i >= 0; i-- {
		img.Image = append(img.Image, img.Image[i])
		img.Delay = append(img.Delay, img.Delay[i])
		img.Disposal = append(img.Disposal, img.Disposal[i])
	}

	img.LoopCount = 0
	return
}

func saveGif(img *gif.GIF, name string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()

	return gif.EncodeAll(f, img)
}

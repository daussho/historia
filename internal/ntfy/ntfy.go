package ntfy

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	mu = &sync.Mutex{}
)

func Send(title, msg string, tags []string) error {
	mu.Lock()
	defer func() {
		time.Sleep(1 * time.Second)
		mu.Unlock()
	}()

	topic := os.Getenv("NTFY_TOPIC")

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("https://ntfy.sh/%s", topic),
		strings.NewReader(msg),
	)
	if err != nil {
		log.Println("ntfy.Send error:", err)
		return err
	}

	req.Header.Set("Title", title)
	req.Header.Set("Tags", strings.Join(tags, ","))

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Println("ntfy.Send error:", err)
		return err
	}

	defer resp.Body.Close()

	return nil
}

func SendError(title, msg string) error {
	return Send(title, msg, []string{"rotating_light"})
}

package hip

import (
	"github.com/tbruyelle/hipchat-go/hipchat"

	"bytes"
	"github.com/pkg/errors"
	"log"
	"strings"
	"sync"
	"time"
)

type HipchatRoomPoller struct {
	cli      *hipchat.Client
	roomId   string
	lastSeen string
}

type HipchatRoomReader struct {
	sync.Mutex
	buffer bytes.Buffer
}

func (r *HipchatRoomReader) Read(p []byte) (n int, err error) {
	r.Lock()
	defer r.Unlock()
	n, _ = r.buffer.Read(p) //ignoring EOF of the underlying buffer as the messages might still come via hipchat polling
	return
}

func (r *HipchatRoomReader) Add(s string) {
	r.Lock()
	defer r.Unlock()
	r.buffer.WriteString(s + "\n")
}

func NewHipchatRoomReader(hipchatToken string, roomId string, messagePrefix string) (*HipchatRoomReader, error) {
	client := hipchat.NewClient(hipchatToken)

	history, response, err := client.Room.Latest(roomId, &hipchat.LatestHistoryOptions{MaxResults: 1})
	err = checkReponseCode(response, err)
	if err != nil {
		return nil, errors.Wrap(err, "hipchat polling err")
	}

	var lastMessageSeen string
	if len(history.Items) > 0 {
		lastMessageSeen = history.Items[0].ID
	}

	reader := &HipchatRoomReader{}

	go func() {
		for {
			time.Sleep(5 * time.Second)
			history, response, err := client.Room.Latest(roomId, &hipchat.LatestHistoryOptions{NotBefore: lastMessageSeen})
			if err = checkReponseCode(response, err); nil != err {
				log.Println("hipchat polling err", err)
				continue
			}
			for _, h := range history.Items {
				if h.ID == lastMessageSeen {
					continue
				}
				lastMessageSeen = h.ID
				if strings.HasPrefix(h.Message, messagePrefix) {
					reader.Add(h.Message[len(messagePrefix):])
				}
			}

		}
	}()

	return reader, nil
}

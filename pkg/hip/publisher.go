package hip

import (
	"github.com/tbruyelle/hipchat-go/hipchat"

	"github.com/pkg/errors"
	"log"
	"strings"
	"sync"
)

type HipchatRoomPrinter struct {
	cli    *hipchat.Client
	roomId string
	wg     sync.WaitGroup
}

func (s *HipchatRoomPrinter) Out(line string) {
	s.throwNotice(&hipchat.NotificationRequest{Notify: false, Message: line, Color: hipchat.ColorGray})
}

func (s *HipchatRoomPrinter) Err(line string) {
	s.throwNotice(&hipchat.NotificationRequest{Notify: true, Message: line, Color: hipchat.ColorRed})
}

func (s *HipchatRoomPrinter) Done(err error) {
	notification := hipchat.NotificationRequest{Notify: true, Message: "done", Color: hipchat.ColorGreen}
	if err != nil {
		notification.Color = hipchat.ColorRed
		notification.Message = "done with error: " + err.Error()
	}
	s.throwNotice(&notification)
	s.wg.Wait()
}

func (s *HipchatRoomPrinter) throwNotice(r *hipchat.NotificationRequest) {
	if "" == strings.TrimSpace(r.Message) {
		return //no point of empty notices
	}
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		err := checkReponseCode(s.cli.Room.Notification(s.roomId, r))
		if err != nil {
			log.Println("Hipchat notification error", err)
		}
	}()
}

//NB: isn't thread safe as does the hipchat authtest which affects all hipchat clients
func NewHipchatRoomPrinter(hipchatToken string, roomId string) (*HipchatRoomPrinter, error) {
	r := HipchatRoomPrinter{cli: hipchat.NewClient(hipchatToken), roomId: roomId}
	hipchat.AuthTest = true //<- this not thread safe
	err := checkReponseCode(r.cli.Room.Notification(roomId, &hipchat.NotificationRequest{Message: "auth test"}))

	hipchat.AuthTest = false
	if err != nil {
		return nil, errors.Wrap(err, "Failed auth check against the room")
	} else {
		return &r, nil
	}
}

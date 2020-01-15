package main

import (
	"fmt"
	"github.com/Rhymen/go-whatsapp"
	"log"
	"time"
	"strconv"
	"strings"
	"os"
)

// historyHandler for acquiring chat history
type historyHandler struct {
	c        *whatsapp.Conn
	messages []string
}

func (h *historyHandler) ShouldCallSynchronously() bool {
	return true
}

// handles and accumulates history's text messages.
// To handle images/documents/videos add corresponding handle functions
func (h *historyHandler) HandleTextMessage(message whatsapp.TextMessage) {
	authorID := "-"
	screenName := "-"
	if message.Info.FromMe {
		authorID = h.c.Info.Wid
		screenName = ""
	} else {
		if message.Info.Source.Participant != nil {
			authorID = *message.Info.Source.Participant
		} else {
			authorID = message.Info.RemoteJid
		}
		if message.Info.Source.PushName != nil {
			screenName = *message.Info.Source.PushName
		}
	}

	date := time.Unix(int64(message.Info.Timestamp), 0)
	h.messages = append(h.messages, fmt.Sprintf("%s	%s (%s): %s", date,
		authorID, screenName, message.Text))
}

func (h *historyHandler) HandleImageMessage(message whatsapp.ImageMessage) {
	authorID := "-"
	screenName := "-"
	if message.Info.FromMe {
		authorID = h.c.Info.Wid
		screenName = ""
	} else {
		if message.Info.Source.Participant != nil {
			authorID = *message.Info.Source.Participant
		} else {
			authorID = message.Info.RemoteJid
		}
		if message.Info.Source.PushName != nil {
			screenName = *message.Info.Source.PushName
		}
	}

	date := time.Unix(int64(message.Info.Timestamp), 0)

    data, err := message.Download();
    if err != nil {
        if err != whatsapp.ErrMediaDownloadFailedWith410 && err != whatsapp.ErrMediaDownloadFailedWith404 {
            return
        }
        _, err := h.c.LoadMediaInfo(message.Info.RemoteJid, message.Info.Id, strconv.FormatBool(message.Info.FromMe));
        if err != nil {
           return
        } else {
            data, err = message.Download()
            if err != nil {
                return
            }
        }
    }

    filename := fmt.Sprintf("/Users/oriol/Downloads/images/%v.%v", message.Info.Id, strings.Split(message.Type, "/")[1])
    file, err := os.Create(filename)
    defer file.Close()
    if err != nil {
        return
    }
    _, err = file.Write(data)
    if err != nil {
        return
    }
    h.messages = append(h.messages, fmt.Sprintf("%s	%s (%s): image received, saved at %v", date,
    		authorID, screenName, filename))

}

func (h *historyHandler) HandleError(err error) {
	log.Printf("Error occured while retrieving chat history: %s", err)
}

func GetHistory(jid string, wac *whatsapp.Conn) []string {
	// create out history handler
	handler := &historyHandler{c: wac}

	// load chat history and pass messages to the history handler to accumulate
	wac.LoadFullChatHistory(jid, 300, time.Millisecond*300, handler)
	return handler.messages
}

func GetAnyHistory(wac *whatsapp.Conn, chats map[string]struct{}) {
	// show list of chats
	var chatSlice []string
	for chat := range chats {
		chatSlice = append(chatSlice, chat)
		fmt.Printf("%d	%s\n", len(chatSlice), chat)
	}

	fmt.Println("Select chat number to get history for:")
	var index = 0
	for index < 1 || index > len(chatSlice) {
		fmt.Scanf("%d", &index)
	}

	// get history for the selected chat
	fmt.Println("Gathering chat history...")
	messages := GetHistory(chatSlice[index-1], wac)
	for _, message := range messages {
		fmt.Println(message)
	}
	fmt.Println(len(messages))
}

package usecase

import (
	"log"
	"talkspace-api/modules/consultation/model"
	doctor "talkspace-api/modules/doctor/model"
	user "talkspace-api/modules/user/model"
	"time"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type Client struct {
	Conn     *websocket.Conn
	Message  chan *Message
	ID       string `json:"id"`
	RoomID   string `json:"room_id"`
	ClientID string `json:"client_id"`
	Role	 string `json:"role"`
}

type Message struct {
	Content  string `json:"content"`
	RoomID   string `json:"room_id"`
	Username	 string `json:"username"`
	Role	 string `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

func (c *Client) WriteMessage(roomID string, db *gorm.DB) {
	defer func() {
		c.Conn.Close()
	}()

	message := &[]model.Message{}
	var username string
	doctor := &doctor.Doctor{}
	user := &user.User{}
	db.Find(&message, "consultation_id = ?", roomID)
	
	for _, m := range *message {
		if m.Role == "doctor" {
			db.Find(&doctor, "id = ?", m.ClientID)
			username = doctor.Fullname
		} else if m.Role == "user" {
			db.Find(&user, "id = ?", m.ClientID)
			username = user.Fullname
		}

		msg := &Message{
			Content:  m.Message,
			RoomID:   m.ConsultationID,
			Username: username,
			Role: m.Role,
			CreatedAt: m.CreatedAt,
		}

		c.Conn.WriteJSON(msg)
	}

	for {
		message, ok := <-c.Message
		if !ok {
			return
		}

		c.Conn.WriteJSON(message)
	}
}

func (c *Client) ReadMessage(hub *Hub, db *gorm.DB) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	var username string
	doctor := &doctor.Doctor{}
	user := &user.User{}

	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		if c.Role == "doctor" {
			db.Find(&doctor, "id = ?", c.ID)
			username = doctor.Fullname
		} else if c.Role == "user" {
			db.Find(&user, "id = ?", c.ID)
			username = user.Fullname
		}
		
		msg := &Message{
			Content:  string(m),
			RoomID:   c.RoomID,
			Username: username,
			Role: c.Role,
			CreatedAt: time.Now(),
		}

		messsage := &model.Message{
			ConsultationID: c.RoomID,
			ClientID: c.ID,
			Message: string(m),
			Role: c.Role,
		}

		db.Create(&messsage)

		hub.Broadcast <- msg
	}
}
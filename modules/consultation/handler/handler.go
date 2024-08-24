package handler

import (
	"fmt"
	"net/http"
	"talkspace-api/middlewares"
	"talkspace-api/modules/consultation/dto"
	"talkspace-api/modules/consultation/model"
	"talkspace-api/modules/consultation/usecase"
	doctor "talkspace-api/modules/doctor/model"
	user "talkspace-api/modules/user/model"
	"talkspace-api/utils/responses"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Handler struct {
	hub *usecase.Hub
	db *gorm.DB
}

func NewHandler(h *usecase.Hub, db *gorm.DB) *Handler {
	return &Handler{
		hub: h,
		db: db,
	}
}

func (h *Handler) CreateRoom(c echo.Context) error {
	var req dto.ConsultationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "error binding request",
		})
	}

	ID, _, _ := middlewares.ExtractToken(c)

	consultation := model.Consultation{
		DoctorID:      req.DoctorID,
		UserID:     ID,
		Status: 	 req.Status,
	}

	if len(h.hub.Rooms) == 0 {
		var rooms []model.Consultation
		h.db.Find(&rooms)
		h.hub.Rooms = make(map[string]*usecase.Room)
		for _, r := range rooms {
			h.hub.Rooms[r.ID] = &usecase.Room{
				ID:      r.ID,
				UserID: r.UserID,
				DoctorID: r.DoctorID,
				Client: make(map[string]*usecase.Client),
			}
		}
	}
	// h.db.Debug().Config.Logger
	check_db := h.db.Find(&model.Consultation{}).Where("user_id = ? AND doctor_id", consultation.UserID, consultation.DoctorID)
	if check_db.Error != nil {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse("room already exists in database"))
	}
	if h.hub.Rooms != nil {
		for _, r := range h.hub.Rooms {
			if r.UserID == consultation.UserID && r.DoctorID == consultation.DoctorID {
				return c.JSON(http.StatusBadRequest, responses.ErrorResponse("room already exists in hub"))
			}	
		}
	}	

	result := h.db.Create(&consultation)
	if result.Error != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "error creating room",
		})
	}

	h.hub.Rooms[consultation.ID] = &usecase.Room{
		ID:      consultation.ID,
		UserID: consultation.UserID,
		DoctorID: consultation.DoctorID,
		Client: make(map[string]*usecase.Client),
	}

	return c.JSON(http.StatusOK, responses.SuccessResponse("create room success", interface{}(nil)))
}


var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) JoinRoom(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "error binding request",
		})
		
	}

	ID, Role, _ := middlewares.ExtractToken(c)

	clientID := ID
	roomID := c.Param("roomId")
	role := Role

	if role == "admin" {
		return c.JSON(http.StatusBadRequest, responses.ErrorResponse("invalid role"))
	}

	cl := &usecase.Client{
		Conn:     conn,
		Message:  make(chan *usecase.Message, 10),
		ID:       clientID,
		RoomID:   roomID,
		Role: role,
	}

	h.hub.Register <- cl

	go cl.WriteMessage(roomID, h.db)
	cl.ReadMessage(h.hub, h.db)

	return nil
}

func (h *Handler) GetRooms(c echo.Context) error {
	roomsRes := make([]dto.RoomRes, 0)
	ID, _, _ := middlewares.ExtractToken(c)

	if len(h.hub.Rooms) == 0 {
		var rooms []model.Consultation
		h.db.Find(&rooms)
		h.hub.Rooms = make(map[string]*usecase.Room)
		for _, r := range rooms {
			h.hub.Rooms[r.ID] = &usecase.Room{
				ID:      r.ID,
				UserID: r.UserID,
				DoctorID: r.DoctorID,
				Client: make(map[string]*usecase.Client),
			}
		}
	}

	fmt.Println(h.hub.Rooms)

	for _, r := range h.hub.Rooms {
		if r.UserID == ID || r.DoctorID == ID {
			var user user.User
			var doctor doctor.Doctor
			h.db.Where("id = ?", r.UserID).Find(&user)
			h.db.Where("id = ?", r.DoctorID).Find(&doctor)
			fmt.Println(r.UserID, r.DoctorID)
			fmt.Println(user)
			roomsRes = append(roomsRes, dto.RoomRes{
				ID:   r.ID,
				DoctorProfilePicture: doctor.ProfilePicture,
				UserProfilePicture: user.ProfilePicture,
				DoctorName: doctor.Fullname,
				UserName: user.Fullname,
			})
		}
	}

	return c.JSON(http.StatusOK, responses.SuccessResponse("get rooms success", roomsRes))
}

func (h *Handler) GetDoctors(c echo.Context) error {
	doctors := []doctor.Doctor{}
	h.db.Find(&doctors)

	doctorMap := make([]dto.DoctorRes, 0)
	for _, d := range doctors {
		doctorMap = append(doctorMap, dto.DoctorRes{
			ID: d.ID,
			Fullname: d.Fullname,
			Email: d.Email,
			ProfilePicture: d.ProfilePicture,
			Role: d.Role,
			Specialist: d.Specialization,
			Experience: d.YearsOfExperience,
			Gender: d.Gender,
			Alumnus: d.Alumnus,
			AboutMe: d.About,
		})
	}

	return c.JSON(http.StatusOK, responses.SuccessResponse("get doctors success",doctorMap))
}

// type ClientRes struct {
// 	ID       string `json:"id"`
// 	Username string `json:"username"`
// }

// func (h *Handler) GetClients(c echo.Context) error {
// 	var clients []ClientRes
// 	roomId := c.Param("roomId")

// 	if _, ok := h.hub.Rooms[roomId]; !ok {
// 		clients = make([]ClientRes, 0)
// 		c.JSON(http.StatusOK, clients)
// 	}

// 	for _, c := range h.hub.Rooms[roomId].Client {
// 		clients = append(clients, ClientRes{
// 			ID:       c.ID,
// 			Username: c.Username,
// 		})
// 	}

// 	return c.JSON(http.StatusOK, clients)
// }

// func (h *Handler) GetChats(c echo.Context) error {
// 	roomID := c.Param("roomId")
// 	var messages []model.Message
// 	config.DB.Where("room_id = ?", roomID).Find(&messages)

// 	return c.JSON(http.StatusOK, messages)
// }
package recycleBox

import (
	"auth-api/internal/adapters/api"
	recycleBoxDomain "auth-api/internal/domain/recycleBox"
	customError "auth-api/internal/error"
	"auth-api/internal/midlleware"
	"auth-api/internal/utils"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
)

const (
	createRecycleBoxURL    = "/recyclebox"
	getRecycleBoxURL       = "/recyclebox/"
	updateRecycleBoxURL    = "/recyclebox/"
	addBottleURL           = "/recyclebox/add-bottle/"
	addBottleWithPointsURL = "/recyclebox/add-bottle-points/"
	GET                    = "GET "
	POST                   = "POST "
	PUT                    = "PUT "
)

type handler struct {
	recycleBoxService recycleBoxDomain.ServiceRecycleBox
}

func NewHandler(service recycleBoxDomain.ServiceRecycleBox) api.Handler {
	return &handler{recycleBoxService: service}
}

func (h *handler) Register(router *http.ServeMux) {
	router.Handle(POST+createRecycleBoxURL, midlleware.TimeoutMiddleware(midlleware.AdminMiddleware(http.HandlerFunc(h.CreateRecycleBox))))
	router.Handle(GET+getRecycleBoxURL, midlleware.TimeoutMiddleware(midlleware.AuthMiddleware(http.HandlerFunc(h.GetRecycleBox))))
	router.Handle(PUT+updateRecycleBoxURL, midlleware.TimeoutMiddleware(midlleware.AuthMiddleware(http.HandlerFunc(h.UpdateRecycleBox))))
	router.Handle(POST+addBottleURL, midlleware.TimeoutMiddleware(midlleware.AuthMiddleware(http.HandlerFunc(h.AddBottle))))
	router.Handle(POST+addBottleWithPointsURL, midlleware.TimeoutMiddleware(midlleware.AuthMiddleware(http.HandlerFunc(h.AddBottleWithPoints))))
}

// CreateRecycleBox handles creating a new recycle box (Admin only)
func (h *handler) CreateRecycleBox(w http.ResponseWriter, r *http.Request) {
	var dtoBox = &recycleBoxDomain.CreateRecycleBoxDTO{}
	if err := json.NewDecoder(r.Body).Decode(dtoBox); err != nil {
		handleJSONDecodeError(w, err)
		return
	}

	box, err := h.recycleBoxService.CreateRecycleBox(r.Context(), dtoBox)
	if err != nil {
		http.Error(w, "Failed to create recycle box", http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	utils.RenderJSON(w, http.StatusCreated, box)
}

// GetRecycleBox handles fetching a recycle box by ID
func (h *handler) GetRecycleBox(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromURL(r, getRecycleBoxURL)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	box, err := h.recycleBoxService.GetRecycleBox(r.Context(), id)
	if err != nil {
		if errors.Is(err, customError.NotFoundError) {
			http.Error(w, "Recycle box not found", http.StatusNotFound)
		} else {
			http.Error(w, "Unexpected error", http.StatusInternalServerError)
			log.Println(err.Error())
		}
		return
	}
	utils.RenderJSON(w, http.StatusOK, box)
}

// UpdateRecycleBox handles updating a recycle box's details
func (h *handler) UpdateRecycleBox(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromURL(r, updateRecycleBoxURL)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var dtoBox = &recycleBoxDomain.UpdateRecycleBoxDTO{}
	if err := json.NewDecoder(r.Body).Decode(dtoBox); err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	box, err := h.recycleBoxService.UpdateRecycleBox(r.Context(), id, dtoBox)
	if err != nil {
		if errors.Is(err, customError.NotFoundError) {
			http.Error(w, "Recycle box not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to update recycle box", http.StatusInternalServerError)
			log.Println(err.Error())
		}
		return
	}
	utils.RenderJSON(w, http.StatusOK, box)
}

// AddBottle handles adding a bottle to the recycle box (User access)
func (h *handler) AddBottle(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromURL(r, addBottleURL)
	if err != nil {
		http.Error(w, "Invalid recycle box ID", http.StatusBadRequest)
		return
	}

	box, err := h.recycleBoxService.AddBottle(r.Context(), id)
	if err != nil {
		if errors.Is(err, customError.NotFoundError) {
			http.Error(w, "Recycle box not found", http.StatusNotFound)
		} else if errors.Is(err, customError.BoxFullError) {
			http.Error(w, "Recycle box is full", http.StatusBadRequest)
		} else {
			http.Error(w, "Unexpected error", http.StatusInternalServerError)
			log.Println(err.Error())
		}
		return
	}
	utils.RenderJSON(w, http.StatusOK, box)
}

// AddBottleWithPoints handles adding a bottle and awarding points to the user (User access)
func (h *handler) AddBottleWithPoints(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromURL(r, addBottleWithPointsURL)
	if err != nil {
		http.Error(w, "Invalid recycle box ID", http.StatusBadRequest)
		return
	}

	// Extract user ID from context
	claims, ok := r.Context().Value("userClaims").(*midlleware.Claims)
	if !ok {
		http.Error(w, "User authentication error", http.StatusUnauthorized)
		return
	}

	box, err := h.recycleBoxService.AddBottleWithPoints(r.Context(), id, claims.UserID)
	if err != nil {
		if errors.Is(err, customError.NotFoundError) {
			http.Error(w, "Recycle box not found", http.StatusNotFound)
		} else if errors.Is(err, customError.BoxFullError) {
			http.Error(w, "Recycle box is full", http.StatusBadRequest)
		} else {
			http.Error(w, "Unexpected error", http.StatusInternalServerError)
			log.Println(err.Error())
		}
		return
	}
	utils.RenderJSON(w, http.StatusOK, box)
}

// Helper function to parse ID from URL
func getIDFromURL(r *http.Request, baseURL string) (int64, error) {
	idStr := r.URL.Path[len(baseURL):]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, errors.New("invalid ID format")
	}
	return id, nil
}

// Helper function to handle JSON decoding errors
func handleJSONDecodeError(w http.ResponseWriter, err error) {
	var unmarshalTypeError *json.UnmarshalTypeError
	var syntaxError *json.SyntaxError
	if errors.As(err, &unmarshalTypeError) {
		http.Error(w, "Invalid request data type", http.StatusBadRequest)
	} else if errors.As(err, &syntaxError) || errors.Is(err, io.ErrUnexpectedEOF) {
		http.Error(w, "Invalid JSON syntax", http.StatusBadRequest)
	} else if errors.Is(err, io.EOF) {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
	} else {
		log.Println(err.Error())
		http.Error(w, "Unexpected error", http.StatusInternalServerError)
	}
}

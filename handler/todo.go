package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// ServeHTTP implements [http.Handler].
func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// log.Println("TODOHandler ServeHTTP called")
	switch r.Method {
	case "POST":
		var req model.CreateTODORequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// log.Println("req", req)
		if req.Subject == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		resp, err := h.Create(r.Context(), &req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	case "GET":
		var req model.ReadTODORequest
		// Set default size
		if req.Size == 0 {
			req.Size = 5
		}

		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if prevIDs, ok := r.Form["prev_id"]; ok && len(prevIDs) > 0 {
			// log.Println("prev_id:", prevIDs[0])
			var prevID int64
			if _, err := fmt.Sscanf(prevIDs[0], "%d", &prevID); err == nil {
				req.PrevID = prevID
			} else {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		if sizes, ok := r.Form["size"]; ok && len(sizes) > 0 {
			log.Println("size:", sizes[0])
			var size int64
			if _, err := fmt.Sscanf(sizes[0], "%d", &size); err == nil {
				req.Size = size
			} else {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			log.Println("req.Size", req.Size)
		}
		// log.Println("req", req)
		// log.Println("ReadTODO handler called with PrevID:", req.PrevID, "Size:", req.Size)
		resp, err := h.Read(r.Context(), &req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	case "PUT":
		var req model.UpdateTODORequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// log.Println("req", req)
		if req.ID == 0 || req.Subject == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		resp, err := h.Update(r.Context(), &req)
		if err != nil {
			if errors.Is(err, &model.ErrNotFound{}) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	case "DELETE":
		var req model.DeleteTODORequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Println("req", req)
		if len(req.IDs) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		resp, err := h.Delete(r.Context(), &req)
		if err != nil {

			// log.Println("Delete req IDs:", req.IDs)
			// log.Println("Delete response:", resp)
			// log.Println("Delete error:", err)
			// log.Printf("err %%s: %s", err.Error())
			// log.Printf("err %%v: %v", err)
			// log.Printf("err %%w: %w", err)

			var notFoundErr *model.ErrNotFound
			if errors.As(err, &notFoundErr) {
				log.Println("Returning 404 Not Found")
				w.WriteHeader(http.StatusNotFound)
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	todo, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}
	return &model.CreateTODOResponse{
		TODO: *todo,
	}, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	// log.Println("Read handler:", req)
	todos, err := h.svc.ReadTODO(ctx, req.PrevID, req.Size)
	if err != nil {
		// log.Println("Read error:", err)
		return nil, err
	}
	result := make([]model.TODO, len(todos))
	for i, todo := range todos {
		// log.Println("todo:", todo)
		result[i] = *todo
	}
	// log.Println("result:", result)
	return &model.ReadTODOResponse{
		TODOs: result,
	}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	todo, err := h.svc.UpdateTODO(ctx, req.ID, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}
	return &model.UpdateTODOResponse{
		TODO: *todo,
	}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	ids := make([]int64, len(req.IDs))
	copy(ids, req.IDs)
	err := h.svc.DeleteTODO(ctx, ids)
	if err != nil {
		return nil, err
	}

	return &model.DeleteTODOResponse{}, nil
}

package file

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"orejametov/service-storage/pkg/logging"
	"orejametov/service-storage/util/responder"
	"github.com/go-chi/chi/v5"
)

func GetFile(service Service, logger logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("GET FILE")

		logger.Debug("get note_uuid from URL")
		noteUUID := r.URL.Query().Get("note_uuid")

		if noteUUID == "" {
			responder.Error(w, r, http.StatusUnprocessableEntity, errors.New("note_uuid query parameter is required"))
			return
		}

		logger.Debug("get fileId from context")
		fileId := chi.URLParam(r, "id")

		f, err := service.GetFile(r.Context(), noteUUID, fileId)
		if err != nil {
			responder.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", f.Name))
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))

		w.Write(f.Bytes)
	}
}

func GetFilesByNoteUUID(service Service, logger logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("GET FILES BY NOTE UUID")
		w.Header().Set("Content-Type", "form/json")

		logger.Debug("get note_uuid from URL")
		noteUUID := r.URL.Query().Get("note_uuid")
		if noteUUID == "" {
			responder.Error(w, r, http.StatusUnprocessableEntity, errors.New("note_uuid query parameter is required"))
			return
		}

		file, err := service.GetFilesByNoteUUID(r.Context(), noteUUID)
		if err != nil {
			responder.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		filesBytes, err := json.Marshal(file)
		if err != nil {
			responder.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(filesBytes)
	}
}

func CreateFile(service Service, logger logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("CREATE FILE")
		w.Header().Set("Content-Type", "form/json")

		// TODO maximum file size
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			responder.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		logger.Debug("decode create upload fileInfo dto")

		files, ok := r.MultipartForm.File["file"]
		if !ok || len(files) == 0 {
			responder.Error(w, r, http.StatusUnprocessableEntity, errors.New("file required"))
			return
		}
		fileInfo := files[0]
		fileReader, err := fileInfo.Open()

		if err != nil {
			responder.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		dto := CreateFileDTO{
			Name: fileInfo.Filename,
			Size: fileInfo.Size,
			Reader: fileReader,
		}

		err = service.Create(r.Context(), r.Form.Get("note_uuid"), dto)
		if err != nil {
			responder.Error(w, r, http.StatusInternalServerError, err)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}	
}

func DeleteFile(service Service, logger logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("DELETE FILE")
		w.Header().Set("Content-Type", "application/json")

		logger.Debug("get fileId from context")
		fileId := chi.URLParam(r, "id")

		noteUUID := r.URL.Query().Get("note_uuid")
		if noteUUID == "" {
			responder.Error(w, r, http.StatusUnprocessableEntity, errors.New("note_uuid query parameter is required"))
			return
		}

		err := service.Delete(r.Context(), noteUUID, fileId)
		if err != nil {
			responder.Error(w, r, http.StatusInternalServerError, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

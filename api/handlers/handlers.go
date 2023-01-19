package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rs/zerolog/hlog"
	"gitlab.com/napspan/SampleCompany/api/models"
	"gitlab.com/napspan/SampleCompany/api/utils"
)

// Index index
func Index(w http.ResponseWriter, r *http.Request) {
	mR := models.MyResponse{Msg: "Requests endpoint /index"}
	hlog.FromRequest(r).Info().Msg(mR.Msg)
	models.GenerateResponse(w, mR, http.StatusOK)
	return
}

// GetAllSavedComputers GetAllSavedComputers
func GetAllSavedComputers(w http.ResponseWriter, r *http.Request) {
	mR := models.MyResponse{Code: 1}
	mR.Code = 200
	mR.Data = models.GetAllComputers()
	models.GenerateResponse(w, mR, http.StatusOK)
}

// GetAllSavedComputers GetAllSavedComputers
func GetAllSavedComputersFromEmployee(w http.ResponseWriter, r *http.Request) {
	mR := models.MyResponse{Code: 1}
	type input struct {
		Employee string `json:"employee"`
	}

	i := input{}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := decoder.Decode(&i)
	if err != nil {
		mR.Msg = "unable to decode your request"
		hlog.FromRequest(r).Error().Caller().Err(err).Msg("")
		models.GenerateResponse(w, mR, 404)
		return
	}
	if i.Employee == "" || len(i.Employee) != 3 {
		err = errors.New("must send something in this field respecting the Employee format, 3 chars")
		hlog.FromRequest(r).Error().Caller().Err(err).Msg("")
		mR.Msg = err.Error()
		models.GenerateResponse(w, mR, http.StatusBadRequest)
		return
	}
	mR.Code = 200
	mR.Data = models.GetAllSavedComputersFromEmployee(i.Employee)
	models.GenerateResponse(w, mR, http.StatusOK)
}

// CreateNewComputer  CreateNewComputer
func CreateNewComputer(w http.ResponseWriter, r *http.Request) {
	mR := models.MyResponse{Code: 1}

	type input struct {
		MAC          string `json:"mac"`
		ComputerName string `json:"computer_name"`
		IP           string `json:"ip"`
		Employee     string `json:"employee,omitempty"`
		Description  string `json:"description,omitempty"`
	}

	i := input{}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := decoder.Decode(&i)
	if err != nil {
		mR.Msg = "unable to decode your request"
		hlog.FromRequest(r).Error().Caller().Err(err).Msg("")
		models.GenerateResponse(w, mR, 500)
		return
	}
	if i.MAC == "" || i.ComputerName == "" || i.IP == "" {
		err = errors.New("either MAC or ComputerName or IP are empty")
		hlog.FromRequest(r).Error().Caller().Err(err).Msg("")
		mR.Msg = err.Error()
		models.GenerateResponse(w, mR, http.StatusBadRequest)
		return
	}
	computerInfo := models.ComputerInfo{
		MAC:          i.MAC,
		ComputerName: i.ComputerName,
		IP:           i.IP,
	}
	if i.Employee != "" {
		if len(i.Employee) != 3 {
			err = errors.New("Employee format is wrong, must be 3 chars in length")
			hlog.FromRequest(r).Error().Caller().Err(err).Msg("")
			mR.Msg = err.Error()
			models.GenerateResponse(w, mR, http.StatusBadRequest)
			return
		}
		computerInfo.Employee = i.Employee
	}
	if i.Description != "" {
		computerInfo.Description = i.Description
	}
	notify := computerInfo.InsertInfo()
	if notify {
		err = utils.NotifyAdmin("warning", i.Employee, "has 3 or more computers")
		if err != nil {
			hlog.FromRequest(r).Error().Caller().Err(err).Msg("")
			models.GenerateResponse(w, mR, 500)
			return
		}
	}

	mR.Msg = "computerInfo Saved"
	mR.Code = 200
	models.GenerateResponse(w, mR, http.StatusOK)
}

// UpdateComputer  UpdateComputer
func UpdateComputer(w http.ResponseWriter, r *http.Request) {
	mR := models.MyResponse{Code: 1}

	type input struct {
		MAC          string `json:"mac"`
		ComputerName string `json:"computer_name"`
		IP           string `json:"ip"`
		Employee     string `json:"employee,omitempty"`
		Description  string `json:"description,omitempty"`
	}

	i := input{}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := decoder.Decode(&i)
	if err != nil {
		mR.Msg = "unable to decode your request"
		hlog.FromRequest(r).Error().Caller().Err(err).Msg("")
		models.GenerateResponse(w, mR, 500)
		return
	}
	if i.MAC == "" || i.ComputerName == "" || i.IP == "" {
		err = errors.New("either MAC or ComputerName or IP are empty")
		hlog.FromRequest(r).Error().Caller().Err(err).Msg("")
		mR.Msg = err.Error()
		models.GenerateResponse(w, mR, http.StatusBadRequest)
		return
	}
	computerInfo := models.ComputerInfo{
		MAC:          i.MAC,
		ComputerName: i.ComputerName,
		IP:           i.IP,
	}
	if i.Employee != "" {
		if len(i.Employee) != 3 {
			err = errors.New("Employee format is wront, must be 3 chars in length")
			hlog.FromRequest(r).Error().Caller().Err(err).Msg("")
			mR.Msg = err.Error()
			models.GenerateResponse(w, mR, http.StatusBadRequest)
			return
		}
		computerInfo.Employee = i.Employee
	}
	if i.Description != "" {
		computerInfo.Description = i.Description
	}
	computerInfo.UpdateInfoByMac()

	mR.Msg = "computerInfo updated"
	mR.Code = 200
	models.GenerateResponse(w, mR, http.StatusOK)
}

// DeleteComputer  DeleteComputer
func DeleteComputer(w http.ResponseWriter, r *http.Request) {
	mR := models.MyResponse{Code: 1}

	type input struct {
		MAC string `json:"mac"`
	}

	i := input{}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := decoder.Decode(&i)
	if err != nil {
		mR.Msg = "unable to decode your request"
		hlog.FromRequest(r).Error().Caller().Err(err).Msg("")
		models.GenerateResponse(w, mR, 404)
		return
	}
	if i.MAC == "" {
		err = errors.New("must send something in this field")
		hlog.FromRequest(r).Error().Caller().Err(err).Msg("")
		mR.Msg = err.Error()
		models.GenerateResponse(w, mR, http.StatusBadRequest)
		return
	}
	computerInfo := models.ComputerInfo{
		MAC: i.MAC,
	}
	computerInfo.DeleteInfoByMac()

	mR.Msg = "computerInfo deleted"
	mR.Code = 200
	models.GenerateResponse(w, mR, http.StatusOK)
}

// AssignComputer AssignComputer
func AssignComputer(w http.ResponseWriter, r *http.Request) {
	mR := models.MyResponse{Code: 1}
	type input struct {
		Employee string `json:"employee"`
		Mac      string `json:"mac"`
	}

	i := input{}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := decoder.Decode(&i)
	if err != nil {
		mR.Msg = "unable to decode your request"
		hlog.FromRequest(r).Error().Caller().Err(err).Msg("")
		models.GenerateResponse(w, mR, 500)
		return
	}
	if i.Employee == "" || len(i.Employee) != 3 {
		err = errors.New("must send something in this field respecting the Employee format, 3 chars")
		hlog.FromRequest(r).Error().Caller().Err(err).Msg("")
		mR.Msg = err.Error()
		models.GenerateResponse(w, mR, http.StatusBadRequest)
		return
	}
	if i.Mac == "" {
		err = errors.New("MAC must be filled")
		hlog.FromRequest(r).Error().Caller().Err(err).Msg("")
		mR.Msg = err.Error()
		models.GenerateResponse(w, mR, http.StatusBadRequest)
		return
	}

	notify := models.AssignComputer(i.Mac, i.Employee)
	if notify {
		err = utils.NotifyAdmin("warning", i.Employee, "has 3 or more computers")
		if err != nil {
			hlog.FromRequest(r).Error().Caller().Err(err).Msg("")
			models.GenerateResponse(w, mR, 500)
			return
		}
	}

	mR.Code = 200
	mR.Data = models.GetAllSavedComputersFromEmployee(i.Employee)
	models.GenerateResponse(w, mR, http.StatusOK)
}

// DisassignComputer DisassignComputer
func DisassignComputer(w http.ResponseWriter, r *http.Request) {
	mR := models.MyResponse{Code: 1}
	type input struct {
		Employee string `json:"employee"`
		Mac      string `json:"mac"`
	}

	i := input{}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := decoder.Decode(&i)
	if err != nil {
		mR.Msg = "unable to decode your request"
		hlog.FromRequest(r).Error().Caller().Err(err).Msg("")
		models.GenerateResponse(w, mR, 500)
		return
	}
	if i.Employee == "" || len(i.Employee) != 3 {
		err = errors.New("must send something in this field respecting the Employee format, 3 chars")
		hlog.FromRequest(r).Error().Caller().Err(err).Msg("")
		mR.Msg = err.Error()
		models.GenerateResponse(w, mR, http.StatusBadRequest)
		return
	}
	if i.Mac == "" {
		err = errors.New("MAC must be filled")
		hlog.FromRequest(r).Error().Caller().Err(err).Msg("")
		mR.Msg = err.Error()
		models.GenerateResponse(w, mR, http.StatusBadRequest)
		return
	}

	models.DisassignComputer(i.Mac, i.Employee)

	mR.Code = 200
	models.GenerateResponse(w, mR, http.StatusOK)
}

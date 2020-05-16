package usecases

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"

	"github.com/bmf-san/gobel-api/app/domain"
	"github.com/bmf-san/goblin"
)

// A TagInteractor is an interactor for a tag.
type TagInteractor struct {
	TagRepository TagRepository
	JSONResponse  JSONResponse
	Logger        Logger
}

// HandleIndex returns a listing of the resource.
func (ti *TagInteractor) HandleIndex(w http.ResponseWriter, r *http.Request) {
	const defaultPage = 1
	const defaultLimit = 10

	count, err := ti.TagRepository.CountAll()
	if err != nil {
		ti.Logger.LogError(err)
		ti.JSONResponse.Error500(w)
		return
	}

	paramPage := r.URL.Query().Get("page")
	var page int
	if paramPage == "" {
		page = defaultPage
	} else {
		page, err = strconv.Atoi(paramPage)
		if err != nil {
			ti.Logger.LogError(err)
			ti.JSONResponse.Error500(w)
			return
		}
	}

	paramLimit := r.URL.Query().Get("limit")
	var limit int
	if paramLimit == "" {
		limit = defaultLimit
	} else {
		limit, err = strconv.Atoi(paramLimit)
		if err != nil {
			ti.Logger.LogError(err)
			ti.JSONResponse.Error500(w)
			return
		}
	}

	var tags domain.Tags
	tags, err = ti.TagRepository.FindAll(page, limit)
	if err != nil {
		ti.Logger.LogError(err)
		ti.JSONResponse.Error500(w)
		return
	}

	var tr TagResponse
	var res []byte
	res, err = json.Marshal(tr.MakeResponseHandleIndex(tags))
	if err != nil {
		ti.Logger.LogError(err)
		ti.JSONResponse.Error500(w)
		return
	}

	pageCount := math.Ceil(float64(count) / float64(limit))

	w.Header().Set("Pagination-Count", strconv.Itoa(count))
	w.Header().Set("Pagination-Pagecount", fmt.Sprint(pageCount))
	w.Header().Set("Pagination-Page", strconv.Itoa(page))
	w.Header().Set("Pagination-Limit", strconv.Itoa(limit))
	ti.JSONResponse.Success200(w, res)
	return
}

// HandleIndexPrivate returns a listing of the resource.
func (ti *TagInteractor) HandleIndexPrivate(w http.ResponseWriter, r *http.Request) {
	const defaultPage = 1
	const defaultLimit = 10

	count, err := ti.TagRepository.CountAll()
	if err != nil {
		ti.Logger.LogError(err)
		ti.JSONResponse.Error500(w)
		return
	}

	paramPage := r.URL.Query().Get("page")
	var page int
	if paramPage == "" {
		page = defaultPage
	} else {
		page, err = strconv.Atoi(paramPage)
		if err != nil {
			ti.Logger.LogError(err)
			ti.JSONResponse.Error500(w)
			return
		}
	}

	paramLimit := r.URL.Query().Get("limit")
	var limit int
	if paramLimit == "" {
		limit = defaultLimit
	} else {
		limit, err = strconv.Atoi(paramLimit)
		if err != nil {
			ti.Logger.LogError(err)
			ti.JSONResponse.Error500(w)
			return
		}
	}

	var tags domain.Tags
	tags, err = ti.TagRepository.FindAll(page, limit)
	if err != nil {
		ti.Logger.LogError(err)
		ti.JSONResponse.Error500(w)
		return
	}

	var tr TagResponse
	var res []byte
	res, err = json.Marshal(tr.MakeResponseHandleIndexPrivate(tags))
	if err != nil {
		ti.Logger.LogError(err)
		ti.JSONResponse.Error500(w)
		return
	}

	pageCount := math.Ceil(float64(count) / float64(limit))

	w.Header().Set("Pagination-Count", strconv.Itoa(count))
	w.Header().Set("Pagination-Pagecount", fmt.Sprint(pageCount))
	w.Header().Set("Pagination-Page", strconv.Itoa(page))
	w.Header().Set("Pagination-Limit", strconv.Itoa(limit))
	ti.JSONResponse.Success200(w, res)
	return
}

// HandleShow display the specified resource.
func (ti *TagInteractor) HandleShow(w http.ResponseWriter, r *http.Request) {
	name := goblin.GetParam(r.Context(), "name")

	var tag domain.Tag
	tag, err := ti.TagRepository.FindByName(name)
	if err != nil {
		ti.Logger.LogError(err)
		ti.JSONResponse.Error500(w)
		return
	}

	var tr TagResponse
	var res []byte
	res, err = json.Marshal(tr.MakeResponseHandleShow(tag))
	if err != nil {
		ti.Logger.LogError(err)
		ti.JSONResponse.Error500(w)
		return
	}

	ti.JSONResponse.Success200(w, res)
	return
}

// HandleShowPrivate display the specified resource.
func (ti *TagInteractor) HandleShowPrivate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(goblin.GetParam(r.Context(), "id"))
	if err != nil {
		ti.Logger.LogError(err)
		ti.JSONResponse.Error500(w)
		return
	}

	var tag domain.Tag
	tag, err = ti.TagRepository.FindByID(id)
	if err != nil {
		ti.Logger.LogError(err)
		ti.JSONResponse.Error500(w)
		return
	}

	var tr TagResponse
	var res []byte
	res, err = json.Marshal(tr.MakeResponseHandleShowPrivate(tag))
	if err != nil {
		ti.Logger.LogError(err)
		ti.JSONResponse.Error500(w)
		return
	}

	ti.JSONResponse.Success200(w, res)
	return
}

// HandleStorePrivate stores a newly created resource in storage.
func (ti *TagInteractor) HandleStorePrivate(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ti.Logger.LogError(err)
		ti.JSONResponse.Error500(w)
		return
	}

	var req RequestTag

	err = json.Unmarshal(body, &req)
	if err != nil {
		ti.Logger.LogError(err)
		ti.JSONResponse.Error500(w)
		return
	}

	err = ti.TagRepository.Save(req)
	if err != nil {
		ti.Logger.LogError(err)
		ti.JSONResponse.Error500(w)
		return
	}

	ti.JSONResponse.Success201(w, []byte("The item was created successfully"))
	return
}

// HandleUpdatePrivate updates the specified resource in storage.
func (ti *TagInteractor) HandleUpdatePrivate(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ti.Logger.LogError(err)
		ti.JSONResponse.Error500(w)
		return
	}

	var req RequestTag

	id, err := strconv.Atoi(goblin.GetParam(r.Context(), "id"))
	if err != nil {
		ti.Logger.LogError(err)
		ti.JSONResponse.Error500(w)
		return
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		ti.Logger.LogError(err)
		ti.JSONResponse.Error500(w)
		return
	}

	err = ti.TagRepository.SaveByID(req, id)
	if err != nil {
		ti.Logger.LogError(err)
		ti.JSONResponse.Error500(w)
		return
	}

	ti.JSONResponse.Success200(w, []byte("The item was updated successfully"))
	return
}

// HandleDestroyPrivate removes the specified resource from storage.
func (ti *TagInteractor) HandleDestroyPrivate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(goblin.GetParam(r.Context(), "id"))
	if err != nil {
		ti.Logger.LogError(err)
		ti.JSONResponse.Error500(w)
		return
	}
	count, err := ti.TagRepository.DeleteByID(id)
	if err != nil {
		ti.Logger.LogError(err)
		ti.JSONResponse.Error500(w)
		return
	}

	if count == 0 {
		ti.JSONResponse.Error404(w)
		return
	}

	ti.JSONResponse.Success200(w, []byte("The item was deleted successfully"))
	return
}

package api

import (
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"
)

func ensureEventCollection(req *rest.Request, params *QueryParams) error {
	if params.CollectionName != "" {
		return nil
	}

	if v := req.FormValue("event_collection"); v != "" {
		params.CollectionName = v
		return nil
	} else {
		return ParamsMissingError
	}
}

func HandleQueryCount(w rest.ResponseWriter, req *rest.Request) {
	if err := RequireReadKey(w, req); err != nil {
		rest.Error(w, err.Error(), err.(StatusError).Code)
		return
	}

	project := currentProject(req)

	var params QueryParams
	if err := req.DecodeJsonPayload(&params); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ensureEventCollection(req, &params)

	result, err := count(project, &params)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteJson(result)
	}
}

func HandleQueryUniqueCount(w rest.ResponseWriter, req *rest.Request) {
	if err := RequireReadKey(w, req); err != nil {
		rest.Error(w, err.Error(), err.(StatusError).Code)
		return
	}

	project := currentProject(req)

	var params struct {
		*QueryParams
		TargetProperty string `json:"target_property"`
	}

	if err := req.DecodeJsonPayload(&params); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ensureEventCollection(req, params.QueryParams)

	result, err := countUnique(project, params.TargetProperty, params.QueryParams)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteJson(result)
	}
}

func HandleQueryMin(w rest.ResponseWriter, req *rest.Request) {
	if err := RequireReadKey(w, req); err != nil {
		rest.Error(w, err.Error(), err.(StatusError).Code)
		return
	}

	project := currentProject(req)

	var params struct {
		*QueryParams
		TargetProperty string `json:"target_property"`
	}

	if err := req.DecodeJsonPayload(&params); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ensureEventCollection(req, params.QueryParams)

	result, err := min(project, params.TargetProperty, params.QueryParams)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteJson(result)
	}
}

func HandleQueryMax(w rest.ResponseWriter, req *rest.Request) {
	if err := RequireReadKey(w, req); err != nil {
		rest.Error(w, err.Error(), err.(StatusError).Code)
		return
	}

	project := currentProject(req)

	var params struct {
		*QueryParams
		TargetProperty string `json:"target_property"`
	}

	if err := req.DecodeJsonPayload(&params); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ensureEventCollection(req, params.QueryParams)

	result, err := max(project, params.TargetProperty, params.QueryParams)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteJson(result)
	}
}

func HandleQuerySum(w rest.ResponseWriter, req *rest.Request) {
	if err := RequireReadKey(w, req); err != nil {
		rest.Error(w, err.Error(), err.(StatusError).Code)
		return
	}

	project := currentProject(req)

	var params struct {
		*QueryParams
		TargetProperty string `json:"target_property"`
	}

	if err := req.DecodeJsonPayload(&params); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ensureEventCollection(req, params.QueryParams)

	result, err := sum(project, params.TargetProperty, params.QueryParams)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteJson(result)
	}
}

func HandleQueryAverage(w rest.ResponseWriter, req *rest.Request) {
	if err := RequireReadKey(w, req); err != nil {
		rest.Error(w, err.Error(), err.(StatusError).Code)
		return
	}

	project := currentProject(req)

	var params struct {
		*QueryParams
		TargetProperty string `json:"target_property"`
	}

	if err := req.DecodeJsonPayload(&params); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ensureEventCollection(req, params.QueryParams)

	result, err := average(project, params.TargetProperty, params.QueryParams)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteJson(result)
	}
}

func HandleQueryPercentile(w rest.ResponseWriter, req *rest.Request) {
	if err := RequireReadKey(w, req); err != nil {
		rest.Error(w, err.Error(), err.(StatusError).Code)
		return
	}

	project := currentProject(req)

	var params struct {
		*QueryParams
		TargetProperty string  `json:"target_property"`
		Percent        float64 `json:"percent"`
	}

	if err := req.DecodeJsonPayload(&params); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ensureEventCollection(req, params.QueryParams)

	result, err := percentile(project, params.TargetProperty, params.Percent, params.QueryParams)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteJson(result)
	}
}

func HandleQueryMedian(w rest.ResponseWriter, req *rest.Request) {
	if err := RequireReadKey(w, req); err != nil {
		rest.Error(w, err.Error(), err.(StatusError).Code)
		return
	}

	project := currentProject(req)

	var params struct {
		*QueryParams
		TargetProperty string  `json:"target_property"`
		Percent        float64 `json:"percent"`
	}

	if err := req.DecodeJsonPayload(&params); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ensureEventCollection(req, params.QueryParams)

	result, err := median(project, params.TargetProperty, params.QueryParams)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteJson(result)
	}
}

func HandleQuerySelectUnique(w rest.ResponseWriter, req *rest.Request) {
	if err := RequireReadKey(w, req); err != nil {
		rest.Error(w, err.Error(), err.(StatusError).Code)
		return
	}

	project := currentProject(req)

	var params struct {
		*QueryParams
		TargetProperty string `json:"target_property"`
	}

	if err := req.DecodeJsonPayload(&params); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ensureEventCollection(req, params.QueryParams)
	if !params.TimeFrame.IsGiven() {
		err := ParamsMissingError
		rest.Error(w, err.Error(), err.Code)
		return
	}

	result, err := selectUnique(project, params.TargetProperty, params.QueryParams)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteJson(result)
	}
}

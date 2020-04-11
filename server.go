package hive

import (
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	gapi "github.com/suborbital/gust/gapi/server"
)

// Server is a hive server
type Server struct {
	*gapi.Server
	h        *Hive
	inFlight map[string]*Result
}

func newServer(h *Hive, opts ...gapi.OptionsModifier) *Server {
	s := gapi.New(opts...)

	server := &Server{
		Server:   s,
		h:        h,
		inFlight: make(map[string]*Result),
	}

	server.POST("/do/:jobtype", server.scheduleHandler())

	return server
}

func (s *Server) scheduleHandler() gapi.HandlerFunc {
	return func(r *http.Request, ctx *gapi.Ctx) (interface{}, error) {
		jobType := ctx.Params.ByName("jobtype")
		if jobType == "" {
			return nil, gapi.E(http.StatusBadRequest, "missing jobtype")
		}

		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, gapi.E(http.StatusInternalServerError, "failed to read request body")
		}
		defer r.Body.Close()

		res := s.h.Do(NewJob(jobType, data))

		then := r.URL.Query().Get("then")
		if then == "true" {
			result, err := res.Then()
			if err != nil {
				return nil, gapi.E(http.StatusNoContent, errors.Wrap(err, "job resulted in error").Error())
			}

			return result, nil
		}

		s.addInFlight(res)

		return []byte(res.ID), nil
	}
}

func (s *Server) addInFlight(r *Result) {
	s.inFlight[r.ID] = r
}

func (s *Server) getInFlight(id string) *Result {
	r, ok := s.inFlight[id]
	if !ok {
		return nil
	}

	return r
}

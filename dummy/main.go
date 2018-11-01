package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/sirupsen/logrus"

	httplogger "github.com/go-http-utils/logger"

	cm "github.com/gopot/concurrent-map"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", NewDummyHandler())

	if err := http.ListenAndServe(":1818", httplogger.Handler(mux, os.Stdout, httplogger.CombineLoggerType)); err != nil {
		logrus.Fatal(err)
	}
}

const PATH_SEPARATOR = "/"

func NewDummyHandler() func(w http.ResponseWriter, req *http.Request) {

	log := logrus.New()

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	log.SetLevel(logrus.TraceLevel)

	baseStorage := &storage{cm: cm.New(1), log: log.WithField("node", "base")}

	return func(w http.ResponseWriter, req *http.Request) {

		path := req.URL.Path
		pathSteps := strings.Split(path, PATH_SEPARATOR)[1:]
		switch req.Method {
		case http.MethodGet:

			elem, found := baseStorage.Get(pathSteps...)
			if !found {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, "Path \"%s\" not found", path)
				return
			}
			// On the leaf
			switch e := elem.(type) {
			case *storage:
				items := e.cm.Items()
				keystr := "Begin:\r\n"
				for k, _ := range items {
					logrus.Debugf("Adding key %#v to response", k)
					keystr += k.(string) + "\r\n"
				}
				w.Write([]byte(keystr))
			case string:
				w.Write([]byte(e + "\r\n"))
			}
		case http.MethodPut:

			defer req.Body.Close()
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				logrus.Errorf("Failed to read body: %e", err)
			}
			baseStorage.Put(string(body), pathSteps...)
		}
	}
}

type storage struct {
	cm  *cm.ConcurrentMap
	log logrus.FieldLogger
}

func (s *storage) Get(steps ...string) (leaf interface{}, found bool) {
	s.log.Debugf("Method GET called with  %#v", steps)
	s.log.Debugf("Looking for leaf %s", steps[0])
	leaf, found = s.cm.Get(steps[0])
	s.log.Debugf("leaf = %#v is found = %v", leaf, found)

	nextSteps := steps[1:]
	s.log.Debugf("NextSteps = %#v", nextSteps)

	if len(nextSteps) == 0 {
		s.log.Debug("There are no next steps - returning")
		return
	}
	s.log.Debugf("Asserting leaf %#v to *storage", leaf)
	if nextStorage, ok := leaf.(*storage); ok {
		s.log.Debugf("Assertion succeed. Getting next steps %#v", nextSteps)

		leaf, found = nextStorage.Get(nextSteps...)
		return
	}
	s.log.Debug("Assertion failed. returning")

	return nil, false
}

func (s *storage) Put(doc string, steps ...string) {

	s.log.Debugf("Method PUT called with \"%s\" and %#v", doc, steps)
	nextSteps := steps[1:]
	if len(nextSteps) == 0 {
		s.log.Infof("Assigning to %s value %v", steps[0], doc)
		s.cm.Set(steps[0], doc)
		return
	}

	s.log.Debugf("Looking for existing next node %s", steps[0])
	nextStorage, found := s.cm.Get(steps[0])
	if !found {
		s.log.Debugf("next node %s not found", steps[0])
		log := s.log.WithField("node", steps[0])
		nextStorage = &storage{cm: cm.New(1), log: log}
		s.log.Debugf("Creating next node %s", steps[0])
		s.cm.Set(steps[0], nextStorage)
	}

	s.log.Debugf("Putting %s into next steps %#v", doc, nextSteps)
	nextStorage.(*storage).Put(doc, nextSteps...)
	return
}

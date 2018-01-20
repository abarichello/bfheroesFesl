package server

import (
	"bytes"
	"fmt"
	"net/http"
	"text/template"

	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

const (
	tplEntitlements = "entitlements.xml"
	tplProducts     = "products.xml"
	tplRelationship = "relationship.xml"
	tplSessionNew   = "session_new.xml"
	tplSession      = "session.xml"
	tplWallets      = "wallets.xml"
	tplStore        = "store.xml"
)

type renderer struct {
	tpl *template.Template
}

func newRenderer() renderer {
	tpl := template.Must(
		template.ParseFiles(
			addPathPrefix(
				tplEntitlements,
				tplProducts,
				tplRelationship,
				tplSessionNew,
				tplSession,
				tplWallets,
				tplStore,
			)...,
		),
	)
	return renderer{tpl}
}

func addPathPrefix(path ...string) []string {
	prefixed := []string{}
	for _, p := range path {
		prefixed = append(prefixed, fmt.Sprintf("server/tpl/%s", p))
	}
	return prefixed
}

func (rdr *renderer) renderXML(w http.ResponseWriter, r *http.Request, name string, data interface{}) {
	buf := new(bytes.Buffer)
	err := rdr.tpl.ExecuteTemplate(buf, name, data)
	if err != nil {
		logrus.WithError(err).Error("Failed render XML", name)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// logrus.Info("Response", buf.String())
	respondXML(w, r, buf.Bytes())
}

func respondXML(w http.ResponseWriter, r *http.Request, v []byte) {
	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	if status, ok := r.Context().Value(render.StatusCtxKey).(int); ok {
		w.WriteHeader(status)
	}
	w.Write(v)
}

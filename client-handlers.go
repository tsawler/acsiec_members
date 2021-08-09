package clienthandlers

import (
	"fmt"
	"github.com/tsawler/goblender/pkg/cache"
	"github.com/tsawler/goblender/pkg/helpers"
	"github.com/tsawler/goblender/pkg/models"
	"github.com/tsawler/goblender/pkg/templates"
	"net/http"
)

var insidePageTemplate = "page.page.tmpl"

// ShowPage shows a page by slug
func ACSIECShowPage(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Query().Get(":slug")

	inCache, err := cache.Has(fmt.Sprintf("page-%s", slug))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	var pg models.Page

	if inCache {
		result, err := cache.Get(fmt.Sprintf("page-%s", slug))
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		pg = result.(models.Page)
	} else {
		p, err := repo.DB.GetPageBySlug(slug)
		if err == models.ErrNoRecord {
			helpers.NotFound(w)
			return
		} else if err != nil {
			helpers.ServerError(w, err)
			return
		}

		err = cache.Set(fmt.Sprintf("page-%s", slug), p)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		pg = p
	}

	// if we have an access level setting, verify that we are allowed to see the page
	if !helpers.CheckPageAccess(r, pg) {
		app.Session.Put(r.Context(), "flash", "Log in first!")
		u := r.URL.Path
		http.Redirect(w, r, fmt.Sprintf("/user/login?target=%s", u), http.StatusFound)
		return
	}

	// check to see if it's a principal page
	isPrincipalPage := dbModel.IsPrincipalPage(pg.ID)
	if isPrincipalPage {
		// make sure user has role
		if !app.Session.Exists(r.Context(), "userID") {
			u := r.URL.Path
			http.Redirect(w, r, fmt.Sprintf("/user/login?target=%s", u), http.StatusFound)
			return
		}
		userID := app.Session.GetInt(r.Context(), "userID")
		// check that user has role
		user, _ := repo.DB.GetUserById(userID)
		if _, ok := user.Roles["principal"]; !ok {
			helpers.ClientError(w, http.StatusUnauthorized)
			return
		}
	}

	// make sure the page is active
	if pg.Active == 0 {
		currentUser, ok := app.Session.Get(r.Context(), "user").(models.User)
		if ok {
			if currentUser.AccessLevel < 3 {
				app.Session.Put(r.Context(), "flash", "Log in first!")
				u := r.URL.Path
				http.Redirect(w, r, fmt.Sprintf("/user/login?target=%s", u), http.StatusFound)
				return
			} else {
				app.Session.Put(r.Context(), "warning", "Note: This page is inactive!")
			}
		} else {
			app.Session.Put(r.Context(), "flash", "Log in first!")
			u := r.URL.Path
			http.Redirect(w, r, fmt.Sprintf("/user/login?target=%s", u), http.StatusFound)
			return
		}
	}
	helpers.Render(w, r, insidePageTemplate, &templates.TemplateData{Page: pg})

}

// PrincipalHandler is an example handler
func PrincipalHandler(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Query().Get(":slug")

	inCache, err := cache.Has(fmt.Sprintf("page-%s", slug))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	var pg models.Page

	if inCache {
		result, err := cache.Get(fmt.Sprintf("page-%s", slug))
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		pg = result.(models.Page)
	} else {
		p, err := repo.DB.GetPageBySlug(slug)
		if err == models.ErrNoRecord {
			helpers.NotFound(w)
			return
		} else if err != nil {
			helpers.ServerError(w, err)
			return
		}

		err = cache.Set(fmt.Sprintf("page-%s", slug), p)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		pg = p
	}
	helpers.Render(w, r, insidePageTemplate, &templates.TemplateData{Page: pg})
}

// CustomShowHome is a sample handler which returns the home page using our local page template for the client,
// and is called from client-routes.go using a route that overrides the one in goBlender. This allows us
// to build custom functionality without having to use non-standard routes.
func CustomShowHome(w http.ResponseWriter, r *http.Request) {
	// do something interesting here, and then render the template
	helpers.Render(w, r, "client-sample.page.tmpl", &templates.TemplateData{})
}

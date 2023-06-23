package logout

import (
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

// Handler for our logout.
func Handler(ctx *gin.Context) {
	logoutUrl, err := url.Parse(os.Getenv("0AUTH2_ISSUER_URL") + os.Getenv("0AUTH2_LOGOUT_PATH"))
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}

	returnTo, err := url.Parse(scheme + "://" + ctx.Request.Host)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	parameters := url.Values{}
	parameters.Add("returnTo", returnTo.String())
	parameters.Add("client_id", os.Getenv("AUTH0_CLIENT_ID"))
	logoutUrl.RawQuery = parameters.Encode()

	ctx.Redirect(http.StatusTemporaryRedirect, logoutUrl.String())
}

package xhttp

import (
	"net/http"
	"net/url"

	"github.com/FTchinese/content-api/pkg/ids"
	"github.com/guregu/null"
)

func UserIDsFromHeader(h http.Header) ids.UserIDs {
	ftcID := h.Get(XUserID)
	unionID := h.Get(XUnionID)

	return ids.UserIDs{
		CompoundID: "",
		FtcID:      null.NewString(ftcID, ftcID != ""),
		UnionID:    null.NewString(unionID, unionID != ""),
	}.MustNormalize()
}

func UserIDsFromQuery(v url.Values) ids.UserIDs {
	ftcId := v.Get("ftc_id")
	unionID := v.Get("union_id")

	return ids.UserIDs{
		CompoundID: "",
		FtcID:      null.NewString(ftcId, ftcId != ""),
		UnionID:    null.NewString(unionID, unionID != ""),
	}.MustNormalize()
}

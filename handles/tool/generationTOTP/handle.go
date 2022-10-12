package generationTOTP

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/anyingiit/UnifyAuthCenter/utils"
)

func Handle(w http.ResponseWriter, r *http.Request) error {
	useDefaultIssureAndDefaultAccountName := false
	type Query struct {
		Issure      string
		AccountName string
	}
	query := Query{
		Issure:      r.FormValue("issure"),
		AccountName: r.FormValue("account_name"),
	}
	if query.Issure == "" || query.AccountName == "" {
		useDefaultIssureAndDefaultAccountName = true
		query.Issure = "unifyauthcenter.com"
		query.AccountName = "Jack A. Doe"
	}
	// fmt.Println(query)

	t, err := template.ParseFiles("./template/tool/generation_TOTP/generation.tmpl")
	if err != nil {
		return err
	}
	// fmt.Println(t)

	secret, pngBase64String, err := utils.GenerationNewTOTP(query.Issure, query.AccountName, 200, 200)

	if err != nil {
		return fmt.Errorf("generation TOTP failed, err: %s", err.Error())
	}

	// fmt.Println(secret, pngBase64String)

	err = t.Execute(w, struct {
		UseDefaultIssureAndDefaultAccountName string
		Issure                                string
		AccountName                           string
		Secret                                string
		/*
			当这里类型为string时, 实际渲染出的网页会显示为`#ZgotmplZ`
			这是因为...
				"ZgotmplZ" is a special value that indicates that unsafe content reached a
				CSS or URL context at runtime. The output of the example will be
					<img src="#ZgotmplZ">
				If the data comes from a trusted source, use content types to exempt it
				from filtering: URL(`javascript:...`).
		*/
		SecretPngBase64 template.URL
	}{
		UseDefaultIssureAndDefaultAccountName: fmt.Sprintf("%t", useDefaultIssureAndDefaultAccountName),
		Issure:                                query.Issure,
		AccountName:                           query.AccountName,
		Secret:                                secret,
		SecretPngBase64:                       template.URL(`data:image/png;base64,` + pngBase64String),
	})

	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

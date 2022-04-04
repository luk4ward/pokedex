package funtranslations

type (
	// TranslateRequest request with a text for Funtranslations API
	TranslateRequest struct {
		Text string `json:"text"`
	}

	// TranslateResponse raw response from Funtranslations
	TranslateResponse struct {
		Contents Contents `json:"contents"`
	}

	// Contents contents of a translation
	Contents struct {
		Text        string `json:"text"`
		Translated  string `json:"translated"`
		Translation string `json:"translation"`
	}
)

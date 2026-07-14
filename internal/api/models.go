package api

type chirpPost struct {
	Body string `json:"body"`
}

type apiError struct {
	Error string `json:"error"`
}

type cleanBodyResp struct {
	CleanedBody string `json:"cleaned_body"`
}

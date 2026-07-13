package main

type ChirpPost struct {
	Body string `json:"body"`
}

type APIError struct {
	Error string `json:"error"`
}

type CleanBodyResp struct {
	CleanedBody string `json:"cleaned_body"`
}

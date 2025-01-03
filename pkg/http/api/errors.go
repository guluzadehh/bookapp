package api

func UnexpectedError() Response {
	return Err("An unexpected error occured.")
}

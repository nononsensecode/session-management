.PHONY: openapi
openapi: openapi_http

.PHONY: openapi_http
openapi_http:
	oapi-codegen -generate types -o ports/openapi_types.gen.go -package ports user.yaml
	oapi-codegen -generate chi-server -o ports/openapi_api.gen.go -package ports user.yaml
package whatsapp

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
)

const (
	apiName string = "whatsapp"
	apiUrl  string = "https://graph.facebook.com/v22.0"
)

// Service stores Service configuration
type Service struct {
	accessToken   string
	httpService   *go_http.Service
	errorResponse ErrorResponse
}

type ServiceConfig struct {
	AccessToken string
}

func NewService(serviceConfig *ServiceConfig) (*Service, *errortools.Error) {
	if serviceConfig == nil {
		return nil, errortools.ErrorMessage("ServiceConfig must not be a nil pointer")
	}

	if serviceConfig.AccessToken == "" {
		return nil, errortools.ErrorMessage("AccessToken not provided")
	}

	httpService, e := go_http.NewService(&go_http.ServiceConfig{})
	if e != nil {
		return nil, e
	}

	return &Service{
		accessToken: serviceConfig.AccessToken,
		httpService: httpService,
	}, nil
}

func (service *Service) Error() Error {
	return service.errorResponse.Error
}

func (service *Service) httpRequest(requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *errortools.Error) {
	// add authentication header
	header := http.Header{}
	header.Set("Authorization", fmt.Sprintf("Bearer %s", service.accessToken))
	(*requestConfig).NonDefaultHeaders = &header

	// add error model
	service.errorResponse = ErrorResponse{}
	(*requestConfig).ErrorModel = &service.errorResponse

	request, response, e := service.httpService.HttpRequest(requestConfig)
	if service.errorResponse.Error.ErrorUserMsg != "" {
		if service.errorResponse.Error.ErrorUserTitle != "" {
			e.SetMessagef("%s: %s", service.errorResponse.Error.ErrorUserTitle, service.errorResponse.Error.ErrorUserMsg)
		} else {
			e.SetMessage(service.errorResponse.Error.ErrorUserMsg)
		}
	} else if service.errorResponse.Error.Message != "" {
		e.SetMessage(service.errorResponse.Error.Message)
	}

	return request, response, e
}

func (service *Service) url(path string) string {
	return fmt.Sprintf("%s/%s", apiUrl, path)
}

func (service *Service) AccessToken(accessToken string) {
	service.accessToken = accessToken
}

func (service Service) ApiName() string {
	return apiName
}

func (service Service) ApiKey() string {
	return service.accessToken
}

func (service Service) ApiCallCount() int64 {
	return service.httpService.RequestCount()
}

func (service Service) ApiReset() {
	service.httpService.ResetRequestCount()
}

package constant

import "net/http"

const (
	StatusOK                  = "OK"
	StatusInvalidParameter    = "INVALID_PARAMETER"
	StatusUnauthorized        = "UNAUTHORIZED"
	StatusNoTradingPermission = "NO_TRADING_PERMISSION"
	StatusNoSubscription      = "NO_SUBSCRIPTION"
	StatusNotFound            = "NOT_FOUND"
	StatusDuplicateCall       = "DUPLICATE_CALL"
	StatusServiceBusy         = "SERVICE_BUSY"
	StatusSystemError         = "SYSTEM_ERROR"
	StatusVendorError         = "VENDOR_ERROR"
	StatusBadGateway          = "BAD_GATEWAY"
	StatusMaintenance         = "MAINTENANCE"
	StatusGatewayTimeout      = "GATEWAY_TIMEOUT"

	// not in docs
	StatusForbidden = "FORBIDDEN"
)

func StatusHTTP(status string) int {
	switch status {
	case StatusOK:
		return http.StatusOK
	case StatusInvalidParameter:
		return http.StatusBadRequest
	case StatusUnauthorized:
		return http.StatusUnauthorized
	case StatusNoTradingPermission:
		return http.StatusForbidden
	case StatusNoSubscription:
		return http.StatusForbidden
	case StatusNotFound:
		return http.StatusNotFound
	case StatusDuplicateCall:
		return http.StatusPreconditionFailed
	case StatusServiceBusy:
		return http.StatusTooManyRequests
	case StatusSystemError:
		return http.StatusInternalServerError
	case StatusVendorError:
		return http.StatusInternalServerError
	case StatusBadGateway:
		return http.StatusBadGateway
	case StatusMaintenance:
		return http.StatusServiceUnavailable
	case StatusGatewayTimeout:
		return http.StatusGatewayTimeout
	case StatusForbidden:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

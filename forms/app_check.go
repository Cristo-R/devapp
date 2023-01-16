package forms

import (
	log "github.com/sirupsen/logrus"

	"gitlab.shoplazza.site/shoplaza/cobra/models"
	"gitlab.shoplazza.site/shoplaza/cobra/utils"
)

type (
	// AppCheckResult display check app result, IsValid display result is valid
	// if IsValid = true, len(InvalidFields) must be to zero
	AppCheckResult struct {
		Title         string            `json:"title"`
		InvalidFields utils.StringArray `json:"invalid_fields"`
		IsValid       bool              `json:"is_valid"`
	}

	AppCheckTestResult struct {
		Results []*AppCheckResult `json:"results"`
	}

	AppCheckSubmitResult struct {
		Results []*AppCheckResult `json:"results"`
	}
)

// NewAppCheckResult create check result, default IsVialed is true
func NewAppCheckResult(title string) *AppCheckResult {
	return &AppCheckResult{
		Title:         title,
		InvalidFields: make([]string, 0),
		IsValid:       true,
	}
}

// NewAppCheckTestResult create check result slice
func NewAppCheckTestResult() *AppCheckTestResult {
	return &AppCheckTestResult{Results: make([]*AppCheckResult, 0)}
}

// NewAppCheckSubmitResult create check result slice
func NewAppCheckSubmitResult() *AppCheckSubmitResult {
	return &AppCheckSubmitResult{Results: make([]*AppCheckResult, 0)}
}

// AppendCheckResult append app check results to test result
func (r *AppCheckTestResult) AppendCheckResult(results ...*AppCheckResult) {
	r.Results = append(r.Results, results...)
}

// AppCheckSubmitResult append check results to submit result
func (r *AppCheckSubmitResult) AppendCheckResult(results ...*AppCheckResult) {
	r.Results = append(r.Results, results...)
}

// AppendInvalidField append invalid field to invalid_fields
func (r *AppCheckResult) AppendInvalidField(field string) {
	if r.InvalidFields == nil {
		r.InvalidFields = make([]string, 0)
	}

	if r.IsValid {
		r.IsValid = false
	}

	r.InvalidFields = append(r.InvalidFields, field)
}

// IsResultValid return check result whether is valid
// valid return true, if app check result valid, len(invalid_fields) is zero
func (r *AppCheckResult) IsResultValid() bool {
	if r != nil {
		if r.IsValid && len(r.InvalidFields) != 0 {
			log.Errorf("unknown app check result. if result is valid, invalid_fields length must be to 0. result: %+v", r)
			return false
		}
		return r.IsValid
	}
	return false
}

// FilterValidAppCheckResult delete valid data in app check results
func FilterValidAppCheckResult(results []*AppCheckResult) []*AppCheckResult {
	ret := make([]*AppCheckResult, 0)
	for _, result := range results {
		if !result.IsResultValid() {
			ret = append(ret, result)
		}
	}
	return ret
}

// FilterRedundancyCheckResult filter redundancy result, only save old check results title
func FilterRedundancyCheckResult(oldResults []*AppCheckResult, newResults []*AppCheckResult) []*AppCheckResult {
	nowResultMap := make(map[string]*AppCheckResult)

	for index, nowResult := range newResults {
		nowResultMap[nowResult.Title] = newResults[index]
	}

	for index, oldResult := range oldResults {
		if value, ok := nowResultMap[oldResult.Title]; ok {
			oldResults[index] = value
		}
	}

	return oldResults
}

// isAllowed check all the result is valid
func isAllowed(results []*AppCheckResult) bool {
	for _, result := range results {
		if !result.IsResultValid() {
			return false
		}
	}
	return true
}

// IsAllowTest test check results whether to allow test
func (r *AppCheckTestResult) IsAllowTest() bool {
	return isAllowed(r.Results)
}

// updateAppCheckResults update only the old data
func updateAppCheckResults(oldResults []*AppCheckResult, newResults []*AppCheckResult) []*AppCheckResult {
	if len(newResults) == 0 {
		return oldResults
	}

	newResultMap := make(map[string]*AppCheckResult)
	for index, nowResult := range newResults {
		newResultMap[nowResult.Title] = newResults[index]
	}

	for index, oldResult := range oldResults {
		if value, ok := newResultMap[oldResult.Title]; ok {
			oldResults[index] = value
		}
	}

	return oldResults
}

// UpdateCurrentResult update current app test result, only update old check results
func (r *AppCheckTestResult) UpdateCurrentResult(newResults []*AppCheckResult) {
	r.Results = updateAppCheckResults(r.Results, newResults)
}

// ParseToAppExtend check test result is allow test, update application extend
// isUpdateToAllowStatus = true, if check result is allow, update app extend is_allow_test update to allow, reason(array type) update to default JSON('[]')
func (r *AppCheckTestResult) ParseToAppExtend(appExtend *models.ApplicationExtend, isUpdateToAllowStatus bool) error {
	if r.IsAllowTest() && isUpdateToAllowStatus {
		appExtend.IsAllowTest = models.AppAllowedStatus
		appExtend.NotAllowTestReason = utils.DefaultJSONValue
	} else {
		reason, err := utils.NewJSON(r)
		if err != nil {
			return err
		}
		appExtend.IsAllowTest = models.AppNotAllowStatus
		appExtend.NotAllowTestReason = reason
	}
	return nil
}

// IsAllowSubmit submit check results whether to allow submit
func (r *AppCheckSubmitResult) IsAllowSubmit() bool {
	return isAllowed(r.Results)
}

// UpdateCurrentResult update current app submit result, only update old check results
func (r *AppCheckSubmitResult) UpdateCurrentResult(newResults []*AppCheckResult) {
	r.Results = updateAppCheckResults(r.Results, newResults)
}

// ParseToAppExtend parse submit result to app extend
// isUpdateToAllowStatus = true, if check result is allow, update app extend is_allow_submit update to allow, reason(array type) update to default JSON('[]')
func (r *AppCheckSubmitResult) ParseToAppExtend(appExtend *models.ApplicationExtend, isUpdateToAllowStatus bool) error {
	if r.IsAllowSubmit() && isUpdateToAllowStatus {
		appExtend.IsAllowSubmit = models.AppAllowedStatus
		appExtend.NotAllowSubmitReason = utils.DefaultJSONValue
	} else {
		reason, err := utils.NewJSON(r)
		if err != nil {
			return err
		}
		appExtend.NotAllowSubmitReason = reason
		appExtend.IsAllowSubmit = models.AppNotAllowStatus
	}
	return nil
}

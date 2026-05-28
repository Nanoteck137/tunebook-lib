package library

type Report struct {
	File      string
	Err       error
	IsWarning bool
}

type Reporter struct {
	Errors      map[string][]Report
	NumErrors   int
	NumWarnings int
}

func (r *Reporter) AddReport(report Report) {
	file := report.File

	errs, exists := r.Errors[file]
	if !exists {
		r.Errors[file] = []Report{report}
	} else {
		errs = append(errs, report)
		r.Errors[file] = errs
	}
}

func (r *Reporter) AddError(file string, err error) {
	r.AddReport(Report{
		File:      file,
		Err:       err,
		IsWarning: false,
	})

	r.NumErrors++
}

func (r *Reporter) AddWarning(file string, err error) {
	r.AddReport(Report{
		File:      file,
		Err:       err,
		IsWarning: true,
	})

	r.NumWarnings++
}

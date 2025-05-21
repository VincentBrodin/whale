package stage

func (r *Stage) SetStyle(style string) error {
	return r.Write(style)
}

func (r *Stage) ResetStyle() error {
	return r.Write(ResetCode)
}

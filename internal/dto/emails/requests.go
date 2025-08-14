package dto

type SendRequest struct {
	FromEmail string `json:"from_email"`
	FromName  string `json:"from_name"`
	Subject   string `json:"subject"`
	ToEmail   string `json:"to_email"`
	Html      string `json:"html"`
	Text      string `json:"text"`
}

func (r *SendRequest) Validate() error {
	if err := r.ValidateFromEmail(); err != nil {
		return err
	}
	if err := r.ValidateFromName(); err != nil {
		return err
	}
	if err := r.ValidateSubject(); err != nil {
		return err
	}
	if err := r.ValidateToEmail(); err != nil {
		return err
	}
	if err := r.ValidateHtml(); err != nil {
		return err
	}
	if err := r.ValidateText(); err != nil {
		return err
	}
	return nil
}

func (r *SendRequest) ValidateFromEmail() error {
	if r.FromEmail == "" {
		return ErrInvalidFromEmail
	}
	return nil
}

func (r *SendRequest) ValidateFromName() error {
	if r.FromName == "" {
		return ErrInvalidFromName
	}
	return nil
}

func (r *SendRequest) ValidateSubject() error {
	if r.Subject == "" {
		return ErrInvalidSubject
	}
	return nil
}

func (r *SendRequest) ValidateToEmail() error {
	if r.ToEmail == "" {
		return ErrInvalidToEmail
	}
	return nil
}

func (r *SendRequest) ValidateHtml() error {
	if r.Html == "" {
		return ErrInvalidHtml
	}
	return nil
}

func (r *SendRequest) ValidateText() error {
	if r.Text == "" {
		return ErrInvalidText
	}
	return nil
}

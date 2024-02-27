package storage

func (a *AuthStorage) CreateUserGoogle(email, username string) error {
	query := `INSERT INTO user(email, username, password) VALUES ($1, $2, $3);`
	_, err := a.db.Exec(query, email, username, "google")
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthStorage) CreateUserGithub(email, username string) error {
	query := `INSERT INTO user(email, username, password) VALUES ($1, $2, $3);`
	_, err := a.db.Exec(query, email, username, "github")
	if err != nil {
		return err
	}
	return nil
}

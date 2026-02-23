package users

type SessionsRepository struct {
	Sessions []*Session
}

func NewSessionsRepository() *SessionsRepository {
	return &SessionsRepository{
		Sessions: []*Session{},
	}
}

func (r *SessionsRepository) Add(session *Session) {
	r.Sessions = append(r.Sessions, session)
}

func (r *SessionsRepository) FindByToken(token string) *Session {
	for _, session := range r.Sessions {
		if session.Token == token {
			return session
		}
	}
	return nil
}

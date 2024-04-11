package database

func (u *User) IsActive() bool {
	if u.Status == AccountStatusDeleted || u.Status == AccountStatusSuspended {
		return false
	}

	return true
}

package models

type ManagementRole struct {
	TableName struct{} `sql:"management_role" pg:",discard_unknown_colums"`
	Id        int64    `sql:",pk"`
	RoleName  string   `sql:"role"`
}

type ManagementUser struct {
	TableName        struct{}        `sql:"management_user"`
	Id               int64           `sql:"management_user_id"`
	Email            string          `sql:"email"`
	Password         string          `sql:"password"`
	ManagementRole   *ManagementRole `pg:"fk:management_role_id"`
	ManagementRoleId int64           `sql:"management_role_id"`
	Verified         bool            `sql:"verified"`
	Enabled          bool            `sql:"enabled"`
}

type UserRegister struct {
	TableName struct{} `sql:"management_user"`
	Id        int64    `sql:"management_user_id"`
	Email     string   `sql:"email"`
	Password  string   `sql:"password"`
}

func (db *Db) GetUser(email string, userId int64) (ManagementUser, error) {
	var user ManagementUser
	err := db.Client.Model(&user).
		Column("management_user.*", "ManagementRole").
		Where("email=? AND management_user_id=?", email, userId).
		Limit(1).
		Select()
	return user, err
}

func (db *Db) GetManagerByEmail(email string) (ManagementUser, error) {
	var user ManagementUser
	err := db.Client.Model(&user).
		Column("management_user.*", "ManagementRole").
		Where("email=?", email).
		Limit(1).
		Select()
	return user, err
}

func (db *Db) RegisterUser(user *UserRegister) (*UserRegister, error) {
	_, err := db.Client.Model(user).
		Column("management_user_id").
		Where("email = ?", user.Email).
		Returning("management_user_id").
		SelectOrInsert()
	return user, err
}

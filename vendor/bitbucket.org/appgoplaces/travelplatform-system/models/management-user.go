package models

import (

)

type ManagementUser struct {
    tableName struct{} `sql:"management_user"`
    Id int64 `sql:"management_user_id"`
    Email string `sql:"email"`
    Password string `sql:"password"`
    Role string `sql:"role"`
}

type UserRegister struct {
    tableName struct{} `sql:"management_user"`
    Id int64 `sql:"management_user_id"`
    Email string `sql:"email"`
    Password string `sql:"password"`
}

func (db *Db) GetManagerByEmail(email string) (ManagementUser, error) {
    var user ManagementUser
    _, err := db.Client.Query(&user, `SELECT * FROM public.management_user WHERE email = ?`, email)
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

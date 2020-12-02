package main

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

// 我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，
// 是否应该 Wrap 这个 error，抛给上层。为什么？应该怎么做请写出代码

var (
	DataNotExist = errors.New("Data Not find")
)

type User struct {
	ID   int
	Name string
}

type QueryParameter struct {
	// ...
}

// model
func mockDBError() ([]User, error) {
	return nil, sql.ErrNoRows
}

func MockQueryUserList(param *QueryParameter) ([]User, error) {
	var users []User
	var err error
	if users, err = mockDBError(); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = DataNotExist
		}

		return users, errors.Wrap(err, "model: query user list meet error")
	}

	return users, err
}

// service
func QueryUserList(param *QueryParameter) ([]User, error) {
	var users []User
	var err error
	if users, err := MockQueryUserList(param); err != nil {
		return users, errors.WithMessage(err, "service: query user list meet error")
	}

	return users, err
}

func main() {
	param := &QueryParameter{}
	users, err := QueryUserList(param)
	if err != nil {
		if errors.Is(err, DataNotExist) {
			fmt.Printf("users not exists: %+v\n", err)
			return
		}
		fmt.Printf("query user list meet error: %+v\n", err)
		return
	}

	fmt.Printf("query user list: %++v\n", users)
}

package userserver

import (
	"context"
	"regexp"
	"time"

	errorpkg "github.com/merefield/grpc-user-api/pkg/error"
	uuidpkg "github.com/merefield/grpc-user-api/pkg/uuid"

	"github.com/go-pg/pg"
	"github.com/twitchtv/twirp"

	"github.com/merefield/grpc-user-api/internal/model"
	pbUser "github.com/merefield/grpc-user-api/proto/user"
)

// Server implements the UserService
type Server struct {
	db  *pg.DB
	sec Securer
}

// New creates an instance of our server
func New(db *pg.DB, sec Securer) *Server {
	return &Server{db: db, sec: sec}
}

// Securer represents password securing service
type Securer interface {
	Hash(string) string
	Password(string, ...string) bool
}

// AddUser gets a user to the in-memory store.
func (s *Server) AddUser(ctx context.Context, user *pbUser.AddUserRequest) (*pbUser.Empty, error) {

	requiredErr := checkRequiredAddAttributes(user)
	if requiredErr != nil {
		return nil, requiredErr
	}

	newUser := &model.User{
		Username: user.Username,
		FullName: user.FullName,
		Email:    user.Email,
		// DisplayName: user.DisplayName,
	}
	_, err := s.db.Model(newUser).Returning("*").Insert()

	returnErr := errorpkg.CheckError(err, "user")
	if returnErr != nil {
		return nil, returnErr
	}

	return &pbUser.Empty{}, nil

	// return &pbUser.User{
	// 	Id:       newUser.ID.String(),
	// 	Username: newUser.Username,
	// 	// DisplayName: newUser.DisplayName,
	// 	FullName: newUser.FullName,
	// 	Email:    newUser.Email,
	// }, nil
}

// GetUser Gets a user from the DB
func (s *Server) GetUser(ctx context.Context, user *pbUser.UserRequest) (*pbUser.UserPublicResponse, error) {
	u, err := getUserModelFromID(user.Id)
	if err != nil {
		return nil, err
	}

	pgerr := s.db.Model(u).
		Column("user.*", "OwnerOfGroups").
		Where("id = ?", u.ID).
		Select()
	twerr := errorpkg.CheckError(pgerr, "user")
	if twerr != nil {
		return nil, twerr
	}

	return &pbUser.UserPublicResponse{
		Username:       u.Username,
		FullName:       u.FullName,
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		Member:         u.Member,
		FollowedGroups: uuidpkg.ConvertUUIDToStrArray(u.FollowedGroups),
		// id:       u.ID.String(),
		//DisplayName: u.DisplayName,
		//NewsletterNotification: u.NewsletterNotification,
		// Email:                  u.Email,
		//OwnerOfGroups:  getUserGroupResponse(u.OwnerOfGroups),
		//FavoriteTracks:         uuidpkg.ConvertUUIDToStrArray(u.FavoriteTracks),
	}, nil
}

// GetUserRestricted intended for privileged roles only supplies more detailed, private info about user.
func (s *Server) GetUserRestricted(ctx context.Context, user *pbUser.UserRequest) (*pbUser.UserPrivateResponse, error) {
	u, err := getUserModelFromID(user.Id)
	if err != nil {
		return nil, err
	}

	pgerr := s.db.Model(u).
		Column("user.*", "OwnerOfGroups").
		Where("id = ?", u.ID).
		Select()
	twerr := errorpkg.CheckError(pgerr, "user")
	if twerr != nil {
		return nil, twerr
	}

	return &pbUser.UserPrivateResponse{
		Id:             u.ID.String(),
		Username:       u.Username,
		FullName:       u.FullName,
		Email:          u.Email,
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		Member:         u.Member,
		RoleId:         u.RoleID,
		TenantId:       u.TenantID,
		FollowedGroups: uuidpkg.ConvertUUIDToStrArray(u.FollowedGroups),
		//DisplayName: u.DisplayName,
		//NewsletterNotification: u.NewsletterNotification,
		//OwnerOfGroups:  getUserGroupResponse(u.OwnerOfGroups),
		//FavoriteTracks:         uuidpkg.ConvertUUIDToStrArray(u.FavoriteTracks),
	}, nil
}

// DeleteUser Deletes a user from the DB
func (s *Server) DeleteUser(ctx context.Context, user *pbUser.UserRequest) (*pbUser.Empty, error) {
	u, err := getUserModelFromID(user.Id)
	if err != nil {
		return nil, err
	}

	_, pgerr := s.db.Model(u).
		Column("user.*").
		Where("id = ?", u.ID).
		Delete()
	twerr := errorpkg.CheckError(pgerr, "user")
	if twerr != nil {
		return nil, twerr
	}

	return &pbUser.Empty{}, nil
}

// UpdateUser gets a user to the in-memory store.
func (s *Server) UpdateUser(ctx context.Context, updateUserRequest *pbUser.UpdateUserRequest) (*pbUser.Empty, error) {
	err := checkRequiredUpdateAttributes(updateUserRequest)

	if err != nil {
		return nil, err
	}

	existingID, err := uuidpkg.GetUUIDFromString(updateUserRequest.Id)

	if err != nil {
		return nil, err
	}

	u := &model.User{
		ID:        existingID,
		Username:  updateUserRequest.Username,
		FullName:  updateUserRequest.FullName,
		Email:     updateUserRequest.Email,
		FirstName: updateUserRequest.FirstName,
		LastName:  updateUserRequest.LastName,
		// DisplayName: user.DisplayName,
	}

	u.UpdatedAt = time.Now()
	_, pgerr := s.db.Model(u).
		Column("updated_at", "username", "full_name", "email", "first_name", "last_name", "member", "newsletter_notification").
		WherePK().
		Returning("*").
		Update()
	twerr := errorpkg.CheckError(pgerr, "user")
	if twerr != nil {
		return nil, twerr
	}
	return &pbUser.Empty{}, nil
}

// UpdateUser gets a user to the in-memory store.
func (s *Server) UpdateUserRestricted(ctx context.Context, updateUserRestrictedRequest *pbUser.UpdateUserRestrictedRequest) (*pbUser.Empty, error) {
	err := checkRequiredRestrictedUpdateAttributes(updateUserRestrictedRequest)

	if err != nil {
		return nil, err
	}

	existingID, err := uuidpkg.GetUUIDFromString(updateUserRestrictedRequest.Id)

	if err != nil {
		return nil, err
	}

	u := &model.User{
		ID:        existingID,
		Username:  updateUserRestrictedRequest.Username,
		FullName:  updateUserRestrictedRequest.FullName,
		Email:     updateUserRestrictedRequest.Email,
		FirstName: updateUserRestrictedRequest.FirstName,
		LastName:  updateUserRestrictedRequest.LastName,
		RoleID:    updateUserRestrictedRequest.RoleId,
		TenantID:  updateUserRestrictedRequest.TenantId,
		// DisplayName: user.DisplayName,
	}

	u.UpdatedAt = time.Now()
	_, pgerr := s.db.Model(u).
		Column("updated_at", "username", "full_name", "email", "first_name", "last_name", "role_id", "tenant_id", "member", "newsletter_notification").
		WherePK().
		Returning("*").
		Update()
	twerr := errorpkg.CheckError(pgerr, "user")
	if twerr != nil {
		return nil, twerr
	}
	return &pbUser.Empty{}, nil
}

// ResetUserPassword reset's a user's password
func (s *Server) ResetUserPassword(ctx context.Context, ResetUserPasswordRequest *pbUser.ResetUserPasswordRequest) (*pbUser.Empty, error) {
	err := checkRequiredResetPasswordAttributes(ResetUserPasswordRequest, s)

	if err != nil {
		return nil, err
	}

	hashedPassword := s.sec.Hash(ResetUserPasswordRequest.Password)

	u := &model.User{
		Email: ResetUserPasswordRequest.Email,
	}

	u.UpdatedAt = time.Now()
	u.Password = hashedPassword
	_, pgerr := s.db.Model(u).
		Column("updated_at", "password").
		Where("email = ?", u.Email).
		Returning("*").
		Update()
	twerr := errorpkg.CheckError(pgerr, "user")
	if twerr != nil {
		return nil, twerr
	}
	return &pbUser.Empty{}, nil
}

// ListUsers lists all users in the store.
func (s *Server) ListUsers(ctx context.Context, Empty *pbUser.Empty) (*pbUser.UserListResponse, error) {

	var users []pbUser.User
	var results pbUser.UserListResponse
	err := s.db.Model(&users).Select()
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		var result pbUser.UserPrivateResponse
		result.Id = user.Id
		result.Username = user.Username
		result.FullName = user.FullName
		result.Email = user.Email
		result.FirstName = user.FirstName
		result.LastName = user.LastName
		result.Member = user.Member
		result.NewsletterNotification = user.NewsletterNotification
		// DisplayName: user.DisplayName,
		results.User = append(results.User, &result)
	}

	return &results, nil
}

func getUserModelFromID(user string) (returneduser *model.User, err error) {
	ID, err := uuidpkg.GetUUIDFromString(user)
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID:                     ID,
		Username:               "",
		FullName:               "",
		Email:                  "",
		FirstName:              "",
		LastName:               "",
		Member:                 false,
		NewsletterNotification: false,
		// DisplayName: user.DisplayName,
	}, nil
}

func getUserModel(user string) (returneduser *model.User, err error) {
	ID, err := uuidpkg.GetUUIDFromString(user)
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID:       ID,
		Username: returneduser.Username,
		// DisplayName: user.DisplayName,
		FullName:               returneduser.FullName,
		Email:                  returneduser.Email,
		FirstName:              returneduser.FirstName,
		LastName:               returneduser.LastName,
		Member:                 returneduser.Member,
		NewsletterNotification: returneduser.NewsletterNotification,
	}, nil
}

func getUserModelforUpdate(user string) (returneduser *model.User, err error) {
	ID, err := uuidpkg.GetUUIDFromString(user)
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID:       ID,
		Username: returneduser.Username,
		// DisplayName: user.DisplayName,
		FullName:               returneduser.FullName,
		Email:                  returneduser.Email,
		FirstName:              returneduser.FirstName,
		LastName:               returneduser.LastName,
		Member:                 returneduser.Member,
		NewsletterNotification: returneduser.NewsletterNotification,
	}, nil
}

func checkRequiredAddAttributes(user *pbUser.AddUserRequest) error {
	if user.Email == "" || user.Username == "" || user.FullName == "" {
		var argument string
		switch {
		case user.Email == "":
			argument = "email"
		case user.Username == "":
			argument = "username"
		case user.FullName == "":
			argument = "full_name"
		}
		return twirp.RequiredArgumentError(argument)
	}
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if re.MatchString(user.Email) == false {
		return twirp.InvalidArgumentError("email", "must be a valid email")
	}
	return nil
}

func checkRequiredUpdateAttributes(user *pbUser.UpdateUserRequest) error {
	if user.Email == "" || user.Username == "" || user.FullName == "" {
		var argument string
		switch {
		case user.Email == "":
			argument = "email"
		case user.Username == "":
			argument = "username"
		case user.FullName == "":
			argument = "full_name"
		}
		return twirp.RequiredArgumentError(argument)
	}
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if re.MatchString(user.Email) == false {
		return twirp.InvalidArgumentError("email", "must be a valid email")
	}
	return nil
}

func checkRequiredRestrictedUpdateAttributes(user *pbUser.UpdateUserRestrictedRequest) error {
	if user.Email == "" || user.Username == "" || user.FullName == "" {
		var argument string
		switch {
		case user.Email == "":
			argument = "email"
		case user.Username == "":
			argument = "username"
		case user.FullName == "":
			argument = "full_name"
		}
		return twirp.RequiredArgumentError(argument)
	}
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if re.MatchString(user.Email) == false {
		return twirp.InvalidArgumentError("email", "must be a valid email")
	}
	return nil
}

func checkRequiredResetPasswordAttributes(user *pbUser.ResetUserPasswordRequest, s *Server) error {
	if user.Email == "" || user.Password == "" {
		var argument string
		switch {
		case user.Email == "":
			argument = "email"
		case user.Password == "":
			argument = "Password"
		}
		return twirp.RequiredArgumentError(argument)
	}
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if re.MatchString(user.Email) == false {
		return twirp.InvalidArgumentError("Email", "must be a valid email")
	}
	if !s.sec.Password(user.Password) {
		return twirp.InvalidArgumentError("Password", "is not strong enough")
	}

	return nil
}

func getUserGroupResponse(ownerOfGroup []model.UserGroup) []*pbUser.RelatedUserGroup {
	groups := make([]*pbUser.RelatedUserGroup, len(ownerOfGroup))
	for i, group := range ownerOfGroup {
		groups[i] = &pbUser.RelatedUserGroup{Id: group.ID.String(), DisplayName: group.DisplayName, Avatar: group.Avatar}
	}
	return groups
}

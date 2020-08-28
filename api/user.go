package api

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

type userInfo struct {
	Name             string   `json:"name"`
	NamePy           string   `json:"name_py"`
	EnName           string   `json:"en_name"`
	EmployeeID       string   `json:"employee_id"`
	EmployeeNo       string   `json:"employee_no"`
	OpenID           string   `json:"open_id"`
	UnionID          string   `json:"union_id"`
	Status           int      `json:"status"`
	EmployeeType     int      `json:"employee_type"`
	Avatar72         string   `json:"avatar_72"`
	Avatar240        string   `json:"avatar_240"`
	Avatar640        string   `json:"avatar_640"`
	AvatarURL        string   `json:"avatar_url"`
	Gender           int      `json:"gender"`
	Email            string   `json:"email"`
	Mobile           string   `json:"mobile"`
	Description      string   `json:"description"`
	Country          string   `json:"country"`
	City             string   `json:"city"`
	WorkStation      string   `json:"work_station"`
	IsTenantManager  bool     `json:"is_tenant_manager"`
	JoinTime         int      `json:"join_time"`
	UpdateTime       int      `json:"update_time"`
	LeaderEmployeeID string   `json:"leader_employee_id"`
	LeaderOpenID     string   `json:"leader_open_id"`
	LeaderUnionID    string   `json:"leader_union_id"`
	Departments      []string `json:"departments"`
	OpenDepartments  []string `json:"open_departments"`
}

func (u *userInfo) IsAdmin() bool {
	return u.IsTenantManager
}

type UserInfosResponse struct {
	BaseResponse
	UserInfos []*userInfo `json:"user_infos"`
}

type User struct {
	app         *App
	accessToken string
}

func NewUser(app *App, accessToken string) *User {
	return &User{
		app:         app,
		accessToken: accessToken,
	}
}

func (u *User) Get(openID string) (*userInfo, error) {
	infos, err := u.List(openID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	length := len(infos.UserInfos)
	if length > 0 {
		return infos.UserInfos[0], nil
	}
	return nil, nil
}

func (u *User) List(openIDs ...string) (resp *UserInfosResponse, err error) {
	length := len(openIDs)
	if length == 0 {
		return nil, nil
	}
	if length > 100 {
		return nil, errors.New("openIDs count max 100")
	}

	params := make([]string, len(openIDs))
	for i := range openIDs {
		params[i] = "open_ids=" + openIDs[i]
	}

	path := fmt.Sprintf("contact/v1/user/batch_get?%s", strings.Join(params, "&"))

	resp = new(UserInfosResponse)
	_, err = NewClientWithBearer(u.app, u.accessToken).Get(path, resp)
	if err != nil {
		err = errors.WithStack(err)
	}
	if !resp.OK() {
		return nil, resp.Error()
	}
	return
}

type AppAdminResponse struct {
	BaseResponse
	Data struct {
		IsAppAdmin bool `json:"is_app_admin"`
	} `json:"data"`
}

func (u *User) IsAppAdmin(openID string) (bool, error) {
	resp := new(AppAdminResponse)

	query := url.Values{
		// "Authorization": []string{u.accessToken},
		"open_id": []string{openID},
	}
	_, err := NewClientWithBearer(u.app, u.accessToken).Get("application/v3/is_user_admin?"+query.Encode(), resp)
	if err != nil {
		return false, errors.WithStack(err)
	}
	if !resp.OK() {
		return false, resp.Error()
	}
	return resp.Data.IsAppAdmin, nil
}

type AdminListResponse struct {
	BaseResponse
	Data struct {
		UserList []struct {
			OpenID string `json:"open_id"`
		} `json:"user_list"`
	} `json:"data"`
}

func (u *User) AdminList() (resp *AdminListResponse, err error) {
	resp = new(AdminListResponse)
	_, err = NewClientWithBearer(u.app, u.accessToken).Get("user/v4/app_admin_user/list", resp)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if !resp.OK() {
		return nil, resp.Error()
	}
	return resp, nil
}

type AccessTokenResponse struct {
	BaseResponse
	Data struct {
		AccessToken      string `json:"access_token"`
		AvatarURL        string `json:"avatar_url"`
		AvatarThumb      string `json:"avatar_thumb"`
		AvatarMiddle     string `json:"avatar_middle"`
		AvatarBig        string `json:"avatar_big"`
		ExpiresIn        int    `json:"expires_in"`
		Name             string `json:"name"`
		EnName           string `json:"en_name"`
		OpenID           string `json:"open_id"`
		TenantKey        string `json:"tenant_key"`
		RefreshExpiresIn int    `json:"refresh_expires_in"`
		RefreshToken     string `json:"refresh_token"`
		TokenType        string `json:"token_type"`
	} `json:"data"`
}

func (u *User) AccessToken(code string) (resp *AccessTokenResponse, err error) {
	body := map[string]interface{}{
		"app_access_token": u.accessToken,
		"grant_type":       "authorization_code",
		"code":             code,
	}
	resp = new(AccessTokenResponse)
	_, err = NewClient(u.app).Post("authen/v1/access_token", body, resp)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	if !resp.OK() {
		return nil, resp.Error()
	}

	return resp, nil
}

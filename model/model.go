package model

import "encoding/json"

type Entity struct {
	EntityId     string `json:"id"`
	Title        string `json:"title"`
	Organization string
	Type         EntityType
}

type User struct {
	UserId      *string `json:"uid"`
	Username    *string `json:"username"`
	Displayname *string `json:"displayname"`
	Userslug    *string `json:"userslug"`
	Picture     *string `json:"picture"`
	IconText    *string `json:"icon:text"`
	IconBgColor *string `json:"icon:bgColor"`
}

type Thread struct {
	ThreadId   *string `json:"tid"`
	Title      *string `json:"title"`
	Posts      []*Post `json:"posts"`
	Timestamp  *int    `json:"timestamp"`
	Content    *string
	Pagination *Pagination `json:"pagination"`
}

type Post struct {
	PostId    *string `json:"pid"`
	UserId    *string `json:"uid"`
	ThreadId  *string `json:"tid"`
	Content   *string `json:"content"`
	ToPostId  *string `json:"toPid"`
	Timestamp *int    `json:"timestamp"`
	Deleted   *bool   `json:"deleted"`
	UserInfo  *User   `json:"user"`
}

type ThreadResponsePayload struct {
	ThreadData *ThreadDTO `json:"topicData"`
	PostData   *PostDTO   `json:"postData"`
}

type PostResponseDAO struct {
	Code    *string  `json:"code"`
	Payload *PostDTO `json:"payload"`
}

type ThreadResponseDTO struct {
	Code    *string                `json:"code"`
	Payload *ThreadResponsePayload `json:"payload"`
}
type UserDTO struct {
	UserId      *json.Number `json:"uid"`
	Username    *string      `json:"username"`
	Displayname *string      `json:"displayname"`
	Userslug    *string      `json:"userslug"`
	Picture     *string      `json:"picture"`
	IconText    *string      `json:"icon:text"`
	IconBgColor *string      `json:"icon:bgColor"`
}

type ThreadDTO struct {
	ThreadId   *json.Number   `json:"tid"`
	Title      *string        `json:"title"`
	Posts      []*PostDTO     `json:"posts"`
	Timestamp  *json.Number   `json:"timestamp"`
	Pagination *PaginationDTO `json:"pagination"`
	PostCount  *json.Number   `json:"postcount"`
}

type PostDTO struct {
	PostId    *json.Number `json:"pid"`
	UserId    *json.Number `json:"uid"`
	ThreadId  *json.Number `json:"tid"`
	Content   *string      `json:"content"`
	ToPostId  *string      `json:"toPid"`
	Timestamp *json.Number `json:"timestamp"`
	Deleted   *json.Number `json:"deleted"`
	UserInfo  *UserDTO     `json:"user"`
}

type PaginationDTO struct {
	CurrentPage *json.Number `json:"currentPage"`
	PageCount   *json.Number `json:"pageCount"`
}

type Pagination struct {
	CurrentPage *int `json:"currentPage"`
	PageCount   *int `json:"pageCount"`
	TotalPosts  *int `json:"totalPosts"`
}

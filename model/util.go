package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/Informasjonsforvaltning/fdk-user-feedback-service/env"
)

const threadTitleTemplate = "%s %s, %s"
const threadLinkedContentTemplate = "Dette er en automatisk opprettet kommentartråd for %s [%s](%s)."
const threadContentTemplate = "Dette er en automatisk opprettet kommentartråd for %s %s."

func fdkLink(entityType EntityType, entityId string) *string {
	var path string
	entityBasePath := entityType.ToPath()
	if entityBasePath == nil || entityId == "" {
		return nil
	}
	path = env.EnvironmentVariables.FdkBaseUri + "/" + *entityBasePath + entityId
	return &path
}

func (thread *Thread) FindThreadPostById(postId *string) (*Post, error) {
	if thread == nil {
		return nil, nil
	}

	var post *Post
	var err error
	for _, element := range thread.Posts {
		if element.PostId != nil && postId != nil && *element.PostId == *postId {
			post = element
			break
		}
	}

	if post == nil {
		err = errors.New("no post with specified id found")
	}

	return post, err
}

func (thread *Thread) FilterDeletedPosts() *Thread {
	if thread == nil {
		return nil
	}

	var posts []*Post
	for _, element := range thread.Posts {
		if element.Deleted == nil || !*element.Deleted {
			posts = append(posts, element)
		}
	}

	return &Thread{
		ThreadId:   thread.ThreadId,
		Title:      thread.Title,
		Posts:      posts,
		Timestamp:  thread.Timestamp,
		Content:    thread.Content,
		Pagination: thread.Pagination,
	}
}

func (threadDto *ThreadDTO) ToThread() *Thread {
	if threadDto == nil {
		return nil
	}

	threadId := NumberPointerToStringPointer(threadDto.ThreadId)
	title := threadDto.Title
	timestamp := NumberPointerToIntPointer(threadDto.Timestamp)
	pagination := ToPagination(threadDto)

	var posts []*Post
	for _, postDto := range threadDto.Posts {
		posts = append(posts, postDto.ToPost())
	}

	return &Thread{
		ThreadId:   threadId,
		Title:      title,
		Posts:      posts,
		Timestamp:  timestamp,
		Pagination: pagination,
	}
}

func ToPagination(threadDto *ThreadDTO) *Pagination {
	if threadDto == nil {
		return nil
	}

	paginationDto := threadDto.Pagination
	if paginationDto == nil {
		return nil
	}

	currentPage := NumberPointerToIntPointer(paginationDto.CurrentPage)
	pageCount := NumberPointerToIntPointer(paginationDto.PageCount)
	totalPosts := NumberPointerToIntPointer(threadDto.PostCount)

	return &Pagination{
		CurrentPage: currentPage,
		PageCount:   pageCount,
		TotalPosts:  totalPosts,
	}
}

func (entity *Entity) ToThread() (*Thread, error) {
	if entity == nil {
		return nil, fmt.Errorf("no entity found")
	}
	if entity.Title == "" {
		return nil, fmt.Errorf("no title for entity of id: %s", entity.EntityId)
	}
	entityLink := fdkLink(entity.Type, entity.EntityId)
	if entityLink == nil {
		title := fmt.Sprintf(threadTitleTemplate, entity.Type.StringNb(), entity.Title, entity.Organization)
		content := fmt.Sprintf(threadContentTemplate, entity.Type.StringNbPlural(), entity.Title)
		return &Thread{
			Title:   &title,
			Content: &content,
		}, nil
	}

	title := fmt.Sprintf(threadTitleTemplate, entity.Type.StringNb(), entity.Title, entity.Organization)
	content := fmt.Sprintf(threadLinkedContentTemplate, entity.Type.StringNbPlural(), entity.Title, *entityLink)
	return &Thread{
		Title:   &title,
		Content: &content,
	}, nil
}

func (postDto *PostDTO) ToPost() *Post {
	if postDto == nil {
		return nil
	}

	postId := NumberPointerToStringPointer(postDto.PostId)
	userId := NumberPointerToStringPointer(postDto.UserId)
	threadId := NumberPointerToStringPointer(postDto.ThreadId)
	index := NumberPointerToStringPointer(postDto.Index)
	timestamp := NumberPointerToIntPointer(postDto.Timestamp)
	deletedVal := NumberPointerToStringPointer(postDto.Deleted)
	deleted := deletedVal != nil && *deletedVal == "1"
	toPostId := postDto.ToPostId
	return &Post{
		PostId:    postId,
		UserId:    userId,
		ThreadId:  threadId,
		Index:     index,
		Content:   postDto.Content,
		ToPostId:  toPostId,
		Timestamp: timestamp,
		Deleted:   &deleted,
		UserInfo:  postDto.UserInfo.ToUser(),
	}
}

func (userDto *UserDTO) ToUser() *User {
	if userDto == nil {
		return nil
	}

	return &User{
		UserId:      NumberPointerToStringPointer(userDto.UserId),
		Username:    userDto.Username,
		Displayname: userDto.Displayname,
		Userslug:    userDto.Userslug,
		Picture:     userDto.Picture,
		IconText:    userDto.IconText,
		IconBgColor: userDto.IconBgColor,
	}
}

func StringPointerToIntPointer(s *string) (*int, error) {
	if s == nil {
		return nil, nil
	}

	intVal, err := strconv.Atoi(*s)
	if err != nil {
		return nil, err
	}

	return &intVal, nil
}

func NumberPointerToStringPointer(num *json.Number) *string {
	if num == nil {
		return nil
	}

	strVal := num.String()

	return &strVal
}

func NumberPointerToIntPointer(num *json.Number) *int {
	if num == nil {
		return nil
	}

	int64Val, err := num.Int64()
	if err != nil {
		return nil
	}

	intVal := int(int64Val)

	return &intVal
}

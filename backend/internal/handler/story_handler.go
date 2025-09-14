package handler

import (

)

type IStoryHandler interface {
	// GenerateStory() error
	// GetStories() error
	// GetStory() error
	// DeleteStory() error
}

type StoryHandler struct {

}

func NewStoryHandler() IStoryHandler {
	return &StoryHandler{}
}
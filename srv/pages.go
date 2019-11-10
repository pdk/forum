package srv

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pdk/forum/model"
	"github.com/pdk/forum/store"
)

// HomePage returns the front page of the application.
func (s Server) HomePage(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		// default/included mux is pretty dumb. just gonna throw this in here,
		// rather than upgrade to a better mux.
		// https://github.com/golang/go/issues/4799
		http.NotFound(w, r)
		return
	}

	s.WritePage(w, "home.html", nil)
}

// SignIn handles new user sign in.
func (s Server) SignIn(w http.ResponseWriter, r *http.Request) {

	name := r.FormValue("name")

	user, err := store.GetOrCreateUserByName(s.DB, name)
	if handleError(w, "cannot find/create user %s: %w", name, err) {
		return
	}

	setSignedInUserName(w, user.Name)

	s.WritePage(w, "welcome.html", map[string]string{
		"name": name,
	})
}

// TopicsPage shows the list of available topics.
func (s Server) TopicsPage(w http.ResponseWriter, r *http.Request) {

	topicList, err := store.QueryTopics(s.DB)
	if handleError(w, "cannot get list of topics: %w", err) {
		return
	}

	s.WritePage(w, "topics.html", map[string]interface{}{
		"topics": topicList,
	})
}

// OneTopicPage shows the threads within one topic.
func (s Server) OneTopicPage(w http.ResponseWriter, r *http.Request) {

	topicID, err := getPathID(r.URL)
	if handleError(w, "cannot identify topic id: %w", err) {
		return
	}

	topic, err := store.GetTopicByID(s.DB, topicID)
	if errorNotFound(w, r, err) || handleError(w, "cannot get topic %d: %w", topicID, err) {
		return
	}

	threads, err := store.QueryThreadsByTopicID(s.DB, topic.ID)
	if handleError(w, "cannot get threads for topic %d: %w", topic.ID, err) {
		return
	}

	s.WritePage(w, "threads.html", map[string]interface{}{
		"topic":   topic,
		"threads": threads,
	})
}

// AddTopic adds a new topic.
func (s Server) AddTopic(w http.ResponseWriter, r *http.Request) {

	topicName := strings.TrimSpace(r.FormValue("name"))
	if s.MaybeUserError(w, len(topicName) == 0, "new topic name must not be blank") {
		return
	}

	user, err := CurrentUser(s.DB, r)
	if handleError(w, "cannot get current user", err) {
		return
	}

	topic := model.NewTopic(user.ID, topicName)
	topic, err = store.CreateTopic(s.DB, topic)
	if handleError(w, "cannot create topic: %w", err) {
		return
	}

	s.WritePage(w, "new-topic.html", map[string]interface{}{
		"topic": topic,
	})
}

// AddPost adds a post to a threed.
func (s Server) AddPost(w http.ResponseWriter, r *http.Request) {

	user, err := CurrentUser(s.DB, r)
	if handleError(w, "cannot get current user", err) {
		return
	}

	body := strings.TrimSpace(r.FormValue("body"))
	if s.MaybeUserError(w, body == "", "Cannot post with blank comment.") {
		return
	}

	threadIDString := r.FormValue("threadID")
	threadID, err := strconv.ParseInt(threadIDString, 10, 64)
	if handleError(w, "cannot parse thread id %s: %w", threadIDString, err) {
		return
	}

	thread, err := store.GetThreadByID(s.DB, threadID)
	if handleError(w, "cannot get thread %d: %w", threadID, err) {
		return
	}

	post := model.NewPost(threadID, user.ID, body)
	post, err = store.CreatePost(s.DB, post)
	if handleError(w, "cannot save new post: %w", err) {
		return
	}

	s.WritePage(w, "new-post.html", map[string]interface{}{
		"thread": thread,
		"post":   post,
	})
}

// AddThread adds a new topic.
func (s Server) AddThread(w http.ResponseWriter, r *http.Request) {

	user, err := CurrentUser(s.DB, r)
	if handleError(w, "cannot get current user", err) {
		return
	}

	topicIDString := r.FormValue("topicID")
	topicID, err := strconv.ParseInt(topicIDString, 10, 64)
	if handleError(w, "cannot parse topic id %s: %w", topicIDString, err) {
		return
	}

	subject := strings.TrimSpace(r.FormValue("subject"))
	body := strings.TrimSpace(r.FormValue("body"))
	if s.MaybeUserError(w, subject == "" || body == "", "To create a thread, both subject and comments must be non-blank.") {
		return
	}

	topic, err := store.GetTopicByID(s.DB, topicID)
	if handleError(w, "cannot get topic to create new thread: %w", err) {
		return
	}

	thread := model.NewThread(topic.ID, user.ID, subject)
	thread, err = store.CreateThread(s.DB, thread)
	if handleError(w, "cannot save new thread: %w", err) {
		return
	}

	post := model.NewPost(thread.ID, user.ID, body)
	post, err = store.CreatePost(s.DB, post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.WritePage(w, "new-thread.html", map[string]interface{}{
		"thread": thread,
		"post":   post,
	})
}

type displayPost struct {
	Body     template.HTML
	UserName string
	PostedAt time.Time
}

// OneThreadPage shows the comments within one thread.
func (s Server) OneThreadPage(w http.ResponseWriter, r *http.Request) {

	threadID, err := getPathID(r.URL)
	if handleError(w, "cannot get thread id: %w", err) {
		return
	}

	thread, err := store.GetThreadByID(s.DB, threadID)
	if handleError(w, "cannot query thread %d: %w", threadID, err) {
		return
	}

	topic, err := store.GetTopicByID(s.DB, thread.TopicID)
	if handleError(w, "cannot query topic %d: %w", thread.TopicID, err) {
		return
	}

	posts, err := store.QueryPostsByThreadID(s.DB, thread.ID)
	if handleError(w, "cannot query posts for thread %d: %w", thread.ID, err) {
		return
	}

	displayPosts := []displayPost{}
	for _, post := range posts {

		user, err := store.GetUserByID(s.DB, post.PostedByID)
		if handleError(w, "cannot get user %d: %w", post.PostedByID, err) {
			return
		}

		displayPosts = append(displayPosts,
			displayPost{
				Body:     bodyAsHTML(post.Body),
				UserName: user.Name,
				PostedAt: post.PostedAt,
			})
	}

	s.WritePage(w, "one-thread.html", map[string]interface{}{
		"topic":  topic,
		"thread": thread,
		"posts":  displayPosts,
	})
}

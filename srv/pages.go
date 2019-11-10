package srv

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log"
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

	err := s.Template.ExecuteTemplate(w, "home.html", nil)
	if err != nil {
		log.Printf("error executing template home.html: %w", err)
	}
}

// SignIn handles new user sign in.
func (s Server) SignIn(w http.ResponseWriter, r *http.Request) {

	name := r.FormValue("name")

	user, err := store.GetOrCreateUserByName(s.DB, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("created user %d, %s, %s", user.ID, user.Name, user.JoinedAt.String())

	setSignedInUserName(w, name)

	err = s.Template.ExecuteTemplate(w, "welcome.html", map[string]string{
		"name": name,
	})
	if err != nil {
		log.Printf("error executing template welcome.html: %w", err)
	}
}

// TopicsPage shows the list of available topics.
func (s Server) TopicsPage(w http.ResponseWriter, r *http.Request) {

	topicList, err := store.QueryTopics(s.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = s.Template.ExecuteTemplate(w, "topics.html", map[string]interface{}{
		"topics": topicList,
	})
	if err != nil {
		log.Printf("error executing template topics.html: %w", err)
	}
}

// OneTopicPage shows the threads within one topic.
func (s Server) OneTopicPage(w http.ResponseWriter, r *http.Request) {

	topicID, err := getPathID(r.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	topic, err := store.GetTopicByID(s.DB, topicID)
	if errors.Is(err, sql.ErrNoRows) {
		http.NotFound(w, r)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	threads, err := store.QueryThreadsByTopicID(s.DB, topic.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = s.Template.ExecuteTemplate(w, "threads.html", map[string]interface{}{
		"topic":   topic,
		"threads": threads,
	})
	if err != nil {
		log.Printf("error executing template threads.html: %w", err)
	}
}

// AddTopic adds a new topic.
func (s Server) AddTopic(w http.ResponseWriter, r *http.Request) {

	topicName := strings.TrimSpace(r.FormValue("name"))
	if len(topicName) == 0 {
		s.ErrorPage(w, "new topic name cannot be blank")
		return
	}

	user, err := CurrentUser(s.DB, r)

	topic := model.NewTopic(user.ID, topicName)
	topic, err = store.CreateTopic(s.DB, topic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = s.Template.ExecuteTemplate(w, "new-topic.html", map[string]interface{}{
		"topic": topic,
	})
	if err != nil {
		log.Printf("error executing template new-topic.html: %w", err)
	}
}

// AddPost adds a post to a threed.
func (s Server) AddPost(w http.ResponseWriter, r *http.Request) {

	user, err := CurrentUser(s.DB, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	threadIDString := r.FormValue("threadID")
	body := strings.TrimSpace(r.FormValue("body"))

	threadID, err := strconv.ParseInt(threadIDString, 10, 64)
	if err != nil {
		err = fmt.Errorf("cannot parse threadID %s: %w", threadIDString, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if body == "" {
		s.ErrorPage(w, "Cannot post comments without any comments.")
		return
	}

	thread, err := store.GetThreadByID(s.DB, threadID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	post := model.NewPost(threadID, user.ID, body)
	post, err = store.CreatePost(s.DB, post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = s.Template.ExecuteTemplate(w, "new-post.html", map[string]interface{}{
		"thread": thread,
		"post":   post,
	})
	if err != nil {
		log.Printf("error executing template new-thread.html: %w", err)
	}
}

// AddThread adds a new topic.
func (s Server) AddThread(w http.ResponseWriter, r *http.Request) {

	user, err := CurrentUser(s.DB, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	topicIDString := r.FormValue("topicID")
	subject := strings.TrimSpace(r.FormValue("subject"))
	body := strings.TrimSpace(r.FormValue("body"))

	topicID, err := strconv.ParseInt(topicIDString, 10, 64)
	if err != nil {
		err = fmt.Errorf("cannot parse topicID %s: %w", topicIDString, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if subject == "" || body == "" {
		s.ErrorPage(w, "To create a thread, both subject and comments are required.")
		return
	}

	topic, err := store.GetTopicByID(s.DB, topicID)
	if err != nil {
		err = fmt.Errorf("cannot get topic to create new thread: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	thread := model.NewThread(topic.ID, user.ID, subject)
	thread, err = store.CreateThread(s.DB, thread)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	post := model.NewPost(thread.ID, user.ID, body)
	post, err = store.CreatePost(s.DB, post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = s.Template.ExecuteTemplate(w, "new-thread.html", map[string]interface{}{
		"thread": thread,
		"post":   post,
	})
	if err != nil {
		log.Printf("error executing template new-thread.html: %w", err)
	}
}

type displayPost struct {
	Body     template.HTML
	UserName string
	PostedAt time.Time
}

// OneThreadPage shows the comments within one thread.
func (s Server) OneThreadPage(w http.ResponseWriter, r *http.Request) {

	threadID, err := getPathID(r.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	thread, err := store.GetThreadByID(s.DB, threadID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	topic, err := store.GetTopicByID(s.DB, thread.TopicID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	posts, err := store.QueryPostsByThreadID(s.DB, thread.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	displayPosts := []displayPost{}
	for _, post := range posts {

		user, err := store.GetUserByID(s.DB, post.PostedByID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		displayPosts = append(displayPosts,
			displayPost{
				Body:     bodyAsHTML(post.Body),
				UserName: user.Name,
				PostedAt: post.PostedAt,
			})
	}

	err = s.Template.ExecuteTemplate(w, "one-thread.html", map[string]interface{}{
		"topic":  topic,
		"thread": thread,
		"posts":  displayPosts,
	})
	if err != nil {
		log.Printf("error executing template one-thread.html: %w", err)
	}
}

// bodyAsHTML is a hacky solution to splitting text into paragraphs and maintaining line breaks.
func bodyAsHTML(body string) template.HTML {

	body = strings.ReplaceAll(body, "\r\n", "\n")
	body = strings.ReplaceAll(body, "\r", "\n")
	paragraphs := strings.Split(body, "\n\n")

	sb := strings.Builder{}

	for _, para := range paragraphs {
		para = strings.TrimSpace(para)
		if para == "" {
			continue
		}

		sb.WriteString("<p>\n")
		para = template.HTMLEscapeString(para)
		para = strings.ReplaceAll(para, "\n", "<br>\n")
		sb.WriteString(para)
		sb.WriteString("\n</p>\n")
	}

	return template.HTML(sb.String())
}

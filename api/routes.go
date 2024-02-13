package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/benharmonics/personal-site-backend/api/requests"
	"github.com/benharmonics/personal-site-backend/chatroom"
	"github.com/benharmonics/personal-site-backend/database"
	"github.com/benharmonics/personal-site-backend/database/models"
	"github.com/benharmonics/personal-site-backend/logging"
	"github.com/benharmonics/personal-site-backend/utils/auth"
	"github.com/benharmonics/personal-site-backend/utils/web"
)

func (s *Server) routes() {
	// We have to set the chatroom database at some point or chat messages will never get saved
	chatroom.SetDatabase(s.db)

	s.HandleFunc("/heartbeat", heartbeat)

	s.HandleFunc("/users/new", cors(createUser(s.db), http.MethodPost))
	s.HandleFunc("/login", cors(login(s.db), http.MethodPost))

	s.HandleFunc("/blogs", cors(getBlogPosts(s.db), http.MethodGet))
	s.HandleFunc("/blogs/id/", cors(getBlogPostByID(s.db), http.MethodGet))
	s.HandleFunc("/blogs/new", cors(newBlogPost(s.db), http.MethodPost))

	s.HandleFunc("/ws/chat/", serveChatroom)

	s.HandleFunc("/images/", serveImages())

	s.HandleFunc("/", cors(http.NotFoundHandler()))
}

func heartbeat(_ http.ResponseWriter, _ *http.Request) {}

func serveChatroom(w http.ResponseWriter, r *http.Request) {
	roomName := strings.Split(r.RequestURI, "/chat/")[1]
	logging.Debug("Serving chatroom", roomName)
	chatroom.ServeChatroom(roomName, w, r)
}

func serveImages() http.HandlerFunc {
	serveImagesHandler := http.StripPrefix("/images/", http.FileServer(http.Dir("static/images")))
	return cors(serveImagesHandler, http.MethodGet)
}

func createUser(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &requests.NewUser{}
		if err := web.DecodeHTTPRequest(r, req); err != nil {
			logAndEmitHTTPError(w, r, http.StatusBadRequest, err.Error())
			return
		}
		user, err := models.NewUser(req.Email, req.Password)
		if err != nil {
			logging.Error("Failed to create a new user:", err)
			logAndEmitHTTPError(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		if err := db.InsertUser(user); err != nil {
			logging.Error("Failed to insert new user into database:", err)
			logAndEmitHTTPError(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		logging.Info("Created new user with ID", user.ID)
		_ = json.NewEncoder(w).Encode(user)
		logging.HTTPOk(r)
	}
}

func login(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &requests.Login{}
		if err := web.DecodeHTTPRequest(r, req); err != nil {
			logAndEmitHTTPError(w, r, http.StatusBadRequest, err.Error())
			return
		}
		user, err := db.FindUser(req.Email)
		if errors.Is(err, mongo.ErrNoDocuments) {
			logging.Error("Invalid email")
			logAndEmitHTTPError(w, r, http.StatusUnauthorized, "invalid email or password")
			return
		} else if err != nil {
			logging.Error("Failed to query database:", err)
			logAndEmitHTTPError(w, r, http.StatusInternalServerError)
			return
		}
		if success, err := auth.ComparePasswordAndHash(req.Password, user.PasswordHash); err != nil {
			logging.Error("Failed to compare password and hash:", err)
			logAndEmitHTTPError(w, r, http.StatusInternalServerError)
			return
		} else if !success {
			logAndEmitHTTPError(w, r, http.StatusUnauthorized, "invalid email or password")
			return
		}
		logging.HTTPOk(r)
	}
}

func getBlogPosts(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		opts := options.Find()
		opts.SetProjection(bson.M{"content": 0})
		opts.SetSort(bson.M{"dateCreated": -1})
		posts, err := db.FindBlogs(bson.M{}, opts)
		if err != nil {
			logging.Error("Failed to get blogs from database:", err)
			logAndEmitHTTPError(w, r, http.StatusFailedDependency)
			return
		}
		_ = json.NewEncoder(w).Encode(posts) // can't error
		logging.HTTPOk(r)
	}
}

func getBlogPostByID(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		routeParts := strings.Split(r.RequestURI, "/") // i.e. /blogs/id/.../...
		if len(routeParts) != 4 {
			logAndEmitHTTPError(w, r, http.StatusNotFound)
			return
		}
		id, err := primitive.ObjectIDFromHex(routeParts[3])
		if err != nil {
			logAndEmitHTTPError(w, r, http.StatusBadRequest)
			return
		}
		post, err := db.FindBlog(bson.M{"_id": id}, nil)
		if err != nil {
			logAndEmitHTTPError(w, r, http.StatusNotFound, "Post not found")
			return
		}
		_ = json.NewEncoder(w).Encode(post)
		logging.HTTPOk(r)
	}
}

func newBlogPost(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := requests.NewBlogPost{}
		if err := web.DecodeHTTPRequest(r, &req); err != nil {
			logAndEmitHTTPError(w, r, http.StatusBadRequest, err.Error())
			return
		}
		post := models.NewBlogPost(models.FromRequest(req))
		if err := db.InsertBlog(post); err != nil {
			logging.Error("Failed to insert blog to database:", err)
			logAndEmitHTTPError(w, r, http.StatusFailedDependency)
			return
		}
		logging.Infof("New Blog: %+v", post)
		logging.HTTPOk(r)
	}
}

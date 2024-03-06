package admin

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"sync"
	"testing/dating/api/pkg/domain"
	"testing/dating/api/pkg/web/common"
	"time"
)

var TokenName = "X-Admin-Token"
var sessionTTL = time.Hour * 24 * 1

type Session struct {
	Token      string        `json:"token"`
	AdminID    uint64        `json:"-"`
	Admin      *domain.Admin `json:"admin"`
	timer      *time.Timer
	ctx        context.Context
	cancelFunc context.CancelFunc
}

type SessionController struct {
	ctx        context.Context
	sessions   map[string]*Session
	tokenName  string
	sessionTTL time.Duration
	lock       sync.RWMutex
}

func NewSessionController(ctx context.Context) *SessionController {
	return &SessionController{
		ctx:      ctx,
		sessions: make(map[string]*Session),
	}
}

func (c *SessionController) Get(token string) *Session {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if session, ok := c.sessions[token]; ok {
		if session.timer != nil {
			session.timer.Reset(sessionTTL)
		}
		return session
	}
	return c.sessions[token]
}

func (c *SessionController) Remove(token string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if session, ok := c.sessions[token]; ok {
		if session.timer != nil {
			session.timer.Stop()
		}
		if session.cancelFunc != nil {
			session.cancelFunc()
		}
	}
	delete(c.sessions, token)
}

func (c *SessionController) New(admin *domain.Admin) *Session {
	c.lock.Lock()
	defer c.lock.Unlock()
	token := uuid.New().String()
	session := &Session{
		Token:   token,
		AdminID: admin.ID,
		Admin:   admin,
	}
	session.ctx, session.cancelFunc = context.WithCancel(c.ctx)
	session.timer = time.NewTimer(sessionTTL)
	go func() {
		defer session.timer.Stop()
		select {
		case <-session.timer.C:
			c.Remove(token)
		case <-session.ctx.Done():
			return
		}
	}()
	c.sessions[token] = session
	return session
}

func (c *SessionController) AuthMiddleware() gin.HandlerFunc {
	return func(g *gin.Context) {
		var err error
		token, err := g.Cookie(TokenName)
		if err != nil || token == "" {
			token = g.Request.Header.Get(TokenName)
			if token == "" {
				g.AbortWithStatusJSON(401, common.WebError{Err: "Unauthorized"})
				return
			}
		}
		session := c.Get(token)
		if session == nil || session.AdminID == 0 {
			g.AbortWithStatusJSON(401, common.WebError{Err: "Unauthorized"})
			return
		}
		log.Info().Interface("admin", session.Admin).Msgf("Admin %d authorized", session.AdminID)
		g.Set("admin_id", session.AdminID)
		g.Set("admin", session.Admin)
		g.Next()
	}
}

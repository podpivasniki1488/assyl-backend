package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/mail"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/podpivasniki1488/assyl-backend/internal/model"
	"github.com/podpivasniki1488/assyl-backend/internal/repository"
	"github.com/podpivasniki1488/assyl-backend/protopb"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"
)

const (
	UsernameTypeNone = iota + 1
	UsernameTypeEmail
	UsernameTypePhone
)

type authService struct {
	repo        *repository.Repository
	secretKey   string
	tracer      trace.Tracer
	redisClient *redis.Client
}

func NewAuthService(repo *repository.Repository, secretKey string, tracer trace.Tracer, redisCli *redis.Client) Auth {
	return &authService{
		repo:        repo,
		secretKey:   secretKey,
		tracer:      tracer,
		redisClient: redisCli,
	}
}

func (a *authService) Confirm(ctx context.Context, username, otpCode string) error {
	ctx, span := a.tracer.Start(ctx, "authService.Confirm")
	defer span.End()

	otp, err := a.redisClient.Get(ctx, username).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return errors.New("no otp code found, or u need register again")
		}

		return err
	}

	if otp != otpCode {
		return fmt.Errorf("otp code is not correct, sended otp %s, redis otp %s", otpCode, otp)
	}

	u, err := a.repo.UserRepo.FindByUsername(ctx, username)
	if err != nil {
		return err
	}

	u.IsApproved = true

	if _, err = a.repo.UserRepo.UpdateUser(ctx, u); err != nil {
		return err
	}

	return nil
}

func (a *authService) Login(ctx context.Context, user model.User) (token string, err error) {
	ctx, span := a.tracer.Start(ctx, "authService.Login")
	defer span.End()

	u, err := a.repo.UserRepo.FindByUsername(ctx, user.Username)
	if err != nil {
		return "", err
	}

	if !u.IsApproved {
		return "", model.ErrUserNotApproved
	}

	if ok := a.comparePasswords(u.Password, user.Password); !ok {
		return "", model.ErrPasswordMatch
	}

	token, err = a.generateJwtToken(u.Username, u.RoleID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (a *authService) Register(ctx context.Context, user model.User) error {
	ctx, span := a.tracer.Start(ctx, "authService.Register")
	defer span.End()

	if _, err := a.repo.UserRepo.FindByUsername(ctx, user.Username); err == nil {
		return model.ErrUserAlreadyExists
	}

	var usernameType int
	switch {
	case a.isEmail(user.Username):
		usernameType = UsernameTypeEmail
	case a.isPhone(user.Username):
		usernameType = UsernameTypePhone
	default:
		usernameType = UsernameTypeNone
	}

	otp := rand.Intn(9000) + 1000

	// send message to
	switch usernameType {
	case UsernameTypeEmail:
		if err := a.repo.EmailRepo.SendEmail(
			ctx,
			[]string{user.Username},
			"OTP Code",
			fmt.Sprintf("enter this otp code in order to confirm registration : %d", otp),
		); err != nil {
			return err
		}
	case UsernameTypeNone:
		// pass
	case UsernameTypePhone:
		return errors.New("not implemented")
	default:
		return errors.New("unknown username type")
	}

	if err := a.redisClient.Set(
		ctx,
		user.Username,
		otp,
		5*time.Minute,
	).Err(); err != nil {
		return err
	}

	hashedPsw, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u := model.User{
		Username:     user.Username,
		Password:     string(hashedPsw),
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		UsernameType: usernameType,
		IsApproved:   false,
		RoleID:       protopb.Role_GUEST,
		ApartmentID:  uuid.Nil,
	}

	if err = a.repo.UserRepo.CreateUser(ctx, &u); err != nil {
		return err
	}

	return nil
}

func (a *authService) isEmail(str string) bool {
	_, err := mail.ParseAddress(str)
	return err == nil
}

func (a *authService) isPhone(str string) bool {
	re := regexp.MustCompile(`(?:^|[^0-9])(1[34578][0-9]{9})(?:$|[^0-9])`)
	submatch := re.FindStringSubmatch(str)
	if len(submatch) < 2 {
		return false
	}
	//match := submatch[1]

	return true
}

func (a *authService) generateJwtToken(username string, role protopb.Role) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"role":     role.String(),
			"issuer":   "jeffry's backend",
		})

	tokenString, err := token.SignedString([]byte(a.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *authService) comparePasswords(hashedPw string, plainPw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPw), []byte(plainPw))

	return err == nil
}

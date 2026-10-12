package service

import (
	"MyMessenger/services/auth/internal/domain"
	"MyMessenger/services/auth/internal/infra/security/code"
	"context"
	"crypto/subtle"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (a *JwtAuth) StartRegistration(ctx context.Context, login, password, email string) (uuid.UUID, time.Duration, error) {
	login, email = domain.NormalizeLogin(login), domain.NormalizeEmail(email)
	log := a.log.Ctx(ctx).With(zap.String("login", login), zap.String("email", email))

	loginTaken, emailTaken, err := a.repo.CheckTaken(ctx, login, email)
	if err != nil {
		log.Error("failed to check login and email", zap.Error(err))
		return uuid.Nil, 0, domain.ErrInternal
	}
	if loginTaken {
		log.Debug("login already taken")
		return uuid.Nil, 0, domain.ErrLoginTaken
	}

	reg := &domain.Registration{
		ID:         uuid.New(),
		Login:      login,
		Email:      email,
		EmailTaken: emailTaken,
	}
	log = log.With(zap.Stringer("registration_id", reg.ID))

	err = a.withEmailCooldown(ctx, log, email, func() error {
		hashed, err := a.hash.Hash(password)
		if err != nil {
			log.Error("failed to hash password", zap.Error(err))
			return domain.ErrInternal
		}
		reg.PasswordHash = hashed
		err = a.registrations.Save(ctx, reg, a.conf.RegistrationTTL)
		if err != nil {
			log.Error("failed to save registration", zap.Error(err))
			return domain.ErrInternal
		}
		return a.sendRegistrationEmail(ctx, log, reg)
	})
	if err != nil {
		return uuid.Nil, 0, err
	}

	log.Info("registration started")
	return reg.ID, a.conf.VerificationCodeTTL, nil
}

func (a *JwtAuth) ResendCode(ctx context.Context, registrationId uuid.UUID) (time.Duration, error) {
	log := a.log.Ctx(ctx).With(zap.Stringer("registration_id", registrationId))

	reg, err := a.getRegistration(ctx, log, registrationId)
	if err != nil {
		return 0, err
	}
	err = a.withEmailCooldown(ctx, log, reg.Email, func() error {
		return a.sendRegistrationEmail(ctx, log, reg)
	})
	if err != nil {
		return 0, err
	}
	return a.conf.VerificationCodeTTL, nil
}

func (a *JwtAuth) ConfirmRegistration(ctx context.Context, registrationId uuid.UUID, vCode string) error {
	log := a.log.Ctx(ctx).With(zap.Stringer("registration_id", registrationId))

	reg, err := a.getRegistration(ctx, log, registrationId)
	if err != nil {
		return err
	}
	storedCode, err := a.registrations.GetCode(ctx, registrationId)
	if errors.Is(err, domain.ErrNotFound) {
		log.Debug("no active verification code")
		return domain.ErrInvalidCode
	}
	if err != nil {
		log.Error("failed to get verification code", zap.Error(err))
		return domain.ErrInternal
	}
	if subtle.ConstantTimeCompare([]byte(storedCode), []byte(vCode)) != 1 {
		log.Debug("verification code mismatch")
		return domain.ErrInvalidCode
	}

	userId := uuid.New()
	err = a.repo.AddNewUser(ctx, userId, reg.Login, reg.PasswordHash, reg.Email)
	if errors.Is(err, domain.ErrLoginTaken) || errors.Is(err, domain.ErrEmailTaken) {
		log.Debug("registration rejected", zap.Error(err))
		return err
	}
	if err != nil {
		log.Error("failed to add user", zap.Error(err))
		return domain.ErrInternal
	}

	err = a.registrations.Delete(ctx, registrationId)
	if err != nil {
		log.Warn("failed to delete registration", zap.Error(err))
	}
	log.Info("user registered", zap.Stringer("user_id", userId))
	return nil
}

func (a *JwtAuth) getRegistration(ctx context.Context, log *zap.Logger, registrationId uuid.UUID) (*domain.Registration, error) {
	reg, err := a.registrations.Get(ctx, registrationId)
	if errors.Is(err, domain.ErrNotFound) {
		log.Debug("registration not found")
		return nil, domain.ErrRegistrationNotFound
	}
	if err != nil {
		log.Error("failed to get registration", zap.Error(err))
		return nil, domain.ErrInternal
	}
	return reg, nil
}

func (a *JwtAuth) withEmailCooldown(ctx context.Context, log *zap.Logger, email string, send func() error) error {
	err := a.registrations.StartCooldown(ctx, email, a.conf.VerificationCodeTTL)
	if errors.Is(err, domain.ErrAlreadyExists) {
		log.Debug("verification code already sent")
		return domain.ErrCodeAlreadySent
	}
	if err != nil {
		log.Error("failed to start email cooldown", zap.Error(err))
		return domain.ErrInternal
	}

	err = send()
	if err != nil {
		stopErr := a.registrations.StopCooldown(ctx, email)
		if stopErr != nil {
			log.Warn("failed to stop email cooldown", zap.Error(stopErr))
		}
		return err
	}
	return nil
}

func (a *JwtAuth) sendRegistrationEmail(ctx context.Context, log *zap.Logger, reg *domain.Registration) error {
	if reg.EmailTaken {
		err := a.pub.PublishRegistrationAttempt(ctx, reg.Email)
		if err != nil {
			log.Error("failed to publish registration attempt", zap.Error(err))
			return domain.ErrInternal
		}
		log.Info("registration attempt on already registered email")
		return nil
	}

	vCode, err := code.GenerateNewCode()
	if err != nil {
		log.Error("failed to generate verification code", zap.Error(err))
		return domain.ErrInternal
	}
	err = a.registrations.SaveCode(ctx, reg.ID, vCode, a.conf.VerificationCodeTTL)
	if err != nil {
		log.Error("failed to save verification code", zap.Error(err))
		return domain.ErrInternal
	}
	err = a.pub.PublishEmailVerification(ctx, reg.Email, vCode)
	if err != nil {
		log.Error("failed to publish email verification", zap.Error(err))
		return domain.ErrInternal
	}
	return nil
}

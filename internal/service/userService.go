package service

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Mukam21/Go_E-Commerce_App/config"
	"github.com/Mukam21/Go_E-Commerce_App/internal/domain"
	"github.com/Mukam21/Go_E-Commerce_App/internal/dto"
	"github.com/Mukam21/Go_E-Commerce_App/internal/helper"
	"github.com/Mukam21/Go_E-Commerce_App/internal/repository"
	"github.com/Mukam21/Go_E-Commerce_App/pkg/notification"
	"gorm.io/gorm"
)

type UserService struct {
	Repo   repository.UserRepository
	CRepo  repository.CatalogRepository
	Auth   helper.Auth
	Config config.AppConfig
}

func (s UserService) SignUp(input dto.UserSignUp) (string, error) {

	hPassword, err := s.Auth.CreateHashedPassword(input.Password)

	if err != nil {
		return "", err
	}

	user, err := s.Repo.CreateUser(domain.User{
		Email:    input.Email,
		Password: hPassword,
		Phone:    input.Phone,
	})

	return s.Auth.GenerateToken(user.ID, user.Email, user.UserType)
}

func (s UserService) findUserByEmail(email string) (*domain.User, error) {

	user, err := s.Repo.FindUser(email)

	return &user, err
}

func (s UserService) Login(email string, password string) (string, error) {

	user, err := s.findUserByEmail(email)

	if err != nil {
		return "", errors.New("user does not exist with the provided email id")
	}

	err = s.Auth.VerifyPassword(password, user.Password)

	if err != nil {
		return "", err
	}

	// generate token

	return s.Auth.GenerateToken(user.ID, user.Email, user.UserType)
}

func (s UserService) isVerifiedUser(id uint) bool {

	currentUser, err := s.Repo.FindUserById(id)

	return err == nil && currentUser.Verified
}

func (s UserService) GetVerificationCode(e domain.User) error {

	//if user already verified
	if s.isVerifiedUser(e.ID) {
		return errors.New("user already verified")
	}

	// generate verification code
	code, err := s.Auth.GenerateCode()
	if err != nil {
		return nil
	}

	// update user
	user := domain.User{
		Expiry: time.Now().Add(30 * time.Minute),
		Code:   code,
	}

	_, err = s.Repo.UpdateUser(e.ID, user)

	if err != nil {
		return errors.New("unable to update verification code")
	}

	user, _ = s.Repo.FindUserById(e.ID) // uytgetdim user.ID

	// send SMS
	notificationClient := notification.NewNotificationClient(s.Config)
	// notificationClient.SendSMS(user.Phone, strconv.Itoa(code))

	msg := fmt.Sprintf("Your verification code is %v", code)
	err = notificationClient.SendSMS(user.Phone, msg)
	if err != nil {
		return errors.New("error on sending SMS: ")
	}

	// return verification code

	return nil
}

func (s UserService) VerifyCode(id uint, code int) error {

	//if user already verified
	if s.isVerifiedUser(id) {
		log.Println("verified...")
		return errors.New("user already verified")
	}

	user, err := s.Repo.FindUserById(id)

	if err != nil {
		return err
	}

	if user.Code != code {
		return errors.New("verification code does not match")
	}

	if !time.Now().Before(user.Expiry) {
		return errors.New("verificatioin code expired")
	}

	updateUser := domain.User{
		Verified: true,
	}

	_, err = s.Repo.UpdateUser(id, updateUser)

	if err != nil {
		return errors.New("unable to to verify user")
	}

	return nil
}

func (s UserService) CreateProfile(id uint, input dto.ProfileInput) error {

	// update user
	user, err := s.Repo.FindUserById(id)

	if err != nil {
		return err
	}
	if input.FirstName != "" {
		user.FirstName = input.FirstName
	}

	if input.LastName != "" {
		user.LastName = input.LastName
	}

	_, err = s.Repo.UpdateUser(id, user)

	if err != nil {
		return err
	}

	// create address

	address := domain.Address{
		AddressLine1: input.AddressInput.AddressLine1,
		AddressLine2: input.AddressInput.AddressLine2,
		City:         input.AddressInput.City,
		Country:      input.AddressInput.Country,
		PostCode:     input.AddressInput.PostCode,
		UserId:       id,
	}

	err = s.Repo.CreateProfile(address)
	if err != nil {
		return err
	}

	return nil
}

func (s UserService) GetProfile(id uint) (*domain.User, error) {

	user, err := s.Repo.FindUserById(id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s UserService) UpdateProfile(id uint, input dto.ProfileInput) error {

	user, err := s.Repo.FindUserById(id)

	if err != nil {
		return err
	}
	if input.FirstName != "" {
		user.FirstName = input.FirstName
	}

	if input.LastName != "" {
		user.LastName = input.LastName
	}

	_, err = s.Repo.UpdateUser(id, user)

	address := domain.Address{
		AddressLine1: input.AddressInput.AddressLine1,
		AddressLine2: input.AddressInput.AddressLine2,
		City:         input.AddressInput.City,
		Country:      input.AddressInput.Country,
		PostCode:     input.AddressInput.PostCode,
		UserId:       id,
	}

	err = s.Repo.UpdateProfile(address)
	if err != nil {
		return nil
	}

	return nil
}

func (s UserService) BecomeSeller(id uint, input dto.SellerInput) (string, error) {

	// find existing user
	user, _ := s.Repo.FindUserById(id)

	if user.UserType == domain.SELLER {
		return "", errors.New("you have already joined seller program")
	}

	// update user
	seller, err := s.Repo.UpdateUser(id, domain.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Phone:     input.PhoneNumber,
		UserType:  domain.SELLER,
	})

	if err != nil {
		return "", err
	}

	// generatting token
	token, err := s.Auth.GenerateToken(user.ID, user.Email, seller.UserType)

	// create bank account information

	err = s.Repo.CreateBankAccount(domain.BankAccount{
		BankAccount: input.BankAccountNumber,
		SwiftCode:   input.SwiftCode,
		PaymentType: input.PaymentType,
		UserId:      id,
	})

	return token, err
}

func (s UserService) FindCart(id uint) ([]domain.Cart, error) {

	cartItems, err := s.Repo.FindCartItems(id)
	log.Printf("error %v", err)

	return cartItems, err
}

func (s UserService) CreateCart(input dto.CreateCartRequest, u domain.User) ([]domain.Cart, error) {

	if input.ProductId <= 0 {
		return nil, errors.New("invalid product id")
	}

	if input.Qty < 1 {
		return nil, errors.New("quantity must be greater than 0")
	}

	cart, err := s.Repo.FindCartItem(u.ID, input.ProductId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if cart.ID > 0 {

		if input.Qty < 1 {
			if err := s.Repo.DeleteCartById(cart.ID); err != nil {
				return nil, errors.New("error deleting cart item")
			}
		} else {
			cart.Qty = input.Qty
			if err := s.Repo.UpdateCart(cart); err != nil {
				return nil, errors.New("error updating cart item")
			}
		}

	} else {

		product, err := s.CRepo.FindProductById(int(input.ProductId))
		if err != nil {
			return nil, errors.New("product not found")
		}

		err = s.Repo.CreateCart(domain.Cart{
			UserId:    u.ID,
			ProductId: input.ProductId,
			Name:      product.Name,
			ImageUrl:  product.ImageUrl,
			Qty:       input.Qty,
			Price:     product.Price,
			SellerId:  uint(product.UserId),
		})

		if err != nil {
			return nil, errors.New("error creating cart item")
		}
	}

	return s.Repo.FindCartItems(u.ID)
}

func (s UserService) CreateOrder(u domain.User) (int, error) {

	cartItems, err := s.Repo.FindCartItems(u.ID)
	if err != nil {
		return 0, errors.New("error on finding cart items")
	}

	if len(cartItems) == 0 {
		return 0, errors.New("cart is empity cannot the order")
	}

	paymentId := "PAY12345"
	txnId := "TXN12345"
	orderRef, _ := helper.RandomNomber(8)

	var amount float64
	var orderItems []domain.OrderItem

	for _, item := range cartItems {
		amount += item.Price * float64(item.Qty)
		orderItems = append(orderItems, domain.OrderItem{
			ProductId: item.ProductId,
			Qty:       item.Qty,
			Price:     item.Price,
			Name:      item.Name,
			ImageUrl:  item.ImageUrl,
			SellerId:  item.SellerId,
		})
	}

	order := domain.Order{
		UserId:         u.ID,
		PaymentId:      paymentId,
		TransactionId:  txnId,
		OrderRefNumber: uint(orderRef),
		Amount:         amount,
		Items:          orderItems,
	}

	err = s.Repo.CreateOrder(order)
	if err != nil {
		return 0, err
	}

	err = s.Repo.DeleteCartItems(u.ID)
	log.Printf("Deleting cart items Error %v", err)

	return orderRef, nil
}

func (s UserService) GetOrders(u domain.User) ([]domain.Order, error) {

	orders, err := s.Repo.FindOrders(u.ID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (s UserService) GetOrderById(id uint, uId uint) (domain.Order, error) {

	order, err := s.Repo.FindOrderById(id, uId)
	if err != nil {
		return order, err
	}

	return order, nil
}

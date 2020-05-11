package services

import (
	"log"

	"github.com/Alexsays/EithermineBot/dtos"
	"github.com/Alexsays/EithermineBot/models"
	"github.com/Alexsays/EithermineBot/repositories"
	"github.com/google/uuid"
)

// CreateUser ...
func CreateUser(user *models.User, repository repositories.UserRepository) dtos.Response {
	uuidResult, err := uuid.NewRandom()

	if err != nil {
		log.Fatalln(err)
	}

	user.ID = uuidResult.String()

	operationResult := repository.Save(user)

	if operationResult.Error != nil {
		return dtos.Response{Success: false, Message: operationResult.Error.Error()}
	}

	var data = operationResult.Result.(*models.User)

	return dtos.Response{Success: true, Data: data}
}

// FindAllUsers ...
func FindAllUsers(repository repositories.UserRepository) dtos.Response {
	operationResult := repository.FindAll()

	if operationResult.Error != nil {
		return dtos.Response{Success: false, Message: operationResult.Error.Error()}
	}

	var datas = operationResult.Result.(*models.Users)

	return dtos.Response{Success: true, Data: datas}
}

// FindOneUserByID ...
func FindOneUserByID(id string, repository repositories.UserRepository) dtos.Response {
	operationResult := repository.FindOneByID(id)

	if operationResult.Error != nil {
		return dtos.Response{Success: false, Message: operationResult.Error.Error()}
	}

	var data = operationResult.Result.(*models.User)

	return dtos.Response{Success: true, Data: data}
}

// FindOneUserByUsername ...
func FindOneUserByUsername(username string, repository repositories.UserRepository) dtos.Response {
	operationResult := repository.FindOneByUsername(username)

	if operationResult.Error != nil {
		return dtos.Response{Success: false, Message: operationResult.Error.Error()}
	}

	var data = operationResult.Result.(*models.User)

	return dtos.Response{Success: true, Data: data}
}

// UpdateUserByID ...
func UpdateUserByID(id string, user *models.User, repository repositories.UserRepository) dtos.Response {
	existingUserResponse := FindOneUserByID(id, repository)

	if !existingUserResponse.Success {
		return existingUserResponse
	}

	existingUser := existingUserResponse.Data.(*models.User)

	existingUser.FirstName = user.FirstName
	existingUser.LastName = user.LastName
	existingUser.Username = user.Username
	existingUser.TelegramID = user.TelegramID
	existingUser.EtherminerToken = user.EtherminerToken

	operationResult := repository.Save(existingUser)

	if operationResult.Error != nil {
		return dtos.Response{Success: false, Message: operationResult.Error.Error()}
	}

	return dtos.Response{Success: true, Data: operationResult.Result}
}

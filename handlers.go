package main

import (
	"encoding/json"
	"fmt"

	"hlc/models"

	"github.com/qiangxue/fasthttp-routing"
)

var (
	emptyJSONObject = []byte("{}")
)

type ValidationError string

func (msg ValidationError) Error() string {
	return string(msg)
}

func getUserByIDHandler(storage *models.Storage) func(ctx *routing.Context) error {
	return func(ctx *routing.Context) error {
		entityID, err := getParamID(ctx)
		if err != nil {
			ctx.SetStatusCode(404)
			return nil
		}

		user, err := storage.GetUserByID(entityID)
		if err != nil {
			return err
		}

		if user == nil {
			ctx.SetStatusCode(404)
			return nil
		}

		if user.Json != nil {
			ctx.Write(user.Json)
			return nil
		}

		return json.NewEncoder(ctx).Encode(user)
	}
}

func getLocationByIDHandler(storage *models.Storage) func(ctx *routing.Context) error {
	return func(ctx *routing.Context) error {
		entityID, err := getParamID(ctx)
		if err != nil {
			ctx.SetStatusCode(404)
			return nil
		}

		location, err := storage.GetLocationByID(entityID)
		if err != nil {
			return err
		}

		if location == nil {
			ctx.SetStatusCode(404)
			return nil
		}

		if location.Json != nil {
			ctx.Write(location.Json)
			return nil
		}

		return json.NewEncoder(ctx).Encode(location)
	}
}

func getVisitByIDHandler(storage *models.Storage) func(ctx *routing.Context) error {
	return func(ctx *routing.Context) error {
		entityID, err := getParamID(ctx)
		if err != nil {
			ctx.SetStatusCode(404)
			return nil
		}

		visit, err := storage.GetVisitByID(entityID)
		if err != nil {
			return err
		}

		if visit == nil {
			ctx.SetStatusCode(404)
			return nil
		}

		if visit.Json != nil {
			ctx.Write(visit.Json)
			return nil
		}

		return json.NewEncoder(ctx).Encode(visit)
	}
}

func getVisitsByUserIDHandler(storage *models.Storage) func(ctx *routing.Context) error {
	return func(ctx *routing.Context) error {
		userID, err := getParamID(ctx)
		if err != nil {
			ctx.SetStatusCode(400)
			return nil
		}

		user, err := storage.GetUserByID(userID)
		if err != nil {
			return err
		}

		if user == nil {
			ctx.SetStatusCode(404)
			return nil
		}

		args := ctx.QueryArgs()

		fromDate, err := getIntParam(args, "fromDate")
		if err != nil {
			ctx.SetStatusCode(400)
			return nil
		}

		toDate, err := getIntParam(args, "toDate")
		if err != nil {
			ctx.SetStatusCode(400)
			return nil
		}

		toDistance, err := getIntParam(args, "toDistance")
		if err != nil {
			ctx.SetStatusCode(400)
			return nil
		}

		country := getStringParam(args, "country")

		var response struct {
			Visits []models.VisitResponse `json:"visits"`
		}

		response.Visits, err = storage.GetVisitsByUser(user, fromDate, toDate, toDistance, country)
		if err != nil {
			return err
		}

		return json.NewEncoder(ctx).Encode(response)
	}
}

func getLocationAvgMarkHandler(storage *models.Storage) func(ctx *routing.Context) error {
	return func(ctx *routing.Context) error {
		locationID, err := getParamID(ctx)
		if err != nil {
			ctx.SetStatusCode(400)
			return nil
		}

		location, err := storage.GetLocationByID(locationID)
		if err != nil {
			ctx.SetStatusCode(400)
			return nil
		}

		if location == nil {
			ctx.SetStatusCode(404)
			return nil
		}

		args := ctx.QueryArgs()

		fromAge, err := getIntParam(args, "fromAge")
		if err != nil {
			ctx.SetStatusCode(400)
			return nil
		}

		toAge, err := getIntParam(args, "toAge")
		if err != nil {
			ctx.SetStatusCode(400)
			return nil
		}

		fromDate, err := getIntParam(args, "fromDate")
		if err != nil {
			ctx.SetStatusCode(400)
			return nil
		}

		toDate, err := getIntParam(args, "toDate")
		if err != nil {
			ctx.SetStatusCode(400)
			return nil
		}

		gender, err := getGenderParam(args)
		if err != nil {
			ctx.SetStatusCode(400)
			return nil
		}

		avg := location.GetAVG(fromDate, toDate, fromAge, toAge, gender)

		fmt.Fprintf(ctx, "{\"avg\": %.5f}", avg)

		return nil
	}
}

func createEntityHandler(storage *models.Storage, entityType models.EntityType) func(ctx *routing.Context) error {
	return func(ctx *routing.Context) error {
		defer ctx.SetConnectionClose()

		body := ctx.PostBody()
		if len(body) == 0 {
			ctx.SetStatusCode(400)
			return nil
		}

		var err error

		switch entityType {
		case models.UserEntityType:
			err = createUserEntity(storage, body)
		case models.LocationEntityType:
			err = createLocationEntity(storage, body)
		case models.VisitsEntityType:
			err = createVisitEntity(storage, body)
		}

		if err != nil {
			if _, ok := err.(ValidationError); ok {
				ctx.SetStatusCode(400)
				return nil
			}

			return err
		}

		_, err = ctx.Write(emptyJSONObject)
		return err
	}
}

func createUserEntity(storage *models.Storage, body []byte) error {
	var userUpdate models.UserUpdate
	err := json.Unmarshal(body, &userUpdate)
	if err != nil {
		return ValidationError(err.Error())
	}

	if !userUpdate.Valid() {
		return ValidationError("user is not valid")
	}

	return storage.SetUser(&models.User{
		ID:        userUpdate.ID.Value,
		Gender:    models.UserGender(userUpdate.Gender.Value),
		LastName:  userUpdate.LastName.Value,
		FirstName: userUpdate.FirstName.Value,
		Email:     userUpdate.Email.Value,
		BirthDate: userUpdate.BirthDate.Value,
	})
}

func createLocationEntity(storage *models.Storage, body []byte) error {
	var locationUpdate models.LocationUpdate
	err := json.Unmarshal(body, &locationUpdate)
	if err != nil {
		return ValidationError(err.Error())
	}

	if !locationUpdate.Valid() {
		return ValidationError("location is not valid")
	}

	return storage.SetLocation(&models.Location{
		ID:       locationUpdate.ID.Value,
		Distance: locationUpdate.Distance.Value,
		City:     locationUpdate.City.Value,
		Place:    locationUpdate.Place.Value,
		Country:  locationUpdate.Country.Value,
	})
}

func createVisitEntity(storage *models.Storage, body []byte) error {
	var visitUpdate models.VisitUpdate
	err := json.Unmarshal(body, &visitUpdate)
	if err != nil {
		return ValidationError(err.Error())
	}

	if !visitUpdate.Valid() {
		return ValidationError("visit is not valid")
	}

	return storage.SetVisit(&models.Visit{
		ID:         visitUpdate.ID.Value,
		UserID:     visitUpdate.UserID.Value,
		Mark:       visitUpdate.Mark.Value,
		LocationID: visitUpdate.LocationID.Value,
		VisitedAt:  visitUpdate.VisitedAt.Value,
	})
}

func updateEntityHandler(storage *models.Storage, entityType models.EntityType) func(ctx *routing.Context) error {
	return func(ctx *routing.Context) error {
		defer ctx.SetConnectionClose()

		entityID, err := getParamID(ctx)
		if err != nil {
			ctx.SetStatusCode(400)
			return nil
		}

		switch entityType {
		case models.UserEntityType:
			user, err := storage.GetUserByID(entityID)
			if err != nil {
				return err
			}

			if user == nil {
				ctx.SetStatusCode(404)
				return nil
			}
		case models.LocationEntityType:
			location, err := storage.GetLocationByID(entityID)
			if err != nil {
				return err
			}

			if location == nil {
				ctx.SetStatusCode(404)
				return nil
			}
		case models.VisitsEntityType:
			visit, err := storage.GetVisitByID(entityID)
			if err != nil {
				return err
			}

			if visit == nil {
				ctx.SetStatusCode(404)
				return nil
			}
		}

		err = updateEntity(storage, entityType, entityID, ctx)
		if err != nil {
			if _, ok := err.(ValidationError); ok {
				ctx.SetStatusCode(400)
				return nil
			}

			return err
		}

		_, err = ctx.Write(emptyJSONObject)
		return err
	}
}

func updateEntity(storage *models.Storage, et models.EntityType, eID int, ctx *routing.Context) error {
	var err error

	body := ctx.PostBody()
	if len(body) == 0 {
		return ValidationError("te")
	}

	switch et {
	case models.UserEntityType:
		err = updateUserEntity(storage, eID, ctx.PostBody())
	case models.LocationEntityType:
		err = updateLocationEntity(storage, eID, ctx.PostBody())
	case models.VisitsEntityType:
		err = updateVisitEntity(storage, eID, ctx.PostBody())
	}

	return err
}

func updateUserEntity(storage *models.Storage, id int, body []byte) error {
	var userUpdate models.UserUpdate
	err := json.Unmarshal(body, &userUpdate)
	if err != nil {
		return err
	}

	if !userUpdate.Valid() {
		return ValidationError("input struct is not valid")
	}

	user, err := storage.GetUserByID(id)
	if err != nil {
		return err
	}

	if userUpdate.BirthDate.Set {
		user.BirthDate = userUpdate.BirthDate.Value
	}
	if userUpdate.Gender.Set {
		user.Gender = models.UserGender(userUpdate.Gender.Value)
	}
	if userUpdate.Email.Set {
		user.Email = userUpdate.Email.Value
	}
	if userUpdate.FirstName.Set {
		user.FirstName = userUpdate.FirstName.Value
	}
	if userUpdate.LastName.Set {
		user.LastName = userUpdate.LastName.Value
	}

	err = storage.SetUser(user)
	if err != nil {
		return err
	}

	return nil
}

func updateLocationEntity(storage *models.Storage, id int, body []byte) error {
	var locationUpdate models.LocationUpdate
	err := json.Unmarshal(body, &locationUpdate)
	if err != nil {
		return err
	}

	if !locationUpdate.Valid() {
		return ValidationError("input struct is not valid")
	}

	location, err := storage.GetLocationByID(id)
	if err != nil {
		return err
	}

	if locationUpdate.Country.Set {
		location.Country = locationUpdate.Country.Value
	}

	if locationUpdate.Place.Set {
		location.Place = locationUpdate.Place.Value
	}

	if locationUpdate.City.Set {
		location.City = locationUpdate.City.Value
	}

	if locationUpdate.Distance.Set {
		location.Distance = locationUpdate.Distance.Value
	}

	err = storage.SetLocation(location)
	if err != nil {
		return err
	}

	return nil
}

func updateVisitEntity(storage *models.Storage, id int, body []byte) error {
	var visitUpdate models.VisitUpdate
	err := json.Unmarshal(body, &visitUpdate)
	if err != nil {
		return err
	}

	if !visitUpdate.Valid() {
		return ValidationError("input struct is not valid")
	}

	visit, err := storage.GetVisitByID(id)
	if err != nil {
		return err
	}

	if visitUpdate.VisitedAt.Set {
		visit.VisitedAt = visitUpdate.VisitedAt.Value
	}

	if visitUpdate.LocationID.Set {
		err = storage.RemoveVisitFromLocation(visit)
		if err != nil {
			return err
		}

		visit.LocationID = visitUpdate.LocationID.Value
	}

	if visitUpdate.Mark.Set {
		visit.Mark = visitUpdate.Mark.Value
	}

	if visitUpdate.UserID.Set {
		err = storage.RemoveVisitFromUser(visit)
		if err != nil {
			return err
		}

		visit.UserID = visitUpdate.UserID.Value
	}

	err = storage.SetVisit(visit)
	if err != nil {
		return err
	}

	return nil
}

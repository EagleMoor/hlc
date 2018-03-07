package models

import (
	"sort"
	"sync"
)

type Storage struct {
	users     map[int]*User
	locations map[int]*Location
	visits    map[int]*Visit

	usersLock     sync.RWMutex
	locationsLock sync.RWMutex
	visitsLock    sync.RWMutex

	update     chan UpdateRow
	updateLock sync.Mutex
}

type UpdateRow struct {
	Type   EntityType
	Entity interface{}
}

func NewStorage() *Storage {
	s := &Storage{
		users:     make(map[int]*User),
		locations: make(map[int]*Location),
		visits:    make(map[int]*Visit),
		update:    make(chan UpdateRow),
	}

	go s.updateFn()

	return s
}

func (s *Storage) updateFn() {
	for {
		row := <-s.update

		if row.Type == UserEntityType {
			u := row.Entity.(*User)
			s.users[u.ID] = u
		}

		if row.Type == LocationEntityType {
			l := row.Entity.(*Location)
			s.locations[l.ID] = l
		}

		if row.Type == VisitsEntityType {
			v := row.Entity.(*Visit)

			user := s.users[v.UserID]
			if user.visits == nil {
				user.visits = make(map[int]*Visit)
			}
			v.User = user
			user.visits[v.ID] = v

			location := s.locations[v.LocationID]
			if location.visits == nil {
				location.visits = make(map[int]*Visit)
			}
			location.visits[v.ID] = v

			s.visits[v.ID] = v
		}
	}
}

func (s *Storage) SetUser(users ...*User) error {
	for _, u := range users {
		s.update <- UpdateRow{
			Type:   UserEntityType,
			Entity: u,
		}
	}

	return nil
}

func (s *Storage) GetUserByID(id int) (*User, error) {
	return s.users[id], nil
}

func (s *Storage) SetLocation(locations ...*Location) error {
	for _, l := range locations {
		s.update <- UpdateRow{
			Type:   LocationEntityType,
			Entity: l,
		}
	}

	return nil
}

func (s *Storage) GetLocationByID(id int) (*Location, error) {
	return s.locations[id], nil
}

func (s *Storage) SetVisit(visits ...*Visit) error {
	for _, v := range visits {
		s.update <- UpdateRow{
			Type:   VisitsEntityType,
			Entity: v,
		}
	}

	return nil
}

func (s *Storage) RemoveVisitFromUser(v *Visit) error {
	s.usersLock.RLock()
	defer s.usersLock.RUnlock()

	user := s.users[v.UserID]
	delete(user.visits, v.ID)

	return nil
}

func (s *Storage) RemoveVisitFromLocation(v *Visit) error {
	s.locationsLock.Lock()
	defer s.locationsLock.Unlock()

	location := s.locations[v.LocationID]
	delete(location.visits, v.ID)

	return nil
}

func (s *Storage) GetVisitByID(id int) (*Visit, error) {
	s.visitsLock.RLock()
	defer s.visitsLock.RUnlock()

	return s.visits[id], nil
}

func (s *Storage) GetVisitsByUser(user *User, fromDate, toDate, toDistance *int, country *string) ([]VisitResponse, error) {
	visitsResponse := make([]VisitResponse, 0, len(user.visits))

	for _, visit := range user.visits {
		if fromDate != nil && visit.VisitedAt <= *fromDate {
			continue
		}

		if toDate != nil && visit.VisitedAt >= *toDate {
			continue
		}

		location, err := s.GetLocationByID(visit.LocationID)
		if err != nil {
			return nil, err
		}

		if toDistance != nil && location.Distance >= *toDistance {
			continue
		}

		if country != nil && location.Country != *country {
			continue
		}

		visitsResponse = append(visitsResponse, VisitResponse{
			Mark:      visit.Mark,
			VisitedAt: visit.VisitedAt,
			Place:     location.Place,
		})
	}

	sort.Sort(ByVisitedAsc(visitsResponse))

	return visitsResponse, nil
}

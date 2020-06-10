package profilesvc

import (
	"context"
	"errors"
	"sync"
)

//go:generate kokgen ./service.go Service

// Service is a simple CRUD interface for user profiles.
type Service interface {
	// @kok(op): "POST /profiles"
	PostProfile(ctx context.Context, profile Profile) (err error)

	// @kok(op): "GET /profiles/{id}"
	// @kok(param): "name:id,in:path"
	GetProfile(ctx context.Context, id string) (profile Profile, err error)

	// @kok(op): "PUT /profiles/{id}"
	// @kok(param): "name:id,in:path"
	PutProfile(ctx context.Context, id string, profile Profile) (err error)

	// @kok(op): "PATCH /profiles/{id}"
	// @kok(param): "name:id,in:path"
	PatchProfile(ctx context.Context, id string, profile Profile) (err error)

	// @kok(op): "DELETE /profiles/{id}"
	// @kok(param): "name:id,in:path"
	DeleteProfile(ctx context.Context, id string) (err error)

	// @kok(op): "GET /profiles/{id}/addresses"
	// @kok(param): "name:id,in:path"
	GetAddresses(ctx context.Context, id string) (addresses []Address, err error)

	// @kok(op): "GET /profiles/{profileID}/addresses/{addressID}"
	// @kok(param): "name:profileID,in:path"
	// @kok(param): "name:addressID,in:path"
	GetAddress(ctx context.Context, profileID string, addressID string) (address Address, err error)

	// @kok(op): "POST /profiles/{profileID}/addresses"
	// @kok(param): "name:profileID,in:path"
	PostAddress(ctx context.Context, profileID string, address Address) (err error)

	// @kok(op): "DELETE /profiles/{profileID}/addresses/{addressID}"
	// @kok(param): "name:profileID,in:path"
	// @kok(param): "name:addressID,in:path"
	DeleteAddress(ctx context.Context, profileID string, addressID string) (err error)
}

// Profile represents a single user profile.
// ID should be globally unique.
type Profile struct {
	ID        string    `json:"id"`
	Name      string    `json:"name,omitempty"`
	Addresses []Address `json:"addresses,omitempty"`
}

// Address is a field of a user profile.
// ID should be unique within the profile (at a minimum).
type Address struct {
	ID       string `json:"id"`
	Location string `json:"location,omitempty"`
}

var (
	ErrInconsistentIDs = errors.New("inconsistent IDs")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotFound        = errors.New("not found")
)

type inmemService struct {
	mtx sync.RWMutex
	m   map[string]Profile
}

func NewInmemService() Service {
	return &inmemService{
		m: map[string]Profile{},
	}
}

func (s *inmemService) PostProfile(ctx context.Context, p Profile) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if _, ok := s.m[p.ID]; ok {
		return ErrAlreadyExists // POST = create, don't overwrite
	}
	s.m[p.ID] = p
	return nil
}

func (s *inmemService) GetProfile(ctx context.Context, id string) (Profile, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	p, ok := s.m[id]
	if !ok {
		return Profile{}, ErrNotFound
	}
	return p, nil
}

func (s *inmemService) PutProfile(ctx context.Context, id string, p Profile) error {
	if id != p.ID {
		return ErrInconsistentIDs
	}
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.m[id] = p // PUT = create or update
	return nil
}

func (s *inmemService) PatchProfile(ctx context.Context, id string, p Profile) error {
	if p.ID != "" && id != p.ID {
		return ErrInconsistentIDs
	}

	s.mtx.Lock()
	defer s.mtx.Unlock()

	existing, ok := s.m[id]
	if !ok {
		return ErrNotFound // PATCH = update existing, don't create
	}

	// We assume that it's not possible to PATCH the ID, and that it's not
	// possible to PATCH any field to its zero value. That is, the zero value
	// means not specified. The way around this is to use e.g. Name *string in
	// the Profile definition. But since this is just a demonstrative example,
	// I'm leaving that out.

	if p.Name != "" {
		existing.Name = p.Name
	}
	if len(p.Addresses) > 0 {
		existing.Addresses = p.Addresses
	}
	s.m[id] = existing
	return nil
}

func (s *inmemService) DeleteProfile(ctx context.Context, id string) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if _, ok := s.m[id]; !ok {
		return ErrNotFound
	}
	delete(s.m, id)
	return nil
}

func (s *inmemService) GetAddresses(ctx context.Context, profileID string) ([]Address, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	p, ok := s.m[profileID]
	if !ok {
		return []Address{}, ErrNotFound
	}
	return p.Addresses, nil
}

func (s *inmemService) GetAddress(ctx context.Context, profileID string, addressID string) (Address, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	p, ok := s.m[profileID]
	if !ok {
		return Address{}, ErrNotFound
	}
	for _, address := range p.Addresses {
		if address.ID == addressID {
			return address, nil
		}
	}
	return Address{}, ErrNotFound
}

func (s *inmemService) PostAddress(ctx context.Context, profileID string, a Address) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	p, ok := s.m[profileID]
	if !ok {
		return ErrNotFound
	}
	for _, address := range p.Addresses {
		if address.ID == a.ID {
			return ErrAlreadyExists
		}
	}
	p.Addresses = append(p.Addresses, a)
	s.m[profileID] = p
	return nil
}

func (s *inmemService) DeleteAddress(ctx context.Context, profileID string, addressID string) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	p, ok := s.m[profileID]
	if !ok {
		return ErrNotFound
	}
	newAddresses := make([]Address, 0, len(p.Addresses))
	for _, address := range p.Addresses {
		if address.ID == addressID {
			continue // delete
		}
		newAddresses = append(newAddresses, address)
	}
	if len(newAddresses) == len(p.Addresses) {
		return ErrNotFound
	}
	p.Addresses = newAddresses
	s.m[profileID] = p
	return nil
}

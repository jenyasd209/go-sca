package repos

import (
	"errors"

	"go-sca/src/database"
	"gorm.io/gorm"
)

var (
	ErrEntryNotFound = errors.New("entry not found")
	ErrNoRowAffected = errors.New("no row affected")
)

type GenericRepo[T any] struct {
	db *database.Database
}

func NewGenericRepo[T any](db *database.Database) *GenericRepo[T] {
	return &GenericRepo[T]{db: db}
}

func (scr *GenericRepo[T]) Create(entity *T) error {
	return scr.db.Create(entity).Error
}

func (scr *GenericRepo[T]) Get(id uint) (*T, error) {
	var entity T
	r := scr.db.First(&entity, id)
	if r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			return nil, ErrEntryNotFound
		}

		return nil, r.Error
	}

	return &entity, nil
}

func (scr *GenericRepo[T]) Update(id uint, entity *T, force bool) error {
	if !force || entity == nil {
		var err error
		entity, err = scr.Get(id)
		if err != nil {
			return err
		}
	}

	return scr.db.Save(entity).Error
}

func (scr *GenericRepo[T]) Delete(id uint) error {
	entity, err := scr.Get(id)
	if err != nil {
		return err
	}

	return scr.db.Delete(entity, id).Error
}

func (scr *GenericRepo[T]) GetAll() ([]*T, error) {
	var entities []*T
	result := scr.db.Find(&entities)
	if result.Error != nil {
		return nil, result.Error
	}

	return entities, nil
}

func (scr *GenericRepo[T]) RawSql(to interface{}, sql string, values ...interface{}) error {
	res := scr.db.Raw(sql, values...)
	if res.Error != nil {
		return res.Error
	}

	return res.Scan(to).Error
}

func (scr *GenericRepo[T]) Exec(sql string, values ...interface{}) error {
	res := scr.db.Exec(sql, values...)
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return ErrNoRowAffected
	}

	return nil
}

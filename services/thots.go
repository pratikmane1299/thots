package services

import (
	"fmt"
	"github.com/pratikmane1299/thots/db"
	"time"
)

type thotsService struct {
	thotsDb db.ThotsDb
}

type ThotPayload struct {
	Thot string `json:"thot"`
}

func (ts thotsService) GetAllThots() ([]db.Thot, error) {
	var thots []db.Thot

	rows, err := ts.thotsDb.Db.Query("select * from thots order by created desc")
	if err != nil {
		return thots, fmt.Errorf("unable to fetch thots %v", err)
	}

	for rows.Next() {
		var thot db.Thot
		err := rows.Scan(
			&thot.Id,
			&thot.Thot,
			&thot.Created,
		)
		if err != nil {
			return thots, err
		}

		thots = append(thots, thot)
	}

	return thots, nil
}

func (ts thotsService) AddThot(thot string) error {
	result, err := ts.thotsDb.Db.Exec("insert into thots (thot, created) values(?, ?)", thot, time.Now())
	if err != nil {
		return fmt.Errorf("error adding a new thot %v", err)
	}

	_, err = result.LastInsertId()
	return err
}

func (ts thotsService) UpdateThot(id string, thot string) error {
	result, err := ts.thotsDb.Db.Exec("update thots set thot=? where id=?", thot, id)
	if err != nil {
		return fmt.Errorf("error updating thot for id %s %v", id, err)
	}

	_, err = result.LastInsertId()
	return err
}

func (ts thotsService) DeleteThot(id string) error {
	result, err := ts.thotsDb.Db.Exec("delete from thots where id = ?", id)
	if err != nil {
		return fmt.Errorf("error deleting thot with %s - %v", id, err)
	}

	_, err = result.LastInsertId()
	return err
}

func NewThotsService(thotsDb db.ThotsDb) thotsService {
	return thotsService{thotsDb: thotsDb}
}

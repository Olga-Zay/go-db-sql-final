package main

import (
	"database/sql"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int64, error) {
	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :createAt)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("createAt", p.CreatedAt))
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	lastNumber, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastNumber, nil
}

func (s ParcelStore) Get(number int64) (Parcel, error) {
	row := s.db.QueryRow("SELECT * FROM parcel WHERE number = :number", sql.Named("number", number))

	p := Parcel{}

	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return p, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	var res []Parcel

	rows, err := s.db.Query("SELECT * FROM parcel WHERE client = :client", sql.Named("client", client))
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		p := Parcel{}

		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return res, err
		}

		res = append(res, p)
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int64, status string) error {
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number))
	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) SetAddress(number int64, address string) error {
	client, err := s.Get(number)
	if err != nil {
		return err
	}

	if client.Status == ParcelStatusRegistered {
		_, err := s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number",
			sql.Named("address", address),
			sql.Named("number", number))
		if err != nil {
			return err
		}
	}

	return nil
}

func (s ParcelStore) Delete(number int64) error {
	parc, err := s.Get(number)
	if err != nil {
		return err
	}

	if parc.Status == ParcelStatusRegistered {
		_, err := s.db.Exec("DELETE FROM parcel WHERE number = :number", sql.Named("number", number))
		if err != nil {
			return err
		}
	}
	return nil
}

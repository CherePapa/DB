package database

import (
	"encoding/json"
	"fmt"

	"go.etcd.io/bbolt"
)

type Medicine struct {
	ID    int
	Name  string
	Price float64
}

func InitDB() (*bbolt.DB, error) {
	// Открытие базы данных
	db, err := bbolt.Open("pharmacy.db", 0600, nil)
	if err != nil {
		return nil, err
	}

	// Создание корзины для хранения данных
	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Medicines"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func InsertMedicine(db *bbolt.DB, id int, name string, price float64) error {
	medicine := Medicine{ID: id, Name: name, Price: price}
	data, err := json.Marshal(medicine)
	if err != nil {
		return err
	}

	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Medicines"))
		return b.Put(itob(id), data)
	})
}

func GetMedicines(db *bbolt.DB) ([]Medicine, error) {
	var medicines []Medicine

	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Medicines"))
		return b.ForEach(func(k, v []byte) error {
			var m Medicine
			err := json.Unmarshal(v, &m)
			if err != nil {
				return err
			}
			medicines = append(medicines, m)
			return nil
		})
	})

	return medicines, err
}

func DeleteMedicine(db *bbolt.DB, id int) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Medicines"))
		return b.Delete(itob(id))
	})
}

func UpdateMedicine(db *bbolt.DB, id int, name string, price float64) error {
	medicine := Medicine{ID: id, Name: name, Price: price}
	data, err := json.Marshal(medicine)
	if err != nil {
		return err
	}

	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Medicines"))
		return b.Put(itob(id), data)
	})
}

// itob возвращает 8-байтовый big-endian представление числа
func itob(v int) []byte {
	b := make([]byte, 8)
	for i := 0; i < 8; i++ {
		b[7-i] = byte(v >> (i * 8))
	}
	return b
}

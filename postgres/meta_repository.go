package postgres

import (
	"database/sql"
	"github.com/google/uuid"
	"pdf_service_api/domain"
)

type metaRepository struct {
	DatabaseHandler DatabaseHandler
}

func NewMetaRepository(db DatabaseHandler) domain.MetaRepository {
	return metaRepository{DatabaseHandler: db}
}

func (m metaRepository) AddMeta(data domain.MetaData) error {
	if err := m.DatabaseHandler.WithConnection(addMetaDataFunction(data)); err != nil {
		return err
	}

	return nil
}

func (m metaRepository) DeleteMeta(data domain.MetaData) error {
	if err := m.DatabaseHandler.WithConnection(removeMetaDataFunction(data)); err != nil {
		return err
	}

	return nil
}

func (m metaRepository) UpdateMeta(data domain.MetaData) error {
	if err := m.DatabaseHandler.WithConnection(updateMetaDataFunction(data)); err != nil {
		return err
	}

	return nil
}

func (m metaRepository) GetMeta(uid uuid.UUID) (domain.MetaData, error) {
	returnedData := &domain.MetaData{}
	callbackFunction := func(data domain.MetaData) error {
		*returnedData = data
		return nil
	}

	if err := m.DatabaseHandler.WithConnection(getMetaDataFunction(uid, callbackFunction)); err != nil {
		return domain.MetaData{}, err
	}

	return *returnedData, nil
}

func addMetaDataFunction(data domain.MetaData) func(db *sql.DB) error {
	return func(db *sql.DB) error {
		SqlStatement := `INSERT INTO documentmeta_table ("Document_UUID", "Number_Of_Pages", "Height", "Width", "Images") values ($1, $2, $3, $4, $5)`
		if _, err := db.Exec(SqlStatement, data.UUID, data.NumberOfPages, data.Height, data.Width, data.Images); err != nil {
			return err
		}

		return nil
	}
}

func removeMetaDataFunction(data domain.MetaData) func(db *sql.DB) error {
	return func(db *sql.DB) error {
		SqlStatement := `DELETE FROM documentmeta_table WHERE "Document_UUID" = $1`
		if _, err := db.Exec(SqlStatement, data.UUID); err != nil {
			return err
		}

		return nil
	}
}

func updateMetaDataFunction(data domain.MetaData) func(db *sql.DB) error {
	return func(db *sql.DB) error {
		SqlStatement := `UPDATE documentmeta_table SET "Number_Of_Pages" = COALESCE($1, "Number_Of_Pages"), "Height" = COALESCE($2, "Height"), "Width" = COALESCE($3, "Width"), "Images" = COALESCE($4, "Images") where "Document_UUID" = $5`
		if _, err := db.Exec(SqlStatement, data.NumberOfPages, data.Height, data.Width, data.Images, data.UUID); err != nil {
			return err
		}

		return nil
	}
}

func getMetaDataFunction(uid uuid.UUID, callback func(data domain.MetaData) error) func(db *sql.DB) error {
	return func(db *sql.DB) error {
		meta := &domain.MetaData{}
		SqlStatement := `SELECT "Document_UUID", "Number_Of_Pages", "Height", "Width", "Images" FROM documentmeta_table where "Document_UUID" = $1`

		row := db.QueryRow(SqlStatement, uid)
		err := row.Scan(&meta.UUID, &meta.NumberOfPages, &meta.Height, &meta.Width, &meta.Images)
		if err != nil {
			return err
		}

		return callback(*meta)
	}
}

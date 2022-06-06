package metadata

import (
	"database/sql"
	"time"

	"github.com/bluele/gcache"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lbryio/lbry.go/v2/extras/errors"
	"github.com/sirupsen/logrus"
)

/*
CREATE TABLE `metadata` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `original_url` text NOT NULL,
  `godycdn_hash` varchar(64) NOT NULL,
  `checksum` varchar(64) NOT NULL,
  `original_size` int(11) NOT NULL,
  `otpimized_size` int(11) NOT NULL,
  `original_mime` varchar(100) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `metadata_godycdn_hash_index` (`godycdn_hash`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
*/

type Manager struct {
	dbConn *sql.DB
	cache  gcache.Cache
}

var instance *Manager

func Init(dsn string) (*Manager, error) {
	if instance != nil {
		return instance, nil
	}
	db, err := connect(dsn)
	if err != nil {
		return nil, err
	}
	instance = &Manager{
		dbConn: db,
		cache:  gcache.New(10000).Expiration(24 * time.Hour).LRU().Build(),
	}
	return instance, nil
}

func connect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	return db, errors.Err(err)
}

type ImageMetadata struct {
	OriginalURL       string `json:"original_url"`
	GodycdnHash       string `json:"godycdn_ash"`
	Checksum          string `json:"checksum"`
	OriginalMimeType  string `json:"original_mime_type"`
	OriginalSize      int    `json:"original_size"`
	OptimizedSize     int    `json:"optimized_size"`
	OptimizedMimeType string `json:"optimized_mime_type"`
}

func (m *Manager) Persist(md *ImageMetadata) error {
	query := "INSERT IGNORE INTO mirage.metadata (original_url, godycdn_hash, checksum, original_size, otpimized_size, original_mime) VALUES (?,?,?,?,?,?)"
	r, err := m.dbConn.Query(query, md.OriginalURL, md.GodycdnHash, md.Checksum, md.OriginalSize, md.OptimizedSize, md.OriginalMimeType)
	if err != nil {
		return errors.Err(err)
	}
	_ = r.Close()
	err = m.cache.Set(md.GodycdnHash, *md)
	if err != nil {
		return errors.Err(err)
	}
	return nil
}

func (m *Manager) Retrieve(godyCdnHash string) (*ImageMetadata, error) {
	cached, err := m.cache.Get(godyCdnHash)
	if err == nil && cached != nil {
		md := cached.(ImageMetadata)
		return &md, nil
	}
	query := "SELECT original_url, godycdn_hash, checksum, original_size, otpimized_size, original_mime FROM metadata WHERE godycdn_hash = ?"
	row := m.dbConn.QueryRow(query, godyCdnHash)
	var md ImageMetadata
	err = row.Scan(&md.OriginalURL, &md.GodycdnHash, &md.Checksum, &md.OriginalSize, &md.OptimizedSize, &md.OriginalMimeType)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Err(err)
	}
	err = m.cache.Set(md.GodycdnHash, md)
	if err != nil {
		logrus.Errorf("failed to cache metadata %s", errors.FullTrace(err))
	}
	return &md, nil
}

func (m *Manager) RetrieveAllForUrl(originalUrl string) ([]*ImageMetadata, error) {
	query := "SELECT original_url, godycdn_hash, checksum, original_size, otpimized_size, original_mime FROM metadata WHERE original_url = ?"
	rows, err := m.dbConn.Query(query, originalUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Err(err)
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	mdSlice := make([]*ImageMetadata, 0, 1)
	for rows.Next() {
		var md ImageMetadata
		err = rows.Scan(&md.OriginalURL, &md.GodycdnHash, &md.Checksum, &md.OriginalSize, &md.OptimizedSize, &md.OriginalMimeType)
		if err != nil {
			return nil, errors.Err(err)
		}
		mdSlice = append(mdSlice, &md)
	}

	return mdSlice, nil
}
func (m *Manager) Delete(md *ImageMetadata) error {
	query := "DELETE FROM metadata WHERE godycdn_hash = ?"
	_, err := m.dbConn.Exec(query, md.GodycdnHash)
	if err != nil {
		return errors.Err(err)
	}
	_ = m.cache.Remove(md.GodycdnHash)
	return nil
}

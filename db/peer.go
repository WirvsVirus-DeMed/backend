package db

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"time"
)

// Peer for caching the Doctors office
type Peer struct {
	IP       net.IP
	Port     uint32
	LastSeen time.Time
}

// Add adds a Peer Object to the Database
func (peer *Peer) Add(db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("insert into peer(ip, port, lastSeen) values(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(peer.IP, peer.Port, peer.LastSeen)
	if err != nil {
		log.Fatal(err)
	}

	tx.Commit()
}

// Delete Peer from Database
func (peer *Peer) Delete(db *sql.DB) {
	_, err := db.Exec("delete from peer where ip=? and port=?", peer.IP, peer.Port)
	if err != nil {
		log.Fatal(err)
	}
}

// Update the Peer object in the Database
func (peer *Peer) Update(db *sql.DB) {
	peer.Delete(db)
	peer.Add(db)
}

// get wrapper to find a set of results
func GetPeers(db *sql.DB, query string, searchStr string) ([]*Peer, error) {
	// SELECT * FROM table WHERE instr(title, searchStr) > 0 OR instr(description, searchStr) > 0 OR searchStr == CAST(pzn as text)
	// SELECT * FROM peer WHERE instr(title, '3') > 0 OR instr(description, '3') > 0 OR '3' == CAST(pzn as text);

	rows, err := db.Query(query, searchStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var peers []*Peer

	for rows.Next() {
		var ip net.IP
		var port uint32
		var lastSeen time.Time

		err = rows.Scan(&ip, &port, &lastSeen)
		if err != nil {
			log.Fatal(err)
		}

		peer := &Peer{ip, port, lastSeen}
		peers = append(peers, peer)
		fmt.Println(ip, port, lastSeen)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return peers, nil
}

// GetAllPeers all the rows of the Database
func GetAllPeers(db *sql.DB) ([]*Peer, error) {
	return GetPeers(db, "select * from peer", "")
}

// DeletePeerTable deletes Peerbacke from Database
func DeletePeerTable(db *sql.DB) {
	_, err := db.Exec("delete from peer")
	if err != nil {
		log.Fatal(err)
	}
}

// CreatePeerTable creates the Peer Table
func CreatePeerTable(db *sql.DB) {
	sqlStmt := `
	create table if not exists peer (ip blob not null, port integer not null, lastSeen timestamp);
	create unique index peer_index on peer (ip, port);`

	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

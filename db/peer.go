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
	ID       int
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

	stmt, err := tx.Prepare("insert into peer(id, ip, port, lastSeen) values(?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(peer.ID, peer.IP, peer.Port, peer.LastSeen)
	if err != nil {
		log.Fatal(err)
	}

	tx.Commit()
}

// Delete Peer from Database
func (peer *Peer) Delete(db *sql.DB) {
	_, err := db.Exec("delete from peer where id=?", peer.ID)
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
func getPeers(db *sql.DB, query string, searchStr int) ([]*Peer, error) {
	// SELECT * FROM table WHERE instr(title, searchStr) > 0 OR instr(description, searchStr) > 0 OR searchStr == CAST(pzn as text)
	// SELECT * FROM peer WHERE instr(title, '3') > 0 OR instr(description, '3') > 0 OR '3' == CAST(pzn as text);

	rows, err := db.Query(query, searchStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var peers []*Peer

	for rows.Next() {
		var id int
		var ip net.IP
		var port uint32
		var lastSeen time.Time

		err = rows.Scan(&id, &ip, &port, &lastSeen)
		if err != nil {
			log.Fatal(err)
		}

		peer := &Peer{id, ip, port, lastSeen}
		peers = append(peers, peer)
		fmt.Println(id, ip, port, lastSeen)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return peers, nil
}

// GetPeer a specific peer from the Database based on the search String or the pzn
func GetPeer(db *sql.DB, id int) ([]*Peer, error) {
	return getPeers(db, "select * from peer where id=?", id)
}

// GetAllPeers all the rows of the Database
func GetAllPeers(db *sql.DB) ([]*Peer, error) {
	return getPeers(db, "select * from peer", 0)
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
	create table peer (id integer not null primary key, ip blob, port integer, lastSeen timestamp);`

	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

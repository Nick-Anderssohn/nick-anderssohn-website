package auth

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

const (
	createUpgradeTableSql = `
CREATE TABLE IF NOT EXISTS upgrade(
  upgradeNumber INTEGER PRIMARY KEY
)`

	insertUpgradeSql = `INSERT INTO upgrade VALUES ($1);`
	selectUpgradeSql   = `SELECT * FROM upgrade WHERE upgradeNumber = $1;`
)

// IMPORTANT: When adding a new upgrade sql statement, add it to the end! Don't change order.
// index = upgrade number
var upgrades = []string{createUserTableSql, createSessionTableSql}

func createUpgradeTable() {
	_, err := db.Exec(createUpgradeTableSql)
	if err != nil {
		panic(err)
	}
}

func runUpgrades() {
	// For each upgrade sql statement
	for upgradeNumber, sqlStatement := range upgrades {
		// Check upgrade has already been ran
		hasRan, err := upgradeHasBeenRan(upgradeNumber)
		if err != nil {
			panic(err)
		}
		if !hasRan {
			log.Println("Upgrade ran:", sqlStatement)
			// Run upgrade
			_, err = db.Exec(sqlStatement)
			if err != nil {
				panic(err)
			}

			// Record in upgrades table
			_, err = db.Exec(insertUpgradeSql, upgradeNumber)
			if err != nil {
				log.Println("could not record upgrade number ", upgradeNumber)
				panic(err)
			}
		}
	}
}

// selectUpgrade returns a file info based off of code
func selectUpgrade(upgradeNumber int) (upNumb int, err error) {
	row := db.QueryRow(selectUpgradeSql, upgradeNumber)
	err = row.Scan(&upNumb)
	return
}

// upgradeHasBeenRan checks if a file exists. Returns an error if there is a problem querying the database.
func upgradeHasBeenRan(upgradeNumber int) (bool, error) {
	_, err := selectUpgrade(upgradeNumber)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
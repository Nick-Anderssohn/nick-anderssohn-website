package db

import (
	"database/sql"
	"log"
)

const (
	insertIntoUpgrades = `INSERT INTO upgrades VALUES ($1);`
	selectUpgradeSql   = `SELECT * FROM upgrades WHERE upgradeNumber = $1;`

	alterFileSizeToBigIntIfNecessary = `
ALTER TABLE files
ALTER COLUMN fileSize TYPE BIGINT;
`
)

// IMPORTANT: When adding a new upgrade sql statement, add it to the end! Don't change order.
// index = upgrade number
var upgrades = []string{alterFileSizeToBigIntIfNecessary}

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
			_, err = conn.Exec(sqlStatement)
			if err != nil {
				panic(err)
			}

			// Record in upgrades table
			_, err = conn.Exec(insertIntoUpgrades, upgradeNumber)
			if err != nil {
				log.Println("could not record upgrade number ", upgradeNumber)
				panic(err)
			}
		}
	}
}

// selectUpgrade returns a file info based off of code
func selectUpgrade(upgradeNumber int) (upNumb int, err error) {
	row := conn.QueryRow(selectUpgradeSql, upgradeNumber)
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

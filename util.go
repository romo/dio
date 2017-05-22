package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
)

// Generate the SHA256 for a commit.
func createCommitID(com commit) [32]byte {
	var b bytes.Buffer
	b.WriteString("tree " + hex.EncodeToString(com.Tree[:]) + "\n")
	if com.Parent != NILSHA256 {
		b.WriteString("parent " + hex.EncodeToString(com.Parent[:]) + "\n")
	}
	b.WriteString("author " + com.AuthorName + " <" + com.AuthorEmail + "> " +
		com.Timestamp.Format(time.UnixDate) + "\n")
	if com.CommitterEmail != "" {
		b.WriteString("committer " + com.CommitterName + " <" + com.CommitterEmail + "> " +
			com.Timestamp.Format(time.UnixDate) + "\n")
	}
	b.WriteString("\n" + com.Message)
	b.WriteByte(0)
	return sha256.Sum256(b.Bytes())
}

// Generate the SHA256 for a tree.
func createDBTreeID(entries []dbTreeEntry) [32]byte {
	var b bytes.Buffer
	for _, j := range entries {
		b.WriteString(string(j.AType))
		b.WriteByte(0)
		b.WriteString(hex.EncodeToString(j.ShaSum[:]))
		b.WriteByte(0)
		b.WriteString(j.Name + "\n")
	}
	return sha256.Sum256(b.Bytes())
}

// Store a set of branches.
func storeBranches(branches []branch) error {
	// Create the storage directory if needed
	_, err := os.Stat(STORAGEDIR)
	if err != nil {
		// As this is just experimental code, we'll assume a failure above means the directory needs creating
		err := os.MkdirAll(STORAGEDIR, os.ModeDir|0755)
		if err != nil {
			log.Printf("Something went wrong when creating the database storage dir: %v\n",
				err.Error())
			return err
		}
	}
	j, err := json.MarshalIndent(branches, "", " ")
	if err != nil {
		log.Printf("Something went wrong when serialising the branch data: %v\n", err.Error())
		return err
	}
	err = ioutil.WriteFile(STORAGEDIR+string(os.PathSeparator)+"branches", j, os.ModePerm)
	if err != nil {
		log.Printf("Something went wrong when writing the branches file: %v\n", err.Error())
		return err
	}
	return nil
}

// Store a database file.
func storeDatabase(db []byte) ([32]byte, error) {
	// Create the storage directory if needed
	_, err := os.Stat(STORAGEDIR + string(os.PathSeparator) + "dbs")
	if err != nil {
		// As this is just experimental code, we'll assume a failure above means the directory needs creating
		err := os.MkdirAll(STORAGEDIR, os.ModeDir|0755)
		if err != nil {
			log.Printf("Something went wrong when creating the database storage dir: %v\n",
				err.Error())
			return NILSHA256, err
		}
	}

	// Create the database file is it doesn't already exist
	s := sha256.Sum256(db)
	p := STORAGEDIR + string(os.PathSeparator) + "dbs" + string(os.PathSeparator) + hex.EncodeToString(s[:])
	f, err := os.Stat(p)
	if err != nil {
		// As this is just experimental code, we'll assume a failure above means the file needs creating
		err = ioutil.WriteFile(p, db, os.ModePerm)
		if err != nil {
			log.Printf("Something went wrong when writing the database file: %v\n", err.Error())
			return NILSHA256, err
		}
		return s, nil
	}

	// The file already exists, so check if the file size matches the buffer size we're intending on writing
	// (Obviously this is just a super lightweight check, not a real world approach)
	if len(db) != int(f.Size()) {
		err = ioutil.WriteFile(p, db, os.ModePerm)
		if err != nil {
			log.Printf("Something went wrong when writing the database file: %v\n", err.Error())
			return NILSHA256, err
		}
	}
	return s, nil
}

// Store an index.
func storeIndex(index []commit) error {
	// Create the storage directory if needed
	_, err := os.Stat(STORAGEDIR)
	if err != nil {
		// As this is just experimental code, we'll assume a failure above means the directory needs creating
		err := os.MkdirAll(STORAGEDIR, os.ModeDir|0755)
		if err != nil {
			log.Printf("Something went wrong when creating the database storage dir: %v\n",
				err.Error())
			return err
		}
	}
	j, err := json.MarshalIndent(index, "", " ")
	if err != nil {
		log.Printf("Something went wrong when serialising the index data: %v\n", err.Error())
		return err
	}
	err = ioutil.WriteFile(STORAGEDIR+string(os.PathSeparator)+"index", j, os.ModePerm)
	if err != nil {
		log.Printf("Something went wrong when writing the index file: %v\n", err.Error())
		return err
	}
	return nil
}
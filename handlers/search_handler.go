package handlers

import (
    "database/sql"
    "encoding/json"
    "net/http"
)

// SearchEnrollments searches enrollments by student name
func SearchEnrollments(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    studentName := r.URL.Query().Get("name")
    
    // VULNERABLE: SQL Injection
    query := "SELECT * FROM enrollments WHERE student_name = '" + studentName + "'"
    
    rows, err := db.Query(query)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    defer rows.Close()
    
    // Process results...
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
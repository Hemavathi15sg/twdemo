package handlers

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/ssh"
)

// AuthHandler demonstrates usage of vulnerable golang.org/x/crypto package
// This is used in the security demo to show CVE remediation
type AuthHandler struct{}

// GenerateSSHKey generates an SSH key pair for authentication
// Uses golang.org/x/crypto v0.14.0 which has CVE-2023-48795 (Terrapin Attack)
// This vulnerability affects SSH protocol sequence number validation
func (h *AuthHandler) GenerateSSHKey(w http.ResponseWriter, r *http.Request) {
	// Generate ED25519 key pair
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		http.Error(w, "Failed to generate key", http.StatusInternalServerError)
		return
	}

	// Convert to SSH format using vulnerable crypto package
	_, err = ssh.NewSignerFromKey(privateKey)
	if err != nil {
		http.Error(w, "Failed to create SSH signer", http.StatusInternalServerError)
		return
	}

	sshPublicKey, err := ssh.NewPublicKey(publicKey)
	if err != nil {
		http.Error(w, "Failed to create SSH public key", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"public_key":  string(ssh.MarshalAuthorizedKey(sshPublicKey)),
		"key_type":    sshPublicKey.Type(),
		"fingerprint": ssh.FingerprintSHA256(sshPublicKey),
		"status":      "Generated using golang.org/x/crypto v0.14.0 (VULNERABLE: CVE-2023-48795)",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

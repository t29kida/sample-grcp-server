package service

import "testing"

func Test_hash_CreateHash(t *testing.T) {
	h := NewHash()

	hash, err := h.CreateHash("password")

	if err != nil {
		t.Errorf("err should not be nil: %v", err)
	}

	if hash == "" {
		t.Error("hash should not empty string")
	}
}

func Test_hash_CompareHash(t *testing.T) {
	h := NewHash()

	t.Run("match", func(t *testing.T) {
		hash, err := h.CreateHash("password")
		if err != nil {
			t.Fatalf("failed to create hash: %v", err)
		}

		match, err := h.CompareHash("password", hash)

		if err != nil {
			t.Errorf("err should be nil: %v", err)
		}

		if !match {
			t.Error("should be matched")
		}
	})

	t.Run("not match", func(t *testing.T) {
		hash, err := h.CreateHash("password")
		if err != nil {
			t.Fatalf("failed to create hash: %v", err)
		}

		match, err := h.CompareHash("not_matched", hash)

		if err != nil {
			t.Errorf("err should be nil: %v", err)
		}

		if match {
			t.Error("should not be matched")
		}
	})
}

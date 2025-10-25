package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
)

func ExampleStorer_Add() {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	mockStorer := NewMockStorer(ctrl)

	// Настраиваем ожидаемое поведение
	mockStorer.EXPECT().
		Add(gomock.Any(), URL{
			Hash: "abc123",
			Link: "https://example.com",
		}, 1).
		Return("abc123", nil)

	ctx := context.Background()
	hash, err := mockStorer.Add(ctx, URL{
		Hash: "abc123",
		Link: "https://example.com",
	}, 1)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("URL added successfully with hash: %s\n", hash)
	}
	// Output: URL added successfully with hash: abc123
}

func ExampleStorer_Add_duplicate() {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	mockStorer := NewMockStorer(ctrl)

	mockStorer.EXPECT().
		Add(gomock.Any(), gomock.Any(), gomock.Any()).
		Return("existing_hash", ErrExistsURL)

	ctx := context.Background()
	hash, err := mockStorer.Add(ctx, URL{
		Hash: "abc123",
		Link: "https://example.com",
	}, 1)

	if errors.Is(err, ErrExistsURL) {
		fmt.Printf("Duplicate URL, existing hash: %s\n", hash)
	}
	// Output: Duplicate URL, existing hash: existing_hash
}

func ExampleStorer_GetByHash() {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	mockStorer := NewMockStorer(ctrl)

	mockStorer.EXPECT().
		GetByHash(gomock.Any(), "test456").
		Return(URL{
			Hash: "test456",
			Link: "https://test.com",
		}, nil)

	ctx := context.Background()
	url, err := mockStorer.GetByHash(ctx, "test456")

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Found URL: %s\n", url.Link)
	}
	// Output: Found URL: https://test.com
}

func ExampleStorer_Ping() {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	mockStorer := NewMockStorer(ctrl)

	mockStorer.EXPECT().
		Ping().
		Return(nil)

	err := mockStorer.Ping()
	if err != nil {
		fmt.Println("Storage unavailable")
	} else {
		fmt.Println("Storage is available")
	}
	// Output: Storage is available
}

func ExampleStorer_BatchAdd() {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	mockStorer := NewMockStorer(ctrl)

	urls := []URL{
		{Hash: "abc", Link: "https://example.com"},
		{Hash: "def", Link: "https://example.org"},
	}

	mockStorer.EXPECT().
		BatchAdd(gomock.Any(), urls, 1).
		Return(nil)

	ctx := context.Background()
	err := mockStorer.BatchAdd(ctx, urls, 1)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Batch URLs added successfully")
	}
	// Output: Batch URLs added successfully
}

func ExampleStorer_GetByUserID() {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	mockStorer := NewMockStorer(ctrl)

	mockStorer.EXPECT().
		GetByUserID(gomock.Any(), 1).
		Return([]URL{}, nil)

	ctx := context.Background()
	urls, err := mockStorer.GetByUserID(ctx, 1)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("User has %d URLs\n", len(urls))
	}
	// Output: User has 0 URLs
}

func ExampleStorer_CreateUser() {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	mockStorer := NewMockStorer(ctrl)

	mockStorer.EXPECT().
		CreateUser(gomock.Any()).
		Return(-1, nil)

	ctx := context.Background()
	userID, err := mockStorer.CreateUser(ctx)

	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	} else {
		fmt.Printf("Created user with ID: %d\n", userID)
	}
	// Output: Created user with ID: -1
}

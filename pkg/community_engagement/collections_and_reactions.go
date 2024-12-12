package community_engagement

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"synnergy_network/pkg/ledger"
	"time"
)

// Collection struct
type Collection struct {
	ID          string    // Unique identifier for the collection
	Name        string    // Name of the collection
	Description string    // Description of the collection
	OwnerID     string    // User ID of the collection's owner
	CreatedAt   time.Time // Timestamp for when the collection was created
	PostIDs     []string  // List of post IDs associated with the collection
}

// Post struct
type Post struct {
	ID        string    // Unique identifier for the post
	Content   string    // Content of the post
	AuthorID  string    // User ID of the post's author
	CreatedAt time.Time // Timestamp for when the post was created
}

// Reaction struct to store information about a reaction to a post or reply
type Reaction struct {
	ID        string    // Unique identifier for the reaction
	PostID    string    // ID of the post associated with the reaction (optional)
	ReplyID   string    // ID of the reply associated with the reaction (optional)
	UserID    string    // ID of the user who reacted
	Type      string    // Type of reaction (like, dislike, etc.)
	CreatedAt time.Time // Timestamp for when the reaction was created
}

// CreateCollection function to create a collection and record it in the ledger
func CreateCollection(userID string, collectionName string, description string) (string, error) {
	// Generate a unique ID for the collection
	collectionID := generateUniqueID(userID + collectionName)

	// Create the collection struct using the ledger.Collection type
	collection := ledger.Collection{
		ID:          collectionID,
		Name:        collectionName,
		Description: description,
		OwnerID:     userID,
		CreatedAt:   time.Now(),
	}

	// Initialize a new ledger instance
	l := &ledger.Ledger{}

	// Record the collection in the ledger
	if err := l.CommunityEngagementLedger.RecordCollection(collection); err != nil {
		return "", fmt.Errorf("failed to record collection: %v", err)
	}

	// Return the collection ID
	return collectionID, nil
}

func AddPostToCollection(userID, collectionID, postID string) error {
	// Initialize a new ledger instance
	l := &ledger.Ledger{}

	// Retrieve the collection
	collection, err := l.CommunityEngagementLedger.FetchCollection(collectionID)
	if err != nil {
		return fmt.Errorf("failed to retrieve collection: %v", err)
	}

	// Verify that the user is the owner of the collection
	if collection.OwnerID != userID {
		return errors.New("unauthorized: only the collection owner can add posts")
	}

	// Add the post to the collection
	err = l.CommunityEngagementLedger.AddPostToCollection(collectionID, postID)
	if err != nil {
		return fmt.Errorf("failed to add post to collection: %v", err)
	}

	return nil
}


// REMOVE_POST_FROM_COLLECTION removes a post from a specified collection.
func RemovePostFromCollection(userID, collectionID, postID string) error {
	// Create a new ledger instance
	l := &ledger.Ledger{}

	// Retrieve the collection
	collection, err := l.CommunityEngagementLedger.FetchCollection(collectionID)
	if err != nil {
		return fmt.Errorf("failed to retrieve collection: %v", err)
	}

	// Verify that the user is the owner of the collection
	if collection.OwnerID != userID {
		return errors.New("unauthorized: only the collection owner can remove posts")
	}

	// Remove the post from the collection
	err = l.CommunityEngagementLedger.RemovePostFromCollection(collectionID, postID)
	if err != nil {
		return fmt.Errorf("failed to remove post from collection: %v", err)
	}

	return nil
}


// listAllCollections retrieves all collections owned by a user
func listAllCollections(userID string) ([]ledger.Collection, error) {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Retrieve collections owned by the user
	collections, err := l.CommunityEngagementLedger.ListUserCollections(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list collections: %v", err)
	}

	return collections, nil
}

// deleteCollection deletes a collection, ensuring only the owner can perform the deletion
func deleteCollection(userID, collectionID string) error {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Retrieve the collection
	collection, err := l.CommunityEngagementLedger.FetchCollection(collectionID)
	if err != nil {
		return fmt.Errorf("failed to retrieve collection: %v", err)
	}

	// Verify that the user is the owner of the collection
	if collection.OwnerID != userID {
		return errors.New("unauthorized: only the collection owner can delete it")
	}

	// Delete the collection
	err = l.CommunityEngagementLedger.DeleteCollection(collectionID)
	if err != nil {
		return fmt.Errorf("failed to delete collection: %v", err)
	}

	return nil
}

// postReaction adds a reaction to a post
func postReaction(userID, postID, reactionType string) error {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Generate a unique ID for the reaction
	reactionID := generateUniqueID(userID + postID + reactionType)

	// Create the reaction struct
	reaction := ledger.Reaction{
		ID:        reactionID,
		PostID:    postID,
		UserID:    userID,
		Type:      reactionType,
		CreatedAt: time.Now(),
	}

	// Record the reaction in the ledger
	if err := l.CommunityEngagementLedger.RecordReaction(reaction); err != nil {
		return fmt.Errorf("failed to record reaction: %v", err)
	}

	return nil
}

// removePostReaction removes a user's reaction from a post
func removePostReaction(userID, postID string) error {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Retrieve the reaction
	reaction, err := l.CommunityEngagementLedger.FetchReactionByUserAndPost(userID, postID)
	if err != nil {
		return fmt.Errorf("failed to retrieve reaction: %v", err)
	}

	// Delete the reaction by its ID
	if err := l.CommunityEngagementLedger.DeleteReaction(reaction.ID); err != nil {
		return fmt.Errorf("failed to delete reaction: %v", err)
	}

	return nil
}

// reactToReply adds a reaction to a reply within a post
func ReactToReply(userID, replyID, reactionType string) error {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Generate a unique ID for the reaction
	reactionID := generateUniqueID(userID + replyID + reactionType)

	// Create the reaction struct
	reaction := ledger.Reaction{
		ID:        reactionID,
		ReplyID:   replyID,
		UserID:    userID,
		Type:      reactionType,
		CreatedAt: time.Now(),
	}

	// Record the reaction in the ledger
	if err := l.CommunityEngagementLedger.RecordReaction(reaction); err != nil {
		return fmt.Errorf("failed to record reaction to reply: %v", err)
	}

	return nil
}

// removeReplyReaction removes a reaction from a reply
func RemoveReplyReaction(userID, replyID string) error {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Retrieve the reaction by user and reply ID
	reaction, err := l.CommunityEngagementLedger.FetchReactionByUserAndReply(userID, replyID)
	if err != nil {
		return fmt.Errorf("failed to retrieve reaction: %v", err)
	}

	// Delete the reaction by its ID
	if err := l.CommunityEngagementLedger.DeleteReaction(reaction.ID); err != nil {
		return fmt.Errorf("failed to delete reaction: %v", err)
	}

	return nil
}

// Utility function to generate unique IDs
func generateUniqueID(data string) string {
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash[:])
}

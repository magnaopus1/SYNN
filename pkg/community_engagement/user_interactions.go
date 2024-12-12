package community_engagement

import (
	"fmt"
	"sync"
	"synnergy_network/pkg/ledger"
	"time"
)

// Mutex for user interactions
var userInteractionMutex sync.Mutex

// PrivateMessage struct to store information about private messages
type PrivateMessage struct {
	ID        string    // Unique identifier for the message
	SenderID  string    // ID of the user sending the message
	ReceiverID string   // ID of the user receiving the message
	Content   string    // Content of the message
	Timestamp time.Time // Time when the message was sent
}

// UserProfile struct to store user profile information
type UserProfile struct {
	UserID   string // Unique identifier for the user
	Username string // Username of the user
	Bio      string // Short bio or description
	Keywords []string // List of keywords for searchability
}

// UserBlock struct to store information about blocks between users
type UserBlock struct {
	RequesterID string // ID of the user who set the block
	TargetUserID string // ID of the user who is blocked
}

// followUser allows one user to follow another
func FollowUser(followerID, followeeID string) error {
	userInteractionMutex.Lock()
	defer userInteractionMutex.Unlock()

	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Record the follow relationship
	if err := l.CommunityEngagementLedger.RecordFollow(followerID, followeeID); err != nil {
		return fmt.Errorf("failed to follow user: %v", err)
	}
	return nil
}

// unfollowUser allows a user to unfollow another
func UnfollowUser(followerID, followeeID string) error {
	userInteractionMutex.Lock()
	defer userInteractionMutex.Unlock()

	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Remove the follow relationship
	if err := l.CommunityEngagementLedger.RemoveFollow(followerID, followeeID); err != nil {
		return fmt.Errorf("failed to unfollow user: %v", err)
	}
	return nil
}

// sendPrivateMessage allows a user to send a private message to another user
func SendPrivateMessage(senderID, receiverID, messageContent string) (string, error) {
	// Generate a unique ID for the message
	messageID := generateUniqueID(senderID + receiverID + messageContent + time.Now().String())

	// Create the private message struct
	privateMessage := ledger.PrivateMessage{
		ID:         messageID,
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    messageContent,
		Timestamp:  time.Now(),
	}

	userInteractionMutex.Lock()
	defer userInteractionMutex.Unlock()

	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Record the private message
	if err := l.CommunityEngagementLedger.RecordPrivateMessage(privateMessage); err != nil {
		return "", fmt.Errorf("failed to send private message: %v", err)
	}
	return messageID, nil
}


// readPrivateMessage retrieves a specific private message
func ReadPrivateMessage(messageID, receiverID string) (ledger.PrivateMessage, error) {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Fetch the private message
	message, err := l.CommunityEngagementLedger.FetchPrivateMessage(messageID, receiverID)
	if err != nil {
		return ledger.PrivateMessage{}, fmt.Errorf("message not found or unauthorized access: %v", err)
	}
	return message, nil
}

// blockUser prevents one user from interacting with another
func BlockUser(requesterID, targetUserID string) error {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Record the block relationship
	if err := l.CommunityEngagementLedger.RecordBlockUser(requesterID, targetUserID); err != nil {
		return fmt.Errorf("failed to block user: %v", err)
	}
	return nil
}

// unblockUser removes a block previously set by a user
func UnblockUser(requesterID, targetUserID string) error {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Remove the block relationship
	if err := l.CommunityEngagementLedger.RemoveBlockUser(requesterID, targetUserID); err != nil {
		return fmt.Errorf("failed to unblock user: %v", err)
	}
	return nil
}

// searchUser finds users based on search criteria such as username or keywords
func SearchUser(query string) ([]ledger.UserProfile, error) {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Search for matching user profiles
	users, err := l.CommunityEngagementLedger.SearchUserProfiles(query)
	if err != nil {
		return nil, fmt.Errorf("failed to search for users: %v", err)
	}
	return users, nil
}

// viewUserProfile retrieves the profile information of a specified user
func ViewUserProfile(userID string) (ledger.UserProfile, error) {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Fetch the user profile
	profile, err := l.CommunityEngagementLedger.FetchUserProfile(userID)
	if err != nil {
		return ledger.UserProfile{}, fmt.Errorf("user profile not found: %v", err)
	}
	return profile, nil
}

// editUserProfile allows a user to update their profile details
func EditUserProfile(userID string, profileUpdates ledger.UserProfile) error {
	userInteractionMutex.Lock()
	defer userInteractionMutex.Unlock()

	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Update the user profile
	if err := l.CommunityEngagementLedger.UpdateUserProfile(userID, profileUpdates); err != nil {
		return fmt.Errorf("failed to edit user profile: %v", err)
	}
	return nil
}

// listUserFollowers returns a list of users following a specified user
func ListUserFollowers(userID string) ([]string, error) {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Fetch followers of the user
	followers, err := l.CommunityEngagementLedger.FetchUserFollowers(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list user followers: %v", err)
	}
	return followers, nil
}

// listUserFollowing returns a list of users that a specified user is following
func ListUserFollowing(userID string) ([]string, error) {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Fetch users the user is following
	following, err := l.CommunityEngagementLedger.FetchUserFollowing(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list users followed: %v", err)
	}
	return following, nil
}

// muteUser prevents notifications from a specific user
func MuteUser(requesterID, targetUserID string) error {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Record the mute relationship
	if err := l.CommunityEngagementLedger.RecordMuteUser(requesterID, targetUserID); err != nil {
		return fmt.Errorf("failed to mute user: %v", err)
	}
	return nil
}

// unmuteUser allows notifications from a specific user again
func UnmuteUser(requesterID, targetUserID string) error {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Remove the mute relationship
	if err := l.CommunityEngagementLedger.RemoveMuteUser(requesterID, targetUserID); err != nil {
		return fmt.Errorf("failed to unmute user: %v", err)
	}
	return nil
}

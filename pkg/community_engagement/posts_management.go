package community_engagement

import (
	"fmt"
	"sync"
	"synnergy_network/pkg/ledger"
	"time"
)

// UserFavorites struct to track favorite posts for each user
type UserFavorites struct {
	UserID   string   // Unique identifier for the user
	PostIDs  []string // List of post IDs marked as favorite by the user
}

// PostReport struct to store information about reported posts
type PostReport struct {
	PostID     string    // ID of the reported post
	ReporterID string    // ID of the user who reported the post
	Reason     string    // Reason for reporting
	Timestamp  time.Time // Time when the report was made
}

// Mutex for post updates
var postMutex sync.Mutex

// createPost initializes a new post and records it in the ledger
func CreatePost(authorID, content string, tags []string) (string, error) {
	// Generate a unique ID for the post
	postID := generateUniqueID(authorID + content + time.Now().String())

	// Create the post struct
	post := ledger.Post{
		ID:        postID,
		AuthorID:  authorID,
		Content:   content,
		Tags:      tags,
		Timestamp: time.Now(),
		Upvotes:   0,
		Downvotes: 0,
	}

	// Record the post in the ledger
	l := &ledger.Ledger{}
	if err := l.CommunityEngagementLedger.RecordPost(post); err != nil {
		return "", fmt.Errorf("failed to record post: %v", err)
	}
	return postID, nil
}

// createReply adds a reply to an existing post
func CreateReply(postID, authorID, content string) (string, error) {
	// Generate a unique ID for the reply
	replyID := generateUniqueID(authorID + content + time.Now().String())

	// Create the reply struct
	reply := ledger.Reply{
		ID:        replyID,
		PostID:    postID,
		AuthorID:  authorID,
		Content:   content,
		Timestamp: time.Now(),
	}

	// Record the reply in the ledger
	l := &ledger.Ledger{}
	if err := l.CommunityEngagementLedger.RecordReply(reply); err != nil {
		return "", fmt.Errorf("failed to record reply: %v", err)
	}
	return replyID, nil
}

// queryPost retrieves a specific post by ID
func QueryPost(postID string) (ledger.Post, error) {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Retrieve the post
	post, err := l.CommunityEngagementLedger.FetchPost(postID)
	if err != nil {
		return ledger.Post{}, fmt.Errorf("post not found: %v", err)
	}
	return post, nil
}

// listAllPosts retrieves a list of all posts
func ListAllPosts() ([]ledger.Post, error) {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Retrieve all posts
	posts, err := l.CommunityEngagementLedger.ListAllPosts()
	if err != nil {
		return nil, fmt.Errorf("failed to list all posts: %v", err)
	}
	return posts, nil
}

// searchPosts retrieves posts containing specific tags or keywords
func SearchPosts(tags []string, keywords []string) ([]ledger.Post, error) {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Search for matching posts
	posts, err := l.CommunityEngagementLedger.SearchPosts(tags, keywords)
	if err != nil {
		return nil, fmt.Errorf("failed to search posts: %v", err)
	}
	return posts, nil
}

// upvotePost increases the upvote count for a post
func UpvotePost(postID string) error {
	// Lock the postMutex to prevent concurrent modifications
	postMutex.Lock()
	defer postMutex.Unlock()

	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Retrieve the post
	post, err := l.CommunityEngagementLedger.FetchPost(postID)
	if err != nil {
		return fmt.Errorf("post not found: %v", err)
	}

	// Increment the upvote count
	post.Upvotes++
	if err := l.CommunityEngagementLedger.UpdatePostVotes(post); err != nil {
		return fmt.Errorf("failed to upvote post: %v", err)
	}

	return nil
}

// downvotePost increases the downvote count for a post
func DownvotePost(postID string) error {
	// Lock the postMutex to prevent concurrent modifications
	postMutex.Lock()
	defer postMutex.Unlock()

	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Retrieve the post
	post, err := l.CommunityEngagementLedger.FetchPost(postID)
	if err != nil {
		return fmt.Errorf("post not found: %v", err)
	}

	// Increment the downvote count
	post.Downvotes++
	if err := l.CommunityEngagementLedger.UpdatePostVotes(post); err != nil {
		return fmt.Errorf("failed to downvote post: %v", err)
	}

	return nil
}

// markPostAsFavorite adds a post to the user's favorites list
func MarkPostAsFavorite(userID, postID string) error {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Add the post to the user's favorites
	if err := l.CommunityEngagementLedger.MarkFavoritePost(userID, postID); err != nil {
		return fmt.Errorf("failed to mark post as favorite: %v", err)
	}
	return nil
}

// removeFavoritePost removes a post from the user's favorites list
func RemoveFavoritePost(userID, postID string) error {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Remove the post from the user's favorites
	if err := l.CommunityEngagementLedger.RemoveFavoritePost(userID, postID); err != nil {
		return fmt.Errorf("failed to remove post from favorites: %v", err)
	}
	return nil
}

// editPost allows an author to edit their post content
func EditPost(postID, newContent string) error {
	// Lock the postMutex to prevent concurrent modifications
	postMutex.Lock()
	defer postMutex.Unlock()

	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Retrieve the post
	post, err := l.CommunityEngagementLedger.FetchPost(postID)
	if err != nil {
		return fmt.Errorf("post not found: %v", err)
	}

	// Update the post content
	post.Content = newContent
	if err := l.CommunityEngagementLedger.UpdatePostContent(post); err != nil {
		return fmt.Errorf("failed to edit post: %v", err)
	}

	return nil
}

// deletePost removes a post and its associated replies
func DeletePost(postID string) error {
	// Lock the postMutex to prevent concurrent modifications
	postMutex.Lock()
	defer postMutex.Unlock()

	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Delete the post and its associated replies
	if err := l.CommunityEngagementLedger.DeletePost(postID); err != nil {
		return fmt.Errorf("failed to delete post: %v", err)
	}
	return nil
}

// deleteReply removes a reply from a post
func DeleteReply(replyID string) error {
	// Lock the postMutex to prevent concurrent modifications
	postMutex.Lock()
	defer postMutex.Unlock()

	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Delete the reply
	if err := l.CommunityEngagementLedger.DeleteReply(replyID); err != nil {
		return fmt.Errorf("failed to delete reply: %v", err)
	}
	return nil
}

// reportPost allows users to report a post for moderation
func ReportPost(postID, reporterID, reason string) error {
	// Create a new report
	report := ledger.PostReport{
		PostID:     postID,
		ReporterID: reporterID,
		Reason:     reason,
		Timestamp:  time.Now(),
	}

	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Record the post report
	if err := l.CommunityEngagementLedger.RecordPostReport(report); err != nil {
		return fmt.Errorf("failed to report post: %v", err)
	}
	return nil
}

// pinPost pins a post for prominent display
func PinPost(adminID, postID string) error {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Mark the post as pinned
	if err := l.CommunityEngagementLedger.PinPost(adminID, postID); err != nil {
		return fmt.Errorf("failed to pin post: %v", err)
	}
	return nil
}

// unpinPost unpins a post, returning it to normal display status
func UnpinPost(adminID, postID string) error {
	// Create a new instance of the ledger
	l := &ledger.Ledger{}

	// Mark the post as unpinned
	if err := l.CommunityEngagementLedger.UnpinPost(adminID, postID); err != nil {
		return fmt.Errorf("failed to unpin post: %v", err)
	}
	return nil
}


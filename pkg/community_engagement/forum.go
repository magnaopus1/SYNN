package community_engagement

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// ForumPost represents a post within the community forum.
type ForumPost struct {
	ID        string    // Unique ID of the post
	Author    string    // The author of the post
	Content   string    // The content of the post
	Timestamp time.Time // Time the post was created
	Replies   []Reply   // Replies to the post
	Hash      string    // Hash to ensure data integrity
}

// Reply represents a reply to a forum post.
type Reply struct {
	ID        string    // Unique ID of the reply
	Author    string    // Author of the reply
	Content   string    // Content of the reply
	Timestamp time.Time // Time the reply was created
	Hash      string    // Hash to ensure reply integrity
}

// ForumManager manages the community forum, posts, and replies.
type ForumManager struct {
	Posts          map[string]*ForumPost // All forum posts
	LedgerInstance *ledger.Ledger        // Reference to the ledger for audit and storage
	mutex          sync.Mutex            // Mutex for thread-safe operations
}

// NewForumManager initializes a new ForumManager with an empty post list and ledger integration.
func NewForumManager(ledgerInstance *ledger.Ledger) *ForumManager {
    return &ForumManager{
        Posts:          make(map[string]*ForumPost),
        LedgerInstance: ledgerInstance,
    }
}

// CreatePost allows a user to create a new forum post.
func (fm *ForumManager) CreatePost(author, content string) (*ForumPost, error) {
    fm.mutex.Lock()
    defer fm.mutex.Unlock()

    postID := generatePostID(author, content)
    post := &ForumPost{
        ID:        postID,
        Author:    author,
        Content:   content,
        Timestamp: time.Now(),
        Replies:   []Reply{},
        Hash:      calculatePostHash(author, content, time.Now()),
    }

    fm.Posts[post.ID] = post
    fmt.Printf("New post created by %s: %s\n", author, post.ID)

    // Encrypt and store the post on the ledger
    encryptedPost, err := common.EncryptPost(post, common.EncryptionKey) // Use encryption.EncryptionKey
    if err != nil {
        return nil, fmt.Errorf("failed to encrypt post: %v", err)
    }

    // Convert encryptedPost (which is []byte) to a base64 encoded string
    encryptedPostStr := base64.StdEncoding.EncodeToString(encryptedPost)

    // Handle both return values from RecordForumPost
    successMsg, err := fm.LedgerInstance.CommunityEngagementLedger.RecordForumPost(post.ID, encryptedPostStr) // Pass string instead of []byte
    if err != nil {
        return nil, fmt.Errorf("failed to store post in ledger: %v", err)
    }

    fmt.Println(successMsg) // Optionally log the success message
    return post, nil
}


// calculateReplyHash generates a hash for a reply based on the postID, author, content, and timestamp.
func calculateReplyHash(postID, author, content string, timestamp time.Time) string {
	data := fmt.Sprintf("%s:%s:%s:%s", postID, author, content, timestamp.Format(time.RFC3339))
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// CreateReply allows a user to reply to an existing forum post.
func (fm *ForumManager) CreateReply(postID, author, content string) (*Reply, error) {
    fm.mutex.Lock()
    defer fm.mutex.Unlock()

    reply := &Reply{
        Author:    author,
        Content:   content,
        Timestamp: time.Now(),
        Hash:      calculateReplyHash(postID, author, content, time.Now()),
    }

    post, exists := fm.Posts[postID]
    if !exists {
        return nil, fmt.Errorf("post %s not found", postID)
    }

    post.Replies = append(post.Replies, *reply)
    fmt.Printf("New reply created by %s for post %s\n", author, postID)

    // Encrypt and store the reply on the ledger
    encryptedReply, err := common.EncryptPost(reply, common.EncryptionKey) // Use encryption.EncryptionKey
    if err != nil {
        return nil, fmt.Errorf("failed to encrypt reply: %v", err)
    }

    // Convert encryptedReply (which is []byte) to a base64 encoded string
    encryptedReplyStr := base64.StdEncoding.EncodeToString(encryptedReply)

    // Handle both return values from RecordForumReply
    successMsg, err := fm.LedgerInstance.CommunityEngagementLedger.RecordForumReply(postID, author, encryptedReplyStr) // Pass string
    if err != nil {
        return nil, fmt.Errorf("failed to store reply in ledger: %v", err)
    }

    fmt.Println(successMsg) // Optionally log the success message
    return reply, nil
}


// QueryPost returns the details of a forum post and its replies.
func (fm *ForumManager) QueryPost(postID string) (*ForumPost, error) {
    fm.mutex.Lock()
    defer fm.mutex.Unlock()

    post, exists := fm.Posts[postID]
    if !exists {
        return nil, fmt.Errorf("post not found: %s", postID)
    }

    // Encrypt the post details (if required) and return the post
    encryptedPost, err := common.EncryptPost(post, common.EncryptionKey) // Corrected to encryption.EncryptionKey
    if err != nil {
        return nil, fmt.Errorf("failed to encrypt post: %v", err)
    }

    // Optionally convert the encrypted post to a base64 string if you want to store or log it as a string
    encryptedPostStr := base64.StdEncoding.EncodeToString(encryptedPost)
    fmt.Printf("Encrypted post: %s\n", encryptedPostStr)

    // Return the unencrypted post, as the function signature expects a *ForumPost
    return post, nil
}

// ListAllPosts returns a list of all forum posts
func (fm *ForumManager) ListAllPosts() ([]*ForumPost, error) {
    fm.mutex.Lock()
    defer fm.mutex.Unlock()

    var postList []*ForumPost
    for _, post := range fm.Posts {
        postList = append(postList, post)
    }

    if len(postList) == 0 {
        return nil, fmt.Errorf("no posts found")
    }

    return postList, nil
}

// SearchPosts allows users to search posts by content or author
func (fm *ForumManager) SearchPosts(query string) ([]*ForumPost, error) {
    fm.mutex.Lock()
    defer fm.mutex.Unlock()

    var result []*ForumPost
    for _, post := range fm.Posts {
        if strings.Contains(post.Content, query) || strings.Contains(post.Author, query) {
            result = append(result, post)
        }
    }

    if len(result) == 0 {
        return nil, fmt.Errorf("no posts match the query: %s", query)
    }

    return result, nil
}

// calculatePostHash calculates the hash of a post or reply to ensure integrity.
func calculatePostHash(author, content string, timestamp time.Time) string {
    hashInput := fmt.Sprintf("%s%s%d", author, content, timestamp.UnixNano())
    hash := sha256.New()
    hash.Write([]byte(hashInput))
    return hex.EncodeToString(hash.Sum(nil))
}

// generatePostID generates a unique ID for a forum post.
func generatePostID(author, content string) string {
    hashInput := fmt.Sprintf("%s%s%d", author, content, time.Now().UnixNano())
    hash := sha256.New()
    hash.Write([]byte(hashInput))
    return hex.EncodeToString(hash.Sum(nil))
}

// generateReplyID generates a unique ID for a reply to a post.
func generateReplyID(author, content string) string {
    hashInput := fmt.Sprintf("%s%s%d", author, content, time.Now().UnixNano())
    hash := sha256.New()
    hash.Write([]byte(hashInput))
    return hex.EncodeToString(hash.Sum(nil))
}

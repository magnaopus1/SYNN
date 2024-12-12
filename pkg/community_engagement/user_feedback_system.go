package community_engagement

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)


// FeedbackSystem manages the feedback from users.
type FeedbackSystem struct {
	Feedbacks       map[string]*Feedback // Store feedback data
	LedgerInstance  *ledger.Ledger       // Ledger integration for secure feedback storage
	mutex           sync.Mutex           // Mutex for thread-safe operations
}

// NewFeedbackSystem initializes a new FeedbackSystem.
func NewFeedbackSystem(ledgerInstance *ledger.Ledger) *FeedbackSystem {
    return &FeedbackSystem{
        Feedbacks:      make(map[string]*Feedback),
        LedgerInstance: ledgerInstance,
    }
}

// SubmitFeedback allows users to submit their feedback.
func (fs *FeedbackSystem) SubmitFeedback(userID, content string) (*Feedback, error) {
    fs.mutex.Lock()
    defer fs.mutex.Unlock()

    // Generate a unique feedback ID
    feedbackID := generateFeedbackID(userID, content)
    newFeedback := &Feedback{
        ID:        feedbackID,
        UserID:    userID,
        Content:   content,
        Submitted: time.Now(),
        Resolved:  false,
        Likes:     0,
        Dislikes:  0,
        Comments:  []Comment{},
    }

    // Add the feedback to FeedbackSystem's map
    fs.Feedbacks[newFeedback.ID] = newFeedback
    fmt.Printf("Feedback from %s submitted.\n", userID)

    // Convert Comments to the ledger.Comment type
    var ledgerComments []ledger.Comment
    for _, comment := range newFeedback.Comments {
        ledgerComments = append(ledgerComments, ledger.Comment{
            UserID:    comment.UserID,
            Content:   comment.Content,
            Submitted: comment.Submitted,
        })
    }

    // Convert to a compatible ledger.Feedback type and store in ledger
    ledgerFeedback := ledger.Feedback{
        ID:        newFeedback.ID,
        UserID:    newFeedback.UserID,
        Content:   newFeedback.Content,
        Submitted: newFeedback.Submitted,
        Resolved:  newFeedback.Resolved,
        Likes:     newFeedback.Likes,
        Dislikes:  newFeedback.Dislikes,
        Comments:  ledgerComments,
    }

    // Store the feedback in the ledger
    err := fs.LedgerInstance.CommunityEngagementLedger.RecordFeedback(ledgerFeedback)
    if err != nil {
        return nil, fmt.Errorf("failed to store feedback in the ledger: %v", err)
    }

    return newFeedback, nil
}






// RetrieveFeedback returns feedback details for a specific feedback ID.
func (fs *FeedbackSystem) RetrieveFeedback(feedbackID string) (*Feedback, error) {
    fs.mutex.Lock()
    defer fs.mutex.Unlock()

    feedback, exists := fs.Feedbacks[feedbackID]
    if !exists {
        return nil, fmt.Errorf("feedback not found: %s", feedbackID)
    }

    // Encrypt the feedback details (for internal processing) but return the original feedback
    encryptedFeedback, err := common.EncryptFeedback(feedback, common.EncryptionKey)
    if err != nil {
        return nil, fmt.Errorf("failed to encrypt feedback: %v", err)
    }

    // Optionally log or store the encrypted feedback, but return the original feedback object
    fmt.Printf("Encrypted feedback: %x\n", encryptedFeedback)

    // Return the original feedback, not the encrypted version
    return feedback, nil
}


// ListAllFeedbacks lists all feedback provided by users.
func (fs *FeedbackSystem) ListAllFeedbacks() ([]*Feedback, error) {
    fs.mutex.Lock()
    defer fs.mutex.Unlock()

    var feedbackList []*Feedback
    for _, feedback := range fs.Feedbacks {
        feedbackList = append(feedbackList, feedback)
    }

    if len(feedbackList) == 0 {
        return nil, fmt.Errorf("no feedback found")
    }

    return feedbackList, nil
}

// SearchFeedback allows users to search feedback based on content or user.
func (fs *FeedbackSystem) SearchFeedback(query string) ([]*Feedback, error) {
    fs.mutex.Lock()
    defer fs.mutex.Unlock()

    var result []*Feedback
    for _, feedback := range fs.Feedbacks {
        if strings.Contains(feedback.Content, query) || strings.Contains(feedback.UserID, query) {
            result = append(result, feedback)
        }
    }

    if len(result) == 0 {
        return nil, fmt.Errorf("no feedback matches the query: %s", query)
    }

    return result, nil
}


// calculateFeedbackHash generates a hash for feedback to ensure its integrity.
func calculateFeedbackHash(userID, feedback string, rating int, timestamp time.Time) string {
    hashInput := fmt.Sprintf("%s%s%d%d", userID, feedback, rating, timestamp.UnixNano())
    hash := sha256.New()
    hash.Write([]byte(hashInput))
    return hex.EncodeToString(hash.Sum(nil))
}

// generateFeedbackID generates a unique ID for feedback.
func generateFeedbackID(userID, feedback string) string {
    hashInput := fmt.Sprintf("%s%s%d", userID, feedback, time.Now().UnixNano())
    hash := sha256.New()
    hash.Write([]byte(hashInput))
    return hex.EncodeToString(hash.Sum(nil))
}

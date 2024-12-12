package ledger

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// RecordForumPost records a new forum post in the community ledger.
func (l *CommunityEngagementLedger) RecordForumPost(author string, content string) (string, error) {
	l.Lock()
	defer l.Unlock()

	// Generate a unique post ID without any arguments
	postID := generateUniqueID()

	// Create a new forum post
	newPost := ForumPost{
		ID:        postID,
		Author:    author,
		Content:   content,
		Timestamp: time.Now(),
		Hash:      generateHash(content), // Optional: Generate hash to ensure integrity
	}

	// Add the post to the forum posts map
	l.ForumPosts[postID] = &newPost


	fmt.Printf("Forum post by %s recorded successfully with Post ID: %s\n", author, postID)
	return postID, nil
}

// RecordForumReply records a reply to an existing forum post in the community ledger.
func (l *CommunityEngagementLedger) RecordForumReply(postID string, author string, content string) (string, error) {
	l.Lock()
	defer l.Unlock()

	// Check if the post exists
	post, postExists := l.ForumPosts[postID]
	if !postExists {
		return "", fmt.Errorf("forum post with ID %s does not exist", postID)
	}

	// Generate a unique reply ID (without any argument)
	replyID := generateUniqueID()

	// Create a new reply
	newReply := Reply{
		ID:        replyID,
		Author:    author,
		Content:   content,
		Timestamp: time.Now(),
		Hash:      generateHash(content), // Optional: Generate hash for reply
	}

	// Add the reply to the post's reply list
	post.Replies = append(post.Replies, newReply)
	l.ForumReplies[postID] = append(l.ForumReplies[postID], newReply)

	fmt.Printf("Reply by %s to Post ID %s recorded successfully with Reply ID: %s\n", author, postID, replyID)
	return replyID, nil
}



// Helper function to generate a unique transaction ID (TxID)
func generateTransactionID() string {
	return fmt.Sprintf("tx_%d", time.Now().UnixNano())
}

// Helper function to generate a simple hash from content (for data integrity)
func generateHash(content string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(content)))
}



func (l *CommunityEngagementLedger) RecordCollection(collection Collection) error {
	// Check if collection ID already exists to avoid duplicates
	if _, exists := l.Collections[collection.ID]; exists {
		return fmt.Errorf("collection with ID %s already exists", collection.ID)
	}

	// Record the collection in the ledger
	l.Collections[collection.ID] = collection
	return nil
}


// FetchCollection retrieves a collection by its ID
func (l *CommunityEngagementLedger) FetchCollection(collectionID string) (Collection, error) {
	collection, exists := l.Collections[collectionID]
	if !exists {
		return Collection{}, fmt.Errorf("collection with ID %s not found", collectionID)
	}
	return collection, nil
}


// AddPostToCollection adds a post ID to a collection's list of post IDs
func (l *CommunityEngagementLedger) AddPostToCollection(collectionID, postID string) error {
	collection, exists := l.Collections[collectionID]
	if !exists {
		return fmt.Errorf("collection with ID %s not found", collectionID)
	}

	// Check if the post already exists in the collection
	for _, id := range collection.PostIDs {
		if id == postID {
			return fmt.Errorf("post with ID %s already in the collection", postID)
		}
	}

	// Add the post ID to the collection's PostIDs list
	collection.PostIDs = append(collection.PostIDs, postID)
	l.Collections[collectionID] = collection
	return nil
}


func (l *CommunityEngagementLedger) RemovePostFromCollection(collectionID, postID string) error {
	collection, exists := l.Collections[collectionID]
	if !exists {
		return fmt.Errorf("collection with ID %s not found", collectionID)
	}

	// Find the post ID in the collection's PostIDs list
	index := -1
	for i, id := range collection.PostIDs {
		if id == postID {
			index = i
			break
		}
	}

	// If postID not found, return an error
	if index == -1 {
		return fmt.Errorf("post with ID %s not found in collection", postID)
	}

	// Remove the post ID from the collection's PostIDs list
	collection.PostIDs = append(collection.PostIDs[:index], collection.PostIDs[index+1:]...)
	l.Collections[collectionID] = collection
	return nil
}


// listUserCollections retrieves all collections owned by the specified userID
func (l *CommunityEngagementLedger) ListUserCollections(userID string) ([]Collection, error) {
	var userCollections []Collection

	for _, collection := range l.Collections {
		if collection.OwnerID == userID {
			userCollections = append(userCollections, collection)
		}
	}

	if len(userCollections) == 0 {
		return nil, fmt.Errorf("no collections found for user with ID %s", userID)
	}

	return userCollections, nil
}


func (l *CommunityEngagementLedger) DeleteCollection(collectionID string) error {
	// Check if the collection exists
	if _, exists := l.Collections[collectionID]; !exists {
		return fmt.Errorf("collection with ID %s not found", collectionID)
	}

	// Delete the collection from the ledger
	delete(l.Collections, collectionID)
	return nil
}


// FetchReactionByUserAndPost retrieves a reaction by a user to a specific post
func (l *CommunityEngagementLedger) FetchReactionByUserAndPost(userID, postID string) (Reaction, error) {
	for _, reaction := range l.Reactions {
		if reaction.UserID == userID && reaction.PostID == postID {
			return reaction, nil
		}
	}
	return Reaction{}, fmt.Errorf("no reaction found for user %s on post %s", userID, postID)
}


func (l *CommunityEngagementLedger) DeleteReaction(reactionID string) error {
	// Check if the reaction exists
	if _, exists := l.Reactions[reactionID]; !exists {
		return fmt.Errorf("reaction with ID %s not found", reactionID)
	}

	// Delete the reaction from the ledger
	delete(l.Reactions, reactionID)
	return nil
}


func (l *CommunityEngagementLedger) RecordReaction(reaction Reaction) error {
	// Check for an existing reaction by the same user on the same reply
	for _, existingReaction := range l.Reactions {
		if existingReaction.ReplyID == reaction.ReplyID && existingReaction.UserID == reaction.UserID {
			return fmt.Errorf("reaction by user %s for reply %s already exists", reaction.UserID, reaction.ReplyID)
		}
	}

	// Record the new reaction
	l.Reactions[reaction.ID] = reaction
	return nil
}


func (l *CommunityEngagementLedger) FetchReactionByUserAndReply(userID, replyID string) (Reaction, error) {
	for _, reaction := range l.Reactions {
		if reaction.UserID == userID && reaction.ReplyID == replyID {
			return reaction, nil
		}
	}
	return Reaction{}, fmt.Errorf("no reaction found for user %s on reply %s", userID, replyID)
}

func (l *CommunityEngagementLedger) RecordEvent(event Event) error {
	// Check if an event with the same ID already exists
	if _, exists := l.CommunityEvents[event.ID]; exists {
		return fmt.Errorf("event with ID %s already exists", event.ID)
	}

	// Record the new event
	l.CommunityEvents[event.ID] = event
	return nil
}


func (l *CommunityEngagementLedger) FetchEvent(eventID string) (Event, error) {
	event, exists := l.CommunityEvents[eventID]
	if !exists {
		return Event{}, fmt.Errorf("event with ID %s not found", eventID)
	}
	return event, nil
}



func (l *CommunityEngagementLedger) UpdateEvent(event Event) error {
	// Check if the event exists
	if _, exists := l.CommunityEvents[event.ID]; !exists {
		return fmt.Errorf("event with ID %s not found", event.ID)
	}

	// Update the event in the ledger
	l.CommunityEvents[event.ID] = event
	return nil
}


func (l *CommunityEngagementLedger) ListAllEvents() ([]Event, error) {
	if len(l.CommunityEvents) == 0 {
		return nil, fmt.Errorf("no events found in the network")
	}

	// Collect all events in a slice
	var events []Event
	for _, event := range l.CommunityEvents {
		events = append(events, event)
	}

	return events, nil
}


func (l *CommunityEngagementLedger) RecordFeedback(feedback Feedback) error {
	// Check if feedback with the same ID already exists
	if _, exists := l.Feedbacks[feedback.ID]; exists {
		return fmt.Errorf("feedback with ID %s already exists", feedback.ID)
	}

	// Record the new feedback
	l.Feedbacks[feedback.ID] = feedback
	return nil
}


func (l *CommunityEngagementLedger) ListAllFeedbacks() ([]Feedback, error) {
	if len(l.Feedbacks) == 0 {
		return nil, fmt.Errorf("no feedback found in the network")
	}

	// Collect all feedback entries in a slice
	var feedbacks []Feedback
	for _, feedback := range l.Feedbacks {
		feedbacks = append(feedbacks, feedback)
	}

	return feedbacks, nil
}


func (l *CommunityEngagementLedger) SearchFeedback(query string) ([]Feedback, error) {
	var results []Feedback
	query = strings.ToLower(query)

	// Search through feedbacks for matching content or userID
	for _, feedback := range l.Feedbacks {
		if strings.Contains(strings.ToLower(feedback.Content), query) || strings.ToLower(feedback.UserID) == query {
			results = append(results, feedback)
		}
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no feedback found matching query: %s", query)
	}
	return results, nil
}


func (l *CommunityEngagementLedger) FetchFeedback(feedbackID string) (Feedback, error) {
	feedback, exists := l.Feedbacks[feedbackID]
	if !exists {
		return Feedback{}, fmt.Errorf("feedback with ID %s not found", feedbackID)
	}
	return feedback, nil
}


func (l *CommunityEngagementLedger) UpdateFeedback(feedback Feedback) error {
	// Check if the feedback exists
	if _, exists := l.Feedbacks[feedback.ID]; !exists {
		return fmt.Errorf("feedback with ID %s not found", feedback.ID)
	}

	// Update the feedback entry
	l.Feedbacks[feedback.ID] = feedback
	return nil
}


func (l *CommunityEngagementLedger) ListResolvedFeedback() ([]Feedback, error) {
	var resolvedFeedback []Feedback

	// Collect all resolved feedback entries
	for _, feedback := range l.Feedbacks {
		if feedback.Resolved {
			resolvedFeedback = append(resolvedFeedback, feedback)
		}
	}

	if len(resolvedFeedback) == 0 {
		return nil, fmt.Errorf("no resolved feedback found")
	}
	return resolvedFeedback, nil
}



func (l *CommunityEngagementLedger) RecordReport(report Report) error {
	// Check if the report already exists by its ID
	if _, exists := l.Reports[report.ID]; exists {
		return fmt.Errorf("report with ID %s already exists", report.ID)
	}

	// Record the new report
	l.Reports[report.ID] = report
	return nil
}

// generateUniqueID generates a globally unique identifier using UUID.
func generateUniqueID() string {
	return uuid.New().String()
}

func (l *CommunityEngagementLedger) UpdateUserStatus(userID, status string) error {
	user, exists := l.Users[userID]
	if !exists {
		return fmt.Errorf("user with ID %s not found", userID)
	}

	// Update user status
	user.Status = status
	l.Users[userID] = user
	return nil
}


func (l *CommunityEngagementLedger) LogModerationAction(adminID, targetID, action, reason string) error {
	logID := generateUniqueID(adminID + targetID + action)
	logEntry := ModerationLog{
		ID:         logID,
		AdminID:    adminID,
		TargetID:   targetID,
		Action:     action,
		Reason:     reason,
		DateLogged: time.Now(),
	}

	// Record the moderation action
	l.ModerationLogs[logID] = logEntry
	return nil
}

func (l *CommunityEngagementLedger) FlagContent(contentID, reason string) error {
	content, exists := l.Contents[contentID]
	if !exists {
		return fmt.Errorf("content with ID %s not found", contentID)
	}

	// Add a flag in the metadata to mark the content as inappropriate
	content.Metadata["flagged_reason"] = reason
	l.Contents[contentID] = content
	return nil
}

func (l *CommunityEngagementLedger) FetchModerationLog() ([]ModerationLog, error) {
	if len(l.ModerationLogs) == 0 {
		return nil, fmt.Errorf("no moderation log entries found")
	}

	// Collect all moderation log entries
	var logs []ModerationLog
	for _, log := range l.ModerationLogs {
		logs = append(logs, log)
	}
	return logs, nil
}


func (l *CommunityEngagementLedger) RecordBugReport(bugReport BugReport) error {
	// Check if a bug report with the same ID already exists
	if _, exists := l.BugReports[bugReport.ID]; exists {
		return fmt.Errorf("bug report with ID %s already exists", bugReport.ID)
	}

	// Record the bug report
	l.BugReports[bugReport.ID] = bugReport
	return nil
}


func (l *CommunityEngagementLedger) ListBugReports() ([]BugReport, error) {
	if len(l.BugReports) == 0 {
		return nil, fmt.Errorf("no bug reports found")
	}

	// Collect all bug reports
	var bugReports []BugReport
	for _, report := range l.BugReports {
		bugReports = append(bugReports, report)
	}
	return bugReports, nil
}

func (l *CommunityEngagementLedger) FetchBugReport(bugID string) (BugReport, error) {
	bugReport, exists := l.BugReports[bugID]
	if !exists {
		return BugReport{}, fmt.Errorf("bug report with ID %s not found", bugID)
	}
	return bugReport, nil
}


func (l *CommunityEngagementLedger) UpdateBugReport(bugReport BugReport) error {
	// Check if the bug report exists
	if _, exists := l.BugReports[bugReport.ID]; !exists {
		return fmt.Errorf("bug report with ID %s not found", bugReport.ID)
	}

	// Update the bug report entry
	l.BugReports[bugReport.ID] = bugReport
	return nil
}


func (l *CommunityEngagementLedger) AddFeedbackComment(feedbackID string, comment Comment) error {
	feedback, exists := l.Feedbacks[feedbackID]
	if !exists {
		return fmt.Errorf("feedback with ID %s not found", feedbackID)
	}

	// Append the comment to the feedback's Comments list
	feedback.Comments = append(feedback.Comments, comment)
	l.Feedbacks[feedbackID] = feedback
	return nil
}

func (l *CommunityEngagementLedger) UpdateFeedbackReaction(feedbackID, userID, reactionType string) error {
	feedback, exists := l.Feedbacks[feedbackID]
	if !exists {
		return fmt.Errorf("feedback with ID %s not found", feedbackID)
	}

	// Update likes or dislikes based on reaction type
	if reactionType == "like" {
		feedback.Likes++
	} else if reactionType == "dislike" {
		feedback.Dislikes++
	} else {
		return fmt.Errorf("invalid reaction type: %s", reactionType)
	}

	l.Feedbacks[feedbackID] = feedback
	return nil
}


// RecordPoll logs a new poll with question and options
func (l *CommunityEngagementLedger) RecordPoll(poll Poll) error {
	// Check if the poll already exists
	if _, exists := l.Polls[poll.ID]; exists {
		return fmt.Errorf("poll with ID %s already exists", poll.ID)
	}

	// Record the new poll
	l.Polls[poll.ID] = poll
	return nil
}

func (l *CommunityEngagementLedger) FetchPoll(pollID string) (Poll, error) {
	poll, exists := l.Polls[pollID]
	if !exists {
		return Poll{}, fmt.Errorf("poll with ID %s not found", pollID)
	}
	return poll, nil
}


// HasUserVoted checks if a user has voted in a poll
func (l *CommunityEngagementLedger) HasUserVoted(pollID, userID string) bool {
    l.Lock()
    defer l.Unlock()
    for _, vote := range l.Votes {
        if vote.PollID == pollID && vote.UserID == userID {
            return true
        }
    }
    return false
}

func (l *CommunityEngagementLedger) RecordVote(pollID, userID, option string) error {
	poll, exists := l.Polls[pollID]
	if !exists {
		return fmt.Errorf("poll with ID %s not found", pollID)
	}

	// Record the user's vote
	poll.Votes[option]++
	poll.VoterList[userID] = option
	l.Polls[pollID] = poll
	return nil
}


func (l *CommunityEngagementLedger) UpdatePollStatus(poll Poll) error {
	// Check if the poll exists
	if _, exists := l.Polls[poll.ID]; !exists {
		return fmt.Errorf("poll with ID %s not found", poll.ID)
	}

	// Update the poll status to closed
	poll.Open = false
	l.Polls[poll.ID] = poll
	return nil
}


func (l *CommunityEngagementLedger) ListAllPolls() ([]Poll, error) {
	if len(l.Polls) == 0 {
		return nil, fmt.Errorf("no polls found")
	}

	// Collect all polls in a slice
	var polls []Poll
	for _, poll := range l.Polls {
		polls = append(polls, poll)
	}
	return polls, nil
}


func (l *CommunityEngagementLedger) RecordPost(post Post) error {
	// Check if the post already exists
	if _, exists := l.ForumPosts[post.ID]; exists {
		return fmt.Errorf("post with ID %s already exists", post.ID)
	}

	// Record the new post
	l.ForumPosts[post.ID] = post
	return nil
}


func (l *CommunityEngagementLedger) RecordReply(reply Reply) error {
	// Check if the post exists for this reply
	if _, exists := l.ForumPosts[reply.PostID]; !exists {
		return fmt.Errorf("post with ID %s not found", reply.PostID)
	}

	// Record the reply
	l.ForumReplies[reply.ID] = reply
	return nil
}

func (l *CommunityEngagementLedger) FetchPost(postID string) (Post, error) {
	post, exists := l.ForumPosts[postID]
	if !exists {
		return Post{}, fmt.Errorf("post with ID %s not found", postID)
	}
	return post, nil
}


// generateReportID generates a unique report ID based on content ID
func generateReportID(contentID string) string {
    return "rep_" + contentID + "_" + time.Now().Format("20060102150405")
}



func (l *CommunityEngagementLedger) ListAllPosts() ([]Post, error) {
	if len(l.ForumPosts) == 0 {
		return nil, fmt.Errorf("no posts found")
	}

	// Collect all posts in a slice
	var posts []Post
	for _, post := range l.ForumPosts {
		posts = append(posts, post)
	}
	return posts, nil
}


// SearchPosts searches for posts containing a specific keyword
func (l *CommunityEngagementLedger) SearchPosts(tags []string, keywords []string) ([]Post, error) {
	var results []Post

	// Convert keywords to lowercase for case-insensitive matching
	keywordSet := make(map[string]bool)
	for _, keyword := range keywords {
		keywordSet[strings.ToLower(keyword)] = true
	}

	// Search for matching posts
	for _, post := range l.ForumPosts {
		tagMatch := len(tags) == 0
		keywordMatch := len(keywords) == 0

		// Check for matching tags
		if !tagMatch {
			for _, tag := range post.Tags {
				for _, searchTag := range tags {
					if strings.EqualFold(tag, searchTag) {
						tagMatch = true
						break
					}
				}
				if tagMatch {
					break
				}
			}
		}

		// Check for matching keywords in content
		if !keywordMatch {
			for word := range keywordSet {
				if strings.Contains(strings.ToLower(post.Content), word) {
					keywordMatch = true
					break
				}
			}
		}

		if tagMatch && keywordMatch {
			results = append(results, post)
		}
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no posts found with specified tags or keywords")
	}

	return results, nil
}


func (l *CommunityEngagementLedger) UpdatePostVotes(post Post) error {
	// Check if the post exists
	if _, exists := l.ForumPosts[post.ID]; !exists {
		return fmt.Errorf("post with ID %s not found", post.ID)
	}

	// Update the post entry
	l.ForumPosts[post.ID] = post
	return nil
}


func (l *CommunityEngagementLedger) MarkFavoritePost(userID, postID string) error {
	// Check if the post exists
	if _, exists := l.ForumPosts[postID]; !exists {
		return fmt.Errorf("post with ID %s not found", postID)
	}

	// Get or initialize the user's favorite posts
	favorites, exists := l.UserFavorites[userID]
	if !exists {
		favorites = UserFavorites{
			UserID:    userID,
			Favorites: []string{},
		}
	}

	// Check if the post is already marked as favorite
	for _, favoritePostID := range favorites.Favorites {
		if favoritePostID == postID {
			return fmt.Errorf("post already marked as favorite")
		}
	}

	// Add the post to the user's favorites
	favorites.Favorites = append(favorites.Favorites, postID)
	l.UserFavorites[userID] = favorites
	return nil
}


func (l *CommunityEngagementLedger) RemoveFavoritePost(userID, postID string) error {
	// Check if the user has favorites
	favorites, exists := l.UserFavorites[userID]
	if !exists {
		return fmt.Errorf("user with ID %s has no favorites", userID)
	}

	// Check if the post is in the user's favorites
	found := false
	for i, favoritePostID := range favorites.Favorites {
		if favoritePostID == postID {
			// Remove the post from favorites
			favorites.Favorites = append(favorites.Favorites[:i], favorites.Favorites[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("post with ID %s not found in user's favorites", postID)
	}

	// Update the user's favorites in the ledger
	l.UserFavorites[userID] = favorites
	return nil
}


func (l *CommunityEngagementLedger) UpdatePostContent(post Post) error {
	// Check if the post exists
	if _, exists := l.ForumPosts[post.ID]; !exists {
		return fmt.Errorf("post with ID %s not found", post.ID)
	}

	// Update the content of the post
	l.ForumPosts[post.ID] = post
	return nil
}


// DeletePost removes a post from the ledger
func (l *CommunityEngagementLedger) DeletePost(postID string) error {
	// Check if the post exists
	if _, exists := l.ForumPosts[postID]; !exists {
		return fmt.Errorf("post with ID %s not found", postID)
	}

	// Delete all replies associated with the post
	for replyID, reply := range l.ForumReplies {
		if reply.PostID == postID {
			delete(l.ForumReplies, replyID)
		}
	}

	// Delete the post itself
	delete(l.ForumPosts, postID)
	return nil
}


// DeleteReply deletes a reply by reply ID
func (l *CommunityEngagementLedger) DeleteReply(replyID string) error {
	// Check if the reply exists
	if _, exists := l.ForumReplies[replyID]; !exists {
		return fmt.Errorf("reply with ID %s not found", replyID)
	}

	// Delete the reply
	delete(l.ForumReplies, replyID)
	return nil
}


func (l *CommunityEngagementLedger) RecordPostReport(report PostReport) error {
	// Check if the post exists
	if _, exists := l.ForumPosts[report.PostID]; !exists {
		return fmt.Errorf("post with ID %s not found", report.PostID)
	}

	// Append the report to the list of reports for the post
	l.PostReports[report.PostID] = append(l.PostReports[report.PostID], report)
	return nil
}

func (l *CommunityEngagementLedger) PinPost(adminID, postID string) error {
	// Check if the post exists
	if _, exists := l.ForumPosts[postID]; !exists {
		return fmt.Errorf("post with ID %s not found", postID)
	}

	// Mark the post as pinned
	l.PinnedPosts[postID] = true
	return nil
}


func (l *CommunityEngagementLedger) UnpinPost(adminID, postID string) error {
	// Check if the post exists
	if _, exists := l.ForumPosts[postID]; !exists {
		return fmt.Errorf("post with ID %s not found", postID)
	}

	// Mark the post as unpinned
	l.PinnedPosts[postID] = false
	return nil
}

func (l *CommunityEngagementLedger) RecordFollow(followerID, followeeID string) error {
	// Check if the follow relationship already exists
	if l.Followings[followerID] == nil {
		l.Followings[followerID] = make(map[string]bool)
	}
	if l.Followings[followerID][followeeID] {
		return fmt.Errorf("user %s is already following user %s", followerID, followeeID)
	}

	// Add the follow relationship
	l.Followings[followerID][followeeID] = true
	return nil
}


// RemoveFollow removes a follow relationship between two users
func (l *CommunityEngagementLedger) RemoveFollow(followerID, followeeID string) error {
	// Check if the follow relationship exists
	if l.Followings[followerID] == nil || !l.Followings[followerID][followeeID] {
		return fmt.Errorf("user %s is not following user %s", followerID, followeeID)
	}

	// Remove the follow relationship
	delete(l.Followings[followerID], followeeID)
	return nil
}


// RecordPrivateMessage logs a private message between two users
func (l *CommunityEngagementLedger) RecordPrivateMessage(message PrivateMessage) error {
	// Add the message to the sender and receiver's message history
	l.PrivateMessages[message.ReceiverID] = append(l.PrivateMessages[message.ReceiverID], message)
	l.PrivateMessages[message.SenderID] = append(l.PrivateMessages[message.SenderID], message)
	return nil
}


// FetchPrivateMessage retrieves a private message by its ID
func (l *CommunityEngagementLedger) FetchPrivateMessage(messageID, receiverID string) (PrivateMessage, error) {
	if messages, exists := l.PrivateMessages[receiverID]; exists {
		if message, found := messages[messageID]; found {
			return message, nil
		}
	}
	return PrivateMessage{}, fmt.Errorf("message with ID %s not found for receiver %s", messageID, receiverID)
}



// RecordBlockUser logs a block action between two users
func (l *CommunityEngagementLedger) RecordBlockUser(requesterID, targetUserID string) error {
	if l.BlockedUsers[requesterID] == nil {
		l.BlockedUsers[requesterID] = make(map[string]bool)
	}
	l.BlockedUsers[requesterID][targetUserID] = true
	return nil
}


// RemoveBlockUser removes a block relationship between two users
func (l *CommunityEngagementLedger) RemoveBlockUser(requesterID, targetUserID string) error {
	if blockedUsers, exists := l.BlockedUsers[requesterID]; exists {
		if _, found := blockedUsers[targetUserID]; found {
			delete(blockedUsers, targetUserID)
			return nil
		}
	}
	return fmt.Errorf("user %s is not blocked by %s", targetUserID, requesterID)
}


// SearchUserProfiles searches for user profiles with specific keywords in their bio
func (l *CommunityEngagementLedger) SearchUserProfiles(query string) ([]UserProfile, error) {
	var results []UserProfile
	query = strings.ToLower(query)

	// Search for matching user profiles by username or keywords
	for _, profile := range l.UserProfiles {
		if strings.Contains(strings.ToLower(profile.Username), query) {
			results = append(results, profile)
		} else {
			for _, keyword := range profile.Keywords {
				if strings.Contains(strings.ToLower(keyword), query) {
					results = append(results, profile)
					break
				}
			}
		}
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no users found matching the query: %s", query)
	}

	return results, nil
}



// FetchUserProfile retrieves a user profile by user ID
func (l *CommunityEngagementLedger) FetchUserProfile(userID string) (UserProfile, error) {
	profile, exists := l.UserProfiles[userID]
	if !exists {
		return UserProfile{}, fmt.Errorf("user profile with ID %s not found", userID)
	}
	return profile, nil
}


// UpdateUserProfile updates a user's profile information
func (l *CommunityEngagementLedger) UpdateUserProfile(userID string, profileUpdates UserProfile) error {
	// Check if the user profile exists
	if _, exists := l.UserProfiles[userID]; !exists {
		return fmt.Errorf("user profile with ID %s not found", userID)
	}

	// Update the profile with new information
	profile := l.UserProfiles[userID]
	if profileUpdates.Username != "" {
		profile.Username = profileUpdates.Username
	}
	if profileUpdates.Bio != "" {
		profile.Bio = profileUpdates.Bio
	}
	if len(profileUpdates.Keywords) > 0 {
		profile.Keywords = profileUpdates.Keywords
	}

	// Save the updated profile back to the map
	l.UserProfiles[userID] = profile
	return nil
}


// FetchUserFollowers retrieves a list of followers for a specific user
func (l *CommunityEngagementLedger) FetchUserFollowers(userID string) ([]string, error) {
	followers, exists := l.Followers[userID]
	if !exists {
		return nil, fmt.Errorf("no followers found for user with ID %s", userID)
	}
	return followers, nil
}


// FetchUserFollowing retrieves a list of users a specific user is following
func (l *CommunityEngagementLedger) FetchUserFollowing(userID string) ([]string, error) {
	following, exists := l.Following[userID]
	if !exists {
		return nil, fmt.Errorf("user with ID %s is not following anyone", userID)
	}
	return following, nil
}

// RecordMuteUser logs a mute action between two users
func (l *CommunityEngagementLedger) RecordMuteUser(requesterID, targetUserID string) error {
	if l.MuteList[requesterID] == nil {
		l.MuteList[requesterID] = make(map[string]bool)
	}
	l.MuteList[requesterID][targetUserID] = true
	return nil
}


// RemoveMuteUser removes a mute relationship between two users
func (l *CommunityEngagementLedger) RemoveMuteUser(requesterID, targetUserID string) error {
	if mutedUsers, exists := l.MuteList[requesterID]; exists {
		if _, found := mutedUsers[targetUserID]; found {
			delete(mutedUsers, targetUserID)
			return nil
		}
	}
	return fmt.Errorf("user %s is not muted by %s", targetUserID, requesterID)
}


// Helper functions
func containsKeyword(content, keyword string) bool {
    return len(content) > 0 && len(keyword) > 0 && len(keyword) <= len(content) && string.Contains(content, keyword)
}

func removeFromSlice(slice []string, value string) []string {
    for i, v := range slice {
        if v == value {
            return append(slice[:i], slice[i+1:]...)
        }
    }
    return slice
}





















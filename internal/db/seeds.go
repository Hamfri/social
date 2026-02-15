package db

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand/v2"
	"social/internal/repository"
)

var usernames = []string{
	"jabari", "faraja", "amani", "zuri", "neema", "imani",
	"bahati", "lulu", "nuru", "rehema", "pendo", "jamila",
	"aisha", "baraka", "juma", "omari", "simba", "sadiki",
	"khamisi", "zuberi", "zawadi", "tumaini", "subira", "kamari",
	"hodari", "makini",
}

var randTitles = []string{
	"Hello World, But Make It a Blog",
	"Why I Finally Started Writing This Blog",
	" Things I Learned Building My First Side Project",
	"A Random Thought That Turned Into a Post",
	"The Day Everything Broke (And What I Did Next)",
	"How Not to Design a Database Schema",
	"Lessons Learned After 30 Days of Consistent Writing",
	"This Post Exists Only for Testing Purposes",
	"What I Wish I Knew Before Becoming a Developer",
	"Debugging at 2 AM: A Love–Hate Story",
	"Small Changes That Made My Code Better",
	"From Idea to Production: A Very Rough Journey",
	"Mistakes Were Made, But Here We Are",
	"The Problem Nobody Warned Me About",
	"Building Features No One Asked For",
	"It Worked on My Machine (Until It Didn’t)",
	"A Quick Brain Dump on Today’s Progress",
	"Refactoring Old Code Without Breaking Everything",
	"Why Simplicity Is Harder Than It Looks",
	"Testing in Production (Just Kidding… Mostly)",
}

var randContent = []string{
	"This is a simple placeholder blog post used for testing purposes. It contains generic text and no real meaning.",
	"In this post, we explore some thoughts about building software, learning from mistakes, and improving over time.",
	"Writing consistently is harder than it looks. This post exists to help test content rendering and storage.",
	"Today I learned that small design decisions can have a big impact later in a project’s lifecycle.",
	"This article talks about nothing in particular but is perfect for testing pagination and limits.",
	"Debugging late at night often leads to unexpected discoveries and even better solutions.",
	"This post documents a fictional journey from an idea to a deployed application.",
	"Sometimes you build features just to see if you can. This is one of those moments.",
	"Refactoring old code can be scary, but it’s usually worth the effort in the long run.",
	"This content is intentionally boring and predictable to make backend testing easier.",
	"Here we describe a problem that doesn’t exist and a solution that wasn’t needed.",
	"Every developer has written a post like this at least once during testing.",
	"This blog entry is part of a dummy dataset for development and QA environments.",
	"A short write-up about progress, blockers, and lessons learned along the way.",
	"This is where insightful commentary would normally go, but today it’s filler text.",
	"Testing edge cases often requires realistic-looking data like this blog post.",
	"This article simulates a real user-generated post for API testing.",
	"Some posts are meant to teach, others are just meant to exist—this is the latter.",
	"This content helps verify create, read, update, and delete functionality.",
	"Thanks for reading this completely fake but very useful test blog post.",
}

var randComments = []string{
	"Nice post!",
	"Great read.",
	"Very helpful, thanks!",
	"I learned something new.",
	"Interesting perspective.",
	"Well written.",
	"Thanks for sharing.",
	"Good explanation.",
	"Clear and concise.",
	"This makes sense.",
	"I agree with this.",
	"Helpful example.",
	"Good points.",
	"Enjoyed reading this.",
	"Straight to the point.",
	"Simple and useful.",
	"Looking forward to more.",
	"Nice work!",
	"Well explained.",
	"Helpful content.",
}

var randTags = []string{
	"golang",
	"backend",
	"programming",
	"webdev",
	"testing",
	"devlife",
}

func Seed(r repository.Repository, db *sql.DB) {
	ctx := context.Background()

	users := generateUsers(len(usernames))

	tx, _ := db.BeginTx(ctx, nil)

	for _, user := range users {
		if err := r.Users.Create(ctx, tx, user); err != nil {

			_ = tx.Rollback()

			fmt.Println("error creating user:", err)

			return
		}
	}

	_ = tx.Commit()

	posts := generatePosts(len(randContent), users)
	for _, post := range posts {
		if err := r.Posts.Create(ctx, post); err != nil {
			fmt.Println("error creating post:", err)
			return
		}
	}

	comments := generateComments(len(posts), users, posts)
	for _, comment := range comments {
		if err := r.Comments.Create(ctx, comment); err != nil {
			fmt.Println("error creating comments:", err)
			return
		}
	}

	userFollows := generateUserFollows(len(usernames))
	shuffle(userFollows)
	for _, userFollow := range userFollows[:rand.IntN(len(userFollows))] {
		if err := r.UserFollows.Follow(ctx, userFollow); err != nil {
			fmt.Println("error creating userfollows:", err)
		}
	}

	fmt.Println("seeding complete")
}

func generateUsers(num int) []*repository.User {
	users := make([]*repository.User, num)

	for i := 0; i < num; i++ {
		username := usernames[i%len(usernames)] // + fmt.Sprintf("%d", i)
		users[i] = &repository.User{
			Username: username,
			Email:    username + "@hadaa.com",
		}

		err := users[i].Password.Set("qwerty1234")
		if err != nil {
			fmt.Printf("error when hashing password for %s", users[i].Email)
			break
		}
	}

	return users
}

// permutation
// P(n, 2) = n * (n - 2)
// usernames = 26, n = 26
// P(26, 2) = 26 * 25 = 650
// (A follows B) != (B follows A)
// No self-follows A != A and no repetitions
func generateUserFollows(userCount int) []*repository.UserFollow {
	var userFollows = make([]*repository.UserFollow, 0, userCount*userCount-1) // length 0, cap = userCount*userCount-1

	for i := 1; i <= userCount; i++ {
		// TERRIBLE idea
		// Never take the address of a loop-scoped variable and store it
		// var userFollow repository.UserFollow
		for x := 1; x <= userCount; x++ {
			if i != x {
				// TERRIBLE idea
				// userFollow.FollowedID = int64(i)
				// userFollow.FollowerID = int64(x)
				// userFollows = append(userFollows, &userFollow)

				userFollows = append(userFollows, &repository.UserFollow{
					FollowedID: int64(i),
					FollowerID: int64(x),
				})
			}
		}
	}

	return userFollows
}

func generatePosts(num int, users []*repository.User) []*repository.Post {
	posts := make([]*repository.Post, num)
	for i := 0; i < num; i++ {
		user := users[rand.IntN(len(users))]

		posts[i] = &repository.Post{
			UserID:  user.ID,
			Title:   randTitles[rand.IntN(len(randTitles))],
			Content: randContent[rand.IntN(len(randContent))],
			Tags:    []string{randTags[rand.IntN(len(randTags))], randTags[rand.IntN(len(randTags))]},
		}
	}

	return posts
}

func generateComments(num int, users []*repository.User, posts []*repository.Post) []*repository.Comment {
	comments := make([]*repository.Comment, num)

	for i := 0; i < num; i++ {
		user := users[rand.IntN(len(users))]
		post := posts[rand.IntN(len(posts))]

		comments[i] = &repository.Comment{
			Comment: randComments[rand.IntN(len(randComments))],
			UserId:  user.ID,
			PostId:  post.ID,
		}
	}

	return comments
}

func shuffle[T any](s []T) {
	r := rand.New(rand.NewPCG(1, 2))
	r.Shuffle(len(s), func(i, j int) {
		s[i], s[j] = s[j], s[i]
	})
}
